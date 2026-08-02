# Arquitetura do Leads Brasil Platform (Printing Press)

## Visão Geral
Plataforma de inteligência de leads (tipo Apollo) gerada automaticamente via `cli-printing-press`.

## Componentes
- `leads-pp-backend` (Go): API que expõe a "Fábrica PP".
- `catalog/leads-brasil.yaml`: Especificação OpenAPI que define a superfície da API.
- `leads-pp-cli` (Gerado): CLI de controle, gerado pelo Printing Press a partir da spec.

## Fluxo de Dados
1. Você define/altera o YAML em `catalog/`.
2. `cli-printing-press` gera o CLI baseado nessa spec.
3. CLI chama a API do backend (`leads-pp-backend`).
4. O backend consulta Casa dos Dados, Company Goat e Contact Goat. O enriquecimento social sempre tenta Scrape Creators; quando ele não retorna dados utilizáveis, tenta Apify somente se o perfil configurar Actor e credenciais.

## Operações externas

O núcleo público só cria e armazena planos de operação. O endpoint `POST /operations/plan` não produz efeito externo e exige um perfil configurado explicitamente. `POST /operations/{plan_id}/apply` exige um ID de plano válido e `approved: true`; cada alvo gera recibo idempotente. Um perfil pode declarar um `operation_adapter_command` executável, que recebe somente o alvo salvo como JSON em stdin. Credenciais e adaptadores de CRM, mensageria ou planilha pertencem ao perfil privado e não são distribuídos neste repositório.