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
		m.page = page.handleInput(msg, page)
	case Typing:
		firstLetter := !page.started
		m.page = page.handleInput(msg, page)
		if firstLetter {
			return m, tickEvery() // Start timer once first key pressed
		}
	case Results:
		m.page = page.handleInput(msg, page)
	case Settings:
		m.page = page.handleInput(msg, page)
	}
	return m, nil
}

func (menu MainMenu) handleInput(msg tea.Msg, page MainMenu) Page {
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
	for i := 0; i < correct.Length(); i++ {
		if !correct.AtIndex(i) {
			// If a letter in current word is incorect, set flag
			correct_word = false
		} else if chars[i] == " " && correct_word {
			correct_words++
		} else {
			correct_word = true
		}
	}
	if correct.Length() == len(chars) && correct_word {
		// If made it to the final character, register the final word
		correct_words++
	}

	minutes := float32(time) / 60.0
	return float32(correct_words) / minutes

}

func finished(page Typing) Results {
	wpm := correct_wpm(page.chars, page.correct, page.time.limit-page.time.remaining)
	return InitResults(wpm, page.mistakes)
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
			if page.cursor >= len(page.chars)-1 {
				return finished(page) // If finished typing all chars
			}
			if msg.String() == string(page.chars[page.cursor]) {
				page.correct.Push(true)
			} else {
				page.correct.Push(false)
				page.mistakes++
			}
			page.cursor++
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
		case "ctrl+r", "r":
			return InitTyping()
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
