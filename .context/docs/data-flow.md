---
type: doc
name: data-flow
description: How data moves through the system and external integrations
category: data-flow
generated: 2026-05-25
updated: 2026-06-17
status: filled
scaffoldVersion: "2.0.0"
---

# Data flow

## Fluxo OrganizeJr

1. Entrada: Lucas adiciona uma demanda comercial em Excel, CSV, texto, DOCX ou Markdown.
2. Contexto: o agente lê `.context/docs/organizejr-commercial-context.md`, a pirâmide de captação e os arquivos citados na demanda.
3. Base EJ: quando a demanda envolve empresas juniores, o agente consulta `Relatório Portal BJ EJs 2025 0101.xlsx` e usa campos como CNPJ, cursos, cidade, federação, membros, faturamento, contratos e cluster.
4. Busca externa: quando o lead é empresa, contabilidade, escola ou instituição fora da base BJ, o agente usa `leads-brasil-pp-cli leads-brasil-platform-search` por nome ou CNAE.
5. Enriquecimento: leads com CNPJ e fit A/B podem ser consultados com `company` e enriquecidos com `enrich`.
6. Qualificação: o agente aplica score por aderência ao ICP, evidência de demanda, acessibilidade, capacidade de pagar/trocar valor e encaixe de abordagem.
7. Saída: ranking, tabela operacional, evidências, mensagens de abordagem e pendências ficam no chat ou em `.context/organizejr/` quando houver arquivo reutilizável.

O fluxo não deve inventar contatos ou dores. Toda saída precisa apontar a fonte local, comando `pp-leads` ou limitação.
