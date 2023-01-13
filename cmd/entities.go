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
