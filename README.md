# Antena

Dashboard de monitoramento read-only para eventos IoT de centrais de alarme de incêndio. Consome o banco remoto Turso/libSQL alimentado pelo **[Signal](../signal)** e exibe os dados em interface web com filtros, paginação e exportação CSV.

## Relação com o Signal

O **Signal** é o serviço de coleta: recebe eventos via TCP das centrais, persiste localmente em SQLite e sincroniza de forma assíncrona com o Turso. O **Antena** é o visualizador: conecta-se diretamente ao Turso e expõe os dados de forma estruturada, sem escrever nada no banco.

```
Central (TCP) → Signal → SQLite → Turso (libSQL)
                                        ↑
                                     Antena (leitura)
```

## Funcionalidades

- Lista de instalações com contagem de eventos
- Eventos paginados com filtros por tipo, central, instalação e dispositivo
- Exportação CSV (separador ponto-e-vírgula, BOM UTF-8, extensão `.xls`)
- Dois modos de deploy: desktop com janela nativa (WebKit) e servidor headless

## Modos de deploy

| Modo | Build | Descrição |
|------|-------|-----------|
| **Desktop** | `make build` | Janela nativa (webview) em porta efêmera |
| **Headless** | `make build-server` | Servidor HTTP puro na porta configurada (`ADDR`) |

## Configuração

Crie `antena.conf` no diretório de trabalho:

```env
TURSO_URL=libsql://seu-banco.turso.io
TURSO_TOKEN=seu-token
ADDR=:4000
```

`ADDR` é opcional (padrão `:4000`).

## Como rodar

```bash
make build && ./build/antena
```

## Builds disponíveis

```bash
make build               # Desktop (webview)
make build-server        # Headless server

# Cross-platform — desktop
make windows             # ./build/antena.exe
make linux               # ./build/antena_linux
make darwin              # ./build/antena_mac

# Cross-platform — headless
make server-windows
make server-linux
make server-darwin

# CSS apenas
make tailwind-build      # build único
make tailwind-watch      # modo watch (desenvolvimento)
```

## Estrutura

```
.
├── cmd/antena/
│   ├── main.go           # HTTP server, handlers, cache de templates
│   ├── desktop.go        # Janela webview (build !headless)
│   └── desktop_stub.go   # Stub headless (build headless)
├── internal/
│   └── models/
│       └── events.go     # EventModel — queries Turso (Installations, All, Count, GetForExport)
├── config/
│   └── config.go         # Config via Viper (antena.conf)
├── ui/
│   ├── html/             # Templates Go (base, pages, partials)
│   ├── static/           # CSS (Tailwind), JS (HTMX)
│   └── efs.go            # embed.FS — ativos embutidos no binário
├── tailwind.config.js
└── Makefile
```

## Rotas

| Método | Path | Descrição |
|--------|------|-----------|
| GET | `/` | Lista de instalações |
| GET | `/installations` | Lista de instalações |
| GET | `/events` | Eventos paginados com filtros |
| GET | `/export` | Formulário de exportação |
| POST | `/export` | Download CSV |
| GET | `/static/*` | Ativos estáticos embutidos |

## Filtros de eventos (query params)

| Parâmetro | Tipo | Matching |
|-----------|------|---------|
| `event_type` | string | LIKE `%valor%` |
| `central_id` | int | exato |
| `inst_id` | string | LIKE `%valor%` |
| `device` | string | LIKE `%valor%` |

## Principais dependências

| Pacote | Uso |
|--------|-----|
| `tursodatabase/libsql-client-go` | Banco remoto Turso |
| `spf13/viper` | Configuração via arquivo `.conf` |
| `webview/webview_go` | Janela nativa (modo desktop) |