package main

import (
	"fmt"
	"log"
	"net/http"

	"connector-demo/auth"
	"connector-demo/config"
	"connector-demo/connectors"
	"connector-demo/connectors/google"
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
	// 根据环境变量决定是否注入测试token
	if config.GetEnv("ENABLE_TEST_TOKENS", "false") == "true" {
		log.Println("检测到测试token配置，正在注入...")
		// 注入Google测试token
		if accessToken := config.GetEnv("TEST_TOKEN_GOOGLE_ACCESS", ""); accessToken != "" {
			googleToken := &utils.TokenInfo{
				AccessToken:  accessToken,
				RefreshToken: config.GetEnv("TEST_TOKEN_GOOGLE_REFRESH", ""),
				Provider:     "google",
				TokenType:    "Bearer",
			}
			tokenManager.SaveToken("1", "google", googleToken)
			log.Println("已注入Google测试token")
		}

		// 注入Slack测试token
		if accessToken := config.GetEnv("TEST_TOKEN_SLACK_ACCESS", ""); accessToken != "" {
			slackToken := &utils.TokenInfo{
				AccessToken:  accessToken,
				RefreshToken: config.GetEnv("TEST_TOKEN_SLACK_REFRESH", ""),
				Provider:     "slack",
				TokenType:    "Bearer",
			}
			tokenManager.SaveToken("1", "slack", slackToken)
			log.Println("已注入Slack测试token")
		}
	}

	// 初始化认证处理器
	authHandler := auth.NewAuthHandler(tokenManager)
	// 创建Google和Slack连接器
	googleService := google.NewGoogleService(tokenManager)
	slackConnector := connectors.NewSlackConnector(tokenManager)

	// 创建Gin路由
	r := gin.Default()
	// 基础路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "连接器演示服务已启动",
			"version": "1.0.0",
			"docs":    "/docs",
		})
	})

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
	}

	// API测试路由组
	api := r.Group("/api")
	{
		// Google API测试
		google := api.Group("/google")
		{
			google.GET("/user-info", func(c *gin.Context) {
				userID := c.Query("user_id")
				if userID == "" {
					c.JSON(400, gin.H{"error": "缺少user_id参数"})
					return
				}

				userInfo, err := googleService.GetUserInfo(userID)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, gin.H{"user_info": userInfo})
			})

			google.GET("/gmail", func(c *gin.Context) {
				userID := c.Query("user_id")
				if userID == "" {
					c.JSON(400, gin.H{"error": "缺少user_id参数"})
					return
				}

				limit := int64(10)
				messages, err := googleService.Gmail.GetInboxMessages(userID, limit)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, gin.H{"messages": messages})
			})

			google.GET("/drive", func(c *gin.Context) {
				userID := c.Query("user_id")
				if userID == "" {
					c.JSON(400, gin.H{"error": "缺少user_id参数"})
					return
				}

				limit := int64(10)
				files, err := googleService.Drive.GetFiles(userID, limit)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, gin.H{"files": files})
			})

			google.GET("/test", func(c *gin.Context) {
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

		// Slack API测试
		slack := api.Group("/slack")
		{
			slack.GET("/user-info", func(c *gin.Context) {
				userID := c.Query("user_id")
				if userID == "" {
					c.JSON(400, gin.H{"error": "缺少user_id参数"})
					return
				}

				userInfo, err := slackConnector.GetUserInfo(userID)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, gin.H{"user_info": userInfo})
			})

			slack.GET("/channels", func(c *gin.Context) {
				userID := c.Query("user_id")
				if userID == "" {
					c.JSON(400, gin.H{"error": "缺少user_id参数"})
					return
				}

				channels, err := slackConnector.ListChannels(userID)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, gin.H{"channels": channels})
			})

			slack.GET("/test", func(c *gin.Context) {
				userID := c.Query("user_id")
				if userID == "" {
					c.JSON(400, gin.H{"error": "缺少user_id参数"})
					return
				}

				if err := slackConnector.TestConnection(userID); err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, gin.H{"message": "Slack连接测试成功"})
			})
		}
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
	log.Printf("📖 文档地址: http://localhost%s/docs", port)
	log.Printf("🔗 认证地址: http://localhost%s/auth/{google|slack}", port)
	log.Printf("🧪 测试地址: http://localhost%s/api", port)

	log.Fatal(http.ListenAndServe(port, r))
}
