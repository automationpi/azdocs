# Azure Test Resources for azdoc

This directory contains scripts to create and destroy Azure resources for testing the azdoc tool.

## Overview

The test environment creates a realistic Azure network setup with:
- **2 Virtual Networks** (Hub-and-Spoke topology)
- **4 Subnets** (2 per VNet)
- **2 Network Security Groups** with security rules
- **1 Route Table** with custom routes
- **VNet Peering** between Hub and Spoke
- **1 Load Balancer** (Basic SKU - Free tier)
- **1 NAT Gateway**
- **2 Public IPs**

## Cost Estimate

**Estimated cost: $1-2 per day** if left running

Cost breakdown:
- VNets, Subnets, NSGs: **FREE**
- Route Tables: **FREE**
- VNet Peering: **$0.01/GB** (minimal if no traffic)
- Basic Load Balancer: **FREE**
- NAT Gateway: **~$1.00/day** (fixed charge) + $0.045/GB processed
- Public IPs: **~$0.10/day** ($3/month)
- Standard Public IP (NAT GW): **~$0.15/day** ($4.50/month)

**💰 Total estimated: ~$1.25-2.00 per day**

> **Note**: Costs are minimal because we're not deploying any compute resources (VMs, containers, etc.)

## Prerequisites

1. **Azure CLI** installed
   ```bash
   # Install on Linux
   curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash

   # Install on macOS
   brew install azure-cli

   # Install on Windows
   # Download from: https://aka.ms/installazurecliwindows
   ```

2. **Azure Account** with active subscription
   - Sign up: https://azure.microsoft.com/free/
   - Free tier includes $200 credit for 30 days

3. **Permissions**: Contributor role on the subscription

## Quick Start

### Step 1: Login to Azure

```bash
# Login interactively
az login

# List your subscriptions
az account list --output table

# Set the subscription you want to use (optional)
az account set --subscription "Your Subscription Name"
```

**💡 Not sure about your subscription ID?** Run this helper:
```bash
./get-subscription-id.sh
```

This will show your current subscription ID and provide copy-paste ready commands!

### Step 2: Deploy Test Resources

```bash
# From the examples directory
cd examples/

# Run deployment script
./deploy-test-resources.sh
```

The script will:
1. Check prerequisites
2. Show you what will be created
3. Ask for confirmation
4. Create all resources (~2-3 minutes)
5. Display summary

**Output example:**
```
==========================================
Azure Test Resources Deployment
==========================================

[INFO] Using subscription: My Subscription (12345678-1234-1234-1234-123456789abc)
[INFO] Creating resource group: azdoc-test-rg in eastus...
[SUCCESS] Resource group created
[INFO] Creating Hub VNet: vnet-hub...
...
==========================================
[SUCCESS] Deployment Complete!
==========================================

Duration: 142s
Resource Group: azdoc-test-rg
Location: eastus

Resources created:
  - 2 Virtual Networks (Hub and Spoke)
  - 4 Subnets
  - 2 Network Security Groups (with rules)
  - 1 Route Table (with UDRs)
  - VNet Peering
  - 2 Public IPs
  - 1 Load Balancer (Basic - Free)
  - 1 NAT Gateway
```

### Step 3: Test azdoc

```bash
# Go back to project root
cd ..

# Get your subscription ID
SUBSCRIPTION_ID=$(az account show --query id -o tsv)

# Run azdoc scan
./azdoc scan --subscription-id $SUBSCRIPTION_ID

# Build documentation
./azdoc build --with-diagrams

# View the results
cat docs/SUBSCRIPTION.md
```

### Step 4: Cleanup Resources

**⚠️ IMPORTANT: Don't forget to cleanup to avoid charges!**

```bash
# From the examples directory
cd examples/

# Run cleanup script
./cleanup-test-resources.sh
```

The script will:
1. Show all resources in the resource group
2. Ask for confirmation (type "DELETE")
3. Delete the resource group and all resources
4. Run in background (~5-10 minutes)

## Configuration

Both scripts support environment variables:

