#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: scripts/run-use-case-enrichment.sh --config path/to/use-case.json [options]

Runs company, company-goat, contact-goat and enrich for each lead in a use-case CSV.

Options:
  --config PATH      Required use-case JSON config.
  --out-dir PATH     Override output directory.
  --limit N          Process only the first N leads.
  --dry-run          Print commands without executing them.
  --setup-auth       Run scripts/setup-external-enrichment-auth.sh first.
  --help             Show this help.
EOF
}

config_path=""
out_dir_override=""
limit=""
dry_run=false
setup_auth=false

while [ "$#" -gt 0 ]; do
  case "$1" in
    --config)
      config_path="${2:-}"
      shift 2
      ;;
    --out-dir)
      out_dir_override="${2:-}"
      shift 2
      ;;
    --limit)
      limit="${2:-}"
      shift 2
      ;;
    --dry-run)
      dry_run=true
      shift
      ;;
    --setup-auth)
      setup_auth=true
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      printf '%s\n' "error: unknown argument: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

if [ -z "$config_path" ]; then
  printf '%s\n' 'error: --config is required' >&2
  usage >&2
  exit 2
fi

script_dir="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(CDPATH= cd -- "$script_dir/.." && pwd)"
config_abs="$(python3 - "$repo_root/$config_path" "$config_path" <<'PY'
import os, sys
candidate = sys.argv[1]
raw = sys.argv[2]
if os.path.isabs(raw):
    print(os.path.abspath(raw))
else:
    print(os.path.abspath(candidate))
PY
)"

if [ ! -f "$config_abs" ]; then
  printf '%s\n' "error: config not found: $config_abs" >&2
  exit 2
fi

readarray -t config_values < <(python3 - "$config_abs" <<'PY'
import json, os, sys
path = sys.argv[1]
with open(path, 'r', encoding='utf-8') as fh:
    data = json.load(fh)
base = os.path.dirname(path)
lead = data.get('lead_table_path', '')
out = data.get('output_dir', '')
name_field = (data.get('field_map') or {}).get('name', 'lead')
cnpj_field = (data.get('field_map') or {}).get('cnpj', 'CNPJ')
label = data.get('label') or data.get('name') or 'use-case'
if lead and not os.path.isabs(lead):
    lead = os.path.abspath(os.path.join(base, lead))
if out and not os.path.isabs(out):
    out = os.path.abspath(os.path.join(base, out))
print(lead)
print(out)
print(name_field)
print(cnpj_field)
print(label)
PY
)

lead_table="${config_values[0]}"
out_dir="${config_values[1]}"
name_field="${config_values[2]}"
cnpj_field="${config_values[3]}"
use_case_label="${config_values[4]}"

if [ -n "$out_dir_override" ]; then
  out_dir="$out_dir_override"
fi
if [ -z "$out_dir" ]; then
  out_dir="$(dirname "$config_abs")/outputs/enrichment"
fi

if [ ! -f "$lead_table" ]; then
  printf '%s\n' "error: lead table not found: $lead_table" >&2
  exit 2
fi

export PP_LEADS_USE_CASE_CONFIG="$config_abs"
mkdir -p "$out_dir"

if [ "$setup_auth" = true ]; then
  "$script_dir/setup-external-enrichment-auth.sh"
fi

cd "$repo_root"

if command -v leads-brasil-pp-cli >/dev/null 2>&1; then
  pp_leads_mode="bin"
elif [ -x "$repo_root/leads-brasil-pp-cli/leads-brasil-pp-cli" ]; then
  pp_leads_mode="repo-bin"
else
  pp_leads_mode="go"
fi

leads-brasil-pp-cli() {
  if [ "$pp_leads_mode" = "bin" ]; then
    command leads-brasil-pp-cli "$@"
    return
  fi
  if [ "$pp_leads_mode" = "repo-bin" ]; then
    "$repo_root/leads-brasil-pp-cli/leads-brasil-pp-cli" "$@"
    return
  fi
  (
    cd "$repo_root/leads-brasil-pp-cli"
    GO111MODULE=on go run ./cmd/leads-brasil-pp-cli "$@"
  )
}

local_api_pid=""
cleanup_local_api() {
  if [ -n "$local_api_pid" ] && kill -0 "$local_api_pid" >/dev/null 2>&1; then
    kill "$local_api_pid" >/dev/null 2>&1 || true
    wait "$local_api_pid" >/dev/null 2>&1 || true
  fi
}
trap cleanup_local_api EXIT

