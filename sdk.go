package brosdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Package brosdk provides a Go SDK for interacting with the Browser Open API service.
//
// The SDK offers:
//   - Client initialization with API key authentication
//   - Configurable HTTP client with timeout support
//   - Optional debug mode for logging requests and responses
//   - Type-safe request/response structures
//   - Methods for user signature management and environment operations
//
// Basic usage:
//
//	client, err := brosdk.NewClient("your-api-key")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Enable debug logging
//	client, err := brosdk.NewClient("your-api-key", brosdk.WithDebug(true))
//
//	// Get user signature
//	sigResp, err := client.GetUserSig(context.Background(), &brosdk.GetUserSigRequest{
//	    CustomerId: "customer123",
//	    Duration:   3600,
//	})
const (
	// DefaultEndpoint is the default API endpoint
	DefaultEndpoint = "https://api.brosdk.com"
)

// Client represents the Browser Open SDK client
type Client struct {
	Endpoint   string
	ApiKey     string
	httpClient *http.Client
	debug      bool
	logger     *log.Logger
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

// WithDebug enables or disables debug logging.
// When enabled, each HTTP request and response will be printed to stderr,
// including method, URL, status code, business code, and elapsed time.
//
// Example:
//
//	client, _ := brosdk.NewClient("key", brosdk.WithDebug(true))
func WithDebug(enabled bool) ClientOption {
	return func(c *Client) {
		c.debug = enabled
	}
}

// WithLogger sets a custom logger for debug output.
// This option only takes effect when debug mode is enabled via WithDebug(true).
// If not set, the SDK logs to stderr with a "[brosdk]" prefix.
func WithLogger(logger *log.Logger) ClientOption {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
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
		logger: log.New(os.Stderr, "[brosdk] ", log.LstdFlags),
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// debugLog prints a debug message if debug mode is enabled
func (c *Client) debugLog(format string, args ...interface{}) {
	if c.debug {
		c.logger.Printf(format, args...)
	}
}

// newRequest creates a new HTTP request with proper headers
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, []byte, error) {
	url := c.Endpoint + path

	var req *http.Request
	var err error
	var rawBody []byte

	if body != nil {
		rawBody, err = json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(rawBody))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	// Set common headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "browser-open-sdk/1.0")

	return req, rawBody, nil
}

// doWithDebug executes the HTTP request, handles debug logging, and returns the response body bytes
func (c *Client) doWithDebug(method, path string, req *http.Request, reqBody []byte) ([]byte, int, error) {
	start := time.Now()

	if c.debug {
		if len(reqBody) > 0 {
			c.debugLog("→ %s %s  body=%s", method, path, string(reqBody))
		} else {
			c.debugLog("→ %s %s", method, path)
		}
	}

	resp, err := c.httpClient.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		if c.debug {
			c.debugLog("✗ %s %s  elapsed=%s  error=%v", method, path, elapsed, err)
		}
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if c.debug {
			c.debugLog("✗ %s %s  elapsed=%s  http_status=%d", method, path, elapsed, resp.StatusCode)
		}
		return nil, resp.StatusCode, fmt.Errorf("non-OK HTTP status: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		if c.debug {
			c.debugLog("✗ %s %s  elapsed=%s  read_error=%v", method, path, elapsed, err)
		}
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	if c.debug {
		c.debugLog("← %s %s  elapsed=%s  http_status=%d  body=%s", method, path, elapsed, resp.StatusCode, string(respBody))
	}

	return respBody, resp.StatusCode, nil
}

// GetUserSig retrieves user signature with the specified parameters
func (c *Client) GetUserSig(ctx context.Context, req *GetUserSigRequest) (*UserSigData, error) {
	const method, path = "POST", "/api/v2/browser/getUserSig"

	httpReq, reqBody, err := c.newRequest(ctx, method, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetUserSig request: %w", err)
	}

	respBody, _, err := c.doWithDebug(method, path, httpReq, reqBody)
	if err != nil {
		return nil, fmt.Errorf("GetUserSig request failed: %w", err)
	}

	var result GetUserSigResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode GetUserSig response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("GetUserSig failed (code=%d): %s", result.Code, result.Msg)
	}

	return &result.Data, nil
}

