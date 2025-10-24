#!/bin/bash

##############################################################################
# Azure Test Resources Deployment Script
#
# This script creates a comprehensive Azure environment for testing azdoc.
# Cost estimate: ~$3-5/day (if left running)
#
# Resources created:
# - 1 Resource Group
# - 2 VNets (Hub and Spoke)
# - 4 Subnets
# - 2 NSGs with sample rules
# - 1 Route Table with UDRs
# - VNet Peering (Hub-Spoke)
# - 3 Public IPs
# - 1 Basic Load Balancer (Free tier)
# - 1 NAT Gateway
# - 1 Virtual Machine (Ubuntu jumpbox)
# - 1 Azure Function App (Python, Consumption Plan)
# - 1 Storage Account (for Functions)
#
# Prerequisites:
# - Azure CLI installed and logged in (az login)
# - Contributor role on subscription
##############################################################################

set -e  # Exit on error
set -u  # Exit on undefined variable

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
RESOURCE_GROUP="${RESOURCE_GROUP:-azdoc-test-rg}"
LOCATION="${LOCATION:-eastus}"
SUBSCRIPTION_ID="${SUBSCRIPTION_ID:-}"

# VNet Configuration
HUB_VNET_NAME="vnet-hub"
HUB_VNET_PREFIX="10.0.0.0/16"
HUB_SUBNET1_NAME="subnet-services"
HUB_SUBNET1_PREFIX="10.0.1.0/24"
HUB_SUBNET2_NAME="subnet-management"
HUB_SUBNET2_PREFIX="10.0.2.0/24"

SPOKE_VNET_NAME="vnet-spoke1"
SPOKE_VNET_PREFIX="10.1.0.0/16"
SPOKE_SUBNET1_NAME="subnet-apps"
SPOKE_SUBNET1_PREFIX="10.1.1.0/24"
SPOKE_SUBNET2_NAME="subnet-data"
SPOKE_SUBNET2_PREFIX="10.1.2.0/24"

# NSG Configuration
NSG_HUB_NAME="nsg-hub"
NSG_SPOKE_NAME="nsg-spoke"

# Route Table
RT_NAME="rt-spoke-to-hub"

# Load Balancer
LB_NAME="lb-internal"
PUBLIC_IP_NAME="pip-lb"

# NAT Gateway
NAT_GW_NAME="natgw-spoke"

# VM Configuration
VM_NAME="vm-hub-jumpbox"
VM_SIZE="Standard_B1s"  # Cheapest VM size
VM_IMAGE="Ubuntu2204"
VM_ADMIN_USER="${VM_ADMIN_USER:-azureuser}"

# Azure Function Configuration
FUNC_APP_NAME="func-spoke-app-${RANDOM}"
STORAGE_ACCOUNT_NAME="stfunc${RANDOM}"
FUNC_PLAN_NAME="plan-func-spoke"

##############################################################################
# Helper Functions
##############################################################################

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."

    # Check if Azure CLI is installed
    if ! command -v az &> /dev/null; then
        log_error "Azure CLI is not installed. Please install it first."
        exit 1
    fi

    # Check if logged in
    if ! az account show &> /dev/null; then
        log_error "Not logged in to Azure. Please run 'az login' first."
        exit 1
    fi

    # Set subscription if provided
    if [ -n "$SUBSCRIPTION_ID" ]; then
        log_info "Setting subscription to: $SUBSCRIPTION_ID"
        az account set --subscription "$SUBSCRIPTION_ID"
    fi

    # Get current subscription
    CURRENT_SUB=$(az account show --query name -o tsv)
    CURRENT_SUB_ID=$(az account show --query id -o tsv)
    log_success "Using subscription: $CURRENT_SUB ($CURRENT_SUB_ID)"
}

##############################################################################
# Resource Creation Functions
##############################################################################

create_resource_group() {
    log_info "Creating resource group: $RESOURCE_GROUP in $LOCATION..."

    if az group show --name "$RESOURCE_GROUP" &> /dev/null; then
        log_warning "Resource group already exists. Skipping."
    else
        az group create \
            --name "$RESOURCE_GROUP" \
            --location "$LOCATION" \
            --tags "Purpose=azdoc-testing" "ManagedBy=azdoc-example-script" \
            --output none
        log_success "Resource group created"
    fi
}

