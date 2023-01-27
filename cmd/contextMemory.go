package cmd

type Entry struct {
	value  string
	tokens int
}

//todo repo to save state
type ContextMemory struct {
	buffer    []*Entry
	maxTokens int
}

func NewContextMemory(maxTokens int) *ContextMemory {
	return &ContextMemory{
		buffer:    []*Entry{},
		maxTokens: maxTokens,
	}
}

func (c *ContextMemory) Add(value *Entry) {
	c.buffer = append(c.buffer, value)

	i := 0
	if c.TokenCount() > c.maxTokens {
		c.buffer = append(c.buffer[:i], c.buffer[i+1:]...)
	}

}

func (c *ContextMemory) All() string {
	all := ""
	for _, v := range c.buffer {
		if v != nil {
			all += v.value
			all += "\n"
		}
	}
	return all
}

func (c *ContextMemory) Reset() {
	c.buffer = []*Entry{}
}

func (c *ContextMemory) TokenCount() int {
	count := 0
	for _, k := range c.buffer {
		if k != nil {
			count += k.tokens
		}
	}
	return count
}
