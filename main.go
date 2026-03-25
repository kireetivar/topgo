package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	mem "github.com/kireetivar/topgo/memory"
	"github.com/kireetivar/topgo/tui"
)

func main() {
	memModel := tui.Model{
		MemUsagePercent: mem.GetMemoryUsage(),
	}
	p := tea.NewProgram(memModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
