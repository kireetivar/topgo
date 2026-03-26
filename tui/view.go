package tui

import (
	"fmt"
	"strings"

	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00ffff"))

	footerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#555555"))
)

func barColour(percent float64) lipgloss.Style {
	switch {
	case percent >= 80:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	case percent >= 60:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff00"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	}
}

func renderBar(percent float64, width int) string {
	filled := int((percent / 100) * float64(width))
	if filled > width {
		filled = width
	}
	filledBar := barColour(percent).Render(strings.Repeat("█", filled))
	emptyBar := strings.Repeat("░", width-filled)
	return filledBar + emptyBar
}

func (m Model) View() string {
	barWidth := m.width - 12
	if barWidth < 10 {
		barWidth = 10
	}

	label := labelStyle.Render("Mem")
	bar := renderBar(m.memUsagePercent, barWidth)
	footer := footerStyle.Render("q: quit")
	return fmt.Sprintf("%s [%s] %4.1f%%\n\n%s", label, bar, m.memUsagePercent, footer)
}
