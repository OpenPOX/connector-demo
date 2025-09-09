package main

import (
	"fmt"
	"log"
	"net/http"

	"connector-demo/auth"
	"connector-demo/config"
	"connector-demo/connectors/google"
	"connector-demo/connectors/slack"
	"connector-demo/routes"
	"connector-demo/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/markbates/goth/providers/google"
	_ "github.com/markbates/goth/providers/slack"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()
	// 初始化OAuth2提供者
	if err := auth.SetupProviders(cfg); err != nil {
		log.Fatalf("初始化OAuth2提供者失败: %v", err)
	}

	// 初始化token管理器
	tokenManager := utils.NewTokenManager()
	// 注入测试 token
	utils.InjectTestTokens(tokenManager)

	// 初始化认证处理器
	authHandler := auth.NewAuthHandler(tokenManager)
	// 创建Google
	googleService := google.NewGoogleService(tokenManager)
	google.SetGoogleService(googleService)
	//Slack连接器
	slackService := slack.NewSlackService(tokenManager)
	slack.SetSlackService(slackService)

	// 创建Gin路由
	r := gin.Default()
	// 基础路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "连接器演示服务已启动",
			"version": "1.0.0",
		})
	})

	routes.RegisterAllModules(r)

	// OAuth2认证路由组
	oauth := r.Group("/auth")
	{
		oauth.GET("/:provider", authHandler.Connect)           // 开始认证
		oauth.GET("/:provider/callback", authHandler.Callback) // 回调处理
	}

	// Token管理路由组
	tokens := r.Group("/tokens")
	{
		tokens.GET("/list", authHandler.GetTokens)                     // 获取用户token列表
		tokens.DELETE("/disconnect/:provider", authHandler.Disconnect) // 断开连接
		tokens.GET("/refresh/:provider", func(c *gin.Context) {        // 刷新 token
			// 假设用户ID从上下文获取（例如通过中间件 auth 解析 JWT）
			userID := c.Query("user_id")
			if userID == "" {
				c.JSON(401, gin.H{"error": "用户未登录"})
				return
			}

			provider := c.Param("provider")
			if provider == "" {
				c.JSON(400, gin.H{"error": "provider不能为空"})
				return
			}

			newToken, err := tokenManager.RefreshToken(userID, provider)
			if err != nil {
				c.JSON(500, gin.H{"error": "刷新token失败", "detail": err.Error()})
				return
			}

			c.JSON(200, gin.H{
				"access_token":  newToken.AccessToken,
				"refresh_token": newToken.RefreshToken,
				"expiry":        newToken.Expiry,
				"provider":      newToken.Provider,
			})
		})
	}

	// 调试路由
	r.GET("/debug/tokens", func(c *gin.Context) {
		tokenManager.PrintAllTokens()
		c.JSON(200, gin.H{"message": "token信息已打印到控制台"})
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"timestamp": fmt.Sprintf("%d", utils.GetCurrentTimestamp()),
		})
	})

	// 启动服务器
	port := ":6767"
	log.Printf("🚀 连接器演示服务启动成功！")
	log.Printf("🔗 认证地址: http://localhost%s/auth/{google|slack}", port)
	log.Printf("🧪 测试地址: http://localhost%s/api", port)

	log.Fatal(http.ListenAndServe(port, r))
}
