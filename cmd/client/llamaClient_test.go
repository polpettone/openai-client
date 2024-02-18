package client

import (
	"encoding/json"
	"fmt"
	"testing"
)

var responseBody = `
	{
	"choices": [
      {
        "finish_reason": "stop",
        "index": 0,
        "message": {
          "content": "Hallo! Wie kann ich Ihnen heute helfen?",
          "role": "assistant"
        }
      }
    ],
    "created": 1708259981,
    "id": "chatcmpl-TSCJHUPtHJ1GM3qr5kbnn4jqRVirAvdW",
    "model": "LLaMA_CPP",
    "object": "chat.completion",
    "usage": {
      "completion_tokens": 14,
      "prompt_tokens": 30,
      "total_tokens": 44
    }
	}
`

func Test_unmarshalling(t *testing.T) {
	var llamaResponse LlamaResponse
	err := json.Unmarshal([]byte(responseBody), &llamaResponse)
	if err != nil {
		t.Errorf("Wanted no err, got |%v|", err)
	}

	fmt.Printf("%v", llamaResponse)
}
