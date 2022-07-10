package src

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

// Send a message every second.
func tickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit // Exit program
		}
	case TickMsg:
		switch page := m.page.(type) {
		case Typing:
			page.time.remaining -= 1
			if page.time.remaining < 0 {
				m.page = InitMainMenu()
				return m, nil
			}
		}
		return m, tickEvery()
	}

	switch page := m.page.(type) {
	case MainMenu:
		m.page = page.handleInput(msg, page)
		switch m.page.(type) {
		case Typing:
			return m, tickEvery()
		default:
			return m, nil
		}
	case Typing:
		m.page = page.handleInput(msg, page)
		return m, nil
	case Settings:
		m.page = page.handleInput(msg, page)
		return m, nil
	}

	return m, nil
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

func (t Typing) handleInput(msg tea.Msg, page Typing) Page {

	return page
}

func (s Settings) handleInput(msg tea.Msg, page Settings) Page {
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
