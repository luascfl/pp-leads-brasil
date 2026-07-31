---
type: doc
name: project-overview
description: High-level overview of the project, its purpose, and key components
category: overview
generated: 2026-07-25
status: unfilled
scaffoldVersion: "2.0.0"
---
## Project Overview

This repository is a public, profile-agnostic lead-intelligence platform. It combines a locally declared lead table, Casa dos Dados, Company Goat, Contact Goat and social-media enrichment to assemble evidence for a lead without embedding a specific CRM or commercial operation.

Private profiles are external integrations. They provide their own data location, field mapping, credentials and external-operation policy through explicit configuration. Public source and documentation must not contain private CRM records, evidence files, automation paths or profile-specific commands.

## Enrichment flow

1. Resolve a lead from the configured local table.
2. Enrich company facts through Casa dos Dados and Company Goat when configured.
3. Enrich professional contact signals through Contact Goat when configured.
4. Attempt Scrape Creators for public social profiles.
5. When that attempt has no usable result, use the configured Apify Actor as a non-blocking fallback.

Provider errors, absent credentials and billing limits are represented in the enrichment payload. They do not invalidate evidence already obtained from other sources.

## External operations

This public repository does not embed CRM synchronization, messaging or contact mutation. Private profiles must declare those operations separately and make their execution auditable, previewable and explicitly approved.

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
