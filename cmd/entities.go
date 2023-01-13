package cmd


type ModelList struct {
	Data    []Model  `json:"data"`
	Object string   `json:"object"`
} 

type Model struct {
	ID           string           `json:"id"`
	Object       string           `json:"object"`
	Created      int              `json:"created"`
	OwnedBy      string           `json:"owned_by"`
	Permission   []ModelPermission `json:"permission"`
	Root         string           `json:"root"`
	Parent       string           `json:"parent"`
}

type ModelPermission struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int    `json:"created"`
	AllowCreateEngine bool   `json:"allow_create_engine"`
	AllowSampling     bool   `json:"allow_sampling"`
	AllowLogprobs     bool   `json:"allow_logprobs"`
	AllowSearchIndices bool  `json:"allow_search_indices"`
	AllowView         bool   `json:"allow_view"`
	AllowFineTuning   bool   `json:"allow_fine_tuning"`
	Organization      string `json:"organization"`
	Group             string `json:"group"`
	IsBlocking        bool   `json:"is_blocking"`
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

