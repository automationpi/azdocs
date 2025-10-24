package renderer

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/automationpi/azdocs/pkg/graph"
	"github.com/automationpi/azdocs/pkg/llm"
)

// DiagramConfig holds diagram renderer configuration
type DiagramConfig struct {
	OutputDir string
	Theme     string
	EnableAI  bool
	OpenAIKey string
}

// DiagramRenderer generates Draw.io diagrams
type DiagramRenderer struct {
	config    DiagramConfig
	llmClient *llm.Client
}

// NewDiagramRenderer creates a new diagram renderer
func NewDiagramRenderer(config DiagramConfig) *DiagramRenderer {
	// Initialize LLM client if AI is enabled
	llmClient := llm.NewClient(llm.Config{
		Enabled: config.EnableAI,
		APIKey:  config.OpenAIKey,
		Model:   "", // Use default model
	})

	return &DiagramRenderer{
		config:    config,
		llmClient: llmClient,
	}
}

// RenderAll generates all diagrams for the topology
func (r *DiagramRenderer) RenderAll(topology *graph.Topology) error {
	if err := os.MkdirAll(r.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create diagrams directory: %w", err)
	}

	// Load resources to generate diagrams
	resources, err := r.loadResources()
	if err != nil {
		return fmt.Errorf("failed to load resources: %w", err)
	}

	// Use LLM for intelligent connection discovery if enabled
	var aiConnections []llm.ConnectionSuggestion
	if r.llmClient.IsEnabled() {
		fmt.Println("  🤖 Analyzing resources with AI for intelligent connections...")
		ctx := context.Background()
		aiConnections, err = r.llmClient.DiscoverConnections(ctx, resources)
		if err != nil {
			fmt.Printf("  ⚠️  AI connection discovery failed: %v\n", err)
			fmt.Println("  ℹ️  Falling back to rule-based connections")
		} else {
			fmt.Printf("  ✅ AI discovered %d connections\n", len(aiConnections))
		}
	}

	// Generate a diagram for each VNet
	vnets := r.filterByType(resources, "microsoft.network/virtualnetworks")
	for _, vnet := range vnets {
		if err := r.generateVNetDiagram(vnet, resources, aiConnections); err != nil {
			return fmt.Errorf("failed to generate diagram for %s: %w", vnet["name"], err)
		}
	}

	// Generate an overview diagram
	if err := r.generateOverviewDiagram(resources); err != nil {
		return fmt.Errorf("failed to generate overview diagram: %w", err)
	}

	return nil
}

// loadResources loads resources from JSON file
func (r *DiagramRenderer) loadResources() ([]map[string]interface{}, error) {
	data, err := os.ReadFile("./data/raw/all-resources.json")
	if err != nil {
		return nil, err
	}

	var resources []map[string]interface{}
	if err := json.Unmarshal(data, &resources); err != nil {
		return nil, err
	}

	return resources, nil
}

// filterByType filters resources by type
func (r *DiagramRenderer) filterByType(resources []map[string]interface{}, resourceType string) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, res := range resources {
		if resType, ok := res["type"].(string); ok {
			if resType == resourceType {
				filtered = append(filtered, res)
			}
		}
	}
	return filtered
}

