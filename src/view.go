package src

import "fmt"

func (m model) View() string {
	var s string

	switch page := m.page.(type) {
	case MainMenu:
		s := "Main Menu"
		// Send to the UI for rendering
		return s
	case Settings:
		s := "Settings\n\n"
		// Iterate over our choices
		for i, choice := range page.choices {

			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if page.cursor == i {
				cursor = ">" // cursor!
			}

			// Is this choice selected?
			checked := " " // not selected
			if _, ok := page.selected[i]; ok {
				checked = "x" // selected!
			}

			// Render the row
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}

		s += "\nPress q to quit.\n"

		// Send to the UI for rendering
		return s
	}

	return s
}
