package confluence

import (
	"context"
	"fmt"
	"net/http"

	"connector-demo/auth"
	"connector-demo/utils"

	confulence "github.com/ctreminiom/go-atlassian/v2/confluence/v2"
	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

// ConfluenceConnector 处理Confluence API调用
type ConfluenceConnector struct {
	tokenManager *utils.TokenManager
}

func NewConfluenceConnector(tm *utils.TokenManager) *ConfluenceConnector {
	return &ConfluenceConnector{tokenManager: tm}
}

// 获取Confluence客户端
func (sc *ConfluenceConnector) GetPages(cloudID string) (*models.PageChunkScheme, *models.ResponseScheme, error) {
	token, exists := sc.tokenManager.GetToken(cloudID, auth.ProviderConfluence)
	if !exists {
		return nil, nil, fmt.Errorf("未找到用户的Confluence token")
	}
	site := "https://api.atlassian.com/ex/confluence/" + cloudID + "/"

	client, err := confulence.New(
		http.DefaultClient,
		site,
	)
	if err != nil {
		return nil, nil, err
	}

	client.Auth.SetBearerToken(token.AccessToken)
	// 示例：获取Confluence页面列表
	return client.Page.Gets(context.Background(), nil, "", 10)
}
