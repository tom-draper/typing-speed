package src

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func InitialModel() model {
	return model{
		page: InitMainMenu(),
	}
}

func (m model) Init() tea.Cmd {
	CallClear()
	// Just return `nil`, which means "no I/O right now, please."
	return tickEvery()
}

func InitMainMenu() MainMenu {
	return MainMenu{
		choices:  []string{"Start", "Settings"},
		selected: make(map[int]struct{}),
	}
}

func InitSettings() Settings {
	return Settings{
		choices:  []string{"Option 1", "Option 2", "Option 3", "Back"},
		selected: make(map[int]struct{}),
	}
}

func InitTyping() Typing {
	return Typing{
		selected: make(map[int]struct{}),
		time: &Time{
			lastUpdated: time.Now(),
			remaining:   5,
		},
	}
}
