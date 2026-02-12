# Browser Open SDK 中文文档

[![Go Reference](https://pkg.go.dev/badge/github.com/browsersdk/brosdk-server-go.svg)](https://pkg.go.dev/github.com/browsersdk/brosdk-server-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/browsersdk/brosdk-server-go)](https://goreportcard.com/report/github.com/browsersdk/brosdk-server-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Browser Open SDK 是一个功能完整的 Go 语言 SDK，用于与 Browser Open API 服务进行交互，提供浏览器环境管理和用户认证的无缝集成。

## 🚀 主要特性

- **🔐 安全认证**: 内置 API 密钥认证，支持 Bearer token
- **🌐 灵活配置**: 可自定义端点、超时时间和 HTTP 客户端
- **⚡ 类型安全**: 强类型的请求/响应结构
- **🔧 模块化设计**: 组件间职责清晰分离
- **🧪 全面测试**: 完整的测试覆盖和模拟能力
- **📦 易于集成**: 安装简单，API 使用直观

## 📦 安装

```bash
go get github.com/browsersdk/brosdk-server-go
```

## 🎯 快速开始

### 基础客户端初始化

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/browsersdk/brosdk-server-go"
)

func main() {
    // 使用 API 密钥创建客户端
    client, err := brosdk.NewClient("your-api-key-here")
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }
    
    // 获取用户签名用于认证
    sigReq := &brosdk.GetUserSigRequest{
        CustomerId: "customer123",
        Duration:   3600, // 1小时（秒）
    }
    
    sigResp, err := client.GetUserSig(context.Background(), sigReq)
    if err != nil {
        log.Fatal("获取用户签名失败:", err)
    }
    
    fmt.Printf("用户签名: %s\n", sigResp.Data.UserSig)
    fmt.Printf("过期时间: %d\n", sigResp.Data.ExpireTime)
}
```

### 高级配置

```go
// 使用自定义配置
client, err := brosdk.NewClient("your-api-key-here",
    brosdk.WithEndpoint("https://custom.api.browser-open.com"),
    brosdk.WithTimeout(60*time.Second),
    brosdk.WithHTTPClient(&http.Client{
        Timeout: 30 * time.Second,
    }),
)
if err != nil {
    log.Fatal("创建客户端失败:", err)
}
```

## 🧪 运行测试

SDK 包含全面的测试覆盖。运行测试：

```bash
# 运行所有测试
go test -v

# 运行带覆盖率的测试
go test -v -cover

# 运行带竞态检测的测试
go test -v -race

# 运行特定测试
go test -v -run TestClient_GetUserSig
```

## 📁 项目结构

```
brosdk-server-go/
├── sdk.go          # 核心 SDK 实现
├── types.go        # 数据结构和类型定义
├── sdk_test.go     # 全面的单元测试
├── README.md       # 英文文档
├── README_zh.md    # 中文文档
├── go.mod          # Go 模块定义
└── go.sum          # 依赖校验和
```

## 🔧 可用方法

### 🔐 GetUserSig - 用户认证

获取用户签名用于安全认证：

```go
req := &brosdk.GetUserSigRequest{
    CustomerId: "customer123",
    Duration:   3600, // 秒
}

