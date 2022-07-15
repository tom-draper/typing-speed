package cmd

import (
	"fmt"
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
	switch msg.(type) {
	case TickMsg:
		switch page := m.page.(type) {
		case Typing:
			page.time.remaining--
			if page.time.remaining < 1 {
				m.page = finished(page)
				return m, nil
			}
			return m, tickEvery()
		}
		return m, nil
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
				// Navigate to a new page
				case "Start":
					return InitTyping(config)
				case "Settings":
					return InitSettings(config)
				default:
					page.selected[page.cursor] = struct{}{}
				}
			}
		}
	}

	return page
}

func correct_wpm(lines []string, correct *Correct, time int) float32 {
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
	if correct.Length() == char && correct_word {
		// If made it to the final character, register the final word
		correct_words++
	}

	minutes := float32(time) / 60.0
	return float32(correct_words) / minutes
}

func finished(page Typing) Results {
	wpm := correct_wpm(page.lines, page.correct, page.time.limit-page.time.remaining)
	accuracy := (float32(page.nCorrect) / (float32(page.nCorrect) + float32(page.nMistakes))) * 100.0
	return InitResults(wpm, accuracy, page.nMistakes)
}

func (typing Typing) handleInput(msg tea.Msg, page Typing, config map[int]struct{}) Page {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return InitMainMenu()
		case "backspace":
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
			if page.cursorLine == len(page.lines)-1 && page.cursor == len(page.lines[len(page.lines)-1])-1 {
				return finished(page) // If finished typing all chars
			}
			if msg.String() == string(page.lines[page.cursorLine][page.cursor]) {
				page.correct.Push(true)
				fmt.Print(page.correct.Length())
				page.nCorrect++
			} else {
				page.correct.Push(false)
				page.nMistakes++
				fmt.Print(page.correct.Length())
			}
			page.cursor++
			if page.cursor >= len(page.lines[page.cursorLine]) {
				// Move to next line
				page.cursor = 0
				page.cursorLine++
			}
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
				// Turn off other options (mutually exclusive)
				if _, ok := page.selected[1]; ok {
					delete(page.selected, 1)
				}
				page.selected[page.cursor] = struct{}{} // Turn on wikipedia
			case "Common words":
				// Turn off other options (mutually exclusive)
				if _, ok := page.selected[0]; ok {
					delete(page.selected, 0)
				}
				page.selected[page.cursor] = struct{}{} // Turn on common words
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
