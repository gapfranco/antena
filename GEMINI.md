# Projeto Antena: Guia de Contexto

Este documento fornece o contexto técnico e as diretrizes de desenvolvimento para os projetos contidos neste repositório.

## 🚀 Visão Geral dos Projetos

### Antena
Visualizador remoto (read-only) conectado diretamente ao banco de dados na nuvem.
- **Arquitetura:** Servidor Web que consulta o Turso diretamente.
- **Diferencial:** Sem persistência local; focado em monitoramento em tempo real.
- **Tecnologias:** Go 1.25, Turso, Viper, HTMX, Tailwind CSS.

---

## 🛠 Comandos de Build e Execução

### Padrão
- **Build:** `make build`
- **Tailwind:** `make tailwind-build`
- **Executar:** Binários gerados em `./bin/`

### Configurações Específicas
- `antena.conf` (Direct Turso access).

---

## 📐 Convenções de Desenvolvimento

- **Backend:** 
    - Uso de `net/http` nativo (Go 1.22+ routing).
    - Modelos em `internal/models` e persistência em `internal/storage`.
    - Configurações via **Viper** (.conf/env).
- **Frontend:**
    - **MPA** com `html/template` e **HTMX** para reatividade parcial.
    - Estilização via **Tailwind CSS** (compilado para `ui/static/css/styles.css`).
    - Ativos embutidos via `embed.FS` em `ui/efs.go`.
- **Banco de Dados:**
    - **SQLite:** Local sem CGO.
    - **Turso:** Remoto 

---

## 📁 Estrutura de Diretórios (Comum)

- `cmd/`: Pontos de entrada (main.go, handlers, routes).
- `internal/`: Lógica de negócio, modelos, storage e parsers (ex: OFX).
- `ui/`: Templates HTML e arquivos estáticos (CSS, JS, Imagens).
- `config/`: Loader de configurações via Viper.
- `bin/`: Binários compilados por plataforma.