// generateVNetDiagram generates a diagram for a single VNet
func (r *DiagramRenderer) generateVNetDiagram(vnet map[string]interface{}, allResources []map[string]interface{}, aiConnections []llm.ConnectionSuggestion) error {
	vnetName := vnet["name"].(string)
	vnetRG := vnet["resourceGroup"].(string)

	// Create Draw.io XML structure
	diagram := NewDrawIODiagram(vnetName)

	// Add VNet container (larger to accommodate subnets)
	diagram.AddAzureContainer(
		vnetName,
		fmt.Sprintf("Virtual Network: %s\nResource Group: %s", vnetName, vnetRG),
		"vnet",
		50, 50, 1100, 600,
	)

	// Find related resources in this VNet's resource group
	relatedResources := r.filterByResourceGroup(allResources, vnetRG)

	// Extract and display subnets from VNet properties
	subnets := r.extractSubnets(vnet)
	subnetY := 120.0
	subnetIDs := make(map[string]string)
	subnetPositions := make(map[string]map[string]float64) // Track x,y for each subnet

	for i, subnet := range subnets {
		subnetName := subnet["name"]
		subnetPrefix := subnet["addressPrefix"]
		subnetID := diagram.AddAzureContainer(
			fmt.Sprintf("subnet-%d", i),
			fmt.Sprintf("Subnet: %s\n%s", subnetName, subnetPrefix),
			"subnet",
			100, subnetY, 350, 120,
		)
		subnetIDs[subnetName] = subnetID

		// Store subnet position for placing resources inside
		subnetPositions[subnetName] = map[string]float64{
			"x":      100,
			"y":      subnetY,
			"width":  350,
			"height": 120,
		}

		subnetY += 140
	}

	// Add NSGs with Azure icons - positioned to the right of subnets
	// Don't connect to VNet as it creates messy lines
	nsgs := r.filterByType(relatedResources, "microsoft.network/networksecuritygroups")
	nsgY := 150.0
	nsgIDs := make(map[string]string)
	for i, nsg := range nsgs {
		nsgName := nsg["name"].(string)
		nsgID := diagram.AddAzureIcon(
			fmt.Sprintf("nsg-%d", i),
			nsgName,
			"nsg",
			520, nsgY,
		)
		nsgIDs[nsgName] = nsgID
		nsgY += 100
	}

	// Add Route Tables with Azure icons
	rts := r.filterByType(relatedResources, "microsoft.network/routetables")
	rtY := 150.0
	rtIDs := make(map[string]string)
	for i, rt := range rts {
		rtName := rt["name"].(string)
		rtID := diagram.AddAzureIcon(
			fmt.Sprintf("rt-%d", i),
			rtName,
			"routetable",
			650, rtY,
		)
		rtIDs[rtName] = rtID
		rtY += 100
	}

	// Add Load Balancers with Azure icons
	lbs := r.filterByType(relatedResources, "microsoft.network/loadbalancers")
	lbY := 150.0
	for i, lb := range lbs {
		lbName := lb["name"].(string)
		diagram.AddAzureIcon(
			fmt.Sprintf("lb-%d", i),
			lbName,
			"loadbalancer",
			780, lbY,
		)
		lbY += 100
	}

	// Add NAT Gateways with Azure icons
	nats := r.filterByType(relatedResources, "microsoft.network/natgateways")
	natY := 150.0
	natIDs := make(map[string]string)
	for i, nat := range nats {
		natName := nat["name"].(string)
		natID := diagram.AddAzureIcon(
			fmt.Sprintf("nat-%d", i),
			natName,
			"natgateway",
			910, natY,
		)
		natIDs[natName] = natID
		natY += 100
	}

	// Add Public IPs with Azure icons
	pips := r.filterByType(relatedResources, "microsoft.network/publicipaddresses")
	pipY := 400.0
	for i, pip := range pips {
		pipName := pip["name"].(string)
		pipID := diagram.AddAzureIcon(
			fmt.Sprintf("pip-%d", i),
			pipName,
			"publicip",
			1040, pipY,
		)

		// Try to connect Public IP to NAT Gateway or Load Balancer
		// This would require parsing properties to find associations
		// For now, just display them
		_ = pipID
		pipY += 80
	}

	// Add Virtual Machines - place inside management subnet
	vms := r.filterByType(relatedResources, "microsoft.compute/virtualmachines")
	vmIDMap := make(map[string]string)
	for i, vm := range vms {
		vmName := vm["name"].(string)

		// Try to place VM in management subnet, otherwise use first available
		var vmX, vmY float64
		if pos, ok := subnetPositions["subnet-management"]; ok {
			vmX = pos["x"] + 20
			vmY = pos["y"] + 40
		} else if pos, ok := subnetPositions["subnet-services"]; ok {
			vmX = pos["x"] + 20
			vmY = pos["y"] + 40
		} else {
			vmX = 200
			vmY = 300
		}

		vmID := diagram.AddAzureIcon(
			fmt.Sprintf("vm-%d", i),
			vmName,
			"vm",
			vmX, vmY,
		)
		vmIDMap[vmName] = vmID
	}

	// Add Network Interfaces - place near VMs
	nics := r.filterByType(relatedResources, "microsoft.network/networkinterfaces")
	nicIDMap := make(map[string]string)
	nicIndex := 0
	for i, nic := range nics {
		nicName := nic["name"].(string)

		// Place NIC in management subnet
		var nicX, nicY float64
		if pos, ok := subnetPositions["subnet-management"]; ok {
			nicX = pos["x"] + 100 + float64(nicIndex*80)
			nicY = pos["y"] + 60
		} else {
			nicX = 280 + float64(i*80)
			nicY = 450
		}

		nicID := diagram.AddAzureIcon(
			fmt.Sprintf("nic-%d", i),
			nicName,
			"nic",
			nicX, nicY,
		)
		nicIDMap[nicName] = nicID
		nicIndex++

		// Connect NIC to VM if VM exists (fallback if AI not enabled)
		if len(aiConnections) == 0 {
			for vmName, vmID := range vmIDMap {
				if strings.Contains(nicName, strings.Split(vmName, "-")[0]) {
					diagram.AddConnection(
						fmt.Sprintf("conn-nic-vm-%d", i),
						"attached to",
						nicID,
						vmID,
						"association",
					)
				}
			}
		}
	}

	// Add Function Apps - place in services/apps subnet
	funcApps := r.filterByType(relatedResources, "microsoft.web/sites")
	funcAppIDMap := make(map[string]string)
	for i, funcApp := range funcApps {
		funcName := funcApp["name"].(string)

		// Try to place in apps/services subnet
		var funcX, funcY float64
		if pos, ok := subnetPositions["subnet-apps"]; ok {
			funcX = pos["x"] + 20
			funcY = pos["y"] + 40
		} else if pos, ok := subnetPositions["subnet-services"]; ok {
			funcX = pos["x"] + 20
			funcY = pos["y"] + 40
		} else {
			funcX = 520
			funcY = 300
		}

		funcID := diagram.AddAzureIcon(
			fmt.Sprintf("func-%d", i),
			funcName,
			"functionapp",
			funcX, funcY,
		)
		funcAppIDMap[funcName] = funcID

		// Connect Function App to storage account if exists (fallback if AI not enabled)
		if len(aiConnections) == 0 {
			storageAccts := r.filterByType(relatedResources, "microsoft.storage/storageaccounts")
			for j, storage := range storageAccts {
				storageName := storage["name"].(string)
				// Check if storage is for this function app
				if strings.Contains(storageName, "stfunc") || strings.Contains(funcName, "func") {
					storageID := fmt.Sprintf("cell-storage-%d", j)
					diagram.AddConnection(
						fmt.Sprintf("conn-func-storage-%d", i),
						"uses",
						funcID,
						storageID,
						"association",
					)
					break // Only connect to one storage account
				}
			}
		}
	}

	// Add Storage Accounts
	storageAccts := r.filterByType(relatedResources, "microsoft.storage/storageaccounts")
	storageY := 450.0
	storageIDMap := make(map[string]string)
	for i, storage := range storageAccts {
		storageName := storage["name"].(string)
		storageID := diagram.AddAzureIcon(
			fmt.Sprintf("storage-%d", i),
			storageName,
			"storage",
			520, storageY,
		)
		storageIDMap[storageName] = storageID
		storageY += 70
	}

	// Apply AI-discovered connections if available
	if len(aiConnections) > 0 {
		fmt.Printf("    🤖 Applying %d AI-discovered connections to %s...\n", len(aiConnections), vnetName)

		// Build a map of resource names to their cell IDs
		resourceIDMap := make(map[string]string)
		for name, id := range vmIDMap {
			resourceIDMap[name] = id
		}
		for name, id := range nicIDMap {
			resourceIDMap[name] = id
		}
		for name, id := range storageIDMap {
			resourceIDMap[name] = id
		}
		for name, id := range funcAppIDMap {
			resourceIDMap[name] = id
		}
		for name, id := range nsgIDs {
			resourceIDMap[name] = id
		}
		for name, id := range natIDs {
			resourceIDMap[name] = id
		}

		// Apply high-confidence AI connections
		connCount := 0
		for _, conn := range aiConnections {
			if conn.Confidence != "high" {
				continue
			}

			// Check if both resources exist in this diagram
			sourceID, sourceExists := resourceIDMap[conn.SourceResource]
			targetID, targetExists := resourceIDMap[conn.TargetResource]

			if sourceExists && targetExists {
				diagram.AddConnection(
					fmt.Sprintf("ai-conn-%d", connCount),
					conn.Label,
					sourceID,
					targetID,
					conn.ConnectionType,
				)
				connCount++
				fmt.Printf("      ✅ %s → %s (%s)\n", conn.SourceResource, conn.TargetResource, conn.Label)
			}
		}

		if connCount > 0 {
			fmt.Printf("    ✨ Applied %d high-confidence AI connections\n", connCount)
		}
	}

	// Save diagram (removed legend as requested)
	filename := filepath.Join(r.config.OutputDir, fmt.Sprintf("%s.drawio", vnetName))
	return diagram.Save(filename)
}

