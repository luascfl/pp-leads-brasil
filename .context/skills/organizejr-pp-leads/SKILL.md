---
type: skill
name: organizejr-pp-leads
skillSlug: organizejr-pp-leads
description: Use esta skill sempre que Lucas Camilo Carvalho, a OrganizeJr, demandas comerciais, planilhas de leads, EJs, empresas juniores, ICPs, abordagem comercial, enriquecimento de contatos ou pp-leads-brasil forem mencionados nesta pasta. Ela orienta como transformar arquivos Excel/texto/docx em pesquisa, enriquecimento, priorização e abordagem de leads para a OrganizeJr, empresa júnior de Psicologia Organizacional da UNEB Salvador.
phases: [P, E, V, C]
generated: 2026-06-17
status: filled
scaffoldVersion: "2.0.0"
---

# OrganizeJr PP Leads

## Objetivo

Use esta skill para operar o `pp-leads-brasil` no contexto comercial da OrganizeJr. O resultado esperado é uma lista de leads qualificados, enriquecidos quando possível, com justificativa de fit e abordagem pronta para Lucas usar.

## Contexto obrigatório

Antes de ranquear ou abordar leads, leia:

1. `.context/docs/organizejr-commercial-context.md`
2. `Pirâmide de captação e ICPs da OrganizeJr.docx`
3. `Relatório Portal BJ EJs 2025 0101.xlsx`, quando a demanda envolver EJs
4. Arquivos de demanda citados pelo usuário ou colocados na pasta, especialmente `.xlsx`, `.csv`, `.docx`, `.txt` e `.md`
5. `leads-brasil-pp-cli/SKILL.md`, se precisar executar ou explicar comandos do pp-leads

Não use memória solta para ICP, nomes, números ou contatos. Reabra a fonte local e cite o arquivo usado.

## Fluxo de trabalho

### 1. Entender a demanda

Extraia dos arquivos ou do pedido:

- segmento desejado: instituições de ensino, EJs/startups, mercado sênior ou outro recorte explícito
- geografia: Salvador, Bahia, Nordeste, Brasil ou lista indicada por Lucas
- quantidade esperada de leads
- objetivo da abordagem: parceria, diagnóstico, capacitação, oficina, indicação, pesquisa, benchmark ou venda direta
- restrições: excluir concorrentes, priorizar contatos institucionais, evitar leads sem CNPJ, usar apenas EJs federadas, usar apenas base local

Se a demanda estiver incompleta, siga o padrão conservador: priorizar Salvador e Bahia, contato institucional, fit alto, baixo risco e próximo passo simples.

### 2. Escolher a fonte certa

- EJs: comece pelo `Relatório Portal BJ EJs 2025 0101.xlsx`.
- EJs de comunicação nacional: filtre o relatório BJ por `EMPRESA_JUNIOR` e `CURSOS` usando termos como Comunicação, Jornalismo, Publicidade, Propaganda, Relações Públicas, Audiovisual, Cinema, Rádio, Mídias, Marketing, Design Gráfico e Produção Cultural. Trate como `ICP 2.3 — EJs de Comunicação com Sobrecarga Criativa e Relacional`, hipótese a validar. Não classifique como alta demanda operacional genérica; procure sinais de fadiga criativa, trabalho always-on, aprovação subjetiva, retrabalho, conflito entre criação/atendimento/estratégia, onboarding frágil, insegurança de feedback, exposição em redes, crise/moderação e perda de padrão criativo por troca de gestão. Use Acesso Comunicação Jr. como case candidato e ESPM Jr. como benchmark, nunca como prova fechada de dor sem validação comercial.
- Empresas e contabilidades: use `pp-leads-brasil` para buscar por nome, CNPJ ou CNAE.
- Escolas, cursinhos e instituições de capacitação: use demanda local quando existir; complemente com `pp-leads-brasil` por nome/CNAE e busca web apenas para validar informação pública.
- Quando faltarem sinais públicos de presença social, consistência de marca ou links ativos, complemente com `scrape-creators-pp-cli`, principalmente para Instagram, LinkedIn, TikTok e YouTube.
- Outras EJs de Psicologia: classifique como benchmark, parceria ou inteligência competitiva, não como lead de venda direta por padrão.

### 3. Rodar pp-leads com segurança

Primeiro descubra o estado real da CLI:

```bash
leads-brasil-pp-cli doctor --json
leads-brasil-pp-cli agent-context --pretty
```

Use sempre `--agent` nas execuções de dados. Ele força JSON compacto, sem cor e sem interação.

Quando `doctor` retornar `auth: not configured`, diferencie backend local e API externa. Se `base_url` for `localhost`/`127.0.0.1`, use `LEADS_BRASIL_BEARER_AUTH=local-dev-token` para desenvolvimento local, porque o backend local deste projeto não valida token real. Se `base_url` for API externa, rode auth interativo, peça o token no terminal, exporte `LEADS_BRASIL_BEARER_AUTH` para a sessão e ofereça salvar com `leads-brasil-pp-cli auth set-token`. Não use `--dry-run` em entregas operacionais; use `--dry-run` só em diagnóstico explícito.

