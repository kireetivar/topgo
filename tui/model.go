package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	mem "github.com/kireetivar/topgo/memory"
)

type Model struct {
	MemUsagePercent float64
	Width           int
}

func (m Model) Init() tea.Cmd {
	return doTick()
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func NewModel() Model {
	return Model{
		MemUsagePercent: mem.GetMemoryUsage(),
		Width:           100,
	}
}
