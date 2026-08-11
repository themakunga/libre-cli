# ⚡ FreeStyle Libre 2 CLI Monitor

![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go)
![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux-lightgrey?style=for-the-badge)
![UI](https://img.shields.io/badge/UI-Bubble%20Tea%20%7C%20Lipgloss-FF0055?style=for-the-badge)
![Config](https://img.shields.io/badge/Config-TOML-8B4513?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Active-brightgreen?style=for-the-badge)

[🇺🇸 English](#-english) | [🇪🇸 Español](#-español)

---

## 🇺🇸 English

### 📖 Overview
A beautiful, highly customizable Terminal-based UI (TUI) to monitor your **FreeStyle Libre 2** glucose data. It fetches data (via LibreLinkUp simulation/API) and displays your current glucose levels, trend arrows, and an ASCII graph of the last few hours right in your terminal.

Built with Golang, Bubble Tea, and Lipgloss.

### ✨ Features
*   **Real-time Monitoring:** Auto-updates every 5 minutes (configurable).
*   **Terminal Graph:** Renders a clean line graph inside the CLI.
*   **Dynamic Colors:** Highlights glucose numbers in green, yellow, or red depending on your target range.
*   **TOML Configuration:** Easy setup and theming (comes with **Tokyo Night** by default).
*   **Cross-Platform:** Works seamlessly on Linux and macOS.

### 🚀 Installation

You can install the application using standard `make` or via the provided bash script.

**Option 1: Using Make (Recommended)**
make install

**Option 2: Using the Installation Script**
./install.sh

### ⚙️ Configuration
On its first run, the app will automatically generate a configuration file located at:
~/.config/libre-cli/config.toml

You can edit this file to change your LibreLinkUp credentials, adjust the auto-update interval, define your target glucose range (min_glucose / max_glucose), and completely customize the theme colors.

### 💻 Usage
Simply run the command from any terminal:
libre-cli

Press q or Ctrl+C to safely exit the application.

---

## 🇪🇸 Español

### 📖 Descripción general
Una hermosa y altamente personalizable interfaz de terminal (TUI) para monitorear tus datos de glucosa del sensor **FreeStyle Libre 2**. Obtiene datos (a través de LibreLinkUp) y muestra tus niveles actuales de glucosa, flechas de tendencia y un gráfico ASCII de las últimas horas directamente en tu terminal.

Construido con Golang, Bubble Tea y Lipgloss.

### ✨ Características
*   **Monitoreo en tiempo real:** Se actualiza automáticamente cada 5 minutos (configurable).
*   **Gráfico en Terminal:** Dibuja un gráfico de líneas limpio dentro de la CLI.
*   **Colores dinámicos:** Resalta los números de glucosa en verde, amarillo o rojo dependiendo de tu rango objetivo.
*   **Configuración TOML:** Fácil configuración y tematización (incluye el tema **Tokyo Night** por defecto).
*   **Multiplataforma:** Funciona sin problemas en Linux y macOS.

### 🚀 Instalación

Puedes instalar la aplicación usando `make` estándar o a través del script bash proporcionado.

**Opción 1: Usando Make (Recomendado)**
make install

**Opción 2: Usando el script de instalación**
./install.sh

### ⚙️ Configuración
En su primera ejecución, la aplicación generará automáticamente un archivo de configuración ubicado en:
~/.config/libre-cli/config.toml

Puedes editar este archivo para cambiar tus credenciales de LibreLinkUp, ajustar el intervalo de actualización, definir tu rango objetivo de glucosa (min_glucose / max_glucose) y personalizar completamente los colores del tema.

### 💻 Uso
Simplemente ejecuta el comando desde cualquier terminal:
libre-cli

Presiona q o Ctrl+C para salir de forma segura de la aplicación.
