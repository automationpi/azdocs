package renderer

import (
	"fmt"
	"strings"

	"github.com/automationpi/azdocs/pkg/analysis"
)

// PriorityAction represents a high-priority action item
type PriorityAction struct {
	Icon     string
	Title    string
	Impact   string
	Severity string
}

// getPriorityActions returns top priority actions across all analyses
func (r *MarkdownRenderer) getPriorityActions(security *analysis.SecurityAnalysis, cost *analysis.CostAnalysis, tagging *analysis.TaggingAnalysis, compliance *analysis.ComplianceAnalysis) []PriorityAction {
	actions := []PriorityAction{}

	// Add security critical items
	for _, finding := range security.Findings {
		if finding.Severity == "Critical" || finding.Severity == "High" {
			actions = append(actions, PriorityAction{
				Icon:     "🔴",
				Title:    finding.Issue,
				Impact:   fmt.Sprintf("Security: %s", finding.Impact),
				Severity: finding.Severity,
			})
		}
	}

	// Add top cost savings
	for _, finding := range cost.Findings {
		if finding.PotentialSavings > 20 {
			actions = append(actions, PriorityAction{
				Icon:     "💰",
				Title:    finding.Issue,
				Impact:   fmt.Sprintf("Save $%.0f/month", finding.PotentialSavings),
				Severity: "High",
			})
		}
	}

	// Add tagging issues
	for _, finding := range tagging.Findings {
		if finding.Severity == "High" && len(finding.Resources) > 5 {
			actions = append(actions, PriorityAction{
				Icon:     "🏷️",
				Title:    finding.Issue,
				Impact:   "Governance: " + finding.Impact,
				Severity: finding.Severity,
			})
		}
	}

	// Add compliance issues
	for _, finding := range compliance.Findings {
		if finding.Severity == "High" {
			actions = append(actions, PriorityAction{
				Icon:     "⚠️",
				Title:    finding.Issue,
				Impact:   finding.Impact,
				Severity: finding.Severity,
			})
		}
	}

	return actions
}

// generateSecuritySection generates the security analysis section
func (r *MarkdownRenderer) generateSecuritySection(content *strings.Builder, security *analysis.SecurityAnalysis) {
	content.WriteString(fmt.Sprintf("**Security Score:** %d/100 - %s\n\n", security.GetSecurityScore(), security.GetSecurityPosture()))

	if len(security.Findings) == 0 {
		content.WriteString("✅ No security issues detected. Excellent security posture!\n\n")
		return
	}

	// Group by category
	findingsByCategory := make(map[string][]analysis.SecurityFinding)
	for _, finding := range security.Findings {
		findingsByCategory[finding.Category] = append(findingsByCategory[finding.Category], finding)
	}

	for category, findings := range findingsByCategory {
		content.WriteString(fmt.Sprintf("### %s Issues (%d)\n\n", category, len(findings)))

		for i, finding := range findings {
			if i >= 10 {
				content.WriteString(fmt.Sprintf("*...and %d more %s issues*\n\n", len(findings)-10, category))
				break
			}

			severityIcon := getSeverityIcon(finding.Severity)
			content.WriteString(fmt.Sprintf("#### %s %s - %s\n\n", severityIcon, finding.Severity, finding.Issue))
			content.WriteString(fmt.Sprintf("**Resource:** %s\n\n", finding.Resource))
			content.WriteString(fmt.Sprintf("**Impact:** %s\n\n", finding.Impact))
			content.WriteString(fmt.Sprintf("**Remediation:** %s\n\n", finding.Remediation))
			content.WriteString("---\n\n")
		}
	}
}

