package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(NewApp(), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		os.Exit(1)
	}
}
