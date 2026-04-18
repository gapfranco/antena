
# Variáveis de Configuração
BINARY_NAME=antena
SRC_PATH=./cmd/antena
OUT_DIR=./build
TAILWIND_BIN=tailwindcss
TW_INPUT=./ui/input.css
TW_OUTPUT=./ui/static/css/styles.css

# Comandos principais
.PHONY: all clean build build-server tailwind-build tailwind-watch \
        windows linux darwin \
        server-windows server-linux server-darwin

all: build ## Compila o projeto

# Regra para gerar o CSS minificado
tailwind-build:
	@echo "Gerando CSS com Tailwind..."
	$(TAILWIND_BIN) -i $(TW_INPUT) -o $(TW_OUTPUT) --minify

# Regra utilitária para desenvolvimento
tailwind-watch:
	$(TAILWIND_BIN) -i $(TW_INPUT) -o $(TW_OUTPUT) --watch

## Desktop (janela nativa via webview) — padrão
build: tailwind-build
	@echo "Construindo antena (desktop)..."
	@mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(BINARY_NAME) $(SRC_PATH)

## Servidor headless (sem janela, para deploy em servidor)
build-server: tailwind-build
	@echo "Construindo antena (headless)..."
	@mkdir -p $(OUT_DIR)
	go build -tags headless -o $(OUT_DIR)/$(BINARY_NAME)-server $(SRC_PATH)

## Desktop nativo — requer compilar no OS alvo (CGo)
windows: tailwind-build
	@mkdir -p $(OUT_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui" -o $(OUT_DIR)/$(BINARY_NAME).exe $(SRC_PATH)

linux: tailwind-build
	@mkdir -p $(OUT_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(OUT_DIR)/$(BINARY_NAME)_linux $(SRC_PATH)

darwin: tailwind-build
	@mkdir -p $(OUT_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(OUT_DIR)/$(BINARY_NAME)_mac $(SRC_PATH)

## Servidor headless cross-platform (sem CGo, pode compilar de qualquer OS)
server-windows: tailwind-build
	@mkdir -p $(OUT_DIR)
	GOOS=windows GOARCH=amd64 go build -tags headless -o $(OUT_DIR)/$(BINARY_NAME)-server.exe $(SRC_PATH)

server-linux: tailwind-build
	@mkdir -p $(OUT_DIR)
	GOOS=linux GOARCH=amd64 go build -tags headless -o $(OUT_DIR)/$(BINARY_NAME)-server_linux $(SRC_PATH)

server-darwin: tailwind-build
	@mkdir -p $(OUT_DIR)
	GOOS=darwin GOARCH=amd64 go build -tags headless -o $(OUT_DIR)/$(BINARY_NAME)-server_mac $(SRC_PATH)

## Limpa a pasta de builds e o CSS gerado
clean:
	@echo "Limpando arquivos..."
	rm -rf $(OUT_DIR)
	rm -f $(TW_OUTPUT)
