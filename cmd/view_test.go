package cmd

import (
	"fmt"
	"testing"

	"github.com/muesli/termenv"
)

func fontStyle(str string, colour termenv.Color, faint bool) termenv.Style {
	if faint {
		return termenv.String(str).Foreground(colour).Faint()
	} else {
		return termenv.String(str).Foreground(colour)
	}
}

func TestAllFontColours(t *testing.T) {
	profile := termenv.ColorProfile()
	foreground := termenv.ForegroundColor()
	fmt.Println("Style:", fontStyle("foreground", foreground, false).String())
	fmt.Println("Style:", fontStyle("foreground faint", foreground, true).String())
	for i := 0; i < 13; i++ {
		fmt.Println("Style:", fontStyle(fmt.Sprintf("Colour %d", i), profile.Color(fmt.Sprintf("%d", i)), false).String())
		fmt.Println("Style:", fontStyle(fmt.Sprintf("Colour %d faint", i), profile.Color(fmt.Sprintf("%d", i)), true).String())
	}
}

func TestModelFontColours(t *testing.T) {
	m := InitialModel()
	fmt.Println("Style:", style("correct", m.styles.correct))
	fmt.Println("Style:", style("normal", m.styles.normal))
	fmt.Println("Style:", style("mistakes", m.styles.mistakes))
	fmt.Println("Style:", style("cursor", m.styles.cursor))
	fmt.Println("Style:", style("highlight", m.styles.highlight))
	fmt.Println("Style:", style("title", m.styles.title))
}

func TestFormatChecked(t *testing.T) {
	selected := make(map[int]struct{})
	selected[2] = struct{}{}
	selected[6] = struct{}{}
	for i := 0; i < 7; i++ {
		s := formatChecked(selected, i)
		if (i == 2 || i == 6) && s != "x" {
			t.Errorf("index %d is checked, but checked string is: [%s]", i, s)
		} else if (i != 2 && i != 6) && s != " " {
			t.Errorf("index %d is not checked, but checked string is: [%s]", i, s)
		}
	}
}
