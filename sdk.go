package brosdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Package sdk provides a Go SDK for interacting with the Browser Open API service.
//
// The SDK offers:
// - Client initialization with API key authentication
// - Configurable HTTP client with timeout support
// - Type-safe request/response structures
// - Methods for user signature management and environment operations
//
// Basic usage:
//
//	client, err := sdk.NewClient("your-api-key")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Get user signature
//	sigReq := &sdk.GetUserSigRequest{
//	    CustomerId: "customer123",
//	    Duration:   3600,
//	}
//	sigResp, err := client.GetUserSig(context.Background(), sigReq)
const (
	// DefaultEndpoint is the default API endpoint
	DefaultEndpoint = "https://api.brosdk.com"
)

// Client represents the Browser Open SDK client
type Client struct {
	Endpoint   string
	ApiKey     string
	httpClient *http.Client
}

// ClientOption defines a function type for configuring the client
type ClientOption func(*Client)

// WithEndpoint sets a custom endpoint for the client
func WithEndpoint(endpoint string) ClientOption {
	return func(c *Client) {
		if endpoint != "" {
			c.Endpoint = endpoint
		}
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithTimeout sets a custom timeout for the HTTP client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{
				Timeout: timeout,
			}
		} else {
			c.httpClient.Timeout = timeout
		}
	}
}

// NewClient creates a new SDK client with the required ApiKey
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey is required")
	}

	client := &Client{
		Endpoint: DefaultEndpoint,
		ApiKey:   apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// newRequest creates a new HTTP request with proper headers
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.Endpoint + path

	var req *http.Request
	var err error

	if body != nil {
		jsonBody, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", marshalErr)
		}
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	// Set common headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "browser-open-sdk/1.0")

	return req, nil
}

// do executes the HTTP request and returns the response
func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK status code: %d", resp.StatusCode)
	}
	return resp, nil
}

// GetUserSig retrieves user signature with the specified parameters
func (c *Client) GetUserSig(ctx context.Context, req *GetUserSigRequest) (*UserSigData, error) {
	httpReq, err := c.newRequest(ctx, "POST", "/api/v2/browser/getUserSig", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetUserSig request: %w", err)
	}

	resp, err := c.do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("GetUserSig request failed: %w", err)
	}
	defer resp.Body.Close()

	var result GetUserSigResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode GetUserSig response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("failed to get user signature: %s", result.Msg)
	}

	return &result.Data, nil
}

// EnvCreate creates a new environment with the specified parameters
func (c *Client) EnvCreate(ctx context.Context, req *EnvInfo) (*EnvInfo, error) {
	httpReq, err := c.newRequest(ctx, "POST", "/api/v2/browser/create", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create EnvCreate request: %w", err)
	}

	resp, err := c.do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("EnvCreate request failed: %w", err)
	}
	defer resp.Body.Close()

	var result EnvResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode EnvCreate response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("failed to create environment: %s", result.Msg)
	}

	return &result.Data, nil
}

// EnvUpdate updates browser environment
func (c *Client) EnvUpdate(ctx context.Context, req *EnvInfo) (*EnvInfo, error) {
	httpReq, err := c.newRequest(ctx, "POST", "/api/v2/browser/update", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateV2 request: %w", err)
	}

	resp, err := c.do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("UpdateV2 request failed: %w", err)
	}
	defer resp.Body.Close()

	var result EnvResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode UpdateV2 response: %w", err)
	}

	fmt.Printf("UpdateV2 response: %+v\n", result)
	if result.Code != 200 {
		return nil, fmt.Errorf("failed to update environment: %s", result.Msg)
	}

	return &result.Data, nil
}

// EnvDestroy deletes browser environment
func (c *Client) EnvDestroy(ctx context.Context, req *EnvDelReq) error {
	httpReq, err := c.newRequest(ctx, "POST", "/api/v2/browser/destroy", req)
	if err != nil {
		return fmt.Errorf("failed to create Destroy request: %w", err)
	}

	resp, err := c.do(httpReq)
	if err != nil {
		return fmt.Errorf("Destroy request failed: %w", err)
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode Destroy response: %w", err)
	}

	if result.Code != 200 {
		return fmt.Errorf("Destroy request failed with code: %s", result.Msg)
	}

	return nil
}

// GetEnvPage gets paginated browser environments
func (c *Client) GetEnvPage(ctx context.Context, req *GetEnvPageReq) (*Page, error) {
	httpReq, err := c.newRequest(ctx, "POST", "/api/v2/browser/page", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create PageV2 request: %w", err)
	}

	resp, err := c.do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("PageV2 request failed: %w", err)
	}
	defer resp.Body.Close()

	var result PageResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode PageV2 response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("PageV2 request failed with code: %s", result.Msg)
	}

	return &result.Data, nil
}
