package main

type counter struct {
	Total int

	done    int
	skipped int
}

func (c *counter) Done() {
	c.done += 1
}

func (c counter) Skipped() {
	c.skipped += 1
}

func (c counter) Result() (int, int) {
	return c.done, c.skipped
}
