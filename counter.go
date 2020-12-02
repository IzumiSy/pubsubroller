package main

import (
	"sync"
)

// mutex counter

type counter struct {
	Total int

	sync.Mutex
	done    int
	skipped int
}

func (c *counter) Done() {
	c.Lock()
	defer c.Unlock()
	c.done += 1
}

func (c *counter) Skipped() {
	c.Lock()
	defer c.Unlock()
	c.skipped += 1
}

func (c *counter) Result() (int, int) {
	return c.done, c.skipped
}
