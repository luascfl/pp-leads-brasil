# State

## Atualizado em
2026-07-30T08:55:00-03:00

## Estado atual
- m15 concluído sem input pendente do usuário.
- m16 concluído: sincronização CRM -> Perplexity Space.
- Story Ralph `m16-us001` concluída: `crm_sheet.py export-perplexity-snapshots` gera Markdown colável no Perplexity Space.

## Decisão de arquitetura
O Google Sheets CRM continua sendo a fonte canônica operacional. O Perplexity Space será usado como dump narrativo compartilhado para conversas presenciais, WhatsApp, abordagens, status inferido e sugestões de projetos aplicáveis. Snapshots do CRM alimentam o Space com o placar atual para reduzir sugestões desatualizadas.

## DoD da story atual
- Comando read-only em `organizejr-pp-leads/ploomes_crm/crm_sheet.py`.
- Export por linhas específicas ou lote limitado.
- Markdown com status, saúde, ICP, decisor, contatos, canais, CNPJ, lacunas, observações e log técnico.
- Validação por sintaxe, help e execução real.

## Evidência final
- `python3 -m py_compile organizejr-pp-leads/ploomes_crm/crm_sheet.py` passou.
- `python3 organizejr-pp-leads/ploomes_crm/crm_sheet.py export-perplexity-snapshots --help` passou e lista o comando.
- `python3 organizejr-pp-leads/ploomes_crm/crm_sheet.py export-perplexity-snapshots --rows 3,4,5,6,105,106 --output organizejr-pp-leads/ploomes_crm/perplexity_space_m16_snapshot.md` gerou 6 snapshots.
- Verificação do Markdown: `lead_sections=6`, `instruction_count=1`.

## Artefato gerado
- `organizejr-pp-leads/ploomes_crm/perplexity_space_m16_snapshot.md`

## Próximo passo
Escolher m17. Opções prováveis: automação de colagem/export no Perplexity Space via navegador persistente, ou geração de mensagens de abordagem por canal a partir dos snapshots m16.
