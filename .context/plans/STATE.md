# State

## Atualizado em
2026-07-31

## Estado atual
- m15 concluído sem input pendente do usuário.
- m16 concluído: sincronização CRM -> Perplexity Space.
- Story Ralph `m16-us001` concluída: `crm_sheet.py export-perplexity-snapshots` gera Markdown colável no Perplexity Space.
- Fechamento de privacidade pós-m16 concluído: `organizejr-pp-leads` foi separado para repo privado e removido do histórico público de `pp-leads-brasil`.
- `scripts/` foi removido do repo público.
- m17 concluído: Perplexity Space OrganizeJr atualizado pelo fluxo confiável `Instruções` + `Arquivos`. Computer, Brain, Skills, conectores, Slack/Teams e links não são tratados como superfícies confiáveis.
- m18 preparação concluída no repo privado: rascunhos de abordagem por canal gerados em `../organizejr-pp-leads/ploomes_crm/outreach_m18_drafts.md`.
- m19-us001 concluído no repo privado: coluna única `Perplexity - tópico do lead`, comandos `perplexity-topic` e snapshots com link do tópico.
- m19-us002 concluído no repo privado: `whatsapp_export.py` abre primeiro a lista `Leads OrganizeJr`, exportou a conversa completa de Ananda Ferreira | ENGETOP e salvou transcript apenas no repo privado.
- m19-us003 concluído no repo privado: `crm_sheet.py crm-review-payload-from-whatsapp` gera JSON revisável a partir do transcript WhatsApp privado. O tópico ENGETOP foi localizado pela busca interna do Perplexity Space, salvo na linha 2 do CRM e incluído no payload revisável. Ajuste de schema posterior migrou `Status Comercial` para `Funil`, com dropdown estrito e payload revisável usando `Funil`, `Observações` e `Próximo passo`. ENGETOP foi aplicada no CRM com `Funil = Contato iniciado`, `Próximo passo = follow-up para agendar reunião diagnóstica` e `Observações` sem resumo comercial redundante.
- m19 atualização operacional: `whatsapp_export.py` foi endurecido para não usar busca global do WhatsApp em `export`/`draft`; o fluxo abre a setinha ao lado de `Favoritas`, seleciona `Leads OrganizeJr` e varre/rola apenas essa lista. A conversa de Gabriella ENGETOP foi exportada no repo privado. O upload equivocado do transcript na aba `Arquivos` do Perplexity Space foi removido; transcripts WhatsApp devem ser colados integralmente no tópico Perplexity do lead, não enviados como arquivos globais do Space.

## Decisão de arquitetura
O Google Sheets CRM continua sendo a fonte canônica operacional. O Perplexity Space será usado como dump narrativo compartilhado para conversas presenciais, WhatsApp, abordagens, status inferido e sugestões de projetos aplicáveis. Snapshots do CRM alimentam o Space com o placar atual para reduzir sugestões desatualizadas.

`organizejr-pp-leads` agora é repositório privado separado em `https://github.com/luascfl/organizejr-pp-leads` e checkout local em `../organizejr-pp-leads`. O repo público `pp-leads-brasil` não deve voltar a receber evidências, planilhas, CRM local, `.clasp.json`, `.env` ou dados comerciais OrganizeJr.

Para automações de WhatsApp, a conversa exportada é a evidência primária. O Perplexity pode organizar, extrair e sugerir atualizações, mas a planilha só deve receber fatos com proveniência clara na conversa ou revisão explícita do usuário. Isso evita transformar inferência do Perplexity em dado cadastral confirmado.

O Space da OrganizeJr deve ser tratado como copiloto fundamentado em documentos. O usuário confirmou que só `Instruções` e `Arquivos` funcionam de forma confiável; portanto, o fluxo não depende de Computer, Brain, Skills, Slack/Teams, conectores ou links.

Schema decidido para o CRM: uma única coluna chamada `Perplexity - tópico do lead`, com um único tópico/sessão principal do Perplexity por lead. A coluna é ponte de navegação/evidência narrativa e não torna o Perplexity fonte canônica cadastral.

Schema decidido para etapa comercial: a coluna antiga `Status Comercial` foi renomeada para `Funil` no Google Sheets CRM, com dropdown canônico (`Pesquisa inicial`, `Qualificado`, `Abordar agora`, `Benchmark`, `Nutrir depois`, `Contato iniciado`, `Interesse sinalizado`, `Reunião/diagnóstico`, `Proposta enviada`, `Aguardando decisão`, `Ganho`, `Perdido`, `Duplicata consolidada`). `Próximo passo` permanece separado da etapa do funil.

