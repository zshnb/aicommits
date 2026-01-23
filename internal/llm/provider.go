package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ProviderConfig 定义初始化 Provider 所需的配置
// 这些字段直接对应 config 包中的内容
type ProviderConfig struct {
	BaseURL               string
	Path                  string
	APIKey                string
	Model                 string
	Language              string
	Timeout               time.Duration
	WithDescription       bool
	SubjectSeparateSymbol string
}

// genericProvider 是通用的 OpenAI 兼容协议实现
type genericProvider struct {
	cfg    ProviderConfig
	client *http.Client
}

// NewProvider 创建通用实例
func NewProvider(cfg ProviderConfig) Client {
	// 确保 BaseURL 格式正确 (移除末尾斜杠，并确保包含 /v1 路径，如果厂商API不需要v1需自行调整逻辑或配置)
	// 大部分兼容接口（DeepSeek, OpenAI, Ollama）通常以 /v1 结尾
	// 为了鲁棒性，我们简单处理：如果 URL 没包含 chat/completions，我们在请求时拼接

	if cfg.Timeout == 0 {
		cfg.Timeout = 60 * time.Second // DeepSeek 有时响应较慢，给大一点超时
	}

	return &genericProvider{
		cfg: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (p *genericProvider) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	// 1. 利用 prompt.go 构建消息
	messages := ConstructMessages(PromptOptions{
		Language:              p.cfg.Language,
		Diff:                  diff,
		WithDescription:       p.cfg.WithDescription,
		SubjectSeparateSymbol: p.cfg.SubjectSeparateSymbol,
	})

	// 2. 构建请求 Payload
	reqBody := ChatRequest{
		Model:    p.cfg.Model,
		Messages: messages,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request failed: %w", err)
	}

	url := p.cfg.BaseURL + p.cfg.Path
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if p.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.cfg.APIKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 5. 处理响应
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("unmarshal response failed: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty response from model")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}
