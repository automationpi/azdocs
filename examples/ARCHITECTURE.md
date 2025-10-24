# Test Environment Architecture

## Overview

This document describes the Azure test environment created by `deploy-test-resources.sh`.

## Network Topology

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Azure Subscription                           │
│                                                                       │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ Resource Group: azdoc-test-rg                               │   │
│  │ Location: East US                                           │   │
│  │                                                             │   │
│  │  ┌────────────────────────────────────────────────────┐    │   │
│  │  │ Hub VNet (vnet-hub)                                │    │   │
│  │  │ Address Space: 10.0.0.0/16                        │    │   │
│  │  │                                                    │    │   │
│  │  │  ┌──────────────────────────────────────────┐    │    │   │
│  │  │  │ subnet-services (10.0.1.0/24)            │    │    │   │
│  │  │  │                                          │    │    │   │
│  │  │  │ ┌──────────────┐  ┌────────────────┐   │    │    │   │
│  │  │  │ │ NSG: nsg-hub │  │ Load Balancer  │   │    │    │   │
│  │  │  │ └──────────────┘  │  (lb-internal) │   │    │    │   │
│  │  │  │                   │   Frontend IP  │   │    │    │   │
│  │  │  │                   │   Backend Pool │   │    │    │   │
│  │  │  │                   │   Health Probe │   │    │    │   │
│  │  │  │                   └────────────────┘   │    │    │   │
│  │  │  └──────────────────────────────────────────┘    │    │   │
│  │  │                                                    │    │   │
│  │  │  ┌──────────────────────────────────────────┐    │    │   │
│  │  │  │ subnet-management (10.0.2.0/24)          │    │    │   │
│  │  │  │                                          │    │    │   │
│  │  │  │ ┌──────────────┐                        │    │    │   │
│  │  │  │ │ NSG: nsg-hub │                        │    │    │   │
│  │  │  │ └──────────────┘                        │    │    │   │
│  │  │  └──────────────────────────────────────────┘    │    │   │
│  │  │                                                    │    │   │
│  │  └────────────────────────┬───────────────────────────┘    │   │
│  │                           │                                │   │
│  │                           │ VNet Peering                   │   │
│  │                           │ (Allow VNet Access +           │   │
│  │                           │  Allow Forwarded Traffic)      │   │
│  │                           │                                │   │
│  │  ┌────────────────────────┴───────────────────────────┐   │   │
│  │  │ Spoke VNet (vnet-spoke1)                           │   │   │
│  │  │ Address Space: 10.1.0.0/16                        │   │   │
│  │  │                                                    │   │   │
│  │  │  ┌──────────────────────────────────────────┐    │   │   │
│  │  │  │ subnet-apps (10.1.1.0/24)                │    │   │   │
│  │  │  │                                          │    │   │   │
│  │  │  │ ┌────────────────┐  ┌────────────────┐ │    │   │   │
│  │  │  │ │ NSG: nsg-spoke │  │ Route Table:   │ │    │   │   │
│  │  │  │ └────────────────┘  │ rt-spoke-to-hub│ │    │   │   │
│  │  │  │                     └────────────────┘ │    │   │   │
│  │  │  └──────────────────────────────────────────┘    │   │   │
│  │  │                                                    │   │   │
│  │  │  ┌──────────────────────────────────────────┐    │   │   │
│  │  │  │ subnet-data (10.1.2.0/24)                │    │   │   │
│  │  │  │                                          │    │   │   │
│  │  │  │ ┌─────────────┐  ┌─────────────────┐   │    │   │   │
│  │  │  │ │ NAT Gateway │──│ Public IP (Std) │   │    │   │   │
│  │  │  │ │ natgw-spoke │  │                 │   │    │   │   │
│  │  │  │ └─────────────┘  └─────────────────┘   │    │   │   │
│  │  │  └──────────────────────────────────────────┘    │   │   │
│  │  │                                                    │   │   │
│  │  └────────────────────────────────────────────────────┘   │   │
│  │                                                             │   │
│  │  ┌────────────────────────────────────────────────────┐   │   │
│  │  │ Public IPs                                         │   │   │
│  │  │  - pip-lb (Basic, Static) → Load Balancer         │   │   │
│  │  │  - natgw-spoke-pip (Standard, Static) → NAT GW    │   │   │
│  │  └────────────────────────────────────────────────────┘   │   │
│  │                                                             │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                       │
└─────────────────────────────────────────────────────────────────────┘
```

## Security Rules

### NSG: nsg-hub

| Priority | Name | Direction | Action | Protocol | Source | Dest | Port |
|----------|------|-----------|--------|----------|--------|------|------|
| 100 | AllowSSHFromManagement | Inbound | Allow | TCP | 10.0.2.0/24 | * | 22 |
| 200 | AllowHTTPSFromInternet | Inbound | Allow | TCP | Internet | * | 443 |

**Applied to:**
- subnet-services (10.0.1.0/24)
- subnet-management (10.0.2.0/24)

### NSG: nsg-spoke

| Priority | Name | Direction | Action | Protocol | Source | Dest | Port |
|----------|------|-----------|--------|----------|--------|------|------|
| 100 | AllowHTTPFromHub | Inbound | Allow | TCP | 10.0.0.0/16 | * | 80,443 |
| 200 | DenyInternetOutbound | Outbound | Deny | * | * | Internet | * |

**Applied to:**
- subnet-apps (10.1.1.0/24)

## Routing

### Route Table: rt-spoke-to-hub

| Name | Address Prefix | Next Hop Type |
|------|----------------|---------------|
| RouteToHub | 10.0.0.0/16 | VnetLocal |
| DefaultRouteToInternet | 0.0.0.0/0 | Internet |

**Associated with:**
- subnet-apps (10.1.1.0/24)

### VNet Peering

#### peer-hub-to-spoke1 (Hub → Spoke)
- Remote VNet: vnet-spoke1 (10.1.0.0/16)
- Allow VNet Access: ✅
- Allow Forwarded Traffic: ✅
- Allow Gateway Transit: ❌
- Use Remote Gateways: ❌
- Peering State: Connected

#### peer-spoke1-to-hub (Spoke → Hub)
- Remote VNet: vnet-hub (10.0.0.0/16)
- Allow VNet Access: ✅
- Allow Forwarded Traffic: ✅
- Allow Gateway Transit: ❌
- Use Remote Gateways: ❌
- Peering State: Connected

## Load Balancer Configuration

### lb-internal (Basic SKU)

**Frontend Configuration:**
- Name: frontend-internal
- Type: Internal
- Subnet: subnet-services (10.0.1.0/24)
- Private IP: Dynamic

**Backend Pool:**
- Name: backend-pool
- Members: (empty - ready for VMs)

**Health Probe:**
- Name: health-probe-http
- Protocol: HTTP
- Port: 80
- Path: /

**Load Balancing Rule:**
- Name: rule-http
- Protocol: TCP
- Frontend Port: 80
- Backend Port: 80
- Backend Pool: backend-pool
- Health Probe: health-probe-http

## NAT Gateway Configuration

### natgw-spoke

**Configuration:**
- SKU: Standard
- Idle Timeout: 4 minutes
- Public IP: natgw-spoke-pip (Standard, Static)

**Associated with:**
- subnet-data (10.1.2.0/24)

**Purpose:**
- Provides outbound Internet connectivity for resources in subnet-data
- All outbound traffic uses the NAT Gateway's public IP
- Inbound connections to private IPs are not allowed

## Resource Tags

All resources are tagged with:
- `Purpose: azdoc-testing`
- `ManagedBy: azdoc-example-script`

Additional tags by resource type:
- **vnet-hub**: `Type: Hub`, `Environment: Test`
- **vnet-spoke1**: `Type: Spoke`, `Environment: Test`
- **nsg-hub**: `AppliesTo: HubVNet`
- **nsg-spoke**: `AppliesTo: SpokeVNet`
- **rt-spoke-to-hub**: `Purpose: SpokeRouting`
- **pip-lb**: `Purpose: LoadBalancer`

## Traffic Flow Examples

### Example 1: Internet to Hub Services

```
Internet → Public IP (pip-lb)
        → Load Balancer (lb-internal)
        → Backend Pool (subnet-services)
        ← NSG: nsg-hub (Allow HTTPS from Internet)