start_local_api() {
  if [ -x "$repo_root/server_bin" ]; then
    "$repo_root/server_bin" >/tmp/pp-leads-server.log 2>&1 &
  else
    GO111MODULE=on go run ./cmd/server >/tmp/pp-leads-server.log 2>&1 &
  fi
  local_api_pid="$!"
  sleep 3
}

doctor_field() {
  python3 - "$1" "$2" <<'PY'
import json, sys
try:
    report = json.loads(sys.argv[1])
except json.JSONDecodeError:
    print("")
    raise SystemExit(0)
value = report.get(sys.argv[2], '')
print(value if isinstance(value, str) else '')
PY
}

is_local_base_url() {
  case "$1" in
    http://127.0.0.1:*|http://localhost:*|http://0.0.0.0:*) return 0 ;;
    *) return 1 ;;
  esac
}

run_interactive_auth() {
  local token save_token
  if [ ! -t 0 ] || [ ! -r /dev/tty ]; then
    printf '%s\n' 'error: autenticação ausente. Rode interativamente ou exporte LEADS_BRASIL_BEARER_AUTH.' >&2
    exit 4
  fi
  printf '%s' 'Cole o token do LEADS_BRASIL_BEARER_AUTH: ' > /dev/tty
  IFS= read -r -s token < /dev/tty
  printf '\n' > /dev/tty
  if [ -z "$token" ]; then
    printf '%s\n' 'error: token vazio.' >&2
    exit 4
  fi
  export LEADS_BRASIL_BEARER_AUTH="$token"
  printf '%s' 'Salvar token em ~/.config/leads-brasil-pp-cli/config.toml? [S/n] ' > /dev/tty
  IFS= read -r save_token < /dev/tty
  case "${save_token:-S}" in
    s|S|sim|SIM|Sim|y|Y|yes|YES|Yes)
      leads-brasil-pp-cli auth set-token "$token"
      ;;
  esac
}

ensure_pp_leads_ready() {
  local doctor_json base_url auth_status api_status
  doctor_json="$(leads-brasil-pp-cli doctor --json || true)"
  base_url="$(doctor_field "$doctor_json" base_url)"
  auth_status="$(doctor_field "$doctor_json" auth)"
  api_status="$(doctor_field "$doctor_json" api)"

  if [ "$auth_status" = 'not configured' ]; then
    if is_local_base_url "$base_url"; then
      export LEADS_BRASIL_BEARER_AUTH="local-dev-token"
    else
      run_interactive_auth
    fi
  fi

  case "$api_status" in
    unreachable*)
      if is_local_base_url "$base_url"; then
        start_local_api
      fi
      ;;
  esac
}

ensure_pp_leads_ready

readarray -t leads < <(python3 - "$lead_table" "$name_field" "$cnpj_field" "$limit" <<'PY'
import csv, sys
path, name_field, cnpj_field, limit_raw = sys.argv[1:5]
limit = int(limit_raw) if limit_raw else None
with open(path, 'r', encoding='utf-8-sig', newline='') as fh:
    rows = csv.DictReader(fh)
    count = 0
    for row in rows:
        name = (row.get(name_field) or '').strip()
        cnpj = (row.get(cnpj_field) or '').strip()
        if not name or not cnpj:
            continue
        print(f"{cnpj}\t{name}")
        count += 1
        if limit is not None and count >= limit:
            break
PY
)

printf '%s\n' "Use case: $use_case_label"
printf '%s\n' "Lead table: $lead_table"
printf '%s\n' "Output dir: $out_dir"

run_or_print() {
  if [ "$dry_run" = true ]; then
    printf '%q ' "$@"
    printf '\n'
    return 0
  fi
  "$@"
}

only_digits() {
  python3 - "$1" <<'PY'
import re, sys
print(re.sub(r'\D+', '', sys.argv[1]))
PY
}

for item in "${leads[@]}"; do
  cnpj="${item%%$'\t'*}"
  name="${item#*$'\t'}"
  digits="$(only_digits "$cnpj")"
  run_or_print leads-brasil-pp-cli company "$cnpj" --agent --select name,cnpj,status,address,contacts,casadosdados,sources,use_case,lead_context --deliver "file:$out_dir/${digits}-company.json"
  run_or_print leads-brasil-pp-cli company-goat "$cnpj" --agent --deliver "file:$out_dir/${digits}-company-goat.json"
  run_or_print leads-brasil-pp-cli contact-goat "$name" --agent --deliver "file:$out_dir/${digits}-contact-goat.json"
  run_or_print leads-brasil-pp-cli enrich "$cnpj" --agent --deliver "file:$out_dir/${digits}-enrich.json"
done
