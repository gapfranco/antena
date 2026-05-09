# Diretrizes do Repositório

## Estrutura do Projeto e Organização dos Módulos

O Antena é um dashboard Go somente leitura para eventos armazenados em Turso/libSQL. O ponto de entrada fica em `cmd/antena/`, com `main.go` responsável pelas rotas HTTP e renderização de templates. Os modos desktop e servidor são separados por build tags em `desktop.go` (`!headless`) e `desktop_stub.go` (`headless`). O acesso a dados fica em `internal/models/events.go`, a configuração em `config/config.go`, e os templates/ativos embutidos em `ui/`. Templates HTML ficam em `ui/html/`, arquivos estáticos em `ui/static/`, e a entrada do Tailwind é `ui/input.css`.

## Comandos de Build, Teste e Desenvolvimento

- `make build`: gera o binário desktop em `./build/antena` e compila o CSS minificado.
- `make build-server`: gera o servidor HTTP headless em `./build/antena-server`.
- `make tailwind-build`: compila `ui/input.css` para `ui/static/css/styles.css`.
- `make tailwind-watch`: observa alterações no Tailwind durante o desenvolvimento da UI.
- `make clean`: remove artefatos de build e CSS gerado.
- `go test ./...`: executa os testes Go quando houver testes no projeto.

Para rodar localmente, crie `antena.conf` no diretório de trabalho:

```env
TURSO_URL=libsql://...
TURSO_TOKEN=...
ADDR=:4000
```

## Estilo de Código e Convenções de Nomes

Use a formatação padrão do Go: execute `gofmt` nos arquivos `.go` modificados antes de commitar. Mantenha nomes de pacotes curtos e em minúsculas. Exporte apenas tipos e funções usados por outros pacotes. Prefira nomes claros para handlers e models, alinhados às rotas ou operações de dados, como handlers de `events` e `EventModel.Count`. Coloque templates em `ui/html/pages` ou `ui/html/partials` conforme sejam páginas completas ou fragmentos reutilizáveis.

## Diretrizes de Testes

Atualmente não há arquivos `*_test.go` commitados nem alvo de teste no Makefile. Adicione testes Go próximos ao código coberto usando a convenção `*_test.go`, depois rode `go test ./...`. Para mudanças em acesso a dados, priorize testes sobre construção de queries e casos de borda como filtros vazios, paginação e saída de exportação.

## Diretrizes de Commit e Pull Request

Os commits recentes usam resumos curtos em português, como `ajuste`, `readme atualizado` e `painel_inicial`. Mantenha mensagens concisas, no imperativo e específicas sobre a alteração. Pull requests devem descrever o comportamento alterado, listar comandos de verificação manual, mencionar premissas de configuração ou schema, e incluir capturas de tela quando templates ou estilos Tailwind mudarem.

## Segurança e Configuração

Não commite `antena.conf`, tokens do Turso, binários gerados ou saídas locais de build. Trate esta aplicação como somente leitura: alterações não devem introduzir escritas no banco sem uma atualização explícita da arquitetura.