resp, err := client.GetUserSig(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("签名: %s\n", resp.Data.UserSig)
fmt.Printf("过期时间: %d\n", resp.Data.ExpireTime)
```

### 🌐 EnvCreate - 创建浏览器环境

创建新的浏览器环境配置：

```go
req := &brosdk.EnvCreateRequest{
    CustomerId:      "customer123",
    EnvName:         "我的浏览器环境",
    UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    System:          "Windows 10",
    Kernel:          "Chrome",
    KernelVersion:   "120.0.0.0",
    EnableCookie:    1,
    Enablenotice:    1,
    Enableopen:      1,
    Enablepic:       1,
    IgnoreCookieErr: 0,
    // 添加其他必要字段...
}

resp, err := client.EnvCreate(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("创建的环境ID: %d\n", resp.Data.EnvId)
```

### 🔄 EnvUpdate - 更新环境 (v2)

更新现有的浏览器环境：

```go
req := &brosdk.EnvCreateRequest{
    EnvId:      123,
    CustomerId: "customer123",
    EnvName:    "更新的环境名称",
    // 根据需要更新其他字段...
}

resp, err := client.EnvUpdate(context.Background(), req)
if err != nil {
    log.Fatal(err)
}
```

### 🗑️ EnvDestroy - 删除环境 (v2)

删除浏览器环境：

```go
req := &brosdk.EnvReq{
    EnvId: 123,
}

resp, err := client.EnvDestroy(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("删除结果: %s\n", resp.Msg)
```

### 📋 GetEnvPage - 列出环境 (v2)

获取分页的浏览器环境列表：

```go
req := &brosdk.GetEnvPageReq{
    ReqPage: brosdk.ReqPage{
        Page:     1,
        PageSize: 20,
    },
    CustomerId: "customer123",
    EnvIds:     []uint64{1, 2, 3}, // 可选过滤
}

resp, err := client.GetEnvPage(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("环境总数: %d\n", resp.Total)
for _, env := range resp.Data {
    fmt.Printf("ID: %d, 名称: %s, 创建时间: %s\n", 
        env.EnvId, env.EnvName, env.CreatedAt)
}
```

## ⚙️ 配置选项

### 自定义端点

```go
client, err := brosdk.NewClient("api-key", 
    brosdk.WithEndpoint("https://your-custom-endpoint.com"))
```

### 自定义超时时间

```go
client, err := brosdk.NewClient("api-key",
    brosdk.WithTimeout(30 * time.Second))
```

### 自定义 HTTP 客户端

```go
customClient := &http.Client{
    Timeout: 10 * time.Second,
    // 添加自定义传输、cookie 等
}
client, err := brosdk.NewClient("api-key",
    brosdk.WithHTTPClient(customClient))
```

## 🛡️ 错误处理

SDK 提供全面的错误处理：

```go
client, err := brosdk.NewClient("") // 空 API 密钥
if err != nil {
    // 处理验证错误
    fmt.Printf("验证错误: %v\n", err)
    return
}

// API 调用错误
resp, err := client.GetUserSig(ctx, req)
if err != nil {
    // 处理 API 错误
    if strings.Contains(err.Error(), "status:") {
        fmt.Printf("API 返回错误状态\n")
    } else if strings.Contains(err.Error(), "request failed") {
        fmt.Printf("发生网络错误\n")
    }
    return
}
```

## 🌐 API 端点

### 版本 1 端点
- `POST /api/usersig` - 用户签名生成
- `POST /api/env` - 环境创建

### 版本 2 端点
- `POST /api/v2/browser/update` - 环境更新
- `POST /api/v2/browser/destroy` - 环境删除
- `POST /api/v2/browser/page` - 环境列表

## 🔒 安全特性

- **Bearer Token 认证**: 自动管理 Authorization 头
- **HTTPS 支持**: 默认安全通信
- **输入验证**: 客户端验证必要参数
- **上下文支持**: 请求取消和超时处理

## 📊 响应结构

所有 API 响应遵循一致的结构：

```go
type Response struct {
    Code  int         `json:"code"`   // 0 表示成功，非零表示错误
    Data  interface{} `json:"data"`   // 响应数据因端点而异
    Msg   string      `json:"msg"`    // 人类可读的消息
    ReqId string      `json:"reqId"`  // 用于调试的请求标识符
}
```

## 🤝 贡献

1. Fork 仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m '添加令人惊叹的特性'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 发起 Pull Request

### 开发设置

```bash
# 克隆仓库
git clone https://github.com/browsersdk/brosdk-server-go.git
cd brosdk-server-go

# 运行测试
go test -v

# 检查代码覆盖率
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# 运行代码检查
golangci-lint run
```

## 📄 许可证

本项目采用 MIT 许可证 - 详情请见 [LICENSE](LICENSE) 文件。

## 🆘 支持

如有问题、疑问或贡献：
- 🐛 [报告问题](https://github.com/browsersdk/brosdk-server-go/issues)
- 💬 [GitHub 讨论](https://github.com/browsersdk/brosdk-server-go/discussions)
- 📧 邮箱: support@browsersdk.com

## 📈 更新日志

### v1.0.0
- 初始发布
- 核心 SDK 功能
- 全面的测试覆盖
- 完整的 API 方法实现
- 文档和示例