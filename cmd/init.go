package cmd

import (
	"errors"
	"math"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
	termbox "github.com/nsf/termbox-go"
)

func terminalDimensions() (int, int) {
	defer termbox.Close()
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	w, h := termbox.Size()
	return w, h
}

func InitialModel() model {
	profile := termenv.ColorProfile()
	foreground := termenv.ForegroundColor()
	w, h := terminalDimensions()

	// Enable default settings - corresponds to options in Settings page
	config := make(map[int]struct{})
	config[0] = struct{}{} // Wikipedia
	config[2] = struct{}{} // 30s timer
	config[5] = struct{}{} // Capitalisation
	config[6] = struct{}{} // Punctuation
	config[7] = struct{}{} // Numbers

	return model{
		page:   InitMainMenu(),
		width:  w,
		height: h,
		styles: Styles{
			correct: func(str string) termenv.Style {
				return termenv.String(str).Foreground(foreground)
			},
			toEnter: func(str string) termenv.Style {
				return termenv.String(str).Foreground(foreground).Faint()
			},
			mistakes: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("1")).Underline()
			},
			cursor: func(str string) termenv.Style {
				return termenv.String(str).Reverse().Bold()
			},
			runningTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2"))
			},
			stoppedTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2")).Faint()
			},
			greener: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("6")).Faint()
			},
			faintGreen: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("10")).Faint()
			},
		},
		config: config,
	}
}

func (m model) Init() tea.Cmd {
	// CallClear()
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func InitMainMenu() MainMenu {
	return MainMenu{
		choices:  []string{"Start", "Settings"},
		selected: make(map[int]struct{}),
	}
}

func InitSettings(config map[int]struct{}) Settings {
	return Settings{
		choices:  []string{"Wikipedia", "Common words", "30s", "60s", "120s", "Capitalisation", "Punctuation", "Numbers", "Back"},
		selected: config,
	}
}

func formatText(text string, width int) []string {
	// Whitespace is removed by wordwrap to preserve with a pipe character
	text = strings.ReplaceAll(text, " ", "|")

	ww := wordwrap.NewWriter(width)
	ww.Breakpoints = []rune{'|'}
	ww.KeepNewlines = false
	ww.Write([]byte(text))
	ww.Close()

	text = strings.ReplaceAll(ww.String(), "|", " ")

	lines := strings.Split(text, "\n")

	return lines
}

func applyConfigFilters(text string, config map[int]struct{}) string {
	// Remove capitalisation
	if _, ok := config[5]; !ok {
		text = strings.ToLower(text)
	}
	// Remove punctuation
	if _, ok := config[6]; !ok {
		re := regexp.MustCompile("[!-/:-@[-`{-~.,?<>']")
		text = re.ReplaceAllString(text, "")
	}
	// Remove numbers
	if _, ok := config[7]; !ok {
		re := regexp.MustCompile(`[0-9]`)
		text = re.ReplaceAllString(text, "")
	}
	return text
}

func typingText(config map[int]struct{}) string {
	var text string
	if _, ok := config[0]; ok {
		text = WikiWords()
	} else if _, ok := config[1]; ok {
		text = CommonWords("/words/common_words.txt")
	}

	text = applyConfigFilters(text, config)
	return text
}

func timeLimit(config map[int]struct{}) (int, error) {
	if _, ok := config[2]; ok {
		return 30, nil
	} else if _, ok := config[3]; ok {
		return 60, nil
	} else if _, ok := config[4]; ok {
		return 120, nil
	}
	return 0, errors.New("error: no time limit config selected")
}

func MaxInt(x int, y int) int {
	if x >= y {
		return x
	}
	return y
}

func maxLen(arr []string) int {
	max := 0
	for _, str := range arr {
		max = MaxInt(max, len(str))
	}
	return max
}

func InitTyping(config map[int]struct{}) Typing {
	width := 50
	text := typingText(config)
	limit, err := timeLimit(config)
	if err != nil {
		panic(err)
	}
	lines := formatText(text, width)
	maxLineLen := maxLen(lines)
	return Typing{
		lines:         lines,
		correct:       NewCorrect(),
		width:         width,
		started:       false,
		cursorLine:    0,
		totalMistakes: 0,
		totalCorrect:  0,
		maxLineLen:    maxLineLen,
		time: &Time{
			lastUpdated: time.Now(),
			limit:       limit,
			remaining:   limit,
		},
		wps: make([]int, limit),
	}
}

func calcPerformance(accuracy float64, recovery float64, wpm float64, mistakes int) float64 {
	ideal := 100.0
	performance := (accuracy*recovery*wpm - float64(mistakes)*0.5) / ideal
	performance = math.Min(performance, 1.0)
	return performance
}

func InitResults(wpms []float64, wpm float64, accuracy float64, mistakes int, recovery float64) Results {
	performance := calcPerformance(accuracy, recovery, wpm, mistakes)
	return Results{
		wpms:        wpms,
		wpm:         wpm,
		accuracy:    accuracy,
		mistakes:    mistakes,
		recovery:    recovery,
		performance: performance,
	}
}
