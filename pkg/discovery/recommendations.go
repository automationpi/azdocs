package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// FetchRecommendations fetches Azure Advisor recommendations
func (c *Client) FetchRecommendations(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
	advisorresources
	| where type == "microsoft.advisor/recommendations"
	| project
		id,
		name,
		category = tostring(properties.category),
		impact = tostring(properties.impact),
		risk = tostring(properties.risk),
		problem = tostring(properties.shortDescription.problem),
		solution = tostring(properties.shortDescription.solution),
		impactedField = tostring(properties.impactedField),
		impactedValue = tostring(properties.impactedValue),
		resourceId = tostring(properties.resourceMetadata.resourceId)
	`

	return QueryResourceGraph(ctx, c.auth, c.config.SubscriptionID, query)
}

// SaveRecommendations saves recommendations to JSON file
func SaveRecommendations(recommendations []map[string]interface{}, outputPath string) error {
	data, err := json.MarshalIndent(recommendations, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal recommendations: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write recommendations file: %w", err)
	}

	return nil
}
