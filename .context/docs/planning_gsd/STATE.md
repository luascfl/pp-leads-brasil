# Estado do planejamento GSD

Atualizado em: 2026-06-17

## Foco atual

Fechar a camada externa real de enriquecimento da OrganizeJr com Casa dos Dados, `company-goat` e `contact-goat`.

## Contexto operacional

- Lucas Camilo Carvalho é assessor comercial da OrganizeJr.
- OrganizeJr é uma empresa júnior de Psicologia Organizacional da UNEB Salvador.
- A pasta recebe demandas comerciais em Excel, texto e documentos.
- As demandas devem ser transformadas em busca de leads, enriquecimento por CNPJ/empresa, priorização por ICP e abordagem personalizada.

## Fontes de verdade deste ciclo

- `.context/docs/organizejr-commercial-context.md`
- `.context/skills/organizejr-pp-leads/SKILL.md`
- `Pirâmide de captação e ICPs da OrganizeJr.docx`
- `Relatório Portal BJ EJs 2025 0101.xlsx`
- `leads-brasil-pp-cli/SKILL.md`


## Demanda comercial ativa

- Lucas fica responsável por prospectar EJs de comunicação em geral em nível nacional.
- Hipótese de ICP: `ICP 2.3 — EJs de Comunicação com Sobrecarga Criativa e Relacional`.
- Tese a validar: a dor de EJs de comunicação não é apenas alta demanda operacional; ela aparece em fadiga criativa, trabalho always-on, aprovações subjetivas, retrabalho, conflito entre criação/atendimento/estratégia, onboarding frágil, insegurança de feedback, exposição em redes, crise/moderação e perda de padrão criativo com troca de gestão.
- Fonte inicial: `Relatório Portal BJ EJs 2025 0101.xlsx`, com filtro nacional por termos de comunicação em `EMPRESA_JUNIOR` e `CURSOS`.
- Case candidato: Acesso Comunicação Jr.; benchmark de maturidade: ESPM Jr.

## Status

- Contexto comercial OrganizeJr: preenchido.
- Skill local OrganizeJr PP Leads: criada e sincronizada.
- Casa dos Dados: cliente real implementado.
- `company-goat` e `contact-goat`: expostos no backend local e no `leads-brasil-pp-cli`.
- `enrich`: consolida lead local, Casa dos Dados, `company-goat`, `contact-goat`, `use_case`, `lead_context` e links de abordagem.
- Separação concluída: `pp-leads-brasil` não hardcode mais OrganizeJr; o caso atual roda via `.context/use-cases/organizejr-ejs-comunicacao/use-case.json`, `scripts/run-use-case-enrichment.sh` e `scripts/setup-external-enrichment-auth.sh`.
