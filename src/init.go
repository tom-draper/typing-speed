package src

import tea "github.com/charmbracelet/bubbletea"

func InitialModel() model {
	return model{
		page: MainMenu{
			choices: []string{"Option 1", "Option 2", "Option 2", "Option 3"},

			// A map which indicates which choices are selected. We're using
			// the  map like a mathematical set. The keys refer to the indexes
			// of the `choices` slice, above.
			selected: make(map[int]struct{}),
		},
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
