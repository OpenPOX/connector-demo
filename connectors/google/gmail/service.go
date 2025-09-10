package gmail

import "fmt"

// Service Gmail数据处理接口
type GmailService struct {
	connector *GmailConnector
}

// NewService 创建新的Gmail服务
func NewService(connector *GmailConnector) *GmailService {
	return &GmailService{
		connector: connector,
	}
}

// GetInboxMessages 获取收件箱邮件
func (s *GmailService) GetInboxMessages(userID string, limit int64) ([]Message, error) {
	messages, err := s.connector.ListMessages(userID, limit)
	if err != nil {
		return nil, fmt.Errorf("获取收件箱邮件失败: %v", err)
	}
	return messages, nil
}

// GetMessageDetail 获取邮件详情
func (s *GmailService) GetMessageDetail(userID string, messageID string) (*Message, error) {
	message, err := s.connector.GetMessage(userID, messageID)
	if err != nil {
		return nil, fmt.Errorf("获取邮件详情失败: %v", err)
	}
	return message, nil
}

// GetUnreadCount 获取未读邮件数量
func (s *GmailService) GetUnreadCount(userID string) (int64, error) {
	// 这里可以实现获取未读邮件数量的逻辑
	// 暂时返回0，后续可以扩展
	return 0, nil
}

// TestConnection 测试Gmail连接
func (s *GmailService) TestConnection(userID string) bool {
	_, err := s.connector.GetUserInfo(userID)
	if err != nil {
		return false
	}
	return true
}