// generateCostSection generates the cost optimization section
func (r *MarkdownRenderer) generateCostSection(content *strings.Builder, cost *analysis.CostAnalysis) {
	content.WriteString(fmt.Sprintf("**Cost Health:** %s (Score: %d/100)\n\n", cost.GetCostHealth(), cost.GetCostScore()))
	content.WriteString(fmt.Sprintf("**Estimated Monthly Cost:** $%.2f\n\n", cost.TotalMonthlyCost))
	content.WriteString(fmt.Sprintf("**Potential Monthly Savings:** $%.2f (%.0f%%)\n\n",
		cost.PotentialMonthlySavings,
		(cost.PotentialMonthlySavings/cost.TotalMonthlyCost)*100))

	if len(cost.Findings) == 0 {
		content.WriteString("✅ No major cost optimization opportunities detected.\n\n")
		return
	}

	// Group by category
	findingsByCategory := make(map[string][]analysis.CostFinding)
	for _, finding := range cost.Findings {
		findingsByCategory[finding.Category] = append(findingsByCategory[finding.Category], finding)
	}

	for category, findings := range findingsByCategory {
		totalSavings := 0.0
		for _, f := range findings {
			totalSavings += f.PotentialSavings
		}

		content.WriteString(fmt.Sprintf("### %s Resources (Save $%.0f/month)\n\n", category, totalSavings))
		content.WriteString("| Resource | Issue | Current Cost | Potential Savings | Remediation |\n")
		content.WriteString("|----------|-------|--------------|-------------------|-------------|\n")

		for i, finding := range findings {
			if i >= 10 {
				content.WriteString(fmt.Sprintf("| ... | *%d more %s items* | - | $%.0f | - |\n", len(findings)-10, category, totalSavings))
				break
			}

			content.WriteString(fmt.Sprintf("| %s | %s | $%.2f | $%.2f | %s |\n",
				finding.Resource,
				finding.Issue,
				finding.CurrentCost,
				finding.PotentialSavings,
				finding.Remediation))
		}
		content.WriteString("\n")
	}
}

// generateTaggingSection generates the tagging compliance section
func (r *MarkdownRenderer) generateTaggingSection(content *strings.Builder, tagging *analysis.TaggingAnalysis) {
	content.WriteString(fmt.Sprintf("**Tagging Health:** %s (Score: %d/100)\n\n", tagging.GetTaggingHealth(), tagging.GetTaggingScore()))
	content.WriteString(fmt.Sprintf("**Compliance Rate:** %.0f%% (%d/%d resources tagged)\n\n",
		tagging.ComplianceRate,
		tagging.TaggedResources,
		tagging.TotalResources))

	content.WriteString("**Required Tags:** ")
	for i, tag := range tagging.RequiredTags {
		if i > 0 {
			content.WriteString(", ")
		}
		content.WriteString(fmt.Sprintf("`%s`", tag))
	}
	content.WriteString("\n\n")

	if len(tagging.Findings) == 0 {
		content.WriteString("✅ Excellent tagging compliance!\n\n")
		return
	}

	// Group by category
	for _, finding := range tagging.Findings {
		severityIcon := getSeverityIcon(finding.Severity)
		content.WriteString(fmt.Sprintf("### %s %s: %s\n\n", severityIcon, finding.Severity, finding.Issue))
		content.WriteString(fmt.Sprintf("**Impact:** %s\n\n", finding.Impact))
		content.WriteString(fmt.Sprintf("**Remediation:** %s\n\n", finding.Remediation))

		if len(finding.Resources) > 0 && len(finding.Resources) <= 10 {
			content.WriteString("**Affected Resources:**\n")
			for _, res := range finding.Resources {
				content.WriteString(fmt.Sprintf("- %s\n", res))
			}
			content.WriteString("\n")
		} else if len(finding.Resources) > 10 {
			content.WriteString(fmt.Sprintf("**Affected Resources:** %d resources (first 10 shown)\n", len(finding.Resources)))
			for i := 0; i < 10; i++ {
				content.WriteString(fmt.Sprintf("- %s\n", finding.Resources[i]))
			}
			content.WriteString(fmt.Sprintf("- *...and %d more*\n\n", len(finding.Resources)-10))
		}

		content.WriteString("---\n\n")
	}
}

func getSeverityIcon(severity string) string {
	switch severity {
	case "Critical":
		return "🔴"
	case "High":
		return "🟠"
	case "Medium":
		return "🟡"
	case "Low":
		return "🔵"
	default:
		return "⚪"
	}
}

