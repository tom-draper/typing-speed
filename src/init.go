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
	return nil
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
		words:   "The quick brown fox jumps over the lazy dog",
		started: false,
		time: &Time{
			lastUpdated: time.Now(),
			remaining:   30,
		},
	}
}
