package cmd

import (
	"time"

	"github.com/muesli/termenv"
)

type MainMenu struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

type Settings struct {
	choices  []string
	cursor   int
	selected map[int]struct{} // Reference to model.config
}

type Time struct {
	lastUpdated time.Time
	limit       int
	remaining   int
}

type Typing struct {
	lines         []string
	correct       *Correct
	width         int
	started       bool
	totalMistakes int
	totalCorrect  int
	cursor        int
	cursorLine    int
	maxLineLen    int
	time          *Time
	wps           []int
	words         int
}

type Results struct {
	wpm         float64
	wpms        []float64
	accuracy    float64
	mistakes    int
	recovery    float64
	performance float64
}

type Page interface {
	view(style Styles, width int, height int) string
}

type Style func(string) termenv.Style

type Styles struct {
	correct   Style
	mistakes  Style
	toEnter   Style
	cursor    Style
	highlight Style
	title     Style
}

type model struct {
	page   Page
	styles Styles
	width  int
	height int
	config map[int]struct{}
}
