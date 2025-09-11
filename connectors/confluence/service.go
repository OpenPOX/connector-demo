package confluence

import (
	"connector-demo/utils"
)

// ConfluenceService 负责封装业务逻辑，调用 SlackConnector
type ConfluenceService struct {
	connector *ConfluenceConnector
}

func NewConfluenceService(tokenManager *utils.TokenManager) *ConfluenceService {
	slackConnector := NewConfluenceConnector(tokenManager)
	return &ConfluenceService{connector: slackConnector}
}