// extractSubnets extracts subnet information from VNet properties
func (r *DiagramRenderer) extractSubnets(vnet map[string]interface{}) []map[string]string {
	var subnets []map[string]string

	// Try to extract properties.subnets
	if props, ok := vnet["properties"].(map[string]interface{}); ok {
		if subnetsArray, ok := props["subnets"].([]interface{}); ok {
			for _, subnetIface := range subnetsArray {
				if subnetMap, ok := subnetIface.(map[string]interface{}); ok {
					subnet := make(map[string]string)

					if name, ok := subnetMap["name"].(string); ok {
						subnet["name"] = name
					}

					if subnetProps, ok := subnetMap["properties"].(map[string]interface{}); ok {
						if prefix, ok := subnetProps["addressPrefix"].(string); ok {
							subnet["addressPrefix"] = prefix
						}
					}

					if subnet["name"] != "" {
						subnets = append(subnets, subnet)
					}
				}
			}
		}
	}

	return subnets
}

// generateOverviewDiagram generates an overview diagram of all resources
func (r *DiagramRenderer) generateOverviewDiagram(resources []map[string]interface{}) error {
	diagram := NewDrawIODiagram("Azure-Overview")

	// Add title
	diagram.AddRectangle(
		"title",
		"Azure Subscription Overview",
		50, 20, 1000, 40,
		"text;html=1;strokeColor=none;fillColor=none;align=center;verticalAlign=middle;whiteSpace=wrap;rounded=0;fontSize=20;fontStyle=1",
	)

	// Group resources by type
	vnets := r.filterByType(resources, "microsoft.network/virtualnetworks")
	nsgs := r.filterByType(resources, "microsoft.network/networksecuritygroups")
	rts := r.filterByType(resources, "microsoft.network/routetables")
	lbs := r.filterByType(resources, "microsoft.network/loadbalancers")
	nats := r.filterByType(resources, "microsoft.network/natgateways")
	pips := r.filterByType(resources, "microsoft.network/publicipaddresses")

	yOffset := 80.0
	xOffset := 100.0

	// Add VNets with Azure icons
	for i, vnet := range vnets {
		vnetName := vnet["name"].(string)
		diagram.AddAzureIcon(
			fmt.Sprintf("vnet-%d", i),
			vnetName,
			"vnet",
			xOffset, yOffset,
		)
		xOffset += 100
	}

	// Add resource type sections
	yOffset += 120
	xOffset = 100

	// Add NSGs
	if len(nsgs) > 0 {
		// Section label
		diagram.AddRectangle(
			"nsg-label",
			fmt.Sprintf("Network Security Groups (%d)", len(nsgs)),
			xOffset, yOffset, 200, 30,
			"text;html=1;strokeColor=none;fillColor=none;align=left;verticalAlign=middle;whiteSpace=wrap;rounded=0;fontSize=14;fontStyle=1",
		)
		yOffset += 40
		for i, nsg := range nsgs {
			nsgName := nsg["name"].(string)
			diagram.AddAzureIcon(
				fmt.Sprintf("nsg-%d", i),
				nsgName,
				"nsg",
				xOffset, yOffset,
			)
			xOffset += 80
			if (i+1)%10 == 0 { // Wrap after 10 items
				xOffset = 100
				yOffset += 70
			}
		}
		yOffset += 80
		xOffset = 100
	}

	// Add Route Tables
	if len(rts) > 0 {
		diagram.AddRectangle(
			"rt-label",
			fmt.Sprintf("Route Tables (%d)", len(rts)),
			xOffset, yOffset, 200, 30,
			"text;html=1;strokeColor=none;fillColor=none;align=left;verticalAlign=middle;whiteSpace=wrap;rounded=0;fontSize=14;fontStyle=1",
		)
		yOffset += 40
		for i, rt := range rts {
			rtName := rt["name"].(string)
			diagram.AddAzureIcon(
				fmt.Sprintf("rt-%d", i),
				rtName,
				"routetable",
				xOffset, yOffset,
			)
			xOffset += 80
			if (i+1)%10 == 0 {
				xOffset = 100
				yOffset += 70
			}
		}
		yOffset += 80
		xOffset = 100
	}

	// Add Load Balancers
	if len(lbs) > 0 {
		diagram.AddRectangle(
			"lb-label",
			fmt.Sprintf("Load Balancers (%d)", len(lbs)),
			xOffset, yOffset, 200, 30,
			"text;html=1;strokeColor=none;fillColor=none;align=left;verticalAlign=middle;whiteSpace=wrap;rounded=0;fontSize=14;fontStyle=1",
		)
		yOffset += 40
		for i, lb := range lbs {
			lbName := lb["name"].(string)
			diagram.AddAzureIcon(
				fmt.Sprintf("lb-%d", i),
				lbName,
				"loadbalancer",
				xOffset, yOffset,
			)
			xOffset += 80
			if (i+1)%10 == 0 {
				xOffset = 100
				yOffset += 70
			}
		}
		yOffset += 80
		xOffset = 100
	}

	// Add NAT Gateways
	if len(nats) > 0 {
		diagram.AddRectangle(
			"nat-label",
			fmt.Sprintf("NAT Gateways (%d)", len(nats)),
			xOffset, yOffset, 200, 30,
			"text;html=1;strokeColor=none;fillColor=none;align=left;verticalAlign=middle;whiteSpace=wrap;rounded=0;fontSize=14;fontStyle=1",
		)
		yOffset += 40
		for i, nat := range nats {
			natName := nat["name"].(string)
			diagram.AddAzureIcon(
				fmt.Sprintf("nat-%d", i),
				natName,
				"natgateway",
				xOffset, yOffset,
			)
			xOffset += 80
			if (i+1)%10 == 0 {
				xOffset = 100
				yOffset += 70
			}
		}
		yOffset += 80
		xOffset = 100
	}

	// Add Public IPs
	if len(pips) > 0 {
		diagram.AddRectangle(
			"pip-label",
			fmt.Sprintf("Public IP Addresses (%d)", len(pips)),
			xOffset, yOffset, 200, 30,
			"text;html=1;strokeColor=none;fillColor=none;align=left;verticalAlign=middle;whiteSpace=wrap;rounded=0;fontSize=14;fontStyle=1",
		)
		yOffset += 40
		for i, pip := range pips {
			pipName := pip["name"].(string)
			diagram.AddAzureIcon(
				fmt.Sprintf("pip-%d", i),
				pipName,
				"publicip",
				xOffset, yOffset,
			)
			xOffset += 80
			if (i+1)%10 == 0 {
				xOffset = 100
				yOffset += 70
			}
		}
	}

	// Add VNet peering connections if there are multiple VNets
	if len(vnets) >= 2 {
		// For simplicity, connect first two VNets (hub-spoke pattern)
		diagram.AddConnection(
			"peering-0-1",
			"VNet Peering",
			"cell-vnet-0",
			"cell-vnet-1",
			"peering",
		)
	}

	// Save diagram (removed legend)
	filename := filepath.Join(r.config.OutputDir, "Overview.drawio")
	return diagram.Save(filename)
}

