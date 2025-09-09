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

	googleGroup.GET("/user-info", func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(400, gin.H{"error": "缺少user_id"})
			return
		}
		userInfo, err := googleService.GetUserInfo(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"user_info": userInfo})
	})

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

		if ok := googleService.TestConnection(userID); !ok {
			c.JSON(500, gin.H{"error": "Google连接测试失败"})
			return
		}
		c.JSON(200, gin.H{"message": "Google连接测试成功"})
	})
}

func init() {
	routes.RegisterModule("google", RegisterRoutes)
}
