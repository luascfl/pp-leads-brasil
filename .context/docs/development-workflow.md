---
type: doc
name: development-workflow
description: Day-to-day engineering processes, branching, and contribution guidelines
category: workflow
generated: 2026-07-25
status: unfilled
scaffoldVersion: "2.0.0"
---
## Development Workflow

This document outlines the day-to-day engineering process for contributing to this repository.

Following these guidelines ensures consistent code quality and smooth collaboration across the team.

## Branching & Releases

**Branching Model**: Feature branches off `main`

- `main` — Production-ready code, always deployable
- `feature/*` — New features and enhancements
- `fix/*` — Bug fixes
- `chore/*` — Maintenance and tooling updates

**Release Process**:
1. Features are developed in branches
2. PRs require review and passing CI
3. Merged PRs are deployed automatically (or tagged for release)

**Versioning**: Semantic versioning (semver) - MAJOR.MINOR.PATCH

## Local Development

**Setup**:
```bash
# Clone and install
git clone <repository-url>
cd <project-name>
npm install
```

**Daily Commands**:
- `npm run dev` — Start development server/watch mode
- `npm run build` — Build for production
- `npm run test` — Run test suite
- `npm run lint` — Check code style


**Telegram people enrichment**:
```bash
# Dry-run, no Telegram read.
python3 organizejr-pp-leads/ploomes_crm/crm_sheet.py telegram-people search --term ENGETOP

# Real Telegram read through manage_telegram, still no contact mutation.
python3 organizejr-pp-leads/ploomes_crm/crm_sheet.py telegram-people search --term ENGETOP --execute-read

# Real Telegram read for a Ploomes lead row without contact data. This returns candidate names and payload candidates, but does not write to Ploomes.
python3 organizejr-pp-leads/ploomes_crm/crm_sheet.py telegram-people search --term ENGETOP --execute-read --ploomes-row 42 --candidate-limit 10
```

The bridge uses `/home/lucas/Downloads/manage_telegram` by default or `MANAGE_TELEGRAM_PATH` when set. Do not copy Telegram automation into this repo; consume the bridge JSON as local enrichment evidence. For a Ploomes row missing `Nome - Pessoa`, pass `--ploomes-row` so the output includes `ploomes_contact_candidates` compatible with `apply-icp-score`. Review and choose the correct person before applying because Telegram may return multiple people for one company term. Contact saving in Telegram stays explicit via `telegram-people apply-save ... --yes`.

**Before Committing**:
```bash
npm run lint && npm run test && npm run build
```

## Code Review Expectations

**PR Requirements**:
- Clear description of changes and motivation
- Tests for new functionality
- Documentation updates for API changes
- Passing CI checks

**Review Checklist**:
- [ ] Code follows project conventions
- [ ] Tests cover the changes adequately
- [ ] No security vulnerabilities introduced
- [ ] Documentation is updated
- [ ] Commit messages follow conventions

**Approval**: At least one approving review required before merge.

See [AGENTS.md](../../AGENTS.md) for AI assistant collaboration guidelines.

## Onboarding Tasks

**First Steps for New Contributors**:
1. Read the [Project Overview](./project-overview.md)
2. Set up local development environment
3. Run the test suite to verify setup
4. Look for issues labeled `good-first-issue` or `help-wanted`

**Helpful Resources**:
- [Architecture Notes](./architecture.md) — System design overview
- [Testing Strategy](./testing-strategy.md) — How to write tests
- [CONTRIBUTING.md](../../CONTRIBUTING.md) — Contribution guidelines

## Related Resources

<!-- Link to related documents for cross-navigation. -->

- [testing-strategy.md](./testing-strategy.md)
- [tooling.md](./tooling.md)
