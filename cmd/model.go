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
	lines              []string
	correct            *Correct
	width              int
	started            bool
	keystrokesCorrect  int
	keystrokesMistakes int
	cursor             int
	cursorLine         int
	maxLineLen         int
	time               *Time
	wps                []int
	words              int
}

type Results struct {
	wpm                float64
	wpms               []float64
	accuracy           float64
	keystrokesMistakes int
	keystrokesCorrect  int
	recovery           float64
	performance        float64
}

type Page interface {
	view(style Styles, width int, height int) string
}

type Style func(string) termenv.Style

type Styles struct {
	normal    Style
	correct   Style
	mistakes  Style
	err       Style
	cursor    Style
	highlight Style
	title     Style
}

type Config struct {
	config    map[int]struct{}
	wikiLinks []string // Temp wiki links
}

type model struct {
	page   Page
	styles Styles
	width  int
	height int
	config Config
}
