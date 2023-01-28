package cmd

type Entry struct {
	Value  string `json:"value"`
	Tokens int    `json:"tokens"`
}

type ContextMemory struct {
	Buffer    []*Entry `json:"buffer"`
	MaxTokens int      `json:"max_tokens"`
	ID        string   `json:"id"`
}

func NewContextMemory(id string, maxTokens int) *ContextMemory {
	return &ContextMemory{
		Buffer:    []*Entry{},
		MaxTokens: maxTokens,
		ID:        id,
	}
}

func (c *ContextMemory) Add(value *Entry) {
	c.Buffer = append(c.Buffer, value)

	i := 0
	for c.TokenCount() > c.MaxTokens {
		c.Buffer = append(c.Buffer[:i], c.Buffer[i+1:]...)
	}
}

func (c *ContextMemory) All() string {
	all := ""
	for _, v := range c.Buffer {
		if v != nil {
			all += v.Value
			all += "\n"
		}
	}
	return all
}

func (c *ContextMemory) Reset() {
	c.Buffer = []*Entry{}
}

func (c *ContextMemory) TokenCount() int {
	count := 0
	for _, k := range c.Buffer {
		if k != nil {
			count += k.Tokens
		}
	}
	return count
}
