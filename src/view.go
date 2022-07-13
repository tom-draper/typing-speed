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

func (typing Typing) view(styles Styles) string {
	var sb strings.Builder

	time := style(fmt.Sprintf("\n\n      %ds\n\n      ", typing.time.remaining), styles.faintGreen)
	sb.WriteString(time)

	var entered strings.Builder
	for i := 0; i < typing.correct.Length(); i++ {
		if typing.correct.AtIndex(i) {
			entered.WriteString(style(typing.words[i], styles.correct))
		} else {
			entered.WriteString(style(typing.words[i], styles.mistakes))
		}
	}
	sb.WriteString(entered.String())

	if typing.cursor < len(typing.words) {
		cursor := style(typing.words[typing.cursor], styles.cursor)
		sb.WriteString(cursor)
		toEnter := style(strings.Join(typing.words[typing.cursor+1:], ""), styles.toEnter)
		sb.WriteString(toEnter)
	}

	return sb.String()
}

func (results Results) view(styles Styles) string {
	var sb strings.Builder

	title := style("\n\n      Result\n", styles.faintGreen)
	sb.WriteString(title)

	sb.WriteString("\n      wpm: ")
	wpm := style(fmt.Sprintf("%.2f", results.wpm), styles.greener)
	sb.WriteString(wpm)

	sb.WriteString("\n      mistakes: ")
	mistakes := style(fmt.Sprintf("%d", results.mistakes), styles.greener)
	sb.WriteString(mistakes)

	return sb.String()
}

func (settings Settings) view(styles Styles) string {
	var sb strings.Builder

	title := style("\n\n      Settings\n\n", styles.faintGreen)
	sb.WriteString(title)

	// Iterate over our choices
	for i, choice := range settings.choices {
		var row string
		if choice == "Back" {
			cursor := "\n     "
			if settings.cursor == i {
				cursor = style("\n    >", styles.greener) // Cursor
			}
			row = fmt.Sprintf("%s %s\n", cursor, choice)
		} else {
			cursor := "     "
			if settings.cursor == i {
				cursor = "    >" // Cursor
			}

			checked := " "
			if _, ok := settings.selected[i]; ok {
				checked = "x" // Choice selected
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
