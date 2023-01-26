package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/polpettone/labor/openai-client/cmd/config"
	"github.com/polpettone/labor/openai-client/pkg"
)

type OpenAIClient struct {
	apiKey string
}

func NewOpenAIClient() (*OpenAIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		return nil, errors.New("Missing API KEY: Provide an env var named OPENAI_API_KEY")
	}

	return &OpenAIClient{apiKey: apiKey}, nil
}

const COMPLETION_URL string = "https://api.openai.com/v1/completions"

func (o *OpenAIClient) Complete(
	question string,
	model string,
	temperature float64,
	maxTokens int) (*TextCompletion, error) {

	config.Logger.Debug().Msgf("Using model: %s \n", model)

	payload := Payload{
		Model:            model,
		Prompt:           question,
		Temperature:      temperature,
		MaxTokens:        maxTokens,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

	authHeader := fmt.Sprintf("Bearer %s", o.apiKey)

	requestBody, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", COMPLETION_URL, bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	startTime := time.Now()
	res, err := client.Do(req)
	endTime := time.Now()
	requestDuration := endTime.Sub(startTime)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var textCompletion TextCompletion
	err = json.Unmarshal([]byte(responseBody), &textCompletion)

	config.HistoryLogger.
		Info().
		Str("url", COMPLETION_URL).
		Int64("response_time_ms", requestDuration.Milliseconds()).
		Int("http_status_code", res.StatusCode).
		RawJSON("request_body", requestBody).
		RawJSON("response_body", responseBody).
		Send()

	if err != nil {
		return nil, err
	}

	return &textCompletion, nil
}

func (o *OpenAIClient) GenerateImage(imageDescription string, imageName string) error {

	payload := ImageCreatingPayload{
		Prompt: imageDescription,
		N:      1,
		Size:   "1024x1024",
	}

	url := "https://api.openai.com/v1/images/generations"
	authHeader := fmt.Sprintf("Bearer %s", o.apiKey)

	// Setze den Request-Body
	requestBody, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	var response ImageResponse
	err = json.Unmarshal([]byte(responseBody), &response)

	if err != nil {
		return err
	}

	err = pkg.DownloadFileFromUrl(response.Data[0].URL, imageName)

	if err != nil {
		return err
	}

	return nil

}

func (o *OpenAIClient) ListModels() ([]Model, error) {

	url := "https://api.openai.com/v1/models"
	authHeader := fmt.Sprintf("Bearer %s", o.apiKey)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)

	var modelList ModelList

	err = json.Unmarshal([]byte(responseBody), &modelList)

	if err != nil {
		return nil, err
	}

	return modelList.Data, nil
}
