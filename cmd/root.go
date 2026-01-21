package cmd

import (
	"aicommits/internal/config"
	"aicommits/internal/git"
	"aicommits/internal/llm"
	"aicommits/internal/ui" // 引入 UI 包
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var shouldStageAll bool
var rootCmd = &cobra.Command{
	Use:   "aicommits",
	Short: "使用AI编写Git提交日志",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 加载配置
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("❌ 配置加载失败: %v\n", err)
			return
		}

		// 检查必要参数
		if cfg.APIKey == "" {
			fmt.Println("❌ 未检测到 API Key。")
			fmt.Println("请先运行配置命令:")
			fmt.Println("  aicommits config")
			return
		}

		if shouldStageAll {
			if err := git.StageAll(); err != nil {
				fmt.Printf("❌ 无法将变更加入暂存区: %v\n", err)
				return
			}
		}

		// 1. 获取 Diff
		diff, err := git.GetStagedDiff()
		if err != nil {
			fmt.Printf("❌ Git错误: %v\n", err)
			return
		}
		if diff == "" {
			fmt.Println("⚠️ 暂存区为空，请先执行 git add")
			return
		}

		// 2. 初始化 LLM Client
		// 这里为了演示方便，配置写死，之后可以用 Viper 做配置文件
		client := llm.NewProvider(llm.ProviderConfig{
			BaseURL:               cfg.BaseURL,
			APIKey:                cfg.APIKey,
			Model:                 cfg.Model,
			Language:              cfg.Language,
			WithDescription:       cfg.WithDescription,
			SubjectSeparateSymbol: cfg.SubjectSeparateSymbol,
		}) // 3. 启动 UI 程序
		// 创建一个带有超时的 Context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		model := ui.NewModel(ctx, client, diff)
		p := tea.NewProgram(model)

		// 运行 UI，它会阻塞直到用户按 Enter/Esc/Ctrl+C
		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("UI 错误: %v\n", err)
			return
		}

		// 4. 处理最终结果
		// 类型断言取回我们的 Model 数据
		m, ok := finalModel.(ui.Model)
		if !ok {
			return
		}

		// 如果用户确认了提交
		if m.Confirmed && m.Msg != "" {
			git.Commit(m.Msg)
		} else {
			fmt.Println("\n🚫 已取消提交")
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}
