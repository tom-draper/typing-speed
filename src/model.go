package src

import "time"

type MainMenu struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

type Settings struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

type Time struct {
	lastUpdated time.Time
	remaining   int
}

type Typing struct {
	words   string
	started bool
	cursor  int // which to-do list item our cursor is pointing at
	time    *Time
}

type Page interface {
	view() string
}

type model struct {
	page Page
}
