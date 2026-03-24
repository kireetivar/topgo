package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	mem "github.com/kireetivar/topgo/memory"
)

type tickMsg time.Time
type model struct {
	memUsagePercent float64
	width           int
}

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

func (m model) Init() tea.Cmd {
	return doTick()
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tickMsg:
		m.memUsagePercent = mem.GetMemoryUsage()
		return m, doTick()
	}
	return m, nil
}

func (m model) View() string {
	barWidth := m.width - 12
	if barWidth < 10 {
		barWidth = 10
	}

	label := labelStyle.Render("Mem")
	bar := renderBar(m.memUsagePercent, barWidth)
	footer := footerStyle.Render("q: quit")
	return fmt.Sprintf("%s [%s] %4.1f%%\n\n%s", label, bar, m.memUsagePercent, footer)
}

func main() {
	memModel := model{
		memUsagePercent: mem.GetMemoryUsage(),
	}
	p := tea.NewProgram(memModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
