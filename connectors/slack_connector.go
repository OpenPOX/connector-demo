package connectors

import (
	"fmt"
	"log"

	"connector-demo/utils"

	"github.com/slack-go/slack"
)

// SlackConnector 处理Slack API调用
type SlackConnector struct {
	tokenManager *utils.TokenManager
}

// NewSlackConnector 创建新的Slack连接器
func NewSlackConnector(tm *utils.TokenManager) *SlackConnector {
	return &SlackConnector{
		tokenManager: tm,
	}
}

// GetUserInfo 获取Slack用户信息
func (sc *SlackConnector) GetUserInfo(userID string) (*slack.User, error) {
	token, exists := sc.tokenManager.GetToken(userID, "slack")
	if !exists {
		return nil, fmt.Errorf("未找到用户的Slack token")
	}

	api := slack.New(token.AccessToken)
	
	// 获取当前用户信息
	authTest, err := api.AuthTest()
	if err != nil {
		return nil, fmt.Errorf("Slack认证测试失败: %v", err)
	}

	user, err := api.GetUserInfo(authTest.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	return user, nil
}

// ListChannels 获取频道列表
func (sc *SlackConnector) ListChannels(userID string) ([]slack.Channel, error) {
	client, err := sc.getClient(userID)
	if err != nil {
		return nil, err
	}

	channels, _, err := client.GetConversations(&slack.GetConversationsParameters{
		Types: []string{"public_channel"},
	})
	return channels, err
}

// GetChannelMessages 获取频道消息
func (sc *SlackConnector) GetChannelMessages(userID string, channelID string, limit int) ([]slack.Message, error) {
	client, err := sc.getClient(userID)
	if err != nil {
		return nil, err
	}

	history, err := client.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	return history.Messages, nil
}

// getClient 获取Slack客户端
func (sc *SlackConnector) getClient(userID string) (*slack.Client, error) {
	token, exists := sc.tokenManager.GetToken(userID, "slack")
	if !exists {
		return nil, fmt.Errorf("未找到用户的Slack token")
	}

	return slack.New(token.AccessToken), nil
}

// ListRecentMessages 获取频道的最近消息
func (sc *SlackConnector) ListRecentMessages(userID, channelID string, limit int) ([]slack.Message, error) {
	return sc.GetChannelMessages(userID, channelID, limit)
}

// SendMessage 发送消息到指定频道
func (sc *SlackConnector) SendMessage(userID, channelID, message string) (string, error) {
	token, exists := sc.tokenManager.GetToken(userID, "slack")
	if !exists {
		return "", fmt.Errorf("未找到用户的Slack token")
	}

	api := slack.New(token.AccessToken)
	
	_, timestamp, err := api.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		return "", fmt.Errorf("发送消息失败: %v", err)
	}

	return timestamp, nil
}

// GetTeamInfo 获取团队信息
func (sc *SlackConnector) GetTeamInfo(userID string) (*slack.TeamInfo, error) {
	token, exists := sc.tokenManager.GetToken(userID, "slack")
	if !exists {
		return nil, fmt.Errorf("未找到用户的Slack token")
	}

	api := slack.New(token.AccessToken)
	
	teamInfo, err := api.GetTeamInfo()
	if err != nil {
		return nil, fmt.Errorf("获取团队信息失败: %v", err)
	}

	return teamInfo, nil
}

// ListUsers 获取团队成员列表
func (sc *SlackConnector) ListUsers(userID string) ([]slack.User, error) {
	token, exists := sc.tokenManager.GetToken(userID, "slack")
	if !exists {
		return nil, fmt.Errorf("未找到用户的Slack token")
	}

	api := slack.New(token.AccessToken)
	
	users, err := api.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %v", err)
	}

	return users, nil
}

// TestConnection 测试Slack连接
func (sc *SlackConnector) TestConnection(userID string) error {
	_, err := sc.GetUserInfo(userID)
	if err != nil {
		return fmt.Errorf("Slack连接测试失败: %v", err)
	}
	log.Printf("用户 %s 的Slack连接测试成功", userID)
	return nil
}

// ChannelInfo 频道信息结构体
type ChannelInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Purpose string `json:"purpose"`
	Topic   string `json:"topic"`
	Members int    `json:"members"`
	IsPrivate bool `json:"is_private"`
}

// GetChannelInfo 获取频道详细信息
func (sc *SlackConnector) GetChannelInfo(userID, channelID string) (*ChannelInfo, error) {
	token, exists := sc.tokenManager.GetToken(userID, "slack")
	if !exists {
		return nil, fmt.Errorf("未找到用户的Slack token")
	}

	api := slack.New(token.AccessToken)
	
	channel, err := api.GetConversationInfo(&slack.GetConversationInfoInput{ChannelID: channelID})
	if err != nil {
		return nil, fmt.Errorf("获取频道信息失败: %v", err)
	}

	return &ChannelInfo{
		ID:        channel.ID,
		Name:      channel.Name,
		Purpose:   channel.Purpose.Value,
		Topic:     channel.Topic.Value,
		Members:   channel.NumMembers,
		IsPrivate: channel.IsPrivate,
	}, nil
}

// GetMessagesWithUserInfo 获取带用户信息的消息列表
func (sc *SlackConnector) GetMessagesWithUserInfo(userID, channelID string, limit int) ([]MessageInfo, error) {
	messages, err := sc.ListRecentMessages(userID, channelID, limit)
	if err != nil {
		return nil, err
	}

	var messageInfos []MessageInfo
	for _, msg := range messages {
		messageInfos = append(messageInfos, MessageInfo{
			Text:      msg.Text,
			User:      msg.User,
			Timestamp: msg.Timestamp,
			Channel:   channelID,
		})
	}

	return messageInfos, nil
}

// MessageInfo 消息信息结构体
type MessageInfo struct {
	Text      string `json:"text"`
	User      string `json:"user"`
	Timestamp string `json:"timestamp"`
	Channel   string `json:"channel"`
}