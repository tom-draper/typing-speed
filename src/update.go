package src

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
	switch msg.(type) {
	case TickMsg:
		switch page := m.page.(type) {
		case Typing:
			page.time.remaining--
			return m, tickEvery()
		default:
			return m, nil
		}
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
		m.page = page.handleInput(msg, page)
		return m, nil
	case Typing:
		firstLetter := !page.started
		m.page = page.handleInput(msg, page)
		if firstLetter {
			return m, tickEvery() // Start timer once first key pressed
		}
		return m, nil
	case Results:
		m.page = page.handleInput(msg, page)
		return m, nil
	case Settings:
		m.page = page.handleInput(msg, page)
		return m, nil
	default:
		return m, nil
	}
}

func (menu MainMenu) handleInput(msg tea.Msg, page MainMenu) Page {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if page.cursor > 0 {
				page.cursor--
			}

		case "down", "j":
			if page.cursor < len(page.choices)-1 {
				page.cursor++
			}

		case "enter", " ":
			_, ok := page.selected[page.cursor]
			if ok {
				delete(page.selected, page.cursor)
			} else {
				switch page.choices[page.cursor] {
				// Navigate to a new page
				case "Start":
					return InitTyping()
				case "Settings":
					return InitSettings()
				default:
					page.selected[page.cursor] = struct{}{}
				}
			}
		}
	}

	return page
}

func correct_wpm(chars []string, correct *Correct, time int) float32 {
	correct_words := 0
	correct_word := true
	for i := range chars {
		if i < correct.Length() && !correct.AtIndex(i) {
			// If a letter in current word is incorect, set flag
			correct_word = false
		} else if chars[i] == " " {
			// Finished word
			if correct_word {
				correct_words++
			}
			correct_word = true // Reset
		}
	}
	// Register the final word
	if correct_word {
		correct_words++
	}
	minutes := float32(time) / 60.0
	return float32(correct_words) / minutes

}

func (typing Typing) handleInput(msg tea.Msg, page Typing) Page {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return InitMainMenu()
		case "backspace":
			page.correct.Pop()
			page.cursor--
		default:
			if page.cursor >= len(page.words)-1 {
				// If finished typing all chars
				wpm := correct_wpm(page.words, page.correct, page.time.limit-page.time.remaining)
				return InitResults(wpm, page.mistakes)
			}
			if msg.String() == string(page.words[page.cursor]) {
				page.correct.Push(true)
				page.cursor++
			} else {
				page.correct.Push(false)
				page.cursor++
				page.mistakes++
			}
		}
	}

	if !page.started {
		page.started = true
	}

	return page
}

func (results Results) handleInput(msg tea.Msg, page Results) Page {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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

		case "up", "k":
			if page.cursor > 0 {
				page.cursor--
			}

		case "down", "j":
			if page.cursor < len(page.choices)-1 {
				page.cursor++
			}

		case "enter", " ":
			_, ok := page.selected[page.cursor]
			if ok {
				delete(page.selected, page.cursor)
			} else {
				switch page.choices[page.cursor] {
				case "Back":
					return InitMainMenu() // Change to main menu page
				default:
					page.selected[page.cursor] = struct{}{}
				}
			}
		}
	}

	return page
}
