package gmail

import (
	"context"
	"encoding/base64"
	"fmt"

	"connector-demo/auth"
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
	ID           string   `json:"id"`
	ThreadID     string   `json:"threadId"`
	LabelIDs     []string `json:"labelIds"`
	Snippet      string   `json:"snippet"`
	Subject      string   `json:"subject"`
	To           string   `json:"to"`
	From         string   `json:"from"`
	Date         string   `json:"date"`
	InternalDate string   `json:"internalDate"`
	AttachmentID string   `json:"attachmentId"`
	Data         string   `json:"data"`
}

// NewGmailConnector 创建新的Gmail连接器
func NewGmailConnector(tm *utils.TokenManager) *GmailConnector {
	return &GmailConnector{
		tokenManager: tm,
	}
}

// GetService 获取Gmail服务客户端
func (gc *GmailConnector) GetService(userID string) (*gmail.Service, error) {
	tokenInfo, exists := gc.tokenManager.GetToken(userID, auth.ProviderGmail)
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

// 封装：将 gmail.Message 转换为本地 Message 对象
func parseGmailMessage(msg *gmail.Message) Message {
	message := Message{
		ID:       msg.Id,
		ThreadID: msg.ThreadId,
		Snippet:  msg.Snippet,
		LabelIDs: msg.LabelIds,
	}

	if msg.Payload != nil {
		// 提取邮件头信息
		for _, header := range msg.Payload.Headers {
			switch header.Name {
			case "Subject":
				message.Subject = header.Value
			case "From":
				message.From = header.Value
			case "Date":
				message.Date = header.Value
			case "To":
				message.To = header.Value
			}
		}

		// 提取邮件正文（body）和附件ID
		var bodyData string
		var attachmentID string
		payload := msg.Payload
		if len(payload.Parts) > 0 {
			for _, part := range payload.Parts {
				if part.MimeType == "text/plain" || part.MimeType == "text/html" {
					if part.Body != nil && part.Body.Data != "" {
						decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
						if err == nil {
							bodyData = string(decoded)
						} else {
							bodyData = part.Body.Data // fallback
						}
					}
				}
				// 提取附件ID
				if part.Filename != "" && part.Body != nil && part.Body.AttachmentId != "" {
					attachmentID = part.Body.AttachmentId
				}
			}
		} else if payload.Body != nil {
			if payload.Body.Data != "" {
				decoded, err := base64.URLEncoding.DecodeString(payload.Body.Data)
				if err == nil {
					bodyData = string(decoded)
				} else {
					bodyData = payload.Body.Data // fallback
				}
			}
			if payload.Body.AttachmentId != "" {
				attachmentID = payload.Body.AttachmentId
			}
		}
		message.Data = bodyData
		message.AttachmentID = attachmentID
	}

	return message
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
		fullMsg, err := service.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			continue // 跳过失败的邮件
		}
		result = append(result, parseGmailMessage(fullMsg))
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

	msg := parseGmailMessage(fullMsg)
	return &msg, nil
}
