package src

import "fmt"

func (m model) View() string {

	switch state := m.page.(type) {
	case MainMenu:
		s := "Main Menu"
	case Settings:
		s := "Settings\n\n"
		// Iterate over our choices
		for i, choice := range m.page.choices {

			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.page.cursor == i {
				cursor = ">" // cursor!
			}

			// Is this choice selected?
			checked := " " // not selected
			if _, ok := m.page.selected[i]; ok {
				checked = "x" // selected!
			}

			// Render the row
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}

		s += "\nPress q to quit.\n"
	}

	// Send the UI for rendering
	return s
}
