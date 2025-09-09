# Google连接器分层架构

## 目录结构

```
connectors/
└── google/
    ├── connector.go          # GoogleConnector统一入口，管理token
    ├── service.go            # GoogleService聚合各子模块
    ├── gmail/
    │   ├── connector.go      # Gmail API 封装
    │   └── service.go        # Gmail 数据处理接口
    └── drive/
        ├── connector.go      # Google Drive API 封装
        └── service.go        # Drive 数据处理接口
```

## 架构说明

### 1. GoogleConnector (统一入口)
- **文件**: `connectors/google/connector.go`
- **职责**: 
  - 管理用户token
  - 提供统一的用户信息获取接口
  - 协调各子模块服务

### 2. GoogleService (聚合服务)
- **文件**: `connectors/google/service.go`
- **职责**:
  - 聚合Gmail和Drive服务
  - 提供高层次的业务接口
  - 简化外部调用

### 3. Gmail模块
- **connector.go**: Gmail API的原始封装
- **service.go**: Gmail业务逻辑处理
- **功能**:
  - 获取邮件列表
  - 获取邮件详情
  - 管理邮件标签

### 4. Drive模块
- **connector.go**: Drive API的原始封装
- **service.go**: Drive业务逻辑处理
- **功能**:
  - 获取文件列表
  - 获取文件详情
  - 文件管理操作

## 使用示例

```go
// 创建Google服务
googleService := google.NewGoogleService(tokenManager)

// 获取用户信息
userInfo, err := googleService.GetUserInfo("user123")

// 获取Gmail邮件
messages, err := googleService.Gmail.GetInboxMessages("user123", 10)

// 获取Drive文件
files, err := googleService.Drive.GetFiles("user123", 10)
```

## 优势

1. **解耦**: 各模块独立，职责清晰
2. **可扩展**: 易于添加新的Google服务模块
3. **可测试**: 各层可独立测试
4. **可维护**: 代码结构清晰，易于维护
5. **复用**: 各模块可在不同场景下复用