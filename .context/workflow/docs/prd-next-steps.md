# Próximos Passos (Next Steps)

O backend do `Leads Brasil Platform` (nosso Apollo brasileiro) está funcional. Ele age como a API que sustenta a sua fábrica.

O que você tem agora:
- Uma API Go (`cmd/server/main.go`) escutando na porta 8080.
- Uma documentação Swagger/OpenAPI da API (`catalog/leads-brasil.yaml`).

O que você faz "agora" no seu computador para ver a mágica acontecer:

## Passo 1: Ligar o Backend
No terminal do seu projeto (`~/Downloads/pp-leads-brasil`), rode a sua API:
```bash
go run ./cmd/server/...
```
*(Ela vai ficar rodando no terminal, esperando chamadas na porta 8080)*

## Passo 2: Gerar o seu CLI com o Printing Press
Agora você vai usar o repositório que você achou no GitHub (`mvanhorn/cli-printing-press`) para **gerar o aplicativo CLI**.

1. Em outro terminal, clone ou entre na pasta do `cli-printing-press` no seu PC.
2. Copie a especificação que criamos para a pasta de catálogo do Printing Press:
```bash
cp ~/Downloads/pp-leads-brasil/catalog/leads-brasil.yaml /caminho/do/seu/cli-printing-press/catalog/
```
3. Use os comandos do Printing Press (conforme o README dele) para gerar o binário:
```bash
# Exemplo fictício baseado no padrão do Printing Press
go run cmd/printing-press/main.go generate --catalog catalog/leads-brasil.yaml
```

## Passo 3: Usar o seu "Apollo"
O Printing Press vai criar um executável (ex: `leads-brasil`). Com o seu backend rodando (Passo 1), você já pode usá-lo:

```bash
# Buscar uma empresa (chama a nossa rota /search)
leads-brasil search --name "Empresa Exemplo"

# Enriquecer usando a Apify (chama nossa rota /enrich)
leads-brasil enrich 12345678000199
```

## Resumo Estratégico
O seu fluxo está dividido:
1.  **Quando você quiser conectar uma nova ferramenta de dados:** Você edita os `clients` e `handlers` no seu projeto Go.
2.  **Quando você quiser mudar os comandos do seu CLI:** Você edita o `leads-brasil.yaml` e roda o Printing Press de novo.