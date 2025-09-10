package google

import (
	"connector-demo/auth"
	"connector-demo/connectors/google/drive"
	"connector-demo/connectors/google/gmail"
	"connector-demo/utils"
)

// GoogleService 聚合各子模块的服务
type GoogleService struct {
	Connector *GoogleConnector
	Gmail     *gmail.GmailService
	Drive     *drive.DriveService
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

// TestConnection 测试Google连接，返回各平台测试状态
func (gs *GoogleService) TestConnection(userID string) map[string]bool {
	platforms := map[string]func(string) bool{
		auth.ProviderGmail:       func(uid string) bool { return gs.Gmail != nil && gs.Gmail.TestConnection(uid) },
		auth.ProviderGoogleDrive: func(uid string) bool { return gs.Drive != nil && gs.Drive.TestConnection(uid) },
	}
	result := make(map[string]bool, len(platforms))
	for k, test := range platforms {
		result[k] = test(userID)
	}
	return result
}
