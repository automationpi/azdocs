package models

import "time"

// Subscription represents an Azure subscription
type Subscription struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	TenantID      string    `json:"tenantId"`
	State         string    `json:"state"`
	Locations     []string  `json:"locations,omitempty"`
	ScannedAt     time.Time `json:"scannedAt"`
	ResourceCount int       `json:"resourceCount"`
}

// Metadata contains provenance information
type Metadata struct {
	Timestamp     time.Time         `json:"timestamp"`
	ToolVersion   string            `json:"toolVersion"`
	Queries       []QueryInfo       `json:"queries,omitempty"`
	SubscriptionID string           `json:"subscriptionId"`
}

// QueryInfo tracks which queries were used
type QueryInfo struct {
	Type      string    `json:"type"`      // e.g., "ResourceGraph", "ARM"
	Query     string    `json:"query"`
	Timestamp time.Time `json:"timestamp"`
	Resources int       `json:"resources"`
}
