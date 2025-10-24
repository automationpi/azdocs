# Examples Directory - Complete Summary

## Overview

The `examples/` directory provides a complete, ready-to-use Azure test environment for testing azdoc without touching production resources.

## What's Included

### 🔧 Automated Scripts

1. **deploy-test-resources.sh**
   - One-command deployment of complete test environment
   - Creates realistic Hub-and-Spoke network topology
   - Includes NSGs, Route Tables, VNet Peering, Load Balancer, NAT Gateway
   - Color-coded output with progress indicators
   - Built-in safety checks and confirmations
   - ~2-3 minute execution time

2. **cleanup-test-resources.sh**
   - Safe deletion with double confirmation
   - Lists all resources before deleting
   - Async deletion (runs in background)
   - Clear progress indicators
   - ~5-10 minute background execution

3. **check-costs.sh**
   - Real-time cost estimation
   - Resource inventory
   - Daily cost breakdown
   - Links to Azure cost management

### 📚 Comprehensive Documentation

1. **README.md** (13 KB)
   - Full user guide
   - Prerequisites and installation
   - Step-by-step walkthrough
   - Troubleshooting guide
   - Cost management tips
   - Azure CLI reference

2. **ARCHITECTURE.md** (10 KB)
   - ASCII network topology diagram
   - Complete security rules documentation
   - Routing table details
   - Traffic flow examples
   - Cost breakdown by resource
   - Resource limits analysis

3. **QUICK_REFERENCE.md** (5 KB)
   - One-liner commands
   - Azure CLI cheat sheet
   - Pricing table
   - Common issues & solutions
   - Time estimates

4. **INDEX.md**
   - Directory overview
   - File descriptions
   - Quick navigation

## Network Topology Created

```
Hub VNet (10.0.0.0/16)
  ├─ subnet-services (10.0.1.0/24)
  │  ├─ NSG: nsg-hub
  │  └─ Load Balancer (internal)
  └─ subnet-management (10.0.2.0/24)
     └─ NSG: nsg-hub
           │
           │ VNet Peering
           ▼
Spoke VNet (10.1.0.0/16)
  ├─ subnet-apps (10.1.1.0/24)
  │  ├─ NSG: nsg-spoke
  │  └─ Route Table: rt-spoke-to-hub
  └─ subnet-data (10.1.2.0/24)
     └─ NAT Gateway: natgw-spoke
```

## Resource Summary

| Resource Type | Count | Cost |
|---------------|-------|------|
| Virtual Networks | 2 | FREE |
| Subnets | 4 | FREE |
| Network Security Groups | 2 | FREE |
| Security Rules | 4 | FREE |
| Route Tables | 1 | FREE |
| Custom Routes (UDRs) | 2 | FREE |
| VNet Peerings | 2 | $0.01/GB |
| Load Balancers | 1 | FREE (Basic) |
| Public IPs | 2 | $0.24/day |
| NAT Gateways | 1 | $1.08/day |
| **TOTAL** | **21** | **~$1.32/day** |

## Features & Highlights

### ✅ Budget-Friendly
- **~$1-2 per day** cost estimate
- Uses free tier services where possible
- No compute resources (VMs, containers)
- Perfect for testing without breaking the bank
- Easy cleanup to stop charges

### ✅ Realistic Topology
- Industry-standard Hub-and-Spoke design
- Proper network segmentation
- Security groups with real-world rules
- Custom routing configuration
- VNet peering with traffic forwarding

### ✅ Production-Ready Scripts
- Idempotent operations
- Error handling and validation
- Progress indicators
- Safety confirmations
- Color-coded output
- Comprehensive logging

### ✅ Excellent Documentation
- Multiple documentation levels
- Visual diagrams
- Traffic flow examples
- Quick reference cards
- Troubleshooting guides

### ✅ Educational Value
- Learn Azure networking concepts
- Understand Hub-and-Spoke topology
- Practice with NSGs and routing
- Explore VNet peering
- Test azdoc functionality

## Usage Flow

```
┌─────────────────────┐
│ 1. Read README.md   │
│    (5-10 minutes)   │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ 2. az login         │
│    Set subscription │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────────────┐
│ 3. ./deploy-test-resources  │
│    (2-3 minutes)           │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────┐
│ 4. Verify in Azure Portal   │
│    (optional)               │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────┐
│ 5. Run azdoc scan & build   │
│    (1 minute)               │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────┐
│ 6. Review generated docs    │
│    (10 minutes)             │
└──────────┬──────────────────┘
           │
           ▼
┌─────────────────────────────┐
│ 7. ./cleanup-test-resources │
│    (5-10 minutes)           │
└─────────────────────────────┘
```

