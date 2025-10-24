package analysis

import (
	"fmt"
	"strings"
)

// TagFinding represents a tagging issue
type TagFinding struct {
	Severity    string   // High, Medium, Low
	Category    string   // Missing, Inconsistent, Invalid
	Resources   []string // Affected resources
	Issue       string   // What's the problem
	Impact      string   // Why it matters
	Remediation string   // How to fix it
}

// TaggingAnalysis contains tagging compliance findings
type TaggingAnalysis struct {
	TotalResources      int
	TaggedResources     int
	UntaggedResources   int
	ComplianceRate      float64
	Findings            []TagFinding
	RequiredTags        []string
	TagValueInconsistencies map[string][]string // tag key -> list of different values
}

// AnalyzeTagging performs tagging compliance analysis
func AnalyzeTagging(resources []map[string]interface{}) *TaggingAnalysis {
	analysis := &TaggingAnalysis{
		TotalResources: len(resources),
		RequiredTags:   []string{"environment", "owner", "cost-center", "application"},
		TagValueInconsistencies: make(map[string][]string),
		Findings:       []TagFinding{},
	}

	// Analyze missing tags
	analysis.analyzeMissingTags(resources)

	// Analyze tag consistency
	analysis.analyzeTagConsistency(resources)

	// Analyze tag value patterns
	analysis.analyzeTagValuePatterns(resources)

	// Calculate compliance rate
	if analysis.TotalResources > 0 {
		analysis.ComplianceRate = (float64(analysis.TaggedResources) / float64(analysis.TotalResources)) * 100
	}

	return analysis
}

func (a *TaggingAnalysis) analyzeMissingTags(resources []map[string]interface{}) {
	missingByTag := make(map[string][]string) // tag name -> list of resources missing it

	for _, res := range resources {
		name, _ := res["name"].(string)
		resType, _ := res["type"].(string)

		// Skip certain system resource types
		if isSystemResource(resType) {
			continue
		}

		tags, ok := res["tags"].(map[string]interface{})
		if !ok || len(tags) == 0 {
			a.UntaggedResources++
			for _, requiredTag := range a.RequiredTags {
				missingByTag[requiredTag] = append(missingByTag[requiredTag], name)
			}
			continue
		}

		a.TaggedResources++

		// Check for missing required tags
		for _, requiredTag := range a.RequiredTags {
			found := false
			for tagKey := range tags {
				if strings.EqualFold(tagKey, requiredTag) {
					found = true
					break
				}
			}
			if !found {
				missingByTag[requiredTag] = append(missingByTag[requiredTag], name)
			}
		}
	}

	// Create findings for missing tags
	for tagName, resourceList := range missingByTag {
		if len(resourceList) > 0 {
			severity := "Medium"
			if tagName == "owner" || tagName == "cost-center" {
				severity = "High"
			}

			a.Findings = append(a.Findings, TagFinding{
				Severity:    severity,
				Category:    "Missing",
				Resources:   resourceList,
				Issue:       fmt.Sprintf("%d resources missing '%s' tag", len(resourceList), tagName),
				Impact:      "Cannot track ownership, cost allocation, or compliance",
				Remediation: fmt.Sprintf("Add '%s' tag to all resources according to tagging policy", tagName),
			})
		}
	}
}

func (a *TaggingAnalysis) analyzeTagConsistency(resources []map[string]interface{}) {
	// Collect all variations of tag keys
	tagKeyVariations := make(map[string]map[string]int) // normalized key -> {actual key -> count}

	for _, res := range resources {
		tags, ok := res["tags"].(map[string]interface{})
		if !ok {
			continue
		}

		for tagKey := range tags {
			normalized := strings.ToLower(tagKey)
			if tagKeyVariations[normalized] == nil {
				tagKeyVariations[normalized] = make(map[string]int)
			}
			tagKeyVariations[normalized][tagKey]++
		}
	}

	// Find inconsistencies
	for normalized, variations := range tagKeyVariations {
		if len(variations) > 1 {
			varList := []string{}
			for variant := range variations {
				varList = append(varList, variant)
			}

			a.Findings = append(a.Findings, TagFinding{
				Severity:    "Low",
				Category:    "Inconsistent",
				Resources:   []string{}, // All resources with these variations
				Issue:       fmt.Sprintf("Tag key '%s' has %d variations: %v", normalized, len(variations), varList),
				Impact:      "Makes filtering and cost reporting difficult",
				Remediation: fmt.Sprintf("Standardize to single format (recommended: '%s')", normalized),
			})
		}
	}
}

func (a *TaggingAnalysis) analyzeTagValuePatterns(resources []map[string]interface{}) {
	// Collect common tag values
	envValues := make(map[string]int)

	for _, res := range resources {
		tags, ok := res["tags"].(map[string]interface{})
		if !ok {
			continue
		}

		// Check environment tag variations
		for tagKey, tagValue := range tags {
			if strings.EqualFold(tagKey, "environment") || strings.EqualFold(tagKey, "env") {
				if val, ok := tagValue.(string); ok {
					normalized := strings.ToLower(val)
					envValues[normalized]++
				}
			}
		}
	}

	// Check for inconsistent environment values
	if len(envValues) > 5 {
		values := []string{}
		for val := range envValues {
			values = append(values, val)
		}

		a.Findings = append(a.Findings, TagFinding{
			Severity:    "Medium",
			Category:    "Inconsistent",
			Resources:   []string{},
			Issue:       fmt.Sprintf("Environment tag has %d different values: %v", len(envValues), values),
			Impact:      "Difficult to filter resources by environment",
			Remediation: "Standardize environment values to: production, staging, development, test",
		})
	}
}

func isSystemResource(resType string) bool {
	systemTypes := []string{
		"microsoft.network/networkwatchers",
		"microsoft.insights/actiongroups",
		"microsoft.operationalinsights/workspaces",
	}

	resTypeLower := strings.ToLower(resType)
	for _, sysType := range systemTypes {
		if resTypeLower == sysType {
			return true
		}
	}

	return false
}

// GetTaggingScore calculates tagging compliance score (0-100)
func (a *TaggingAnalysis) GetTaggingScore() int {
	score := int(a.ComplianceRate)

	// Penalize for findings
	for _, finding := range a.Findings {
		switch finding.Severity {
		case "High":
			score -= 10
		case "Medium":
			score -= 5
		case "Low":
			score -= 2
		}
	}

	if score < 0 {
		score = 0
	}

	return score
}

// GetTaggingHealth returns a human-readable tagging health status
func (a *TaggingAnalysis) GetTaggingHealth() string {
	score := a.GetTaggingScore()

	if score >= 90 {
		return "✅ EXCELLENT"
	} else if score >= 75 {
		return "✅ GOOD"
	} else if score >= 50 {
		return "⚠️ NEEDS ATTENTION"
	} else {
		return "🔴 POOR"
	}
}
