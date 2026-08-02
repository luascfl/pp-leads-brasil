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

## m21: External-operation contract

**Goal:** let explicit private profiles perform bounded external mutations through a single auditable plan and receipt contract.

### Phase 21.1: Plan, apply and receipt boundary

**Story:** `m21-us001` — Introduce generic external operation plans.

**Scope**
- Define immutable operation plans, target diffs and receipts.
- Make `plan` side-effect free and make `apply` require plan ID plus explicit approval.
- Restrict adapters to profile-owned integrations and prevent enrichment providers from mutating systems.
- Implement idempotent replay and observable partial failure.

**Dependencies:** m20 completed. Detailed execution: [`m21-external-operation-contract/PLAN.md`](./m21-external-operation-contract/PLAN.md).

**Definition of done**
- Plan/apply operations are bounded by stored target diffs and expiry.
- Every target writes a receipt with evidence and before/after values.
- Public tests pass without a private profile checkout.

## m22: Audited CRM synchronization

**Goal:** implement a private Google Sheets → Ploomes adapter that is gated by m21 plans, explicit approval and durable per-field receipts.

### Phase 22.1: Private adapter and reconciliation

**Story:** `m22-us001` — Plan, apply and reconcile explicit CRM field changes.

**Scope**
- Keep Google Sheets as the operational source of truth and Ploomes as the execution mirror.
- Validate a private schema contract, row lock, immutable lead identity, Ploomes ID and before value before every mutation.
- Build plans only from explicitly selected rows, with one field-level target and evidence reference per change.
- Remove or disable any uncontrolled timed Apps Script sender before enabling the adapter.
- Restrict initial mutation scope to validated standard company/person fields; defer custom Ploomes fields.

**Dependencies:** m21 completed. Detailed execution: [`m22-audited-crm-sync/PLAN.md`](./m22-audited-crm-sync/PLAN.md).

**Definition of done**
- A private adapter applies only m21-approved field-level targets and records reconciliation evidence.
- Schema, lock, identity and before-value conflicts fail closed.
- The public checkout remains profile-agnostic and its test suite passes without the private sibling repository.

**Follow-on phases, not yet planned**
- m23: immutable lead identity and field-level provenance repair.
