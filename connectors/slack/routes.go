package slack

import (
	"connector-demo/routes"

	"github.com/gin-gonic/gin"
)

var slackService *SlackService

func SetSlackService(s *SlackService) {
	slackService = s
}

func RegisterRoutes(rg *gin.RouterGroup) {
	slackGroup := rg.Group("/slack")

	slackGroup.GET("/user-info", func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(400, gin.H{"error": "缺少user_id"})
			return
		}
		info, err := slackService.GetUserInfo(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"user_info": info})
	})

	slackGroup.GET("/channels", func(c *gin.Context) {
		userID := c.Query("user_id")
		channels, err := slackService.ListChannels(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"channels": channels})
	})

	slackGroup.GET("/test", func(c *gin.Context) {
		userID := c.Query("user_id")
		if !slackService.TestConnection(userID) {
			c.JSON(500, gin.H{"error": "Slack连接测试失败"})
			return
		}
		c.JSON(200, gin.H{"message": "Slack连接测试成功"})
	})
}

// 自动注册到 routes 模块
func init() {
	routes.RegisterModule("slack", RegisterRoutes)
}
