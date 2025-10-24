package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/automationpi/azdocs/pkg/auth"
	"github.com/automationpi/azdocs/pkg/cache"
)

// Config holds discovery configuration
type Config struct {
	SubscriptionID string
	TenantID       string
	Concurrency    int
	ShowProgress   bool
}

// Client handles Azure resource discovery
type Client struct {
	auth   *auth.AzureAuthenticator
	cache  *cache.Cache
	config Config
}

// NewDiscoveryClient creates a new discovery client
func NewDiscoveryClient(auth *auth.AzureAuthenticator, cache *cache.Cache, config Config) *Client {
	return &Client{
		auth:   auth,
		cache:  cache,
		config: config,
	}
}

// Result contains discovery results
type Result struct {
	SubscriptionID string                 `json:"subscriptionId"`
	Timestamp      string                 `json:"timestamp"`
	Stats          Stats                  `json:"stats"`
	RawData        map[string]interface{} `json:"rawData,omitempty"`
}

// Stats contains resource statistics
type Stats struct {
	TotalResources int `json:"totalResources"`
	VNets          int `json:"vnets"`
	Subnets        int `json:"subnets"`
	NSGs           int `json:"nsgs"`
	RouteTables    int `json:"routeTables"`
}

// Discover performs resource discovery
func (c *Client) Discover(ctx context.Context) (*Result, error) {
	// Use Azure Resource Graph for fast discovery
	result := &Result{
		SubscriptionID: c.config.SubscriptionID,
		Timestamp:      fmt.Sprintf("%v", time.Now().UTC().Format(time.RFC3339)),
		Stats:          Stats{},
		RawData:        make(map[string]interface{}),
	}

	// Query all resources using Resource Graph
	resources, err := c.queryResourceGraph(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query Resource Graph: %w", err)
	}

	// Count resources by type
	result.Stats.TotalResources = len(resources)

	for _, res := range resources {
		resType, ok := res["type"].(string)
		if !ok {
			continue
		}

		// Resource Graph returns lowercase types, normalize them
		switch {
		case resType == "microsoft.network/virtualnetworks" || resType == "Microsoft.Network/virtualNetworks":
			result.Stats.VNets++
		case resType == "microsoft.network/virtualnetworks/subnets" || resType == "Microsoft.Network/virtualNetworks/subnets":
			result.Stats.Subnets++
		case resType == "microsoft.network/networksecuritygroups" || resType == "Microsoft.Network/networkSecurityGroups":
			result.Stats.NSGs++
		case resType == "microsoft.network/routetables" || resType == "Microsoft.Network/routeTables":
			result.Stats.RouteTables++
		}
	}

	result.RawData["resources"] = resources

	return result, nil
}

// queryResourceGraph queries Azure Resource Graph for all resources
func (c *Client) queryResourceGraph(ctx context.Context) ([]map[string]interface{}, error) {
	return c.getAllResources(ctx)
}

// SaveToDirectory saves discovery results to a directory
func (r *Result) SaveToDirectory(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save raw data
	rawDir := filepath.Join(dir, "raw")
	if err := os.MkdirAll(rawDir, 0755); err != nil {
		return err
	}

	// Save all resources to raw/all-resources.json
	if resources, ok := r.RawData["resources"]; ok {
		resourcesPath := filepath.Join(rawDir, "all-resources.json")
		resourcesData, err := json.MarshalIndent(resources, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal resources: %w", err)
		}
		if err := os.WriteFile(resourcesPath, resourcesData, 0644); err != nil {
			return fmt.Errorf("failed to write resources: %w", err)
		}
	}

	// Save metadata
	metaPath := filepath.Join(dir, "metadata.json")
	metaData, err := json.MarshalIndent(map[string]interface{}{
		"subscriptionId": r.SubscriptionID,
		"timestamp":      r.Timestamp,
		"stats":          r.Stats,
	}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(metaPath, metaData, 0644)
}
