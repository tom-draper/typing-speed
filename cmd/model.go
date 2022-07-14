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
	selected map[int]struct{}
}

type Time struct {
	lastUpdated time.Time
	limit       int
	remaining   int
}

type Typing struct {
	chars    []string
	correct  *Correct
	width    int
	started  bool
	mistakes int
	cursor   int
	time     *Time
}

type Results struct {
	wpm      float32
	accuracy float32
	mistakes int
}

type Page interface {
	view(style Styles, width int, height int) string
}

type Style func(string) termenv.Style

type Styles struct {
	correct      Style
	toEnter      Style
	mistakes     Style
	cursor       Style
	runningTimer Style
	stoppedTimer Style
	greener      Style
	faintGreen   Style
}

type model struct {
	page   Page
	styles Styles
	width  int
	height int
}