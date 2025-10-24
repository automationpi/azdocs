package models

import "time"

// RouteTable represents a route table
type RouteTable struct {
	ID                        string            `json:"id"`
	Name                      string            `json:"name"`
	Location                  string            `json:"location"`
	ResourceGroup             string            `json:"resourceGroup"`
	Routes                    []Route           `json:"routes"`
	DisableBGPRoutePropagation bool             `json:"disableBgpRoutePropagation"`
	Subnets                   []string          `json:"subnets,omitempty"` // Subnet IDs
	Tags                      map[string]string `json:"tags,omitempty"`
}

// Route represents a user-defined route
type Route struct {
	Name              string `json:"name"`
	AddressPrefix     string `json:"addressPrefix"`
	NextHopType       string `json:"nextHopType"` // VirtualNetworkGateway, VnetLocal, Internet, VirtualAppliance, None
	NextHopIPAddress  string `json:"nextHopIpAddress,omitempty"`
	HasBGPOverride    bool   `json:"hasBgpOverride"`
}

// EffectiveRoutes represents the effective routes for a NIC
type EffectiveRoutes struct {
	NICID          string          `json:"nicId"`
	NICName        string          `json:"nicName"`
	SubnetID       string          `json:"subnetId"`
	Routes         []EffectiveRoute `json:"routes"`
	SampledAt      time.Time       `json:"sampledAt"`
	TopNPrefixes   int             `json:"topNPrefixes,omitempty"` // How many routes to include
}

// EffectiveRoute represents a single effective route
type EffectiveRoute struct {
	Name                   string   `json:"name,omitempty"`
	Source                 string   `json:"source"` // Default, User, VirtualNetworkGateway, VirtualNetworkServiceEndpoint
	State                  string   `json:"state"`  // Active, Invalid
	AddressPrefixes        []string `json:"addressPrefixes"`
	NextHopType            string   `json:"nextHopType"`
	NextHopIPAddresses     []string `json:"nextHopIpAddresses,omitempty"`
	DisableBGPRouteProp    bool     `json:"disableBgpRoutePropagation"`
}

// Gateway represents VPN or ExpressRoute Gateway
type Gateway struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	Location           string            `json:"location"`
	ResourceGroup      string            `json:"resourceGroup"`
	Type               string            `json:"type"` // Vpn or ExpressRoute
	VNetID             string            `json:"vnetId"`
	SubnetID           string            `json:"subnetId"` // GatewaySubnet
	SKU                string            `json:"sku"`
	VPNType            string            `json:"vpnType,omitempty"`    // RouteBased or PolicyBased
	EnableBGP          bool              `json:"enableBgp"`
	ActiveActive       bool              `json:"activeActive"`
	PublicIPs          []string          `json:"publicIps"`
	PrivateIP          string            `json:"privateIp,omitempty"`
	Connections        []GatewayConnection `json:"connections,omitempty"`
	BGPSettings        *BGPSettings      `json:"bgpSettings,omitempty"`
	Tags               map[string]string `json:"tags,omitempty"`
}

// GatewayConnection represents a connection (VPN or ExpressRoute)
type GatewayConnection struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Type                   string `json:"type"` // Vnet2Vnet, IPsec, ExpressRoute
	ConnectionStatus       string `json:"connectionStatus"`
	LocalNetworkGatewayID  string `json:"localNetworkGatewayId,omitempty"`
	Peer                   string `json:"peer,omitempty"`
	SharedKey              bool   `json:"hasSharedKey"` // Just indicate presence, not value
	RoutingWeight          int    `json:"routingWeight,omitempty"`
	EnableBGP              bool   `json:"enableBgp"`
}

// BGPSettings represents BGP configuration
type BGPSettings struct {
	ASN               int64    `json:"asn"`
	BGPPeeringAddress string   `json:"bgpPeeringAddress"`
	PeerWeight        int      `json:"peerWeight"`
}

// NATGateway represents a NAT Gateway
type NATGateway struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Location          string            `json:"location"`
	ResourceGroup     string            `json:"resourceGroup"`
	SKU               string            `json:"sku"`
	PublicIPs         []string          `json:"publicIps"`
	PublicIPPrefixes  []string          `json:"publicIpPrefixes,omitempty"`
	IdleTimeoutMins   int               `json:"idleTimeoutMinutes"`
	Subnets           []string          `json:"subnets,omitempty"` // Subnet IDs
	Tags              map[string]string `json:"tags,omitempty"`
}
