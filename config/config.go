package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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

// LoadConfig 从环境变量加载配置，支持.env文件
func LoadConfig() *Config {
	// 尝试加载.env文件，如果存在则加载
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件,将使用环境变量")
	}

	return &Config{
		GoogleClientID:     GetEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: GetEnv("GOOGLE_CLIENT_SECRET", ""),
		SlackClientID:      GetEnv("SLACK_CLIENT_ID", ""),
		SlackClientSecret:  GetEnv("SLACK_CLIENT_SECRET", ""),
		RedirectURL:        GetEnv("REDIRECT_URL", "http://localhost:6767"),
		ServerPort:         GetEnv("PORT", "6767"),
	}
}

// GetEnv 获取环境变量，如果不存在则返回默认值
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
