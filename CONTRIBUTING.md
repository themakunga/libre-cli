# Guía de Contribución a libre-cli

¡Gracias por tu interés en contribuir a `libre-cli`! Este proyecto es de código abierto y cualquier tipo de ayuda (reportar errores, proponer funcionalidades o enviar código) es muy bienvenido.

## ¿Cómo puedo contribuir?

### 1. Reportar un Bug o Problema
Antes de crear un nuevo reporte:
- Revisa las [Issues existentes](https://github.com/TU_USUARIO/libre-cli/issues) para verificar si alguien más ya reportó el problema.
- Si no existe, utiliza la plantilla de **Bug Report** proporcionada e incluye:
  - Pasos exactos para reproducir el fallo.
  - La captura de pantalla o log de la vista de Debug (tecla `d`). *Por favor, oculta tus credenciales o tokens.*

### 2. Proponer una Nueva Funcionalidad
- Abre una Issue usando la plantilla **Feature Request**.
- Explica claramente qué problema resuelve la nueva función y cómo te gustaría que se viera en la interfaz TUI.

### 3. Enviar un Pull Request (PR)
1. **Haz un Fork** del repositorio a tu cuenta.
2. **Crea una rama descriptiva** para tu cambio:
   ```bash
   git checkout -b feat/nueva-vista-de-graficos
   # o
   git checkout -b fix/error-redireccion-region
3. **Sigue las convenciones de Commits (Conventional Commits):**
    Para que el generador automático de Changelog y Releases funcione, nombra tus commits de la siguiente forma:

        fix: corregido cálculo de promedio en estadísticas

        feat: agregada compatibilidad con orden horizontal

        docs: actualizado el archivo README

    Prueba tu código: Asegúrate de ejecutar go fmt o dejar que los hooks de pre-commit validen la calidad del código.

    Abre un Pull Request hacia la rama main del repositorio original.

_Nota: Todos los contribuyentes deben seguir nuestro (./CODE_OF_CONDUCT)[Código de Conducta.]_
