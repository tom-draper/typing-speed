package src

import (
	"errors"
	"sync"
)

type Correct struct {
	lock  sync.Mutex // you don't have to do this if you don't want thread safety
	stack []bool
}

func NewStack() *Correct {
	return &Correct{sync.Mutex{}, make([]bool, 0)}
}

func (c *Correct) Push(v bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.stack = append(c.stack, v)
}

func (c *Correct) Pop() (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	l := len(c.stack)
	if l == 0 {
		return false, errors.New("empty correct stack")
	}

	res := c.stack[l-1]
	c.stack = c.stack[:l-1]
	return res, nil
}

func (c *Correct) AtIndex(index int) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if index > len(c.stack) {
		return false, errors.New("index value exceeds stack length")
	}
	return c.stack[index], nil
}

func (c *Correct) Length() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.stack)
}