Regra operacional adicionada: nunca localizar tópico Perplexity por histórico do Chrome, cache, local storage ou banco local do navegador. O link da coluna `Perplexity - tópico do lead` só pode ser salvo depois de abrir o resultado pela busca interna do Perplexity Space e observar a URL final no navegador.

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

## Evidência final m19-us001
- Repo privado `organizejr-pp-leads`: commit `1165948 feat(crm): track perplexity topic per lead` pushado em `main`.
- Google Sheets CRM: coluna `Perplexity - tópico do lead` criada na aba `Clientes`, coluna `AQ`.
- `python3 -m py_compile ../organizejr-pp-leads/ploomes_crm/crm_sheet.py` passou.
- `perplexity-topic ensure-column` confirma a coluna em `AQ`.
- `perplexity-topic missing --limit 1` lista leads sem tópico principal.
- `perplexity-topic get --row 3` retorna Consultagro com tópico vazio.
- `export-perplexity-snapshots --rows 3 --obs-limit 80` inclui `Perplexity tópico: sem tópico registrado`.

## Evidência final m19-us002
- Repo privado `organizejr-pp-leads`: commit `47a1be3 fix(crm): use whatsapp leads filter before export` pushado em `main`.
- `python3 -m py_compile ploomes_crm/whatsapp_export.py` passou.
- `whatsapp_export.py sync-status` retornou `ok=true`, `status=complete_no_sync_card`, `percent=100`.
- `whatsapp_export.py export --instruction "exportar toda a conversa de Ananda da Engetop"` retornou `ok=true`, `used_preferred_filter=true`, `messages=10`.
- Transcript privado salvo em `../organizejr-pp-leads/leads/engetop/whatsapp/ananda/full-2026-07-30.md`, conversation key `whatsapp:engetop:ananda-ferreira-engetop`, SHA-256 `07b948e5c23dc96153e56a7fe1dad158686d361d93ebf3e02e1a499727ed1387`.

## Evidência final m19-us003
- Repo privado `organizejr-pp-leads`: commit `4ed62e2 feat(crm): generate whatsapp review payload` pushado em `main`.
- Repo privado `organizejr-pp-leads`: commit `6568787 feat(crm): migrate commercial status to funil` pushado em `main`.
- Repo privado `organizejr-pp-leads`: commit `c4820ff feat(crm): support next step updates` pushado em `main`.
- `python3 -m py_compile ploomes_crm/crm_sheet.py` passou.
- `crm-review-payload-from-whatsapp --help` lista `--input`, `--perplexity-topic-url` e `--output`.
- Busca interna do Perplexity Space abriu o tópico ENGETOP em `https://www.perplexity.ai/search/40cee0f7-3e0b-4842-ae03-e9a33f4abcb2`; título observado: `Engetop vou falar com ela e sobre ela nas próximas mensagens, estou no EDL Bahia`.
- `perplexity-topic set --row 2 --url https://www.perplexity.ai/search/40cee0f7-3e0b-4842-ae03-e9a33f4abcb2` salvou o link no Google Sheets CRM.
- `perplexity-topic get --row 2` confirmou `lead=ENGETOP` e o mesmo `perplexity_topic_url`.
- `funil-schema --yes` renomeou `Status Comercial` para `Funil` na coluna `AP`, aplicou dropdown estrito e normalizou 10 valores.
- `analyze-health` rodou com a nova coluna `Funil` e reaplicou os dropdowns de `ICP Considerado` e `Funil`.
- `crm-review-payload-from-whatsapp --input leads/engetop/whatsapp/ananda/full-2026-07-30.md --perplexity-topic-url https://www.perplexity.ai/search/40cee0f7-3e0b-4842-ae03-e9a33f4abcb2 --output leads/engetop/whatsapp/ananda/crm-review-payload-2026-07-30.json` gerou payload com `row=2`, `lead_slug=engetop`, `audit.message_count=10`, updates `Funil`, `Observações` e `Próximo passo`, todos com `requires_review=true`.
- Google Sheets CRM linha 2 ENGETOP: coluna `AI` renomeada de `Prontidão comercial` para `Próximo passo`; `Funil` atualizado para `Contato iniciado`; `Próximo passo` atualizado para follow-up com Ananda para agendar reunião diagnóstica; `Observações` recebeu apenas registro contextual, sem repetir a ação.
- Atualização ENGETOP/Gabriella: `whatsapp_export.py export --lead ENGETOP --contact Gabriella --mode full` retornou `ok=true`, `used_preferred_filter=true`, `messages=51`, `conversation_key=whatsapp:engetop:gabriella-engetop`. Transcript privado: `../organizejr-pp-leads/leads/engetop/whatsapp/gabriella/full-2026-07-30.md`. O upload equivocado de `full-2026-07-30.md` na aba `Arquivos` do Perplexity Space foi removido; após a remoção, a UI voltou a `Files`, `12 itens`, sem listar o arquivo. O transcript foi colado no tópico ENGETOP em 2 partes, e a resposta final observada no mesmo URL retornou mensagens enviadas, respostas recebidas, próximos passos, objeções, campos para CRM, inferências e sugestão.

