package cmd

import (
	"aicommits/internal/git"
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aicommits",
	Short: "使用AI编写Git提交日志",
	Long:  `一个基于大语言模型的CLI工具，读取暂存区的Diff并自动生成符合规范的Commit Message。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 正在分析代码变更...")

		// 1. 获取Git Diff
		diff, err := git.GetStagedDiff()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		if diff == "" {
			fmt.Println("⚠️ 暂存区(Staged)为空，请先执行 git add")
			return
		}

		// 暂时先打印Diff长度，证明读取成功
		fmt.Printf("✅ 成功获取Diff，长度为: %d 字符\n", len(diff))
		fmt.Println("🔜 下一步：将Diff发送给LLM生成日志...")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
