#!/bin/bash

##############################################################################
# Get Azure Subscription ID Helper Script
#
# This script helps you find and use your Azure subscription ID
##############################################################################

set -e

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[TIP]${NC} $1"
}

echo ""
echo "=========================================="
echo "Azure Subscription ID Helper"
echo "=========================================="
echo ""

# Check if Azure CLI is installed
if ! command -v az &> /dev/null; then
    echo "❌ Azure CLI is not installed."
    echo ""
    echo "Install it with:"
    echo "  # Linux"
    echo "  curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash"
    echo ""
    echo "  # macOS"
    echo "  brew install azure-cli"
    echo ""
    echo "  # Windows"
    echo "  Download from: https://aka.ms/installazurecliwindows"
    echo ""
    exit 1
fi

# Check if logged in
if ! az account show &> /dev/null; then
    echo "❌ Not logged in to Azure."
    echo ""
    log_info "Please login first:"
    echo "  az login"
    echo ""
    exit 1
fi

# Get current subscription
CURRENT_SUB_ID=$(az account show --query id -o tsv)
CURRENT_SUB_NAME=$(az account show --query name -o tsv)
CURRENT_SUB_STATE=$(az account show --query state -o tsv)

log_success "You are logged in!"
echo ""

log_info "Current Active Subscription:"
echo "  Name: $CURRENT_SUB_NAME"
echo "  ID:   $CURRENT_SUB_ID"
echo "  State: $CURRENT_SUB_STATE"
echo ""

# Count total subscriptions
TOTAL_SUBS=$(az account list --query "length(@)" -o tsv)

if [ "$TOTAL_SUBS" -gt 1 ]; then
    log_info "You have $TOTAL_SUBS subscriptions available"
    echo ""
    echo "All your subscriptions:"
    echo ""
    az account list --query "[].{Name:name, ID:id, State:state, IsDefault:isDefault}" --output table
    echo ""
    log_warning "To switch subscriptions, use:"
    echo "  az account set --subscription \"<subscription-name-or-id>\""
    echo ""
fi

echo "=========================================="
log_info "Quick Commands"
echo "=========================================="
echo ""

echo "1️⃣  Export as environment variable:"
echo "   export SUBSCRIPTION_ID=\"$CURRENT_SUB_ID\""
echo ""

echo "2️⃣  Use with azdoc scan:"
echo "   ./azdoc scan --subscription-id \"$CURRENT_SUB_ID\""
echo ""

echo "3️⃣  Use with azdoc all:"
echo "   ./azdoc all --subscription-id \"$CURRENT_SUB_ID\""
echo ""

echo "4️⃣  Use in deploy script:"
echo "   export SUBSCRIPTION_ID=\"$CURRENT_SUB_ID\""
echo "   ./examples/deploy-test-resources.sh"
echo ""

echo "5️⃣  Quick one-liner (uses current subscription):"
echo "   ./azdoc scan --subscription-id \$(az account show --query id -o tsv)"
echo ""

echo "=========================================="
log_info "Copy-Paste Ready Commands"
echo "=========================================="
echo ""

cat << EOF
# Set subscription ID as environment variable
export SUBSCRIPTION_ID="$CURRENT_SUB_ID"

# Run azdoc doctor
./azdoc doctor --subscription-id "\$SUBSCRIPTION_ID"

# Deploy test resources
cd examples
export SUBSCRIPTION_ID="$CURRENT_SUB_ID"
./deploy-test-resources.sh

# Run azdoc scan and build
cd ..
./azdoc scan --subscription-id "\$SUBSCRIPTION_ID"
./azdoc build

# Cleanup test resources
cd examples
./cleanup-test-resources.sh
EOF

echo ""
echo "=========================================="
log_success "Done!"
echo "=========================================="
echo ""