// generateAISecurityInsights generates AI-powered security insights section
func (r *MarkdownRenderer) generateAISecurityInsights(content *strings.Builder, insights interface{}) {
	// Type assertion to get the insights slice
	securityInsights, ok := insights.([]interface{})
	if !ok || len(securityInsights) == 0 {
		return
	}

	content.WriteString("### 🤖 AI-Powered Security Insights\n\n")
	content.WriteString("*Generated using advanced AI analysis of your security posture*\n\n")

	for i, insightData := range securityInsights {
		insight, ok := insightData.(map[string]interface{})
		if !ok {
			continue
		}

		priority, _ := insight["priority"].(float64)
		severity, _ := insight["severity"].(string)
		title, _ := insight["title"].(string)
		description, _ := insight["description"].(string)
		riskLevel, _ := insight["risk_level"].(string)
		impact, _ := insight["impact"].(string)
		category, _ := insight["category"].(string)

		severityIcon := getSeverityIcon(severity)

		content.WriteString(fmt.Sprintf("#### %s Priority %d: %s\n\n", severityIcon, int(priority), title))
		content.WriteString(fmt.Sprintf("**Category:** %s | **Severity:** %s\n\n", category, severity))
		content.WriteString(fmt.Sprintf("**Description:** %s\n\n", description))
		content.WriteString(fmt.Sprintf("**Risk Level:** %s\n\n", riskLevel))
		content.WriteString(fmt.Sprintf("**Impact:** %s\n\n", impact))

		// Add recommendations
		if recs, ok := insight["recommendations"].([]interface{}); ok && len(recs) > 0 {
			content.WriteString("**Recommendations:**\n")
			for _, rec := range recs {
				if recStr, ok := rec.(string); ok {
					content.WriteString(fmt.Sprintf("- %s\n", recStr))
				}
			}
			content.WriteString("\n")
		}

		if i < len(securityInsights)-1 {
			content.WriteString("---\n\n")
		}
	}
}

// generateAICostInsights generates AI-powered cost optimization insights section
func (r *MarkdownRenderer) generateAICostInsights(content *strings.Builder, insights interface{}) {
	// Type assertion to get the insights slice
	costInsights, ok := insights.([]interface{})
	if !ok || len(costInsights) == 0 {
		return
	}

	content.WriteString("### 🤖 AI-Powered Cost Optimization Insights\n\n")
	content.WriteString("*Generated using advanced AI analysis of your cost patterns*\n\n")

	totalEstimatedSavings := 0.0
	for _, insightData := range costInsights {
		if insight, ok := insightData.(map[string]interface{}); ok {
			if savings, ok := insight["estimated_savings"].(float64); ok {
				totalEstimatedSavings += savings
			}
		}
	}

	content.WriteString(fmt.Sprintf("**💡 Total Estimated Additional Savings:** $%.2f/month\n\n", totalEstimatedSavings))

	for i, insightData := range costInsights {
		insight, ok := insightData.(map[string]interface{})
		if !ok {
			continue
		}

		priority, _ := insight["priority"].(float64)
		title, _ := insight["title"].(string)
		description, _ := insight["description"].(string)
		savings, _ := insight["estimated_savings"].(float64)
		effort, _ := insight["effort"].(string)
		category, _ := insight["category"].(string)

		effortIcon := "🟢"
		if effort == "Medium" {
			effortIcon = "🟡"
		} else if effort == "High" {
			effortIcon = "🔴"
		}

		content.WriteString(fmt.Sprintf("#### 💰 Priority %d: %s\n\n", int(priority), title))
		content.WriteString(fmt.Sprintf("**Category:** %s | **Effort:** %s %s | **Savings:** $%.2f/month\n\n", category, effortIcon, effort, savings))
		content.WriteString(fmt.Sprintf("**Description:** %s\n\n", description))

		// Add recommendations
		if recs, ok := insight["recommendations"].([]interface{}); ok && len(recs) > 0 {
			content.WriteString("**Recommendations:**\n")
			for _, rec := range recs {
				if recStr, ok := rec.(string); ok {
					content.WriteString(fmt.Sprintf("- %s\n", recStr))
				}
			}
			content.WriteString("\n")
		}

		if i < len(costInsights)-1 {
			content.WriteString("---\n\n")
		}
	}
}
