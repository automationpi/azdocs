package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// Client handles LLM API interactions
type Client struct {
	openai  *openai.Client
	enabled bool
	model   string
}

// Config holds LLM configuration
type Config struct {
	Enabled bool
	APIKey  string
	Model   string
}

// NewClient creates a new LLM client
func NewClient(config Config) *Client {
	if !config.Enabled {
		return &Client{enabled: false}
	}

	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	if apiKey == "" {
		fmt.Println("Warning: OpenAI API key not found. AI features disabled.")
		return &Client{enabled: false}
	}

	model := config.Model
	if model == "" {
		model = openai.GPT4TurboPreview // Use GPT-4 Turbo for better structured output
	}

	return &Client{
		openai:  openai.NewClient(apiKey),
		enabled: true,
		model:   model,
	}
}

// IsEnabled returns whether AI features are enabled
func (c *Client) IsEnabled() bool {
	return c.enabled
}

// ConnectionSuggestion represents a suggested connection between resources
type ConnectionSuggestion struct {
	SourceResource string `json:"source_resource"`
	TargetResource string `json:"target_resource"`
	ConnectionType string `json:"connection_type"` // "association", "peering", "routing", "natgw"
	Label          string `json:"label"`
	Confidence     string `json:"confidence"` // "high", "medium", "low"
	Reason         string `json:"reason"`
}

// LayoutSuggestion represents positioning suggestions for resources
type LayoutSuggestion struct {
	ResourceName string  `json:"resource_name"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	SubnetName   string  `json:"subnet_name,omitempty"`
	Grouping     string  `json:"grouping"` // "subnet", "tier", "function"
	Reason       string  `json:"reason"`
}

// DiscoverConnections uses LLM to analyze Azure resources and suggest connections
func (c *Client) DiscoverConnections(ctx context.Context, resources []map[string]interface{}) ([]ConnectionSuggestion, error) {
	if !c.enabled {
		return nil, nil
	}

	// Prepare resource data for LLM
	resourceData, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resources: %w", err)
	}

	prompt := fmt.Sprintf(`You are an Azure architecture expert. Analyze these Azure resources and identify all logical connections and relationships between them.

Resources:
%s

Instructions:
1. Identify relationships based on Azure properties (e.g., NICs attached to VMs, Function Apps using Storage Accounts, NSG associations, VNet peerings)
2. Look for subnet associations, VNet integrations, public IP attachments
3. Consider naming conventions (e.g., "vm-hub-jumpbox-nic" likely connects to "vm-hub-jumpbox")
4. Return ONLY valid JSON array of connections

Return format (valid JSON only, no markdown):
[
  {
    "source_resource": "resource-name-1",
    "target_resource": "resource-name-2",
    "connection_type": "association|peering|routing|natgw",
    "label": "short description",
    "confidence": "high|medium|low",
    "reason": "why this connection exists"
  }
]`, string(resourceData))

	resp, err := c.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an Azure infrastructure analysis expert. You provide JSON responses only.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3, // Lower temperature for more deterministic output
		MaxTokens:   2000,
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse JSON response
	var suggestions []ConnectionSuggestion
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &suggestions); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w\nResponse: %s", err, content)
	}

	return suggestions, nil
}

// OptimizeLayout uses LLM to suggest optimal positioning for resources
func (c *Client) OptimizeLayout(ctx context.Context, resources []map[string]interface{}, subnets []map[string]string) ([]LayoutSuggestion, error) {
	if !c.enabled {
		return nil, nil
	}

	// Prepare data
	resourceData, _ := json.MarshalIndent(resources, "", "  ")
	subnetData, _ := json.MarshalIndent(subnets, "", "  ")

	prompt := fmt.Sprintf(`You are an Azure network diagram expert. Analyze these resources and suggest optimal positioning for a clear, professional diagram.

Resources:
%s

Subnets:
%s

Diagram constraints:
- VNet container: 1100x600 pixels, starting at (50, 50)
- Subnet containers: 350x120 pixels each, stacked vertically
- Resources should be inside subnets when applicable
- Leave space between resources (minimum 80px)

Instructions:
1. Place compute resources (VMs, Function Apps) INSIDE appropriate subnets
2. Place network resources (NSGs, Route Tables, LBs, NAT GWs) outside subnets but near them
3. Group related resources together
4. Avoid overlapping
5. Create visual hierarchy (most important resources prominently placed)

Return ONLY valid JSON array (no markdown):
[
  {
    "resource_name": "exact-resource-name",
    "x": 100.0,
    "y": 150.0,
    "subnet_name": "subnet-name-if-inside",
    "grouping": "subnet|tier|function",
    "reason": "why positioned here"
  }
]`, string(resourceData), string(subnetData))

	resp, err := c.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a diagram layout optimization expert. You provide JSON responses only.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.2,
		MaxTokens:   3000,
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse JSON response
	var suggestions []LayoutSuggestion
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &suggestions); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w\nResponse: %s", err, content)
	}

	return suggestions, nil
}

// SecurityInsight represents AI-generated security analysis
type SecurityInsight struct {
	Category        string   `json:"category"`
	Severity        string   `json:"severity"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	RiskLevel       string   `json:"risk_level"`
	Impact          string   `json:"impact"`
	Recommendations []string `json:"recommendations"`
	Priority        int      `json:"priority"`
}