```bash
# Custom resource group name
export RESOURCE_GROUP="my-custom-rg"

# Custom location
export LOCATION="westus2"

# Custom subscription ID
export SUBSCRIPTION_ID="12345678-1234-1234-1234-123456789abc"

# Then run the script
./deploy-test-resources.sh
```

## What Gets Created

### Network Topology

```
┌─────────────────────────────────────────────┐
│ Hub VNet (10.0.0.0/16)                      │
│                                             │
│  ┌────────────────────────────────────┐    │
│  │ subnet-services (10.0.1.0/24)      │    │
│  │ - NSG: nsg-hub                     │    │
│  │ - Load Balancer (Internal)         │    │
│  └────────────────────────────────────┘    │
│                                             │
│  ┌────────────────────────────────────┐    │
│  │ subnet-management (10.0.2.0/24)    │    │
│  │ - NSG: nsg-hub                     │    │
│  └────────────────────────────────────┘    │
└──────────────┬──────────────────────────────┘
               │
               │ VNet Peering
               │
┌──────────────┴──────────────────────────────┐
│ Spoke VNet (10.1.0.0/16)                    │
│                                             │
│  ┌────────────────────────────────────┐    │
│  │ subnet-apps (10.1.1.0/24)          │    │
│  │ - NSG: nsg-spoke                   │    │
│  │ - Route Table: rt-spoke-to-hub     │    │
│  └────────────────────────────────────┘    │
│                                             │
│  ┌────────────────────────────────────┐    │
│  │ subnet-data (10.1.2.0/24)          │    │
│  │ - NAT Gateway                      │    │
│  └────────────────────────────────────┘    │
└─────────────────────────────────────────────┘
```

### Resources Details

#### Virtual Networks

| Name | Address Space | Subnets |
|------|---------------|---------|
| vnet-hub | 10.0.0.0/16 | subnet-services (10.0.1.0/24)<br>subnet-management (10.0.2.0/24) |
| vnet-spoke1 | 10.1.0.0/16 | subnet-apps (10.1.1.0/24)<br>subnet-data (10.1.2.0/24) |

#### Network Security Groups

**nsg-hub** (attached to Hub subnets):
- **Rule 100**: Allow SSH from management subnet (10.0.2.0/24)
- **Rule 200**: Allow HTTPS from Internet

**nsg-spoke** (attached to Spoke app subnet):
- **Rule 100**: Allow HTTP/HTTPS from Hub VNet
- **Rule 200**: Deny all Internet outbound

#### Route Table

**rt-spoke-to-hub** (attached to spoke subnet-apps):
- **RouteToHub**: 10.0.0.0/16 → VnetLocal
- **DefaultRouteToInternet**: 0.0.0.0/0 → Internet

#### VNet Peering

- **peer-hub-to-spoke1**: Hub → Spoke (Allow VNet access, Allow forwarded traffic)
- **peer-spoke1-to-hub**: Spoke → Hub (Allow VNet access, Allow forwarded traffic)

#### Load Balancer

**lb-internal** (Basic SKU):
- **Frontend**: Internal IP in subnet-services
- **Backend Pool**: (empty - ready for VMs)
- **Health Probe**: HTTP on port 80
- **Rule**: Port 80 TCP

#### NAT Gateway

**natgw-spoke**:
- Attached to subnet-data
- Provides outbound Internet connectivity
- 4-minute idle timeout

## Troubleshooting

### Authentication Issues

```bash
# Error: "az: command not found"
# Solution: Install Azure CLI (see Prerequisites)

# Error: "Please run 'az login' to setup account"
# Solution: Login to Azure
az login

# Error: "The subscription ... doesn't exist"
# Solution: List and set correct subscription
az account list --output table
az account set --subscription "Your Subscription Name"
```

### Permission Issues

```bash
# Error: "The client ... does not have authorization"
# Solution: You need Contributor role on the subscription

# Check your role assignments
az role assignment list --assignee $(az account show --query user.name -o tsv)

# Ask your admin to grant Contributor role
```

### Resource Already Exists

