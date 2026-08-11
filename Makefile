APP_NAME = libre-cli
BUILD_DIR = bin
INSTALL_DIR = /usr/local/bin

.PHONY: all build install clean

all: build

build:
	@echo "==> Construyendo $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)
	@echo "==> ¡Construcción completada en $(BUILD_DIR)/$(APP_NAME)!"

install: build
	@echo "==> Instalando $(APP_NAME) en $(INSTALL_DIR) (puede pedir contraseña)..."
	@sudo cp $(BUILD_DIR)/$(APP_NAME) $(INSTALL_DIR)/$(APP_NAME)
	@sudo chmod +x $(INSTALL_DIR)/$(APP_NAME)
	@echo "==> ¡Instalación exitosa! Ejecuta '$(APP_NAME)' desde cualquier terminal."

clean:
	@echo "==> Limpiando binarios..."
	@rm -rf $(BUILD_DIR)
