package google

import (
	"connector-demo/utils"
)

// GoogleConnector Google连接器统一入口，管理token
type GoogleConnector struct {
	tokenManager *utils.TokenManager
}

// NewGoogleConnector 创建新的Google连接器
func NewGoogleConnector(tm *utils.TokenManager) *GoogleConnector {
	return &GoogleConnector{
		tokenManager: tm,
	}
}

// getOAuth2Service 获取OAuth2服务客户端
// func (gc *GoogleConnector) getOAuth2Service(userID string) (*oauth2.Service, error) {
// 	tokenInfo, exists := gc.tokenManager.GetToken(userID, "google")
// 	if !exists {
// 		return nil, fmt.Errorf("未找到Google访问令牌")
// 	}

// 	// 创建OAuth2客户端
// 	client := utils.CreateOAuth2Client(tokenInfo.AccessToken)

// 	// 创建OAuth2服务
// 	service, err := oauth2.NewService(context.Background(), option.WithHTTPClient(client))
// 	if err != nil {
// 		return nil, fmt.Errorf("创建OAuth2服务失败: %v", err)
// 	}

// 	return service, nil
// }
