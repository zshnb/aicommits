package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 结构体定义了我们的配置项
type Config struct {
	Provider              string `mapstructure:"provider"`
	APIKey                string `mapstructure:"api_key"`
	Model                 string `mapstructure:"model"`
	BaseURL               string `mapstructure:"base_url"`
	Language              string `mapstructure:"language"`
	WithDescription       bool   `mapstructure:"with_description"`
	SubjectSeparateSymbol string `mapstructure:"subject_separate_symbol"`
}

// init 初始化 Viper 配置
func init() {
	// 配置文件名 (不带后缀)
	viper.SetConfigName(".aicommits")
	// 配置文件类型
	viper.SetConfigType("yaml")
	// 查找路径: 用户主目录
	home, _ := os.UserHomeDir()
	viper.AddConfigPath(home)
}

// Load 读取配置
func Load() (*Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		// 如果是“未找到配置文件”错误，返回空配置即可，不报错
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return &Config{
				// 设置默认值
				Model:   "gpt-5-nano",
				BaseURL: "https://api.openai.com/v1",
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Set 更新某一项配置并保存到磁盘
func Set(key, value string) error {
	viper.Set(key, value)

	// 确保配置文件存在，如果不存在则创建
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			home, _ := os.UserHomeDir()
			configPath := filepath.Join(home, ".aicommits.yaml")
			// 创建空文件
			if _, err := os.Create(configPath); err != nil {
				return err
			}
			// 重新读取一遍以便 Viper 绑定文件
			viper.SetConfigFile(configPath)
		}
	}

	return viper.WriteConfig()
}

func Save(cfg *Config) error {
	viper.Set("provider", cfg.Provider)
	viper.Set("api_key", cfg.APIKey)
	viper.Set("model", cfg.Model)
	viper.Set("base_url", cfg.BaseURL)
	viper.Set("language", cfg.Language)
	viper.Set("with_description", cfg.WithDescription)
	viper.Set("subject_separate_symbol", cfg.SubjectSeparateSymbol)

	// 确保文件存在
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			home, _ := os.UserHomeDir()
			configPath := filepath.Join(home, ".aicommits.yaml")
			os.Create(configPath)
			viper.SetConfigFile(configPath)
		}
	}
	return viper.WriteConfig()
}

func GetPrintable() string {
	key := viper.GetString("api_key")
	if len(key) > 8 {
		key = key[:4] + "..." + key[len(key)-4:]
	} else if key != "" {
		key = "***"
	} else {
		key = "(未设置)"
	}

	return fmt.Sprintf(`
Current Configuration:
  Provider: %s
  Model:    %s
  API Key:  %s
  Subject Separate Symbol: %s
`, viper.GetString("provider"), viper.GetString("model"), key, viper.GetString("subject_separate_symbol"))
}
