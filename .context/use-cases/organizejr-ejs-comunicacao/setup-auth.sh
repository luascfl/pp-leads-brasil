#!/usr/bin/env bash
set -euo pipefail

script_dir="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(CDPATH= cd -- "$script_dir/../../.." && pwd)"

exec "$repo_root/scripts/setup-external-enrichment-auth.sh" "$@"
