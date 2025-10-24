package analysis

import (
	"fmt"
	"strings"
)

// CostFinding represents a cost optimization opportunity
type CostFinding struct {
	Severity         string  // High, Medium, Low
	Category         string  // Idle, Oversized, Orphaned, etc.
	Resource         string  // Resource name
	Issue            string  // What's the problem
	CurrentCost      float64 // Estimated current monthly cost
	PotentialSavings float64 // Estimated monthly savings
	Remediation      string  // How to fix it
}

// CostAnalysis contains cost optimization findings
type CostAnalysis struct {
	TotalMonthlyCost      float64
	PotentialMonthlySavings float64
	Findings              []CostFinding
}

// AnalyzeCost performs cost optimization analysis
func AnalyzeCost(resources []map[string]interface{}) *CostAnalysis {
	analysis := &CostAnalysis{
		Findings: []CostFinding{},
	}

	// Analyze orphaned resources
	analysis.analyzeOrphanedResources(resources)

	// Analyze idle resources
	analysis.analyzeIdleResources(resources)

	// Analyze oversized resources (based on SKU)
	analysis.analyzeOversizedResources(resources)

	// Analyze storage tier optimization
	analysis.analyzeStorageTiers(resources)

	// Calculate totals
	for _, finding := range analysis.Findings {
		analysis.PotentialMonthlySavings += finding.PotentialSavings
	}

	// Estimate total cost (rough estimates)
	analysis.estimateTotalCost(resources)

	return analysis
}

func (a *CostAnalysis) analyzeOrphanedResources(resources []map[string]interface{}) {
	// Build map of attached disks
	attachedDisks := make(map[string]bool)
	for _, res := range resources {
		resType, _ := res["type"].(string)
		if strings.ToLower(resType) != "microsoft.compute/virtualmachines" {
			continue
		}

		props, ok := res["properties"].(map[string]interface{})
		if !ok {
			continue
		}

		// Check OS disk
		if osDisk, ok := props["storageProfile"].(map[string]interface{}); ok {
			if managedDisk, ok := osDisk["osDisk"].(map[string]interface{}); ok {
				if diskRef, ok := managedDisk["managedDisk"].(map[string]interface{}); ok {
					if id, ok := diskRef["id"].(string); ok {
						attachedDisks[id] = true
					}
				}
			}
		}
	}

	// Check for orphaned disks
	for _, res := range resources {
		resType, _ := res["type"].(string)
		if strings.ToLower(resType) != "microsoft.compute/disks" {
			continue
		}

		id, _ := res["id"].(string)
		name, _ := res["name"].(string)

		if !attachedDisks[id] {
			props, ok := res["properties"].(map[string]interface{})
			diskSize := 128 // default
			if ok {
				if size, ok := props["diskSizeGB"].(float64); ok {
					diskSize = int(size)
				}
			}

			// Estimate cost: ~$0.12/GB/month for Standard HDD
			estimatedCost := float64(diskSize) * 0.12

			a.Findings = append(a.Findings, CostFinding{
				Severity:         "Medium",
				Category:         "Orphaned",
				Resource:         name,
				Issue:            "Disk is not attached to any VM",
				CurrentCost:      estimatedCost,
				PotentialSavings: estimatedCost,
				Remediation:      "Delete orphaned disk or attach to a VM if needed",
			})
		}
	}

	// Check for orphaned public IPs
	usedPublicIPs := make(map[string]bool)
	for _, res := range resources {
		resType, _ := res["type"].(string)
		if strings.ToLower(resType) != "microsoft.network/networkinterfaces" {
			continue
		}

		props, ok := res["properties"].(map[string]interface{})
		if !ok {
			continue
		}

		if ipConfigs, ok := props["ipConfigurations"].([]interface{}); ok {
			for _, ipConfigIface := range ipConfigs {
				if ipConfig, ok := ipConfigIface.(map[string]interface{}); ok {
					if ipProps, ok := ipConfig["properties"].(map[string]interface{}); ok {
						if pubIP, ok := ipProps["publicIPAddress"].(map[string]interface{}); ok {
							if id, ok := pubIP["id"].(string); ok {
								usedPublicIPs[id] = true
							}
						}
					}
				}
			}
		}
	}

	for _, res := range resources {
		resType, _ := res["type"].(string)
		if strings.ToLower(resType) != "microsoft.network/publicipaddresses" {
			continue
		}

		id, _ := res["id"].(string)
		name, _ := res["name"].(string)

		if !usedPublicIPs[id] {
			// Static IPs cost ~$3.65/month
			a.Findings = append(a.Findings, CostFinding{
				Severity:         "Low",
				Category:         "Orphaned",
				Resource:         name,
				Issue:            "Public IP is not associated with any resource",
				CurrentCost:      3.65,
				PotentialSavings: 3.65,
				Remediation:      "Delete unused public IP or associate with a resource",
			})
		}
	}
}

func (a *CostAnalysis) analyzeIdleResources(resources []map[string]interface{}) {
	for _, res := range resources {
		resType, _ := res["type"].(string)
		name, _ := res["name"].(string)

		if strings.ToLower(resType) == "microsoft.compute/virtualmachines" {
			props, ok := res["properties"].(map[string]interface{})
			if !ok {
				continue
			}

			// Check power state (if available in properties)
			// Note: This requires extended properties from Azure Resource Graph
			// For now, we'll add a finding suggesting to check for idle VMs

			vmSize := "Standard_B2s" // default
			if hardwareProfile, ok := props["hardwareProfile"].(map[string]interface{}); ok {
				if size, ok := hardwareProfile["vmSize"].(string); ok {
					vmSize = size
				}
			}

			estimatedCost := estimateVMCost(vmSize)

			a.Findings = append(a.Findings, CostFinding{
				Severity:         "Low",
				Category:         "Idle",
				Resource:         name,
				Issue:            "VM may be idle or underutilized (requires metrics analysis)",
				CurrentCost:      estimatedCost,
				PotentialSavings: estimatedCost * 0.7, // Could save 70% by deallocating
				Remediation:      "Review VM metrics. If CPU <5% avg, consider deallocating when not in use or downsizing",
			})
		}
	}
}