create_hub_vnet() {
    log_info "Creating Hub VNet: $HUB_VNET_NAME..."

    az network vnet create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$HUB_VNET_NAME" \
        --address-prefix "$HUB_VNET_PREFIX" \
        --subnet-name "$HUB_SUBNET1_NAME" \
        --subnet-prefix "$HUB_SUBNET1_PREFIX" \
        --tags "Type=Hub" "Environment=Test" \
        --output none

    log_info "Creating second subnet in Hub VNet..."
    az network vnet subnet create \
        --resource-group "$RESOURCE_GROUP" \
        --vnet-name "$HUB_VNET_NAME" \
        --name "$HUB_SUBNET2_NAME" \
        --address-prefix "$HUB_SUBNET2_PREFIX" \
        --output none

    log_success "Hub VNet created with 2 subnets"
}

create_spoke_vnet() {
    log_info "Creating Spoke VNet: $SPOKE_VNET_NAME..."

    az network vnet create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$SPOKE_VNET_NAME" \
        --address-prefix "$SPOKE_VNET_PREFIX" \
        --subnet-name "$SPOKE_SUBNET1_NAME" \
        --subnet-prefix "$SPOKE_SUBNET1_PREFIX" \
        --tags "Type=Spoke" "Environment=Test" \
        --output none

    log_info "Creating second subnet in Spoke VNet..."
    az network vnet subnet create \
        --resource-group "$RESOURCE_GROUP" \
        --vnet-name "$SPOKE_VNET_NAME" \
        --name "$SPOKE_SUBNET2_NAME" \
        --address-prefix "$SPOKE_SUBNET2_PREFIX" \
        --output none

    log_success "Spoke VNet created with 2 subnets"
}

create_nsgs() {
    log_info "Creating NSG for Hub: $NSG_HUB_NAME..."

    az network nsg create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$NSG_HUB_NAME" \
        --tags "AppliesTo=HubVNet" \
        --output none

    # Add sample rules to Hub NSG
    log_info "Adding rules to Hub NSG..."

    az network nsg rule create \
        --resource-group "$RESOURCE_GROUP" \
        --nsg-name "$NSG_HUB_NAME" \
        --name "AllowSSHFromManagement" \
        --priority 100 \
        --direction Inbound \
        --access Allow \
        --protocol Tcp \
        --source-address-prefixes "10.0.2.0/24" \
        --source-port-ranges "*" \
        --destination-address-prefixes "*" \
        --destination-port-ranges 22 \
        --description "Allow SSH from management subnet" \
        --output none

    az network nsg rule create \
        --resource-group "$RESOURCE_GROUP" \
        --nsg-name "$NSG_HUB_NAME" \
        --name "AllowHTTPSFromInternet" \
        --priority 200 \
        --direction Inbound \
        --access Allow \
        --protocol Tcp \
        --source-address-prefixes "Internet" \
        --source-port-ranges "*" \
        --destination-address-prefixes "*" \
        --destination-port-ranges 443 \
        --description "Allow HTTPS from Internet" \
        --output none

    log_info "Creating NSG for Spoke: $NSG_SPOKE_NAME..."

    az network nsg create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$NSG_SPOKE_NAME" \
        --tags "AppliesTo=SpokeVNet" \
        --output none

    # Add sample rules to Spoke NSG
    log_info "Adding rules to Spoke NSG..."

    az network nsg rule create \
        --resource-group "$RESOURCE_GROUP" \
        --nsg-name "$NSG_SPOKE_NAME" \
        --name "AllowHTTPFromHub" \
        --priority 100 \
        --direction Inbound \
        --access Allow \
        --protocol Tcp \
        --source-address-prefixes "$HUB_VNET_PREFIX" \
        --source-port-ranges "*" \
        --destination-address-prefixes "*" \
        --destination-port-ranges 80 443 \
        --description "Allow HTTP/HTTPS from Hub" \
        --output none

    az network nsg rule create \
        --resource-group "$RESOURCE_GROUP" \
        --nsg-name "$NSG_SPOKE_NAME" \
        --name "DenyInternetOutbound" \
        --priority 200 \
        --direction Outbound \
        --access Deny \
        --protocol "*" \
        --source-address-prefixes "*" \
        --source-port-ranges "*" \
        --destination-address-prefixes "Internet" \
        --destination-port-ranges "*" \
        --description "Deny all Internet outbound" \
        --output none

    log_success "NSGs created with sample rules"
}

