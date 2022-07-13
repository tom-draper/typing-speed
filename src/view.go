package src

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	return m.page.view(m.styles)
}

func (menu MainMenu) view(styles Styles) string {
	var sb strings.Builder

	title := style("\n\n      Main Menu\n\n", styles.faintGreen)
	sb.WriteString(title)
	// Send to the UI for rendering
	for i, choice := range menu.choices {

		// Is the cursor pointing at this choice?
		cursor := "     " // no cursor
		if menu.cursor == i {
			cursor = style("    >", styles.greener) // cursor!
		}

		// Render the row
		row := fmt.Sprintf("%s %s\n", cursor, choice)
		sb.WriteString(row)
	}

	exit_instr := style("\n      Press Esc to exit.\n", styles.toEnter)
	sb.WriteString(exit_instr)

	return sb.String()
}

func (t Typing) view(styles Styles) string {
	var sb strings.Builder

	time := style(fmt.Sprintf("\n\n      %ds\n\n      ", t.time.remaining), styles.faintGreen)
	sb.WriteString(time)

	entered := style("the quick brown fox ju", styles.correct)
	sb.WriteString(entered)

	cursor := style("m", styles.cursor)
	sb.WriteString(cursor)

	toEnter := style("ped over the lazy dog", styles.toEnter)
	// toEnter := style(fmt.Sprintf("\n\n      %s", t.words), styles.toEnter)
	sb.WriteString(toEnter)

	return sb.String()
}

func (s Settings) view(styles Styles) string {
	var sb strings.Builder

	title := style("\n\n      Settings\n\n", styles.faintGreen)
	sb.WriteString(title)

	// Iterate over our choices
	for i, choice := range s.choices {
		var row string
		if choice == "Back" {
			// Is the cursor pointing at this choice?
			cursor := "\n     " // no cursor
			if s.cursor == i {
				cursor = style("\n    >", styles.greener) // cursor!
			}
			row = fmt.Sprintf("%s %s\n", cursor, choice)
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
			row = fmt.Sprintf("%s [%s] %s\n", style(cursor, styles.greener), checked, choice)
		}
		sb.WriteString(row)

	}
	return sb.String()
}

func style(s string, style Style) string {
	return style(s).String()
}
