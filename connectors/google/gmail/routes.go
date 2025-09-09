package gmail

import "github.com/gin-gonic/gin"

var gmailService *GmailService

func SetGmailService(gs *GmailService) {
	gmailService = gs
}

func RegisterRoutes(rg *gin.RouterGroup) {
	gmailGroup := rg.Group("/gmail")

	gmailGroup.GET("/inbox", func(c *gin.Context) {
		userID := c.Query("user_id")
		messages, _ := gmailService.GetInboxMessages(userID, 10)
		c.JSON(200, gin.H{"messages": messages})
	})
}
