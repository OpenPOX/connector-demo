package slack

// import (
// 	"connector-demo/routes"

// 	"github.com/gin-gonic/gin"
// )

// var slackService *SlackService

// func SetSlackService(ss *SlackService) {
// 	slackService = ss
// }

// func RegisterRoutes(rg *gin.RouterGroup) {
// 	slackGroup := rg.Group("/slack")

// 	slackGroup.GET("/channels", func(c *gin.Context) {
// 		userID := c.Query("user_id")
// 		channels, _ := slackService.ListChannels(userID)
// 		c.JSON(200, gin.H{"channels": channels})
// 	})
// }

// func init() {
// 	routes.RegisterModule("slack", RegisterRoutes)
// }
