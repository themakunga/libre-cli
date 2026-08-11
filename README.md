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

### Opción 1: Descargar Binario Precompilado (Recomendado)
No necesitas tener Go instalado. Puedes descargar la versión ya compilada para tu sistema operativo desde la sección de **Releases**:

👉 **[Ir a las Versiones Compiladas (Releases) ➔](https://github.com/TU_USUARIO/libre-cli/releases)**

1. Ve a la pestaña de [Releases](https://github.com/TU_USUARIO/libre-cli/releases) y descarga el archivo comprimido ejecutable correspondiente a tu sistema:
   - **macOS:** `libre-cli_Darwin_x86_64.tar.gz` (Intel) o `libre-cli_Darwin_arm64.tar.gz` (Apple Silicon M1/M2/M3).
   - **Linux:** `libre-cli_Linux_x86_64.tar.gz` o `libre-cli_Linux_arm64.tar.gz`.
   - **Windows:** `libre-cli_Windows_x86_64.zip`.
2. Extrae el archivo comprimido.
3. En macOS / Linux, concede permisos de ejecución y muévelo a tu ruta global (opcional):
   ```bash
   chmod +x libre-cli
   sudo mv libre-cli /usr/local/bin/
   ```
4. Ejecuta el comando en tu terminal:
  ```bash
    libre-cli
  ```

Opción 2: Compilar desde el Código Fuente

Si prefieres compilar la aplicación tú mismo utilizando el compilador de Go (requiere Go 1.20+):

```bash
git clone [https://github.com/TU_USUARIO/libre-cli.git](https://github.com/TU_USUARIO/libre-cli.git)
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
