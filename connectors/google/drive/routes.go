package drive

import "github.com/gin-gonic/gin"

var driveService *DriveService

func SetDriveService(ds *DriveService) {
	driveService = ds
}

func RegisterRoutes(rg *gin.RouterGroup) {
	driveGroup := rg.Group("/drive")

	driveGroup.GET("/files", func(c *gin.Context) {
		userID := c.Query("user_id")
		files, _ := driveService.GetFiles(userID, 10)
		c.JSON(200, gin.H{"files": files})
	})
}
