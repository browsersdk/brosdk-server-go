# Browser Open SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/browsersdk/brosdk-server-go.svg)](https://pkg.go.dev/github.com/browsersdk/brosdk-server-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/browsersdk/brosdk-server-go)](https://goreportcard.com/report/github.com/browsersdk/brosdk-server-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

📖 [中文文档](README_zh.md) | English Documentation

A comprehensive Go SDK for interacting with the Browser Open API service, providing seamless integration for browser environment management and user authentication.

## 🚀 Features

- **🔐 Secure Authentication**: Built-in API key authentication with Bearer token support
- **🌐 Flexible Configuration**: Customizable endpoints, timeouts, and HTTP clients
- **🔍 Debug Mode**: Built-in structured logging for request/response debugging
- **⚡ Type-Safe Operations**: Strongly typed request/response structures
- **🔧 Modular Design**: Clean separation of concerns across multiple components
- **🧪 Comprehensive Testing**: Full test coverage with mocking capabilities
- **📦 Easy Integration**: Simple installation and straightforward API usage

## 📦 Installation

```bash
go get github.com/browsersdk/brosdk-server-go
```

## 🎯 Quick Start

### Basic Client Initialization

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
    // Create client with required API key
    client, err := brosdk.NewClient("your-api-key-here")
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }
    
    // Get user signature for authentication
    sigReq := &brosdk.GetUserSigRequest{
        CustomerId: "customer123",
        Duration:   3600, // 1 hour in seconds
    }
    
    sigResp, err := client.GetUserSig(context.Background(), sigReq)
    if err != nil {
        log.Fatal("Failed to get user signature:", err)
    }
    
    fmt.Printf("User Signature: %s\n", sigResp.Data.UserSig)
    fmt.Printf("Expires at: %d\n", sigResp.Data.ExpireTime)
}
```

### Advanced Configuration

```go
// Configure with custom settings
client, err := brosdk.NewClient("your-api-key-here",
    brosdk.WithEndpoint("https://custom.api.browser-open.com"),
    brosdk.WithTimeout(60*time.Second),
    brosdk.WithHTTPClient(&http.Client{
        Timeout: 30 * time.Second,
    }),
)
if err != nil {
    log.Fatal("Failed to create client:", err)
}
```

### Debug Mode

Enable debug mode to log all HTTP requests and responses:

```go
// Enable debug mode (logs to stderr by default)
client, err := brosdk.NewClient("your-api-key-here",
    brosdk.WithDebug(true),
)

// Enable debug mode with custom logger
import "log"

logger := log.New(os.Stdout, "[brosdk] ", log.LstdFlags)
client, err := brosdk.NewClient("your-api-key-here",
    brosdk.WithDebug(true),
    brosdk.WithLogger(logger),
)

// Example debug output:
// [brosdk] 2026/03/20 12:00:00 → POST /api/v2/browser/create  body={"envId":"...","finger":{...}}
// [brosdk] 2026/03/20 12:00:00 ← POST /api/v2/browser/create  elapsed=123ms  http_status=200  body={"code":200,"data":{...}}
```

Debug mode helps you troubleshoot API issues by logging:
- Request method and URL path
- Request body (JSON)
- Response HTTP status code
- Response body (JSON)
- Request duration/elapsed time
- Error details (if any)

## 🧪 Running Tests

The SDK includes comprehensive test coverage. Run tests with:

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -cover

# Run tests with race detection
go test -v -race

# Run specific test
go test -v -run TestClient_GetUserSig
```

## 📁 Project Structure

```
brosdk-server-go/
├── sdk.go          # Core SDK implementation
├── types.go        # Data structures and type definitions
├── sdk_test.go     # Comprehensive unit tests
├── README.md       # Documentation
├── go.mod          # Go module definition
└── go.sum          # Dependency checksums
```

## 🔧 Available Methods

### 🔐 GetUserSig - User Authentication

Retrieve user signature for secure authentication:

```go
req := &brosdk.GetUserSigRequest{
    CustomerId: "customer123",
    Duration:   3600, // seconds
}

resp, err := client.GetUserSig(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Signature: %s\n", resp.UserSig)
fmt.Printf("Expires: %d\n", resp.ExpireTime)
```

### 🌐 EnvCreate - Create Browser Environment

Create a new browser environment configuration:

```go
req := &brosdk.CreateEnv{
    CustomerId:  "customer123",
    EnvName:     "My Browser Environment",
    Serial:      "001",
    Finger: brosdk.Finger{
        System:        "Windows 11",
        Kernel:        "Chrome",
        KernelVersion: "140",
    },
}

resp, err := client.EnvCreate(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created Environment ID: %s\n", resp.EnvId)
```

### 🔄 EnvUpdate - Update Environment (v2)

Update an existing browser environment:

```go
envName := "Updated Environment Name"
finger := brosdk.Finger{
        System:        "Windows",
        Kernel:        "Chrome",
        KernelVersion: "140",
}
req := &brosdk.UpdateEnv{
    EnvId:   "123",
    EnvName: &envName,
    Finger:  &finger,
}

resp, err := client.EnvUpdate(context.Background(), req)
if err != nil {
    log.Fatal(err)
}
```

`EnvUpdate` regenerates fingerprints only when fingerprint settings or generation context actually change. The region parameter remains for backward compatibility only; IP-related fingerprint values use the proxy IP when configured and fall back to the client IP of the create or fingerprint update request.

### EnvUpdateEnvMeta - Update Environment Metadata Only

Update the name, serial, third-party data, topic, or proxy settings without changing the fingerprint:

