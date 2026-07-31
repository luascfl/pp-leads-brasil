# Requirements

## Milestone m20: public profile boundary

### Objective
Make `pp-leads-brasil` a public, profile-agnostic lead-intelligence platform. Private profiles, including OrganizeJr, must be opt-in and discovered only through explicit configuration.

### In scope
- Remove automatic OrganizeJr profile discovery from the public CLI.
- Keep explicit profile loading through `PP_LEADS_USE_CASE_CONFIG` and `PP_LEADS_ICP_DIR`.
- Replace public tests that read private OrganizeJr paths with public fixtures.
- Make `doctor` report whether an explicit profile is configured, without finding sibling directories.
- Remove public documentation that presents private CRM, evidence or commercial workflows as public architecture.
- Define a canonical-context cleanup plan for stale legacy context directories without deleting ambiguous material automatically.
- Make Scrape Creators the required primary social-enrichment stage and Apify a configured, non-blocking fallback.

### Out of scope
- CRM, Google Sheets, Ploomes, WhatsApp, Perplexity, Telegram and OrganizeJr private evidence.
- External-provider integrations beyond the configured Scrape Creators → Apify fallback.
- Automated writes to external systems.
- Deleting legacy context artifacts during this first story.

### Acceptance requirements
1. A clean public checkout runs `go test ./...` without `../organizejr-pp-leads` or any private lead table.
2. No public command sets `PP_LEADS_ICP_DIR` by scanning its current directory, executable directory or siblings.
3. Explicit profile configuration continues to load a profile and resolve its declared paths.
4. Public documentation describes profiles as external/private opt-in integrations, not bundled CRM behavior.
5. The story records exact legacy context locations and a non-destructive migration decision.
6. Social enrichment always attempts Scrape Creators; an unavailable result attempts Apify only with explicit Actor configuration, and both provider failures remain non-blocking.

### Validation
- `go test ./...`
- CLI smoke test without profile configuration
- CLI smoke test with a public test profile
- Search for hard-coded OrganizeJr profile activation in public source and tests
- `graphify update .` after source changes
