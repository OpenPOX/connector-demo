package confluence

import (
	"connector-demo/routes"

	"github.com/gin-gonic/gin"
)

var confluenceService *ConfluenceService

func SetConfluenceService(s *ConfluenceService) {
	confluenceService = s
}

func RegisterRoutes(rg *gin.RouterGroup) {

}

// 自动注册到 routes 模块
func init() {
	routes.RegisterModule("confluence", RegisterRoutes)
}
