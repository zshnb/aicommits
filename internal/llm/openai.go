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

// OpenAIConfig 配置结构体
type OpenAIConfig struct {
	APIKey  string
	BaseURL string
	Model   string
	Timeout time.Duration
}

type openAIClient struct {
	config OpenAIConfig
	client *http.Client
}

// NewOpenAIClient 创建一个新的客户端实例
func NewOpenAIClient(cfg OpenAIConfig) Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.openai.com/v1"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	return &openAIClient{
		config: cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

// 定义请求和响应结构体 (只定义我们需要的部分)
type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *openAIClient) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	// 1. 准备 System Prompt
	systemPrompt := `你是一个资深的Git专家。请根据用户的git diff内容，生成一个简洁、符合Conventional Commits规范的提交消息。
	格式要求：
	1. 第一行是标题（<type>: <subject>），不超过50个字符。
	2. 空一行。
	3. 第三行开始是详细描述（可选），每行不超过72个字符。
	4. 只返回提交消息本身，不要包含任何markdown代码块（如 '''）或解释性文字。`

	// 2. 构建请求体
	reqBody := chatRequest{
		Model: c.config.Model,
		Messages: []message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: fmt.Sprintf("这是我的代码变更:\n\n%s", diff)},
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// 3. 创建 HTTP 请求
	url := strings.TrimRight(c.config.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	// 4. 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API调用失败: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 5. 解析响应
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API错误 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var result chatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("模型未返回任何内容")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}