## Artefatos gerados
- `../organizejr-pp-leads/ploomes_crm/perplexity_space_m16_snapshot.md`
- `../organizejr-pp-leads/ploomes_crm/perplexity_space_instructions_m17_candidate.md`
- `../organizejr-pp-leads/ploomes_crm/perplexity_space_m17_evidence.md`
- `../organizejr-pp-leads/ploomes_crm/outreach_m18_drafts.md`
- `../organizejr-pp-leads/ploomes_crm/whatsapp_m19_requirements.md`
- `../organizejr-pp-leads/leads/engetop/whatsapp/gabriella/full-2026-07-30.md`
- `../organizejr-pp-leads/ploomes_crm/whatsapp_gabriella_perplexity_evidence.md`

Próximo passo m19: usar o Perplexity Space para revisar conjuntamente as conversas de Ananda e Gabriella antes de sugerir qualquer mudança adicional no CRM da ENGETOP; manter a regra de que automações WhatsApp para leads comerciais usam listas/filtros, nunca busca global.

## Milestone m20

- `m20-us001` concluída: a CLI pública não descobre perfil OrganizeJr; testes usam fixtures temporárias públicas e perfil externo só é carregado por `PP_LEADS_USE_CASE_CONFIG` ou `PP_LEADS_ICP_DIR`.
- Enriquecimento social agora sempre tenta Scrape Creators. Em falha de billing, credencial, timeout ou ausência de dados utilizáveis, usa Apify apenas se `APIFY_TOKEN` e `APIFY_SOCIAL_ACTOR_ID` estiverem configurados. Falhas de ambos permanecem não bloqueantes e são registradas no payload.
- Verificação: `go test ./...` passou, incluindo o contrato de fallback Scrape Creators → Apify. A busca pública não encontrou descoberta implícita de perfil privado. Os documentos públicos de overview, tooling e workflow foram limpos de comandos e caminhos privados. `graphify update .` reextraiu 721 arquivos de código.
- Os legados `.context/docs/planning_gsd` e `.context/prd_ralph` continuam inventariados e não foram removidos por falta de decisão explícita de migração.

## Milestone m21

- `m21-us001` concluída: `internal/operation` implementa planos imutáveis com digest SHA-256, expiração, alvos com evidência e recibos duráveis por alvo. Aplicação exige aprovação explícita e é idempotente para alvos já aplicados.
- `POST /v1/operations/plan` só cria plano em perfil explícito; `POST /v1/operations/{plan_id}/apply` processa apenas o plano salvo. Sem adapter privado configurado, a aplicação registra falha por alvo, sem efeito externo.
- A CLI expõe `operation plan --input operation.json` e `operation apply <plan-id> --yes`; `--agent` isolado não autoriza mutação.
- Verificação: `go test ./...` e `cd leads-brasil-pp-cli && go test ./...` passaram. `graphify update .` reextraiu 726 arquivos de código.

## Milestone m22

- `m22-us001` em execução: o núcleo público aceita um `operation_adapter_command` declarado pelo perfil e entrega apenas o alvo imutável via stdin. Não transporta credenciais nem implementa Ploomes.
- O perfil privado OrganizeJr agora contém `crm_schema.json`, um construtor de planos somente a partir de snapshots de linhas explicitamente selecionadas e um adapter que exige `--yes`, revalida o valor remoto antes do `PATCH` e grava ledger local privado.
- Pendente para conclusão: conectar o leitor real de Sheets/Ploomes aos snapshots, verificar/remover gatilhos Apps Script e executar uma alteração de teste aprovada com reconciliação read-only.
