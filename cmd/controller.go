package cmd

import (
	"fmt"
	"time"
)

func Questioner(
	question string,
	model string,
	temperature float64,
	maxTokens int) (string, error) {
	client, err := NewOpenAIClient()

	if err != nil {
		return "", nil
	}

	response, err := client.Ask(question, model, temperature, maxTokens)

	if err != nil {
		return "", err
	}

	var result string
	for _, v := range response.Choices {
		result = fmt.Sprintf("%s\n", v.Text)
	}

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
