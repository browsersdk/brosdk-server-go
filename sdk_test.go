package brosdk

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

// MockRoundTripper implements http.RoundTripper for testing
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		opts    []ClientOption
		wantErr bool
	}{
		{
			name:    "valid api key",
			apiKey:  "test-key",
			wantErr: false,
		},
		{
			name:    "empty api key",
			apiKey:  "",
			wantErr: true,
		},
		{
			name:    "with custom endpoint",
			apiKey:  "test-key",
			opts:    []ClientOption{WithEndpoint("https://custom.example.com")},
			wantErr: false,
		},
		{
			name:    "with custom timeout",
			apiKey:  "test-key",
			opts:    []ClientOption{WithTimeout(10 * time.Second)},
			wantErr: false,
		},
		{
			name:    "with custom http client",
			apiKey:  "test-key",
			opts:    []ClientOption{WithHTTPClient(&http.Client{Timeout: 5 * time.Second})},
			wantErr: false,
		},
		{
			name:    "with debug enabled",
			apiKey:  "test-key",
			opts:    []ClientOption{WithDebug(true)},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.apiKey, tt.opts...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewClient() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}

			if client == nil {
				t.Fatal("NewClient() returned nil client")
			}

			if client.ApiKey != tt.apiKey {
				t.Errorf("NewClient() ApiKey = %v, want %v", client.ApiKey, tt.apiKey)
			}

			if client.httpClient == nil {
				t.Error("NewClient() httpClient is nil")
			}
		})
	}
}

func TestWithEndpoint(t *testing.T) {
	client := &Client{}
	option := WithEndpoint("https://test.example.com")
	option(client)

	if client.Endpoint != "https://test.example.com" {
		t.Errorf("WithEndpoint() = %v, want %v", client.Endpoint, "https://test.example.com")
	}

	// Test with empty endpoint (should not change)
	oldEndpoint := client.Endpoint
	option = WithEndpoint("")
	option(client)
	if client.Endpoint != oldEndpoint {
		t.Errorf("WithEndpoint() with empty string changed endpoint from %v to %v", oldEndpoint, client.Endpoint)
	}
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{Timeout: 10 * time.Second}
	client := &Client{}
	option := WithHTTPClient(customClient)
	option(client)

	if client.httpClient != customClient {
		t.Error("WithHTTPClient() did not set the custom HTTP client")
	}

	// Test with nil client (should not change)
	option = WithHTTPClient(nil)
	option(client)
	if client.httpClient != customClient {
		t.Error("WithHTTPClient() with nil should not change existing client")
	}
}

func TestWithTimeout(t *testing.T) {
	client := &Client{}
	timeout := 5 * time.Second
	option := WithTimeout(timeout)
	option(client)

	if client.httpClient == nil {
		t.Fatal("WithTimeout() httpClient is nil")
	}

	if client.httpClient.Timeout != timeout {
		t.Errorf("WithTimeout() = %v, want %v", client.httpClient.Timeout, timeout)
	}

	// Test updating existing client
	newTimeout := 10 * time.Second
	option = WithTimeout(newTimeout)
	option(client)

	if client.httpClient.Timeout != newTimeout {
		t.Errorf("WithTimeout() update = %v, want %v", client.httpClient.Timeout, newTimeout)
	}
}

func TestClient_newRequest(t *testing.T) {
	client, err := NewClient("test-key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	body := &GetUserSigRequest{CustomerId: "test", Duration: 3600}

	req, rawBody, err := client.newRequest(ctx, "POST", "/test", body)
	if err != nil {
		t.Fatalf("newRequest() error = %v", err)
	}

	// Test request properties
	if req.Method != "POST" {
		t.Errorf("newRequest() Method = %v, want POST", req.Method)
	}

	expectedURL := DefaultEndpoint + "/test"
	if req.URL.String() != expectedURL {
		t.Errorf("newRequest() URL = %v, want %v", req.URL.String(), expectedURL)
	}

	// Test headers
	authHeader := req.Header.Get("Authorization")
	if authHeader != "Bearer test-key" {
		t.Errorf("newRequest() Authorization header = %v, want Bearer test-key", authHeader)
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("newRequest() Content-Type = %v, want application/json", contentType)
	}

	userAgent := req.Header.Get("User-Agent")
	if userAgent != "browser-open-sdk/1.0" {
		t.Errorf("newRequest() User-Agent = %v, want browser-open-sdk/1.0", userAgent)
	}

	// Test context
	if req.Context() != ctx {
		t.Error("newRequest() context not set correctly")
	}

	// Verify rawBody is valid JSON
	var parsedFromRaw GetUserSigRequest
	if err := json.Unmarshal(rawBody, &parsedFromRaw); err != nil {
		t.Fatalf("newRequest() rawBody is not valid JSON: %v", err)
	}
	if parsedFromRaw.CustomerId != "test" || parsedFromRaw.Duration != 3600 {
		t.Errorf("newRequest() rawBody = %+v, want CustomerId=test, Duration=3600", parsedFromRaw)
	}

	// Test body content from request
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("Failed to read request body: %v", err)
	}

	var parsedBody GetUserSigRequest
	if err := json.Unmarshal(bodyBytes, &parsedBody); err != nil {
		t.Fatalf("Failed to unmarshal request body: %v", err)
	}

	if parsedBody.CustomerId != "test" || parsedBody.Duration != 3600 {
		t.Errorf("newRequest() body = %+v, want CustomerId=test, Duration=3600", parsedBody)
	}
}

