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

func (typing *Typing) runCountdown(m model) {
	for {
		diff := time.Since(typing.timeLastUpdated)
		if diff > time.Second {
			typing.timeRemaining -= 1
			typing.timeLastUpdated = time.Now()
			m.View()
			if typing.timeRemaining == 0 {
				m.Update(tea.Msg("timerfinished"))
			}
		}
	}
}

func InitTyping(m model) Typing {
	typing := Typing{
		choices:         []string{"Option 1", "Option 2", "Option 3", "Back"},
		selected:        make(map[int]struct{}),
		timeLastUpdated: time.Now(),
		timeRemaining:   5,
	}
	go typing.runCountdown(m)
	return typing
}
