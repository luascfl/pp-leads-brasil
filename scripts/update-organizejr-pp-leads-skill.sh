#!/usr/bin/env bash
set -euo pipefail

skill_slug="${SKILL_SLUG:-organizejr-pp-leads}"
dry_run=false
all_known=false
declare -a target_roots=()
declare -a exact_targets=()

usage() {
  printf '%s\n' "Usage: $0 [--dry-run] [--all-known] [--target-root DIR] [--target DIR]"
  printf '%s\n' ""
  printf '%s\n' "Updates the installed $skill_slug skill from the newest project copy."
  printf '%s\n' ""
  printf '%s\n' "Source candidates, newest SKILL.md wins:"
  printf '%s\n' "  ./organizejr-pp-leads/SKILL.md"
  printf '%s\n' "  ./.context/skills/organizejr-pp-leads/SKILL.md"
  printf '%s\n' ""
  printf '%s\n' "Default target:"
  printf '%s\n' "  ~/.codex/skills/organizejr-pp-leads"
  printf '%s\n' ""
  printf '%s\n' "Options:"
  printf '%s\n' "  --dry-run          Show what would change without writing."
  printf '%s\n' "  --all-known        Also update existing ~/.claude/skills and ~/.gemini/skills roots."
  printf '%s\n' "  --target-root DIR  Install under DIR/$skill_slug. Repeatable."
  printf '%s\n' "  --target DIR       Install to exact target skill directory. Repeatable."
  printf '%s\n' "  -h, --help         Show this help."
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --dry-run)
      dry_run=true
      shift
      ;;
    --all-known)
      all_known=true
      shift
      ;;
    --target-root)
      if [ "$#" -lt 2 ]; then
        printf '%s\n' "error: --target-root needs a directory" >&2
        exit 2
      fi
      target_roots+=("$2")
      shift 2
      ;;
    --target)
      if [ "$#" -lt 2 ]; then
        printf '%s\n' "error: --target needs a directory" >&2
        exit 2
      fi
      exact_targets+=("$2")
      shift 2
      ;;
    -h|--help)
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

script_dir="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(CDPATH= cd -- "$script_dir/.." && pwd)"

declare -a source_dirs=(
  "$repo_root/$skill_slug"
  "$repo_root/.context/skills/$skill_slug"
)

latest_dir=""
for dir in "${source_dirs[@]}"; do
  skill_file="$dir/SKILL.md"
  if [ ! -f "$skill_file" ]; then
    continue
  fi
  if [ -z "$latest_dir" ] || [ "$skill_file" -nt "$latest_dir/SKILL.md" ]; then
    latest_dir="$dir"
  fi
done

if [ -z "$latest_dir" ]; then
  printf '%s\n' "error: no project copy found for $skill_slug" >&2
  exit 1
fi

read -r source_sha _ < <(sha256sum "$latest_dir/SKILL.md")
printf '%s\n' "source: $latest_dir"
printf '%s\n' "source_sha256: $source_sha"

if [ "${#target_roots[@]}" -eq 0 ] && [ "${#exact_targets[@]}" -eq 0 ]; then
  target_roots+=("$HOME/.codex/skills")
fi

if [ "$all_known" = true ]; then
  for root in "$HOME/.claude/skills" "$HOME/.gemini/skills"; do
    if [ -d "$root" ]; then
      target_roots+=("$root")
    fi
  done
fi

install_one() {
  src="$1"
  dst="$2"

  if [ "$src" = "$dst" ]; then
    printf '%s\n' "skip: source and target are the same: $dst"
    return 0
  fi

  if [ "$dry_run" = true ]; then
    printf '%s\n' "would update: $dst"
    return 0
  fi

  mkdir -p "$(dirname -- "$dst")"

  if [ -e "$dst" ]; then
    if [ -f "$dst/SKILL.md" ] && cmp -s "$src/SKILL.md" "$dst/SKILL.md"; then
      rm -rf "$dst"
    else
      backup="$dst.backup.$(date +%Y%m%d%H%M%S)"
      mv "$dst" "$backup"
      printf '%s\n' "backup: $backup"
    fi
  fi

  mkdir -p "$dst"
  cp -a "$src"/. "$dst"/
  printf '%s\n' "updated: $dst"
}

for root in "${target_roots[@]}"; do
  install_one "$latest_dir" "$root/$skill_slug"
done

for target in "${exact_targets[@]}"; do
  install_one "$latest_dir" "$target"
done