Se `doctor` retornar `api: unreachable` com `base_url` local, tente iniciar `server_bin` ou `go run ./cmd/server`. Se retornar unreachable com API externa, corrija `base_url` em `~/.config/leads-brasil-pp-cli/config.toml`. O backend local agora lê a tabela a partir do caso de uso configurado por `PP_LEADS_USE_CASE_CONFIG` ou `PP_LEADS_LOCAL_TABLE`; para a OrganizeJr, use `.context/use-cases/organizejr-ejs-comunicacao/use-case.json`. `company` retorna dados locais, `use_case`, `lead_context` e tenta anexar Casa dos Dados quando `CASA_DADOS_API_KEY` ou `PP_LEADS_CASA_DADOS_API_KEY` estiver configurada; `enrich` retorna o payload genérico do caso de uso, dados Casa dos Dados disponíveis e links `mailto:`/`wa.me` quando houver e-mail ou telefone.

Comandos principais:

```bash
leads-brasil-pp-cli leads-brasil-platform-search --name "contabilidade Salvador" --agent --select name,cnpj,city,state,status
leads-brasil-pp-cli leads-brasil-platform-search --cnae "6920-6/01" --agent --select name,cnpj,city,state,status
leads-brasil-pp-cli company "34.434.241/0001-60" --agent --select name,cnpj,status,address,contacts,casadosdados,sources,use_case,lead_context
leads-brasil-pp-cli company-goat "34.434.241/0001-60" --agent
leads-brasil-pp-cli contact-goat "Nome da EJ ou empresa" --agent
leads-brasil-pp-cli enrich "34.434.241/0001-60" --agent --deliver file:.context/use-cases/organizejr-ejs-comunicacao/outputs/enrichment/34-434-241-0001-60.json
```

Exemplos úteis com `scrape-creators-pp-cli`:

```bash
scrape-creators-pp-cli doctor
scrape-creators-pp-cli instagram list-profile "perfil_da_ej" --agent
scrape-creators-pp-cli linkedin list-company "nome-da-ej" --agent
scrape-creators-pp-cli bio resolve "https://linktr.ee/exemplo" --agent
```

Antes de usar `contact-goat` com fontes externas reais, rode:

```bash
.context/use-cases/organizejr-ejs-comunicacao/setup-auth.sh
```

Ou, sem wrapper do caso:

```bash
scripts/setup-external-enrichment-auth.sh
```

Esse setup cobre:

### 4. Qualificar pelo ICP da OrganizeJr

Use a pontuação de `.context/organizejr/metodologia-score-leads-organizejr.md` e o resumo de `.context/docs/organizejr-commercial-context.md`:

- 30, aderência ao ICP
- 20, evidência de demanda
- 20, acessibilidade
- 15, capacidade de pagar ou trocar valor
- 15, encaixe de abordagem

Classes:

- A: abordar agora, score 75 ou mais
- B: nutrir, score 55 a 74
- C: monitorar, score abaixo de 55
- Benchmark: referência, parceria ou inteligência competitiva

Sempre explique o score em uma frase curta com evidência da fonte.

### 5. Preparar abordagem

Para cada lead A ou B, gere uma abordagem curta com:

1. evidência objetiva, por exemplo cidade, curso, contratos, crescimento, porte, CNPJ, site, dor descrita ou relação com UNIJr-BA
2. hipótese de dor ligada a Psicologia Organizacional, por exemplo engajamento, comunicação, liderança, cultura, gestão de tempo, carreira, permanência ou RH sem estrutura
3. oferta da OrganizeJr, por exemplo diagnóstico, capacitação, oficina, roda de conversa, trilha, piloto ou conversa de 20 minutos
4. próximo passo claro

Tom por camada:

- Base: educativo, institucional e colaborativo.
- Meio: de par para par, prático e ligado à rotina de gestão.
- Topo: consultivo, com prova social e risco baixo.

## Formato de saída

Entregue uma tabela principal e um plano curto.

Colunas mínimas:

| lead | camada | ICP | cidade/UF | CNPJ | contato | fonte | evidência | score | classe | oferta sugerida | primeira mensagem | próximo passo | risco/pendência |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |

Depois da tabela, inclua:

- filtros usados
- comandos `pp-leads` executados, quando houver
- arquivos consultados
- leads excluídos e motivo, quando relevante
- pendências, como contato ausente, autenticação ausente ou necessidade de validação manual

Salve artefatos em `.context/organizejr/` quando a entrega gerar arquivo reutilizável.

## Regras de qualidade

- Não trate a OrganizeJr. como lead de venda.
- Não invente e-mail, telefone, cargo, faturamento, dor ou vínculo.
- Prefira contatos institucionais.
- Não faça disparo em massa. A saída deve apoiar abordagem personalizada.
- Não use dado sensível ou pessoal sem necessidade comercial clara.
- Se o usuário pedir uma lista operacional, entregue arquivo `.csv` ou `.xlsx`; se pedir estratégia, entregue `.md` com tabela e mensagens.
- Se a demanda explícita de Lucas contrariar a pirâmide, explique o tradeoff e siga a demanda explícita.

## Exemplos de acionamento

- "usa o pp-leads para encontrar contabilidades em Salvador para a OrganizeJr"
- "pega essa planilha de demandas e monta leads e abordagem"
- "quero EJs da Bahia com mais chance de comprar capacitação de gestão de tempo"
- "enriquece esses CNPJs e cria mensagem para abordagem"
- "acha escolas para feira de profissões usando a pirâmide da OrganizeJr"
