package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/browsersdk/brosdk-server-go"
)

func main() {
	// Initialize client
	client, err := brosdk.NewClient("your-api-key-here")
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	fmt.Println("=== Browser Open SDK Examples ===\n")

	// Example 1: Get User Signature
	if err := getUserSignatureExample(client); err != nil {
		log.Printf("GetUserSig example failed: %v", err)
	}

	// Example 2: Create Environment
	if err := createEnvironmentExample(client); err != nil {
		log.Printf("EnvCreate example failed: %v", err)
	}

	// Example 3: List Environments
	if err := listEnvironmentsExample(client); err != nil {
		log.Printf("GetEnvPage example failed: %v", err)
	}
}

func getUserSignatureExample(client *brosdk.Client) error {
	fmt.Println("1. Getting User Signature...")
	
	req := &brosdk.GetUserSigRequest{
		CustomerId: "demo-customer",
		Duration:   3600, // 1 hour
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetUserSig(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get user signature: %w", err)
	}

	fmt.Printf("   ✓ Success! UserSig: %s\n", resp.Data.UserSig)
	fmt.Printf("   ✓ Expires at: %d\n", resp.Data.ExpireTime)
	fmt.Printf("   ✓ Request ID: %s\n\n", resp.ReqId)
	
	return nil
}

func createEnvironmentExample(client *brosdk.Client) error {
	fmt.Println("2. Creating Browser Environment...")
	
	req := &brosdk.EnvCreateRequest{
		CustomerId:      "demo-customer",
		EnvName:         "Demo Browser Environment",
		UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		System:          "Windows 10",
		Kernel:          "Chrome",
		KernelVersion:   "120.0.0.0",
		DPI:             "96",
		DeviceName:      "Demo PC",
		Mac:             "00:11:22:33:44:55",
		PublicIp:        "192.168.1.100",
		Zone:            "UTC",
		EnableCookie:    1,
		Enablenotice:    1,
		Enableopen:      1,
		Enablepic:       1,
		IgnoreCookieErr: 0,
		Geographic: brosdk.Geographic{
			Enable:    1,
			Latitude:  "39.9042",
			Longitude: "116.4074",
			Accuracy:  "high",
			Useip:     1,
		},
		FontList: []string{
			"Arial", "Helvetica", "Times New Roman", "Courier New",
		},
		Language: []string{"en-US", "zh-CN"},
		ScanPort: []int{80, 443, 8080},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := client.EnvCreate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}

	fmt.Printf("   ✓ Success! Environment ID: %d\n", resp.Data.EnvId)
	fmt.Printf("   ✓ Environment Name: %s\n", resp.Data.EnvName)
	fmt.Printf("   ✓ Request ID: %s\n\n", resp.ReqId)
	
	return nil
}

func listEnvironmentsExample(client *brosdk.Client) error {
	fmt.Println("3. Listing Browser Environments...")
	
	req := &brosdk.GetEnvPageReq{
		ReqPage: brosdk.ReqPage{
			Page:     1,
			PageSize: 10,
		},
		CustomerId: "demo-customer",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetEnvPage(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to list environments: %w", err)
	}

	fmt.Printf("   ✓ Success! Total environments: %d\n", resp.Total)
	fmt.Printf("   ✓ Retrieved %d environments\n", len(resp.Data))
	fmt.Printf("   ✓ Request ID: %s\n\n", resp.ReqId)
	
	if len(resp.Data) > 0 {
		fmt.Println("   Environment Details:")
		for i, env := range resp.Data {
			fmt.Printf("     %d. ID: %d, Name: %s, Created: %s\n", 
				i+1, env.EnvId, env.EnvName, env.CreatedAt)
		}
	}
	
	return nil
}