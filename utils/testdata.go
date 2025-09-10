package utils

import (
	"log"

	"connector-demo/config"
)

// InjectTestTokens 注入测试 token
func InjectTestTokens(tm *TokenManager) {
	if config.GetEnv("ENABLE_TEST_TOKENS", "false") != "true" {
		return
	}

	log.Println("检测到测试token配置，正在注入...")

	// 注入 Google 测试 token
	if accessToken := config.GetEnv("TEST_TOKEN_GOOGLE_ACCESS", ""); accessToken != "" {
		googleToken := &TokenInfo{
			AccessToken:  accessToken,
			RefreshToken: config.GetEnv("TEST_TOKEN_GOOGLE_REFRESH", ""),
			Provider:     "google",
			TokenType:    "Bearer",
		}
		tm.SaveToken("1", "gmail", googleToken)
		tm.SaveToken("1", "google-drive", googleToken)

		tm.RefreshToken("1", "gmail")
		tm.RefreshToken("1", "google-drive")
		log.Println("已注入Google测试token")
	}

	// 注入 Slack 测试 token
	if accessToken := config.GetEnv("TEST_TOKEN_SLACK_ACCESS", ""); accessToken != "" {
		slackToken := &TokenInfo{
			AccessToken:  accessToken,
			RefreshToken: config.GetEnv("TEST_TOKEN_SLACK_REFRESH", ""),
			Provider:     "slack",
			TokenType:    "Bearer",
		}
		tm.SaveToken("1", "slack", slackToken)
		tm.RefreshToken("1", "slack")
		log.Println("已注入Slack测试token")
	}
}
