package cmd

type OPENAI_ChatRequest struct {
	Prompt           string   `json:"prompt"`
	MaxTokens        int      `json:"max_tokens"`
	Model            string   `json:"model"`
	Temperature      float64  `json:"temperature"`
	TopP             float64  `json:"top_p"`
	FrequencyPenalty float64  `json:"frequency_penalty"`
	PresencePenalty  float64  `json:"presence_penalty"`
	Stop             []string `json:"stop"`
	User             string   `json:"user"`
}

type OPENAI_ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type OPENAI_ImageRequest struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
	User           string `json:"user"`
}

type OPENAI_ImageEditRequest struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
	User           string `json:"user"`
	Image          string `json:"image"`
	Mask           string `json:"mask"`
}

type OPENAI_ImageResponse struct {
	Created int64 `json:"created"`
	Data    []struct {
		URL string `json:"url"`
	} `json:"data"`
}

type OPENAI_ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

type OPENAI_ChatCompletionResponse struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Created int             `json:"created"`
	Choices []OPENAI_Choice `json:"choices"`
	Usage   map[string]int  `json:"usage"`
}

type OPENAI_Choice struct {
	Index        int            `json:"index"`
	Message      OPENAI_Message `json:"message"`
	FinishReason string         `json:"finish_reason"`
}

type OPENAI_Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OPENAI_ChatCompletionRequest struct {
	Model            string           `json:"model"`
	Messages         []OPENAI_Message `json:"messages"`
	MaxTokens        int              `json:"max_tokens"`
	Temperature      float64          `json:"temperature"`
	TopP             float64          `json:"top_p"`
	N                int              `json:"n"`
	Stream           bool             `json:"stream"`
	PresencePenalty  float64          `json:"presence_penalty"`
	FrequencyPenalty float64          `json:"frequency_penalty"`
	Stop             []string         `json:"stop"`
	User             string           `json:"user"`
}