func (a *CostAnalysis) analyzeOversizedResources(resources []map[string]interface{}) {
	for _, res := range resources {
		resType, _ := res["type"].(string)
		name, _ := res["name"].(string)

		switch strings.ToLower(resType) {
		case "microsoft.compute/virtualmachines":
			props, ok := res["properties"].(map[string]interface{})
			if !ok {
				continue
			}

			vmSize := ""
			if hardwareProfile, ok := props["hardwareProfile"].(map[string]interface{}); ok {
				if size, ok := hardwareProfile["vmSize"].(string); ok {
					vmSize = size
				}
			}

			// Check if using expensive SKUs that could be downgraded
			if strings.Contains(vmSize, "Standard_D") || strings.Contains(vmSize, "Standard_E") {
				currentCost := estimateVMCost(vmSize)
				potentialSavings := currentCost * 0.4 // 40% savings with B-series

				a.Findings = append(a.Findings, CostFinding{
					Severity:         "Medium",
					Category:         "Oversized",
					Resource:         name,
					Issue:            fmt.Sprintf("VM using %s - may be oversized for workload", vmSize),
					CurrentCost:      currentCost,
					PotentialSavings: potentialSavings,
					Remediation:      "Analyze CPU/Memory metrics. Consider B-series or smaller size if utilization <30%",
				})
			}
		}
	}
}

func (a *CostAnalysis) analyzeStorageTiers(resources []map[string]interface{}) {
	for _, res := range resources {
		resType, _ := res["type"].(string)
		name, _ := res["name"].(string)

		if strings.ToLower(resType) == "microsoft.storage/storageaccounts" {
			props, ok := res["properties"].(map[string]interface{})
			if !ok {
				continue
			}

			accessTier, _ := props["accessTier"].(string)
			if accessTier == "Hot" {
				// Suggest reviewing if data is accessed infrequently
				a.Findings = append(a.Findings, CostFinding{
					Severity:         "Low",
					Category:         "StorageTier",
					Resource:         name,
					Issue:            "Storage account using Hot tier - review access patterns",
					CurrentCost:      50, // Estimated
					PotentialSavings: 25, // 50% savings with Cool tier
					Remediation:      "If data is accessed <1x/month, move to Cool tier. For archival, use Archive tier",
				})
			}
		}
	}
}

func (a *CostAnalysis) estimateTotalCost(resources []map[string]interface{}) {
	totalCost := 0.0

	for _, res := range resources {
		resType, _ := res["type"].(string)

		switch strings.ToLower(resType) {
		case "microsoft.compute/virtualmachines":
			props, _ := res["properties"].(map[string]interface{})
			vmSize := "Standard_B2s"
			if props != nil {
				if hw, ok := props["hardwareProfile"].(map[string]interface{}); ok {
					if size, ok := hw["vmSize"].(string); ok {
						vmSize = size
					}
				}
			}
			totalCost += estimateVMCost(vmSize)

		case "microsoft.storage/storageaccounts":
			totalCost += 50 // Rough estimate

		case "microsoft.network/natgateways":
			totalCost += 33 // ~$33/month

		case "microsoft.network/publicipaddresses":
			totalCost += 3.65

		case "microsoft.compute/disks":
			props, _ := res["properties"].(map[string]interface{})
			diskSize := 128.0
			if props != nil {
				if size, ok := props["diskSizeGB"].(float64); ok {
					diskSize = size
				}
			}
			totalCost += diskSize * 0.12
		}
	}

	a.TotalMonthlyCost = totalCost
}

func estimateVMCost(vmSize string) float64 {
	// Rough monthly cost estimates for common VM sizes
	costMap := map[string]float64{
		"Standard_B1s":  7.5,
		"Standard_B2s":  30,
		"Standard_B2ms": 60,
		"Standard_D2s_v3": 70,
		"Standard_D4s_v3": 140,
		"Standard_E2s_v3": 87,
		"Standard_E4s_v3": 175,
	}

	if cost, ok := costMap[vmSize]; ok {
		return cost
	}

	// Default estimate if size not in map
	return 50
}

// GetCostScore calculates cost optimization score (0-100)
func (a *CostAnalysis) GetCostScore() int {
	if a.TotalMonthlyCost == 0 {
		return 100
	}

	optimizationPct := (a.PotentialMonthlySavings / a.TotalMonthlyCost) * 100

	if optimizationPct < 10 {
		return 100 // Excellent, <10% potential savings
	} else if optimizationPct < 25 {
		return 80 // Good, 10-25% potential savings
	} else if optimizationPct < 40 {
		return 60 // Moderate, 25-40% potential savings
	} else {
		return 40 // Poor, >40% potential savings
	}
}

// GetCostHealth returns a human-readable cost health status
func (a *CostAnalysis) GetCostHealth() string {
	score := a.GetCostScore()

	if score >= 90 {
		return "✅ EXCELLENT"
	} else if score >= 75 {
		return "✅ GOOD"
	} else if score >= 50 {
		return "⚠️ NEEDS OPTIMIZATION"
	} else {
		return "🔴 HIGH WASTE"
	}
}
