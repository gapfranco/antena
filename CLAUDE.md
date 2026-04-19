# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build          # Desktop binary → ./build/antena (includes Tailwind CSS build)
make build-server   # Headless server binary → ./build/antena-server (-tags headless)
make tailwind-build # Generate minified CSS: ui/input.css → ui/static/css/styles.css
make tailwind-watch # Watch Tailwind CSS changes (use during development)
make clean          # Remove ./build and generated CSS

# Desktop (requires native OS runner + CGo + WebKit)
make windows        # → ./build/antena.exe
make linux          # → ./build/antena_linux
make darwin         # → ./build/antena_mac

# Headless server (pure Go, cross-compile from any OS)
make server-linux   # → ./build/antena-server_linux
make server-windows # → ./build/antena-server.exe
make server-darwin  # → ./build/antena-server_mac
```

No test or lint targets exist in the Makefile.

### Running locally

Requires `antena.conf` in the working directory (gitignored):
```
TURSO_URL=libsql://...
TURSO_TOKEN=<jwt>
ADDR=:4000
```

```bash
make build && ./build/antena   # Serves on :4000 by default
```

## Architecture

**Antena** is a read-only monitoring dashboard for a Turso (cloud SQLite) database of IoT events. It uses standard Go `net/http` (Go 1.22+ routing), server-rendered HTML templates, HTMX for partial updates, and Tailwind CSS.

### Desktop vs headless

The app has two build modes controlled by the `headless` build tag:

- **Desktop** (default, `!headless`): `desktop.go` starts the HTTP server on a random loopback port and opens a native window via `webview_go` (uses OS WebKit/WebView2). Requires CGo and OS-specific WebKit libs.
- **Headless** (`-tags headless`): `desktop_stub.go` simply calls `runServer()` — pure Go, no CGo, cross-compilable.

`main.go` always calls `runDesktop(app, addr)`; the build tag selects which implementation is compiled in.

### Key files

- `cmd/antena/main.go` — HTTP server (`runServer`), all route handlers, template rendering, startup sequence
- `cmd/antena/desktop.go` — webview desktop window (`//go:build !headless`)
- `cmd/antena/desktop_stub.go` — headless stub (`//go:build headless`)
- `internal/models/events.go` — `EventModel` with DB queries: `Installations()`, `All()` (paginated/filtered), `Count()`, `GetForExport()`
- `config/config.go` — Viper config loader (reads `antena.conf` then env vars)
- `ui/efs.go` — declares the `embed.FS` for static assets and templates
- `ui/html/` — Go `html/template` files: `base.html`, `pages/` (`installations.html`, `events.html`, `export.html`, `nav.html`), `partials/` (`nav.html`)

### Request flow

1. `main.go` loads config → opens Turso DB → parses all templates into a cache → starts HTTP server
2. Each handler queries `EventModel`, populates a `templateData` struct, and renders the appropriate cached template
3. Static assets served from embedded `ui/static/` FS

### Routes

| Method | Path | Handler |
|--------|------|---------|
| GET | `/` | installations list |
| GET | `/installations` | installations list |
| GET | `/events` | paginated event list with filters |
| GET | `/export` | export form |
| POST | `/export` | CSV download (semicolon-separated, UTF-8 BOM, `.xls` extension) |
| GET | `/static/*` | embedded static assets |

### Data models

`Event` fields: `ID`, `Central`, `Link`, `DeviceId`, `EventType`, `Local`, `Device`, `DeviceType`, `TsUnixMs`, `InstId`, `TypeId`

`Installation` fields: `InstId`, `EventCount`

`templateData` fields: `CurrentYear`, `IsAuthenticated`, `ActiveMenu`, `Events`, `Installations`, `Centrals`, `SearchType`, `SearchCentral`, `SearchInstID`, `SearchDevice`, `CurrentPage`, `TotalPages`, `HasNextPage`, `HasPrevPage`

### Event list filters (query params)

`event_type`, `central_id`, `inst_id`, `device` — all support partial matching (LIKE `%value%`), except `central_id` (exact int match).

### Template functions

Registered in `main.go`: `humanDate`, `add`, `formatUnixMs`

### Config

Viper reads `antena.conf` (env format) first, then falls back to environment variables. Required keys: `TURSO_URL`, `TURSO_TOKEN`. Optional: `ADDR` (default `:4000`).