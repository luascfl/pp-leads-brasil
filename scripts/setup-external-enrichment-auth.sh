#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: scripts/setup-external-enrichment-auth.sh [--check]

Interactive helper for Casa dos Dados and external enrichers.

Options:
  --check   Only print doctor output, never prompt.
  --help    Show this help.
EOF
}

check_only=false
case "${1:-}" in
  --check)
    check_only=true
    ;;
  --help|-h)
    usage
    exit 0
    ;;
  "")
    ;;
  *)
    printf '%s\n' "error: unknown argument: $1" >&2
    usage >&2
    exit 2
    ;;
esac

need_bin() {
  if ! command -v "$1" >/dev/null 2>&1; then
    printf '%s\n' "error: comando ausente: $1" >&2
    exit 127
  fi
}

need_bin contact-goat-pp-cli

run_doctors() {
  printf '%s\n' '== contact-goat doctor =='
  contact-goat-pp-cli doctor || true
  if command -v company-goat-pp-cli >/dev/null 2>&1; then
    printf '\n%s\n' '== company-goat doctor =='
    company-goat-pp-cli doctor || true
  fi
  if command -v scrape-creators-pp-cli >/dev/null 2>&1; then
    printf '\n%s\n' '== scrape-creators doctor =='
    scrape-creators-pp-cli doctor || true
  fi
}

prompt_secret_env() {
  local var_name="$1"
  local prompt="$2"
  local value
  if [ -n "${!var_name:-}" ]; then
    return 0
  fi
  printf '%s' "$prompt" > /dev/tty
  IFS= read -r -s value < /dev/tty
  printf '\n' > /dev/tty
  if [ -n "$value" ]; then
    export "$var_name=$value"
    printf '%s\n' "$var_name exportada para esta sessão."
  fi
}

run_doctors

if [ "$check_only" = true ]; then
  exit 0
fi

if [ ! -t 0 ] || [ ! -r /dev/tty ]; then
  printf '%s\n' 'error: modo interativo indisponível; use --check ou rode em um terminal.' >&2
  exit 4
fi

if [ -z "${CASA_DADOS_API_KEY:-${PP_LEADS_CASA_DADOS_API_KEY:-}}" ]; then
  printf '%s' 'Cole CASA_DADOS_API_KEY, ou deixe vazio para pular: ' > /dev/tty
  prompt_secret_env CASA_DADOS_API_KEY ''
fi

printf '%s' 'Rodar login do Happenstance pelo Chrome agora? [s/N] ' > /dev/tty
IFS= read -r happenstance_login < /dev/tty
case "${happenstance_login:-N}" in
  s|S|sim|SIM|Sim|y|Y|yes|YES|Yes)
    contact-goat-pp-cli auth login --chrome --service happenstance
    ;;
esac

printf '%s' 'Rodar login do LinkedIn MCP agora? [s/N] ' > /dev/tty
IFS= read -r linkedin_login < /dev/tty
case "${linkedin_login:-N}" in
  s|S|sim|SIM|Sim|y|Y|yes|YES|Yes)
    if command -v uvx >/dev/null 2>&1; then
      uvx linkedin-scraper-mcp@latest --login
    else
      printf '%s\n' 'uvx não encontrado; pulei o login do LinkedIn MCP.'
    fi
    ;;
esac

if [ -z "${DEEPLINE_API_KEY:-}" ]; then
  printf '%s' 'Cole DEEPLINE_API_KEY, ou deixe vazio para pular: ' > /dev/tty
  prompt_secret_env DEEPLINE_API_KEY ''
fi

printf '%s' 'Configurar BYOK do Hunter agora? [s/N] ' > /dev/tty
IFS= read -r hunter_byok < /dev/tty
case "${hunter_byok:-N}" in
  s|S|sim|SIM|Sim|y|Y|yes|YES|Yes)
    contact-goat-pp-cli config byok set hunter HUNTER_API_KEY
    ;;
esac

printf '%s' 'Configurar BYOK do Apollo agora? [s/N] ' > /dev/tty
IFS= read -r apollo_byok < /dev/tty
case "${apollo_byok:-N}" in
  s|S|sim|SIM|Sim|y|Y|yes|YES|Yes)
    contact-goat-pp-cli config byok set apollo APOLLO_API_KEY
    ;;
esac

printf '\n%s\n' '== Doctor final =='
run_doctors

printf '\n%s\n' 'Se quiser persistir variáveis no shell:'
printf '%s\n' 'export CASA_DADOS_API_KEY="..."'
printf '%s\n' 'export DEEPLINE_API_KEY="..."'
printf '%s\n' 'export HAPPENSTANCE_API_KEY="..."'
printf '%s\n' 'export HUNTER_API_KEY="..."'
printf '%s\n' 'export APOLLO_API_KEY="..."'
