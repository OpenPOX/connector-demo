package google

import (
	"context"
	"fmt"

	"connector-demo/utils"
	"google.golang.org/api/option"
	"google.golang.org/api/oauth2/v2"
)

// GoogleConnector Google连接器统一入口，管理token
type GoogleConnector struct {
	tokenManager   *utils.TokenManager
}

// UserInfo 用户信息结构体
type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Picture string `json:"picture"`
}

// NewGoogleConnector 创建新的Google连接器
func NewGoogleConnector(tm *utils.TokenManager) *GoogleConnector {
	return &GoogleConnector{
		tokenManager: tm,
	}
}

// GetUserInfo 获取用户信息，使用OAuth2 API
func (gc *GoogleConnector) GetUserInfo(userID string) (*UserInfo, error) {
	service, err := gc.getOAuth2Service(userID)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	return &UserInfo{
		ID:      userInfo.Id,
		Name:    userInfo.Name,
		Email:   userInfo.Email,
		Picture: userInfo.Picture,
	}, nil
}

// getOAuth2Service 获取OAuth2服务客户端
func (gc *GoogleConnector) getOAuth2Service(userID string) (*oauth2.Service, error) {
	tokenInfo, exists := gc.tokenManager.GetToken(userID, "google")
	if !exists {
		return nil, fmt.Errorf("未找到Google访问令牌")
	}

	// 创建OAuth2客户端
	client := utils.CreateOAuth2Client(tokenInfo.AccessToken)
	
	// 创建OAuth2服务
	service, err := oauth2.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("创建OAuth2服务失败: %v", err)
	}

	return service, nil
}

// TestConnection 测试Google连接
func (gc *GoogleConnector) TestConnection(userID string) error {
	_, err := gc.GetUserInfo(userID)
	return err
}