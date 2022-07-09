package main

import (
	"fmt"
	"os"
	"typing/src"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(src.InitialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
