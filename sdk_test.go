package brosdk

import (
	"context"
	"encoding/json"
	"io"
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

	req, err := client.newRequest(ctx, "POST", "/test", body)
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

	// Test body content
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

	req, err := client.newRequest(context.Background(), "GET", "/test", nil)
	if err != nil {
		t.Fatalf("newRequest() error = %v", err)
	}

	if req.Body != nil {
		t.Error("newRequest() with nil body should have nil Body")
	}
}

func TestClient_GetUserSig_Success(t *testing.T) {
	// Mock HTTP client
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request
			if req.URL.Path != "/api/usersig" {
				t.Errorf("Expected path /api/usersig, got %s", req.URL.Path)
			}

			// Return mock response
			responseBody := `{
				"code": 0,
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

	if !strings.Contains(err.Error(), "status: 400") {
		t.Errorf("GetUserSig() error = %v, want error containing 'status: 400'", err)
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
			if req.URL.Path != "/api/env" {
				t.Errorf("Expected path /api/env, got %s", req.URL.Path)
			}

			responseBody := `{
				"code": 0,
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
		UserAgent:  "test-user-agent",
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
				"code": 0,
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
				"code": 0,
				"data": [
					{
						"envId": "1",
						"customerId": "customer1",
						"envName": "Environment 1",
						"createdAt": "2023-01-01T00:00:00Z",
						"updatedAt": "2023-01-01T00:00:00Z"
					},
					{
						"envId": "2",
						"customerId": "customer1",
						"envName": "Environment 2",
						"createdAt": "2023-01-02T00:00:00Z",
						"updatedAt": "2023-01-02T00:00:00Z"
					}
				],
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
