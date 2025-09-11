package auth

import (
	"fmt"

	"connector-demo/auth/providers/confluence"
	"connector-demo/config"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/slack"
)

const (
	ProviderGmail       = "gmail"
	ProviderGoogleDrive = "google-drive"
	ProviderSlack       = "slack"
	ProviderConfluence  = "confluence"
)

var SupportedProviders = []string{ProviderGmail, ProviderGoogleDrive, ProviderSlack, ProviderConfluence}

// SetupProviders 配置OAuth2提供者
func SetupProviders(cfg *config.Config) error {
	providers := []goth.Provider{}

	// 配置Google提供者
	if cfg.GoogleClientID != "" && cfg.GoogleClientSecret != "" {
		gmailProvider := google.New(
			cfg.GoogleClientID,
			cfg.GoogleClientSecret,
			fmt.Sprintf("%s/auth/gmail/callback", cfg.RedirectURL),
			"openid",
			"email",
			"profile",
			"https://www.googleapis.com/auth/gmail.readonly",
		)
		gmailProvider.SetName(ProviderGmail)
		gmailProvider.SetAccessType("offline")
		gmailProvider.SetPrompt("consent")
		providers = append(providers, gmailProvider)

		googleDriveProvider := google.New(
			cfg.GoogleClientID,
			cfg.GoogleClientSecret,
			fmt.Sprintf("%s/auth/google-drive/callback", cfg.RedirectURL),
			"openid",
			"email",
			"profile",
			"https://www.googleapis.com/auth/drive.readonly",
		)
		googleDriveProvider.SetName(ProviderGoogleDrive)
		googleDriveProvider.SetAccessType("offline")
		googleDriveProvider.SetPrompt("consent")
		providers = append(providers, googleDriveProvider)
	}

	// 配置Slack提供者
	if cfg.SlackClientID != "" && cfg.SlackClientSecret != "" {
		slackProvider := slack.New(
			cfg.SlackClientID,
			cfg.SlackClientSecret,
			fmt.Sprintf("%s/auth/slack/callback", cfg.RedirectURL),
			"channels:read",
			"groups:read",
			"im:read",
			"mpim:read",
			"channels:history",
			"groups:history",
			"im:history",
			"mpim:history",
			"users:read",
		)
		providers = append(providers, slackProvider)
	}

	if cfg.ConfluenceClientID != "" && cfg.ConfluenceClientSecret != "" {
		confluenceProvider := confluence.New(
			cfg.ConfluenceClientID,
			cfg.ConfluenceClientSecret,
			fmt.Sprintf("%s/auth/confluence/callback", cfg.RedirectURL),
		)
		providers = append(providers, confluenceProvider)
	}

	if len(providers) == 0 {
		return fmt.Errorf("没有配置任何OAuth2提供者")
	}

	goth.UseProviders(providers...)
	return nil
}

// GetProvider 获取指定平台的提供者
func GetProvider(platform string) (goth.Provider, error) {
	provider, err := goth.GetProvider(platform)
	if err != nil {
		return nil, fmt.Errorf("未找到提供者: %s", platform)
	}
	return provider, nil
}

// GetSupportedProviders 获取所有支持的提供者
func GetSupportedProviders() []string {
	return SupportedProviders
}

// IsSupportedProvider 检查是否支持指定平台
func IsSupportedProvider(platform string) bool {
	supported := GetSupportedProviders()
	for _, p := range supported {
		if p == platform {
			return true
		}
	}
	return false
}
