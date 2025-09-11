package confluence

import (
	"context"
	"fmt"
	"net/http"

	"connector-demo/auth"
	"connector-demo/utils"

	confulence "github.com/ctreminiom/go-atlassian/v2/confluence/v2"
)

// ConfluenceConnector 处理Confluence API调用
type ConfluenceConnector struct {
	tokenManager *utils.TokenManager
}

func NewConfluenceConnector(tm *utils.TokenManager) *ConfluenceConnector {
	return &ConfluenceConnector{tokenManager: tm}
}

// 获取Confluence客户端
func (sc *ConfluenceConnector) getPages(userID string) (interface{}, error) {
	token, exists := sc.tokenManager.GetToken(userID, auth.ProviderConfluence)
	if !exists {
		return nil, fmt.Errorf("未找到用户的Confluence token")
	}
	site := "https://your-domain.atlassian.net/wiki" // 替换为你的Confluence站点URL
	client, err := confulence.New(http.DefaultClient, site)
	if err != nil {
		return nil, err
	}
	client.Auth.SetBearerToken(token.AccessToken)
	// 示例：获取Confluence页面列表
	pages, _, err := client.Page.Gets(context.Background(), nil, "", 0)
	if err != nil {
		return nil, err
	}
	return pages, nil
}