// filterByResourceGroup filters resources by resource group
func (r *DiagramRenderer) filterByResourceGroup(resources []map[string]interface{}, resourceGroup string) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, res := range resources {
		if rg, ok := res["resourceGroup"].(string); ok {
			if rg == resourceGroup {
				filtered = append(filtered, res)
			}
		}
	}
	return filtered
}

// DrawIODiagram represents a Draw.io diagram
type DrawIODiagram struct {
	XMLName xml.Name       `xml:"mxfile"`
	Host    string         `xml:"host,attr"`
	Type    string         `xml:"type,attr"`
	Version string         `xml:"version,attr"`
	Diagram DrawIODiagramPage `xml:"diagram"`
}

// DrawIODiagramPage represents a diagram page
type DrawIODiagramPage struct {
	XMLName xml.Name      `xml:"diagram"`
	Name    string        `xml:"name,attr"`
	ID      string        `xml:"id,attr"`
	Model   DrawIOModel   `xml:"mxGraphModel"`
}

// DrawIOModel represents the graph model
type DrawIOModel struct {
	XMLName xml.Name    `xml:"mxGraphModel"`
	Root    DrawIORoot  `xml:"root"`
}

// DrawIORoot represents the root of the graph
type DrawIORoot struct {
	XMLName xml.Name      `xml:"root"`
	Cells   []DrawIOCell  `xml:"mxCell"`
}

