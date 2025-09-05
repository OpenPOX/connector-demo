# OAuth2连接器服务

这是一个基于Gin框架的Go服务，用于统一接入第三方平台的OAuth2认证，支持Google和Slack平台。

## 功能特性

- ✅ Google OAuth2认证（Gmail/Drive）
- ✅ Slack OAuth2认证
- ✅ 内存token管理（开发环境）
- ✅ 可扩展的provider架构
- ✅ RESTful API接口
- ✅ 完整的API文档和测试接口

## 项目结构

```
connector-demo/
├── main.go                 # 主程序入口
├── go.mod                  # 依赖管理
├── config/
│   └── config.go          # 配置管理
├── auth/
│   ├── auth.go            # OAuth2路由处理
│   └── provider.go        # goth提供者配置
├── connectors/
│   ├── google_connector.go    # Google API连接器
│   └── slack_connector.go     # Slack API连接器
├── utils/
│   ├── token_manager.go   # token管理器
│   └── oauth_client.go    # OAuth2客户端工具
└── README.md              # 项目文档
```

## 快速开始

### 1. 环境配置

设置必要的环境变量：

```bash
# Google OAuth2配置
export GOOGLE_CLIENT_ID="your_google_client_id"
export GOOGLE_CLIENT_SECRET="your_google_client_secret"

# Slack OAuth2配置
export SLACK_CLIENT_ID="your_slack_client_id"
export SLACK_CLIENT_SECRET="your_slack_client_secret"

# 可选配置
export REDIRECT_URL="http://localhost:8080"  # 回调URL
export PORT="8080"                             # 服务端口
```

### 2. 安装依赖

```bash
go mod tidy
go mod download
```

### 3. 启动服务

```bash
go run main.go
```

## API文档

### 基础接口

- `GET /` - 服务状态和信息
- `GET /health` - 健康检查

### OAuth2认证

- `GET /auth/connect/:platform` - 开始OAuth2流程
  - 平台: `google`, `slack`
- `GET /auth/callback/:platform` - OAuth2回调处理

### Token管理

- `GET /tokens?user_id={user_id}` - 获取用户token
- `DELETE /tokens/disconnect/:platform?user_id={user_id}` - 断开连接

### API测试

#### Google API
- `GET /api/google/test?user_id={user_id}` - 测试连接
- `GET /api/google/gmail?user_id={user_id}` - 获取Gmail邮件列表
- `GET /api/google/drive?user_id={user_id}` - 获取Drive文件列表

#### Slack API
- `GET /api/slack/test?user_id={user_id}` - 测试连接
- `GET /api/slack/channels?user_id={user_id}` - 获取频道列表
- `GET /api/slack/messages?user_id={user_id}&channel_id={channel_id}` - 获取消息列表

### 调试接口
- `GET /debug/tokens` - 查看所有token（调试用）

## 使用示例

### 1. 连接Google

1. 访问: `http://localhost:8080/auth/connect/google`
2. 授权后，系统会在控制台打印access_token和refresh_token
3. 使用返回的用户ID测试API: `http://localhost:8080/api/google/test?user_id=用户邮箱`

### 2. 连接Slack

1. 访问: `http://localhost:8080/auth/connect/slack`
2. 授权后，系统会在控制台打印access_token和refresh_token
3. 使用返回的用户ID测试API: `http://localhost:8080/api/slack/test?user_id=用户ID`

## 开发指南

### 添加新的Provider

1. 在`auth/provider.go`中添加新的provider配置
2. 在`connectors/`目录下创建新的连接器文件
3. 在`main.go`的`setupRoutes`中添加对应的API路由

### 示例：添加Notion支持

```go
// 在auth/provider.go中添加
providers = append(providers, notion.New(
    cfg.NotionClientID,
    cfg.NotionClientSecret,
    fmt.Sprintf("%s/auth/callback/notion", cfg.RedirectURL),
))
```

### Token存储扩展

当前使用内存存储token，生产环境建议：

1. 实现数据库存储（Redis/PostgreSQL）
2. 添加token刷新机制
3. 实现用户会话管理

## 注意事项

- 当前token存储在内存中，重启服务会丢失
- 请确保回调URL与OAuth应用配置一致
- 开发环境建议使用`http://localhost:8080`
- 生产环境请使用HTTPS

## 环境变量模板

创建`.env`文件：

```bash
# Google OAuth2
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# Slack OAuth2  
SLACK_CLIENT_ID=your_slack_client_id
SLACK_CLIENT_SECRET=your_slack_client_secret

# 服务配置
REDIRECT_URL=http://localhost:8080
PORT=8080
GIN_MODE=debug
```

## 许可证

MIT License