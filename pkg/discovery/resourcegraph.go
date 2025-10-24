package discovery

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph"
	"github.com/automationpi/azdocs/pkg/auth"
)

// QueryResourceGraph queries Azure Resource Graph for resources
func QueryResourceGraph(ctx context.Context, authClient *auth.AzureAuthenticator, subscriptionID string, query string) ([]map[string]interface{}, error) {
	// Create Resource Graph client
	client, err := armresourcegraph.NewClient(authClient.GetCredential(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Resource Graph client: %w", err)
	}

	// Build query request
	request := armresourcegraph.QueryRequest{
		Subscriptions: []*string{&subscriptionID},
		Query:         &query,
	}

	// Execute query
	response, err := client.Resources(ctx, request, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Parse results
	if response.Data == nil {
		return []map[string]interface{}{}, nil
	}

	// Convert response data to slice of maps
	var results []map[string]interface{}

	// The Data field is an interface{}, need to handle it properly
	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
	}

	return results, nil
}

// GetAllResources retrieves all resources in a subscription
func (c *Client) getAllResources(ctx context.Context) ([]map[string]interface{}, error) {
	query := "Resources | project id, name, type, location, resourceGroup, tags, properties"

	return QueryResourceGraph(ctx, c.auth, c.config.SubscriptionID, query)
}

// GetVNetDetails retrieves detailed VNet information including subnets
func (c *Client) GetVNetDetails(ctx context.Context) ([]map[string]interface{}, error) {
	query := `Resources
| where type == 'microsoft.network/virtualnetworks'
| project id, name, type, location, resourceGroup, tags, properties`

	return QueryResourceGraph(ctx, c.auth, c.config.SubscriptionID, query)
}

// GetVNetPeerings retrieves VNet peering information
func (c *Client) GetVNetPeerings(ctx context.Context) ([]map[string]interface{}, error) {
	query := `Resources
| where type == 'microsoft.network/virtualnetworks'
| mv-expand peering = properties.virtualNetworkPeerings
| project vnetId = id, vnetName = name, peeringName = peering.name,
  remoteVNetId = tostring(peering.properties.remoteVirtualNetwork.id),
  peeringState = tostring(peering.properties.peeringState)`

	return QueryResourceGraph(ctx, c.auth, c.config.SubscriptionID, query)
}
