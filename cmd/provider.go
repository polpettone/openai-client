package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/polpettone/openai-client/cmd/config"
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

func saveContextMemory(context *ContextMemory) error {
	data, err := json.Marshal(context)
	if err != nil {
		return err
	}

	err = os.WriteFile(context.ID, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadContextMemory(contextID string) (*ContextMemory, error) {

	content, err := os.ReadFile(contextID)
	if err != nil {
		return nil, err
	}

	contextMemory := &ContextMemory{}
	err = json.Unmarshal(content, contextMemory)
	if err != nil {
		return nil, err
	}

	return contextMemory, nil
}

func (p *Provider) ClearContext() error {
	p.contextMemory.Reset()
	err := saveContextMemory(p.contextMemory)
	if err != nil {
		return err
	}
	return nil
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
		if _, err := os.Stat(p.contextMemory.ID); err == nil {
			p.contextMemory, err = loadContextMemory(p.contextMemory.ID)
			if err != nil {
				return "", err
			}
			prompt = fmt.Sprintf("%s \n %s", p.contextMemory.All(), text)
		} else {
			config.FileLogger.
				Info().
				Str("msg",
					fmt.Sprintf("no context memory found with id: %s",
						p.contextMemory.ID))
		}
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
		entry := &Entry{Value: text, Tokens: promptTokens}
		p.contextMemory.Add(entry)
	}

	var result string
	for _, v := range response.Choices {
		result = fmt.Sprintf("%s\n", v.Text)
		if p.contextEnabled {
			entry := &Entry{Value: result, Tokens: completionTokens}
			p.contextMemory.Add(entry)
		}
	}

	if p.contextEnabled {
		err = saveContextMemory(p.contextMemory)
	}

	if err != nil {
		return "", err
	}

	config.ContextMemoryLogger.
		Info().
		Str("memory", p.contextMemory.All()).
		Int("token_count", p.contextMemory.TokenCount()).
		Int("max_tokens", p.contextMemory.MaxTokens).
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
