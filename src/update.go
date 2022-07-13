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
			if page.time.remaining < 0 {
				m.page = InitMainMenu()
				return m, nil
			}
		}
		return m, tickEvery()
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
			if msg.String() == string(page.words[page.cursor]) {
				page.correct.Push(true)
				page.cursor++
			} else {
				page.correct.Push(false)
				page.cursor++
			}
			print(typing.correct.Length())
		}
	}

	if !page.started {
		page.started = true
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