// DrawIOCell represents a cell in the diagram
type DrawIOCell struct {
	XMLName  xml.Name         `xml:"mxCell"`
	ID       string           `xml:"id,attr"`
	Value    string           `xml:"value,attr,omitempty"`
	Style    string           `xml:"style,attr,omitempty"`
	Vertex   string           `xml:"vertex,attr,omitempty"`
	Edge     string           `xml:"edge,attr,omitempty"`
	Source   string           `xml:"source,attr,omitempty"`
	Target   string           `xml:"target,attr,omitempty"`
	Parent   string           `xml:"parent,attr,omitempty"`
	Geometry *DrawIOGeometry  `xml:"mxGeometry,omitempty"`
}

// DrawIOGeometry represents geometry of a cell
type DrawIOGeometry struct {
	XMLName xml.Name `xml:"mxGeometry"`
	X       float64  `xml:"x,attr"`
	Y       float64  `xml:"y,attr"`
	Width   float64  `xml:"width,attr"`
	Height  float64  `xml:"height,attr"`
	As      string   `xml:"as,attr"`
}

// NewDrawIODiagram creates a new Draw.io diagram
func NewDrawIODiagram(name string) *DrawIODiagram {
	diagram := &DrawIODiagram{
		Host:    "app.diagrams.net",
		Type:    "device",
		Version: "21.6.5",
		Diagram: DrawIODiagramPage{
			Name: name,
			ID:   "diagram-1",
			Model: DrawIOModel{
				Root: DrawIORoot{
					Cells: []DrawIOCell{
						{ID: "0"},
						{ID: "1", Parent: "0"},
					},
				},
			},
		},
	}
	return diagram
}

