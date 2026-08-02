# m22: audited CRM synchronization

## Goal

Make a private Google Sheets → Ploomes adapter use the public m21 plan, approval and receipt contract without exposing CRM identifiers, credentials, lead data or private automation in `pp-leads-brasil`.

## Decision

Google Sheets remains the operational and analytical source of truth. Ploomes remains the execution mirror for basic company and person records. Until Ploomes FieldKeys are explicitly mapped and verified, analytical fields remain Sheets-only and must not be encoded into `Note` as a covert structured sync.

The adapter lives only in the private OrganizeJr profile. The public platform keeps the generic `operation.Adapter` interface and receives no CRM-specific dependency.

## Contract

### Read and plan

1. Read explicitly selected sheet rows and the matching Ploomes records.
2. Validate a declared sheet schema before any write-capable path. The schema names each column, owner, allowed writer, sync direction, overwrite rule and Ploomes destination.
3. Reject rows that are locked, have no immutable `lead_enrichment_id`, lack Ploomes IDs where an update is requested, have a missing or renamed required header, or have a source conflict.
4. Build one `operation.TargetChange` per field. Each key is stable and includes `lead_enrichment_id`, Ploomes record ID and Ploomes field key or standard field name. Evidence references the private dossier or source record, never embeds raw content.
5. Persist the m21 plan. Planning may write only a private audit-preview record. It must not call a Ploomes write endpoint or alter Sheet business data.

### Review and apply

1. The operator reviews the recorded before/after diff and selects the stored plan ID.
2. Apply only through `operation apply <plan-id> --yes` or the equivalent approved API request.
3. The adapter rechecks the target identity, before value and schema hash immediately before the Ploomes mutation. A changed value becomes a conflict receipt, not an overwrite.
4. On success, append a private audit record with plan digest, target key, record IDs, field, before/after, evidence reference, operator, timestamp, Ploomes response ID and result.
5. Reapply returns existing successful receipts. Failed targets remain explicit and retryable only through the stored plan while it is unexpired.

## Scope

- Add a private `crm_schema.json` manifest and schema validator.
- Add a private plan builder from explicit row numbers, never from all eligible rows.
- Add a private Ploomes adapter implementing `operation.Adapter` for standard company/person fields.
- Add a private, append-only audit log and per-row lock check.
- Disable or prove inactive the 15-minute `enviarAoPloomes` trigger before enabling the adapter.
- Reconcile each successful target with read-only Google Sheets and Ploomes reads.

## Field policy

| Field group | m22 policy |
| --- | --- |
| Standard company/person fields, Ploomes IDs, phones | Eligible after schema and identity validation. |
| `Funil`, `Próximo passo`, `Score ICP (0-10)`, `ICP Considerado` | Sheets-only until FieldKeys and type semantics are explicitly mapped. |
| `Observações`, Perplexity URL, social links, dossier text | Evidence/reference only. No automatic destructive overwrite. |
| `Status Importação`, `Erro Importação`, `Última sincronização` | Adapter control/audit fields, never business-data evidence. |

## Non-goals

- No public Ploomes or Google Sheets client, configuration, fixture, credential or lead record.
- No `enviarAoPloomes`, `sincronizarBidirecionalPloomes`, scheduled trigger or all-eligible-row execution.
- No creation of Ploomes tasks, messages, WhatsApp drafts or contacts.
- No custom-field write before a separately approved FieldKey mapping story.
- No source-of-truth reversal from Ploomes into analytical Sheet columns.

## Implementation order

1. **Schema and trigger gate**: declare the schema, add a validator and inspect Apps Script triggers. Remove or disable the timed trigger if present.
2. **Private plan builder**: read selected rows, Ploomes records and evidence references; emit only m21-compatible targets and an audit preview.
3. **Private adapter**: revalidate schema, identity and before value; update one standard field per target; expose deterministic test doubles.
4. **Audit and reconciliation**: append receipts to the private ledger, reread both systems and surface conflicts.
5. **Custom fields decision**: discover `GET /Fields` FieldKeys and add a separate plan only if values, types, ownership and overwrite policy are approved.

## Acceptance criteria

- A clean public checkout still passes `go test ./...` without the private profile.
- A malformed/missing Sheet header, lock, missing immutable identity, missing Ploomes ID or before-value conflict prevents Ploomes mutation.
- Planning an explicit selection creates a side-effect-free m21 plan containing only field-level diffs with evidence references.
- Apply without `--yes`, after expiry or with altered schema/data fails closed.
- One successful target produces a durable operation receipt and a private audit entry; replay does not invoke Ploomes again.
- A Ploomes failure affects only its target and is visible in reconciliation output.
- No scheduled or bulk Apps Script route can bypass the plan-and-apply contract.

## Verification

- Unit tests for schema validation, stable target keys, selection bounds, locks, missing IDs, before-value conflicts, failed target retries and successful replay.
- Fixture integration test with a fake Sheets reader and fake Ploomes client proving plan → approve → receipt → reconciliation.
- Manual read-only check of actual Apps Script triggers before any production apply.
- One manually approved non-sensitive test field in the private profile, followed by readback from Ploomes and Google Sheets.
