# Examples Directory Index

## Files

### Scripts (Executable)

1. **deploy-test-resources.sh** (16 KB)
   - Creates complete Azure test environment
   - Includes: 2 VNets, 4 Subnets, NSGs, Route Tables, Peering, LB, NAT Gateway
   - Cost: ~$1-2/day
   - Runtime: 2-3 minutes

2. **cleanup-test-resources.sh** (5 KB)
   - Safely deletes all test resources
   - Requires explicit confirmation ("DELETE")
   - Runtime: 5-10 minutes (background)

3. **check-costs.sh** (4 KB)
   - Shows current resource deployment
   - Estimates daily costs
   - Lists all resources in resource group

### Documentation

1. **README.md** (13 KB)
   - Complete guide for test environment
   - Prerequisites, installation, usage
   - Cost estimates and troubleshooting

2. **ARCHITECTURE.md** (10 KB)
   - Visual network topology diagram
   - Security rules documentation
   - Routing configuration
   - Traffic flow examples
   - Resource limits and cost breakdown

3. **QUICK_REFERENCE.md** (5 KB)
   - One-line commands for common tasks
   - Azure CLI quick reference
   - Pricing table
   - Time estimates

4. **INDEX.md** (this file)
   - Overview of all files in examples directory

## Quick Commands

```bash
# Deploy everything
./deploy-test-resources.sh

# Check what's running
./check-costs.sh

# Cleanup everything
./cleanup-test-resources.sh
```

## Files by Purpose

### For New Users
1. Start with: **README.md**
2. Run: **deploy-test-resources.sh**
3. Then: Run azdoc from parent directory

### For Quick Reference
1. **QUICK_REFERENCE.md** - All commands in one place
2. **check-costs.sh** - See what's deployed

### For Understanding Architecture
1. **ARCHITECTURE.md** - Complete topology and config
2. View in Azure Portal for visual confirmation

### For Cleanup
1. **cleanup-test-resources.sh** - Delete everything safely

## Total Size

- Scripts: ~25 KB
- Documentation: ~28 KB
- **Total: ~53 KB**

## What You'll Create

- **2 Virtual Networks** (Hub-and-Spoke)
- **4 Subnets** (2 per VNet)
- **2 Network Security Groups** (4 rules total)
- **1 Route Table** (2 custom routes)
- **VNet Peering** (bidirectional)
- **1 Load Balancer** (Basic - Free tier)
- **1 NAT Gateway**
- **2 Public IPs**

## Cost Summary

| Component | Cost/Day |
|-----------|----------|
| Network resources | FREE |
| Basic Load Balancer | FREE |
| Public IPs | $0.24 |
| NAT Gateway | $1.08+ |
| **Total** | **~$1.32+/day** |

## Time Estimates

| Task | Duration |
|------|----------|
| Read README | 10 min |
| Deploy resources | 3 min |
| Run azdoc | 1 min |
| Review docs | 10 min |
| Cleanup | 10 min |
| **Total** | **~34 min** |

## Learn More

- Parent README: `../README.md`
- Implementation guide: `../IMPLEMENTATION.md`
- Project summary: `../PROJECT_SUMMARY.md`

---

**Ready to start?** Run `./deploy-test-resources.sh`
