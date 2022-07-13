package src

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func InitialModel() model {
	profile := termenv.ColorProfile()
	foreground := termenv.ForegroundColor()
	return model{
		page: InitMainMenu(),
		styles: Styles{
			correct: func(str string) termenv.Style {
				return termenv.String(str).Foreground(foreground)
			},
			toEnter: func(str string) termenv.Style {
				return termenv.String(str).Foreground(foreground).Faint()
			},
			mistakes: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("1")).Underline()
			},
			cursor: func(str string) termenv.Style {
				return termenv.String(str).Reverse().Bold()
			},
			runningTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2"))
			},
			stoppedTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2")).Faint()
			},
			greener: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("6")).Faint()
			},
			faintGreen: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("10")).Faint()
			},
		},
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
		words:   strings.Split("The quick brown fox jumps over the lazy dog", ""),
		started: false,
		time: &Time{
			lastUpdated: time.Now(),
			remaining:   30,
		},
	}
}