func TestClient_newRequest_NoBody(t *testing.T) {
	client, err := NewClient("test-key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, rawBody, err := client.newRequest(context.Background(), "GET", "/test", nil)
	if err != nil {
		t.Fatalf("newRequest() error = %v", err)
	}

	if req.Body != nil {
		t.Error("newRequest() with nil body should have nil Body")
	}

	if rawBody != nil {
		t.Error("newRequest() with nil body should return nil rawBody")
	}
}

func TestClient_GetUserSig_Success(t *testing.T) {
	// Mock HTTP client
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request path
			if req.URL.Path != "/api/v2/browser/getUserSig" {
				t.Errorf("Expected path /api/v2/browser/getUserSig, got %s", req.URL.Path)
			}

			// Return mock response
			responseBody := `{
				"code": 200,
				"data": {
					"expireTime": 1234567890,
					"userSig": "test-sig"
				},
				"msg": "success",
				"reqId": "test-req-id"
			}`

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, err := NewClient("test-key", WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &GetUserSigRequest{
		CustomerId: "test-customer",
		Duration:   3600,
	}

	resp, err := client.GetUserSig(context.Background(), req)
	if err != nil {
		t.Fatalf("GetUserSig() error = %v", err)
	}

	if resp.UserSig != "test-sig" {
		t.Errorf("GetUserSig() userSig = %v, want test-sig", resp.UserSig)
	}

	if resp.ExpireTime != 1234567890 {
		t.Errorf("GetUserSig() expireTime = %v, want 1234567890", resp.ExpireTime)
	}
}

func TestClient_GetUserSig_HTTPError(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return nil, &mockHTTPError{message: "connection failed"}
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, err := NewClient("test-key", WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &GetUserSigRequest{CustomerId: "test", Duration: 3600}
	_, err = client.GetUserSig(context.Background(), req)

	if err == nil {
		t.Error("GetUserSig() should return error for HTTP failure")
	}

	if !strings.Contains(err.Error(), "request failed") {
		t.Errorf("GetUserSig() error = %v, want error containing 'request failed'", err)
	}
}

func TestClient_GetUserSig_NonOKStatus(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(strings.NewReader("Bad Request")),
				Header:     make(http.Header),
			}, nil
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, err := NewClient("test-key", WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &GetUserSigRequest{CustomerId: "test", Duration: 3600}
	_, err = client.GetUserSig(context.Background(), req)

	if err == nil {
		t.Error("GetUserSig() should return error for non-200 status")
	}

	if !strings.Contains(err.Error(), "non-OK HTTP status") {
		t.Errorf("GetUserSig() error = %v, want error containing 'non-OK HTTP status'", err)
	}
}

func TestClient_GetUserSig_InvalidJSON(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("invalid json")),
				Header:     make(http.Header),
			}, nil
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, err := NewClient("test-key", WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &GetUserSigRequest{CustomerId: "test", Duration: 3600}
	_, err = client.GetUserSig(context.Background(), req)

	if err == nil {
		t.Error("GetUserSig() should return error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "failed to decode") {
		t.Errorf("GetUserSig() error = %v, want error containing 'failed to decode'", err)
	}
}

func TestClient_EnvCreate_Success(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v2/browser/create" {
				t.Errorf("Expected path /api/v2/browser/create, got %s", req.URL.Path)
			}

			responseBody := `{
				"code": 200,
				"data": {
					"envId": "123",
					"envName": "Test Environment",
					"customerId": "test-customer"
				},
				"msg": "success",
				"reqId": "test-req-id"
			}`

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, err := NewClient("test-key", WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &EnvInfo{
		CustomerId: "test-customer",
		EnvName:    "Test Environment",
	}

	resp, err := client.EnvCreate(context.Background(), req)
	if err != nil {
		t.Fatalf("EnvCreate() error = %v", err)
	}

	if resp.EnvId != "123" {
		t.Errorf("EnvCreate() envId = %v, want 123", resp.EnvId)
	}

	if resp.EnvName != "Test Environment" {
		t.Errorf("EnvCreate() envName = %v, want Test Environment", resp.EnvName)
	}
}

func TestClient_EnvUpdate_Success(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v2/browser/update" {
				t.Errorf("Expected path /api/v2/browser/update, got %s", req.URL.Path)
			}

			responseBody := `{
				"code": 200,
				"data": {
					"envId": "456",
					"envName": "Updated Environment"
				},
				"msg": "success",
				"reqId": "test-req-id"
			}`

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, err := NewClient("test-key", WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &EnvInfo{
		EnvId:      "456",
		CustomerId: "test-customer",
		EnvName:    "Updated Environment",
	}

	resp, err := client.EnvUpdate(context.Background(), req)
	if err != nil {
		t.Fatalf("EnvUpdate() error = %v", err)
	}

	if resp.EnvId != "456" {
		t.Errorf("EnvUpdate() envId = %v, want 456", resp.EnvId)
	}

}

func TestClient_EnvDestroy_Success(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v2/browser/destroy" {
				t.Errorf("Expected path /api/v2/browser/destroy, got %s", req.URL.Path)
			}

			responseBody := `{
				"code": 200,
				"data": null,
				"msg": "Environment destroyed successfully",
				"reqId": "test-req-id"
			}`

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, err := NewClient("test-key", WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &EnvDelReq{EnvId: "789"}

	err = client.EnvDestroy(context.Background(), req)
	if err != nil {
		t.Fatalf("EnvDestroy() error = %v", err)
	}
}

func TestClient_GetEnvPage_Success(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v2/browser/page" {
				t.Errorf("Expected path /api/v2/browser/page, got %s", req.URL.Path)
			}

			responseBody := `{
				"code": 200,
				"data": {
					"list": [
						{
							"envId": "1",
							"customerId": "customer1",
							"envName": "Environment 1"
						},
						{
							"envId": "2",
							"customerId": "customer1",
							"envName": "Environment 2"
						}
					],
					"total": 2,
					"pageSize": 10,
					"currentPage": 1
				},
				"msg": "success",
				"reqId": "test-req-id",
				"total": 2
			}`

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, err := NewClient("test-key", WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := &GetEnvPageReq{
		ReqPage: ReqPage{
			Page:     1,
			PageSize: 10,
		},
		CustomerId: "customer1",
	}

	resp, err := client.GetEnvPage(context.Background(), req)
	if err != nil {
		t.Fatalf("GetEnvPage() error = %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("GetEnvPage() total = %v, want 2", resp.Total)
	}

	if resp.List[0].EnvName != "Environment 1" {
		t.Errorf("GetEnvPage() first env name = %v, want Environment 1", resp.List[0].EnvName)
	}
}

// Helper types for testing
type mockHTTPError struct {
	message string
}

func (e *mockHTTPError) Error() string {
	return e.message
}

func TestWithDebug(t *testing.T) {
	client, err := NewClient("test-key", WithDebug(true))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if !client.debug {
		t.Error("WithDebug(true) should set debug = true")
	}

	// Disable debug
	opt := WithDebug(false)
	opt(client)
	if client.debug {
		t.Error("WithDebug(false) should set debug = false")
	}
}

func TestWithLogger(t *testing.T) {
	var buf bytes.Buffer
	customLogger := log.New(&buf, "[test] ", 0)

	client, err := NewClient("test-key", WithDebug(true), WithLogger(customLogger))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.logger != customLogger {
		t.Error("WithLogger() did not set the custom logger")
	}

	// Verify logger works
	client.debugLog("hello %s", "world")
	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("WithLogger() logger output = %q, want to contain 'hello world'", buf.String())
	}

	// Test with nil logger (should not change)
	opt := WithLogger(nil)
	opt(client)
	if client.logger != customLogger {
		t.Error("WithLogger(nil) should not change existing logger")
	}
}

func TestClient_DebugLogging(t *testing.T) {
	var buf bytes.Buffer
	customLogger := log.New(&buf, "", 0)

	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"code":200,"data":{"expireTime":1234567890,"userSig":"debug-sig"},"msg":"OK","reqId":"req-debug"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}

	client, err := NewClient("test-key",
		WithDebug(true),
		WithLogger(customLogger),
		WithHTTPClient(&http.Client{Transport: mockTransport}),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetUserSig(context.Background(), &GetUserSigRequest{
		CustomerId: "debug-customer",
		Duration:   3600,
	})
	if err != nil {
		t.Fatalf("GetUserSig() error = %v", err)
	}

	output := buf.String()

	// Outgoing request log should contain method and path
	if !strings.Contains(output, "POST") {
		t.Errorf("debug log should contain 'POST', got: %s", output)
	}
	if !strings.Contains(output, "/api/v2/browser/getUserSig") {
		t.Errorf("debug log should contain path, got: %s", output)
	}

	// Response log should contain status 200
	if !strings.Contains(output, "200") {
		t.Errorf("debug log should contain HTTP status '200', got: %s", output)
	}

	// Response log should contain elapsed time
	if !strings.Contains(output, "elapsed=") {
		t.Errorf("debug log should contain 'elapsed=', got: %s", output)
	}
}

func TestClient_DebugLogging_Disabled(t *testing.T) {
	var buf bytes.Buffer
	customLogger := log.New(&buf, "", 0)

	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := `{"code":200,"data":{"expireTime":1234567890,"userSig":"no-debug-sig"},"msg":"OK","reqId":"req-1"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}

	// Debug disabled (default)
	client, err := NewClient("test-key",
		WithLogger(customLogger),
		WithHTTPClient(&http.Client{Transport: mockTransport}),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetUserSig(context.Background(), &GetUserSigRequest{
		CustomerId: "test",
		Duration:   3600,
	})
	if err != nil {
		t.Fatalf("GetUserSig() error = %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("debug logging disabled: expected no output, got: %s", buf.String())
	}
}

func TestClient_GetUiFingerList_Success(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v2/browser/getUiFingerList" {
				t.Errorf("Expected path /api/v2/browser/getUiFingerList, got %s", req.URL.Path)
			}

			responseBody := `{
	"reqId": "09030b4f-6c66-4c75-bd77-2b742c9a6be4",
	"code": 200,
	"msg": "OK",
	"data": {
		"chromeKernelVersion": [{
			"id": 10140,
			"version": "140"
		}, {
			"id": 10138,
			"version": "138"
		}, {
			"id": 10134,
			"version": "134"
		}, {
			"id": 10131,
			"version": "131"
		}, {
			"id": 10127,
			"version": "127"
		}, {
			"id": 3,
			"version": "122"
		}, {
			"id": 2,
			"version": "119"
		}, {
			"id": 1,
			"version": "107"
		}],
		"firefoxKernelversion": [{
			"id": 20130,
			"version": "130"
		}],
		"system": {
			"Windows": [{
				"name": "Windows",
				"system": "All Windows",
				"browser": "Windows NT 10.0",
				"platformVersion": "15.0.0"
			}, {
				"name": "Windows",
				"system": "Windows 11",
				"browser": "Windows NT 10.0",
				"platformVersion": "15.0.0"
			}, {
				"name": "Windows",
				"system": "Windows 10",
				"browser": "Windows NT 10.0",
				"platformVersion": "10.0.0"
			}, {
				"name": "Windows",
				"system": "Windows 8.1",
				"browser": "Windows NT 6.3",
				"platformVersion": "8.1.0"
			}, {
				"name": "Windows",
				"system": "Windows 8",
				"browser": "Windows NT 6.2",
				"platformVersion": "8.0.0"
			}, {
				"name": "Windows",
				"system": "Windows 7",
				"browser": "Windows NT 6.1",
				"platformVersion": "7.0.0"
			}],
			"Linux": [{
				"name": "Linux",
				"system": "Linux",
				"browser": "X11",
				"platformVersion": ""
			}]
		},
		"chromeUAversion": ["140", "138", "134", "131", "127", "126", "125", "124", "123", "122", "121", "120", "119", "118", "117", "116", "115", "114", "113", "112", "111", "110", "109", "108", "107", "106", "105", "104", "103", "102", "101", "100"],
		"firefoxUAversion": ["130"],
		"language": [{
			"country": "非洲地区",
			"ab": "af",
			"ecountry": "Afrikaans",
			"name": "南非荷兰语",
			"code": ""
		}, {
			"country": "亚美尼亚",
			"ab": "am",
			"ecountry": "Amharic",
			"name": "阿姆哈拉语",
			"code": "AM"
		}, {
			"country": "尼加拉瓜",
			"ab": "an",
			"ecountry": "Aragonese",
			"name": "阿拉贡语",
			"code": "NI"
		}, {
			"country": "西班牙",
			"ab": "ast",
			"ecountry": "Asturian",
			"name": "阿斯图里亚斯语",
			"code": "ES"
		}, {
			"country": "阿塞拜疆",
			"ab": "az",
			"ecountry": "Azerbaijani",
			"name": "阿塞拜疆语",
			"code": "AZ"
		}, {
			"country": "保加利亚",
			"ab": "bg",
			"ecountry": "Bulgarian",
			"name": "保加利亚语",
			"code": "BG"
		}, {
			"country": "文莱",
			"ab": "bn",
			"ecountry": "Bengali",
			"name": "孟加拉语",
			"code": "BN"
		}, {
			"country": "巴哈马",
			"ab": "bs",
			"ecountry": "Bosnian",
			"name": "波斯尼亚语",
			"code": "BS"
		}, {
			"country": "加拿大",
			"ab": "ca",
			"ecountry": "Catalan",
			"name": "加泰隆语",
			"code": "CA"
		}, {
			"country": "菲律宾",
			"ab": "ceb",
			"ecountry": "Cebuano",
			"name": "丹麦语",
			"code": "PH"
		}, {
			"country": "哥伦比亚",
			"ab": "co",
			"ecountry": "Corsican",
			"name": "科西嘉语",
			"code": "CO"
		}, {
			"country": "捷克",
			"ab": "cs",
			"ecountry": "Czech",
			"name": "捷克语",
			"code": "CZ"
		}, {
			"country": "塞浦路斯",
			"ab": "cy",
			"ecountry": "Welsh",
			"name": "威尔士语",
			"code": "CY"
		}, {
			"country": "丹麦",
			"ab": "da",
			"ecountry": "Danish",
			"name": "丹麦语",
			"code": "DK"
		}, {
			"country": "德国",
			"ab": "de",
			"ecountry": "German",
			"name": "德语",
			"code": "DE"
		}, {
			"country": "奥地利",
			"ab": "de-AT",
			"ecountry": "German (Austria)",
			"name": "德语-奥地利变体",
			"code": "AT"
		}, {
			"country": "瑞士",
			"ab": "de-CH",
			"ecountry": "German (Switzerland)",
			"name": "德语-瑞士变体",
			"code": "CH"
		}, {
			"country": "德国",
			"ab": "de-DE",
			"ecountry": "German (Germany)",
			"name": "德语-德国变体",
			"code": "DE"
		}, {
			"country": "列支敦士登",
			"ab": "de-LI",
			"ecountry": "German (Liechtenstein)",
			"name": "德语-列支敦士登变体",
			"code": "LI"
		}, {
			"country": "希腊",
			"ab": "el",
			"ecountry": "Greek",
			"name": "现代希腊语",
			"code": "GR"
		}, {
			"country": "英国",
			"ab": "en-US",
			"ecountry": "English",
			"name": "英语",
			"code": "GB"
		}, {
			"country": "澳大利亚",
			"ab": "en-AU",
			"ecountry": "English (Australia)",
			"name": "澳大利亚-英语",
			"code": "AU"
		}, {
			"country": "加拿大",
			"ab": "en-CA",
			"ecountry": "English (Canada)",
			"name": "加拿大-英语",
			"code": "CA"
		}, {
			"country": "英国",
			"ab": "en-GB",
			"ecountry": "English (UK)",
			"name": "英语",
			"code": "GB"
		}, {
			"country": "英国",
			"ab": "en-GB-oxendict",
			"ecountry": "English (UK, OED spelling)",
			"name": "英国-英语",
			"code": "GB"
		}, {
			"country": "印度",
			"ab": "en-IN",
			"ecountry": "English (India)",
			"name": "印度-英语",
			"code": "IN"
		}, {
			"country": "新西兰",
			"ab": "en-NZ",
			"ecountry": "English (New Zealand)",
			"name": "新西兰-英语",
			"code": "NZ"
		}, {
			"country": "南非",
			"ab": "en-ZA",
			"ecountry": "English (South Africa)",
			"name": "南非-英语",
			"code": "ZA"
		}, {
			"country": "西班牙",
			"ab": "eo",
			"ecountry": "Esperanto",
			"name": "世界语",
			"code": "ES"
		}, {
			"country": "西班牙",
			"ab": "es",
			"ecountry": "Spanish",
			"name": "西班牙语",
			"code": "ES"
		}, {
			"country": "墨西哥",
			"ab": "es-419",
			"ecountry": "Spanish (Latin America)",
			"name": "墨西哥-西班牙语",
			"code": "MX"
		}, {
			"country": "阿根廷",
			"ab": "es-AR",
			"ecountry": "Spanish (Argentina)",
			"name": "阿根廷-西班牙语",
			"code": "AR"
		}, {
			"country": "智利",
			"ab": "es-CL",
			"ecountry": "Spanish (Chile)",
			"name": "智利-西班牙语",
			"code": "CL"
		}, {
			"country": "哥伦比亚",
			"ab": "es-CO",
			"ecountry": "Spanish (Colombia)",
			"name": "哥伦比亚-西班牙语",
			"code": "CO"
		}, {
			"country": "哥斯达黎加",
			"ab": "es-CR",
			"ecountry": "Spanish (Costa Rica)",
			"name": "哥斯达黎加-西班牙语",
			"code": "CR"
		}, {
			"country": "西班牙",
			"ab": "es-ES",
			"ecountry": "Spanish (Spain)",
			"name": "西班牙-西班牙语",
			"code": "ES"
		}, {
			"country": "洪都拉斯",
			"ab": "es-HN",
			"ecountry": "Spanish (Honduras)",
			"name": "洪都拉斯-西班牙语",
			"code": "HN"
		}, {
			"country": "墨西哥",
			"ab": "es-MX",
			"ecountry": "Spanish (Mexico)",
			"name": "墨西哥-西班牙语",
			"code": "MX"
		}, {
			"country": "秘鲁",
			"ab": "es-PE",
			"ecountry": "Spanish (Peru)",
			"name": "秘鲁-西班牙语",
			"code": "PE"
		}, {
			"country": "乌拉圭",
			"ab": "es-UY",
			"ecountry": "Spanish (Uruguay)",
			"name": "乌拉圭-西班牙语",
			"code": "UY"
		}, {
			"country": "委内瑞拉",
			"ab": "es-VE",
			"ecountry": "Spanish (Venezuela)",
			"name": "委内瑞拉-西班牙语",
			"code": "VE"
		}, {
			"country": "爱沙尼亚",
			"ab": "et",
			"ecountry": "Estonian",
			"name": "爱沙尼亚语",
			"code": "EE"
		}, {
			"country": "西班牙",
			"ab": "eu",
			"ecountry": "Basque",
			"name": "英语",
			"code": "ES"
		}, {
			"country": "伊朗",
			"ab": "fa",
			"ecountry": "Persian",
			"name": "波斯语",
			"code": "IR"
		}, {
			"country": "芬兰",
			"ab": "fi",
			"ecountry": "Finnish",
			"name": "芬兰语",
			"code": "FI"
		}, {
			"country": "菲律宾",
			"ab": "fil",
			"ecountry": "Filipino",
			"name": "菲律宾语",
			"code": "PH"
		}, {
			"country": "丹麦",
			"ab": "fo",
			"ecountry": "Faroese",
			"name": "法罗语",
			"code": "DK"
		}, {
			"country": "法国",
			"ab": "fr",
			"ecountry": "French",
			"name": "法语",
			"code": "FR"
		}, {
			"country": "加拿大",
			"ab": "fr-CA",
			"ecountry": "French (Canada)",
			"name": "加拿大-法语",
			"code": "CA"
		}, {
			"country": "瑞士",
			"ab": "fr-CH",
			"ecountry": "French (Switzerland)",
			"name": "瑞士-法语",
			"code": "CH"
		}, {
			"country": "法国",
			"ab": "fr-FR",
			"ecountry": "French (France)",
			"name": "法国-法语",
			"code": "FR"
		}, {
			"country": "荷兰",
			"ab": "fy",
			"ecountry": "Frisian",
			"name": "弗里西亚语",
			"code": "NL"
		}, {
			"country": "爱尔兰",
			"ab": "ga",
			"ecountry": "Irish",
			"name": "爱尔兰语",
			"code": "IE"
		}, {
			"country": "英国",
			"ab": "gd",
			"ecountry": "Scots Gaelic",
			"name": "苏格兰盖尔语",
			"code": "GB"
		}, {
			"country": "西班牙",
			"ab": "gl",
			"ecountry": "Galician",
			"name": "加利西亚语",
			"code": "ES"
		}, {
			"country": "巴拉圭",
			"ab": "gn",
			"ecountry": "Guarani",
			"name": "瓜拉尼语",
			"code": "PY"
		}, {
			"country": "印度",
			"ab": "gu",
			"ecountry": "Gujarati",
			"name": "古吉拉特语",
			"code": "IN"
		}, {
			"country": "尼日利亚",
			"ab": "ha",
			"ecountry": "Hausa",
			"name": "豪萨语",
			"code": "NG"
		}, {
			"country": "马来西亚",
			"ab": "haw",
			"ecountry": "Hawaiian",
			"name": "夏威夷语",
			"code": "MY"
		}, {
			"country": "印度",
			"ab": "hi",
			"ecountry": "Hindi",
			"name": "印地语",
			"code": "IN"
		}, {
			"country": "越南",
			"ab": "hmn",
			"ecountry": "Hmong",
			"name": "苗族语",
			"code": "VN"
		}, {
			"country": "克罗地亚",
			"ab": "hr",
			"ecountry": "Croatian",
			"name": "克罗地亚语",
			"code": "HR"
		}, {
			"country": "加拿大",
			"ab": "ht",
			"ecountry": "Haitian Creole",
			"name": "海地克里奥尔语",
			"code": "CA"
		}, {
			"country": "匈牙利",
			"ab": "hu",
			"ecountry": "Hungarian",
			"name": "匈牙利语",
			"code": "HU"
		}, {
			"country": "亚美尼亚",
			"ab": "hy",
			"ecountry": "Armenian",
			"name": "亚美尼亚语",
			"code": "AM"
		}, {
			"country": "印度尼西亚",
			"ab": "id",
			"ecountry": "Indonesian",
			"name": "印度尼西亚语",
			"code": "ID"
		}, {
			"country": "印度尼西亚",
			"ab": "id",
			"ecountry": "Indonesian",
			"name": "印度尼西亚语",
			"code": "ID"
		}, {
			"country": "尼日利亚",
			"ab": "ig",
			"ecountry": "Igbo",
			"name": "伊博语",
			"code": "NG"
		}, {
			"country": "冰岛",
			"ab": "is",
			"ecountry": "Icelandic",
			"name": "冰岛语",
			"code": "IS"
		}, {
			"country": "意大利",
			"ab": "it",
			"ecountry": "Italian",
			"name": "意大利语",
			"code": "IT"
		}, {
			"country": "瑞士",
			"ab": "it-CH",
			"ecountry": "Italian (Switzerland)",
			"name": "意大利语",
			"code": "CH"
		}, {
			"country": "意大利",
			"ab": "it-IT",
			"ecountry": "Italian (Italy)",
			"name": "意大利语",
			"code": "IT"
		}, {
			"country": "日本",
			"ab": "ja",
			"ecountry": "Japanese",
			"name": "日语",
			"code": "JP"
		}, {
			"country": "印度尼西亚",
			"ab": "jv",
			"ecountry": "Javanese",
			"name": "爪哇语",
			"code": "ID"
		}, {
			"country": "格鲁吉亚",
			"ab": "ka",
			"ecountry": "Georgian",
			"name": "格鲁吉亚语",
			"code": "GE"
		}, {
			"country": "哈萨克斯坦",
			"ab": "kk",
			"ecountry": "Kazakh",
			"name": "哈萨克语",
			"code": "KAZ"
		}, {
			"country": "柬埔寨",
			"ab": "km",
			"ecountry": "Cambodian",
			"name": "高棉语",
			"code": "KH"
		}, {
			"country": "印度",
			"ab": "kn",
			"ecountry": "Kannada",
			"name": "卡纳达语",
			"code": "IN"
		}, {
			"country": "韩国",
			"ab": "ko",
			"ecountry": "Korean",
			"name": "韩语",
			"code": "KR"
		}, {
			"country": "印度",
			"ab": "kok",
			"ecountry": "Konkani",
			"name": "昆卡尼语",
			"code": "IN"
		}, {
			"country": "伊朗",
			"ab": "ku",
			"ecountry": "Kurdish",
			"name": "库尔德语",
			"code": "IR"
		}, {
			"country": "吉尔吉斯斯坦",
			"ab": "ky",
			"ecountry": "Kyrgyz",
			"name": "吉尔吉斯语",
			"code": "KG"
		}, {
			"country": "意大利",
			"ab": "la",
			"ecountry": "Latin",
			"name": "拉丁语",
			"code": "IT"
		}, {
			"country": "卢森堡",
			"ab": "lb",
			"ecountry": "Luxembourgish",
			"name": "卢森堡语",
			"code": "LU"
		}, {
			"country": "刚果金",
			"ab": "ln",
			"ecountry": "Lingala",
			"name": "林加拉语",
			"code": "CD"
		}, {
			"country": "刚果布",
			"ab": "ln",
			"ecountry": "Lingala",
			"name": "林加拉语",
			"code": "CG"
		}, {
			"country": "老挝",
			"ab": "lo",
			"ecountry": "Laothian",
			"name": "老挝语",
			"code": "LA"
		}, {
			"country": "立陶宛",
			"ab": "lt",
			"ecountry": "Lithuanian",
			"name": "立陶宛语",
			"code": "LT"
		}, {
			"country": "拉脱维亚",
			"ab": "lv",
			"ecountry": "Latvian",
			"name": "拉脱维亚语",
			"code": "LV"
		}, {
			"country": "马尔代夫",
			"ab": "mg",
			"ecountry": "Malagasy",
			"name": "马尔加什语",
			"code": "MV"
		}, {
			"country": "新西兰",
			"ab": "mi",
			"ecountry": "Maori",
			"name": "毛利语",
			"code": "NZ"
		}, {
			"country": "北马其顿",
			"ab": "mk",
			"ecountry": "Macedonian",
			"name": "马其顿语",
			"code": "MK"
		}, {
			"country": "印度",
			"ab": "ml",
			"ecountry": "Malayalam",
			"name": "马拉亚拉姆语",
			"code": "IN"
		}, {
			"country": "蒙古",
			"ab": "mn",
			"ecountry": "Mongolian",
			"name": "蒙古语",
			"code": "MNG"
		}, {
			"country": "摩尔多瓦",
			"ab": "mo",
			"ecountry": "Moldavian",
			"name": "摩尔多瓦语",
			"code": "MD"
		}, {
			"country": "印度",
			"ab": "mr",
			"ecountry": "Marathi",
			"name": "马拉地语",
			"code": "IN"
		}, {
			"country": "马来西亚",
			"ab": "ms",
			"ecountry": "Malay",
			"name": "马来语",
			"code": "MY"
		}, {
			"country": "马耳他",
			"ab": "mt",
			"ecountry": "Maltese",
			"name": "马耳他语",
			"code": "MT"
		}, {
			"country": "缅甸",
			"ab": "my",
			"ecountry": "Burmese",
			"name": "缅甸语",
			"code": "MM"
		}, {
			"country": "挪威",
			"ab": "nb",
			"ecountry": "Norwegian (Bokmal)",
			"name": "挪威语",
			"code": "NO"
		}, {
			"country": "尼泊尔",
			"ab": "ne",
			"ecountry": "Nepali",
			"name": "尼泊尔语",
			"code": "NP"
		}, {
			"country": "荷兰",
			"ab": "nl",
			"ecountry": "Dutch",
			"name": "荷兰语",
			"code": "NL"
		}, {
			"country": "挪威",
			"ab": "nn",
			"ecountry": "Norwegian (Nynorsk)",
			"name": "新挪威语",
			"code": "NO"
		}, {
			"country": "挪威",
			"ab": "no",
			"ecountry": "Norwegian",
			"name": "挪威语",
			"code": "NO"
		}, {
			"country": "加拿大",
			"ab": "ny",
			"ecountry": "Nyanja",
			"name": "尼扬扎语",
			"code": "CA"
		}, {
			"country": "法国",
			"ab": "oc",
			"ecountry": "Occitan",
			"name": "奥克语",
			"code": "FR"
		}, {
			"country": "埃塞俄比亚",
			"ab": "om",
			"ecountry": "Oromo",
			"name": "奥罗莫语",
			"code": "ET"
		}, {
			"country": "印度",
			"ab": "or",
			"ecountry": "Odia (Oriya)",
			"name": "奥利亚语",
			"code": "IN"
		}, {
			"country": "巴基斯坦",
			"ab": "pa",
			"ecountry": "Punjabi",
			"name": "旁遮普语",
			"code": "PK"
		}, {
			"country": "波兰",
			"ab": "pl",
			"ecountry": "Polish",
			"name": "波兰语",
			"code": "PL"
		}, {
			"country": "阿富汗",
			"ab": "ps",
			"ecountry": "Pashto",
			"name": "普什图语",
			"code": "AF"
		}, {
			"country": "葡萄牙",
			"ab": "pt",
			"ecountry": "Portuguese",
			"name": "葡萄牙语",
			"code": "PT"
		}, {
			"country": "巴西",
			"ab": "pt-BR",
			"ecountry": "Portuguese (Brazil)",
			"name": "巴西-葡萄牙语",
			"code": "BR"
		}, {
			"country": "葡萄牙",
			"ab": "pt-PT",
			"ecountry": "Portuguese (Portugal)",
			"name": "葡萄牙-葡萄牙语",
			"code": "PT"
		}, {
			"country": "阿根廷",
			"ab": "qu",
			"ecountry": "Quechua",
			"name": "克丘亚语",
			"code": "AR"
		}, {
			"country": "瑞士",
			"ab": "rm",
			"ecountry": "Romansh",
			"name": "罗曼什语",
			"code": "CH"
		}, {
			"country": "罗马尼亚",
			"ab": "ro",
			"ecountry": "Romanian",
			"name": "罗马尼亚语",
			"code": "RO"
		}, {
			"country": "俄罗斯",
			"ab": "ru",
			"ecountry": "Russian",
			"name": "俄语",
			"code": "RU"
		}, {
			"country": "卢旺达",
			"ab": "rw",
			"ecountry": "Kinyarwanda",
			"name": "基尼亚鲁旺达语",
			"code": "RW"
		}, {
			"country": "印度",
			"ab": "sd",
			"ecountry": "Sindhi",
			"name": "信德语",
			"code": "IN"
		}, {
			"country": "塞尔维亚",
			"ab": "sh",
			"ecountry": "Serbo-Croatian",
			"name": "塞尔维亚-克罗地亚语",
			"code": "RS"
		}, {
			"country": "斯里兰卡",
			"ab": "si",
			"ecountry": "Sinhalese",
			"name": "僧伽罗语",
			"code": "LK"
		}, {
			"country": "斯洛伐克",
			"ab": "sk",
			"ecountry": "Slovak",
			"name": "斯洛伐克语",
			"code": "SK"
		}, {
			"country": "斯洛文尼亚",
			"ab": "sl",
			"ecountry": "Slovenian",
			"name": "斯洛文尼亚语",
			"code": "SI"
		}, {
			"country": "萨摩亚",
			"ab": "sm",
			"ecountry": "Samoan",
			"name": "萨摩亚语",
			"code": "WS"
		}, {
			"country": "津巴布韦",
			"ab": "sn",
			"ecountry": "Shona",
			"name": "绍纳语",
			"code": "ZW"
		}, {
			"country": "索马里",
			"ab": "so",
			"ecountry": "Somali",
			"name": "索马里语",
			"code": "SO"
		}, {
			"country": "阿尔巴尼亚",
			"ab": "sq",
			"ecountry": "Albanian",
			"name": "阿尔巴尼亚语",
			"code": "AL"
		}, {
			"country": "塞尔维亚",
			"ab": "sr",
			"ecountry": "Serbian",
			"name": "塞尔维亚语",
			"code": "RS"
		}, {
			"country": "南非",
			"ab": "st",
			"ecountry": "Sesotho",
			"name": "塞索托语",
			"code": "ZA"
		}, {
			"country": "印度尼西亚",
			"ab": "su",
			"ecountry": "Sundanese",
			"name": "巽他语",
			"code": "ID"
		}, {
			"country": "瑞典",
			"ab": "sv",
			"ecountry": "Swedish",
			"name": "瑞典语",
			"code": "SE"
		}, {
			"country": "坦桑尼亚",
			"ab": "sw",
			"ecountry": "Swahili",
			"name": "斯瓦希里语",
			"code": "TZ"
		}, {
			"country": "印度",
			"ab": "ta",
			"ecountry": "Tamil",
			"name": "泰米尔语",
			"code": "IN"
		}, {
			"country": "印度",
			"ab": "te",
			"ecountry": "Telugu",
			"name": "泰卢固语",
			"code": "IN"
		}, {
			"country": "塔吉克斯坦",
			"ab": "tg",
			"ecountry": "Tajik",
			"name": "塔吉克语",
			"code": "TJ"
		}, {
			"country": "泰国",
			"ab": "th",
			"ecountry": "Thai",
			"name": "泰语",
			"code": "TH"
		}, {
			"country": "埃塞俄比亚",
			"ab": "ti",
			"ecountry": "Tigrinya",
			"name": "提格利尼亚语",
			"code": "ET"
		}, {
			"country": "土库曼斯坦",
			"ab": "tk",
			"ecountry": "Turkmen",
			"name": "土库曼语",
			"code": "TM"
		}, {
			"country": "南非",
			"ab": "tn",
			"ecountry": "Tswana",
			"name": "茨瓦纳语",
			"code": "ZA"
		}, {
			"country": "土耳其",
			"ab": "tr",
			"ecountry": "Turkish",
			"name": "土耳其语",
			"code": "TR"
		}, {
			"country": "乌克兰",
			"ab": "tt",
			"ecountry": "Tatar",
			"name": "鞑靼语",
			"code": "UA"
		}, {
			"country": "加纳",
			"ab": "tw",
			"ecountry": "Twi",
			"name": "特威语",
			"code": "GH"
		}, {
			"country": "乌兹别克斯坦",
			"ab": "ug",
			"ecountry": "Uyghur",
			"name": "维吾尔语",
			"code": "UZ"
		}, {
			"country": "乌克兰",
			"ab": "uk",
			"ecountry": "Ukrainian",
			"name": "乌克兰语",
			"code": "UA"
		}, {
			"country": "印度",
			"ab": "ur",
			"ecountry": "Urdu",
			"name": "乌尔都语",
			"code": "IN"
		}, {
			"country": "乌兹别克斯坦",
			"ab": "uz",
			"ecountry": "Uzbek",
			"name": "乌兹别克语",
			"code": "UZ"
		}, {
			"country": "越南",
			"ab": "vi",
			"ecountry": "Vietnamese",
			"name": "越南语",
			"code": "VN"
		}, {
			"country": "比利时",
			"ab": "wa",
			"ecountry": "Walloon",
			"name": "瓦隆语",
			"code": "BE"
		}, {
			"country": "塞内加尔",
			"ab": "wo",
			"ecountry": "Wolof",
			"name": "沃洛夫语",
			"code": "SN"
		}, {
			"country": "南非",
			"ab": "xh",
			"ecountry": "Xhosa",
			"name": "科萨语",
			"code": "ZA"
		}, {
			"country": "德国",
			"ab": "yi",
			"ecountry": "Yiddish",
			"name": "意第绪语",
			"code": "DE"
		}, {
			"country": "尼日利亚",
			"ab": "yo",
			"ecountry": "Yoruba",
			"name": "约鲁巴语",
			"code": "NG"
		}, {
			"country": "中国",
			"ab": "zh-CN",
			"ecountry": "Chinese",
			"name": "中文",
			"code": "CN"
		}, {
			"country": "中国",
			"ab": "zh",
			"ecountry": "Chinese (China)",
			"name": "中国（大陆）-中文",
			"code": "CN"
		}, {
			"country": "中国",
			"ab": "zh-HK",
			"ecountry": "Chinese (Hong Kong)",
			"name": "中国香港-中文",
			"code": "HK"
		}, {
			"country": "中国",
			"ab": "zh-TW",
			"ecountry": "Chinese (Taiwan)",
			"name": "中国台湾-中文",
			"code": "TW"
		}, {
			"country": "新加坡",
			"ab": "zh-SG",
			"ecountry": "Singapore",
			"name": "新加坡-中文",
			"code": "SG"
		}, {
			"country": "新加坡",
			"ab": "en-SG",
			"ecountry": "Singapore",
			"name": "新加坡-英语",
			"code": "SG"
		}, {
			"country": "沙特阿拉伯",
			"ab": "ar-SA",
			"ecountry": "Arabic",
			"name": "沙特阿拉伯-阿拉伯语",
			"code": "SA"
		}, {
			"country": "阿联酋",
			"ab": "ar-AE",
			"ecountry": "United Arab Emirates",
			"name": "阿拉伯联合酋长国-阿拉伯语",
			"code": "AE"
		}, {
			"country": "巴林",
			"ab": "ar-BH",
			"ecountry": "Bahrian",
			"name": "巴林-阿拉伯语",
			"code": "BH"
		}, {
			"country": "阿尔及利亚",
			"ab": "ar-DZ",
			"ecountry": "Algeria",
			"name": "阿尔及利亚-阿拉伯语",
			"code": "DZ"
		}, {
			"country": "埃及",
			"ab": "ar-EG",
			"ecountry": "Egypt",
			"name": "埃及-阿拉伯语",
			"code": "EG"
		}, {
			"country": "伊拉克",
			"ab": "ar-IQ",
			"ecountry": "Iraq",
			"name": "伊拉克-阿拉伯语",
			"code": "IQ"
		}, {
			"country": "约旦",
			"ab": "ar-JO",
			"ecountry": "Jordan",
			"name": "约旦-阿拉伯语",
			"code": "JO"
		}, {
			"country": "科威特",
			"ab": "ar-KW",
			"ecountry": "Kuwait",
			"name": "科威特-阿拉伯语",
			"code": "KW"
		}, {
			"country": "黎巴嫩",
			"ab": "ar-LB",
			"ecountry": "Lebanon",
			"name": "黎巴嫩-阿拉伯语",
			"code": "LB"
		}, {
			"country": "利比亚",
			"ab": "ar-LY",
			"ecountry": "Libyan Arab Jm",
			"name": "利比亚-阿拉伯语",
			"code": "LY"
		}, {
			"country": "摩洛哥",
			"ab": "ar-MA",
			"ecountry": "Morocco",
			"name": "摩洛哥-阿拉伯语",
			"code": "MA"
		}, {
			"country": "阿曼",
			"ab": "ar-OM",
			"ecountry": "Oman",
			"name": "阿曼-阿拉伯语",
			"code": "OM"
		}, {
			"country": "卡塔尔",
			"ab": "ar-QA",
			"ecountry": "Qatar",
			"name": "卡塔尔-阿拉伯语",
			"code": "QA"
		}, {
			"country": "叙利亚",
			"ab": "ar-SY",
			"ecountry": "Syrian",
			"name": "叙利亚-阿拉伯语",
			"code": "SY"
		}, {
			"country": "突尼斯",
			"ab": "ar-TN",
			"ecountry": "Tunisia",
			"name": "突尼斯-阿拉伯语",
			"code": "TN"
		}, {
			"country": "也门",
			"ab": "ar-YE",
			"ecountry": "Yemen",
			"name": "也门-阿拉伯语",
			"code": "YE"
		}, {
			"country": "伯利兹",
			"ab": "en-BZ",
			"ecountry": "Belize",
			"name": "伯利兹-英语",
			"code": "BZ"
		}, {
			"country": "牙买加",
			"ab": "en-JM",
			"ecountry": "Jamaica",
			"name": "牙买加-英语",
			"code": "JM"
		}, {
			"country": "特立尼达和多巴哥",
			"ab": "en-TT",
			"ecountry": "Trinidad and Tobago",
			"name": "特立尼达和多巴哥-英语",
			"code": "TT"
		}, {
			"country": "多米尼加",
			"ab": "es-DO",
			"ecountry": "Dominican",
			"name": "多米尼加-西班牙语",
			"code": "DO"
		}, {
			"country": "巴拿马",
			"ab": "es-PA",
			"ecountry": "Panama",
			"name": "巴拿马-西班牙语",
			"code": "PA"
		}, {
			"country": "波多黎各",
			"ab": "es-PR",
			"ecountry": "Puerto Rico",
			"name": "波多黎各-西班牙语",
			"code": ""
		}, {
			"country": "萨尔瓦多",
			"ab": "es-SV",
			"ecountry": "El Salvador",
			"name": "萨尔瓦多-西班牙语",
			"code": "SV"
		}, {
			"country": "法罗群岛",
			"ab": "fo-FO",
			"ecountry": "Faroe Islands",
			"name": "法罗群岛-法罗语",
			"code": "FO"
		}, {
			"country": "摩纳哥",
			"ab": "fr-MC",
			"ecountry": "Monaco",
			"name": "摩纳哥-法语",
			"code": "MC"
		}, {
			"country": "以色列",
			"ab": "he-IL",
			"ecountry": "Israel",
			"name": "以色列-希伯来语",
			"code": "IL"
		}, {
			"country": "格陵兰",
			"ab": "kl-GL",
			"ecountry": "Greenland",
			"name": "格陵兰-格陵兰语",
			"code": "GL"
		}, {
			"country": "蒙古国",
			"ab": "mn-Mong",
			"ecountry": "Mongolia",
			"name": "蒙古国-蒙古语",
			"code": "MN"
		}, {
			"country": "危地马拉",
			"ab": "qut-GT",
			"ecountry": "Guatemala",
			"name": "危地马拉-基切语",
			"code": "GT"
		}, {
			"country": "玻利维亚",
			"ab": "quz-BO",
			"ecountry": "Bolivia",
			"name": "玻利维亚-克丘亚语",
			"code": "BO"
		}, {
			"country": "厄瓜多尔",
			"ab": "quz-EC",
			"ecountry": "Ecuador",
			"name": "厄瓜多尔-克丘亚语",
			"code": "EC"
		}, {
			"country": "肯尼亚",
			"ab": "sw-KE",
			"ecountry": "Kenya",
			"name": "肯尼亚-斯瓦希里语",
			"code": "KE"
		}, {
			"country": "南非",
			"ab": "zu",
			"ecountry": "Zulu",
			"name": "祖鲁语",
			"code": "ZA"
		}],
		"zone": ["GMT-12:00 Etc/GMT+12", "GMT-11:00 Etc/GMT+11", "GMT-11:00 Pacific/Midway", "GMT-11:00 Pacific/Niue", "GMT-11:00 Pacific/Pago_Pago", "GMT-10:00 America/Adak", "GMT-10:00 Etc/GMT+10", "GMT-10:00 HST", "GMT-10:00 Pacific/Honolulu", "GMT-10:00 Pacific/Rarotonga", "GMT-10:00 Pacific/Tahiti", "GMT-09:30 Pacific/Marquesas", "GMT-09:00 America/Metlakatla", "GMT-09:00 America/Anchorage", "GMT-09:00 America/Juneau", "GMT-09:00 America/Nome", "GMT-09:00 America/Sitka", "GMT-09:00 America/Yakutat", "GMT-09:00 Etc/GMT+9", "GMT-09:00 Pacific/Gambier", "GMT-08:00 America/Los_Angeles", "GMT-08:00 America/Tijuana", "GMT-08:00 America/Vancouver", "GMT-08:00 Etc/GMT+8", "GMT-08:00 PST8PDT", "GMT-08:00 Pacific/Pitcairn", "GMT-07:00 America/Boise", "GMT-07:00 America/Cambridge_Bay", "GMT-07:00 America/Chihuahua", "GMT-07:00 America/Creston", "GMT-07:00 America/Dawson", "GMT-07:00 America/Dawson_Creek", "GMT-07:00 America/Denver", "GMT-07:00 America/Edmonton", "GMT-07:00 America/Fort_Nelson", "GMT-07:00 America/Hermosillo", "GMT-07:00 America/Inuvik", "GMT-07:00 America/Mazatlan", "GMT-07:00 America/Ojinaga", "GMT-07:00 America/Phoenix", "GMT-07:00 America/Whitehorse", "GMT-07:00 America/Yellowknife", "GMT-07:00 Etc/GMT+7", "GMT-07:00 MST", "GMT-07:00 MST7MDT", "GMT-06:00 America/Bahia_Banderas", "GMT-06:00 America/Belize", "GMT-06:00 America/Chicago", "GMT-06:00 America/Costa_Rica", "GMT-06:00 America/El_Salvador", "GMT-06:00 America/Guatemala", "GMT-06:00 America/Indiana/Knox", "GMT-06:00 America/Indiana/Tell_City", "GMT-06:00 America/Managua", "GMT-06:00 America/Matamoros", "GMT-06:00 America/Menominee", "GMT-06:00 America/Merida", "GMT-06:00 America/Mexico_City", "GMT-06:00 America/Monterrey", "GMT-06:00 America/North_Dakota/Beulah", "GMT-06:00 America/North_Dakota/Center", "GMT-06:00 America/North_Dakota/New_Salem", "GMT-06:00 America/Rainy_River", "GMT-06:00 America/Rankin_Inlet", "GMT-06:00 America/Regina", "GMT-06:00 America/Resolute", "GMT-06:00 America/Swift_Current", "GMT-06:00 America/Tegucigalpa", "GMT-06:00 America/Winnipeg", "GMT-06:00 CST6CDT", "GMT-06:00 Etc/GMT+6", "GMT-06:00 Pacific/Galapagos", "GMT-05:00 America/Atikokan", "GMT-05:00 America/Bogota", "GMT-05:00 America/Cancun", "GMT-05:00 America/Cayman", "GMT-05:00 America/Detroit", "GMT-05:00 America/Eirunepe", "GMT-05:00 America/Grand_Turk", "GMT-05:00 America/Guayaquil", "GMT-05:00 America/Havana", "GMT-05:00 America/Indiana/Indianapolis", "GMT-05:00 America/Indiana/Marengo", "GMT-05:00 America/Indiana/Petersburg", "GMT-05:00 America/Indiana/Vevay", "GMT-05:00 America/Indiana/Vincennes", "GMT-05:00 America/Indiana/Winamac", "GMT-05:00 America/Indianapolis", "GMT-05:00 America/Iqaluit", "GMT-05:00 America/Jamaica", "GMT-05:00 America/Kentucky/Louisville", "GMT-05:00 America/Kentucky/Monticello", "GMT-05:00 America/Lima", "GMT-05:00 America/Montreal", "GMT-05:00 America/Nassau", "GMT-05:00 America/New_York", "GMT-05:00 America/Nipigon", "GMT-05:00 America/Panama", "GMT-05:00 America/Pangnirtung", "GMT-05:00 America/Port-au-Prince", "GMT-05:00 America/Rio_Branco", "GMT-05:00 America/Thunder_Bay", "GMT-05:00 America/Toronto", "GMT-05:00 EST", "GMT-05:00 EST5EDT", "GMT-05:00 Etc/GMT+5", "GMT-05:00 Pacific/Easter", "GMT-04:00 America/Anguilla", "GMT-04:00 America/Antigua", "GMT-04:00 America/Aruba", "GMT-04:00 America/Barbados", "GMT-04:00 America/Blanc-Sablon", "GMT-04:00 America/Boa_Vista", "GMT-04:00 America/Campo_Grande", "GMT-04:00 America/Caracas", "GMT-04:00 America/Cuiaba", "GMT-04:00 America/Curacao", "GMT-04:00 America/Dominica", "GMT-04:00 America/Glace_Bay", "GMT-04:00 America/Goose_Bay", "GMT-04:00 America/Grenada", "GMT-04:00 America/Guadeloupe", "GMT-04:00 America/Guyana", "GMT-04:00 America/Halifax", "GMT-04:00 America/Kralendijk", "GMT-04:00 America/La_Paz", "GMT-04:00 America/Lower_Princes", "GMT-04:00 America/Manaus", "GMT-04:00 America/Marigot", "GMT-04:00 America/Martinique", "GMT-04:00 America/Moncton", "GMT-04:00 America/Montserrat", "GMT-04:00 America/Port_of_Spain", "GMT-04:00 America/Porto_Velho", "GMT-04:00 America/Puerto_Rico", "GMT-04:00 America/Santo_Domingo", "GMT-04:00 America/St_Barthelemy", "GMT-04:00 America/St_Kitts", "GMT-04:00 America/St_Lucia", "GMT-04:00 America/St_Thomas", "GMT-04:00 America/St_Vincent", "GMT-04:00 America/Thule", "GMT-04:00 America/Tortola", "GMT-04:00 Atlantic/Bermuda", "GMT-04:00 Etc/GMT+4", "GMT-03:30 America/St_Johns", "GMT-03:00 America/Araguaina", "GMT-03:00 America/Argentina/Buenos_Aires", "GMT-03:00 America/Argentina/Catamarca", "GMT-03:00 America/Argentina/Cordoba", "GMT-03:00 America/Argentina/Jujuy", "GMT-03:00 America/Argentina/La_Rioja", "GMT-03:00 America/Argentina/Mendoza", "GMT-03:00 America/Argentina/Rio_Gallegos", "GMT-03:00 America/Argentina/Salta", "GMT-03:00 America/Argentina/San_Juan", "GMT-03:00 America/Argentina/San_Luis", "GMT-03:00 America/Argentina/Tucuman", "GMT-03:00 America/Argentina/Ushuaia", "GMT-03:00 America/Asuncion", "GMT-03:00 America/Bahia", "GMT-03:00 America/Belem", "GMT-03:00 America/Cayenne", "GMT-03:00 America/Fortaleza", "GMT-03:00 America/Godthab", "GMT-03:00 America/Maceio", "GMT-03:00 America/Miquelon", "GMT-03:00 America/Montevideo", "GMT-03:00 America/Nuuk", "GMT-03:00 America/Paramaribo", "GMT-03:00 America/Punta_Arenas", "GMT-03:00 America/Recife", "GMT-03:00 America/Santarem", "GMT-03:00 America/Santiago", "GMT-03:00 America/Sao_Paulo", "GMT-03:00 Antarctica/Palmer", "GMT-03:00 Antarctica/Rothera", "GMT-03:00 Atlantic/Stanley", "GMT-03:00 Etc/GMT+3", "GMT-02:00 America/Noronha", "GMT-02:00 Atlantic/South_Georgia", "GMT-02:00 Etc/GMT+2", "GMT-01:00 America/Scoresbysund", "GMT-01:00 Atlantic/Azores", "GMT-01:00 Atlantic/Cape_Verde", "GMT-01:00 Etc/GMT+1", "GMT+00:00 Africa/Abidjan", "GMT+00:00 Africa/Accra", "GMT+00:00 Africa/Bamako", "GMT+00:00 Africa/Banjul", "GMT+00:00 Africa/Bissau", "GMT+00:00 Africa/Conakry", "GMT+00:00 Africa/Dakar", "GMT+00:00 Africa/Freetown", "GMT+00:00 Africa/Lome", "GMT+00:00 Africa/Monrovia", "GMT+00:00 Africa/Nouakchott", "GMT+00:00 Africa/Ouagadougou", "GMT+00:00 Africa/Sao_Tome", "GMT+00:00 America/Danmarkshavn", "GMT+00:00 Antarctica/Troll", "GMT+00:00 Atlantic/Canary", "GMT+00:00 Atlantic/Faroe", "GMT+00:00 Atlantic/Madeira", "GMT+00:00 Atlantic/Reykjavik", "GMT+00:00 Atlantic/St_Helena", "GMT+00:00 Etc/GMT", "GMT+00:00 Etc/GMT+0", "GMT+00:00 Etc/GMT-0", "GMT+00:00 Etc/GMT0", "GMT+00:00 Etc/Greenwich", "GMT+00:00 Etc/Universal", "GMT+00:00 Etc/Zulu", "GMT+00:00 Europe/Dublin", "GMT+00:00 Europe/Guernsey", "GMT+00:00 Europe/Isle_of_Man", "GMT+00:00 Europe/Jersey", "GMT+00:00 Europe/Lisbon", "GMT+00:00 Europe/London", "GMT+00:00 GMT", "GMT+00:00 UTC", "GMT+00:00 WET", "GMT+01:00 Africa/Algiers", "GMT+01:00 Africa/Bangui", "GMT+01:00 Africa/Brazzaville", "GMT+01:00 Africa/Casablanca", "GMT+01:00 Africa/Ceuta", "GMT+01:00 Africa/Douala", "GMT+01:00 Africa/El_Aaiun", "GMT+01:00 Africa/Kinshasa", "GMT+01:00 Africa/Lagos", "GMT+01:00 Africa/Libreville", "GMT+01:00 Africa/Luanda", "GMT+01:00 Africa/Malabo", "GMT+01:00 Africa/Ndjamena", "GMT+01:00 Africa/Niamey", "GMT+01:00 Africa/Porto-Novo", "GMT+01:00 Africa/Tunis", "GMT+01:00 Arctic/Longyearbyen", "GMT+01:00 CET", "GMT+01:00 Etc/GMT-1", "GMT+01:00 Europe/Amsterdam", "GMT+01:00 Europe/Andorra", "GMT+01:00 Europe/Belgrade", "GMT+01:00 Europe/Berlin", "GMT+01:00 Europe/Bratislava", "GMT+01:00 Europe/Brussels", "GMT+01:00 Europe/Budapest", "GMT+01:00 Europe/Busingen", "GMT+01:00 Europe/Copenhagen", "GMT+01:00 Europe/Gibraltar", "GMT+01:00 Europe/Ljubljana", "GMT+01:00 Europe/Luxembourg", "GMT+01:00 Europe/Madrid", "GMT+01:00 Europe/Malta", "GMT+01:00 Europe/Monaco", "GMT+01:00 Europe/Oslo", "GMT+01:00 Europe/Paris", "GMT+01:00 Europe/Podgorica", "GMT+01:00 Europe/Prague", "GMT+01:00 Europe/Rome", "GMT+01:00 Europe/San_Marino", "GMT+01:00 Europe/Sarajevo", "GMT+01:00 Europe/Skopje", "GMT+01:00 Europe/Stockholm", "GMT+01:00 Europe/Tirane", "GMT+01:00 Europe/Vaduz", "GMT+01:00 Europe/Vatican", "GMT+01:00 Europe/Vienna", "GMT+01:00 Europe/Warsaw", "GMT+01:00 Europe/Zagreb", "GMT+01:00 Europe/Zurich", "GMT+01:00 MET", "GMT+02:00 Africa/Blantyre", "GMT+02:00 Africa/Bujumbura", "GMT+02:00 Africa/Cairo", "GMT+02:00 Africa/Gaborone", "GMT+02:00 Africa/Harare", "GMT+02:00 Africa/Johannesburg", "GMT+02:00 Africa/Khartoum", "GMT+02:00 Africa/Kigali", "GMT+02:00 Africa/Lubumbashi", "GMT+02:00 Africa/Lusaka", "GMT+02:00 Africa/Maputo", "GMT+02:00 Africa/Maseru", "GMT+02:00 Africa/Mbabane", "GMT+02:00 Africa/Tripoli", "GMT+02:00 Africa/Windhoek", "GMT+02:00 Asia/Amman", "GMT+02:00 Asia/Beirut", "GMT+02:00 Asia/Damascus", "GMT+02:00 Asia/Famagusta", "GMT+02:00 Asia/Gaza", "GMT+02:00 Asia/Hebron", "GMT+02:00 Asia/Jerusalem", "GMT+02:00 Asia/Nicosia", "GMT+02:00 EET", "GMT+02:00 Etc/GMT-2", "GMT+02:00 Europe/Athens", "GMT+02:00 Europe/Bucharest", "GMT+02:00 Europe/Chisinau", "GMT+02:00 Europe/Helsinki", "GMT+02:00 Europe/Kaliningrad", "GMT+02:00 Europe/Kiev", "GMT+02:00 Europe/Mariehamn", "GMT+02:00 Europe/Nicosia", "GMT+02:00 Europe/Riga", "GMT+02:00 Europe/Sofia", "GMT+02:00 Europe/Tallinn", "GMT+02:00 Europe/Uzhgorod", "GMT+02:00 Europe/Vilnius", "GMT+02:00 Europe/Zaporozhye", "GMT+03:00 Africa/Addis_Ababa", "GMT+03:00 Africa/Asmara", "GMT+03:00 Africa/Dar_es_Salaam", "GMT+03:00 Africa/Djibouti", "GMT+03:00 Africa/Juba", "GMT+03:00 Africa/Kampala", "GMT+03:00 Africa/Mogadishu", "GMT+03:00 Africa/Nairobi", "GMT+03:00 Antarctica/Syowa", "GMT+03:00 Asia/Aden", "GMT+03:00 Asia/Baghdad", "GMT+03:00 Asia/Bahrain", "GMT+03:00 Asia/Istanbul", "GMT+03:00 Asia/Kuwait", "GMT+03:00 Asia/Qatar", "GMT+03:00 Asia/Riyadh", "GMT+03:00 Etc/GMT-3", "GMT+03:00 Europe/Istanbul", "GMT+03:00 Europe/Kirov", "GMT+03:00 Europe/Minsk", "GMT+03:00 Europe/Moscow", "GMT+03:00 Europe/Simferopol", "GMT+03:00 Indian/Antananarivo", "GMT+03:00 Indian/Comoro", "GMT+03:00 Indian/Mayotte", "GMT+03:30 Asia/Tehran", "GMT+04:00 Asia/Baku", "GMT+04:00 Asia/Dubai", "GMT+04:00 Asia/Muscat", "GMT+04:00 Asia/Tbilisi", "GMT+04:00 Asia/Yerevan", "GMT+04:00 Etc/GMT-4", "GMT+04:00 Europe/Astrakhan", "GMT+04:00 Europe/Samara", "GMT+04:00 Europe/Saratov", "GMT+04:00 Europe/Ulyanovsk", "GMT+04:00 Europe/Volgograd", "GMT+04:00 Indian/Mahe", "GMT+04:00 Indian/Mauritius", "GMT+04:00 Indian/Reunion", "GMT+04:30 Asia/Kabul", "GMT+05:00 Antarctica/Mawson", "GMT+05:00 Asia/Aqtau", "GMT+05:00 Asia/Aqtobe", "GMT+05:00 Asia/Ashgabat", "GMT+05:00 Asia/Atyrau", "GMT+05:00 Asia/Dushanbe", "GMT+05:00 Asia/Karachi", "GMT+05:00 Asia/Oral", "GMT+05:00 Asia/Qyzylorda", "GMT+05:00 Asia/Samarkand", "GMT+05:00 Asia/Tashkent", "GMT+05:00 Asia/Yekaterinburg", "GMT+05:00 Etc/GMT-5", "GMT+05:00 Indian/Kerguelen", "GMT+05:00 Indian/Maldives", "GMT+05:30 Asia/Calcutta", "GMT+05:30 Asia/Colombo", "GMT+05:30 Asia/Kolkata", "GMT+05:45 Asia/Kathmandu", "GMT+05:45 Asia/Katmandu", "GMT+06:00 Antarctica/Vostok", "GMT+06:00 Asia/Almaty", "GMT+06:00 Asia/Bishkek", "GMT+06:00 Asia/Dhaka", "GMT+06:00 Asia/Omsk", "GMT+06:00 Asia/Qostanay", "GMT+06:00 Asia/Thimphu", "GMT+06:00 Asia/Urumqi", "GMT+06:00 Etc/GMT-6", "GMT+06:00 Indian/Chagos", "GMT+06:30 Asia/Yangon", "GMT+06:30 Indian/Cocos", "GMT+07:00 Antarctica/Davis", "GMT+07:00 Asia/Bangkok", "GMT+07:00 Asia/Barnaul", "GMT+07:00 Asia/Ho_Chi_Minh", "GMT+07:00 Asia/Hovd", "GMT+07:00 Asia/Jakarta", "GMT+07:00 Asia/Krasnoyarsk", "GMT+07:00 Asia/Novokuznetsk", "GMT+07:00 Asia/Novosibirsk", "GMT+07:00 Asia/Phnom Penh", "GMT+07:00 Asia/Pontianak", "GMT+07:00 Asia/Tomsk", "GMT+07:00 Asia/Vientiane", "GMT+07:00 Etc/GMT-7", "GMT+07:00 Indian/Christmas", "GMT+08:00 Asia/Brunei", "GMT+08:00 Asia/Choibalsan", "GMT+08:00 Asia/Hong_Kong", "GMT+08:00 Asia/Irkutsk", "GMT+08:00 Asia/Kuala_Lumpur", "GMT+08:00 Asia/Kuching", "GMT+08:00 Asia/Macau", "GMT+08:00 Asia/Makassar", "GMT+08:00 Asia/Manila", "GMT+08:00 Asia/Shanghai", "GMT+08:00 Asia/Singapore", "GMT+08:00 Asia/Taipei", "GMT+08:00 Asia/Ulaanbaatar", "GMT+08:00 Australia/Perth", "GMT+08:00 Etc/GMT-8", "GMT+08:45 Australia/Eucla", "GMT+09:00 Asia/Chita", "GMT+09:00 Asia/Dili", "GMT+09:00 Asia/Jayapura", "GMT+09:00 Asia/Khandyga", "GMT+09:00 Asia/Pyongyang", "GMT+09:00 Asia/Seoul", "GMT+09:00 Asia/Tokyo", "GMT+09:00 Asia/Yakutsk", "GMT+09:00 Etc/GMT-9", "GMT+09:00 Pacific/Palau", "GMT+09:30 Australia/Darwin", "GMT+10:00 Antarctica/DumontDUrville", "GMT+10:00 Asia/Ust-Nera", "GMT+10:00 Asia/Vladivostok", "GMT+10:00 Australia/Brisbane", "GMT+10:00 Australia/Lindeman", "GMT+10:00 Etc/GMT-10", "GMT+10:00 Pacific/Chuuk", "GMT+10:00 Pacific/Guam", "GMT+10:00 Pacific/Port_Moresby", "GMT+10:00 Pacific/Saipan", "GMT+10:30 Australia/Adelaide", "GMT+10:30 Australia/Broken_Hill", "GMT+11:00 Antarctica/Casey", "GMT+11:00 Antarctica/Macquarie", "GMT+11:00 Asia/Magadan", "GMT+11:00 Asia/Sakhalin", "GMT+11:00 Asia/Srednekolymsk", "GMT+11:00 Australia/Currie", "GMT+11:00 Australia/Hobart", "GMT+11:00 Australia/Lord_Howe", "GMT+11:00 Australia/Melbourne", "GMT+11:00 Australia/Sydney", "GMT+11:00 Etc/GMT-11", "GMT+11:00 Pacific/Bougainville", "GMT+11:00 Pacific/Efate", "GMT+11:00 Pacific/Guadalcanal", "GMT+11:00 Pacific/Kosrae", "GMT+11:00 Pacific/Noumea", "GMT+11:00 Pacific/Pohnpei", "GMT+12:00 Asia/Anadyr", "GMT+12:00 Asia/Kamchatka", "GMT+12:00 Etc/GMT-12", "GMT+12:00 Pacific/Fiji", "GMT+12:00 Pacific/Funafuti", "GMT+12:00 Pacific/Kwajalein", "GMT+12:00 Pacific/Majuro", "GMT+12:00 Pacific/Nauru", "GMT+12:00 Pacific/Norfolk", "GMT+12:00 Pacific/Tarawa", "GMT+12:00 Pacific/Wake", "GMT+12:00 Pacific/Wallis", "GMT+13:00 Antarctica/McMurdo", "GMT+13:00 Etc/GMT-13", "GMT+13:00 Pacific/Auckland", "GMT+13:00 Pacific/Enderbury", "GMT+13:00 Pacific/Fakaofo", "GMT+13:00 Pacific/Tongatapu", "GMT+13:45 Pacific/Chatham", "GMT+14:00 Etc/GMT-14", "GMT+14:00 Pacific/Apia", "GMT+14:00 Pacific/Kiritimati"],
		"dpi": {
			"Linux": ["default", "800x600", "1024x600", "1024x640", "1024x768", "1152x864", "1280x720", "1280x768", "1280x800", "1280x960", "1280x1024", "1360x768", "1366x768", "1400x1050", "1400x900", "1440x900", "1536x864", "1600x900", "1600x1200", "1680x1050", "1920x1080", "1920x1200", "2048x1152", "2304x1440", "2560x1440", "2560x1600", "2880x1800", "4096x2304", "5120x2880"],
			"MacOS": ["default", "2048x1152", "2304x1440", "2560x1440", "2560x1600", "2880x1800", "4096x2304", "5120x2880"],
			"Windows": ["default", "800x600", "1024x600", "1024x640", "1024x768", "1152x864", "1280x720", "1280x768", "1280x800", "1280x960", "1280x1024", "1360x768", "1366x768", "1400x1050", "1400x900", "1440x900", "1536x864", "1600x900", "1600x1200", "1680x1050", "1920x1080", "1920x1200", "2048x1152", "2304x1440", "2560x1440", "2560x1600", "2880x1800", "4096x2304", "5120x2880"]
		},
		"webgl": {
			"Linux": ["Google Inc. (AMD)", "Google Inc. (Intel)", "Google Inc. (NVIDIA)"],
			"MacOS": ["Google Inc. (Apple)", "Google Inc. (ATI Technologies Inc.)", "Google Inc. (NVIDIA)"],
			"Windows": ["Google Inc. (AMD)", "Google Inc. (Intel)", "Google Inc. (NVIDIA)"]
		},
		"cpu": {
			"Linux": [2, 4, 8, 12, 16, 20, 24],
			"MacOS": [8, 12, 16, 20],
			"Windows": [4, 8, 12, 16, 20, 24]
		},
		"mem": {
			"Linux": [2, 4, 8],
			"MacOS": [8],
			"Windows": [4, 8]
		},
		"region": {
			"AD": "安道尔",
			"AE": "阿联酋",
			"AF": "阿富汗",
			"AG": "安提瓜和巴布达",
			"AI": "安圭拉",
			"AL": "阿尔巴尼亚",
			"AM": "亚美尼亚",
			"AO": "安哥拉",
			"AR": "阿根廷",
			"AS": "美属萨摩亚",
			"AT": "奥地利",
			"AU": "澳大利亚",
			"AW": "阿鲁巴",
			"AX": "奥兰群岛",
			"AZ": "阿塞拜疆",
			"BA": "波黑",
			"BB": "巴巴多斯",
			"BD": "孟加拉",
			"BE": "比利时",
			"BF": "布基纳法索",
			"BG": "保加利亚",
			"BH": "巴林",
			"BI": "布隆迪",
			"BJ": "贝宁",
			"BM": "百慕大",
			"BN": "文莱",
			"BO": "玻利维亚",
			"BQ": "荷兰加勒比区",
			"BR": "巴西",
			"BS": "巴哈马",
			"BT": "不丹",
			"BW": "博茨瓦纳",
			"BY": "白俄罗斯",
			"BZ": "伯利兹",
			"CA": "加拿大",
			"CD": "刚果金",
			"CF": "中非",
			"CG": "刚果布",
			"CH": "瑞士",
			"CI": "科特迪瓦",
			"CL": "智利",
			"CM": "喀麦隆",
			"CN": "中国",
			"CO": "哥伦比亚",
			"CR": "哥斯达黎加",
			"CU": "古巴",
			"CV": "佛得角",
			"CW": "库拉索",
			"CY": "塞浦路斯",
			"CZ": "捷克",
			"DE": "德国",
			"DJ": "吉布提",
			"DK": "丹麦",
			"DM": "多米尼克",
			"DO": "多米尼加",
			"DZ": "阿尔及利亚",
			"EC": "厄瓜多尔",
			"EE": "爱沙尼亚",
			"EG": "埃及",
			"ER": "厄立特里亚",
			"ES": "西班牙",
			"ET": "埃塞俄比亚",
			"FI": "芬兰",
			"FJ": "斐济群岛",
			"FM": "密克罗尼西亚联邦",
			"FO": "法罗群岛",
			"FR": "法国",
			"GA": "加蓬",
			"GB": "英国",
			"GD": "格林纳达",
			"GE": "格鲁吉亚",
			"GF": "法属圭亚那",
			"GG": "根西岛",
			"GH": "加纳",
			"GI": "直布罗陀",
			"GL": "格陵兰",
			"GM": "冈比亚",
			"GN": "几内亚",
			"GP": "瓜德罗普",
			"GQ": "赤道几内亚",
			"GR": "希腊",
			"GT": "危地马拉",
			"GU": "关岛",
			"GW": "几内亚比绍",
			"GY": "圭亚那",
			"HK": "中国香港",
			"HN": "洪都拉斯",
			"HR": "克罗地亚",
			"HT": "海地",
			"HU": "匈牙利",
			"ID": "印尼",
			"IE": "爱尔兰",
			"IL": "以色列",
			"IM": "马恩岛",
			"IN": "印度",
			"IQ": "伊拉克",
			"IR": "伊朗",
			"IS": "冰岛",
			"IT": "意大利",
			"JE": "泽西岛",
			"JM": "牙买加",
			"JO": "约旦",
			"JP": "日本",
			"KE": "肯尼亚",
			"KG": "吉尔吉斯斯坦",
			"KH": "柬埔寨",
			"KI": "基里巴斯",
			"KM": "科摩罗",
			"KN": "圣基茨和尼维斯",
			"KR": "韩国",
			"KW": "科威特",
			"KY": "开曼群岛",
			"KZ": "哈萨克斯坦",
			"LA": "老挝",
			"LB": "黎巴嫩",
			"LC": "圣卢西亚",
			"LI": "列支敦士登",
			"LK": "斯里兰卡",
			"LR": "利比里亚",
			"LS": "莱索托",
			"LT": "立陶宛",
			"LU": "卢森堡",
			"LV": "拉脱维亚",
			"LY": "利比亚",
			"MA": "摩洛哥",
			"MC": "摩纳哥",
			"MD": "摩尔多瓦",
			"ME": "黑山",
			"MF": "法属圣马丁",
			"MG": "马达加斯加",
			"MH": "马绍尔群岛",
			"MK": "马其顿",
			"ML": "马里",
			"MM": "缅甸",
			"MN": "蒙古国蒙古",
			"MO": "中国澳门",
			"MP": "北马里亚纳群岛",
			"MQ": "马提尼克",
			"MR": "毛里塔尼亚",
			"MT": "马耳他",
			"MU": "毛里求斯",
			"MV": "马尔代夫",
			"MW": "马拉维",
			"MX": "墨西哥",
			"MY": "马来西亚",
			"MZ": "莫桑比克",
			"NA": "纳米比亚",
			"NC": "新喀里多尼亚",
			"NE": "尼日尔",
			"NG": "尼日利亚",
			"NI": "尼加拉瓜",
			"NL": "荷兰",
			"NO": "挪威",
			"NP": "尼泊尔",
			"NZ": "新西兰",
			"OM": "阿曼",
			"PA": "巴拿马",
			"PE": "秘鲁",
			"PF": "法属波利尼西亚",
			"PG": "巴布亚新几内亚",
			"PH": "菲律宾",
			"PK": "巴基斯坦",
			"PL": "波兰",
			"PM": "圣皮埃尔和密克隆",
			"PR": "波多黎各",
			"PS": "巴勒斯坦",
			"PT": "葡萄牙",
			"PW": "帕劳",
			"PY": "巴拉圭",
			"QA": "卡塔尔",
			"RE": "留尼汪",
			"RO": "罗马尼亚",
			"RS": "塞尔维亚",
			"RU": "俄罗斯",
			"RW": "卢旺达",
			"SA": "沙特阿拉伯",
			"SB": "所罗门群岛",
			"SC": "塞舌尔",
			"SD": "苏丹",
			"SE": "瑞典",
			"SG": "新加坡",
			"SI": "斯洛文尼亚",
			"SK": "斯洛伐克",
			"SL": "塞拉利昂",
			"SM": "圣马力诺",
			"SN": "塞内加尔",
			"SO": "索马里",
			"SR": "苏里南",
			"SS": "南苏丹",
			"ST": "圣多美和普林西比",
			"SV": "萨尔瓦多",
			"SX": "荷属圣马丁",
			"SY": "叙利亚",
			"SZ": "斯威士兰",
			"TC": "特克斯和凯科斯群岛",
			"TD": "乍得",
			"TG": "多哥",
			"TH": "泰国",
			"TJ": "塔吉克斯坦",
			"TL": "东帝汶",
			"TM": "土库曼斯坦",
			"TN": "突尼斯",
			"TO": "汤加",
			"TR": "土耳其",
			"TT": "特立尼达和多巴哥",
			"TV": "图瓦卢",
			"TW": "中国台湾",
			"TZ": "坦桑尼亚",
			"UA": "乌克兰",
			"UG": "乌干达",
			"US": "美国",
			"UY": "乌拉圭",
			"UZ": "乌兹别克斯坦",
			"VA": "梵蒂冈",
			"VC": "圣文森特和格林纳丁斯",
			"VE": "委内瑞拉",
			"VG": "英属维尔京群岛",
			"VI": "美属维尔京群岛",
			"VN": "越南",
			"VU": "瓦努阿图",
			"WF": "瓦利斯和富图纳",
			"WS": "萨摩亚",
			"XK": "科索沃",
			"YE": "也门",
			"YT": "马约特",
			"ZA": "南非",
			"ZM": "赞比亚",
			"ZW": "津巴布韦"
		}
	}
}`

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}

	httpClient := &http.Client{Transport: mockTransport}
	client, err := NewClient("test-key", WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	resp, err := client.GetUiFingerList(context.Background())
	if err != nil {
		t.Fatalf("GetEnvPage() error = %v", err)
	}

	if len(resp.ChromeKernelversion) <= 0 {
		t.Errorf("GetUiFingerList() chromeKernelversion = %v, want > 0", resp.ChromeKernelversion)
	}

}
