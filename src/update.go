package src

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit // Exit program
		}
	}

	if msg == "timerfinished" {
		m.page = InitMainMenu()
		return m, nil
	}

	switch page := m.page.(type) {
	case MainMenu:
		m.page = page.handleInput(msg, page, m)
		return m, nil
	case Typing:
		m.page = page.handleInput(msg, page)
		return m, nil
	case Settings:
		m.page = page.handleInput(msg, page)
		return m, nil
	}

	return m, nil
}

func (m MainMenu) handleInput(msg tea.Msg, page MainMenu, model model) Page {
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
					return InitTyping(model)
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
				page.selected[page.cursor] = struct{}{}
			}
		}
	}

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
