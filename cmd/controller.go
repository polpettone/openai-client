package cmd

import (
	"fmt"
	"log"
)

func Questioner(question string) (string, error) {

	fmt.Printf("Question: %s \n", question)
	fmt.Printf("Wait a moment...\n")

	client, err := NewOpenAIClient()

	if err != nil {
		log.Fatal(err)
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

func TextCompletionNew() {

	err := callCompletion("The first thing you should know about javascript is")

	if err != nil {
		log.Fatal(err)
	}

}
