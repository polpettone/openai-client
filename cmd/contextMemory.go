package cmd

type Entry struct {
	value  string
	tokens int
}

//todo repo to save state
type ContextMemory struct {
	buffer []*Entry
	size   int
	tail   int
}

func NewContextMemory(size int) *ContextMemory {
	return &ContextMemory{
		size:   size,
		buffer: make([]*Entry, size),
	}
}

func (c *ContextMemory) Add(value *Entry) {
	c.buffer[c.tail] = value
	c.tail = (c.tail + 1) % c.size
}

func (c *ContextMemory) All() string {
	all := ""
	for n := 0; n < c.size; n++ {
		if c.buffer[n] != nil {
			all += c.buffer[n].value
			all += "\n"
		}
	}
	return all
}

func (c *ContextMemory) Reset() {
	c.buffer = make([]*Entry, c.size)
}
