package llm

import "fmt"

// PromptOptions 定义构建提示词所需的参数
type PromptOptions struct {
	Language        string // "cn" 或 "en"
	Diff            string // Git diff 内容
	WithDescription bool
	WithAppendix    bool
}

const (
	systemPromptTpl = `
<role>
You are an expert developer and git specialist.
</role>
<goal>
Your task is to generate a concise and standardized git commit message based on the provided code changes (diff).
</goal>
<context>
please follow below type definition
- build: Used for modifying the project build system, such as changing dependencies, external interfaces, or upgrading Node versions.
- chore: Used for modifying non-business code, such as changing build processes or tool configurations.
- ci: Used for modifying the Continuous Integration process, such as changing Travis, Jenkins, or other workflow configurations.
- docs: Used for modifying documentation, such as changing README files or API documentation.
- style: Used for modifying code style, such as adjusting indentation, spaces, blank lines, etc.
- refactor: Used for code refactoring, such as modifying code structure, variable names, or function names without changing functional logic.
- perf: Used for performance optimization, such as improving code performance or reducing memory usage.
</context>
<restriction>
- Use the Conventional Commits format: <type>[optional scope]: <subject>
- The subject line **MUST** be less than 100 characters.
- Do NOT include markdown blocks (like ''' or code fences). Just return the raw message.
%s
</restriction>
`

	withDescriptionPrompt = "- Provide a detailed description body around 3 - 5 lines, each line **MUST** be less than 72 char. Leave a blank line after the subject."
	langInstructionCN     = "- The commit message **MUST** be written in Simplified Chinese (简体中文)."
	langInstructionEN     = "- The commit message **MUST** be written in English."
)

func ConstructMessages(opts PromptOptions) []Message {
	moreInstruction := langInstructionEN
	if opts.Language == "cn" {
		moreInstruction = langInstructionCN
	}

	if opts.WithDescription {
		moreInstruction += "\n" + withDescriptionPrompt
	}

	// 2. 组装 System Prompt
	finalSystemPrompt := fmt.Sprintf(systemPromptTpl, moreInstruction)

	// 3. 返回消息结构
	return []Message{
		{Role: "system", Content: finalSystemPrompt},
		{Role: "user", Content: fmt.Sprintf("Here is the git diff output:\n\n%s", opts.Diff)},
	}
}