// AddRectangle adds a rectangle to the diagram
func (d *DrawIODiagram) AddRectangle(id, value string, x, y, width, height float64, style string) string {
	cellID := fmt.Sprintf("cell-%s", id)
	cell := DrawIOCell{
		ID:     cellID,
		Value:  value,
		Style:  style,
		Vertex: "1",
		Parent: "1",
		Geometry: &DrawIOGeometry{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
			As:     "geometry",
		},
	}
	d.Diagram.Model.Root.Cells = append(d.Diagram.Model.Root.Cells, cell)
	return cellID
}

// AddAzureIcon adds an Azure icon with label to the diagram
func (d *DrawIODiagram) AddAzureIcon(id, label, iconType string, x, y float64) string {
	cellID := fmt.Sprintf("cell-%s", id)

	// Azure icon styles based on resource type using Azure 2023 icon set
	var style string
	switch iconType {
	case "vnet":
		// Azure Virtual Network icon - using shape notation
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/networking/Virtual_Networks.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	case "nsg":
		// Azure Network Security Group icon
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/networking/Network_Security_Groups.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	case "routetable":
		// Azure Route Table icon
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/networking/Route_Tables.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	case "loadbalancer":
		// Azure Load Balancer icon
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/networking/Load_Balancers.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	case "natgateway":
		// Azure NAT Gateway icon
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/networking/NAT.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	case "publicip":
		// Azure Public IP icon
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/networking/Public_IP_Addresses.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	case "vm":
		// Azure Virtual Machine icon
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/compute/Virtual_Machine.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	case "nic":
		// Azure Network Interface icon
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/networking/Network_Interfaces.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	case "functionapp":
		// Azure Function App icon
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/compute/Function_Apps.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	case "storage":
		// Azure Storage Account icon
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/storage/Storage_Accounts.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	case "appserviceplan":
		// Azure App Service Plan icon
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/compute/App_Service_Plans.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	default:
		// Generic Azure resource - use a cloud icon
		style = "aspect=fixed;html=1;points=[];align=center;image;fontSize=12;image=img/lib/azure2/general/Azure.svg;labelPosition=bottom;verticalLabelPosition=top;verticalAlign=bottom;"
	}

	cell := DrawIOCell{
		ID:     cellID,
		Value:  label,
		Style:  style,
		Vertex: "1",
		Parent: "1",
		Geometry: &DrawIOGeometry{
			X:      x,
			Y:      y,
			Width:  50,
			Height: 50,
			As:     "geometry",
		},
	}
	d.Diagram.Model.Root.Cells = append(d.Diagram.Model.Root.Cells, cell)
	return cellID
}

