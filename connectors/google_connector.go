package connectors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"connector-demo/utils"
)

// GoogleConnector 处理Google API调用
type GoogleConnector struct {
	tokenManager *utils.TokenManager
}

// NewGoogleConnector 创建Google连接器
func NewGoogleConnector(tm *utils.TokenManager) *GoogleConnector {
	return &GoogleConnector{
		tokenManager: tm,
	}
}

// UserInfo 用户信息结构体
type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// GetUserInfo 获取Google用户信息
func (gc *GoogleConnector) GetUserInfo(userID string) (*UserInfo, error) {
	tokenInfo, exists := gc.tokenManager.GetToken(userID, "google")
	if !exists {
		return nil, fmt.Errorf("token not found")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+tokenInfo.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// TestConnection 测试Google API连接
func (gc *GoogleConnector) TestConnection(userID string) error {
	_, err := gc.GetUserInfo(userID)
	return err
}