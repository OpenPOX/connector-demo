package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default() // 创建默认的 Gin 引擎

	// 定义路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})

	// 启动服务器，默认监听 8080 端口
	r.Run(":6767")
}
