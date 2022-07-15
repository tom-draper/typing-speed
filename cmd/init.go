package cmd

import (
	"regexp"
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

	// Corresponds to options in Settings page
	config := make(map[int]struct{})
	config[0] = struct{}{} // Wikipedia
	config[2] = struct{}{} // Capitalisation
	config[3] = struct{}{} // Punctuation
	config[4] = struct{}{} // Numbers

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
		config: config,
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

func InitSettings(config map[int]struct{}) Settings {
	settings := Settings{
		choices:  []string{"Wikipedia", "Common words", "Capitalisation", "Punctuation", "Numbers", "Back"},
		selected: config,
	}
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

func applyConfigFilters(text string, config map[int]struct{}) string {
	var ok bool
	_, ok = config[2]
	if !ok {
		// Remove capitalisation
		text = strings.ToLower(text)
	}
	_, ok = config[3]
	if !ok {
		// Remove punctuation
		re := regexp.MustCompile("[!-/:-@[-`{-~.,?<>']")
		text = re.ReplaceAllString(text, "")
	}
	_, ok = config[4]
	if !ok {
		// Remove numbers
		re := regexp.MustCompile(`[0-9]`)
		text = re.ReplaceAllString(text, "")
	}
	return text
}

func InitTyping(config map[int]struct{}) Typing {
	width := 50
	text := wiki_words()

	text = applyConfigFilters(text, config)

	return Typing{
		chars:      formatText(text),
		correct:    NewCorrect(),
		width:      width,
		started:    false,
		cursorLine: 0,
		nMistakes:  0,
		nCorrect:   0,
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
