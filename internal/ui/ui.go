package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"

	"libre-cli/internal/api"
	"libre-cli/internal/config"
	"libre-cli/internal/mascot"
)

type tickMsg time.Time

type ViewMode int

const (
	ModeMain ViewMode = iota
	ModeSensors
	ModeStats
	ModeDebug
)

type Model struct {
	config      config.Config
	glucose     float64
	trend       int
	history     []float64
	sensors     []api.Sensor
	debugData   api.DebugInfo
	lastUpdated time.Time
	loading     bool
	viewMode    ViewMode
	termWidth   int
	termHeight  int
	err         error
	mascotFrame int // current animation frame for Gotita
}

// mascotAnimMsg is fired by the mascot animation ticker.
type mascotAnimMsg time.Time

// mascotTickCmd schedules the next mascot animation frame after duration d.
func mascotTickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return mascotAnimMsg(t)
	})
}

func New(cfg config.Config) Model {
	return Model{
		config:     cfg,
		loading:    true,
		history:    []float64{},
		termWidth:  80,
		termHeight: 24,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchGlucoseCmd(),
		tickEvery(m.config.App.UpdateIntervalMinutes),
		mascotTickCmd(mascot.AnimInterval(mascot.StateLoading)),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "i":
			if m.viewMode == ModeSensors {
				m.viewMode = ModeMain
			} else {
				m.viewMode = ModeSensors
			}
			return m, nil
		case "s":
			if m.viewMode == ModeStats {
				m.viewMode = ModeMain
			} else {
				m.viewMode = ModeStats
			}
			return m, nil
		case "d":
			if m.viewMode == ModeDebug {
				m.viewMode = ModeMain
			} else {
				m.viewMode = ModeDebug
			}
			return m, nil
		}

	case mascotAnimMsg:
		m.mascotFrame++
		state := mascot.GetState(m.glucose, m.config.App.MinGlucose, m.config.App.MaxGlucose, m.loading)
		return m, mascotTickCmd(mascot.AnimInterval(state))

	case tickMsg:
		return m, tea.Batch(
			m.fetchGlucoseCmd(),
			tickEvery(m.config.App.UpdateIntervalMinutes),
		)

	case api.GlucoseData:
		m.glucose = msg.Current
		m.trend = msg.Trend
		m.history = msg.History
		m.sensors = msg.Sensors
		m.debugData = msg.Debug
		m.lastUpdated = time.Now()
		m.loading = false
		return m, nil

	case error:
		m.err = msg
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPresiona 'q' para salir.", m.err)
	}

	theme := m.config.Theme

	if m.loading {
		gotita := mascot.Render(mascot.StateLoading, m.mascotFrame, lipgloss.Color(theme.Accent1))
		loadingMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Accent1)).
			MarginLeft(2).
			Render("Conectando con la\nAPI LibreLinkUp...")
		row := lipgloss.JoinHorizontal(lipgloss.Center, gotita, loadingMsg)
		return lipgloss.NewStyle().Margin(1, 2).Render(row)
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Accent2)).MarginBottom(1)
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Bg)).Background(lipgloss.Color(theme.Accent1)).Padding(0, 1).MarginTop(1)

	footerStr := " [i] Sensores | [s] Stats | [d] Debug | [q] Salir "
	footerBlock := footerStyle.Render(fmt.Sprintf("Actualizado: %s |%s", m.lastUpdated.Format("15:04:05"), footerStr))

	// ==========================================
	// VISTAS SECUNDARIAS (Debug, Stats, Sensors)
	// ==========================================
	if m.viewMode == ModeDebug {
		var debugLines []string
		debugLines = append(debugLines, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Critical)).Render("🛠 Inspector de Datos API"))
		debugLines = append(debugLines, fmt.Sprintf("Base URL: %s", m.debugData.BaseURL))
		debugLines = append(debugLines, fmt.Sprintf("Patient ID: %s", m.debugData.PatientID))
		debugLines = append(debugLines, fmt.Sprintf("HTTP Code: %d", m.debugData.LastHTTPCode))

		tokShort := m.debugData.Token
		if len(tokShort) > 20 {
			tokShort = tokShort[:20] + "..."
		}
		debugLines = append(debugLines, fmt.Sprintf("Token: %s", tokShort))

		jsonSnippet := m.debugData.RawJSON
		if len(jsonSnippet) > 400 {
			jsonSnippet = jsonSnippet[:400] + "\n... [Respuesta Truncada]"
		}

		jsonBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Fg)).
			BorderStyle(lipgloss.NormalBorder()).
			Padding(1).
			Render(jsonSnippet)

		debugLines = append(debugLines, "\nJSON Crudo de Glucosa:\n"+jsonBox)

		box := lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(1, 2).Render(strings.Join(debugLines, "\n"))
		ui := lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render("⚡ Modo Debug"), box, footerBlock)
		return lipgloss.NewStyle().Margin(1, 2).Render(ui)
	}

	if m.viewMode == ModeStats {
		var statsLines []string
		statsLines = append(statsLines, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Accent1)).Render("📊 Estadísticas de Lecturas Reales"))

		if len(m.history) == 0 {
			statsLines = append(statsLines, "No hay suficiente historial cargado.")
		} else {
			var total, min, max float64
			min = m.history[0]
			max = m.history[0]
			inRangeCount := 0

			minRange := m.config.App.MinGlucose
			maxRange := m.config.App.MaxGlucose

			for _, v := range m.history {
				total += v
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
				if v >= minRange && v <= maxRange {
					inRangeCount++
				}
			}

			avg := total / float64(len(m.history))
			tir := (float64(inRangeCount) / float64(len(m.history))) * 100
			eA1c := (avg + 46.7) / 28.7

			statsLines = append(statsLines, fmt.Sprintf("Total Mediciones: %d", len(m.history)))
			statsLines = append(statsLines, fmt.Sprintf("Promedio: %.1f mg/dL", avg))
			statsLines = append(statsLines, fmt.Sprintf("Mínimo: %.1f mg/dL | Máximo: %.1f mg/dL", min, max))
			statsLines = append(statsLines, fmt.Sprintf("Tiempo en Rango (TIR): %.1f%% (%d-%d mg/dL)", tir, int(minRange), int(maxRange)))
			statsLines = append(statsLines, fmt.Sprintf("HbA1c Estimada (eA1c): %.2f%%", eA1c))
		}

		box := lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(1, 2).Render(strings.Join(statsLines, "\n\n"))
		ui := lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render("⚡ Resumen Estadístico"), box, footerBlock)
		return lipgloss.NewStyle().Margin(1, 2).Render(ui)
	}

	if m.viewMode == ModeSensors {
		var sensorLines []string
		sensorLines = append(sensorLines, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Accent1)).Render("Últimos Sensores Registrados:\n"))

		for i, s := range m.sensors {
			statusColor := lipgloss.Color(theme.Good)
			statusText := "Activo"

			if s.DaysLeft <= 0 {
				statusColor = lipgloss.Color(theme.Critical)
				statusText = "Finalizado"
			} else if s.DaysLeft <= 3 {
				statusColor = lipgloss.Color(theme.Warning)
			}

			idStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Accent2)).Render(s.ID)
			daysStyle := lipgloss.NewStyle().Foreground(statusColor).Render(fmt.Sprintf("%d días (%s)", s.DaysLeft, statusText))

			line := fmt.Sprintf("%d. Sensor: %s | Inicio: %s | Estado: %s",
				i+1, idStyle, s.StartDate.Format("02 Jan 2006"), daysStyle)

			sensorLines = append(sensorLines, line)
		}

		box := lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(1, 2).Render(strings.Join(sensorLines, "\n\n"))
		ui := lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render("⚡ Info de Sensores"), box, footerBlock)
		return lipgloss.NewStyle().Margin(1, 2).Render(ui)
	}

	// ==========================================
	// VISTA PRINCIPAL (Alturas fijas por modo)
	// ==========================================
	appCfg := m.config.App
	glucoseColor := lipgloss.Color(theme.Good)
	if m.glucose < appCfg.MinGlucose || m.glucose > appCfg.MaxGlucose {
		glucoseColor = lipgloss.Color(theme.Critical)
	}

	// ── Gotita: mascota reactiva ───────────────────────────────────────────
	mascotState := mascot.GetState(m.glucose, appCfg.MinGlucose, appCfg.MaxGlucose, false)
	mascotColor := mascot.ColorFor(mascotState, theme.Good, theme.Warning, theme.Critical, theme.Accent1)
	mascotBlock := mascot.Render(mascotState, m.mascotFrame, mascotColor)
	availWidth := m.termWidth - 4
	if availWidth < 7 {
		availWidth = 7
	}
	mascotCentered := lipgloss.NewStyle().
		Width(availWidth).
		Align(lipgloss.Center).
		Render(mascotBlock)

	valueStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(glucoseColor).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Accent1)).
		Padding(0, 2)

	trendArrow := "→"
	switch m.trend {
	case 1:
		trendArrow = "↓↓"
	case 2:
		trendArrow = "↓"
	case 3:
		trendArrow = "→"
	case 4:
		trendArrow = "↗"
	case 5:
		trendArrow = "↑↑"
	}

	currentBlock := valueStyle.Render(fmt.Sprintf("%d mg/dL %s", int(math.Round(m.glucose)), trendArrow))

	var content string
	var fixedHeight int

	// LÓGICA RESPONSIVE DE LAYOUT
	if m.termWidth >= 65 {
		// ------------------------------------------
		// MODO HORIZONTAL (Ancho >= 65) -> Altura fija: 5
		// ------------------------------------------
		fixedHeight = 5

		graphText := ""
		if len(m.history) > 0 {
			graph := asciigraph.Plot(m.history,
				asciigraph.Height(3),
				asciigraph.Width(35),
				asciigraph.Caption("Tendencia"),
			)
			graphText = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Fg)).MarginLeft(2).Render(graph)
		}

		content = lipgloss.JoinHorizontal(
			lipgloss.Center,
			currentBlock,
			graphText,
		)

	} else {
		// ------------------------------------------
		// MODO VERTICAL (Ancho < 65) -> Altura fija: 8
		// ------------------------------------------
		fixedHeight = 8

		graphText := ""
		if len(m.history) > 0 {
			graph := asciigraph.Plot(m.history,
				asciigraph.Height(4),
				asciigraph.Width(28),
				asciigraph.Caption("Tendencia"),
			)
			graphText = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Fg)).MarginTop(1).Render(graph)
		}

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			currentBlock,
			graphText,
		)
	}

	// Caja contenedora con altura estricta asignada
	mainBoxStyle := lipgloss.NewStyle().
		Height(fixedHeight).
		Align(lipgloss.Center, lipgloss.Center)

	mainContentBox := mainBoxStyle.Render(content)

	ui := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("⚡ FreeStyle Libre Monitor"),
		mascotCentered,
		mainContentBox,
		footerBlock,
	)

	return lipgloss.NewStyle().Margin(1, 2).Render(ui)
}

func tickEvery(minutes int) tea.Cmd {
	duration := time.Duration(minutes) * time.Minute
	if duration == 0 {
		duration = 5 * time.Minute
	}
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) fetchGlucoseCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := api.Fetch(m.config.App.Email, m.config.App.Password, m.config.App.Region)
		if err != nil {
			return err
		}
		return data
	}
}
