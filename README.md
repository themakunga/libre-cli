# ⚡ libre-cli

Un monitor interactivo en tiempo real para la terminal (TUI) para visualizar datos de glucosa de sensores **FreeStyle Libre** consumidos a través de la API no oficial de **LibreLinkUp**.

![Licencia MIT](https://img.shields.io/badge/License-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8.svg)
![GitHub Release](https://img.shields.io/github/v/release/TU_USUARIO/libre-cli?color=brightgreen)

---

> ⚠️ **AVISO LEGAL Y DESCARGO DE RESPONSABILIDAD (DISCLAIMER)**
>
> 1. **Sin Afiliación Comercial:** Este proyecto **NO** está afiliado, asociado, autorizado, respaldado ni conectado de ninguna manera oficialmente con **Abbott Laboratories**, **FreeStyle Libre**, **LibreLinkUp**, ni con ninguna de sus filiales o subsidiarias. El nombre *FreeStyle Libre*, así como las marcas y logotipos relacionados, son marcas registradas de sus respectivos propietarios.
> 2. **Uso Exclusivamente Personal y Educativo:** Esta herramienta ha sido desarrollada con fines exclusivamente educativos, de aprendizaje técnico y de **uso estrictamente personal**. No debe utilizarse en entornos de producción médica ni para la automatización de decisiones de salud.
> 3. **No es un Dispositivo Médico:** **`libre-cli` NO es un dispositivo médico** ni un software de diagnóstico. No debes tomar decisiones médicas ni de dosificación (como administración de insulina) basándote en los datos visualizados en esta aplicación. Utiliza siempre las aplicaciones oficiales aprobadas por las autoridades de salud y los dispositivos de medición directa recomendados por tu médico.

---

## 🚀 Características
- **Gráfico en tiempo real:** Visualiza la tendencia de tu glucosa mediante gráficos ASCII interactivos (`asciigraph`).
- **Diseño Adaptable (Responsive):** La interfaz reacomoda automáticamente los paneles en horizontal o vertical según la altura de tu consola.
- **Múltiples Vistas:**
  - **Principal:** Nivel de glucosa, flecha de tendencia y gráfico.
  - **Sensores (`i`):** Estado y días restantes del sensor activo.
  - **Estadísticas (`s`):** Tiempo en Rango (TIR), promedios y eA1c estimada.
  - **Debug (`d`):** Inspector del JSON crudo y códigos HTTP.

---

## 📦 Instalación

### Opción 1: Binario precompilado (Recomendado)

No necesitas tener Go instalado. Descarga el binario para tu sistema desde la página de **[Releases](https://github.com/TU_USUARIO/libre-cli/releases/latest)**.

El pipeline genera los siguientes archivos por cada versión:

| Plataforma | Archivo |
|---|---|
| macOS — Apple Silicon (M1/M2/M3) | `libre-cli-darwin-arm64.tar.gz` |
| macOS — Intel | `libre-cli-darwin-amd64.tar.gz` |
| Linux — x86_64 | `libre-cli-linux-amd64.tar.gz` |
| Linux — ARM64 (Raspberry Pi, etc.) | `libre-cli-linux-arm64.tar.gz` |
| Windows — x86_64 | `libre-cli-windows-amd64.zip` |
| Windows — ARM64 | `libre-cli-windows-arm64.zip` |

---

#### 🍎 macOS

```bash
# Apple Silicon (M1/M2/M3)
curl -L https://github.com/TU_USUARIO/libre-cli/releases/latest/download/libre-cli-darwin-arm64.tar.gz | tar xz

# Intel
curl -L https://github.com/TU_USUARIO/libre-cli/releases/latest/download/libre-cli-darwin-amd64.tar.gz | tar xz

# Dar permisos y mover al PATH
chmod +x libre-cli
sudo mv libre-cli /usr/local/bin/

# Primera ejecución: si macOS bloquea el binario por "desarrollador no verificado"
xattr -d com.apple.quarantine /usr/local/bin/libre-cli
```

#### 🐧 Linux

```bash
# x86_64
curl -L https://github.com/TU_USUARIO/libre-cli/releases/latest/download/libre-cli-linux-amd64.tar.gz | tar xz

# ARM64
curl -L https://github.com/TU_USUARIO/libre-cli/releases/latest/download/libre-cli-linux-arm64.tar.gz | tar xz

# Dar permisos y mover al PATH
chmod +x libre-cli
sudo mv libre-cli /usr/local/bin/
```

#### 🪟 Windows

Descarga `libre-cli-windows-amd64.zip` desde [Releases](https://github.com/TU_USUARIO/libre-cli/releases/latest) y descomprime con el Explorador de archivos, o desde PowerShell:

```powershell
# Descomprimir
Expand-Archive libre-cli-windows-amd64.zip -DestinationPath .

# Ejecutar directamente
.\libre-cli.exe

# Opcional: mover a una carpeta incluida en el PATH (ej. C:\Tools)
Move-Item libre-cli.exe C:\Tools\libre-cli.exe
```

> **Nota Windows:** si Windows Defender bloquea el ejecutable, ve a *Configuración → Seguridad de Windows → Protección contra virus → Historial de protección* y permite el archivo manualmente.

#### ✅ Verificar integridad del binario (opcional)

Cada release incluye un archivo `.sha256` junto al binario para verificar que la descarga no fue alterada:

```bash
# Linux / macOS
sha256sum -c libre-cli-linux-amd64.tar.gz.sha256

# Windows (PowerShell)
Get-FileHash libre-cli-windows-amd64.zip -Algorithm SHA256
# Compara el hash con el contenido del archivo .sha256 del release
```

---

### Opción 2: Compilar desde el código fuente

Requiere [Go 1.20+](https://go.dev/dl/):

```bash
git clone https://github.com/TU_USUARIO/libre-cli.git
cd libre-cli
go build -o libre-cli ./cmd/libre-cli
./libre-cli
```

## ⚙️ Configuración

Para utilizar la aplicación debes tener una cuenta activa en la app LibreLinkUp configurada como seguidor del paciente.
Al ejecutar el programa por primera vez, se generará automáticamente un archivo de configuración en `~/.config/libre-cli/config.toml` (o en `%APPDATA%\libre-cli\config.toml` en Windows).
Abre dicho archivo y edítalo con tus credenciales:

```toml
[app]
email = "tu_correo_seguidor@email.com"
password = "tu_password"
region = "cl" # Opciones: cl (Chile/LatAm), eu (Europa), us (EEUU)
update_interval_minutes = 5
min_glucose = 70.0
max_glucose = 180.0
```

## ⌨️ Controles de Teclado

Tecla | Acción
--- | ---
s | "Abre / Cierra la vista de Estadísticas (TIR, Promedio, eA1c)."
i | Abre / Cierra la vista de Sensores activos.
d| Abre / Cierra la vista de Debug / Inspector JSON.
q / Ctrl+C | Sale de la aplicación.

## 🤝 Contribuciones
¡Las contribuciones son bienvenidas! Por favor, lee nuestro Código de Conducta y revisa la Guía de Contribución antes de enviar un Pull Request.

## 📄 Licencia

Este proyecto se distribuye bajo los términos de la Licencia MIT. Consulta el archivo LICENSE para obtener más detalles.
