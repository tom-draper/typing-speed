package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	return m.page.view(m.styles, m.width, m.height)
}

func (menu MainMenu) view(styles Styles, width int, height int) string {
	var sb strings.Builder

	title := style("Main Menu\n\n", styles.faintGreen)
	sb.WriteString(title)
	// Send to the UI for rendering
	for i, choice := range menu.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if menu.cursor == i {
			cursor = style(">", styles.greener) // cursor!
		}

		// Render the row
		row := fmt.Sprintf("%s %s\n", cursor, choice)
		sb.WriteString(row)
	}

	exit_instr := style("\nPress Esc to exit.\n", styles.toEnter)
	sb.WriteString(exit_instr)

	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(sb.String())

	return lipgloss.Place(width-10, height, lipgloss.Center, lipgloss.Center, s)
}

func (typing Typing) view(styles Styles, width int, height int) string {
	var sb strings.Builder

	time := style(fmt.Sprintf("%ds", typing.time.remaining), styles.faintGreen)
	sb.WriteString(time)
	sb.WriteString("\n")
	// Ensure full width is used by a line taken to anchor centering
	sb.WriteString(strings.Repeat("", width))
	sb.WriteString("\n")

	charsProcessed := 0
	for i := 0; i < len(typing.lines); i++ {
		if i > typing.cursorLine+4 {
			break
		}
		for j := 0; j < len(typing.lines[i]); j++ {
			if i < typing.cursorLine || (i == typing.cursorLine && j < typing.cursor) {
				// Entered chars
				entered := style(string(typing.lines[i][j]), styles.correct)
				if !typing.correct.AtIndex(charsProcessed) {
					entered = style(string(typing.lines[i][j]), styles.mistakes)
				}
				sb.WriteString(entered)
			} else if j == typing.cursor && i == typing.cursorLine {
				cursor := style(string(typing.lines[i][j]), styles.cursor)
				sb.WriteString(cursor)
			} else {
				toEnter := style(string(typing.lines[i][j]), styles.toEnter)
				sb.WriteString(toEnter)
			}
			charsProcessed++
		}
		sb.WriteString("\n")
	}

	// var entered strings.Builder
	// for i := 0; i < typing.correct.Length(); i++ {
	// 	if typing.correct.AtIndex(i) {
	// 		entered.WriteString(style(typing.lines[i], styles.correct))
	// 	} else {
	// 		entered.WriteString(style(typing.lines[i], styles.mistakes))
	// 	}
	// }
	// sb.WriteString(entered.String())

	// if typing.cursor < len(typing.lines) {
	// 	// Cursor
	// 	cursor := style(typing.lines[typing.cursor], styles.cursor)
	// 	sb.WriteString(cursor)
	// 	// To enter
	// 	newlines := 3
	// 	if typing.cursorLine < 2 {
	// 		newlines = 4
	// 	}
	// 	endingIdx := endingIdx(typing.lines, typing.cursor, newlines)
	// 	toEnter := style(strings.Join(typing.lines[typing.cursor+1:endingIdx], ""), styles.toEnter)
	// 	sb.WriteString(toEnter)
	// }

	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(sb.String())

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, s)
}

func (results Results) view(styles Styles, width int, height int) string {
	var sb strings.Builder

	title := style("Result\n", styles.faintGreen)
	sb.WriteString(title)

	sb.WriteString("\nWPM: ")
	wpm := style(fmt.Sprintf("%.2f", results.wpm), styles.greener)
	sb.WriteString(wpm)

	sb.WriteString("\nAccuracy: ")
	accuracy := style(fmt.Sprintf("%.2f", results.accuracy)+"%", styles.greener)
	sb.WriteString(accuracy)

	sb.WriteString("\nMistakes: ")
	mistakes := style(fmt.Sprintf("%d", results.mistakes), styles.greener)
	sb.WriteString(mistakes)
	sb.WriteString("\n")

	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(sb.String())

	return lipgloss.Place(width-13, height, lipgloss.Center, lipgloss.Center, s)
}

func (settings Settings) view(styles Styles, width int, height int) string {
	var sb strings.Builder

	title := style("Settings\n\n", styles.faintGreen)
	sb.WriteString(title)

	// Iterate over our choices
	for i, choice := range settings.choices {
		var row string
		if choice == "Wikipedia" {
			cursor := " "
			if settings.cursor == i {
				cursor = style(">", styles.greener) // Cursor
			}

			colouredChoice := choice
			_, ok := settings.selected[i]
			if ok {
				colouredChoice = style(choice, styles.greener)
			}
			row = fmt.Sprintf("%s %s\n", cursor, colouredChoice)
		} else if choice == "Common words" {
			cursor := " "
			if settings.cursor == i {
				cursor = style(">", styles.greener) // Cursor
			}

			colouredChoice := choice
			_, ok := settings.selected[i]
			if ok {
				colouredChoice = style(choice, styles.greener)
			}
			row = fmt.Sprintf("%s %s\n\n", cursor, colouredChoice)
		} else if choice == "Back" {
			cursor := "\n "
			if settings.cursor == i {
				cursor = style("\n>", styles.greener) // Cursor
			}
			row = fmt.Sprintf("%s %s\n", cursor, choice)
		} else {
			cursor := " "
			if settings.cursor == i {
				cursor = ">" // Cursor
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

	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(sb.String())

	return lipgloss.Place(width-9, height, lipgloss.Center, lipgloss.Center, s)
}

func style(s string, style Style) string {
	return style(s).String()
}
