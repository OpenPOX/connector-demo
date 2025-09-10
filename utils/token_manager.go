package utils

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/markbates/goth"
)

// TokenInfo 存储token信息
type TokenInfo struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
	Provider     string    `json:"provider,omitempty"`
}

// TokenManager 管理用户的token
// 注意：这是内存存储，生产环境应使用数据库
type TokenManager struct {
	tokens map[string]map[string]*TokenInfo // userID -> platform -> token
	mu     sync.RWMutex
}

// NewTokenManager 创建新的token管理器
func NewTokenManager() *TokenManager {
	return &TokenManager{
		tokens: make(map[string]map[string]*TokenInfo),
	}
}

// SaveToken 保存用户的token
func (tm *TokenManager) SaveToken(userID, platform string, token *TokenInfo) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.tokens[userID] == nil {
		tm.tokens[userID] = make(map[string]*TokenInfo)
	}
	tm.tokens[userID][platform] = token
	return nil
}

// GetToken 获取用户的token（自动刷新过期token）
func (tm *TokenManager) GetToken(userID, platform string) (*TokenInfo, bool) {
	tm.mu.RLock()
	userTokens, ok := tm.tokens[userID]
	tm.mu.RUnlock()
	if !ok {
		return nil, false
	}
	token, ok := userTokens[platform]
	if !ok || token == nil {
		return nil, false
	}
	// 判断token是否过期
	if !token.Expiry.IsZero() && token.Expiry.Before(time.Now()) {
		// 已过期，尝试刷新
		newToken, err := tm.RefreshToken(userID, platform)
		if err != nil {
			return nil, false
		}
		token = newToken
	}
	return token, true
}

func (tm *TokenManager) RefreshToken(userID, platform string) (*TokenInfo, error) {
	// 获取当前 token
	tm.mu.RLock()
	userTokens, ok := tm.tokens[userID]
	tm.mu.RUnlock()
	if !ok {
		return nil, errors.New("token not found")
	}
	token, ok := userTokens[platform]
	if !ok || token == nil {
		return nil, errors.New("token not found")
	}

	// 获取  provider
	provider, err := goth.GetProvider(platform)
	if err != nil {
		// provider 获取失败
		log.Printf("provider not found: %v", err)
		return nil, err
	}

	// 使用 goth provider 的 RefreshToken 方法
	newOAuthToken, err := provider.RefreshToken(token.RefreshToken)
	if err != nil {
		return nil, err
	}
	if newOAuthToken == nil {
		return nil, errors.New("refresh token returned nil")
	}

	// 更新内存 token
	newToken := &TokenInfo{
		AccessToken:  newOAuthToken.AccessToken,
		RefreshToken: newOAuthToken.RefreshToken,
		Expiry:       newOAuthToken.Expiry,
		TokenType:    newOAuthToken.TokenType,
		Provider:     platform,
	}

	err = tm.SaveToken(userID, platform, newToken)
	if err != nil {
		log.Printf("failed to save refreshed token: %v", err)
	}

	return newToken, nil
}

// GetAllTokens 获取用户的所有token
func (tm *TokenManager) GetAllTokens(userID string) map[string]*TokenInfo {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	if tokens, ok := tm.tokens[userID]; ok {
		// 返回副本避免并发问题
		result := make(map[string]*TokenInfo)
		for k, v := range tokens {
			result[k] = v
		}
		return result
	}
	return make(map[string]*TokenInfo)
}

// DeleteToken 删除用户的token
func (tm *TokenManager) DeleteToken(userID, platform string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if userTokens, ok := tm.tokens[userID]; ok {
		delete(userTokens, platform)
		if len(userTokens) == 0 {
			delete(tm.tokens, userID)
		}
	}
}

// PrintAllTokens 打印所有token（调试用）
func (tm *TokenManager) PrintAllTokens() {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	for userID, platforms := range tm.tokens {
		for platform, token := range platforms {
			log.Printf("User: %s Platform: %s AccessToken: %s RefreshToken: %s", userID, platform, token.AccessToken, token.RefreshToken)
		}
	}
}
