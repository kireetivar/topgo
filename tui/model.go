package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kireetivar/topgo/cpu"
)

type Model struct {
	memUsagePercent float64
	cpuUsagePercent float64
	width           int
	cpuStat         *cpu.CPUStat
	err             error
}

func (m Model) Init() tea.Cmd {
	return doTick()
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
		_ = msg
		return m, m.fetchAllData()
	case dataMsg:
		m.memUsagePercent = msg.memUsagePercent
		m.cpuUsagePercent = msg.cpuUsagePercent
		return m, doTick()
	case errMsg:
		m.err = msg.err
		return m, doTick()
	}
	return m, nil
}

func NewModel() Model {
	return Model{
		memUsagePercent: 0,
		cpuUsagePercent: 0,
		cpuStat:         &cpu.CPUStat{},
	}
}
