package main

import (
	"fmt"
	"os"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kireetivar/topgo/tui"
)

func main() {
	if runtime.GOOS != "linux" {
		fmt.Printf("Error: topgo relies on the Linux /proc filesystem.\n")
		fmt.Printf("It cannot run on %s.\n", runtime.GOOS)
		os.Exit(1)
	}

	initial := tui.NewModel()
	p := tea.NewProgram(initial, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
