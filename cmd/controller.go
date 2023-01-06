package cmd

import (
	"fmt"
	"strings"
	"time"
)

func Questioner(question string) (string, error) {

	fmt.Printf("Question: %s \n", question)
	fmt.Printf("Wait a moment...\n")

	client, err := NewOpenAIClient()

	if err != nil {
		return "", nil
	}

	response, err := client.Ask(question)

	if err != nil {
		return "", err
	}

	var result string
	for k, v := range response.Choices {
		result = fmt.Sprintf("%d: \n %s \n", k, v.Text)
	}

	return result, nil

}

func ImageGenerator(query string) (string, error) {

	fmt.Printf("generate image from: %s \n", query)
	fmt.Printf("Wait a moment...\n")

	client, err := NewOpenAIClient()

	if err != nil {
		return "", err
	}

	imageName := generateName(strings.Split(query, " ")[0])
	err = client.GenerateImage(query, imageName)

	if err != nil {
		return "", err
	}

	return imageName, nil
}

func generateName(value string) string {
	// Hol die aktuelle Zeit in Millisekunden
	milliseconds := time.Now().UnixNano() / int64(time.Millisecond)

	return fmt.Sprintf("%s-%d", value, milliseconds)
}
