package utils

import (
	"log"
	"sync"
	"time"
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

// GetToken 获取用户的token
func (tm *TokenManager) GetToken(userID, platform string) (*TokenInfo, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	if userTokens, ok := tm.tokens[userID]; ok {
		if token, ok := userTokens[platform]; ok {
			return token, true
		}
	}
	return nil, false
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