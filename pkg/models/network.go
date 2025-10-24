package models

// VNet represents a Virtual Network
type VNet struct {
	ID                   string          `json:"id"`
	Name                 string          `json:"name"`
	Location             string          `json:"location"`
	ResourceGroup        string          `json:"resourceGroup"`
	AddressSpaces        []string        `json:"addressSpaces"`
	DNSServers           []string        `json:"dnsServers,omitempty"`
	Subnets              []Subnet        `json:"subnets"`
	Peerings             []VNetPeering   `json:"peerings,omitempty"`
	Tags                 map[string]string `json:"tags,omitempty"`
	EnableDDoSProtection bool            `json:"enableDDoSProtection"`
}

// Subnet represents a subnet within a VNet
type Subnet struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	AddressPrefix        string   `json:"addressPrefix"`
	NSGRef               string   `json:"nsgRef,omitempty"`        // Reference to NSG ID
	RouteTableRef        string   `json:"routeTableRef,omitempty"` // Reference to RouteTable ID
	PrivateEndpoints     []string `json:"privateEndpoints,omitempty"`
	ServiceEndpoints     []string `json:"serviceEndpoints,omitempty"`
	Delegations          []string `json:"delegations,omitempty"`
	NATGatewayRef        string   `json:"natGatewayRef,omitempty"`
	ConnectedNICs        []string `json:"connectedNICs,omitempty"` // NIC IDs
}

// VNetPeering represents a VNet peering connection
type VNetPeering struct {
	ID                        string `json:"id"`
	Name                      string `json:"name"`
	RemoteVNetID              string `json:"remoteVNetId"`
	RemoteVNetName            string `json:"remoteVNetName"`
	AllowVNetAccess           bool   `json:"allowVNetAccess"`
	AllowForwardedTraffic     bool   `json:"allowForwardedTraffic"`
	AllowGatewayTransit       bool   `json:"allowGatewayTransit"`
	UseRemoteGateways         bool   `json:"useRemoteGateways"`
	PeeringState              string `json:"peeringState"`
	RemoteAddressSpaces       []string `json:"remoteAddressSpaces,omitempty"`
}

// NetworkInterface represents a NIC
type NetworkInterface struct {
	ID                   string              `json:"id"`
	Name                 string              `json:"name"`
	Location             string              `json:"location"`
	ResourceGroup        string              `json:"resourceGroup"`
	SubnetID             string              `json:"subnetId"`
	PrivateIP            string              `json:"privateIp"`
	PrivateIPAllocation  string              `json:"privateIpAllocation"` // Static or Dynamic
	PublicIPRef          string              `json:"publicIpRef,omitempty"`
	NSGRef               string              `json:"nsgRef,omitempty"`
	EnableIPForwarding   bool                `json:"enableIpForwarding"`
	EnableAcceleratedNet bool                `json:"enableAcceleratedNetworking"`
	IPConfigurations     []IPConfiguration   `json:"ipConfigurations"`
	AttachedTo           *AttachedResource   `json:"attachedTo,omitempty"` // VM, VMSS, etc.
}

// IPConfiguration represents an IP configuration on a NIC
type IPConfiguration struct {
	Name                    string `json:"name"`
	PrivateIP               string `json:"privateIp"`
	PrivateIPAllocation     string `json:"privateIpAllocation"`
	SubnetID                string `json:"subnetId"`
	PublicIPRef             string `json:"publicIpRef,omitempty"`
	LoadBalancerBackendRefs []string `json:"loadBalancerBackendRefs,omitempty"`
	AppGatewayBackendRefs   []string `json:"appGatewayBackendRefs,omitempty"`
	Primary                 bool   `json:"primary"`
}

// AttachedResource describes what a NIC is attached to
type AttachedResource struct {
	Type string `json:"type"` // VM, VMSS, etc.
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PublicIP represents a public IP address
type PublicIP struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Location          string            `json:"location"`
	ResourceGroup     string            `json:"resourceGroup"`
	IPAddress         string            `json:"ipAddress,omitempty"`
	AllocationMethod  string            `json:"allocationMethod"` // Static or Dynamic
	SKU               string            `json:"sku"`              // Basic or Standard
	Version           string            `json:"version"`          // IPv4 or IPv6
	DNSName           string            `json:"dnsName,omitempty"`
	Tags              map[string]string `json:"tags,omitempty"`
}

// PrivateEndpoint represents a private endpoint
type PrivateEndpoint struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	Location               string   `json:"location"`
	ResourceGroup          string   `json:"resourceGroup"`
	SubnetID               string   `json:"subnetId"`
	PrivateIP              string   `json:"privateIp"`
	PrivateLinkServiceID   string   `json:"privateLinkServiceId"`
	PrivateDNSZoneGroups   []string `json:"privateDnsZoneGroups,omitempty"`
	CustomDNSConfigs       []CustomDNSConfig `json:"customDnsConfigs,omitempty"`
}

// CustomDNSConfig represents custom DNS configuration
type CustomDNSConfig struct {
	FQDN        string   `json:"fqdn"`
	IPAddresses []string `json:"ipAddresses"`
}
