# State

## Atualizado em
2026-07-30T10:04:45-03:00

## Estado atual
- m15 concluído sem input pendente do usuário.
- m16 concluído: sincronização CRM -> Perplexity Space.
- Story Ralph `m16-us001` concluída: `crm_sheet.py export-perplexity-snapshots` gera Markdown colável no Perplexity Space.
- Fechamento de privacidade pós-m16 concluído: `organizejr-pp-leads` foi separado para repo privado e removido do histórico público de `pp-leads-brasil`.
- `scripts/` foi removido do repo público.

## Decisão de arquitetura
O Google Sheets CRM continua sendo a fonte canônica operacional. O Perplexity Space será usado como dump narrativo compartilhado para conversas presenciais, WhatsApp, abordagens, status inferido e sugestões de projetos aplicáveis. Snapshots do CRM alimentam o Space com o placar atual para reduzir sugestões desatualizadas.

`organizejr-pp-leads` agora é repositório privado separado em `https://github.com/luascfl/organizejr-pp-leads` e checkout local em `../organizejr-pp-leads`. O repo público `pp-leads-brasil` não deve voltar a receber evidências, planilhas, CRM local, `.clasp.json`, `.env` ou dados comerciais OrganizeJr.

## DoD da story m16
- Comando read-only em `../organizejr-pp-leads/ploomes_crm/crm_sheet.py`.
- Export por linhas específicas ou lote limitado.
- Markdown com status, saúde, ICP, decisor, contatos, canais, CNPJ, lacunas, observações e log técnico.
- Validação por sintaxe, help e execução real.

## Evidência final m16
- `python3 -m py_compile ../organizejr-pp-leads/ploomes_crm/crm_sheet.py` passou antes do split.
- `python3 ../organizejr-pp-leads/ploomes_crm/crm_sheet.py export-perplexity-snapshots --help` passou antes do split e lista o comando.
- `python3 ../organizejr-pp-leads/ploomes_crm/crm_sheet.py export-perplexity-snapshots --rows 3,4,5,6,105,106 --output ../organizejr-pp-leads/ploomes_crm/perplexity_space_m16_snapshot.md` gerou 6 snapshots antes do split.
- Verificação do Markdown: `lead_sections=6`, `instruction_count=1`.
- Repo privado `luascfl/organizejr-pp-leads` criado e verificado como `PRIVATE`.
- Repo público `luascfl/pp-leads-brasil` purgado: `tree_path_hits=0`, `history_path_hits=0` para `organizejr-pp-leads/`.

## Artefato gerado
- `../organizejr-pp-leads/ploomes_crm/perplexity_space_m16_snapshot.md`

## Próximo passo
m17 aprovado: sincronização operacional CRM -> Perplexity Space via navegador persistente. Escopo: abrir o Space da OrganizeJr no Perplexity com Chrome/CDP persistente, colar snapshot m16, validar recebimento do contexto e registrar fluxo repetível.

## Passo seguinte
m18 aprovado: geração de mensagens de abordagem por canal a partir dos snapshots m16/m17. Escopo provável: WhatsApp, Instagram, e-mail e LinkedIn, preservando diferença entre sugestão, mensagem enviada, resposta do lead, inferência e dado confirmado.
