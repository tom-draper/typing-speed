package cmd

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
	termbox "github.com/nsf/termbox-go"
)

func terminalDimensions() (int, int) {
	defer termbox.Close()
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	w, h := termbox.Size()
	return w, h
}

func InitialModel() model {
	profile := termenv.ColorProfile()
	foreground := termenv.ForegroundColor()
	w, h := terminalDimensions()

	return model{
		page:   InitMainMenu(),
		width:  w,
		height: h,
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
	settings := Settings{
		choices:  []string{"Wikipedia", "Common words", "Capitalisation", "Punctuation", "Numbers", "Back"},
		selected: make(map[int]struct{}),
	}
	settings.selected[0] = struct{}{}
	settings.selected[2] = struct{}{}
	settings.selected[3] = struct{}{}
	settings.selected[4] = struct{}{}
	return settings
}

func formatText(text string) []string {
	text = strings.ReplaceAll(text, " ", "|")

	ww := wordwrap.NewWriter(50)
	ww.Breakpoints = []rune{'|'}
	ww.KeepNewlines = false
	ww.Write([]byte(text))
	ww.Close()

	text = strings.ReplaceAll(ww.String(), "|", " ")

	chars := strings.Split(text, "")

	return chars
}

func InitTyping() Typing {
	width := 50
	text := wiki_words()
	return Typing{
		chars:    formatText(text),
		correct:  NewCorrect(),
		width:    width,
		started:  false,
		mistakes: 0,
		time: &Time{
			lastUpdated: time.Now(),
			limit:       30,
			remaining:   30,
		},
	}
}

func InitResults(wpm float32, accuracy float32, mistakes int) Results {
	return Results{
		wpm:      wpm,
		accuracy: accuracy,
		mistakes: mistakes,
	}
}
