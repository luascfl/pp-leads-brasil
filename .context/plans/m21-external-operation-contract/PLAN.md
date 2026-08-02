# m21: external-operation contract

## Goal

Allow an explicitly configured private profile to propose and apply a bounded external mutation without embedding any CRM, messaging, or profile-specific integration in the public platform.

## Contract

An external operation is a two-step workflow:

1. `plan` validates a declared profile, target set, proposed field changes, source evidence and operation kind. It persists an immutable operation plan and returns its ID, digest, targets, diff, expiration and dry-run result. It has no external side effect.
2. `apply` accepts only a previously persisted, unexpired plan ID and an explicit approval flag. It applies only its recorded targets, writes one receipt per target, and returns a partial-failure result without replaying successful targets.

The contract applies to profile-owned adapters only. Scrape Creators, Apify, Casa dos Dados, Company Goat and Contact Goat provide evidence; they never apply a mutation.

## Scope

- Add generic plan, target diff, receipt and result types under the public operation layer.
- Require a configured profile for all external operations.
- Add a file-backed local plan and receipt store with canonical JSON digests.
- Expose separate plan and apply API/CLI operations. `apply` requires the plan ID and `--yes`.
- Make replay idempotent by receipt and target key.
- Return structured failures for missing profile, empty target set, expired plan, digest mismatch, unapproved apply and adapter failure.

## Non-goals

- No bundled Google Sheets, Ploomes, WhatsApp, Telegram, email or contact adapter.
- No generic `--agent` authorization for mutation.
- No target discovery during apply.
- No direct write from enrichment or social-provider output.
- No secret persistence in plans or receipts.

## Implementation slices

### 1. Operation domain and persistence

Add immutable `OperationPlan`, `TargetChange`, `OperationReceipt` and `OperationResult` types. A plan contains profile identity, kind, source references, ordered targets, proposed field diffs, digest and expiry. The store writes plans and receipts in a profile-owned output directory, rejecting path traversal and secret-shaped fields.

### 2. Validation and adapter boundary

Add a generic `OperationAdapter` that receives a single recorded target change. The validator rejects empty plans, duplicate targets, missing evidence, profile mismatch and unsupported kinds before persistence. The public factory holds only this interface, never a private adapter implementation.

### 3. Plan and apply endpoints

Add `POST /operations/plan` and `POST /operations/{plan_id}/apply`. Planning is dry-run only. Apply requires a header or request field proving explicit confirmation and fails closed when absent. The server verifies the stored digest and expiration before invoking adapters.

### 4. CLI and receipts

Generate or extend CLI commands so `operation plan` emits machine-readable diffs and `operation apply <id> --yes` emits receipts. `--agent` must not imply approval for these commands. A retry returns recorded receipts for completed target keys and invokes the adapter only for failed or pending targets.

### 5. Tests and documentation

Test empty target rejection, out-of-scope target rejection, missing approval, expired plan, unchanged replay, single-target partial failure and absence of private profile paths. Document the profile manifest extension and the evidence-only boundary for enrichment providers.

## Acceptance criteria

- Planning has no external side effect and returns an immutable digest.
- Applying requires explicit plan ID plus approval and cannot alter an unplanned target.
- Each attempted target has a durable receipt with before/after values, evidence reference, operator, timestamp and result.
- Reapplying a completed target is idempotent.
- A partial failure is observable and does not replay successful targets.
- `go test ./...` passes without a private profile checkout.