// AddAzureContainer adds an Azure container (like VNet or Subnet) to the diagram
func (d *DrawIODiagram) AddAzureContainer(id, label, containerType string, x, y, width, height float64) string {
	cellID := fmt.Sprintf("cell-%s", id)

	var style string
	switch containerType {
	case "vnet":
		// Azure Virtual Network container style
		style = "rounded=1;whiteSpace=wrap;html=1;fillColor=#dae8fc;strokeColor=#0078D4;strokeWidth=2;dashed=1;dashPattern=5 5;fontSize=14;fontStyle=1;"
	case "subnet":
		// Azure Subnet container style
		style = "rounded=0;whiteSpace=wrap;html=1;fillColor=#e1d5e7;strokeColor=#9673a6;strokeWidth=1;dashed=1;dashPattern=3 3;fontSize=12;"
	case "resourcegroup":
		// Azure Resource Group container style
		style = "rounded=1;whiteSpace=wrap;html=1;fillColor=#f5f5f5;strokeColor=#666666;strokeWidth=2;fontSize=16;fontStyle=1;"
	default:
		style = "rounded=1;whiteSpace=wrap;html=1;fillColor=#f5f5f5;strokeColor=#666666;"
	}

	cell := DrawIOCell{
		ID:     cellID,
		Value:  label,
		Style:  style,
		Vertex: "1",
		Parent: "1",
		Geometry: &DrawIOGeometry{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
			As:     "geometry",
		},
	}
	d.Diagram.Model.Root.Cells = append(d.Diagram.Model.Root.Cells, cell)
	return cellID
}

