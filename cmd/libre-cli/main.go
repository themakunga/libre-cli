package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"libre-cli/internal/config"
	"libre-cli/internal/ui"
)

func main() {
	// 1. Cargar configuración desde ~/.config/libre-cli/config.toml
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error fatal cargando configuracion: %v\n", err)
		os.Exit(1)
	}

	// 2. Inicializar el modelo de Bubble Tea con la configuración cargada
	appModel := ui.New(cfg)

	// 3. Ejecutar la app en pantalla completa alternativa
	p := tea.NewProgram(appModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error iniciando la app: %v\n", err)
		os.Exit(1)
	}
}