// GenerateSecurityInsights analyzes security findings and generates AI insights
func (c *Client) GenerateSecurityInsights(ctx context.Context, securityFindings interface{}) ([]SecurityInsight, error) {
	if !c.enabled {
		return nil, nil
	}

	// Marshal findings to JSON
	findingsData, err := json.MarshalIndent(securityFindings, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal security findings: %w", err)
	}

	prompt := fmt.Sprintf(`You are a cybersecurity expert specializing in Azure security. Analyze these security findings and provide actionable insights.

Security Findings:
%s

Instructions:
1. Identify the most critical security risks based on the findings
2. Group related findings into broader security themes
3. Provide specific, actionable recommendations
4. Prioritize based on risk level and ease of remediation
5. Consider Azure best practices and Well-Architected Framework

Return ONLY valid JSON array (no markdown):
[
  {
    "category": "Network Security|Identity|Encryption|Compliance",
    "severity": "Critical|High|Medium|Low",
    "title": "concise insight title",
    "description": "detailed explanation of the security concern",
    "risk_level": "description of potential risk",
    "impact": "what could happen if not addressed",
    "recommendations": ["specific action 1", "specific action 2"],
    "priority": 1-10
  }
]

Focus on the top 5-8 most important insights only.`, string(findingsData))

	resp, err := c.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an Azure security expert. You provide JSON responses only with actionable security insights.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   2500,
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse JSON response
	var insights []SecurityInsight
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &insights); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w\nResponse: %s", err, content)
	}

	return insights, nil
}

// CostOptimizationInsight represents AI-generated cost optimization analysis
type CostOptimizationInsight struct {
	Category        string   `json:"category"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	EstimatedSavings float64 `json:"estimated_savings"`
	Effort          string   `json:"effort"`
	Recommendations []string `json:"recommendations"`
	Priority        int      `json:"priority"`
}

// GenerateCostOptimizationInsights analyzes cost findings and generates AI insights
func (c *Client) GenerateCostOptimizationInsights(ctx context.Context, costFindings interface{}) ([]CostOptimizationInsight, error) {
	if !c.enabled {
		return nil, nil
	}

	// Marshal findings to JSON
	findingsData, err := json.MarshalIndent(costFindings, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cost findings: %w", err)
	}

	prompt := fmt.Sprintf(`You are a cloud cost optimization expert specializing in Azure. Analyze these cost findings and provide strategic insights.

Cost Findings:
%s

Instructions:
1. Identify the highest-impact cost optimization opportunities
2. Group related findings into optimization themes
3. Provide specific, actionable recommendations
4. Estimate effort level (Low/Medium/High) for each recommendation
5. Consider Azure pricing models, reserved instances, and architectural improvements

Return ONLY valid JSON array (no markdown):
[
  {
    "category": "Compute|Storage|Network|Database|Licensing",
    "title": "concise optimization opportunity",
    "description": "detailed explanation of the cost issue and opportunity",
    "estimated_savings": 123.45,
    "effort": "Low|Medium|High",
    "recommendations": ["specific action 1", "specific action 2"],
    "priority": 1-10
  }
]

Focus on the top 5-8 most impactful insights only.`, string(findingsData))

	resp, err := c.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an Azure cost optimization expert. You provide JSON responses only with actionable cost-saving insights.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   2500,
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse JSON response
	var insights []CostOptimizationInsight
	content := resp.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &insights); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w\nResponse: %s", err, content)
	}

	return insights, nil
}
