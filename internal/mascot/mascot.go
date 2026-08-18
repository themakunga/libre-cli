// Package mascot implements "Gotita", the animated mascot for libre-cli.
// Gotita reacts to the user's glucose levels with different expressions
// and animation speeds.
package mascot

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// State represents the mascot's emotional state.
type State int

const (
	StateLoading      State = iota // fetching data from the API
	StateNormal                    // glucose within the configured range
	StateWarning                   // slightly outside the range
	StateCriticalHigh              // dangerously high glucose
	StateCriticalLow               // dangerously low glucose
)

// GetState derives the mascot state from the current glucose reading.
//
// Thresholds:
//   - CriticalLow  : glucose < minGlucose - 15
//   - CriticalHigh : glucose > maxGlucose + 40
//   - Warning       : glucose outside [minGlucose, maxGlucose]
//   - Normal        : glucose inside [minGlucose, maxGlucose]
func GetState(glucose, minGlucose, maxGlucose float64, loading bool) State {
	if loading {
		return StateLoading
	}
	if glucose > 0 && glucose < minGlucose-15 {
		return StateCriticalLow
	}
	if glucose > maxGlucose+40 {
		return StateCriticalHigh
	}
	if glucose < minGlucose || glucose > maxGlucose {
		return StateWarning
	}
	return StateNormal
}

// AnimInterval returns the animation tick duration for a given state.
// Critical states animate faster to convey urgency.
func AnimInterval(s State) time.Duration {
	switch s {
	case StateCriticalHigh, StateCriticalLow:
		return 350 * time.Millisecond
	case StateWarning:
		return 600 * time.Millisecond
	case StateLoading:
		return 500 * time.Millisecond
	default: // StateNormal
		return 800 * time.Millisecond
	}
}

// ColorFor returns the appropriate theme color for a given mascot state.
func ColorFor(s State, good, warning, critical, accent string) lipgloss.Color {
	switch s {
	case StateNormal:
		return lipgloss.Color(good)
	case StateWarning:
		return lipgloss.Color(warning)
	case StateCriticalHigh, StateCriticalLow:
		return lipgloss.Color(critical)
	default: // StateLoading
		return lipgloss.Color(accent)
	}
}

// Render returns a styled 4-line mascot block.
//
// The block is always 7 characters wide:
//
//	╭─────╮   ← border top
//	│ ^_^ │   ← face (changes with state + frame)
//	╰─────╯   ← border bottom
//	 Gotita    ← label (changes with state, blinks on critical)
func Render(s State, frame int, color lipgloss.Color) string {
	frame = frame % 4

	var face, label string

	switch s {
	case StateLoading:
		// Pulsing eyes: dim → medium → bright → medium
		eyes := []string{". .", "o o", "O O", "o o"}
		face = fmt.Sprintf("│ %s │", eyes[frame])
		label = "  ...  "

	case StateNormal:
		// Occasional blink on frame 3
		if frame == 3 {
			face = "│ -_- │"
		} else {
			face = "│ ^_^ │"
		}
		label = " Gotita"

	case StateWarning:
		// Alternating worried eyes
		if frame%2 == 0 {
			face = "│ o_o │"
		} else {
			face = "│ O_o │"
		}
		label = " ¡Ojo! "

	case StateCriticalHigh:
		// Fast blinking alarmed face + label
		if frame%2 == 0 {
			face = "│ >_< │"
			label = "¡ALTO! "
		} else {
			face = "│ x_x │"
			label = "  !!!  "
		}

	case StateCriticalLow:
		// Fast blinking alarmed face + label
		if frame%2 == 0 {
			face = "│ >_< │"
			label = "¡BAJO! "
		} else {
			face = "│ x_x │"
			label = "  !!!  "
		}
	}

	bodyStyle := lipgloss.NewStyle().Foreground(color)
	labelStyle := lipgloss.NewStyle().Foreground(color).Bold(true)

	lines := []string{
		bodyStyle.Render("╭─────╮"),
		bodyStyle.Render(face),
		bodyStyle.Render("╰─────╯"),
		labelStyle.Render(label),
	}

	return strings.Join(lines, "\n")
}
