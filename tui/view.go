package tui

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kireetivar/topgo/process"
)

var (
	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00ffff"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#ffffff"))

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444"))

	footerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#555555"))

	errStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ff0000"))
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

var numCPU = runtime.NumCPU()

func loadColor(val float64) lipgloss.Style {
	ratio := val / float64(numCPU)
	switch {
	case ratio >= 1.0:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	case ratio >= 0.7:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff00"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	}
}

func renderProcessTable(processes []process.Process) string {
	var builder strings.Builder

	for _, proc := range processes {
		fdDisplay := fmt.Sprintf("%d", proc.FDs)
		if proc.FDs == -1 {
			fdDisplay = "-"
		}

		fmt.Fprintf(&builder, "%-8d %-20s %4s %5d %8.1f %8.1f %4s   %-s\n",
			proc.PID, proc.Name, proc.State, proc.NumThreads, proc.CPU, proc.Mem, fdDisplay, proc.Cmd)
	}
	return builder.String()
}

func (m Model) View() string {
	statsWidth := 16
	barWidth := m.width - 6 - statsWidth // 6 = label(3) + " [" + "] "
	if barWidth < 10 {
		barWidth = 10
	}
	label := labelStyle.Render("Mem")
	bar := renderBar(m.memUsagePercent, barWidth)
	usedMemory := (m.memUsagePercent / 100) * m.totalMemory
	memStats := fmt.Sprintf("%.1f/%.1f GB", usedMemory, m.totalMemory)

	cpuLabel := labelStyle.Render("CPU")
	cpuBar := renderBar(m.cpuUsagePercent, barWidth)
	cpuStats := fmt.Sprintf("%-16s", fmt.Sprintf("%.1f%%", m.cpuUsagePercent))

	swpLabel := labelStyle.Render("Swp")
	swpBar := renderBar(m.swapPercentage, barWidth)
	usedSwap := (m.swapPercentage / 100) * m.swapTotal
	swpStats := fmt.Sprintf("%.1f/%.1f GB", usedSwap, m.swapTotal)

	header := fmt.Sprintf("%s [%s] %s\n%s [%s] %s\n%s [%s] %s", label, bar, memStats, swpLabel, swpBar, swpStats, cpuLabel, cpuBar, cpuStats)

	hours := int(m.uptime.Hours())
	days := hours / 24
	hours = hours % 24
	mins := int(m.uptime.Minutes()) % 60
	uptimeStr := infoStyle.Render(fmt.Sprintf("up %dd %dh %dm", days, hours, mins))

	loadStr := infoStyle.Render("Load: ") + loadColor(m.loadAvg[0]).Render(fmt.Sprintf("%.2f", m.loadAvg[0])) + "  " + loadColor(m.loadAvg[1]).Render(fmt.Sprintf("%.2f", m.loadAvg[1])) + "  " +
		loadColor(m.loadAvg[2]).Render(fmt.Sprintf("%.2f", m.loadAvg[2]))

	infoLine := uptimeStr + "    " + loadStr

	visibleRows := m.getVisibleRows()
	rawHeader := fmt.Sprintf("%-8s %-20s %4s %5s %8s %8s %4s   %-s", "PID", "Name", "ST", "THR", "CPU", "MEM", "FD", "CMD")
	tableHeader := tableHeaderStyle.Render(rawHeader)
	separator := separatorStyle.Render(strings.Repeat("─", len(rawHeader)))
	var processTable string
	if visibleRows > 0 && len(m.processes) > 0 && m.offset+visibleRows <= len(m.processes) {
		processTable = renderProcessTable(m.processes[m.offset : m.offset+visibleRows])
	} else if visibleRows > 0 && len(m.processes) > 0 {
		processTable = renderProcessTable(m.processes[m.offset:])
	}
	sortIndicator := "cpu"
	if m.sortBy == process.SortByMem {
		sortIndicator = "mem"
	}
	footer := footerStyle.Render(fmt.Sprintf("q: quit  c/m: sort by [%s]", sortIndicator))

	if m.err != nil {
		footer += "   " + errStyle.Render("ERROR: "+m.err.Error())
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, infoLine, "", tableHeader, separator, processTable, footer)
}