associate_nsgs() {
    log_info "Associating NSGs with subnets..."

    # Associate Hub NSG with Hub subnets
    az network vnet subnet update \
        --resource-group "$RESOURCE_GROUP" \
        --vnet-name "$HUB_VNET_NAME" \
        --name "$HUB_SUBNET1_NAME" \
        --network-security-group "$NSG_HUB_NAME" \
        --output none

    # Associate Spoke NSG with Spoke subnets
    az network vnet subnet update \
        --resource-group "$RESOURCE_GROUP" \
        --vnet-name "$SPOKE_VNET_NAME" \
        --name "$SPOKE_SUBNET1_NAME" \
        --network-security-group "$NSG_SPOKE_NAME" \
        --output none

    log_success "NSGs associated with subnets"
}

create_route_table() {
    log_info "Creating route table: $RT_NAME..."

    az network route-table create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$RT_NAME" \
        --tags "Purpose=SpokeRouting" \
        --output none

    # Add sample routes
    log_info "Adding custom routes..."

    az network route-table route create \
        --resource-group "$RESOURCE_GROUP" \
        --route-table-name "$RT_NAME" \
        --name "RouteToHub" \
        --address-prefix "$HUB_VNET_PREFIX" \
        --next-hop-type VnetLocal \
        --output none

    az network route-table route create \
        --resource-group "$RESOURCE_GROUP" \
        --route-table-name "$RT_NAME" \
        --name "DefaultRouteToInternet" \
        --address-prefix "0.0.0.0/0" \
        --next-hop-type Internet \
        --output none

    # Associate with spoke subnet
    log_info "Associating route table with spoke subnet..."
    az network vnet subnet update \
        --resource-group "$RESOURCE_GROUP" \
        --vnet-name "$SPOKE_VNET_NAME" \
        --name "$SPOKE_SUBNET1_NAME" \
        --route-table "$RT_NAME" \
        --output none

    log_success "Route table created and associated"
}

create_vnet_peering() {
    log_info "Creating VNet peering: Hub <-> Spoke..."

    # Get VNet IDs
    HUB_VNET_ID=$(az network vnet show \
        --resource-group "$RESOURCE_GROUP" \
        --name "$HUB_VNET_NAME" \
        --query id -o tsv)

    SPOKE_VNET_ID=$(az network vnet show \
        --resource-group "$RESOURCE_GROUP" \
        --name "$SPOKE_VNET_NAME" \
        --query id -o tsv)

    # Create peering: Hub -> Spoke
    az network vnet peering create \
        --resource-group "$RESOURCE_GROUP" \
        --name "peer-hub-to-spoke1" \
        --vnet-name "$HUB_VNET_NAME" \
        --remote-vnet "$SPOKE_VNET_ID" \
        --allow-vnet-access \
        --allow-forwarded-traffic \
        --output none

    # Create peering: Spoke -> Hub
    az network vnet peering create \
        --resource-group "$RESOURCE_GROUP" \
        --name "peer-spoke1-to-hub" \
        --vnet-name "$SPOKE_VNET_NAME" \
        --remote-vnet "$HUB_VNET_ID" \
        --allow-vnet-access \
        --allow-forwarded-traffic \
        --output none

    log_success "VNet peering created"
}

