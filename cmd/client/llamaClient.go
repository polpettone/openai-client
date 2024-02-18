package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/polpettone/openai-client/cmd/logging"
)

type LlamaClient struct {
}

func NewLlamaClient() *LlamaClient {
	return &LlamaClient{}
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type LlamaRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type LlamaResponse struct {
	Choices []struct {
		FinishReason string  `json:"finish_reason"`
		Index        int     `json:"index"`
		Message      Message `json:"message"`
	} `json:"choices"`
	ID     string `json:"id"`
	Model  string `json:"model"`
	Object string `json:"object"`

	Usage struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

const llamaURL string = "http://localhost:8080/v1/chat/completions"

func (o *LlamaClient) Complete(question string) (*LlamaResponse, error) {

	request := LlamaRequest{
		Messages: []Message{{Content: question, Role: "user"}},
		Model:    "LLaMA_CPP",
	}
	authHeader := fmt.Sprintf("Bearer %s", "no-key")
	requestBody, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", llamaURL, bytes.NewBuffer(requestBody))

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

	if res.StatusCode != 200 {
		return nil, errors.New(string(responseBody))
	}

	var llamaResponse LlamaResponse
	err = json.Unmarshal(responseBody, &llamaResponse)

	if err != nil {
		logging.Logger.Info().Str("err", err.Error()).Send()
	}

	logging.HistoryLogger.
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

	return &llamaResponse, nil
}
