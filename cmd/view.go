package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
)

func (m model) View() string {
	return m.page.view(m.styles, m.width, m.height)
}

func (menu MainMenu) view(styles Styles, width int, height int) string {
	var sb strings.Builder

	title := style("  Main Menu\n\n", styles.title)
	sb.WriteString(title)

	for i, choice := range menu.choices {
		cursor := formatCursor(menu.cursor, i, styles)
		row := fmt.Sprintf("%s %s\n", cursor, choice)
		sb.WriteString(row)
	}

	exit_instr := style("\n  Press Esc to exit.\n", styles.normal)
	sb.WriteString(exit_instr)

	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(sb.String())

	return lipgloss.Place(width-9, height, lipgloss.Center, lipgloss.Center, s)
}

func (typing Typing) view(styles Styles, width int, height int) string {
	var sb strings.Builder

	time := style(fmt.Sprintf("%ds", typing.time.remaining), styles.title)
	sb.WriteString(time)
	sb.WriteString("\n")
	// Ensure full width is used by a line taken to anchor centering
	sb.WriteString(strings.Repeat(" ", typing.maxLineLen))
	sb.WriteString("\n")

	charsProcessed := 0
	for i := 0; i < len(typing.lines); i++ {
		if distantPastLine(i, typing.cursorLine) {
			break // Skip printing lines to enter far away from current line
		}
		// Insert char-by-char from current line
		for j := 0; j < len(typing.lines[i]); j++ {
			if distantFutureLine(i, typing.cursorLine) {
				charsProcessed += len(typing.lines[i])
				break // Skip printing entered lines more then 2 away from current line
			} else if i < typing.cursorLine || (i == typing.cursorLine && j < typing.cursor) {
				// Entered chars
				entered := style(string(typing.lines[i][j]), styles.green)
				if !typing.correct.AtIndex(charsProcessed) {
					entered = style(string(typing.lines[i][j]), styles.redUnderline)
				}
				sb.WriteString(entered)
			} else if j == typing.cursor && i == typing.cursorLine {
				// Cursor
				cursor := style(string(typing.lines[i][j]), styles.cursor)
				sb.WriteString(cursor)
			} else {
				// Chars to enter
				toEnter := style(string(typing.lines[i][j]), styles.normal)
				sb.WriteString(toEnter)
			}
			charsProcessed++
			if j == len(typing.lines[i])-1 {
				sb.WriteString("\n")
			}
		}
	}

	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(sb.String())

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, s)
}

func distantPastLine(line int, cursorLine int) bool {
	return (cursorLine == 0 && line > cursorLine+3) || (cursorLine > 0 && line > cursorLine+2)
}

func distantFutureLine(line int, cursorLine int) bool {
	return line < cursorLine-1
}

func (results Results) view(styles Styles, width int, height int) string {
	var sb strings.Builder

	title := style("Results\n", styles.title)
	sb.WriteString(title)

	graph := plotWpms(results.wpms, width)
	sb.WriteString(graph)

	sb.WriteString("\nWPM: ")
	wpm := style(fmt.Sprintf("%.2f", results.wpm), styles.highlight)
	sb.WriteString(wpm)

	sb.WriteString("   Accuracy: ")
	accuracy := style(fmt.Sprintf("%.2f", results.accuracy*100)+"%", styles.highlight)
	sb.WriteString(accuracy)

	sb.WriteString("   Keystrokes: ")
	correct := style(fmt.Sprintf("%d", results.keystrokesCorrect), styles.highlight)
	sb.WriteString(correct)

	divider := style(" | ", styles.normal)
	sb.WriteString(divider)

	mistakes := style(fmt.Sprintf("%d", results.keystrokesMistakes), styles.red)
	sb.WriteString(mistakes)

	sb.WriteString("   Recovery: ")
	recovery := style(fmt.Sprintf("%.2f", results.recovery*100)+"%", styles.highlight)
	sb.WriteString(recovery)

	performanceLabel := "\n\nPerformance: "
	totalBars := 58
	// totalBars := int(float64(width)*60) - len(performanceLabel)
	bars := int(results.performance * float64(totalBars))
	sb.WriteString(performanceLabel)
	sb.WriteString(style(strings.Repeat("|", bars), styles.highlight))
	sb.WriteString(style(strings.Repeat("|", totalBars-bars), styles.normal))

	restart := style("\n\nPress r to restart.", styles.normal)
	sb.WriteString(restart)

	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(sb.String())

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, s)
}

func plotWpms(wpms []float64, width int) string {
	wpmGraph := asciigraph.Plot(
		wpms,
		asciigraph.Precision(0),
		asciigraph.Height(6),
		asciigraph.Width(63),
		asciigraph.CaptionColor(2),
		asciigraph.LabelColor(2),
	)

	return lipgloss.NewStyle().Padding(1, 0, 1, 1).Render(wpmGraph)
}

func (settings Settings) view(styles Styles, width int, height int) string {
	var sb strings.Builder

	title := style("  Settings\n\n", styles.title)
	sb.WriteString(title)

	// Iterate over our choices
	for i, choice := range settings.choices {
		var row string
		if choice == "Wikipedia" {
			cursor := formatCursor(settings.cursor, i, styles)
			colouredChoice := formatColouredChoice(choice, settings.selected, i, styles)
			row = fmt.Sprintf("%s %s\n", cursor, colouredChoice)
		} else if choice == "Common words" {
			cursor := formatCursor(settings.cursor, i, styles)
			colouredChoice := formatColouredChoice(choice, settings.selected, i, styles)
			row = fmt.Sprintf("%s %s\n\n", cursor, colouredChoice)
		} else if choice == "30s" {
			cursor := formatCursor(settings.cursor, i, styles)
			colouredChoice := formatColouredChoice(choice, settings.selected, i, styles)
			row = fmt.Sprintf("  Time:  %s %s", cursor, colouredChoice)
		} else if choice == "60s" {
			cursor := formatCursor(settings.cursor, i, styles)
			colouredChoice := formatColouredChoice(choice, settings.selected, i, styles)
			row = fmt.Sprintf("  %s %s", cursor, colouredChoice)
		} else if choice == "120s" {
			cursor := formatCursor(settings.cursor, i, styles)
			colouredChoice := formatColouredChoice(choice, settings.selected, i, styles)
			row = fmt.Sprintf("  %s %s\n\n", cursor, colouredChoice)
		} else if choice == "Back" {
			cursor := formatCursor(settings.cursor, i, styles)
			row = fmt.Sprintf("\n%s %s\n", cursor, choice)
		} else {
			cursor := formatCursor(settings.cursor, i, styles)
			checked := formatChecked(settings.selected, i)
			row = fmt.Sprintf("%s [%s] %s\n", style(cursor, styles.highlight), checked, choice)
		}
		sb.WriteString(row)
	}

	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(sb.String())

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, s)
}

func formatCursor(cursor int, current int, styles Styles) string {
	cursorStr := " " // No cursor
	if cursor == current {
		cursorStr = style(">", styles.highlight) // Cursor
	}
	return cursorStr
}

func formatChecked(selected map[int]struct{}, idx int) string {
	checked := " "
	if _, ok := selected[idx]; ok {
		checked = "x" // Choice selected
	}
	return checked
}

func formatColouredChoice(choice string, selected map[int]struct{}, idx int, styles Styles) string {
	colouredChoice := choice
	if _, ok := selected[idx]; ok {
		colouredChoice = style(choice, styles.highlight)
	}
	return colouredChoice
}

func style(s string, style Style) string {
	return style(s).String()
}