```go
envName := "New Environment Name"
proxy := "" // clear the proxy by passing a non-nil pointer
topic := "work"
topicConfig := json.RawMessage(`{"color":"blue"}`)

resp, err := client.EnvUpdateEnvMeta(context.Background(), &brosdk.UpdateEnvMeta{
    EnvId:       "123",
    EnvName:     &envName,
    Proxy:       &proxy,
    Topic:       &topic,
    TopicConfig: &topicConfig,
})
```

### 🗑️ EnvDestroy - Delete Environment (v2)

Delete a browser environment:

```go
req := &brosdk.EnvDelReq{
    EnvId: "123",
}

err := client.EnvDestroy(context.Background(), req)
if err != nil {
    log.Fatal(err)
}
```

### 📋 GetEnvPage - List Environments (v2)

Get paginated list of browser environments:

```go
req := &brosdk.GetEnvPageReq{
    ReqPage: brosdk.ReqPage{
        Page:     1,
        PageSize: 20,
    },
    CustomerId: "customer123",
    EnvIds:     []string{"1", "2", "3"}, // Optional filtering
}

resp, err := client.GetEnvPage(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total environments: %d\n", resp.Total)
for _, env := range resp.List {
    fmt.Printf("ID: %s, Name: %s\n", env.EnvId, env.EnvName)
}
```

## ⚙️ Configuration Options

### Custom Endpoint

```go
client, err := brosdk.NewClient("api-key",
    brosdk.WithEndpoint("https://your-custom-endpoint.com"))
```

### Custom Timeout

```go
client, err := brosdk.NewClient("api-key",
    brosdk.WithTimeout(30 * time.Second))
```

### Custom HTTP Client

```go
customClient := &http.Client{
    Timeout: 10 * time.Second,
    // Add custom transport, cookies, etc.
}
client, err := brosdk.NewClient("api-key",
    brosdk.WithHTTPClient(customClient))
```

### Debug Mode

```go
// Enable debug logging
client, err := brosdk.NewClient("api-key",
    brosdk.WithDebug(true),
)

// Custom logger for debug output
logger := log.New(os.Stdout, "[brosdk] ", log.LstdFlags)
client, err := brosdk.NewClient("api-key",
    brosdk.WithDebug(true),
    brosdk.WithLogger(logger),
)
```

## 🛡️ Error Handling

The SDK provides comprehensive error handling:

```go
client, err := brosdk.NewClient("") // Empty API key
if err != nil {
    // Handle validation error
    fmt.Printf("Validation error: %v\n", err)
    return
}

// API call errors
resp, err := client.GetUserSig(ctx, req)
if err != nil {
    // Handle API errors
    if strings.Contains(err.Error(), "status:") {
        fmt.Printf("API returned error status\n")
    } else if strings.Contains(err.Error(), "request failed") {
        fmt.Printf("Network error occurred\n")
    }
    return
}
```

## 🌐 API Endpoints

### Version 1 Endpoints
- `POST /api/usersig` - User signature generation
- `POST /api/env` - Environment creation

### Version 2 Endpoints
- `POST /api/v2/browser/create` - Environment creation
- `POST /api/v2/browser/update` - Environment updates
- `POST /api/v2/browser/updateEnv` - Environment metadata updates
- `POST /api/v2/browser/destroy` - Environment deletion
- `POST /api/v2/browser/page` - Environment listing
- `GET /api/v2/browser/getUiFingerList` - Fingerprint configuration options
- `GET /api/v2/browser/platformList` - Kernel platforms
- `GET /api/v2/browser/archList` - Kernel architectures
- `GET /api/v2/browser/kernelIdList` - Kernel identifiers
- `POST /api/v2/browser/kernelList` - Available kernels
- `POST /api/v2/browser/getGlobalFinger` - Global fingerprint settings
- `POST /api/v2/browser/setGlobalFinger` - Update global fingerprint settings

## 🔒 Security Features

- **Bearer Token Authentication**: Automatic Authorization header management
- **HTTPS Support**: Secure communication by default
- **Input Validation**: Client-side validation of required parameters
- **Context Support**: Request cancellation and timeout handling

## 📊 Response Structure

All API responses follow a consistent structure:

```go
type Response struct {
    Code  int         `json:"code"`   // 0 for success, non-zero for errors
    Data  interface{} `json:"data"`   // Response data varies by endpoint
    Msg   string      `json:"msg"`    // Human-readable message
    ReqId string      `json:"reqId"`  // Request identifier for debugging
}
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
# Clone the repository
git clone https://github.com/browsersdk/brosdk-server-go.git
cd brosdk-server-go

# Run tests
go test -v

# Check code coverage
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run linting
golangci-lint run
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For issues, questions, or contributions:
- 🐛 [Report Issues](https://github.com/browsersdk/brosdk-server-go/issues)
- 💬 [GitHub Discussions](https://github.com/browsersdk/brosdk-server-go/discussions)
- 📧 Email: support@browsersdk.com

## 📈 Changelog

### v1.1.0
- ✨ Add debug mode with structured logging
- ✨ Add `WithDebug()` option to enable/disable debug logging
- ✨ Add `WithLogger()` option for custom log output
- 🐛 Fix incorrect error messages in `GetUiFingerList`
- 🐛 Fix missing JSON tag on `CONAB.Must` field
- 🐛 Remove leftover debug `fmt.Printf` statements
- 🧪 Add comprehensive debug mode tests

### v1.0.0
- Initial release
- Core SDK functionality
- Comprehensive test coverage
- Full API method implementations
- Documentation and examples
