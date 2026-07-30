# State

## Atualizado em
2026-07-30T10:42:00-03:00

## Estado atual
- m15 concluído sem input pendente do usuário.
- m16 concluído: sincronização CRM -> Perplexity Space.
- Story Ralph `m16-us001` concluída: `crm_sheet.py export-perplexity-snapshots` gera Markdown colável no Perplexity Space.
- Fechamento de privacidade pós-m16 concluído: `organizejr-pp-leads` foi separado para repo privado e removido do histórico público de `pp-leads-brasil`.
- `scripts/` foi removido do repo público.
- m17 concluído: Perplexity Space OrganizeJr atualizado pelo fluxo confiável `Instruções` + `Arquivos`. Computer, Brain, Skills, conectores, Slack/Teams e links não são tratados como superfícies confiáveis.
- m18 preparação concluída no repo privado: rascunhos de abordagem por canal gerados em `../organizejr-pp-leads/ploomes_crm/outreach_m18_drafts.md`.
- m19 preparação concluída no repo privado: requisitos WhatsApp -> Perplexity -> CRM registrados em `../organizejr-pp-leads/ploomes_crm/whatsapp_m19_requirements.md`.

## Decisão de arquitetura
O Google Sheets CRM continua sendo a fonte canônica operacional. O Perplexity Space será usado como dump narrativo compartilhado para conversas presenciais, WhatsApp, abordagens, status inferido e sugestões de projetos aplicáveis. Snapshots do CRM alimentam o Space com o placar atual para reduzir sugestões desatualizadas.

`organizejr-pp-leads` agora é repositório privado separado em `https://github.com/luascfl/organizejr-pp-leads` e checkout local em `../organizejr-pp-leads`. O repo público `pp-leads-brasil` não deve voltar a receber evidências, planilhas, CRM local, `.clasp.json`, `.env` ou dados comerciais OrganizeJr.

Para automações de WhatsApp, a conversa exportada é a evidência primária. O Perplexity pode organizar, extrair e sugerir atualizações, mas a planilha só deve receber fatos com proveniência clara na conversa ou revisão explícita do usuário. Isso evita transformar inferência do Perplexity em dado cadastral confirmado.

O Space da OrganizeJr deve ser tratado como copiloto fundamentado em documentos. O usuário confirmou que só `Instruções` e `Arquivos` funcionam de forma confiável; portanto, o fluxo não depende de Computer, Brain, Skills, Slack/Teams, conectores ou links.

Schema decidido para o CRM: uma única coluna chamada `Perplexity - tópico do lead`, com um único tópico/sessão principal do Perplexity por lead. A coluna é ponte de navegação/evidência narrativa e não torna o Perplexity fonte canônica cadastral.

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

## Evidência final m17
- Chrome persistente/CDP abriu o Perplexity autenticado.
- Cloudflare foi resolvido manualmente pelo usuário.
- Instruções do Space foram substituídas pelo texto validado e persistiram após reload com 3053 caracteres.
- Arquivos enviados ao Space: `perplexity_space_m16_snapshot.md`, `outreach_m18_drafts.md`, `perplexity_space_instructions_m17_candidate.md`, `whatsapp_m19_requirements.md`.
- A aba Arquivos listou `12 itens` e os quatro arquivos novos.
- Teste por sessão não foi usado como critério final porque a UI caiu em Computer e bloqueio de créditos, reforçando que o fluxo confiável é `Instruções` + `Arquivos`.

## Artefatos gerados
- `../organizejr-pp-leads/ploomes_crm/perplexity_space_m16_snapshot.md`
- `../organizejr-pp-leads/ploomes_crm/perplexity_space_instructions_m17_candidate.md`
- `../organizejr-pp-leads/ploomes_crm/perplexity_space_m17_evidence.md`
- `../organizejr-pp-leads/ploomes_crm/outreach_m18_drafts.md`
- `../organizejr-pp-leads/ploomes_crm/whatsapp_m19_requirements.md`

## Próximo passo
Executar m18 ou m19 conforme prioridade operacional: usar `outreach_m18_drafts.md` para disparos comerciais revisados, ou iniciar implementação do export WhatsApp completo/incremental seguindo `whatsapp_m19_requirements.md`.
