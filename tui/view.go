package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	mem "github.com/kireetivar/topgo/memory"
)

type tickMsg time.Time

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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	case tickMsg:
		m.MemUsagePercent = mem.GetMemoryUsage()
		return m, doTick()
	}
	return m, nil
}

func (m Model) View() string {
	barWidth := m.Width - 12
	if barWidth < 10 {
		barWidth = 10
	}

	label := labelStyle.Render("Mem")
	bar := renderBar(m.MemUsagePercent, barWidth)
	footer := footerStyle.Render("q: quit")
	return fmt.Sprintf("%s [%s] %4.1f%%\n\n%s", label, bar, m.MemUsagePercent, footer)
}
