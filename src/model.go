package src

import tea "github.com/charmbracelet/bubbletea"

type MainMenu struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

type MainMenuSelection interface {
	handleInput(msg tea.Msg, menu MainMenu) Page
}

type Settings struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

type Page interface{}

type model struct {
	page Page
}
