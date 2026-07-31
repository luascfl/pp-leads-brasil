# Pesquisa, pp-apify, pp-company-goat e pp-contact-goat

Atualizado em: 2026-07-31

## Conclusão operacional

No ecossistema Printing Press, `pp-apify`, `pp-company-goat` e `pp-contact-goat` são três ferramentas diferentes:

1. `pp-apify` é a camada para operar a plataforma Apify: rodar Actors, controlar custo, buscar datasets e manter store local.
2. `pp-company-goat` é uma ferramenta de pesquisa de empresas/startups. Não é apenas um Actor Apify. Ela agrega sinais de SEC Form D, GitHub, Hacker News, Companies House, YC, Wikidata e DNS/RDAP.
3. `pp-contact-goat` é uma ferramenta de enriquecimento de contatos e relacionamento. A fonte pública encontrada descreve LinkedIn MCP, Happenstance, Deepline e store local de contatos. Não encontrei evidência de que ele seja um Actor Apify específico.

## pp-apify

Fonte local consultada: `skill://pp-apify`.

Função:

- operar qualquer Actor da Apify;
- rodar Actor com controle de custo;
- buscar itens em datasets locais;
- manter SQLite local com resultados passados;
- comparar Actors;
- gerar digest/newsletter;
- aplicar schedules.

Uso no `pp-leads-brasil`:

- Scrape Creators é a primeira tentativa obrigatória de enriquecimento social por lead.
- Apify é o fallback para quando Scrape Creators não retorna dados utilizáveis por billing, credencial, timeout ou erro de provedor.
- O fallback requer `APIFY_TOKEN` e `APIFY_SOCIAL_ACTOR_ID`; sem ambos, ele registra indisponibilidade sem interromper o enriquecimento.
- `APIFY_SOCIAL_MAX_COST` é opcional e repassado como `--max-cost`, para recusar uma execução cuja projeção exceda o teto.
- Nenhum dos dois provedores substitui `company-goat` ou `contact-goat`.

Comandos base:

```bash
apify-pp-cli doctor
apify-pp-cli store --search "company goat" --agent
apify-pp-cli run <actor-id> --input @input.json --agent --max-cost 0.50
```

Achado: `apify-pp-cli store --search "company goat"` e `apify-pp-cli store --search "contact goat"` retornaram vazio na busca pública da Store. Isso sugere que `pp-company-goat` e `pp-contact-goat` não são, por padrão, Actors públicos com esses nomes.

## pp-company-goat

Fonte principal: repositório `mvanhorn/printing-press-library`, catálogo oficial.

Descrição encontrada:

> `company-goat`: Look up startups across SEC Form D, GitHub, Hacker News, Companies House, YC, Wikidata, and DNS in one command. Built for scouts, founders, BD, and AI agents researching small/midsize startups.

Fonte: https://github.com/mvanhorn/printing-press-library

Também foi encontrado um Actor Apify relacionado, mas não equivalente:

- Actor: `changeable_peddler/startup-funding-signal-report`
- URL: https://apify.com/changeable_peddler/startup-funding-signal-report
- O próprio listing diz ser inspirado por Printing Press `company-goat`.
- Sinais: SEC Form D, Hacker News, GitHub, RDAP/domain age e diretórios públicos.
- Uso provável: pesquisa de startups, funding e tração, não EJs brasileiras.

Implicação para OrganizeJr:

- `company-goat` puro é menos aderente para EJs brasileiras de comunicação.
- Para EJs, o equivalente local deve combinar Casa dos Dados, relatório BJ, site da EJ, redes sociais e contexto OrganizeJr.
- Actor Apify inspirado em company-goat pode servir para empresas/startups, não para a lista de EJs atual.

## pp-contact-goat

Fonte principal: Printing Press homepage e marketplace de skills.

Descrição encontrada:

- “Super LinkedIn” para terminal.
- Pesquisa e enriquecimento de contatos.
- Fontes: LinkedIn MCP server, Happenstance, Deepline e store local.
- Fluxo: buscar no LinkedIn, cruzar Happenstance para warm intros, pagar Deepline para e-mail verificado quando necessário.

Fontes:

- https://printingpress.dev/
- https://claudemarketplaces.com/skills/mvanhorn/printing-press-library/pp-contact-goat

## Uso da plataforma por perfis externos

Perfis externos podem combinar as fontes desta plataforma com seus próprios ICPs, dados locais e processos comerciais. Esses perfis não fazem parte do repositório público e devem ser ativados apenas por configuração explícita.

### Camada 1, tabela local declarada pelo perfil

Uso:

- prover contexto de lead, campos e identificadores conhecidos;
- fornecer a base para busca por nome e CNPJ;
- manter dados proprietários fora do repositório público.

### Camada 2, Casa dos Dados

Uso:

- confirmar CNPJ;
- razão social;
- situação cadastral;
- endereço;
- CNAE;
- quadro societário;
- telefone e e-mail quando disponíveis.

### Camada 3, Company Goat

Uso:

- pesquisa ampliada de empresa e domínio;
- sinais públicos de legitimidade, tração, presença técnica e diretórios.

### Camada 4, Contact Goat

Uso:

- encontrar contato profissional relevante;
- mapear possível warm intro;
- enriquecer canal profissional quando o canal institucional não resolve.

### Camada 5, enriquecimento social

Uso:

- Scrape Creators é sempre tentado primeiro para localizar e obter dados públicos de LinkedIn, Instagram e outros canais.
- Apify roda apenas como fallback configurado por perfil.
- Cada resposta informa `social_provider` e registra o status de Scrape Creators e, quando acionado, do Apify.
- Billing, ausência de credencial e falha do provedor não invalidam os dados já obtidos por outras camadas.

## Decisão de arquitetura

O cliente `internal/client/casadados` resolve a tabela local declarada e consulta Casa dos Dados quando configurado. `internal/client/pp` orquestra as fontes disponíveis e reporta quais camadas foram executadas, falharam ou ficaram indisponíveis. Nenhum cliente deve embutir o nome, os arquivos ou o fluxo operacional de um perfil.

## Variáveis e credenciais esperadas

- `CASA_DADOS_API_KEY` ou `PP_LEADS_CASA_DADOS_API_KEY`, para Casa dos Dados.
- Credenciais próprias do Contact Goat, conforme sua integração configurada.
- `SCRAPE_CREATORS_API_KEY`, para a tentativa primária de enriquecimento social.
- `APIFY_TOKEN` e `APIFY_SOCIAL_ACTOR_ID`, para o fallback social por Actor. `APIFY_SOCIAL_MAX_COST` define, quando desejado, o teto de custo projetado por execução.

## Risco

Contact Goat pode lidar com dados pessoais e profissionais. Perfis devem priorizar fonte pública, minimização de dados e finalidade comercial legítima.
