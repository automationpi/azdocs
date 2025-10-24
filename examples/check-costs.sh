#!/bin/bash

##############################################################################
# Azure Cost Checker Script
#
# This script shows estimated and actual costs for the test resources
##############################################################################

set -e

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

RESOURCE_GROUP="${RESOURCE_GROUP:-azdoc-test-rg}"

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

echo ""
echo "=========================================="
echo "Azure Cost Estimation"
echo "=========================================="
echo ""

# Check if resource group exists
if ! az group show --name "$RESOURCE_GROUP" &> /dev/null; then
    log_warning "Resource group '$RESOURCE_GROUP' does not exist."
    log_info "No resources to check."
    exit 0
fi

log_info "Resource Group: $RESOURCE_GROUP"
echo ""

# Get resource group creation time
CREATION_TIME=$(az group show --name "$RESOURCE_GROUP" --query "properties.provisioningState" -o tsv)
log_info "Status: $CREATION_TIME"
echo ""

# List resources with their types
log_info "Resources deployed:"
echo ""

az resource list \
    --resource-group "$RESOURCE_GROUP" \
    --query "[].{Name:name, Type:type}" \
    --output table

echo ""

# Count resources by type
log_info "Resource count by type:"
echo ""

VNET_COUNT=$(az network vnet list -g "$RESOURCE_GROUP" --query "length(@)" -o tsv 2>/dev/null || echo 0)
NSG_COUNT=$(az network nsg list -g "$RESOURCE_GROUP" --query "length(@)" -o tsv 2>/dev/null || echo 0)
RT_COUNT=$(az network route-table list -g "$RESOURCE_GROUP" --query "length(@)" -o tsv 2>/dev/null || echo 0)
LB_COUNT=$(az network lb list -g "$RESOURCE_GROUP" --query "length(@)" -o tsv 2>/dev/null || echo 0)
PIP_COUNT=$(az network public-ip list -g "$RESOURCE_GROUP" --query "length(@)" -o tsv 2>/dev/null || echo 0)
NAT_COUNT=$(az network nat gateway list -g "$RESOURCE_GROUP" --query "length(@)" -o tsv 2>/dev/null || echo 0)

echo "  Virtual Networks: $VNET_COUNT"
echo "  Network Security Groups: $NSG_COUNT"
echo "  Route Tables: $RT_COUNT"
echo "  Load Balancers: $LB_COUNT"
echo "  Public IPs: $PIP_COUNT"
echo "  NAT Gateways: $NAT_COUNT"

echo ""
echo "=========================================="
echo "Estimated Daily Costs"
echo "=========================================="
echo ""

TOTAL_DAILY=0

# VNets, NSGs, Route Tables are free
log_success "Virtual Networks ($VNET_COUNT): FREE"
log_success "Network Security Groups ($NSG_COUNT): FREE"
log_success "Route Tables ($RT_COUNT): FREE"

# Basic Load Balancer is free
if [ "$LB_COUNT" -gt 0 ]; then
    log_success "Basic Load Balancer: FREE"
fi

# Public IPs (Basic) - $0.0036/hour = ~$0.09/day each
if [ "$PIP_COUNT" -gt 0 ]; then
    PIP_COST=$(echo "$PIP_COUNT * 0.09" | bc)
    TOTAL_DAILY=$(echo "$TOTAL_DAILY + $PIP_COST" | bc)
    echo "  Public IPs (Basic): \$${PIP_COST}/day"
fi

# NAT Gateway - $0.045/hour = ~$1.08/day
if [ "$NAT_COUNT" -gt 0 ]; then
    # Fixed charge per gateway
    NAT_FIXED=$(echo "$NAT_COUNT * 1.08" | bc)
    TOTAL_DAILY=$(echo "$TOTAL_DAILY + $NAT_FIXED" | bc)
    echo "  NAT Gateway (fixed): \$${NAT_FIXED}/day"

    # Standard Public IP for NAT
    NAT_PIP=$(echo "$NAT_COUNT * 0.15" | bc)
    TOTAL_DAILY=$(echo "$TOTAL_DAILY + $NAT_PIP" | bc)
    echo "  NAT Gateway Public IP: \$${NAT_PIP}/day"

    log_warning "NAT Gateway data processing: ~\$0.045/GB processed"
fi

# VNet peering - charged per GB (minimal without VMs)
log_warning "VNet Peering: \$0.01/GB transferred (minimal without VMs)"

echo ""
echo "=========================================="
log_info "Estimated total: ~\$${TOTAL_DAILY}/day"
echo "=========================================="
echo ""

log_warning "Note: Costs are estimates based on East US pricing"
log_warning "Actual costs may vary based on region and usage"
echo ""

log_info "To see actual costs:"
echo "  az consumption usage list --start-date YYYY-MM-DD --end-date YYYY-MM-DD"
echo ""
log_info "Or visit Azure Portal:"
echo "  Cost Management + Billing → Cost Analysis"
echo ""

log_info "To stop incurring charges, run:"
echo "  ./cleanup-test-resources.sh"
echo ""
