package slack

import (
	"connector-demo/routes"
	"strconv"

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

	// 获取消息列表
	slackGroup.GET("/messages/:channel_id", func(c *gin.Context) {
		userID := c.Query("user_id")
		channelID := c.Param("channel_id")
		if userID == "" || channelID == "" {
			c.JSON(400, gin.H{"error": "缺少 user_id 或 channel_id"})
			return
		}

		// limit 参数，可选，默认 100
		limitStr := c.Query("limit")
		limit := 100
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}

		// oldest / latest 参数，可选，Slack ts 字符串
		oldest := c.Query("oldest")
		latest := c.Query("latest")

		messages, err := slackService.ListMessages(userID, channelID, limit, oldest, latest)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"messages": messages})
	})
}

// 自动注册到 routes 模块
func init() {
	routes.RegisterModule("slack", RegisterRoutes)
}
