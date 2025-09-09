package gmail

import (
	"context"
	"fmt"

	"connector-demo/utils"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GmailConnector Gmail API封装
type GmailConnector struct {
	tokenManager *utils.TokenManager
}

// Message Gmail邮件信息
type Message struct {
	ID       string   `json:"id"`
	ThreadID string   `json:"threadId"`
	Snippet  string   `json:"snippet"`
	Subject  string   `json:"subject"`
	From     string   `json:"from"`
	Date     string   `json:"date"`
	LabelIDs []string `json:"labelIds"`
}

// NewGmailConnector 创建新的Gmail连接器
func NewGmailConnector(tm *utils.TokenManager) *GmailConnector {
	return &GmailConnector{
		tokenManager: tm,
	}
}

// GetService 获取Gmail服务客户端
func (gc *GmailConnector) GetService(userID string) (*gmail.Service, error) {
	tokenInfo, exists := gc.tokenManager.GetToken(userID, "google")
	if !exists {
		return nil, fmt.Errorf("未找到Google访问令牌")
	}

	// 创建OAuth2客户端
	client := utils.CreateOAuth2Client(tokenInfo.AccessToken)
	
	// 创建Gmail服务
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("创建Gmail服务失败: %v", err)
	}

	return service, nil
}

// ListMessages 获取邮件列表
func (gc *GmailConnector) ListMessages(userID string, maxResults int64) ([]Message, error) {
	service, err := gc.GetService(userID)
	if err != nil {
		return nil, err
	}

	// 获取邮件列表
	messages, err := service.Users.Messages.List("me").MaxResults(maxResults).Do()
	if err != nil {
		return nil, fmt.Errorf("获取邮件列表失败: %v", err)
	}

	var result []Message
	for _, msg := range messages.Messages {
		// 获取邮件详情
		fullMsg, err := service.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			continue // 跳过失败的邮件
		}

		message := Message{
			ID:       msg.Id,
			ThreadID: msg.ThreadId,
			Snippet:  fullMsg.Snippet,
			LabelIDs: fullMsg.LabelIds,
		}

		// 提取邮件头信息
		for _, header := range fullMsg.Payload.Headers {
			switch header.Name {
			case "Subject":
				message.Subject = header.Value
			case "From":
				message.From = header.Value
			case "Date":
				message.Date = header.Value
			}
		}

		result = append(result, message)
	}

	return result, nil
}

// GetMessage 获取单封邮件详情
func (gc *GmailConnector) GetMessage(userID string, messageID string) (*Message, error) {
	service, err := gc.GetService(userID)
	if err != nil {
		return nil, err
	}

	fullMsg, err := service.Users.Messages.Get("me", messageID).Do()
	if err != nil {
		return nil, fmt.Errorf("获取邮件详情失败: %v", err)
	}

	message := Message{
		ID:       fullMsg.Id,
		ThreadID: fullMsg.ThreadId,
		Snippet:  fullMsg.Snippet,
		LabelIDs: fullMsg.LabelIds,
	}

	// 提取邮件头信息
	for _, header := range fullMsg.Payload.Headers {
		switch header.Name {
		case "Subject":
			message.Subject = header.Value
		case "From":
			message.From = header.Value
		case "Date":
			message.Date = header.Value
		}
	}

	return &message, nil
}