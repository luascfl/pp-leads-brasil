#!/usr/bin/env bash
set -euo pipefail

need_bin() {
  if ! command -v "$1" >/dev/null 2>&1; then
    printf '%s\n' "error: comando ausente: $1" >&2
    exit 127
  fi
}

need_bin contact-goat-pp-cli
need_bin python3

printf '%s\n' '== Doctor inicial =='
contact-goat-pp-cli doctor || true

if [ -t 0 ] && [ -r /dev/tty ]; then
  if [ -z "${CASA_DADOS_API_KEY:-${PP_LEADS_CASA_DADOS_API_KEY:-}}" ]; then
    printf '%s' 'CASA_DADOS_API_KEY não está definida. Quer exportar agora? [s/N] ' > /dev/tty
    IFS= read -r answer < /dev/tty
    case "${answer:-N}" in
      s|S|sim|SIM|Sim|y|Y|yes|YES|Yes)
        printf '%s' 'Cole CASA_DADOS_API_KEY: ' > /dev/tty
        IFS= read -r -s casa_key < /dev/tty
        printf '\n' > /dev/tty
        if [ -n "$casa_key" ]; then
          export CASA_DADOS_API_KEY="$casa_key"
          printf '%s\n' 'CASA_DADOS_API_KEY exportada para esta sessão.'
        fi
        unset casa_key
        ;;
    esac
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
      uvx linkedin-scraper-mcp@latest --login
      ;;
  esac

  if [ -z "${DEEPLINE_API_KEY:-}" ]; then
    printf '%s' 'DEEPLINE_API_KEY não está definida. Quer exportar agora? [s/N] ' > /dev/tty
    IFS= read -r deepline_answer < /dev/tty
    case "${deepline_answer:-N}" in
      s|S|sim|SIM|Sim|y|Y|yes|YES|Yes)
        printf '%s' 'Cole DEEPLINE_API_KEY: ' > /dev/tty
        IFS= read -r -s deepline_key < /dev/tty
        printf '\n' > /dev/tty
        if [ -n "$deepline_key" ]; then
          export DEEPLINE_API_KEY="$deepline_key"
          printf '%s\n' 'DEEPLINE_API_KEY exportada para esta sessão.'
        fi
        unset deepline_key
        ;;
    esac
  fi

  printf '%s' 'Configurar BYOK do Hunter em contact-goat agora? [s/N] ' > /dev/tty
  IFS= read -r hunter_byok < /dev/tty
  case "${hunter_byok:-N}" in
    s|S|sim|SIM|Sim|y|Y|yes|YES|Yes)
      contact-goat-pp-cli config byok set hunter HUNTER_API_KEY
      ;;
  esac

  printf '%s' 'Configurar BYOK do Apollo em contact-goat agora? [s/N] ' > /dev/tty
  IFS= read -r apollo_byok < /dev/tty
  case "${apollo_byok:-N}" in
    s|S|sim|SIM|Sim|y|Y|yes|YES|Yes)
      contact-goat-pp-cli config byok set apollo APOLLO_API_KEY
      ;;
  esac
fi

printf '%s\n' ''
printf '%s\n' '== Doctor final =='
contact-goat-pp-cli doctor || true

printf '%s\n' ''
printf '%s\n' 'Se quiser persistir variáveis no shell, adicione ao seu ~/.bashrc ou ~/.zshrc:'
printf '%s\n' 'export CASA_DADOS_API_KEY="..."'
printf '%s\n' 'export DEEPLINE_API_KEY="..."'
printf '%s\n' 'export HAPPENSTANCE_API_KEY="..."   # opcional, para bearer surface'
printf '%s\n' 'export HUNTER_API_KEY="..."         # opcional, se usar BYOK'
printf '%s\n' 'export APOLLO_API_KEY="..."         # opcional, se usar BYOK'
