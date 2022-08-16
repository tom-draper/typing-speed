package cmd

import (
	"errors"
	"sync"
)

type Correct struct {
	lock  sync.Mutex
	stack []bool
}

func NewCorrect() *Correct {
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

func (c *Correct) AtIndex(index int) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if index > len(c.stack) {
		panic(errors.New("index value exceeds stack length"))
	}
	return c.stack[index]
}

func (c *Correct) Length() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.stack)
}

func (c *Correct) FinalAccuracy() float32 {
	c.lock.Lock()
	defer c.lock.Unlock()
	correct := 0
	for i := 0; i < len(c.stack); i++ {
		if c.stack[i] {
			correct++
		}
	}
	return float32(correct) / float32(len(c.stack)) * 100
}

func (c *Correct) Mistakes() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	mistakes := 0
	for i := 0; i < len(c.stack); i++ {
		if !c.stack[i] {
			mistakes++
		}
	}
	return mistakes
}
