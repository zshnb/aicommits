package llm

import "fmt"

// PromptOptions 定义构建提示词所需的参数
type PromptOptions struct {
	Language string // "cn" 或 "en"
	Diff     string // Git diff 内容
}

const (
	systemPromptTpl = `You are an expert developer and git specialist.
Your task is to generate a concise and standardized git commit message based on the provided code changes (diff).

Format Requirements:
1. Use the Conventional Commits format: <type>: <subject>
2. The subject line must be less than 50 characters.
3. Leave a blank line after the subject.
4. Provide a detailed description body (wrapping at 72 chars) if the changes are complex.
5. Do NOT include markdown blocks (like ''' or code fences). Just return the raw message.
6. %s`

	langInstructionCN = "The commit message MUST be written in Simplified Chinese (简体中文)."
	langInstructionEN = "The commit message MUST be written in English."
)

// ConstructMessages 构建发送给 LLM 的消息列表
func ConstructMessages(opts PromptOptions) []Message {
	// 1. 确定语言指令
	langInstruction := langInstructionEN
	if opts.Language == "cn" {
		langInstruction = langInstructionCN
	}

	// 2. 组装 System Prompt
	finalSystemPrompt := fmt.Sprintf(systemPromptTpl, langInstruction)

	// 3. 返回消息结构
	return []Message{
		{Role: "system", Content: finalSystemPrompt},
		{Role: "user", Content: fmt.Sprintf("Here is the git diff output:\n\n%s", opts.Diff)},
	}
}
