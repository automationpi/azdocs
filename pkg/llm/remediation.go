package llm

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// RemediationStep represents a detailed remediation step
type RemediationStep struct {
	Step        int    `json:"step"`
	Action      string `json:"action"`
	Command     string `json:"command,omitempty"`
	Description string `json:"description"`
}

// RecommendationWithRemediation combines Azure recommendation with AI-generated remediation
type RecommendationWithRemediation struct {
	Category          string            `json:"category"`
	Impact            string            `json:"impact"`
	Risk              string            `json:"risk"`
	Problem           string            `json:"problem"`
	AzureSolution     string            `json:"azure_solution"`
	RemediationSteps  []RemediationStep `json:"remediation_steps"`
	EstimatedTime     string            `json:"estimated_time"`
	PreRequisites     []string          `json:"prerequisites"`
	ValidationCommand string            `json:"validation_command,omitempty"`
}

// GenerateRemediationSteps uses LLM to create detailed remediation steps
func (c *Client) GenerateRemediationSteps(ctx context.Context, recommendations []interface{}) ([]RecommendationWithRemediation, error) {
	if !c.enabled {
		return nil, nil
	}

	if len(recommendations) == 0 {
		return nil, nil
	}

	recData, err := json.MarshalIndent(recommendations, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal recommendations: %w", err)
	}

	prompt := fmt.Sprintf(`You are an Azure cloud operations expert. Analyze these Azure Advisor recommendations and provide detailed, actionable remediation steps.

Azure Advisor Recommendations:
%s

Instructions:
1. For each recommendation, provide step-by-step remediation instructions
2. Include specific Azure CLI, PowerShell, or Portal commands where applicable
3. Estimate time required for remediation
4. List any prerequisites needed before starting
5. Provide validation commands to confirm the fix
6. Be specific and actionable - assume the user is a DevOps engineer
7. Focus on automation and infrastructure-as-code approaches where possible

Return ONLY valid JSON array (no markdown):
[
  {
    "category": "Security/Cost/Performance/Reliability/OperationalExcellence",
    "impact": "High/Medium/Low",
    "risk": "Warning/Error/None",
    "problem": "Clear description of the issue",
    "azure_solution": "Azure's recommended solution",
    "remediation_steps": [
      {
        "step": 1,
        "action": "What to do",
        "command": "az cli command or portal steps",
        "description": "Why this step is needed"
      }
    ],
    "estimated_time": "5 minutes / 30 minutes / 2 hours",
    "prerequisites": ["List", "of", "requirements"],
    "validation_command": "Command to verify the fix worked"
  }
]`, string(recData))

	resp, err := c.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert Azure DevOps engineer. You provide JSON responses only with detailed, practical remediation steps for Azure issues.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3, // Lower temperature for more precise technical instructions
		MaxTokens:   2500,
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse JSON response
	var remediations []RecommendationWithRemediation
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &remediations); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w\nResponse: %s", err, content)
	}

	return remediations, nil
}
