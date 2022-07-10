package src

import (
	"fmt"
)

func (m model) View() string {
	return m.page.view()
}

func (m MainMenu) view() string {
	s := "\n\n      Main Menu\n\n"
	// Send to the UI for rendering
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := "     " // no cursor
		if m.cursor == i {
			cursor = "    >" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\n      Press Esc to exit.\n"
	return s
}

func (t Typing) view() string {
	s := "\n\n      " + fmt.Sprint(t.time.remaining) + "s\n\n" + t.words
	return s
}

func (s Settings) view() string {
	str := "\n\n      Settings\n\n"
	// Iterate over our choices
	for i, choice := range s.choices {
		if choice == "Back" {
			// Is the cursor pointing at this choice?
			cursor := "\n     " // no cursor
			if s.cursor == i {
				cursor = "\n    >" // cursor!
			}
			str += fmt.Sprintf("%s %s\n", cursor, choice)
		} else {
			// Is the cursor pointing at this choice?
			cursor := "     " // no cursor
			if s.cursor == i {
				cursor = "    >" // cursor!
			}

			// Is this choice selected?
			checked := " " // not selected
			if _, ok := s.selected[i]; ok {
				checked = "x" // selected!
			}
			// Render the row
			str += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}

	}
	return str
}
