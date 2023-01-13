package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/polpettone/labor/openai-client/cmd/config"
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

type ImageCreatingPayload struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ImageResponse struct {
	Created int                `json:"created"`
	Data    []ImageResponseURL `json:"data"`
}

type ImageResponseURL struct {
	URL string `json:"url"`
}

type TextCompletion struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Text         string      `json:"text"`
	Index        int         `json:"index"`
	LogProbs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Payload struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	Temperature      float64 `json:"temperature"`
	MaxTokens        int     `json:"max_tokens"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

var models = []string{
	"code-cushman-001",
	"code-davinci-002",
	"text-davinci-003",
}

const COMPLETION_URL string = "https://api.openai.com/v1/completions"

func (o *OpenAIClient) Ask(
	question string,
	model string,
	temperature float64) (*TextCompletion, error) {

	config.Logger.Debug().Msgf("Using model: %s \n", model)

	payload := Payload{
		Model:            model,
		Prompt:           question,
		Temperature:      temperature,
		MaxTokens:        1000,
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

	err = downloadImageFromUrl(response.Data[0].URL, imageName)

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
func downloadImageFromUrl(url string, imageName string) error {

	fmt.Println("Download Image ...")
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Fehler beim Herunterladen des Bildes:", res.Status)
		return nil
	}

	file, err := os.Create(imageName)
	if err != nil {
		fmt.Println("Fehler beim Erstellen der Datei:", err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		fmt.Println("Fehler beim Schreiben des Bildes:", err)
		return err
	}

	fmt.Println("Bild erfolgreich heruntergeladen.")
	return nil
}
