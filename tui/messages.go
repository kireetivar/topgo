package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kireetivar/topgo/memory"
)

type tickMsg time.Time

type dataMsg struct {
	memUsagePercent float64
	cpuUsagePercent float64
}

type errMsg struct{ err error }

func (m Model) fetchAllData() tea.Cmd {
	return func() tea.Msg {
		memUsage, err := memory.GetMemoryUsage()
		if err != nil {
			return errMsg{err: err}
		}
		cpuUsage, err := m.cpuStat.GetCPUUsage()
		if err != nil {
			return errMsg{err: err}
		}
		return dataMsg{
			memUsagePercent: memUsage,
			cpuUsagePercent: cpuUsage,
		}
	}
}
