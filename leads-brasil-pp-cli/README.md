# Leads Brasil CLI

Platform for Brazilian Lead Intelligence (CNPJs, Companies, People)

## Install

The recommended path installs both the `leads-brasil-pp-cli` binary and the `pp-leads-brasil` agent skill (Claude Code, Codex, Cursor, Gemini CLI, GitHub Copilot, and other agents supported by the upstream [`skills`](https://github.com/vercel-labs/skills) CLI) in one shot:

```bash
npx -y @mvanhorn/printing-press-library install leads-brasil
```

For CLI only (no skill):

```bash
npx -y @mvanhorn/printing-press-library install leads-brasil --cli-only
```

For skill only — installs the skill into the same agents as the default command above, but skips the CLI binary (use this to update or reinstall just the skill):

```bash
npx -y @mvanhorn/printing-press-library install leads-brasil --skill-only
```

To constrain the skill install to one or more specific agents (repeatable — agent names match the [`skills`](https://github.com/vercel-labs/skills) CLI):

```bash
npx -y @mvanhorn/printing-press-library install leads-brasil --agent claude-code
npx -y @mvanhorn/printing-press-library install leads-brasil --agent claude-code --agent codex
```

### Without Node

The generated install path is category-agnostic until this CLI is published. If `npx` is not available before publish, install Node or use the category-specific Go fallback from the public-library entry after publish.

### Pre-built binary

Download a pre-built binary for your platform from the [latest release](https://github.com/mvanhorn/printing-press-library/releases/tag/leads-brasil-current). On macOS, clear the Gatekeeper quarantine: `xattr -d com.apple.quarantine <binary>`. On Unix, mark it executable: `chmod +x <binary>`.

<!-- pp-hermes-install-anchor -->
## Install for Hermes

From the Hermes CLI:

```bash
hermes skills install mvanhorn/printing-press-library/cli-skills/pp-leads-brasil --force
```

Inside a Hermes chat session:

```bash
/skills install mvanhorn/printing-press-library/cli-skills/pp-leads-brasil --force
```

## Install for OpenClaw

Tell your OpenClaw agent (copy this):

```
Install the pp-leads-brasil skill from https://github.com/mvanhorn/printing-press-library/tree/main/cli-skills/pp-leads-brasil. The skill defines how its required CLI can be installed.
```

## Use with Claude Desktop

This CLI ships an [MCPB](https://github.com/modelcontextprotocol/mcpb) bundle — Claude Desktop's standard format for one-click MCP extension installs (no JSON config required).

To install:

1. Download the `.mcpb` for your platform from the [latest release](https://github.com/mvanhorn/printing-press-library/releases/tag/leads-brasil-current).
2. Double-click the `.mcpb` file. Claude Desktop opens and walks you through the install.
3. Fill in `LEADS_BRASIL_BEARER_AUTH` when Claude Desktop prompts you.

Requires Claude Desktop 1.0.0 or later. Pre-built bundles ship for macOS Apple Silicon (`darwin-arm64`) and Windows (`amd64`, `arm64`); for other platforms, use the manual config below.

<details>
<summary>Manual JSON config (advanced)</summary>

If you can't use the MCPB bundle (older Claude Desktop, unsupported platform), install the MCP binary and configure it manually.


Install the MCP binary from this CLI's published public-library entry or pre-built release.

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "leads-brasil": {
      "command": "leads-brasil-pp-mcp",
      "env": {
        "LEADS_BRASIL_BEARER_AUTH": "<your-key>"
      }
    }
  }
}
```

</details>

## Quick Start

### 1. Install

See [Install](#install) above.

### 2. Set Up Credentials

Get your access token from your API provider's developer portal, then store it:

```bash
leads-brasil-pp-cli auth set-token YOUR_TOKEN_HERE
```

Or set it via environment variable:

```bash
export LEADS_BRASIL_BEARER_AUTH="your-token-here"
```

### 3. Verify Setup

```bash
leads-brasil-pp-cli doctor
```

This checks your configuration and credentials.

### 4. Try Your First Command

```bash
leads-brasil-pp-cli company mock-value
```

## Usage

Run `leads-brasil-pp-cli --help` for the full command reference and flag list.

## Commands

### company

Manage company

- **`leads-brasil-pp-cli company <cnpj>`** - Get company data by CNPJ

### enrich

Manage enrich

- **`leads-brasil-pp-cli enrich <cnpj>`** - Trigger enrichment for a company via PP Suite

### leads-brasil-platform-search

Manage leads brasil platform search

- **`leads-brasil-pp-cli leads-brasil-platform-search`** - Search for companies by name or filter

### operation

Create a side-effect-free plan from a profile-owned JSON input, then apply it only with separately supplied approval:

```bash
leads-brasil-pp-cli operation plan --input operation.json
leads-brasil-pp-cli operation apply <plan-id> --yes
```

`--agent` does not approve an external mutation. `operation apply` requires an explicit `--yes`.


## Output Formats

```bash
# Human-readable table (default in terminal, JSON when piped)
leads-brasil-pp-cli company mock-value

# JSON for scripting and agents
leads-brasil-pp-cli company mock-value --json

# Filter to specific fields
leads-brasil-pp-cli company mock-value --json --select id,name,status

# Dry run — show the request without sending
leads-brasil-pp-cli company mock-value --dry-run

# Agent mode — JSON + compact + no prompts in one flag
leads-brasil-pp-cli company mock-value --agent
```

## Agent Usage

This CLI is designed for AI agent consumption:

- **Non-interactive** - never prompts, every input is a flag
- **Pipeable** - `--json` output to stdout, errors to stderr
- **Filterable** - `--select id,name` returns only fields you need
- **Previewable** - `--dry-run` shows the request without sending
- **Explicit retries** - add `--idempotent` to create retries when a no-op success is acceptable
- **Confirmable** - `--yes` for explicit confirmation of destructive actions
- **Piped input** - write commands can accept structured input when their help lists `--stdin`
- **Agent-safe by default** - no colors or formatting unless `--human-friendly` is set

Exit codes: `0` success, `2` usage error, `3` not found, `4` auth error, `5` API error, `7` rate limited, `10` config error.

## Health Check

```bash
leads-brasil-pp-cli doctor
```

Verifies configuration, credentials, and connectivity to the API.

## Configuration

Config file: `~/.config/leads-brasil-platform-pp-cli/config.toml`

Static request headers can be configured under `headers`; per-command header overrides take precedence.

Environment variables:

| Name | Kind | Required | Description |
| --- | --- | --- | --- |
| `LEADS_BRASIL_BEARER_AUTH` | per_call | Yes | Set to your API credential. |

## Troubleshooting
**Authentication errors (exit code 4)**
- Run `leads-brasil-pp-cli doctor` to check credentials
- Verify the environment variable is set: `echo $LEADS_BRASIL_BEARER_AUTH`
**Not found errors (exit code 3)**
- Check the resource ID is correct
- Run the `list` command to see available items

---

Generated by [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press)
