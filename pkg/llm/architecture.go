package llm

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// ArchitectureDescription represents AI-generated description
type ArchitectureDescription struct {
	Overview              string            `json:"overview"`
	ResourceGroupInsights map[string]string `json:"resource_group_insights"`
	KeyFindings           []string          `json:"key_findings"`
}

// GenerateArchitectureDescription uses LLM to create natural language descriptions
func (c *Client) GenerateArchitectureDescription(ctx context.Context, resources []map[string]interface{}) (*ArchitectureDescription, error) {
	if !c.enabled {
		return nil, nil
	}

	// Group resources by resource group for better analysis
	rgGroups := make(map[string][]map[string]interface{})
	for _, res := range resources {
		if rg, ok := res["resourceGroup"].(string); ok {
			rgGroups[rg] = append(rgGroups[rg], res)
		}
	}

	resourceData, err := json.MarshalIndent(rgGroups, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resources: %w", err)
	}

	prompt := fmt.Sprintf(`You are an Azure architecture expert. Analyze this Azure subscription and provide insightful descriptions.

Resources grouped by Resource Group:
%s

Instructions:
1. Provide a high-level overview of the entire architecture (2-3 sentences)
2. For each resource group, provide a 1-2 sentence description of its purpose and key resources
3. Identify 3-5 key findings or notable aspects of this architecture
4. Write in professional but friendly tone
5. Focus on the "why" and "what" rather than just listing resources

Return ONLY valid JSON (no markdown):
{
  "overview": "High-level architecture summary",
  "resource_group_insights": {
    "resource-group-name": "Description of this resource group's purpose and contents"
  },
  "key_findings": [
    "Notable finding 1",
    "Notable finding 2"
  ]
}`, string(resourceData))

	resp, err := c.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an Azure architecture documentation expert. You provide JSON responses only with insightful, professional descriptions.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.4, // Slightly creative for better descriptions
		MaxTokens:   1500,
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse JSON response
	var description ArchitectureDescription
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &description); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w\nResponse: %s", err, content)
	}

	return &description, nil
}
