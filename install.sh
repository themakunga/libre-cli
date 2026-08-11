#!/bin/bash
set -e

APP_NAME="libre-cli"
INSTALL_DIR="/usr/local/bin"

echo "======================================"
echo " Instalador de $APP_NAME"
echo "======================================"

# Verificar si Go está instalado
if ! command -v go &> /dev/null; then
    echo "❌ Error: 'go' no está instalado. Instálalo primero (https://go.dev/doc/install)."
    exit 1
fi

echo "🔨 Descargando dependencias..."
go mod tidy

echo "🔨 Construyendo el binario..."
go build -o $APP_NAME ./cmd/$APP_NAME

echo "🚀 Instalando en $INSTALL_DIR (se requieren permisos de administrador)..."
sudo mv $APP_NAME $INSTALL_DIR/$APP_NAME
sudo chmod +x $INSTALL_DIR/$APP_NAME

echo "✅ ¡Instalación completada exitosamente!"
echo "➡️  Ahora puedes usar el comando: $APP_NAME"
