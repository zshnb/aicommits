package llm

import "context"

// Client 定义了所有 LLM 提供商必须实现的通用接口
// 无论是 OpenAI, DeepSeek 还是 Ollama，都必须满足这个契约
type Client interface {
	GenerateCommitMessage(ctx context.Context, diff string) (string, error)
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}
