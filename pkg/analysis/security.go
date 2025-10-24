package analysis

import (
	"fmt"
	"strings"
)

// SecurityFinding represents a security issue
type SecurityFinding struct {
	Severity    string   // Critical, High, Medium, Low
	Category    string   // NSG, PublicExposure, Encryption, etc.
	Resource    string   // Resource name
	Issue       string   // What's the problem
	Impact      string   // Why it matters
	Remediation string   // How to fix it
	References  []string // Related resources
}

// SecurityAnalysis contains all security findings
type SecurityAnalysis struct {
	CriticalCount int
	HighCount     int
	MediumCount   int
	LowCount      int
	Findings      []SecurityFinding
}

// AnalyzeSecurity performs comprehensive security analysis
func AnalyzeSecurity(resources []map[string]interface{}) *SecurityAnalysis {
	analysis := &SecurityAnalysis{
		Findings: []SecurityFinding{},
	}

	// Analyze NSG rules
	analysis.analyzeNSGRules(resources)

	// Analyze public exposure
	analysis.analyzePublicExposure(resources)

	// Analyze encryption
	analysis.analyzeEncryption(resources)

	// Analyze network isolation
	analysis.analyzeNetworkIsolation(resources)

	// Count by severity
	for _, finding := range analysis.Findings {
		switch finding.Severity {
		case "Critical":
			analysis.CriticalCount++
		case "High":
			analysis.HighCount++
		case "Medium":
			analysis.MediumCount++
		case "Low":
			analysis.LowCount++
		}
	}

	return analysis
}

func (a *SecurityAnalysis) analyzeNSGRules(resources []map[string]interface{}) {
	for _, res := range resources {
		resType, _ := res["type"].(string)
		if strings.ToLower(resType) != "microsoft.network/networksecuritygroups" {
			continue
		}

		nsgName, _ := res["name"].(string)
		props, ok := res["properties"].(map[string]interface{})
		if !ok {
			continue
		}

		securityRules, ok := props["securityRules"].([]interface{})
		if !ok {
			continue
		}

		for _, ruleIface := range securityRules {
			rule, ok := ruleIface.(map[string]interface{})
			if !ok {
				continue
			}

			ruleProps, ok := rule["properties"].(map[string]interface{})
			if !ok {
				continue
			}

			// Check for overly permissive rules
			direction, _ := ruleProps["direction"].(string)
			access, _ := ruleProps["access"].(string)
			sourcePrefix, _ := ruleProps["sourceAddressPrefix"].(string)
			destPort, _ := ruleProps["destinationPortRange"].(string)
			ruleName, _ := rule["name"].(string)

			if direction == "Inbound" && access == "Allow" {
				// Check for 0.0.0.0/0 or * or Internet
				if sourcePrefix == "0.0.0.0/0" || sourcePrefix == "*" || sourcePrefix == "Internet" {
					severity := "Medium"
					issue := fmt.Sprintf("NSG rule '%s' allows inbound traffic from Internet (%s)", ruleName, sourcePrefix)

					// Critical if dangerous ports are open
					if isDangerousPort(destPort) {
						severity = "Critical"
						issue = fmt.Sprintf("NSG rule '%s' allows %s from Internet (port %s)", ruleName, destPort, destPort)
					}

					a.Findings = append(a.Findings, SecurityFinding{
						Severity:    severity,
						Category:    "NSG",
						Resource:    nsgName,
						Issue:       issue,
						Impact:      "Resources may be exposed to attacks from the Internet",
						Remediation: "Restrict source IP ranges to known trusted networks. Use Azure Bastion for management access.",
					})
				}
			}
		}
	}
}

func (a *SecurityAnalysis) analyzePublicExposure(resources []map[string]interface{}) {
	// Build map of public IPs
	publicIPs := make(map[string]string) // resource ID -> public IP
	for _, res := range resources {
		resType, _ := res["type"].(string)
		if strings.ToLower(resType) != "microsoft.network/publicipaddresses" {
			continue
		}

		id, _ := res["id"].(string)
		props, ok := res["properties"].(map[string]interface{})
		if !ok {
			continue
		}

		ipAddress, _ := props["ipAddress"].(string)
		if ipAddress != "" {
			publicIPs[id] = ipAddress
		}
	}

	// Check VMs with public IPs
	for _, res := range resources {
		resType, _ := res["type"].(string)
		if strings.ToLower(resType) != "microsoft.compute/virtualmachines" {
			continue
		}

		vmName, _ := res["name"].(string)

		// Check if VM has public IP attached via NIC
		// This would require checking NICs, for now flag all VMs with public exposure potential
		a.Findings = append(a.Findings, SecurityFinding{
			Severity:    "High",
			Category:    "PublicExposure",
			Resource:    vmName,
			Issue:       "Virtual Machine may have public IP exposure",
			Impact:      "Direct internet exposure increases attack surface",
			Remediation: "Use Azure Bastion for secure remote access instead of public IPs",
		})
	}
}

