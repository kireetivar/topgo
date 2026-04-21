package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kireetivar/topgo/cpu"
	"github.com/kireetivar/topgo/process"
)

type Model struct {
	memUsagePercent float64
	cpuUsagePercent float64
	totalMemory     float64
	width           int
	height          int
	offset          int
	cpuStat         *cpu.CPUStat
	processes       []process.Process
	processTracker  *process.ProcessTracker
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
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "down", "j":
			visibleRows := m.getVisibleRows()
			if m.offset < len(m.processes)-visibleRows {
				m.offset++
			}
		case "up", "k":
			if m.offset > 0 {
				m.offset--
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		return m, m.fetchAllData()
	case dataMsg:
		m.memUsagePercent = msg.memUsagePercent
		m.cpuUsagePercent = msg.cpuUsagePercent
		m.totalMemory = msg.totalMemory
		m.processes = msg.processes
		maxOffset := max(len(m.processes)-m.getVisibleRows(), 0)
		if m.offset > maxOffset {
			m.offset = maxOffset
		}
		return m, doTick()
	case errMsg:
		m.err = msg.err
		return m, doTick()
	}
	return m, nil
}

func (m Model) getVisibleRows() int {
	return m.height - 6
}

func NewModel() Model {
	return Model{
		memUsagePercent: 0,
		cpuUsagePercent: 0,
		cpuStat:         &cpu.CPUStat{},
		processTracker:  process.NewProcessTracker(),
	}
}
