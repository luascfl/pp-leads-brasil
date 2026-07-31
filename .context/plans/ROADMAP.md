# Roadmap

## m20: Public profile boundary

**Goal:** restore a strict public/private boundary while preserving explicit external profile support.

### Phase 20.1: Explicit profile resolution

**Story:** `m20-us001` — Remove implicit OrganizeJr discovery and isolate public fixtures.

**Scope**
- Delete the public CLI's automatic profile discovery path.
- Preserve `usecase.Config` loading from explicitly supplied configuration.
- Replace private-path test fixtures with repository-owned public fixtures.
- Update `doctor` profile reporting.
- Correct public architecture and data-flow documentation.
- Make Scrape Creators the required social-enrichment stage with a configured, non-blocking Apify fallback.
- Record, but do not delete, stale legacy context artifacts.

**Dependencies:** none.

**Definition of done**
- All m20 requirements pass.
- `go test ./...` passes from a checkout without the private sibling repository.
- Graphify output is incrementally updated after code changes.
- Canonical state, Ralph PRD and technical context contain validation evidence.

**Follow-on phases, not yet planned**
- m21: gated external-operation contract for private profiles.
- m22: audited Google Sheets and Ploomes synchronization.
- m23: immutable lead identity and field-level provenance repair.
