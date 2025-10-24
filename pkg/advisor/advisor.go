package advisor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/automationpi/azdocs/pkg/auth"
)

// Recommendation represents an Azure Advisor recommendation
type Recommendation struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Category         string `json:"category"`
	Impact           string `json:"impact"`
	Risk             string `json:"risk"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	RecommendedAction string `json:"recommended_action"`
	ImpactedField    string `json:"impacted_field"`
	ImpactedValue    string `json:"impacted_value"`
	ResourceGroup    string `json:"resource_group"`
}

// Client handles Azure Advisor API interactions
type Client struct {
	subscriptionID string
	credential     *azidentity.DefaultAzureCredential
}

// NewClient creates a new Advisor client
func NewClient(subscriptionID string) (*Client, error) {
	cred, err := auth.GetCredential()
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure credentials: %w", err)
	}

	return &Client{
		subscriptionID: subscriptionID,
		credential:     cred,
	}, nil
}

// GetRecommendations fetches all Azure Advisor recommendations using Resource Graph
func (c *Client) GetRecommendations(ctx context.Context) ([]Recommendation, error) {
	// Use Azure Resource Graph to query Advisor recommendations
	query := `
	advisorresources
	| where type == "microsoft.advisor/recommendations"
	| project
		id,
		name,
		category = tostring(properties.category),
		impact = tostring(properties.impact),
		risk = tostring(properties.risk),
		shortDescription = tostring(properties.shortDescription.problem),
		recommendedAction = tostring(properties.shortDescription.solution),
		impactedField = tostring(properties.impactedField),
		impactedValue = tostring(properties.impactedValue),
		resourceGroup = tostring(properties.resourceMetadata.resourceId)
	`

	// Import the resource graph client
	// For now, we'll create a placeholder implementation
	// The actual implementation should use pkg/discovery/resourcegraph.go pattern

	var recommendations []Recommendation

	// TODO: Implement actual Resource Graph query
	// This would follow the same pattern as in pkg/discovery/resourcegraph.go
	// For now, return empty array - we'll implement full integration in next step

	return recommendations, nil
}

// ParseRecommendationResponse parses the Azure Advisor API response
func ParseRecommendationResponse(data []byte) ([]Recommendation, error) {
	var response struct {
		Value []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Properties struct {
				Category         string `json:"category"`
				Impact           string `json:"impact"`
				Risk             string `json:"risk"`
				ShortDescription struct {
					Problem  string `json:"problem"`
					Solution string `json:"solution"`
				} `json:"shortDescription"`
				ExtendedProperties map[string]string `json:"extendedProperties"`
				ImpactedField      string            `json:"impactedField"`
				ImpactedValue      string            `json:"impactedValue"`
				ResourceMetadata   struct {
					ResourceID string `json:"resourceId"`
				} `json:"resourceMetadata"`
			} `json:"properties"`
		} `json:"value"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse recommendations: %w", err)
	}

	var recommendations []Recommendation
	for _, item := range response.Value {
		rec := Recommendation{
			ID:                item.ID,
			Name:              item.Name,
			Category:          item.Properties.Category,
			Impact:            item.Properties.Impact,
			Risk:              item.Properties.Risk,
			ShortDescription:  item.Properties.ShortDescription.Problem,
			RecommendedAction: item.Properties.ShortDescription.Solution,
			ImpactedField:     item.Properties.ImpactedField,
			ImpactedValue:     item.Properties.ImpactedValue,
		}

		// Extract resource group from resource ID if available
		if item.Properties.ResourceMetadata.ResourceID != "" {
			// Parse: /subscriptions/{sub}/resourceGroups/{rg}/...
			rec.ResourceGroup = extractResourceGroup(item.Properties.ResourceMetadata.ResourceID)
		}

		recommendations = append(recommendations, rec)
	}

	return recommendations, nil
}

func extractResourceGroup(resourceID string) string {
	// Simple extraction from resource ID
	// Expected format: /subscriptions/{sub}/resourceGroups/{rg}/providers/...
	parts := make([]string, 0)
	current := ""
	for _, char := range resourceID {
		if char == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}

	// Find "resourceGroups" and return the next part
	for i, part := range parts {
		if part == "resourceGroups" && i+1 < len(parts) {
			return parts[i+1]
		}
	}

	return ""
}
