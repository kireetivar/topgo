package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	memUsagePercent float64
	width           int
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(doTick(), fetchAllData())
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tickMsg:
		return m, fetchAllData()
	case dataMsg:
		m.memUsagePercent = msg.memUsagePercent
		return m, doTick()
	}
	return m, nil
}

func NewModel() Model {
	return Model{}
}
