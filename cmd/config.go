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
	Short: "Config",
	Run: func(cmd *cobra.Command, args []string) {
		interactiveConfig()
	},
}

var providerModels = map[string][]string{
	"deepseek": {"deepseek-chat", "deepseek-reasoner"},
	"openai":   {"gpt-5-nano", "gpt-5-mini", "gpt-5.1", "gpt-4o"},
	"grok":     {"grok-4-1-fast-non-reasoning", "grok-4-1-fast-reasoning", "grok-code-fast-1"}, // 常用本地模型
	"claude":   {"claude-sonnet-4-5", "claude-haiku-4-5", "claude-opus-4-5"},                   // 常用本地模型
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
		provider              = currentCfg.Provider
		apiKey                = currentCfg.APIKey
		model                 = currentCfg.Model
		baseURL               = currentCfg.BaseURL
		path                  = currentCfg.Path
		language              = currentCfg.Language
		withDescription       = currentCfg.WithDescription
		subjectSeparateSymbol = currentCfg.SubjectSeparateSymbol
	)

	// 如果是第一次配置，设置一些默认值
	if language == "" {
		language = "cn"
	}
	if provider == "" {
		provider = "deepseek"
	}
	if subjectSeparateSymbol == "" {
		subjectSeparateSymbol = ","
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
					huh.NewOption("Grok", "grok"),
					huh.NewOption("Claude", "claude"),
				).
				Value(&provider),
		),
	).Run()

	if err != nil {
		fmt.Println("❌ 配置已取消")
		return
	}

	switch provider {
	case "deepseek":
		baseURL = "https://api.deepseek.com"
		path = "/chat/completions"
	case "openai":
		baseURL = "https://api.openai.com"
		path = "/v1/chat/completions"
	case "grok":
		baseURL = "https://api.x.ai"
		path = "/v1/responses"
	case "claude":
		baseURL = "https://api.anthropic.com"
		path = "/v1/messages"
	}

	// 准备模型选项
	var modelOptions []huh.Option[string]
	if models, ok := providerModels[provider]; ok {
		for _, m := range models {
			val := m
			modelOptions = append(modelOptions, huh.NewOption(m, val))
		}
	}

	modelOptions = append(modelOptions, huh.NewOption("其他模型", "manual"))

	if provider != "custom" {
		var selectedModel string
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(fmt.Sprintf("选择 %s 模型", provider)).
					Options(modelOptions...).
					Value(&selectedModel),
			),
		).Run()
		if err != nil {
			return
		}

		if selectedModel != "manual" {
			model = selectedModel
		} else {
			model = ""
		}
	} else {
		model = ""
	}

	formFields := []huh.Field{
		huh.NewInput().
			Title("API Key").
			Value(&apiKey).
			EchoMode(huh.EchoModePassword),
		huh.NewSelect[string]().
			Title("提交日志语言").
			Options(
				huh.NewOption("中文", "cn"),
				huh.NewOption("English", "en"),
			).
			Value(&language),
		huh.NewConfirm().
			Title("是否生成详细描述?").
			Value(&withDescription),
		huh.NewInput().
			Title("分割符").
			Value(&subjectSeparateSymbol),
	}

	if model == "" {
		formFields = append(formFields, huh.NewInput().
			Title("请输入模型名称").
			Placeholder("e.g. gpt-4-turbo").
			Value(&model))
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
		Provider:              provider,
		APIKey:                apiKey,
		BaseURL:               baseURL,
		Path:                  path,
		Model:                 model,
		Language:              language,
		SubjectSeparateSymbol: subjectSeparateSymbol,
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
			"api_key":                 true,
			"model":                   true,
			"base_url":                true,
			"subject_separate_symbol": true,
		}

		if !validKeys[key] {
			fmt.Printf("❌ 无效的配置项: %s\n仅支持: api_key, model, base_url, subject_separate_symbol\n", key)
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
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.Flags().BoolVarP(&shouldStageAll, "add", "a", false, "Stage all files before commit")
}
