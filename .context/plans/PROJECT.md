# Project plan

## Projeto
pp-leads-brasil, stack público de pesquisa, enrichment e automação comercial. Dados e evidências OrganizeJr vivem no repo privado `luascfl/organizejr-pp-leads`.

## Estado macro
- m13, consolidação de duplicatas locais, concluído.
- m14, backfill cadastral prioritário, concluído.
- m15, expansão de contato com evidência para empresas não resolvidas, concluído.
- m16, sincronização CRM -> Perplexity Space, concluído.
- Pós-m16, split privado OrganizeJr e purge do histórico público, concluído.
- m17, sincronização operacional no Perplexity Space via instruções e arquivos, concluído.
- m18, geração de mensagens de abordagem por canal a partir dos snapshots, preparação concluída no repo privado.
- m19, integração CRM Sheet -> Perplexity -> WhatsApp -> CRM, planejada como próximo ciclo operacional.

## Milestone concluído, m17

### Objetivo
Transformar o snapshot Markdown gerado no m16 em contexto realmente disponível no Perplexity Space da OrganizeJr usando o fluxo confiável `Instruções` + `Arquivos`, sem tornar o Perplexity fonte canônica do CRM.

### Dependências
- Checkout privado `../organizejr-pp-leads` com `ploomes_crm/crm_sheet.py` e artefatos CRM.
- Snapshot m16 em `../organizejr-pp-leads/ploomes_crm/perplexity_space_m16_snapshot.md`.
- Chrome persistente/CDP autenticado no Perplexity.
- Perplexity Space operando apenas com `Instruções` e `Arquivos` como superfícies confiáveis.
- Google Sheets CRM permanece como fonte canônica operacional.

### Critérios de conclusão
- O Perplexity Space da OrganizeJr recebeu o snapshot m16 como arquivo.
- A automação usou navegador persistente e não copiou perfil completo do Chrome.
- As instruções salvas diferenciam sugestão, mensagem enviada, resposta do lead, inferência e dado confirmado.
- Existe verificação observável de que as instruções persistiram e os arquivos ficaram listados no Space.
- O procedimento registrado para repetição é atualizar instruções duráveis e subir arquivos Markdown.

## Milestone preparado, m18

### Objetivo
Gerar mensagens comerciais por canal a partir dos snapshots atualizados, preservando o CRM como fonte final e usando o Perplexity como apoio narrativo.

### Canais previstos
- WhatsApp.
- Instagram.
- E-mail.
- LinkedIn.

### Critérios esperados
- Mensagens curtas e acionáveis por lead/canal.
- Separação explícita entre sugestão, mensagem enviada, resposta do lead, inferência e dado confirmado.
- Nenhum campo cadastral do CRM é preenchido apenas por sugestão do Perplexity.

## Milestone planejado, m19

### Objetivo
Fechar a integração operacional CRM Sheet -> Perplexity -> WhatsApp -> CRM. Cada linha do CRM terá uma única coluna `Perplexity - tópico do lead`, apontando para uma única sessão/tópico principal por lead. O sistema deve exportar snapshots do CRM, preservar o link do tópico no snapshot, exportar conversas WhatsApp completas ou incrementais, atualizar o tópico do lead no Perplexity e produzir atualização revisável para a planilha.

### Escopo funcional
- Garantir schema do CRM com a coluna única `Perplexity - tópico do lead`.
- Implementar leitura/escrita revisável desse link no `crm_sheet.py`, sem sobrescrever valor existente sem confirmação explícita.
- Incluir o link do tópico Perplexity nos snapshots exportados para o Space.
- Resolver contato e lead por linguagem natural, por exemplo Ananda + Engetop.
- Considerar primeiro a lista/filtro do WhatsApp chamada `Leads OrganizeJr` para localizar conversas comerciais, antes de cair para busca global.
- Exportar conversa completa do WhatsApp quando o usuário pedir histórico total.
- Detectar e exportar apenas novas mensagens quando o usuário indicar que houve atualização.
- Enviar transcript/delta para o tópico único do lead no Perplexity Space como contexto comercial atualizado.
- Pedir ao Perplexity extração estruturada de fatos comerciais, próximos passos, objeções, status, canal e dados confirmados.
- Gerar payload revisável para `crm_sheet.py`, sem aplicar silenciosamente fatos ambíguos.

