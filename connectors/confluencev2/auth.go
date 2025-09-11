// Package confluencev2 implements the Confluence V2 API client
package confluencev2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	APIAuthToken               = "https://auth.atlassian.com/oauth/token"
	APIAuthAccessibleResources = "https://api.atlassian.com/oauth/token/accessible-resources"
)

// OAuth2 Confluence OAuth2服务
type OAuth2 struct {
	client       http.Client
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

// AuthAccessToken 认证访问令牌
type AuthAccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// NewOAuth2 创建OAuth2认证
func NewOAuth2(clientID, secret, redirectURI string) *OAuth2 {
	return &OAuth2{
		client:       http.Client{},
		ClientID:     clientID,
		ClientSecret: secret,
		RedirectURI:  redirectURI,
	}
}

// AuthorizationCallBack 授权回调
func AuthorizationCallBack(ctx *gin.Context) {
	code := ctx.MustGet("code").(string)
	state := ctx.MustGet("state").(string)
	fmt.Println(code, state)
	ctx.JSON(200, nil)
}

func (o *OAuth2) AuthorizationCode(code string) (*AuthAccessToken, error) {
	api := APIAuthToken
	out := &AuthAccessToken{}
	body := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     o.ClientID,
		"client_secret": o.ClientSecret,
		"code":          code,
		"redirect_uri":  o.RedirectURI,
	}
	bs, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := o.client.Post(api, "application/json", bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, err
	}
	return out, nil
}

// RefreshToken 使用刷新令牌获取新的访问令牌
func (o OAuth2) RefreshToken(refreshToken string) (*AuthAccessToken, error) {
	api := APIAuthToken
	out := &AuthAccessToken{}
	body := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     o.ClientID,
		"client_secret": o.ClientSecret,
		"refresh_token": refreshToken,
	}
	bs, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := o.client.Post(api, "application/json", bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, err
	}
	return out, nil
}

// AccessibleResourcesResponse 访问资源响应
type AccessibleResourcesResponse []struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes"`
	AvatarURL string   `json:"avatarUrl"`
}

// AccessibleResources 获取可访问的资源
func (o OAuth2) AccessibleResources(accessToken string) (AccessibleResourcesResponse, error) {
	api := APIAuthAccessibleResources
	out := AccessibleResourcesResponse{}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}
