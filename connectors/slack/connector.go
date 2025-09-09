package slack

import (
	"connector-demo/utils"
	"fmt"

	"github.com/slack-go/slack"
)

// SlackConnector 处理Slack API调用
type SlackConnector struct {
	tokenManager *utils.TokenManager
}

func NewSlackConnector(tm *utils.TokenManager) *SlackConnector {
	return &SlackConnector{tokenManager: tm}
}

// 获取Slack客户端
func (sc *SlackConnector) getClient(userID string) (*slack.Client, error) {
	token, exists := sc.tokenManager.GetToken(userID, "slack")
	if !exists {
		return nil, fmt.Errorf("未找到用户的Slack token")
	}
	return slack.New(token.AccessToken), nil
}

// 原始API调用封装
func (sc *SlackConnector) GetUserInfo(userID string) (*slack.User, error) {
	client, err := sc.getClient(userID)
	if err != nil {
		return nil, err
	}
	authTest, err := client.AuthTest()
	if err != nil {
		return nil, fmt.Errorf("Slack认证测试失败: %v", err)
	}
	return client.GetUserInfo(authTest.UserID)
}

func (sc *SlackConnector) ListChannels(userID string) ([]slack.Channel, error) {
	client, err := sc.getClient(userID)
	if err != nil {
		return nil, err
	}
	channels, _, err := client.GetConversations(&slack.GetConversationsParameters{Types: []string{"public_channel"}})
	return channels, err
}

// 更多原始API方法保持在这里，例如：GetChannelMessages、SendMessage 等
