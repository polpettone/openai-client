package cmd

import (
	"fmt"
	"time"

	"github.com/polpettone/labor/openai-client/cmd/config"
)

type Provider struct {
	contextMemory  *ContextMemory
	contextEnabled bool
}

func NewProvider(maxTokens int, contextEnabled bool, contextMemoryID string) *Provider {

	return &Provider{
		contextMemory:  NewContextMemory(contextMemoryID, maxTokens),
		contextEnabled: contextEnabled,
	}

}

func (p *Provider) ClearContext() {
	p.contextMemory.Reset()
}

func (p *Provider) Prompt(
	text string,
	model string,
	temperature float64,
	maxTokens int) (string, error) {
	client, err := NewOpenAIClient()

	if err != nil {
		return "", nil
	}

	prompt := text

	if p.contextEnabled {
		prompt = fmt.Sprintf("%s \n %s", p.contextMemory.All(), text)
	}

	response, err := client.Complete(
		prompt,
		model,
		temperature,
		maxTokens)

	if err != nil {
		return "", err
	}

	promptTokens := response.Usage.PromptTokens
	completionTokens := response.Usage.CompletionTokens

	if p.contextEnabled {
		entry := &Entry{value: text, tokens: promptTokens}
		p.contextMemory.Add(entry)
	}

	var result string
	for _, v := range response.Choices {
		result = fmt.Sprintf("%s\n", v.Text)
		if p.contextEnabled {
			entry := &Entry{value: result, tokens: completionTokens}
			p.contextMemory.Add(entry)
		}
	}

	config.ContextMemoryLogger.
		Info().
		Str("memory", p.contextMemory.All()).
		Int("token_count", p.contextMemory.TokenCount()).
		Int("max_tokens", p.contextMemory.maxTokens).
		Send()

	return result, nil

}

func ImageGenerator(query string, imageName string) (string, error) {
	client, err := NewOpenAIClient()
	if err != nil {
		return "", err
	}

	err = client.GenerateImage(query, imageName)

	if err != nil {
		return "", err
	}

	return imageName, nil
}

func ListModels() ([]Model, error) {

	client, err := NewOpenAIClient()
	if err != nil {
		return nil, err
	}
	models, err := client.ListModels()
	if err != nil {
		return nil, err
	}
	return models, nil
}

func generateName(value string) string {
	time.Sleep(10 * time.Millisecond)
	milliseconds := time.Now().UnixNano() / int64(time.Millisecond)

	return fmt.Sprintf("%s-%d", value, milliseconds)
}
