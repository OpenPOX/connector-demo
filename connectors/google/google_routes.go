package google

import (
	"connector-demo/connectors/google/drive"
	"connector-demo/connectors/google/gmail"

	"connector-demo/routes"

	"github.com/gin-gonic/gin"
)

var googleService *GoogleService

func SetGoogleService(gs *GoogleService) {
	googleService = gs
}

func RegisterRoutes(rg *gin.RouterGroup) {
	googleGroup := rg.Group("/google")

	// 注册 Gmail 子路由
	gmail.SetGmailService(googleService.Gmail)
	gmail.RegisterRoutes(googleGroup)

	// 注册 Drive 子路由
	drive.SetDriveService(googleService.Drive)
	drive.RegisterRoutes(googleGroup)

	googleGroup.GET("/test", func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(400, gin.H{"error": "缺少user_id参数"})
			return
		}

		status := googleService.TestConnection(userID)
		c.JSON(200, status)
	})

}

func init() {
	routes.RegisterModule("google", RegisterRoutes)
}
