---
type: doc
name: data-flow
description: How data moves through the system and external integrations
category: data-flow
generated: 2026-07-24
status: unfilled
scaffoldVersion: "2.0.0"
---
## Data Flow & Integrations

This document describes how data flows through the system, including internal processing and external integrations.

Understanding data flow helps with debugging, performance optimization, and maintaining system reliability.

## Module Dependencies

Module dependency overview:

- **Entry Layer** → Services, Utils
- **Services** → Data Access, External APIs
- **Data Access** → Database, Cache

*See [`codebase-map.json`](./codebase-map.json) for detailed dependency graphs.*

## Service Layer

Key services in the system:

- **[ServiceName]** — [Purpose] (`src/services/path.ts`)

*See [`codebase-map.json`](./codebase-map.json) for complete service listings.*

## High-level flow

```mermaid
flowchart LR
    A[Lead table in Google Sheets] --> B[crm_sheet.py analyze-health]
    C[Local evidence: Perplexity markdown, goat JSONs, web checks] --> D[JSON batches]
    D --> E[crm_sheet.py apply-icp-score]
    E --> A
    A --> F[Health dashboard and enrichment queue]
```

**Data flow steps**:
1. `crm_sheet.py` reads the `Clientes` Google Sheet through Google Sheets OAuth credentials.
2. `analyze-health` normalizes each row, preserves existing generated columns, computes CRM health, and writes structured health/status columns back to the sheet.
3. Enrichment evidence is converted into JSON batches containing `row`, `score`, `icp_considerado`, `resumo`, `justificativa` and optional contact fields.
4. `apply-icp-score` updates existing rows in place, fills empty contact cells only, appends chronological observations, and can append new lead rows with `row: "new"`.
5. The final dashboard reports coverage for person, actionable channel, company e-mail, CNPJ, ICP, score and observations.

## Internal Movement

<!-- Describe how modules collaborate (queues, events, RPC calls, shared databases). -->

_Add descriptive content here (optional)._

## External Integrations

**External Services**:

| Service | Purpose | Auth Method |
|---------|---------|-------------|
| [Service] | [Purpose] | [API Key/OAuth/etc.] |

*Document error handling and retry strategies for each integration.*

## Observability & Failure Modes

<!-- Describe metrics, traces, or logs that monitor the flow. Note backoff, dead-letter, or compensating actions. -->

_Add descriptive content here (optional)._

## Related Resources

<!-- Link to related documents for cross-navigation. -->

- [architecture.md](./architecture.md)
