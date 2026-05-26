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
4. O backend invoca a Fábrica PP (Apify/Casa dos Dados).