create_public_ip() {
    log_info "Creating public IP: $PUBLIC_IP_NAME..."

    az network public-ip create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$PUBLIC_IP_NAME" \
        --sku Basic \
        --allocation-method Static \
        --tags "Purpose=LoadBalancer" \
        --output none

    PUBLIC_IP_ADDRESS=$(az network public-ip show \
        --resource-group "$RESOURCE_GROUP" \
        --name "$PUBLIC_IP_NAME" \
        --query ipAddress -o tsv)

    log_success "Public IP created: $PUBLIC_IP_ADDRESS"
}

create_load_balancer() {
    log_info "Creating basic load balancer: $LB_NAME..."

    # Create load balancer
    az network lb create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$LB_NAME" \
        --sku Basic \
        --vnet-name "$HUB_VNET_NAME" \
        --subnet "$HUB_SUBNET1_NAME" \
        --frontend-ip-name "frontend-internal" \
        --backend-pool-name "backend-pool" \
        --output none

    # Add health probe
    az network lb probe create \
        --resource-group "$RESOURCE_GROUP" \
        --lb-name "$LB_NAME" \
        --name "health-probe-http" \
        --protocol Http \
        --port 80 \
        --path "/" \
        --output none

    # Add load balancing rule
    az network lb rule create \
        --resource-group "$RESOURCE_GROUP" \
        --lb-name "$LB_NAME" \
        --name "rule-http" \
        --protocol Tcp \
        --frontend-port 80 \
        --backend-port 80 \
        --frontend-ip-name "frontend-internal" \
        --backend-pool-name "backend-pool" \
        --probe-name "health-probe-http" \
        --output none

    log_success "Load balancer created (Basic SKU - Free tier)"
}

create_nat_gateway() {
    log_info "Creating NAT Gateway: $NAT_GW_NAME..."

    # Create public IP for NAT Gateway
    az network public-ip create \
        --resource-group "$RESOURCE_GROUP" \
        --name "${NAT_GW_NAME}-pip" \
        --sku Standard \
        --allocation-method Static \
        --output none

    # Create NAT Gateway
    az network nat gateway create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$NAT_GW_NAME" \
        --public-ip-addresses "${NAT_GW_NAME}-pip" \
        --idle-timeout 4 \
        --output none

    # Associate with spoke data subnet
    az network vnet subnet update \
        --resource-group "$RESOURCE_GROUP" \
        --vnet-name "$SPOKE_VNET_NAME" \
        --name "$SPOKE_SUBNET2_NAME" \
        --nat-gateway "$NAT_GW_NAME" \
        --output none

    log_success "NAT Gateway created and associated"
}

create_virtual_machines() {
    log_info "Creating Virtual Machine: $VM_NAME..."

    # Create a Public IP for the VM
    VM_PIP_NAME="${VM_NAME}-pip"
    log_info "Creating Public IP for VM..."
    az network public-ip create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$VM_PIP_NAME" \
        --sku Standard \
        --allocation-method Static \
        --tags "Purpose=VMAccess" \
        --output none

    # Create Network Interface
    VM_NIC_NAME="${VM_NAME}-nic"
    log_info "Creating Network Interface for VM..."
    az network nic create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$VM_NIC_NAME" \
        --vnet-name "$HUB_VNET_NAME" \
        --subnet "$HUB_SUBNET2_NAME" \
        --public-ip-address "$VM_PIP_NAME" \
        --network-security-group "$NSG_HUB_NAME" \
        --output none

    # Create the VM
    log_info "Creating Ubuntu VM (this may take a few minutes)..."
    log_warning "VM will be created with SSH key authentication"
    log_warning "If you don't have SSH keys, they will be generated in ~/.ssh/"

    az vm create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$VM_NAME" \
        --size "$VM_SIZE" \
        --image "$VM_IMAGE" \
        --admin-username "$VM_ADMIN_USER" \
        --generate-ssh-keys \
        --nics "$VM_NIC_NAME" \
        --tags "Purpose=Jumpbox" "Environment=Test" \
        --output none

    log_success "Virtual Machine created"
}