// EnvCreate creates a new environment with the specified parameters
func (c *Client) EnvCreate(ctx context.Context, req *CreateEnv) (*EnvInfo, error) {
	const method, path = "POST", "/api/v2/browser/create"

	httpReq, reqBody, err := c.newRequest(ctx, method, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create EnvCreate request: %w", err)
	}

	respBody, _, err := c.doWithDebug(method, path, httpReq, reqBody)
	if err != nil {
		return nil, fmt.Errorf("EnvCreate request failed: %w", err)
	}

	var result EnvResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode EnvCreate response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("EnvCreate failed (code=%d): %s", result.Code, result.Msg)
	}

	return &result.Data, nil
}

// EnvUpdate updates an existing browser environment
func (c *Client) EnvUpdate(ctx context.Context, req *UpdateEnv) (*EnvInfo, error) {
	const method, path = "POST", "/api/v2/browser/update"

	httpReq, reqBody, err := c.newRequest(ctx, method, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create EnvUpdate request: %w", err)
	}

	respBody, _, err := c.doWithDebug(method, path, httpReq, reqBody)
	if err != nil {
		return nil, fmt.Errorf("EnvUpdate request failed: %w", err)
	}

	var result EnvResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode EnvUpdate response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("EnvUpdate failed (code=%d): %s", result.Code, result.Msg)
	}

	return &result.Data, nil
}

// EnvUpdate updates an existing browser environment with a new name
func (c *Client) EnvUpdateEnvMeta(ctx context.Context, req *UpdateEnvMeta) (*EnvInfo, error) {
	const method, path = "POST", "/api/v2/browser/updateEnv"

	httpReq, reqBody, err := c.newRequest(ctx, method, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create EnvUpdate request: %w", err)
	}

	respBody, _, err := c.doWithDebug(method, path, httpReq, reqBody)
	if err != nil {
		return nil, fmt.Errorf("EnvUpdate request failed: %w", err)
	}

	var result EnvResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode EnvUpdate response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("EnvUpdate failed (code=%d): %s", result.Code, result.Msg)
	}

	return &result.Data, nil
}

// EnvDestroy deletes a browser environment
func (c *Client) EnvDestroy(ctx context.Context, req *EnvDelReq) error {
	const method, path = "POST", "/api/v2/browser/destroy"

	httpReq, reqBody, err := c.newRequest(ctx, method, path, req)
	if err != nil {
		return fmt.Errorf("failed to create EnvDestroy request: %w", err)
	}

	respBody, _, err := c.doWithDebug(method, path, httpReq, reqBody)
	if err != nil {
		return fmt.Errorf("EnvDestroy request failed: %w", err)
	}

	var result Response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("failed to decode EnvDestroy response: %w", err)
	}

	if result.Code != 200 {
		return fmt.Errorf("EnvDestroy failed (code=%d): %s", result.Code, result.Msg)
	}

	return nil
}

// GetEnvPage gets paginated browser environments
func (c *Client) GetEnvPage(ctx context.Context, req *GetEnvPageReq) (*Page, error) {
	const method, path = "POST", "/api/v2/browser/page"

	httpReq, reqBody, err := c.newRequest(ctx, method, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEnvPage request: %w", err)
	}

	respBody, _, err := c.doWithDebug(method, path, httpReq, reqBody)
	if err != nil {
		return nil, fmt.Errorf("GetEnvPage request failed: %w", err)
	}

	var result PageResp
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode GetEnvPage response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("GetEnvPage failed (code=%d): %s", result.Code, result.Msg)
	}

	return &result.Data, nil
}

// GetUiFingerList retrieves the list of available fingerprint parameters
func (c *Client) GetUiFingerList(ctx context.Context) (*GetUiFingerList, error) {
	const method, path = "GET", "/api/v2/browser/getUiFingerList"

	httpReq, _, err := c.newRequest(ctx, method, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetUiFingerList request: %w", err)
	}

	respBody, _, err := c.doWithDebug(method, path, httpReq, nil)
	if err != nil {
		return nil, fmt.Errorf("GetUiFingerList request failed: %w", err)
	}

	var result GetUiFingerListResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode GetUiFingerList response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("GetUiFingerList failed (code=%d): %s", result.Code, result.Msg)
	}

	return &result.Data, nil
}
