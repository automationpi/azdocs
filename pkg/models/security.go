package models

// NSG represents a Network Security Group
type NSG struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Rules         []NSGRule         `json:"rules"`
	Subnets       []string          `json:"subnets,omitempty"`       // Subnet IDs this NSG is attached to
	NICs          []string          `json:"nics,omitempty"`          // NIC IDs this NSG is attached to
	Tags          map[string]string `json:"tags,omitempty"`
}

// NSGRule represents a security rule
type NSGRule struct {
	Name                     string   `json:"name"`
	Priority                 int      `json:"priority"`
	Direction                string   `json:"direction"` // Inbound or Outbound
	Access                   string   `json:"access"`    // Allow or Deny
	Protocol                 string   `json:"protocol"`  // TCP, UDP, ICMP, *
	SourceAddressPrefixes    []string `json:"sourceAddressPrefixes"`
	SourcePortRanges         []string `json:"sourcePortRanges"`
	DestAddressPrefixes      []string `json:"destAddressPrefixes"`
	DestPortRanges           []string `json:"destPortRanges"`
	SourceASGs               []string `json:"sourceASGs,omitempty"`
	DestASGs                 []string `json:"destASGs,omitempty"`
	Description              string   `json:"description,omitempty"`
}

// ASG represents an Application Security Group
type ASG struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Tags          map[string]string `json:"tags,omitempty"`
}

// AzureFirewall represents Azure Firewall
type AzureFirewall struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`
	Location          string               `json:"location"`
	ResourceGroup     string               `json:"resourceGroup"`
	SKU               string               `json:"sku"` // Standard or Premium
	Tier              string               `json:"tier"`
	VNetID            string               `json:"vnetId"`
	SubnetID          string               `json:"subnetId"` // AzureFirewallSubnet
	ManagementSubnetID string              `json:"managementSubnetId,omitempty"`
	PublicIPs         []string             `json:"publicIps"`
	PrivateIP         string               `json:"privateIp"`
	FirewallPolicyID  string               `json:"firewallPolicyId,omitempty"`
	ThreatIntelMode   string               `json:"threatIntelMode,omitempty"`
	Rules             []FirewallRuleCollection `json:"rules,omitempty"`
	Tags              map[string]string    `json:"tags,omitempty"`
}

// FirewallRuleCollection represents a collection of firewall rules
type FirewallRuleCollection struct {
	Name     string         `json:"name"`
	Priority int            `json:"priority"`
	Action   string         `json:"action"` // Allow or Deny
	Type     string         `json:"type"`   // Network, Application, NAT
	Rules    []FirewallRule `json:"rules"`
}

// FirewallRule represents a firewall rule
type FirewallRule struct {
	Name                string   `json:"name"`
	RuleType            string   `json:"ruleType"` // Network, Application, NAT
	Protocols           []string `json:"protocols"`
	SourceAddresses     []string `json:"sourceAddresses,omitempty"`
	DestAddresses       []string `json:"destAddresses,omitempty"`
	DestPorts           []string `json:"destPorts,omitempty"`
	DestFQDNs           []string `json:"destFqdns,omitempty"`
	TargetFQDNs         []string `json:"targetFqdns,omitempty"`
	TranslatedAddress   string   `json:"translatedAddress,omitempty"`
	TranslatedPort      string   `json:"translatedPort,omitempty"`
}
