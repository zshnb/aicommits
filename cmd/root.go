package cmd

import (
	"aicommits/internal/git"
	"aicommits/internal/llm" // 引入新包
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aicommits",
	Short: "使用AI编写Git提交日志",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 检查环境变量 (为了MVP快速验证)
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Println("❌ 错误: 未设置 OPENAI_API_KEY 环境变量")
			fmt.Println("提示: export OPENAI_API_KEY='sk-...'")
			return
		}

		fmt.Println("🚀 正在分析代码变更...")

		// 2. 获取 Diff
		diff, err := git.GetStagedDiff()
		if err != nil {
			fmt.Printf("❌ Git错误: %v\n", err)
			return
		}
		if diff == "" {
			fmt.Println("⚠️ 暂存区为空，请先执行 git add")
			return
		}

		// 3. 初始化 LLM 客户端
		// 这里演示如何配置为 DeepSeek (只需要改 BaseURL 和 Model)
		// 如果你想用官方 OpenAI，就把 BaseURL 留空，Model 改为 "gpt-3.5-turbo"
		client := llm.NewOpenAIClient(llm.OpenAIConfig{
			APIKey: apiKey,
			Model:  "gpt-5-nano", // 示例：DeepSeek 模型
		})

		fmt.Println("⏳ 正在请求 AI 生成日志...")

		// 4. 调用接口
		msg, err := client.GenerateCommitMessage(context.Background(), diff)
		if err != nil {
			fmt.Printf("❌ AI生成失败: %v\n", err)
			return
		}

		// 5. 输出结果
		fmt.Println("\n------------------------------------------------")
		fmt.Println(msg)
		fmt.Println("------------------------------------------------")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
