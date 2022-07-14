package main

import (
	"fmt"
	"os"
	"typing/cmd"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(cmd.InitialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