// AddConnection adds an arrow/edge between two cells
func (d *DrawIODiagram) AddConnection(id, label, sourceID, targetID, connectionType string) string {
	cellID := fmt.Sprintf("cell-edge-%s", id)

	var style string
	switch connectionType {
	case "peering":
		// VNet peering - bidirectional with dashed line
		style = "endArrow=classic;startArrow=classic;html=1;rounded=0;strokeWidth=2;strokeColor=#0078D4;dashed=1;dashPattern=5 5;"
	case "association":
		// NSG/Route Table association - solid arrow
		style = "endArrow=classic;html=1;rounded=0;strokeWidth=1.5;strokeColor=#666666;"
	case "routing":
		// Route table routing - dotted arrow
		style = "endArrow=classic;html=1;rounded=0;strokeWidth=1.5;strokeColor=#9673a6;dashed=1;dashPattern=3 3;"
	case "natgw":
		// NAT Gateway association
		style = "endArrow=classic;html=1;rounded=0;strokeWidth=1.5;strokeColor=#d79b00;"
	default:
		// Generic connection
		style = "endArrow=classic;html=1;rounded=0;strokeWidth=1;strokeColor=#666666;"
	}

	cell := DrawIOCell{
		ID:     cellID,
		Value:  label,
		Style:  style,
		Edge:   "1",
		Source: sourceID,
		Target: targetID,
		Parent: "1",
		Geometry: &DrawIOGeometry{
			As: "geometry",
		},
	}
	d.Diagram.Model.Root.Cells = append(d.Diagram.Model.Root.Cells, cell)
	return cellID
}

// AddLegend adds a legend to the diagram
func (d *DrawIODiagram) AddLegend(x, y float64) {
	// Legend container
	d.AddRectangle(
		"legend-container",
		"Legend",
		x, y, 200, 250,
		"rounded=1;whiteSpace=wrap;html=1;fillColor=#f5f5f5;strokeColor=#666666;strokeWidth=1;fontSize=12;fontStyle=1;align=left;verticalAlign=top;spacingLeft=5;spacingTop=5;",
	)

	legendY := y + 30
	legendItems := []struct {
		icon  string
		label string
	}{
		{"vnet", "Virtual Network"},
		{"nsg", "Network Security Group"},
		{"routetable", "Route Table"},
		{"loadbalancer", "Load Balancer"},
		{"natgateway", "NAT Gateway"},
		{"publicip", "Public IP Address"},
	}

	for i, item := range legendItems {
		// Add icon
		d.AddAzureIcon(
			fmt.Sprintf("legend-icon-%d", i),
			"",
			item.icon,
			x+10, legendY,
		)
		// Add label
		d.AddRectangle(
			fmt.Sprintf("legend-label-%d", i),
			item.label,
			x+65, legendY+10, 120, 20,
			"text;html=1;strokeColor=none;fillColor=none;align=left;verticalAlign=middle;whiteSpace=wrap;rounded=0;fontSize=10;",
		)
		legendY += 35
	}
}

// Save saves the diagram to a file
func (d *DrawIODiagram) Save(filename string) error {
	data, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal XML: %w", err)
	}

	// Add XML declaration
	xmlData := []byte(xml.Header + string(data))

	return os.WriteFile(filename, xmlData, 0644)
}
