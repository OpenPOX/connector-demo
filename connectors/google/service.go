package google

import (
	"connector-demo/connectors/google/drive"
	"connector-demo/connectors/google/gmail"
	"connector-demo/utils"
)

// GoogleService 聚合各子模块的服务
type GoogleService struct {
	Connector *GoogleConnector
	Gmail     *gmail.Service
	Drive     *drive.Service
}

// NewGoogleService 创建新的Google服务聚合
func NewGoogleService(tokenManager *utils.TokenManager) *GoogleService {
	googleConnector := NewGoogleConnector(tokenManager)
	// 创建各个连接器
	gmailConnector := gmail.NewGmailConnector(tokenManager)
	driveConnector := drive.NewDriveConnector(tokenManager)

	// 创建各个服务
	gmailService := gmail.NewService(gmailConnector)
	driveService := drive.NewService(driveConnector)

	return &GoogleService{
		Connector: googleConnector,
		Gmail:     gmailService,
		Drive:     driveService,
	}
}

// GetUserInfo 获取用户信息
func (gs *GoogleService) GetUserInfo(userID string) (*UserInfo, error) {
	return gs.Connector.GetUserInfo(userID)
}

// TestConnection 测试Google连接
func (gs *GoogleService) TestConnection(userID string) bool {
	if err := gs.Connector.TestConnection(userID); err != nil {
		return false
	}

	// 可选：测试 Gmail/Drive 子模块
	// if gs.Gmail != nil && !gs.Gmail.Test(userID) {
	// 	return false
	// }
	// if gs.Drive != nil && !gs.Drive.Test(userID) {
	// 	return false
	// }

	return true
}