### Regra de fonte
A conversa WhatsApp é a evidência primária. O Perplexity organiza e extrai, mas não vira fonte final sem rastreabilidade. Atualizações na planilha exigem citação do trecho da conversa, confirmação explícita do usuário ou campo claramente declarado pela outra pessoa.

### Critérios esperados
- Estado por conversa para saber último export processado.
- Export completo e export incremental separados.
- Transcript salvo apenas no repo privado ou armazenamento local privado.
- Prompt do Perplexity diferencia mensagem enviada, resposta recebida, inferência, sugestão e dado confirmado.
- Payload de CRM marca origem como WhatsApp, inclui evidência/resumo e preserva o link `Perplexity - tópico do lead`.

### Plano de execução m19

1. **Schema CRM**
   - Verificar se a planilha tem a coluna `Perplexity - tópico do lead`.
   - Criar comando seguro para detectar coluna ausente e orientar criação.
   - Só escrever link quando o usuário passar URL explícita ou confirmar criação do tópico.

2. **Comandos `crm_sheet.py`**
   - `perplexity-topic get --row N` para ver o link atual.
   - `perplexity-topic set --row N --url ...` para registrar o tópico principal do lead.
   - `perplexity-topic missing` para listar leads sem tópico.
   - `export-perplexity-snapshots` deve incluir a coluna quando existir.

### Status m19-us001
- Concluído no repo privado `organizejr-pp-leads`, commit `1165948`.
- Coluna `Perplexity - tópico do lead` criada na aba `Clientes`, coluna `AQ`.
- `crm_sheet.py perplexity-topic get|set|missing|ensure-column` implementado.
- `export-perplexity-snapshots` inclui `Perplexity tópico` em cada lead.
- `apply-icp-score` aceita `perplexity_topic_url` e `perplexity_topic` sem sobrescrever célula existente em leads já cadastrados.

3. **Regra de tópico único**
   - Uma sessão/tópico principal por lead.
   - Nova pesquisa, WhatsApp, follow-up ou proposta entra no mesmo tópico.
   - O tópico não confirma dado cadastral sozinho.

4. **WhatsApp -> Perplexity**
   - Localizar conversa primeiro na lista/filtro `Leads OrganizeJr`.
   - Exportar histórico completo ou delta novo.
   - Salvar transcript/delta apenas no repo privado.
   - Usar o link do CRM para abrir o tópico do lead e anexar/colar o novo contexto.

### Status m19-us002
- Concluído no repo privado `organizejr-pp-leads`, commit `47a1be3`.
- Novo script privado: `ploomes_crm/whatsapp_export.py`.
- Comandos disponíveis: `resolve`, `sync-status`, `export --dry-run`, `export`.
- Preflight obrigatório validado: `sync-status` retornou `ok=true`, `status=complete_no_sync_card`, `percent=100`.
- Export real concluído para Ananda/ENGETOP: `used_preferred_filter=true`, `messages=10`.
- Transcript privado salvo em `../organizejr-pp-leads/leads/engetop/whatsapp/ananda/full-2026-07-30.md`.

5. **Perplexity -> CRM revisável**
   - Gerar payload com fatos, inferências, resposta do lead, mensagem enviada, evidência textual e próximo passo.
   - Nunca aplicar fato ambíguo automaticamente.
   - Atualizar planilha só com evidência rastreável ou confirmação explícita.

### Status m19-us003
- Concluído no repo privado `organizejr-pp-leads`, commit `4ed62e2`.
- Novo comando read-only: `crm_sheet.py crm-review-payload-from-whatsapp`.
- O comando lê transcript WhatsApp privado e gera JSON revisável com `Status Comercial`, `Observações`, `Próximo passo`, campos proibidos por falta de evidência e warnings de revisão.
- Busca interna do Perplexity Space confirmou o tópico ENGETOP em `https://www.perplexity.ai/search/40cee0f7-3e0b-4842-ae03-e9a33f4abcb2` e o link foi salvo na linha 2 do CRM.
- Validação real com Ananda/ENGETOP regenerou `leads/engetop/whatsapp/ananda/crm-review-payload-2026-07-30.json` com `perplexity_topic_url`, sem aplicar atualização automática na planilha.

