package auth

import (
	"connector-demo/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// AuthHandler 处理OAuth2认证相关请求
type AuthHandler struct {
	tokenManager *utils.TokenManager
}

// NewAuthHandler 创建新的认证处理器
func NewAuthHandler(tm *utils.TokenManager) *AuthHandler {
	return &AuthHandler{
		tokenManager: tm,
	}
}

// Connect 处理连接请求，重定向到OAuth2授权页面
func (ah *AuthHandler) Connect(c *gin.Context) {
	provider := c.Param("provider")
	if !IsSupportedProvider(provider) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("不支持的平台: %s", provider)})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// Callback 处理OAuth2回调
func (ah *AuthHandler) Callback(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(400, gin.H{"error": "provider参数不能为空"})
		return
	}

	// 设置provider
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	// 完成OAuth2认证
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("OAuth2认证失败: %v", err)})
		return
	}

	// 获取用户ID（这里使用用户邮箱作为唯一标识）
	userID := user.Email
	if userID == "" {
		userID = user.UserID
	}

	// 保存token信息
	tokenInfo := &utils.TokenInfo{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Expiry:       user.ExpiresAt,
		Provider:     provider,
	}

	if err := ah.tokenManager.SaveToken(userID, provider, tokenInfo); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("保存token失败: %v", err)})
		return
	}

	// c.JSON(200, gin.H{
	// 	"message":  fmt.Sprintf("%s 认证成功", provider),
	// 	"user_id":  userID,
	// 	"provider": provider,
	// })

	// OAuth2成功后，重定向到前端页面
	frontendURL := fmt.Sprintf("http://localhost:6767")
	c.Redirect(302, frontendURL)
}

// GetTokens 获取用户的所有token
func (ah *AuthHandler) GetTokens(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少user_id参数"})
		return
	}

	tokens := ah.tokenManager.GetAllTokens(userID)
	c.JSON(http.StatusOK, gin.H{"user_id": userID, "tokens": tokens})
}

// Disconnect 断开指定平台的连接
func (ah *AuthHandler) Disconnect(c *gin.Context) {
	userID := c.Query("user_id")
	provider := c.Param("provider")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少user_id参数"})
		return
	}

	if !IsSupportedProvider(provider) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("不支持的平台: %s", provider)})
		return
	}

	ah.tokenManager.DeleteToken(userID, provider)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("已断开与 %s 的连接", provider)})
}

// addProviderToContext 将provider添加到请求上下文中
func addProviderToContext(r *http.Request, provider string) *http.Request {
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()
	return r
}
