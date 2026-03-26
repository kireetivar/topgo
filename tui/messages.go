package tui

import (
	"time"
	tea "github.com/charmbracelet/bubbletea"
	mem "github.com/kireetivar/topgo/memory"
)

type tickMsg time.Time

type dataMsg struct {
	memUsagePercent	float64
}

func fetchAllData() tea.Cmd {
	return func() tea.Msg {
		return  dataMsg {
			memUsagePercent : mem.GetMemoryUsage(),
		}
	}
}