```

### Example 2: Management to Services

```
subnet-management (10.0.2.0/24)
        → subnet-services (10.0.1.0/24)
        ← NSG: nsg-hub (Allow SSH from Management)
```

### Example 3: Hub to Spoke

```
vnet-hub (10.0.0.0/16)
        → VNet Peering
        → vnet-spoke1 (10.1.0.0/16)
        → subnet-apps (10.1.1.0/24)
        ← NSG: nsg-spoke (Allow HTTP/HTTPS from Hub)
```

### Example 4: Spoke Apps Outbound

```
subnet-apps (10.1.1.0/24)
        → Route Table: rt-spoke-to-hub (Default Route)
        → Internet
        ← NSG: nsg-spoke (BLOCKED - Deny Internet Outbound)
```

### Example 5: Spoke Data Outbound (via NAT)

```
subnet-data (10.1.2.0/24)
        → NAT Gateway (natgw-spoke)
        → Public IP (natgw-spoke-pip)
        → Internet
```

## Address Space Allocation

| Network | Address Space | Available IPs |
|---------|---------------|---------------|
| **Hub VNet** | 10.0.0.0/16 | 65,536 |
| ├─ subnet-services | 10.0.1.0/24 | 251 (5 reserved by Azure) |
| └─ subnet-management | 10.0.2.0/24 | 251 |
| **Spoke VNet** | 10.1.0.0/16 | 65,536 |
| ├─ subnet-apps | 10.1.1.0/24 | 251 |
| └─ subnet-data | 10.1.2.0/24 | 251 |

**Total IP Space:** 131,072 IPs across both VNets

## Cost Breakdown (Est.)

| Resource | Quantity | Daily Cost | Notes |
|----------|----------|------------|-------|
| Virtual Networks | 2 | $0 | FREE |
| Subnets | 4 | $0 | FREE |
| NSGs | 2 | $0 | FREE |
| NSG Rules | 4 | $0 | FREE |
| Route Tables | 1 | $0 | FREE |
| Routes (UDRs) | 2 | $0 | FREE |
| VNet Peering | 1 pair | ~$0 | $0.01/GB (minimal without VMs) |
| Basic Load Balancer | 1 | $0 | FREE |
| Public IP (Basic) | 1 | $0.09 | $2.60/month |
| NAT Gateway | 1 | $1.08 | $32.40/month fixed |
| Public IP (Standard) | 1 | $0.15 | $4.50/month |
| Data Processing | varies | $0+ | $0.045/GB via NAT GW |
| **TOTAL** | | **~$1.32/day** | **~$40/month** |

## Resource Limits Used

| Limit Type | Used | Default Limit | % Used |
|------------|------|---------------|--------|
| VNets per subscription | 2 | 1,000 | 0.2% |
| Subnets per VNet | 2 each | 3,000 | <0.1% |
| NSGs per subscription | 2 | 5,000 | <0.1% |
| NSG rules per NSG | 2 each | 1,000 | 0.2% |
| Route tables per subscription | 1 | 200 | 0.5% |
| Public IPs per subscription | 2 | 1,000 | 0.2% |
| Load Balancers per subscription | 1 | 1,000 | 0.1% |

## Cleanup Order

When running `cleanup-test-resources.sh`, resources are deleted in this order:

1. **NAT Gateway** (releases association with subnet)
2. **Load Balancer** (releases association with subnet)
3. **VNet Peering** (both directions)
4. **Route Table** (releases association with subnet)
5. **Network Security Groups** (releases association with subnets)
6. **Public IPs** (released from NAT GW and LB)
7. **Subnets** (part of VNet deletion)
8. **Virtual Networks**
9. **Resource Group** (final cleanup)

> Note: Azure handles dependencies automatically when deleting the resource group.

## Monitoring & Diagnostics

Enable these for production use (not included in test setup):

- **Network Watcher**: Connection monitor, packet capture
- **NSG Flow Logs**: Traffic analytics
- **Diagnostic Settings**: Send logs to Log Analytics
- **Azure Monitor**: Metrics and alerts
- **Traffic Analytics**: Visualize network traffic patterns

## Next Steps

1. Deploy: `./deploy-test-resources.sh`
2. Verify in Azure Portal
3. Run azdoc: `../azdoc scan --subscription-id <id>`
4. Review generated documentation
5. Cleanup: `./cleanup-test-resources.sh`

## References

- [Azure Virtual Network](https://docs.microsoft.com/azure/virtual-network/)
- [Network Security Groups](https://docs.microsoft.com/azure/virtual-network/network-security-groups-overview)
- [Azure Load Balancer](https://docs.microsoft.com/azure/load-balancer/)
- [Azure NAT Gateway](https://docs.microsoft.com/azure/virtual-network/nat-gateway/)
- [VNet Peering](https://docs.microsoft.com/azure/virtual-network/virtual-network-peering-overview)
