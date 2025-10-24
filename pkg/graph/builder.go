package graph

import (
	"encoding/json"
	"fmt"
	"os"
)

// Builder constructs topology graphs from normalized data
type Builder struct {
	data interface{}
}

// NewBuilder creates a new graph builder
func NewBuilder(data interface{}) *Builder {
	return &Builder{data: data}
}

// Build constructs the topology graph
func (b *Builder) Build() (*Topology, error) {
	// TODO: Implement actual graph building logic
	return &Topology{}, nil
}

// Topology represents the network topology graph
type Topology struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Node represents a network node (VNet, Subnet, etc.)
type Node struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// Edge represents a connection between nodes
type Edge struct {
	From string                 `json:"from"`
	To   string                 `json:"to"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data,omitempty"`
}

// SaveToFile saves the topology to a file
func (t *Topology) SaveToFile(path string) error {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadNormalizedData loads previously normalized data
func LoadNormalizedData(dir string) (interface{}, error) {
	// TODO: Implement loading logic
	metaPath := fmt.Sprintf("%s/metadata.json", dir)
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}
