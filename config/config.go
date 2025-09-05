package config

import (
	"os"
)

// Config 存储所有配置信息
type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	SlackClientID      string
	SlackClientSecret  string
	RedirectURL        string
	ServerPort         string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	return &Config{
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		SlackClientID:      getEnv("SLACK_CLIENT_ID", ""),
		SlackClientSecret:  getEnv("SLACK_CLIENT_SECRET", ""),
		RedirectURL:        getEnv("REDIRECT_URL", "http://localhost:6767"),
		ServerPort:         getEnv("PORT", "8080"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsValid 检查配置是否有效
func (c *Config) IsValid() bool {
	return c.GoogleClientID != "" && c.GoogleClientSecret != "" &&
		c.SlackClientID != "" && c.SlackClientSecret != ""
}
