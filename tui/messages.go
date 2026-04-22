package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kireetivar/topgo/cpu"
	"github.com/kireetivar/topgo/memory"
	"github.com/kireetivar/topgo/process"
	"github.com/kireetivar/topgo/sysinfo"
)

type tickMsg time.Time

type dataMsg struct {
	memUsagePercent float64
	cpuUsagePercent float64
	totalMemory     float64
	swapTotal       float64
	swapUsage       float64
	uptime          time.Duration
	loadAvg         [3]float64
	processes       []process.Process
}

type errMsg struct{ err error }

func (m Model) fetchAllData() tea.Cmd {
	return func() tea.Msg {
		memStats, err := memory.GetMemoryUsage()
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
		uptime, err := sysinfo.GetUptime()
		if err != nil {
			return errMsg{err: err}
		}
		load, err := sysinfo.GetLoadAvg()
		if err != nil {
			return errMsg{err: err}
		}
		return dataMsg{
			memUsagePercent: memStats.UsagePercentage,
			cpuUsagePercent: cpuUsage,
			totalMemory:     memStats.TotalGB,
			swapUsage:       memStats.SwapPercentage,
			swapTotal:       memStats.SwapTotalGB,
			uptime:          uptime,
			loadAvg:         load,
			processes:       processes,
		}
	}
}
