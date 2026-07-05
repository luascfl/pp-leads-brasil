# Ralph PRD context

## Produto

`pp-leads-brasil` é a base local para inteligência de leads brasileiros, com busca por empresas/CNPJs e enriquecimento via PP Suite.

## Uso OrganizeJr

A persona principal desta pasta é Lucas Camilo Carvalho, assessor comercial da OrganizeJr, empresa júnior de Psicologia Organizacional da UNEB Salvador.

O ciclo incremental atual adiciona uma camada comercial local para:

1. ler demandas de Lucas em Excel, texto e documentos;
2. considerar a pirâmide de captação e ICPs da OrganizeJr;
3. usar o relatório nacional de EJs do Portal BJ 2025 como base de leads e benchmark;
4. usar `pp-leads-brasil` para buscar, consultar e enriquecer empresas e EJs por CNPJ/nome/CNAE;
5. gerar ranking, justificativa de fit e abordagem personalizada.

## Story ativa

**external-enrichment-readiness**

Como Lucas, assessor comercial da OrganizeJr, quero que `pp-leads-brasil` consolide enriquecimento local, Casa dos Dados, `company-goat` e `contact-goat` para que a lista final de prospecção já saia com fontes, score, mensagem pronta e links de abordagem.

Critérios de aceite:

- `internal/client/casadados` consulta Casa dos Dados real quando houver chave e aceita tabela local via `PP_LEADS_LOCAL_TABLE` ou `PP_LEADS_USE_CASE_CONFIG`;
- `internal/client/pp` orquestra lead local, Casa dos Dados, `company-goat` e `contact-goat` sem hardcode de saída da OrganizeJr;
- `leads-brasil-pp-cli` expõe comandos `company-goat` e `contact-goat`;
- `enrich` retorna fontes usadas, `use_case`, `lead_context` e links `mailto:`/`wa.me` quando houver dados;
- script `scripts/run-use-case-enrichment.sh` roda `company`, `company-goat`, `contact-goat` e `enrich` por lead a partir de um `use-case.json`.
