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
				m.page = finished(page)
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
		m.page = page.handleInput(msg, page, m.config)
	case Typing:
		firstLetter := !page.started
		m.page = page.handleInput(msg, page, m.config)
		if firstLetter {
			return m, tickEvery() // Start timer once first key pressed
		}
	case Results:
		m.page = page.handleInput(msg, page, m.config)
	case Settings:
		m.page = page.handleInput(msg, page)
	}
	return m, nil
}

func (menu MainMenu) handleInput(msg tea.Msg, page MainMenu, config map[int]struct{}) Page {
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

func correct_wpm(lines []string, correct *Correct, time int) float64 {
	char := 0
	correct_words := 0
	correct_word := true
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines[i]); j++ {
			if char >= correct.Length() {
				break // Reached the last character entered
			} else if !correct.AtIndex(char) {
				// If a letter in current word is incorect, set flag
				correct_word = false
			} else if lines[i][j] == ' ' { // Finished a word
				if correct_word || j == len(lines[i])-1 {
					correct_words++
				} else {
					correct_word = true // Reset flag
				}
			}
			char++
		}
	}
	// If made it to the final character, register the final word
	if correct.Length() == char && correct_word {
		correct_words++
	}

	minutes := float64(time) / 60.0
	return float64(correct_words) / minutes
}

func words_per_min_from_sec(wps []int) []float64 {
	wpms := make([]float64, len(wps))
	for i := range wps {
		// Multiply second to get minutes, and divide by average word length (5)
		wpms[i] = float64(wps[i]) * 60 / 5
	}
	return wpms
}

func finished(page Typing) Results {
	wpms := words_per_min_from_sec(page.wps)
	wpm := correct_wpm(page.lines, page.correct, page.time.limit-page.time.remaining)
	accuracy := (float64(page.nCorrect) / (float64(page.nCorrect) + float64(page.nMistakes))) * 100.0
	return InitResults(wpms, wpm, accuracy, page.nMistakes)
}

func (typing Typing) handleInput(msg tea.Msg, page Typing, config map[int]struct{}) Page {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return InitMainMenu()
		case "backspace":
			// If on the first char of a line
			if page.cursor == 0 {
				if page.cursorLine > 0 {
					page.cursorLine--
					page.cursor = len(page.lines[page.cursorLine])
				}
			} else {
				page.correct.Pop()
				page.cursor--
			}
		default:
			// Check if typed last char
			if page.cursorLine == len(page.lines)-1 && page.cursor == len(page.lines[len(page.lines)-1])-1 {
				return finished(page)
			}
			// Check whether input char correct
			if msg.String() == string(page.lines[page.cursorLine][page.cursor]) {
				page.correct.Push(true)
				page.nCorrect++
			} else {
				page.correct.Push(false)
				page.nMistakes++
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

func (results Results) handleInput(msg tea.Msg, page Results, config map[int]struct{}) Page {
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

func (settings Settings) handleInput(msg tea.Msg, page Settings) Page {
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
				delete(page.selected, 1)
				page.selected[page.cursor] = struct{}{} // Turn on wikipedia
			case "Common words":
				// Turn off other word collections
				delete(page.selected, 0)
				page.selected[page.cursor] = struct{}{} // Turn on common words
			case "30s":
				// Turn off 60s and 120s
				delete(page.selected, 3)
				delete(page.selected, 4)
				page.selected[page.cursor] = struct{}{}
			case "60s":
				// Turn off 60s and 120s
				delete(page.selected, 2)
				delete(page.selected, 4)
				page.selected[page.cursor] = struct{}{}
			case "120s":
				// Turn off 60s and 120s
				delete(page.selected, 2)
				delete(page.selected, 3)
				page.selected[page.cursor] = struct{}{}
			case "Back":
				return InitMainMenu() // Change to main menu page
			default:
				// Toggle option
				if _, ok := page.selected[page.cursor]; ok {
					delete(page.selected, page.cursor)
				} else {
					page.selected[page.cursor] = struct{}{}
				}
			}

		}
	}

	return page
}
