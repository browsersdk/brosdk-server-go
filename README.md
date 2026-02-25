# Browser Open SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/browsersdk/brosdk-server-go.svg)](https://pkg.go.dev/github.com/browsersdk/brosdk-server-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/browsersdk/brosdk-server-go)](https://goreportcard.com/report/github.com/browsersdk/brosdk-server-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

📖 [中文文档](README_zh.md) | English Documentation

A comprehensive Go SDK for interacting with the Browser Open API service, providing seamless integration for browser environment management and user authentication.

## 🚀 Features

- **🔐 Secure Authentication**: Built-in API key authentication with Bearer token support
- **🌐 Flexible Configuration**: Customizable endpoints, timeouts, and HTTP clients
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

fmt.Printf("Signature: %s\n", resp.Data.UserSig)
fmt.Printf("Expires: %d\n", resp.Data.ExpireTime)
```

### 🌐 EnvCreate - Create Browser Environment

Create a new browser environment configuration:

```go
req := &brosdk.EnvRequest{
    CustomerId:      "customer123",
    EnvName:         "My Browser Environment",
    UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    System:          "Windows 10",
    Kernel:          "Chrome",
    KernelVersion:   "120.0.0.0",
    EnableCookie:    1,
    Enablenotice:    1,
    Enableopen:      1,
    Enablepic:       1,
    IgnoreCookieErr: 0,
    // Add other required fields...
}

resp, err := client.EnvCreate(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created Environment ID: %d\n", resp.Data.EnvId)
```

### 🔄 EnvUpdate - Update Environment (v2)

Update an existing browser environment:

```go
req := &brosdk.EnvRequest{
    EnvId:      123,
    CustomerId: "customer123",
    EnvName:    "Updated Environment Name",
    // Update other fields as needed...
}

resp, err := client.EnvUpdate(context.Background(), req)
if err != nil {
    log.Fatal(err)
}
```

### 🗑️ EnvDestroy - Delete Environment (v2)

Delete a browser environment:

```go
req := &brosdk.EnvReq{
    EnvId: 123,
}

resp, err := client.EnvDestroy(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Deletion result: %s\n", resp.Msg)
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
    EnvIds:     []uint64{1, 2, 3}, // Optional filtering
}

resp, err := client.GetEnvPage(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total environments: %d\n", resp.Total)
for _, env := range resp.Data {
    fmt.Printf("ID: %d, Name: %s, Created: %s\n", 
        env.EnvId, env.EnvName, env.CreatedAt)
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
- `POST /api/v2/browser/update` - Environment updates
- `POST /api/v2/browser/destroy` - Environment deletion
- `POST /api/v2/browser/page` - Environment listing

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

### v1.0.0
- Initial release
- Core SDK functionality
- Comprehensive test coverage
- Full API method implementations
- Documentation and examples