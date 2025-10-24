package analysis

import (
	"fmt"
	"strings"
)

// ComplianceFinding represents a compliance issue
type ComplianceFinding struct {
	Category    string   // DR, Monitoring, Backup
	Severity    string   // Critical, High, Medium, Low
	Resources   []string // Affected resources
	Issue       string
	Impact      string
	Remediation string
}

// ComplianceAnalysis contains DR and monitoring findings
type ComplianceAnalysis struct {
	BackupCoverage      float64 // Percentage of VMs with backup
	MonitoringCoverage  float64 // Percentage of resources with diagnostics
	Findings            []ComplianceFinding
}

// AnalyzeCompliance performs DR and monitoring analysis
func AnalyzeCompliance(resources []map[string]interface{}) *ComplianceAnalysis {
	analysis := &ComplianceAnalysis{
		Findings: []ComplianceFinding{},
	}

	analysis.analyzeBackup(resources)
	analysis.analyzeMonitoring(resources)
	analysis.analyzeGeoRedundancy(resources)

	return analysis
}

func (a *ComplianceAnalysis) analyzeBackup(resources []map[string]interface{}) {
	vmsWithoutBackup := []string{}
	totalVMs := 0

	for _, res := range resources {
		resType, _ := res["type"].(string)
		if strings.ToLower(resType) != "microsoft.compute/virtualmachines" {
			continue
		}

		totalVMs++
		name, _ := res["name"].(string)

		// In real implementation, check Azure Backup vault associations
		// For now, assume VMs without backup
		vmsWithoutBackup = append(vmsWithoutBackup, name)
	}

	if totalVMs > 0 {
		a.BackupCoverage = float64(totalVMs-len(vmsWithoutBackup)) / float64(totalVMs) * 100

		if len(vmsWithoutBackup) > 0 {
			a.Findings = append(a.Findings, ComplianceFinding{
				Category:    "Backup",
				Severity:    "High",
				Resources:   vmsWithoutBackup,
				Issue:       fmt.Sprintf("%d VMs without Azure Backup configured", len(vmsWithoutBackup)),
				Impact:      "Risk of data loss if VM fails or is corrupted",
				Remediation: "Configure Azure Backup with appropriate retention policy (7-30 days recommended)",
			})
		}
	}
}

func (a *ComplianceAnalysis) analyzeMonitoring(resources []map[string]interface{}) {
	resourcesWithoutDiagnostics := []string{}
	totalMonitorable := 0

	for _, res := range resources {
		resType, _ := res["type"].(string)
		name, _ := res["name"].(string)

		// Only check monitorable resources
		if !isMonitorable(resType) {
			continue
		}

		totalMonitorable++

		// In real implementation, check diagnostic settings
		// For now, assume resources need diagnostics
		resourcesWithoutDiagnostics = append(resourcesWithoutDiagnostics, fmt.Sprintf("%s (%s)", name, resType))
	}

	if totalMonitorable > 0 {
		a.MonitoringCoverage = float64(totalMonitorable-len(resourcesWithoutDiagnostics)) / float64(totalMonitorable) * 100

		if len(resourcesWithoutDiagnostics) > 5 {
			a.Findings = append(a.Findings, ComplianceFinding{
				Category:    "Monitoring",
				Severity:    "Medium",
				Resources:   resourcesWithoutDiagnostics[:5], // Show first 5
				Issue:       fmt.Sprintf("%d resources without diagnostic settings (showing first 5)", len(resourcesWithoutDiagnostics)),
				Impact:      "Limited visibility into resource health and performance",
				Remediation: "Enable diagnostic settings to send logs to Log Analytics workspace",
			})
		} else if len(resourcesWithoutDiagnostics) > 0 {
			a.Findings = append(a.Findings, ComplianceFinding{
				Category:    "Monitoring",
				Severity:    "Medium",
				Resources:   resourcesWithoutDiagnostics,
				Issue:       fmt.Sprintf("%d resources without diagnostic settings", len(resourcesWithoutDiagnostics)),
				Impact:      "Limited visibility into resource health and performance",
				Remediation: "Enable diagnostic settings to send logs to Log Analytics workspace",
			})
		}
	}
}

func (a *ComplianceAnalysis) analyzeGeoRedundancy(resources []map[string]interface{}) {
	nonRedundantStorage := []string{}

	for _, res := range resources {
		resType, _ := res["type"].(string)
		if strings.ToLower(resType) != "microsoft.storage/storageaccounts" {
			continue
		}

		name, _ := res["name"].(string)
		sku, ok := res["sku"].(map[string]interface{})
		if !ok {
			continue
		}

		skuName, _ := sku["name"].(string)

		// Check if using LRS (not geo-redundant)
		if strings.Contains(skuName, "LRS") {
			nonRedundantStorage = append(nonRedundantStorage, name)
		}
	}

	if len(nonRedundantStorage) > 0 {
		a.Findings = append(a.Findings, ComplianceFinding{
			Category:    "DR",
			Severity:    "Medium",
			Resources:   nonRedundantStorage,
			Issue:       fmt.Sprintf("%d storage accounts using LRS (locally redundant)", len(nonRedundantStorage)),
			Impact:      "Data not protected against regional outages",
			Remediation: "Consider GRS (Geo-Redundant Storage) or GZRS for critical data",
		})
	}
}

func isMonitorable(resType string) bool {
	monitorableTypes := []string{
		"microsoft.compute/virtualmachines",
		"microsoft.storage/storageaccounts",
		"microsoft.network/applicationgateways",
		"microsoft.network/loadbalancers",
		"microsoft.web/sites",
		"microsoft.sql/servers",
	}

	resTypeLower := strings.ToLower(resType)
	for _, monType := range monitorableTypes {
		if resTypeLower == monType {
			return true
		}
	}

	return false
}

// GetComplianceScore calculates overall compliance score (0-100)
func (a *ComplianceAnalysis) GetComplianceScore() int {
	score := 100

	// Weight backup and monitoring coverage
	score -= int((100 - a.BackupCoverage) * 0.3)
	score -= int((100 - a.MonitoringCoverage) * 0.2)

	// Penalize for findings
	for _, finding := range a.Findings {
		switch finding.Severity {
		case "Critical":
			score -= 15
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

// GetComplianceHealth returns a human-readable compliance status
func (a *ComplianceAnalysis) GetComplianceHealth() string {
	score := a.GetComplianceScore()

	if score >= 90 {
		return "✅ EXCELLENT"
	} else if score >= 75 {
		return "✅ GOOD"
	} else if score >= 50 {
		return "⚠️ NEEDS ATTENTION"
	} else {
		return "🔴 CRITICAL"
	}
}
