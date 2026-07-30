# Project plan

## Projeto
pp-leads-brasil, stack público de pesquisa, enrichment e automação comercial. Dados e evidências OrganizeJr vivem no repo privado `luascfl/organizejr-pp-leads`.

## Estado macro
- m13, consolidação de duplicatas locais, concluído.
- m14, backfill cadastral prioritário, concluído.
- m15, expansão de contato com evidência para empresas não resolvidas, concluído.
- m16, sincronização CRM -> Perplexity Space, concluído.
- Pós-m16, split privado OrganizeJr e purge do histórico público, concluído.
- m17, sincronização operacional no Perplexity Space via navegador persistente, aprovado como próximo milestone.
- m18, geração de mensagens de abordagem por canal a partir dos snapshots, aprovado como milestone seguinte.
- m19, automação WhatsApp -> Perplexity -> CRM, aprovado como milestone posterior.

## Milestone ativo, m17

### Objetivo
Transformar o snapshot Markdown gerado no m16 em contexto realmente disponível no Perplexity Space da OrganizeJr, usando navegador persistente autenticado, sem tornar o Perplexity fonte canônica do CRM.

### Dependências
- Checkout privado `../organizejr-pp-leads` com `ploomes_crm/crm_sheet.py` e artefatos CRM.
- Snapshot m16 em `../organizejr-pp-leads/ploomes_crm/perplexity_space_m16_snapshot.md`.
- Chrome persistente/CDP autenticado no Perplexity.
- Google Sheets CRM permanece como fonte canônica operacional.

### Critérios de conclusão
- O Perplexity Space da OrganizeJr recebe o snapshot m16.
- A automação usa navegador persistente e não copia perfil completo do Chrome.
- O fluxo diferencia sugestão, mensagem enviada, resposta do lead, inferência e dado confirmado.
- Existe verificação observável de que o contexto foi colado/salvo no Space.
- O procedimento fica registrado para repetição por lote.

## Milestone seguinte, m18

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

## Milestone posterior, m19

### Objetivo
Automatizar o ciclo WhatsApp -> Perplexity -> CRM para conversas comerciais por lead e contato. O usuário deve poder pedir algo como "exportar toda a conversa de Ananda da Engetop" ou "houve novas mensagens em Ananda Engetop", e o sistema deve exportar a conversa completa ou delta, atualizar o Perplexity Space e produzir atualização revisável para a planilha.

### Escopo funcional
- Resolver contato e lead por linguagem natural, por exemplo Ananda + Engetop.
- Considerar primeiro a lista/filtro do WhatsApp chamada `Leads OrganizeJr` para localizar conversas comerciais, antes de cair para busca global.
- Exportar conversa completa do WhatsApp quando o usuário pedir histórico total.
- Detectar e exportar apenas novas mensagens quando o usuário indicar que houve atualização.
- Enviar transcript/delta para o Perplexity Space como contexto comercial do lead.
- Pedir ao Perplexity extração estruturada de fatos comerciais, próximos passos, objeções, status, canal e dados confirmados.
- Gerar payload revisável para `crm_sheet.py`, sem aplicar silenciosamente fatos ambíguos.

### Regra de fonte
A conversa WhatsApp é a evidência primária. O Perplexity organiza e extrai, mas não vira fonte final sem rastreabilidade. Atualizações na planilha exigem citação do trecho da conversa, confirmação explícita do usuário ou campo claramente declarado pela outra pessoa.

### Critérios esperados
- Estado por conversa para saber último export processado.
- Export completo e export incremental separados.
- Transcript salvo apenas no repo privado ou armazenamento local privado.
- Prompt do Perplexity diferencia mensagem enviada, resposta recebida, inferência, sugestão e dado confirmado.
- Payload de CRM marca origem como WhatsApp e inclui evidência/resumo.
