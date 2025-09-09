package slack

import (
	"connector-demo/utils"
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

// SlackConnector 处理Slack API调用
type SlackConnector struct {
	tokenManager *utils.TokenManager
}

// SlackFile 封装 Slack 消息中附件信息
type SlackFile struct {
	ID       string // 文件 ID
	Name     string // 文件名
	MimeType string // 文件类型
	URL      string // 下载链接（需 token）
}

// SlackMessage 封装 Slack 消息
type SlackMessage struct {
	ID                 string // 全局唯一 ID，格式：ts + channelID
	ChannelID          string // 消息所在频道 / 对话 ID
	ChannelName        string // 频道 / 对话名称
	ChannelType        string
	UserID             string      // 发送者 ID
	UserName           string      // 发送者名字
	Text               string      // 消息文本内容
	Files              []SlackFile // 消息中附件，可选
	Timestamp          string      // 消息时间
	ThreadTS           string      // 所属线程 ID，如果是线程消息
	EditedTime         *time.Time  // 可选，消息被编辑时间
	Reactions          []string    // 可选，消息表情
	Permalink          string      // 可选，消息永久链接
	BotID              string      // 可选，机器人发送者 ID
	IsPinned           bool        // 可选，是否置顶
	ThreadRepliesCount int         // 可选，线程回复数
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

// ListMessages 获取指定 channel 的历史消息
func (sc *SlackConnector) ListMessages(userID, channelID string, limit int, oldest, latest string) ([]SlackMessage, error) {
	client, err := sc.getClient(userID)
	if err != nil {
		return nil, err
	}

	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     limit,
		Oldest:    oldest, // 可选，RFC3339 或 Slack ts 格式
		Latest:    latest, // 可选
		Inclusive: false,
	}

	history, err := client.GetConversationHistory(params)
	if err != nil {
		return nil, fmt.Errorf("获取 channel 消息失败: %v", err)
	}

	var messages []SlackMessage
	for _, m := range history.Messages {
		msg := SlackMessage{
			ID:          fmt.Sprintf("%s:%s", m.Timestamp, channelID), // 全局唯一 ID
			ChannelID:   channelID,
			ChannelName: "",               // 可在调用前通过 channel list 映射填充
			ChannelType: "public_channel", // 填充频道类型 public_channel/private_channel/im/mpim
			UserID:      m.User,
			UserName:    "", // 可在调用前通过 users.list 映射填充
			Text:        m.Text,
			Timestamp:   m.Timestamp,
			ThreadTS:    m.ThreadTimestamp,
		}

		// 处理文件
		if len(m.Files) > 0 {
			for _, f := range m.Files {
				msg.Files = append(msg.Files, SlackFile{
					ID:       f.ID,
					Name:     f.Name,
					MimeType: f.Mimetype,
					URL:      f.URLPrivate,
				})
			}
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

// 更多原始API方法保持在这里，例如：GetChannelMessages、SendMessage 等
