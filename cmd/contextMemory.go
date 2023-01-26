package cmd

//todo repo to save state
type ContextMemory struct {
	buffer []string
	size   int
	tail   int
}

func NewContextMemory(size int) *ContextMemory {
	return &ContextMemory{
		size:   size,
		buffer: make([]string, size),
	}
}

func (c *ContextMemory) Add(value string) {
	c.buffer[c.tail] = value
	c.tail = (c.tail + 1) % c.size
}

func (c *ContextMemory) All() []string {
	return c.buffer
}

func (c *ContextMemory) Reset() {
	c.buffer = make([]string, c.size)
}
