# Pesquisa, pp-apify, pp-company-goat e pp-contact-goat

Atualizado em: 2026-06-17

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

- deve ser o adaptador para chamar Actors Apify quando o enriquecimento exigir scraping, datasets ou execução recorrente;
- não deve substituir `company-goat` ou `contact-goat`; ele é a infraestrutura para rodar Actors quando o fluxo precisar.

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

Implicação para OrganizeJr:

- Para EJs, `contact-goat` deve ser usado depois de `company-goat`/Casa dos Dados, quando o lead já foi priorizado.
- Ele não deve ser usado para disparo em massa.
- Deve procurar pessoa/cargo responsável por presidência, comercial, projetos, pessoas, atendimento ou diretoria da EJ.
- O output esperado para Lucas deve incluir e-mail verificado ou provável, link de perfil, caminho de warm intro, fonte e confiança.

## Como fica a arquitetura correta do pp-leads-brasil

### Camada 1, base local OrganizeJr

Fontes:

- `Relatório Portal BJ EJs 2025 0101.xlsx`
- `.context/organizejr/lead-table-2026-06-17-ejs-comunicacao-lucas.csv`
- `.context/organizejr/metodologia-score-leads-organizejr.md`
- `.context/docs/organizejr-commercial-context.md`

Uso:

- ranking, ICP, evidência, mensagem e próximos passos.

### Camada 2, Casa dos Dados

Fonte:

- API Casa dos Dados, com `api-key`.

Uso:

- confirmar CNPJ;
- razão social;
- situação cadastral;
- endereço;
- CNAE;
- quadro societário;
- telefone/e-mail quando disponível.

### Camada 3, company-goat

Uso ideal:

- pesquisa ampliada de empresa/domínio;
- sinais de legitimidade, tração, presença técnica e diretórios;
- mais útil para startups e empresas sênior do que para EJs do relatório BJ.

### Camada 4, contact-goat

Uso ideal:

- achar pessoa responsável e contato individual/profissional;
- mapear warm intro;
- enriquecer e-mail quando o canal institucional não resolve.

### Camada 5, pp-apify

Uso ideal:

- rodar Actors específicos quando for preciso scraping, dataset ou automação recorrente;
- registrar custo e output;
- não assumir que todo goat é um Actor Apify.

## Decisão para o código atual

O `pp-leads-brasil` não deve fingir que `company-goat` e `contact-goat` são apenas endpoints locais simples.

Implementação recomendada:

1. `internal/client/casadados`: cliente real da Casa dos Dados.
2. `internal/client/pp`: orquestrador que chama:
   - base local OrganizeJr;
   - Casa dos Dados;
   - `pp-company-goat` quando instalado;
   - `pp-contact-goat` quando instalado;
   - `apify-pp-cli` apenas para Actors Apify explícitos.
3. `enrich`: deve retornar quais camadas rodaram, quais falharam, custo/pendência e fontes usadas.

Atualização de implementação em 2026-06-17:

- `internal/client/casadados` já consulta Casa dos Dados real quando a chave está configurada.
- `internal/client/pp` já orquestra base local, Casa dos Dados, `company-goat` e `contact-goat`.
- `leads-brasil-pp-cli` já expõe comandos `company-goat` e `contact-goat`.
- Script local de preparação das credenciais: `.context/organizejr/setup-contact-goat-auth.sh`
- `apify-pp-cli` continua reservado para Actors Apify explícitos, não para substituir os goats.

## Variáveis/credenciais esperadas

- `CASA_DADOS_API_KEY` ou `PP_LEADS_CASA_DADOS_API_KEY`, para Casa dos Dados.
- Credenciais próprias do `pp-contact-goat`, conforme skill dele: LinkedIn MCP/Happenstance/Deepline.
- `APIFY_TOKEN`, apenas se rodar Actors via `pp-apify`.

## Risco

`pp-contact-goat` pode lidar com dados pessoais/profissionais. Para a OrganizeJr, use com escopo conservador: contatos institucionais primeiro; contato individual só quando houver necessidade comercial clara e fonte profissional pública.
