# Project plan

## Projeto
pp-leads-brasil, stack de pesquisa, enrichment e CRM comercial da OrganizeJr.

## Estado macro
- m13, consolidação de duplicatas locais, concluído.
- m14, backfill cadastral prioritário, concluído.
- m15, expansão de contato com evidência para empresas não resolvidas, concluído.
- m16, sincronização CRM -> Perplexity Space, ativo.

## Milestone ativo, m16

### Objetivo
Transformar o estado canônico mínimo do Google Sheets CRM em snapshots curtos, coláveis no Perplexity Space da OrganizeJr, para que o Space responda melhor a dumps comerciais desestruturados sem virar fonte canônica do CRM.

### Dependências
- Google Sheets CRM acessível via `crm_sheet.py`.
- Ledger `organizejr-pp-leads/ploomes_crm/enrichment_attempts.jsonl` preservando bloqueios anti-loop.
- Evidências locais em `organizejr-pp-leads/leads/<lead>/` quando existirem.
- Perplexity continua como dump narrativo e compartilhado, não como fonte final para campos cadastrais.

### Critérios de conclusão
- Existe comando read-only para exportar snapshots Markdown do CRM.
- Snapshot diferencia status, lacunas, observações, canais, CNPJ e log técnico.
- Output pode ser salvo como artefato local e colado no Perplexity Space.
- Validação técnica registrada em `.context/plans/STATE.md`.