create_azure_functions() {
    log_info "Creating Azure Function App: $FUNC_APP_NAME..."

    # Create Storage Account (required for Azure Functions)
    log_info "Creating Storage Account for Functions..."
    az storage account create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$STORAGE_ACCOUNT_NAME" \
        --location "$LOCATION" \
        --sku Standard_LRS \
        --kind StorageV2 \
        --tags "Purpose=FunctionAppStorage" \
        --output none

    # Create Function App Plan (Consumption - cheapest)
    log_info "Creating Function App Service Plan..."
    az functionapp plan create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$FUNC_PLAN_NAME" \
        --location "$LOCATION" \
        --sku Y1 \
        --is-linux true \
        --output none

    # Create Function App
    log_info "Creating Function App..."
    az functionapp create \
        --resource-group "$RESOURCE_GROUP" \
        --name "$FUNC_APP_NAME" \
        --storage-account "$STORAGE_ACCOUNT_NAME" \
        --plan "$FUNC_PLAN_NAME" \
        --runtime python \
        --runtime-version 3.11 \
        --functions-version 4 \
        --os-type Linux \
        --tags "Purpose=ServerlessApp" "Environment=Test" \
        --output none

    # Enable VNet integration for the Function App
    log_info "Enabling VNet integration for Function App..."
    az functionapp vnet-integration add \
        --resource-group "$RESOURCE_GROUP" \
        --name "$FUNC_APP_NAME" \
        --vnet "$SPOKE_VNET_NAME" \
        --subnet "$SPOKE_SUBNET1_NAME" \
        --output none || log_warning "VNet integration failed (may require Premium plan)"

    log_success "Azure Function App created"
}

##############################################################################
# Main Execution
##############################################################################

main() {
    echo ""
    echo "=========================================="
    echo "Azure Test Resources Deployment"
    echo "=========================================="
    echo ""

    log_info "This script will create Azure resources for testing azdoc"
    log_warning "Estimated cost: ~\$3-5/day if left running (with VM and Functions)"
    log_warning "Don't forget to run cleanup-test-resources.sh when done!"
    echo ""

    # Ask for confirmation
    read -p "Do you want to continue? (yes/no): " -r
    echo
    if [[ ! $REPLY =~ ^[Yy]es$ ]]; then
        log_info "Deployment cancelled"
        exit 0
    fi

    # Start deployment
    START_TIME=$(date +%s)

    check_prerequisites

    echo ""
    log_info "Starting resource deployment..."
    echo ""

    create_resource_group
    create_hub_vnet
    create_spoke_vnet
    create_nsgs
    associate_nsgs
    create_route_table
    create_vnet_peering
    create_public_ip
    create_load_balancer
    create_nat_gateway
    create_virtual_machines
    create_azure_functions

    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))

    echo ""
    echo "=========================================="
    log_success "Deployment Complete!"
    echo "=========================================="
    echo ""
    log_info "Duration: ${DURATION}s"
    log_info "Resource Group: $RESOURCE_GROUP"
    log_info "Location: $LOCATION"
    log_info "Subscription: $CURRENT_SUB_ID"
    echo ""
    log_info "Resources created:"
    echo "  - 2 Virtual Networks (Hub and Spoke)"
    echo "  - 4 Subnets"
    echo "  - 2 Network Security Groups (with rules)"
    echo "  - 1 Route Table (with UDRs)"
    echo "  - VNet Peering"
    echo "  - 3 Public IPs"
    echo "  - 1 Load Balancer (Basic - Free)"
    echo "  - 1 NAT Gateway"
    echo "  - 1 Virtual Machine (Standard_B1s Ubuntu)"
    echo "  - 1 Network Interface"
    echo "  - 1 Azure Function App (Python 3.11, Consumption Plan)"
    echo "  - 1 Storage Account (for Functions)"
    echo "  - 1 App Service Plan (Consumption)"
    echo ""
    log_info "Next steps:"
    echo "  1. Run azdoc to document these resources:"
    echo "     cd .."
    echo "     ./azdoc scan --subscription-id $CURRENT_SUB_ID"
    echo "     ./azdoc build"
    echo ""
    echo "  2. When done testing, cleanup resources:"
    echo "     ./examples/cleanup-test-resources.sh"
    echo ""
}

# Run main function
main "$@"
