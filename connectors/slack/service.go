package slack

import (
	"connector-demo/utils"
	"log"

	"github.com/slack-go/slack"
)

// SlackService 负责封装业务逻辑，调用 SlackConnector
type SlackService struct {
	connector *SlackConnector
}

func NewSlackService(tokenManager *utils.TokenManager) *SlackService {
	slackConnector := NewSlackConnector(tokenManager)
	return &SlackService{connector: slackConnector}
}

// 获取用户信息
func (s *SlackService) GetUserInfo(userID string) (*slack.User, error) {
	return s.connector.GetUserInfo(userID)
}

// 获取频道列表
func (s *SlackService) ListChannels(userID string) ([]slack.Channel, error) {
	return s.connector.ListChannels(userID)
}

// 获取指定频道的历史消息
func (s *SlackService) ListMessages(userID, channelID string, limit int, oldest, latest string) ([]SlackMessage, error) {
	return s.connector.ListMessages(userID, channelID, limit, oldest, latest)
}

// 测试连接，返回bool
func (s *SlackService) TestConnection(userID string) bool {
	_, err := s.connector.GetUserInfo(userID)
	if err != nil {
		log.Printf("Slack连接测试失败: %v", err)
		return false
	}
	log.Printf("Slack连接测试成功: %s", userID)
	return true
}

// 可以在这里封装更多组合业务逻辑，如 GetMessagesWithUserInfo、发送消息等
