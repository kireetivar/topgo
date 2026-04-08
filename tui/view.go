package tui

import (
	"fmt"
	"strings"

	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/kireetivar/topgo/process"
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

func renderProcessTable(processes []process.Process) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%-10s %-20s %-10s %-10s\n", "PID", "Name", "CPU", "MEM")
	for _, proc := range processes {
		fmt.Fprintf(&builder, "%-10d %-20s %-10.1f %10.1f\n", proc.PID, proc.Name, proc.CPU, proc.Mem)
	}
	return builder.String()
}

func (m Model) View() string {
	barWidth := m.width - 12
	if barWidth < 10 {
		barWidth = 10
	}
	if m.err != nil {
		return m.err.Error()
	}
	label := labelStyle.Render("Mem")
	bar := renderBar(m.memUsagePercent, barWidth)
	cpuLabel := labelStyle.Render("CPU")
	cpuBar := renderBar(m.cpuUsagePercent, barWidth)
	header := fmt.Sprintf("%s [%s] %4.1f%%\n\n%s [%s] %4.1f%%", label, bar, m.memUsagePercent, cpuLabel, cpuBar, m.cpuUsagePercent)
	visibleRows := m.height - 6
	var processTable string
	if visibleRows > 0 && len(m.processes) > 0 && m.offset+visibleRows <= len(m.processes) {
		processTable = renderProcessTable(m.processes[m.offset : m.offset+visibleRows])
	} else if visibleRows > 0 && len(m.processes) > 0 {
		processTable = renderProcessTable(m.processes[m.offset:])
	}
	footer := footerStyle.Render("q: quit")
	return lipgloss.JoinVertical(lipgloss.Left, header, processTable, footer)
}