## Quick Start (TL;DR)

```bash
# 1. Navigate to examples
cd examples

# 2. Deploy (costs ~$1-2/day)
./deploy-test-resources.sh

# 3. Test azdoc
cd ..
./azdoc scan --subscription-id $(az account show --query id -o tsv)
./azdoc build

# 4. Cleanup when done
cd examples
./cleanup-test-resources.sh
```

## Security Considerations

### NSG Rules Included

**Hub NSG:**
- ✅ Allow SSH from management subnet
- ✅ Allow HTTPS from Internet
- ✅ Azure default rules

**Spoke NSG:**
- ✅ Allow HTTP/HTTPS from Hub
- ❌ Deny all Internet outbound
- ✅ Azure default rules

### Network Isolation

- Management subnet isolated
- Spoke apps can only reach Hub
- Spoke data uses NAT Gateway for outbound
- No public IPs on subnets (except via LB/NAT)

### Best Practices Demonstrated

- Subnet-level NSG associations
- Custom route tables
- NAT Gateway for outbound traffic
- Hub-and-Spoke for network isolation
- Principle of least privilege

## Educational Topics Covered

1. **Virtual Networks**: VNet creation, address spaces
2. **Subnets**: Segmentation, address planning
3. **NSGs**: Security rules, priorities, tags
4. **Routing**: UDRs, route tables, next-hop types
5. **VNet Peering**: Connectivity, traffic forwarding
6. **Load Balancing**: Frontend/backend, health probes
7. **NAT Gateway**: Outbound Internet, public IPs
8. **Azure CLI**: Resource management, queries
9. **Cost Management**: Estimation, optimization
10. **IaC Concepts**: Automation, idempotency

## Advanced Use Cases

### Multi-Region Testing
```bash
export RESOURCE_GROUP="azdoc-test-eastus"
export LOCATION="eastus"
./deploy-test-resources.sh

export RESOURCE_GROUP="azdoc-test-westus"
export LOCATION="westus2"
./deploy-test-resources.sh
```

### Custom IP Ranges
Edit script variables:
- `HUB_VNET_PREFIX`
- `SPOKE_VNET_PREFIX`
- Subnet prefixes

### Extended Topology
Add to scripts:
- Additional spoke VNets
- Azure Firewall
- Application Gateway
- VPN Gateway
- Private Endpoints

## Troubleshooting

### Common Issues

1. **"az: command not found"**
   - Install Azure CLI
   - See: README.md Prerequisites

2. **"Please run 'az login'"**
   - Run `az login`
   - Set subscription

3. **"Resource already exists"**
   - Cleanup first: `./cleanup-test-resources.sh`
   - Or use different RESOURCE_GROUP

4. **"Quota exceeded"**
   - Try different region
   - Request quota increase

5. **High costs**
   - Cleanup immediately
   - Check for forgotten resources

## Cost Optimization Tips

1. **Delete when not in use** - Don't leave running overnight
2. **Use Basic SKUs** - Already done in scripts
3. **Avoid compute** - No VMs in test environment
4. **Monitor daily** - Run `check-costs.sh` regularly
5. **Set budget alerts** - Azure Portal → Cost Management

## Files Checklist

- [x] deploy-test-resources.sh (16 KB)
- [x] cleanup-test-resources.sh (5 KB)
- [x] check-costs.sh (4 KB)
- [x] README.md (13 KB)
- [x] ARCHITECTURE.md (10 KB)
- [x] QUICK_REFERENCE.md (5 KB)
- [x] INDEX.md (2 KB)

**Total: 7 files, ~55 KB**

## Integration with azdoc

### Purpose
These examples provide realistic Azure resources for:
- Testing azdoc scanning
- Verifying discovery logic
- Testing graph building
- Validating Markdown rendering
- Testing diagram generation

### Workflow
1. Deploy test resources
2. Run azdoc scan
3. Verify all resources discovered
4. Check documentation accuracy
5. Test diagram generation
6. Iterate on azdoc features
7. Cleanup resources

## Next Steps

1. ✅ Deploy test environment
2. ✅ Run azdoc scan
3. ✅ Review generated documentation
4. 🔄 Implement remaining azdoc features
5. ✅ Cleanup test resources

## Support & Feedback

- Issues: https://github.com/naveen/azdoc/issues
- Discussions: https://github.com/naveen/azdoc/discussions
- Azure Support: https://azure.microsoft.com/support/

## License

These example scripts are part of the azdoc project and are provided under the MIT License.

---

**🚀 Ready to get started?** Head to `examples/README.md` for detailed instructions!