func (a *SecurityAnalysis) analyzeEncryption(resources []map[string]interface{}) {
	for _, res := range resources {
		resType, _ := res["type"].(string)
		resName, _ := res["name"].(string)

		switch strings.ToLower(resType) {
		case "microsoft.storage/storageaccounts":
			props, ok := res["properties"].(map[string]interface{})
			if !ok {
				continue
			}

			encryption, ok := props["encryption"].(map[string]interface{})
			if !ok {
				a.Findings = append(a.Findings, SecurityFinding{
					Severity:    "Medium",
					Category:    "Encryption",
					Resource:    resName,
					Issue:       "Storage account encryption status unclear",
					Impact:      "Data at rest may not be encrypted",
					Remediation: "Enable storage account encryption with customer-managed keys",
				})
				continue
			}

			services, _ := encryption["services"].(map[string]interface{})
			if services == nil {
				a.Findings = append(a.Findings, SecurityFinding{
					Severity:    "Low",
					Category:    "Encryption",
					Resource:    resName,
					Issue:       "Storage encryption services not configured",
					Impact:      "Some storage services may not be encrypted",
					Remediation: "Enable encryption for Blob, File, Table, and Queue services",
				})
			}

		case "microsoft.compute/disks":
			props, ok := res["properties"].(map[string]interface{})
			if !ok {
				continue
			}

			encryptionSettings, ok := props["encryptionSettingsCollection"].(map[string]interface{})
			if !ok || encryptionSettings == nil {
				a.Findings = append(a.Findings, SecurityFinding{
					Severity:    "Low",
					Category:    "Encryption",
					Resource:    resName,
					Issue:       "Disk encryption not configured",
					Impact:      "Disk data at rest is not encrypted with customer-managed keys",
					Remediation: "Enable Azure Disk Encryption (ADE) or use encryption at host",
				})
			}
		}
	}
}

func (a *SecurityAnalysis) analyzeNetworkIsolation(resources []map[string]interface{}) {
	// Check for resources without private endpoints or service endpoints
	for _, res := range resources {
		resType, _ := res["type"].(string)
		resName, _ := res["name"].(string)

		switch strings.ToLower(resType) {
		case "microsoft.storage/storageaccounts":
			props, ok := res["properties"].(map[string]interface{})
			if !ok {
				continue
			}

			// Check for public network access
			networkAcls, ok := props["networkAcls"].(map[string]interface{})
			if ok {
				defaultAction, _ := networkAcls["defaultAction"].(string)
				if defaultAction == "Allow" {
					a.Findings = append(a.Findings, SecurityFinding{
						Severity:    "Medium",
						Category:    "NetworkIsolation",
						Resource:    resName,
						Issue:       "Storage account allows public network access",
						Impact:      "Data can be accessed from any network",
						Remediation: "Configure network ACLs to deny by default and use private endpoints",
					})
				}
			}
		}
	}
}

func isDangerousPort(portRange string) bool {
	dangerousPorts := []string{"22", "3389", "1433", "3306", "5432", "27017", "6379"}
	for _, port := range dangerousPorts {
		if strings.Contains(portRange, port) {
			return true
		}
	}
	return false
}

// GetSecurityScore calculates overall security score (0-100)
func (a *SecurityAnalysis) GetSecurityScore() int {
	if len(a.Findings) == 0 {
		return 100
	}

	// Weighted scoring: Critical=-20, High=-10, Medium=-5, Low=-2
	score := 100
	score -= a.CriticalCount * 20
	score -= a.HighCount * 10
	score -= a.MediumCount * 5
	score -= a.LowCount * 2

	if score < 0 {
		score = 0
	}

	return score
}

// GetSecurityPosture returns a human-readable security posture
func (a *SecurityAnalysis) GetSecurityPosture() string {
	score := a.GetSecurityScore()

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
