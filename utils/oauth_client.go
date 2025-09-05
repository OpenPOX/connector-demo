package utils

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

// HTTPClient 接口用于支持不同的HTTP客户端
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CreateOAuth2Client 创建OAuth2 HTTP客户端
func CreateOAuth2Client(accessToken string) HTTPClient {
	token := &oauth2.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}

	// 创建一个空的OAuth2配置，因为我们只需要token
	// 实际的token验证和刷新由各个SDK处理
	config := &oauth2.Config{}
	
	return config.Client(context.Background(), token)
}

// CreateStandardClient 创建标准HTTP客户端（用于测试）
func CreateStandardClient() *http.Client {
	return &http.Client{}
}