#!/bin/bash

##############################################################################
# Azure Test Resources Cleanup Script
#
# This script deletes all resources created by deploy-test-resources.sh
#
# WARNING: This will permanently delete the resource group and all resources!
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
SUBSCRIPTION_ID="${SUBSCRIPTION_ID:-}"

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
        log_error "Azure CLI is not installed."
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
    log_info "Using subscription: $CURRENT_SUB ($CURRENT_SUB_ID)"
}

check_resource_group_exists() {
    if ! az group show --name "$RESOURCE_GROUP" &> /dev/null; then
        log_warning "Resource group '$RESOURCE_GROUP' does not exist."
        log_info "Nothing to clean up."
        exit 0
    fi
}

list_resources() {
    log_info "Resources in '$RESOURCE_GROUP':"
    echo ""

    # List all resources
    az resource list \
        --resource-group "$RESOURCE_GROUP" \
        --query "[].{Name:name, Type:type, Location:location}" \
        --output table

    echo ""

    # Count resources
    RESOURCE_COUNT=$(az resource list --resource-group "$RESOURCE_GROUP" --query "length(@)" -o tsv)
    log_info "Total resources: $RESOURCE_COUNT"
}

delete_resource_group() {
    log_warning "Deleting resource group: $RESOURCE_GROUP"
    log_warning "This will delete ALL resources in the group!"

    # Delete without waiting (async)
    az group delete \
        --name "$RESOURCE_GROUP" \
        --yes \
        --no-wait

    log_success "Deletion initiated (running in background)"
    log_info "You can monitor progress with:"
    echo "  az group show --name $RESOURCE_GROUP"
}

##############################################################################
# Main Execution
##############################################################################

main() {
    echo ""
    echo "=========================================="
    echo "Azure Test Resources Cleanup"
    echo "=========================================="
    echo ""

    log_warning "This script will DELETE the following resource group and ALL resources:"
    log_warning "  Resource Group: $RESOURCE_GROUP"
    echo ""

    check_prerequisites
    check_resource_group_exists

    echo ""
    list_resources
    echo ""

    # Ask for confirmation
    log_error "WARNING: This action CANNOT be undone!"
    echo ""
    read -p "Type 'DELETE' to confirm deletion: " -r
    echo

    if [[ $REPLY != "DELETE" ]]; then
        log_info "Cleanup cancelled"
        exit 0
    fi

    # Additional confirmation
    read -p "Are you absolutely sure? (yes/no): " -r
    echo
    if [[ ! $REPLY =~ ^[Yy]es$ ]]; then
        log_info "Cleanup cancelled"
        exit 0
    fi

    # Start cleanup
    START_TIME=$(date +%s)

    echo ""
    log_info "Starting cleanup..."
    echo ""

    delete_resource_group

    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))

    echo ""
    echo "=========================================="
    log_success "Cleanup Initiated!"
    echo "=========================================="
    echo ""
    log_info "Duration: ${DURATION}s"
    log_info "Resource Group: $RESOURCE_GROUP"
    echo ""
    log_info "The deletion is running in the background."
    log_info "It may take 5-10 minutes to complete."
    echo ""
    log_info "Check deletion status with:"
    echo "  az group show --name $RESOURCE_GROUP"
    echo ""
    log_info "List all resource groups:"
    echo "  az group list --output table"
    echo ""
}

# Run main function
main "$@"
