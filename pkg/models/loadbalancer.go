package models

// LoadBalancer represents an Azure Load Balancer
type LoadBalancer struct {
	ID               string                  `json:"id"`
	Name             string                  `json:"name"`
	Location         string                  `json:"location"`
	ResourceGroup    string                  `json:"resourceGroup"`
	SKU              string                  `json:"sku"` // Basic or Standard
	Type             string                  `json:"type"` // Public or Internal
	FrontendIPs      []LoadBalancerFrontend  `json:"frontendIps"`
	BackendPools     []LoadBalancerBackend   `json:"backendPools"`
	LoadBalancingRules []LBRule              `json:"loadBalancingRules"`
	Probes           []LBProbe               `json:"probes"`
	InboundNATRules  []LBNATRule             `json:"inboundNatRules,omitempty"`
	OutboundRules    []LBOutboundRule        `json:"outboundRules,omitempty"`
	Tags             map[string]string       `json:"tags,omitempty"`
}

// LoadBalancerFrontend represents a frontend IP configuration
type LoadBalancerFrontend struct {
	Name            string `json:"name"`
	PrivateIP       string `json:"privateIp,omitempty"`
	PublicIPRef     string `json:"publicIpRef,omitempty"`
	SubnetID        string `json:"subnetId,omitempty"`
	Zones           []string `json:"zones,omitempty"`
}

// LoadBalancerBackend represents a backend address pool
type LoadBalancerBackend struct {
	Name      string   `json:"name"`
	NICRefs   []string `json:"nicRefs,omitempty"`   // NIC IP config IDs
	Addresses []BackendAddress `json:"addresses,omitempty"`
}

// BackendAddress represents an address in a backend pool
type BackendAddress struct {
	Name      string `json:"name"`
	IPAddress string `json:"ipAddress"`
	VNetID    string `json:"vnetId,omitempty"`
}

// LBRule represents a load balancing rule
type LBRule struct {
	Name                string `json:"name"`
	Protocol            string `json:"protocol"` // TCP or UDP
	FrontendPort        int    `json:"frontendPort"`
	BackendPort         int    `json:"backendPort"`
	FrontendIPRef       string `json:"frontendIpRef"`
	BackendPoolRef      string `json:"backendPoolRef"`
	ProbeRef            string `json:"probeRef"`
	EnableFloatingIP    bool   `json:"enableFloatingIp"`
	IdleTimeoutMins     int    `json:"idleTimeoutMinutes"`
	LoadDistribution    string `json:"loadDistribution,omitempty"` // Default, SourceIP, SourceIPProtocol
	DisableOutboundSNAT bool   `json:"disableOutboundSnat"`
}

// LBProbe represents a health probe
type LBProbe struct {
	Name            string `json:"name"`
	Protocol        string `json:"protocol"` // HTTP, HTTPS, TCP
	Port            int    `json:"port"`
	Path            string `json:"path,omitempty"`
	IntervalSecs    int    `json:"intervalSeconds"`
	NumProbes       int    `json:"numberOfProbes"`
}

// LBNATRule represents an inbound NAT rule
type LBNATRule struct {
	Name              string `json:"name"`
	Protocol          string `json:"protocol"`
	FrontendPort      int    `json:"frontendPort"`
	BackendPort       int    `json:"backendPort"`
	FrontendIPRef     string `json:"frontendIpRef"`
	NICRef            string `json:"nicRef,omitempty"`
	EnableFloatingIP  bool   `json:"enableFloatingIp"`
	IdleTimeoutMins   int    `json:"idleTimeoutMinutes"`
}

// LBOutboundRule represents an outbound rule
type LBOutboundRule struct {
	Name                 string   `json:"name"`
	Protocol             string   `json:"protocol"` // TCP, UDP, All
	FrontendIPRefs       []string `json:"frontendIpRefs"`
	BackendPoolRef       string   `json:"backendPoolRef"`
	IdleTimeoutMins      int      `json:"idleTimeoutMinutes"`
	AllocatedOutboundPorts int    `json:"allocatedOutboundPorts,omitempty"`
	EnableTCPReset       bool     `json:"enableTcpReset"`
}

// ApplicationGateway represents an Application Gateway
type ApplicationGateway struct {
	ID                  string                    `json:"id"`
	Name                string                    `json:"name"`
	Location            string                    `json:"location"`
	ResourceGroup       string                    `json:"resourceGroup"`
	SKU                 string                    `json:"sku"` // Standard_v2, WAF_v2
	Tier                string                    `json:"tier"`
	Capacity            int                       `json:"capacity"`
	VNetID              string                    `json:"vnetId"`
	SubnetID            string                    `json:"subnetId"`
	PublicIPs           []string                  `json:"publicIps,omitempty"`
	PrivateIPs          []string                  `json:"privateIps,omitempty"`
	FrontendPorts       []AppGWFrontendPort       `json:"frontendPorts"`
	BackendPools        []AppGWBackendPool        `json:"backendPools"`
	HTTPListeners       []AppGWListener           `json:"httpListeners"`
	RequestRoutingRules []AppGWRoutingRule        `json:"requestRoutingRules"`
	Probes              []AppGWProbe              `json:"probes,omitempty"`
	WAFEnabled          bool                      `json:"wafEnabled"`
	WAFMode             string                    `json:"wafMode,omitempty"` // Detection or Prevention
	Tags                map[string]string         `json:"tags,omitempty"`
}

// AppGWFrontendPort represents a frontend port
type AppGWFrontendPort struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

// AppGWBackendPool represents a backend address pool
type AppGWBackendPool struct {
	Name      string             `json:"name"`
	Addresses []AppGWBackendAddress `json:"addresses"`
}

// AppGWBackendAddress represents a backend address
type AppGWBackendAddress struct {
	FQDN      string `json:"fqdn,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
}

// AppGWListener represents an HTTP listener
type AppGWListener struct {
	Name              string `json:"name"`
	Protocol          string `json:"protocol"` // HTTP or HTTPS
	FrontendIPRef     string `json:"frontendIpRef"`
	FrontendPortRef   string `json:"frontendPortRef"`
	HostName          string `json:"hostName,omitempty"`
	RequireServerNameIndication bool `json:"requireSni"`
	SSLCertificateRef string `json:"sslCertificateRef,omitempty"`
}

// AppGWRoutingRule represents a request routing rule
type AppGWRoutingRule struct {
	Name               string `json:"name"`
	RuleType           string `json:"ruleType"` // Basic or PathBasedRouting
	Priority           int    `json:"priority"`
	ListenerRef        string `json:"listenerRef"`
	BackendPoolRef     string `json:"backendPoolRef,omitempty"`
	BackendHTTPSettings string `json:"backendHttpSettings,omitempty"`
	RedirectConfigRef  string `json:"redirectConfigRef,omitempty"`
}

// AppGWProbe represents a custom health probe
type AppGWProbe struct {
	Name            string `json:"name"`
	Protocol        string `json:"protocol"`
	Host            string `json:"host,omitempty"`
	Path            string `json:"path"`
	IntervalSecs    int    `json:"intervalSeconds"`
	TimeoutSecs     int    `json:"timeoutSeconds"`
	UnhealthyThreshold int `json:"unhealthyThreshold"`
	PickHostNameFromBackend bool `json:"pickHostNameFromBackend"`
}
