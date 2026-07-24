---
type: doc
name: project-overview
description: High-level overview of the project, its purpose, and key components
category: overview
generated: 2026-07-24
status: unfilled
scaffoldVersion: "2.0.0"
---
## Project overview

This repository supports OrganizeJr lead prospecting and CRM enrichment for Brazilian empresas juniores and related education/commercial leads. The active workflow turns local evidence, Perplexity exports, public web data and Google Sheets CRM rows into structured enrichment: ICP classification, score, contact channels, CNPJ, observations and enrichment queue diagnostics.

The current operational entry point for the CRM loop is `organizejr-pp-leads/ploomes_crm/crm_sheet.py`. It reads and writes the `Clientes` Google Sheet, calculates CRM health through `analyze-health`, lists missing ICP/score through `list-icp-pending`, and injects score/contact evidence through `apply-icp-score`.

## Codebase Reference

> **Detailed Analysis**: For complete symbol counts, architecture layers, and dependency graphs, see [`codebase-map.json`](./codebase-map.json).

## Quick Facts

- **Root**: `./`
- **Primary Language**: [Language] ([X] files)
- **Entry Point**: `src/index.ts` or `src/main.ts`
- **Full Analysis**: [`codebase-map.json`](./codebase-map.json)

## Entry Points

- **Main Entry**: `src/index.ts` - Primary module exports
- **CLI**: `src/cli.ts` - Command-line interface (if applicable)
- **Server**: `src/server.ts` - HTTP server entry (if applicable)

## Key Exports

See [`codebase-map.json`](./codebase-map.json) for the complete list of exported symbols.

Key public APIs:
- [List main exported classes/functions]

## File Structure & Code Organization

- `src/` — Source code and main application logic
- `tests/` or `__tests__/` — Test files and fixtures
- `dist/` or `build/` — Compiled output (gitignored)
- `docs/` — Documentation files
- `scripts/` — Build and utility scripts

## Technology Stack Summary

**Runtime**: Node.js

**Language**: TypeScript/JavaScript

**Build Tools**:
- TypeScript compiler (tsc) or bundler (esbuild, webpack, etc.)
- Package manager: npm/yarn/pnpm

**Code Quality**:
- Linting: ESLint
- Formatting: Prettier
- Type checking: TypeScript strict mode

## Core Framework Stack

<!-- Document core frameworks per layer (backend, frontend, data, messaging). Mention architectural patterns enforced by these frameworks. -->

_Add descriptive content here (optional)._

## UI & Interaction Libraries

<!-- List UI kits, CLI interaction helpers, or design system dependencies. Note theming, accessibility, or localization considerations. -->

_Add descriptive content here (optional)._

## Development Tools Overview

See [Tooling](./tooling.md) for detailed development environment setup.

**Essential Commands**:
- `npm install` — Install dependencies
- `npm run build` — Build the project
- `npm run test` — Run tests
- `npm run dev` — Start development mode

## Getting Started Checklist

1. Clone the repository
2. Install dependencies: `npm install`
3. Copy environment template: `cp .env.example .env` (if applicable)
4. Run tests to verify setup: `npm run test`
5. Start development: `npm run dev`
6. Review [Development Workflow](./development-workflow.md) for day-to-day tasks

## Next Steps

- Review [Architecture](./architecture.md) for system design details
- See [Development Workflow](./development-workflow.md) for contribution guidelines
- Check [Testing Strategy](./testing-strategy.md) for quality requirements

## Related Resources

<!-- Link to related documents for cross-navigation. -->

- [architecture.md](./architecture.md)
- [development-workflow.md](./development-workflow.md)
- [tooling.md](./tooling.md)
- [codebase-map.json](./codebase-map.json)
