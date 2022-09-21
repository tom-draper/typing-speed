package cmd

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

func tickEvery() tea.Cmd {
	// Send a message every second
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width != 0 && msg.Height != 0 {
			m.width = msg.Width
			m.height = msg.Height
		}
		return m, nil
	case TickMsg:
		switch page := m.page.(type) {
		case Typing:
			page.time.remaining--
			if page.time.remaining < 1 {
				m.page = showResults(page)
				return m, nil // Timer finished - stop tick events
			}
			return m, tickEvery()
		}
		return m, nil // No longer on timer page - stop tick events
	}

	switch page := m.page.(type) {
	case MainMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			}
		}
		m.page = page.handleInput(msg, m.config)
	case Typing:
		firstLetter := !page.started
		m.page = page.handleInput(msg, m.config)
		if firstLetter {
			return m, tickEvery() // Start timer once first key pressed
		}
	case Results:
		m.page = page.handleInput(msg, m.config)
	case Settings:
		m.page = page.handleInput(msg, m.config)
	}
	return m, nil
}

func (page MainMenu) handleInput(msg tea.Msg, config Config) Page {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k", "w":
			if page.cursor > 0 {
				page.cursor--
			}

		case "down", "j", "s":
			if page.cursor < len(page.choices)-1 {
				page.cursor++
			}

		case "enter", " ":
			if _, ok := page.selected[page.cursor]; ok {
				delete(page.selected, page.cursor)
			} else {
				switch page.choices[page.cursor] {
				case "Start":
					return InitTyping(config) // Navigate to Typing page
				case "Settings":
					return InitSettings(config) // Navigate to Settings page
				default:
					page.selected[page.cursor] = struct{}{}
				}
			}
		}
	}

	return page
}

func correctWpm(lines []string, correct *Correct, time int) float64 {
	char := 0
	correctWords := 0
	correctWord := true
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines[i]); j++ {
			if char >= correct.Length() {
				break // Reached the last character entered
			} else if !correct.AtIndex(char) {
				// If a letter in current word is incorect, set flag
				correctWord = false
			} else if lines[i][j] == ' ' { // Finished a word
				if correctWord || j == len(lines[i])-1 {
					correctWords++
				} else {
					correctWord = true // Reset flag
				}
			}
			char++
		}
	}
	// If made it to the final character, register the final word
	if correct.Length() == char && correctWord {
		correctWords++
	}

	minutes := float64(time) / 60.0
	return float64(correctWords) / minutes
}

func wordsPerMinFromSec(wps []int, avg_word_len float64) []float64 {
	wpms := make([]float64, len(wps))
	for i := range wps {
		// Multiply second to get minutes, and divide by average word length (5)
		wpms[i] = float64(wps[i]) * 60 / avg_word_len
	}
	return wpms
}

func showResults(page Typing) Results {
	avgWordLen := float64(page.correct.Length()) / float64(page.words)
	wpms := wordsPerMinFromSec(page.wps, avgWordLen)
	wpm := correctWpm(page.lines, page.correct, page.time.limit-page.time.remaining)
	accuracy := float64(page.keystrokesCorrect) / (float64(page.keystrokesCorrect) + float64(page.keystrokesMistakes))
	remainingMistakes := page.correct.Mistakes() // Uncorrected mistakes
	recovery := 1.0
	if page.keystrokesMistakes > 0 {
		recovery = 1.0 - (float64(remainingMistakes) / float64(page.keystrokesMistakes))
	}
	return InitResults(wpms, wpm, accuracy, page.keystrokesCorrect, page.keystrokesMistakes, recovery)
}

func (page Typing) handleInput(msg tea.Msg, config Config) Page {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return InitMainMenu()
		case "backspace", "ctrl+backspace":
			// If on the first char of a line
			if page.cursor == 0 {
				if page.cursorLine > 0 {
					page.cursorLine--
					page.cursor = len(page.lines[page.cursorLine]) - 1
				}
			} else {
				// Register word completed
				if page.lines[page.cursorLine][page.cursor] == ' ' {
					page.words--
				}
				page.correct.Pop()
				page.cursor--
			}
		default:
			// Check if typed last char
			if page.cursorLine == len(page.lines)-1 && page.cursor == len(page.lines[len(page.lines)-1])-1 {
				return showResults(page)
			}
			// If encountered a space, increment words completed
			if page.lines[page.cursorLine][page.cursor] == ' ' {
				page.words++
			}
			// Check whether entered char is correct
			if msg.String() == string(page.lines[page.cursorLine][page.cursor]) {
				page.correct.Push(true)
				page.keystrokesCorrect++
			} else {
				page.correct.Push(false)
				page.keystrokesMistakes++
			}
			page.cursor++
			// Check if move to next line
			if page.cursor >= len(page.lines[page.cursorLine]) {
				page.cursor = 0
				page.cursorLine++
			}
			page.wps[page.time.limit-page.time.remaining]++
		}
	}

	if !page.started {
		page.started = true
	}

	return page
}

func (page Results) handleInput(msg tea.Msg, config Config) Page {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r", "r":
			return InitTyping(config)
		case "ctrl+c", "esc", "enter":
			return InitMainMenu()
		}
	}

	return page
}

func (page Settings) handleInput(msg tea.Msg, config Config) Page {
	// For method consistency, page.selected and config reference the same map
	// Modifications made to config are reflected in page.selected and vice versa
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return InitMainMenu()

		case "up", "k", "w":
			if page.cursor > 0 {
				page.cursor--
			}

		case "down", "j", "s":
			if page.cursor < len(page.choices)-1 {
				page.cursor++
			}

		case "enter", " ":
			switch page.choices[page.cursor] {
			case "Wikipedia":
				// Turn off other word collections
				delete(config.config, 1)
				config.config[page.cursor] = struct{}{} // Turn on wikipedia
			case "Common words":
				// Turn off other word collections
				delete(config.config, 0)
				config.config[page.cursor] = struct{}{} // Turn on common words
			case "30s":
				// Turn off 60s and 120s
				delete(config.config, 3)
				delete(config.config, 4)
				config.config[page.cursor] = struct{}{}
			case "60s":
				// Turn off 60s and 120s
				delete(config.config, 2)
				delete(config.config, 4)
				config.config[page.cursor] = struct{}{}
			case "120s":
				// Turn off 60s and 120s
				delete(config.config, 2)
				delete(config.config, 3)
				config.config[page.cursor] = struct{}{}
			case "Back":
				return InitMainMenu() // Change to main menu page
			default:
				// Toggle option
				if _, ok := config.config[page.cursor]; ok {
					delete(config.config, page.cursor)
				} else {
					config.config[page.cursor] = struct{}{}
				}
			}

		}
	}

	return page
}
