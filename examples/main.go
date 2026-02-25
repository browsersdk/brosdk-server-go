package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/browsersdk/brosdk-server-go"
)

var apiKey = "5Ij4QwXjzEGsMCtEmxUs6hI4nectJeeYhhkdchpMZD0cgmGnQLtvQoLXoVZJ1TQg"

func main() {

	// Initialize client
	client, err := brosdk.NewClient(apiKey, brosdk.WithEndpoint("http://192.168.0.188:9988"))
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	fmt.Println("=== Browser Open SDK Examples ===")

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
		Duration:   3600 * 24, // 1 hour
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetUserSig(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get user signature: %w", err)
	}

	fmt.Printf("   ✓ Success! UserSig: %s\n", resp.UserSig)
	fmt.Printf("   ✓ Expires at: %d\n", resp.ExpireTime)

	return nil
}

func createEnvironmentExample(client *brosdk.Client) error {
	fmt.Println("2. Creating Browser Environment...")

	req := &brosdk.EnvInfo{
		CustomerId:      "demo-customer",
		EnvName:         "Demo Browser Environment",
		UserAgent:       "",
		System:          "Windows 10",
		Kernel:          "Chrome",
		KernelVersion:   "134",
		Dpi:             "96",
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
			UseIP:     1,
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

	fmt.Printf("   ✓ Success! Environment ID: %d\n", resp.EnvId)
	fmt.Printf("   ✓ Environment Name: %s\n", resp.EnvName)

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
	fmt.Printf("   ✓ Retrieved %d environments\n", len(resp.List))

	if len(resp.List) > 0 {
		fmt.Println("   Environment Details:")
		for i, env := range resp.List {
			fmt.Printf("     %d. ID: %d, Name: %s\n",
				i+1, env.EnvId, env.EnvName)
		}
	}

	return nil
}
