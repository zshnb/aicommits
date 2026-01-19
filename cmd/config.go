package cmd

import (
	"aicommits/internal/config"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置工具 (交互式)",
	Run: func(cmd *cobra.Command, args []string) {
		interactiveConfig()
	},
}

// interactiveConfig 启动交互式表单
func interactiveConfig() {
	// 1. 读取现有配置作为默认值
	currentCfg, _ := config.Load()
	if currentCfg == nil {
		currentCfg = &config.Config{}
	}

	// 定义表单绑定的变量
	var (
		provider = currentCfg.Provider
		apiKey   = currentCfg.APIKey
		model    = currentCfg.Model
		baseURL  = currentCfg.BaseURL
		language = currentCfg.Language
	)

	// 如果是第一次配置，设置一些默认值
	if language == "" {
		language = "cn"
	}
	if provider == "" {
		provider = "deepseek"
	}

	// 2. 第一步：选择提供商
	// 我们将表单分成两组，因为更改提供商会影响后续的默认值
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("选择 AI 提供商").
				Options(
					huh.NewOption("DeepSeek", "deepseek"),
					huh.NewOption("OpenAI", "openai"),
					huh.NewOption("Ollama (本地)", "ollama"),
					huh.NewOption("自定义", "custom"),
				).
				Value(&provider),
		),
	).Run()

	if err != nil {
		fmt.Println("❌ 配置已取消")
		return
	}

	// 根据提供商自动填充默认值 (如果用户原本没填过)
	switch provider {
	case "deepseek":
		if baseURL == "" {
			baseURL = "https://api.deepseek.com"
		}
		if model == "" {
			model = "deepseek-chat"
		}
	case "openai":
		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}
		if model == "" {
			model = "gpt-5-nano"
		}
	case "ollama":
		if baseURL == "" {
			baseURL = "http://localhost:11434/v1"
		}
		if model == "" {
			model = "llama3"
		}
	}

	// 3. 第二步：填写详细信息
	// 根据是否需要 API Key 动态调整表单
	formFields := []huh.Field{
		huh.NewInput().
			Title("API Endpoint (Base URL)").
			Value(&baseURL),

		huh.NewInput().
			Title("模型名称 (Model)").
			Value(&model),

		huh.NewSelect[string]().
			Title("提交日志语言").
			Options(
				huh.NewOption("中文 (Chinese)", "cn"),
				huh.NewOption("English", "en"),
			).
			Value(&language),
	}

	// Ollama 通常不需要 API Key，其他需要
	if provider != "ollama" {
		// 在 BaseURL 之前插入 API Key 输入框
		apiKeyField := huh.NewInput().
			Title("API Key").
			Value(&apiKey).
			EchoMode(huh.EchoModePassword). // 密码掩码模式
			Description("DeepSeek 或 OpenAI 的密钥")

		// 插入到切片头部
		formFields = append([]huh.Field{apiKeyField}, formFields...)
	}

	err = huh.NewForm(
		huh.NewGroup(formFields...),
	).Run()

	if err != nil {
		fmt.Println("❌ 配置已取消")
		return
	}

	// 4. 保存配置
	newConfig := &config.Config{
		Provider: provider,
		APIKey:   apiKey,
		BaseURL:  baseURL,
		Model:    model,
		Language: language,
	}

	if err := config.Save(newConfig); err != nil {
		fmt.Printf("❌ 保存失败: %v\n", err)
		return
	}

	fmt.Println("✅ 配置已成功保存到 ~/.aicommits.yaml")
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "设置配置项",
	Args:  cobra.ExactArgs(2), // 必须传入 key 和 value
	Run: func(cmd *cobra.Command, args []string) {
		key := strings.ToLower(args[0])
		val := args[1]

		// 简单的校验
		validKeys := map[string]bool{
			"api_key":  true,
			"model":    true,
			"base_url": true,
		}

		if !validKeys[key] {
			fmt.Printf("❌ 无效的配置项: %s\n仅支持: api_key, model, base_url\n", key)
			return
		}

		if err := config.Set(key, val); err != nil {
			fmt.Printf("❌ 保存配置失败: %v\n", err)
			return
		}

		fmt.Printf("✅ 已更新 %s\n", key)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "查看当前配置",
	Run: func(cmd *cobra.Command, args []string) {
		// 触发加载
		config.Load()
		fmt.Println(config.GetPrintable())
	},
}

func init() {
	// 将子命令注册到 config 命令
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(listCmd)

	// 将 config 命令注册到根命令
	rootCmd.AddCommand(configCmd)
}
