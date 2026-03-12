
# Variáveis de Configuração
BINARY_NAME=antena
SRC_PATH=./cmd/antena
OUT_DIR=./bin
TAILWIND_BIN=tailwindcss
TW_INPUT=./ui/input.css
TW_OUTPUT=./ui/static/css/styles.css

# Comandos principais
.PHONY: all clean build tailwind-build tailwind-watch

all: build ## Compila o projeto

# Regra para gerar o CSS minificado
tailwind-build:
	@echo "Gerando CSS com Tailwind..."
	$(TAILWIND_BIN) -i $(TW_INPUT) -o $(TW_OUTPUT) --minify

# Regra utilitária para desenvolvimento
tailwind-watch:
	$(TAILWIND_BIN) -i $(TW_INPUT) -o $(TW_OUTPUT) --watch

## Compilação
build: tailwind-build
	@echo "Construindo antena..."
	@mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(BINARY_NAME) $(SRC_PATH)

## Compilação para Windows (amd64)
# Adicionamos tailwind-build como dependência
windows: tailwind-build
	@echo "Construindo para Windows..."
	@mkdir -p $(OUT_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(OUT_DIR)/$(BINARY_NAME).exe $(SRC_PATH)

## Compilação para Linux (amd64)
linux: tailwind-build
	@echo "Construindo para Linux..."
	@mkdir -p $(OUT_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(OUT_DIR)/$(BINARY_NAME)_linux $(SRC_PATH)

## Compilação para macOS (Darwin amd64/arm64)
darwin: tailwind-build
	@echo "Construindo para macOS..."
	@mkdir -p $(OUT_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(OUT_DIR)/$(BINARY_NAME)_mac $(SRC_PATH)

## Limpa a pasta de builds e o CSS gerado
clean:
	@echo "Limpando arquivos..."
	rm -rf $(OUT_DIR)
	rm -f $(TW_OUTPUT)
