# Quick Reference Guide

## One-Line Commands

### Deploy Everything
```bash
cd examples && ./deploy-test-resources.sh
```

### Check Costs
```bash
cd examples && ./check-costs.sh
```

### Test azdoc
```bash
./azdoc scan --subscription-id $(az account show --query id -o tsv) && ./azdoc build
```

### Cleanup Everything
```bash
cd examples && ./cleanup-test-resources.sh
```

## Azure CLI Quick Reference

### Authentication
```bash
az login                                    # Login interactively
az account list --output table              # List subscriptions
az account set --subscription "My Sub"      # Set active subscription
az account show                             # Show current subscription
```

### Resource Groups
```bash
az group list --output table                          # List all resource groups
az group show --name azdoc-test-rg                    # Show specific group
az group delete --name azdoc-test-rg --yes --no-wait # Delete group
```

### Virtual Networks
```bash
az network vnet list -g azdoc-test-rg --output table              # List VNets
az network vnet show -g azdoc-test-rg -n vnet-hub                 # Show VNet details
az network vnet subnet list -g azdoc-test-rg --vnet-name vnet-hub # List subnets
```

### Network Security Groups
```bash
az network nsg list -g azdoc-test-rg --output table                    # List NSGs
az network nsg rule list -g azdoc-test-rg --nsg-name nsg-hub -o table  # List rules
```

### Route Tables
```bash
az network route-table list -g azdoc-test-rg --output table                          # List route tables
az network route-table route list -g azdoc-test-rg --route-table-name rt-spoke-to-hub # List routes
```

### VNet Peering
```bash
az network vnet peering list -g azdoc-test-rg --vnet-name vnet-hub --output table
```

### Load Balancers
```bash
az network lb list -g azdoc-test-rg --output table
az network lb show -g azdoc-test-rg -n lb-internal
```

### Public IPs
```bash
az network public-ip list -g azdoc-test-rg --output table
```

### NAT Gateway
```bash
az network nat gateway list -g azdoc-test-rg --output table
```

## Environment Variables

```bash
export RESOURCE_GROUP="azdoc-test-rg"              # Resource group name
export LOCATION="eastus"                           # Azure region
export SUBSCRIPTION_ID="your-sub-id"               # Subscription ID
```

## Cost Commands

```bash
# Show resource costs
az consumption usage list \
  --start-date $(date -d "7 days ago" +%Y-%m-%d) \
  --end-date $(date +%Y-%m-%d) \
  --output table

# Show budget
az consumption budget list

# Show usage details
az consumption usage list \
  --query "[?contains(instanceName, 'azdoc')]" \
  --output table
```

## Debugging Commands

```bash
# Get deployment status
az deployment group list -g azdoc-test-rg --output table

# Show resource group activities
az monitor activity-log list \
  --resource-group azdoc-test-rg \
  --max-events 50 \
  --output table

# Check service health
az rest --method get \
  --uri https://management.azure.com/subscriptions/$(az account show --query id -o tsv)/providers/Microsoft.ResourceHealth/availabilityStatuses?api-version=2015-01-01
```

## Useful Queries

### Get all resource IDs
```bash
az resource list -g azdoc-test-rg --query "[].id" -o tsv
```

### Get all resource locations
```bash
az resource list -g azdoc-test-rg --query "[].location" -o tsv | sort -u
```

### Get all resource types
```bash
az resource list -g azdoc-test-rg --query "[].type" -o tsv | sort -u
```

### Count resources by type
```bash
az resource list -g azdoc-test-rg \
  --query "[].type" -o tsv | sort | uniq -c
```

### Get all tags
```bash
az resource list -g azdoc-test-rg \
  --query "[].tags" -o json
```

## Common Issues & Solutions

### Issue: "az: command not found"
```bash
# Install Azure CLI
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
```

### Issue: "Please run 'az login'"
```bash
az login
az account show
```

### Issue: "Resource already exists"
```bash
# Clean up first
./cleanup-test-resources.sh
# Then deploy again
./deploy-test-resources.sh
```

### Issue: "Quota exceeded"
```bash
# Try different region
export LOCATION="westus2"
./deploy-test-resources.sh
```

### Issue: "Permission denied"
```bash
# Check permissions
az role assignment list --assignee $(az account show --query user.name -o tsv)
# You need "Contributor" role
```

## Pricing Reference (East US)

| Resource | Cost |
|----------|------|
| Virtual Network | FREE |
| Subnet | FREE |
| NSG | FREE |
| Route Table | FREE |
| VNet Peering | $0.01/GB transferred |
| Basic Load Balancer | FREE |
| Public IP (Basic) | $0.0036/hour (~$2.60/month) |
| Public IP (Standard) | $0.0054/hour (~$3.90/month) |
| NAT Gateway | $0.045/hour (~$32.40/month) |
| NAT Gateway Data Processing | $0.045/GB |

**Total for test setup: ~$1-2/day**

## Time Estimates

| Action | Time |
|--------|------|
| Deploy resources | 2-3 minutes |
| Cleanup resources | 5-10 minutes |
| Run azdoc scan | 30-60 seconds |
| Build documentation | 5-10 seconds |

## Next Steps

1. ✅ Deploy test resources
2. ✅ Verify in Azure Portal
3. ✅ Run azdoc scan
4. ✅ View generated docs
5. ✅ Cleanup resources

## Links

- [Azure CLI Documentation](https://docs.microsoft.com/cli/azure/)
- [Azure Pricing Calculator](https://azure.microsoft.com/pricing/calculator/)
- [Azure Free Account](https://azure.microsoft.com/free/)
- [azdoc GitHub](https://github.com/naveen/azdoc)
