package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kireetivar/topgo/cpu"
	"github.com/kireetivar/topgo/memory"
	"github.com/kireetivar/topgo/process"
)

type tickMsg time.Time

type dataMsg struct {
	memUsagePercent float64
	cpuUsagePercent float64
	totalMemory     float64
	processes       []process.Process
}

type errMsg struct{ err error }

func (m Model) fetchAllData() tea.Cmd {
	return func() tea.Msg {
		memUsage, memTotal, err := memory.GetMemoryUsage()
		if err != nil {
			return errMsg{err: err}
		}
		curtotal, curidle, err := cpu.ReadTotalCPUTicks()
		if err != nil {
			return errMsg{err: err}
		}
		cpuUsage, err := m.cpuStat.GetCPUUsage(curtotal, curidle)
		if err != nil {
			return errMsg{err: err}
		}
		processes, err := m.processTracker.GetProcessList(curtotal, m.sortBy)
		if err != nil {
			return errMsg{err: err}
		}
		return dataMsg{
			memUsagePercent: memUsage,
			cpuUsagePercent: cpuUsage,
			totalMemory:     memTotal,
			processes:       processes,
		}
	}
}
