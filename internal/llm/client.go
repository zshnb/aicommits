package llm

import "context"

// Client 定义了所有 LLM 提供商必须实现的通用接口
// 无论是 OpenAI, DeepSeek 还是 Ollama，都必须满足这个契约
type Client interface {
	// GenerateCommitMessage 根据代码差异生成提交信息
	GenerateCommitMessage(ctx context.Context, diff string) (string, error)
}