```bash
# Error: "Resource group already exists"
# Solution: The script will skip creation and continue

# If you want to start fresh, cleanup first
./cleanup-test-resources.sh
./deploy-test-resources.sh
```

### Quota Limits

```bash
# Error: "Operation could not be completed as it results in exceeding approved ... quota"
# Solution: Choose a different region or request quota increase

# Try a different region
export LOCATION="westus2"
./deploy-test-resources.sh

# Check quotas
az vm list-usage --location eastus --output table
```

## Verify Resources

### Using Azure CLI

```bash
# List all resources in the group
az resource list --resource-group azdoc-test-rg --output table

# Show VNets
az network vnet list --resource-group azdoc-test-rg --output table

# Show peering status
az network vnet peering list \
  --resource-group azdoc-test-rg \
  --vnet-name vnet-hub \
  --output table

# Show NSG rules
az network nsg rule list \
  --resource-group azdoc-test-rg \
  --nsg-name nsg-hub \
  --output table

# Show routes
az network route-table route list \
  --resource-group azdoc-test-rg \
  --route-table-name rt-spoke-to-hub \
  --output table
```

### Using Azure Portal

1. Go to https://portal.azure.com
2. Navigate to Resource Groups
3. Click on `azdoc-test-rg`
4. Explore resources and visualize topology

## Advanced Usage

### Deploy to Multiple Regions

```bash
# Region 1
export RESOURCE_GROUP="azdoc-test-eastus"
export LOCATION="eastus"
./deploy-test-resources.sh

# Region 2
export RESOURCE_GROUP="azdoc-test-westus"
export LOCATION="westus2"
./deploy-test-resources.sh
```

### Custom IP Ranges

Edit `deploy-test-resources.sh` and modify:

```bash
# Hub VNet
HUB_VNET_PREFIX="192.168.0.0/16"
HUB_SUBNET1_PREFIX="192.168.1.0/24"
HUB_SUBNET2_PREFIX="192.168.2.0/24"

# Spoke VNet
SPOKE_VNET_PREFIX="172.16.0.0/16"
SPOKE_SUBNET1_PREFIX="172.16.1.0/24"
SPOKE_SUBNET2_PREFIX="172.16.2.0/24"
```

### Add More Resources

The scripts are designed to be extended. You can add:
- More VNets and peerings
- Azure Firewall
- Application Gateway
- Virtual Machines (will increase cost significantly)
- Private Endpoints
- VPN Gateway or ExpressRoute

## Cost Management

### Monitor Costs

```bash
# View current costs
az consumption usage list --output table

# Set up budget alerts in Azure Portal
# Cost Management + Billing → Budgets → Add
```

### Minimize Costs

1. **Delete when not in use**: Run cleanup script every night
2. **Use Basic SKUs**: Already done in the script
3. **Avoid compute resources**: Script doesn't create VMs
4. **Use free tier**: Many services are free (VNets, NSGs, etc.)

### Free Trial Benefits

New Azure accounts get:
- **$200 credit** for 30 days
- **12 months** of select free services
- **Always free** services (25+ products)

Learn more: https://azure.microsoft.com/free/

## Automation

### Scheduled Deployment and Cleanup

```bash
# Create cron job to cleanup at night (6 PM)
crontab -e

# Add this line (adjust path)
0 18 * * * cd /path/to/azdoc/examples && ./cleanup-test-resources.sh <<< $'DELETE\nyes'

# Deploy in the morning (8 AM)
0 8 * * * cd /path/to/azdoc/examples && ./deploy-test-resources.sh <<< 'yes'
```

## Next Steps

After deploying test resources:

1. **Run azdoc** to generate documentation
2. **Explore the output** (Markdown + diagrams)
3. **Modify resources** and re-run azdoc to see changes
4. **Implement remaining azdoc features** using this test environment
5. **Cleanup resources** when done to avoid charges

## Support

- **Script issues**: Open issue at https://github.com/naveen/azdoc/issues
- **Azure billing**: https://azure.microsoft.com/support/
- **Azure CLI help**: `az --help` or https://docs.microsoft.com/cli/azure/

## License

These example scripts are part of the azdoc project and are provided under the MIT License.
