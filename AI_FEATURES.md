# AI-Enhanced Diagram Generation

azdoc now supports AI-powered diagram generation using OpenAI's GPT-4 for intelligent layout optimization and connection discovery.

## Features

### 1. **Intelligent Connection Discovery**
The LLM analyzes Azure resource properties to automatically identify:
- VM-to-NIC attachments
- Function App storage dependencies
- VNet peering relationships
- NSG subnet associations
- Public IP attachments
- NAT Gateway associations

### 2. **Smart Layout Optimization**
The LLM suggests optimal positioning for resources based on:
- Subnet boundaries
- Resource relationships
- Visual hierarchy (important resources prominently placed)
- Avoiding overlaps
- Creating logical groupings

## Setup

### Prerequisites
- OpenAI API key (get one at https://platform.openai.com/api-keys)
- Set environment variable:
  ```bash
  export OPENAI_API_KEY="sk-..."
  ```

### Usage

**Enable AI features with flag:**
```bash
./azdoc build --with-diagrams --enable-ai
```

**Or provide API key directly:**
```bash
./azdoc build --with-diagrams --enable-ai --openai-key="sk-..."
```

**Run complete pipeline with AI:**
```bash
./azdoc all --subscription-id <YOUR_SUB_ID> --with-diagrams --enable-ai
```

## How It Works

### Connection Discovery Flow:
1. azdoc scans Azure resources and caches properties
2. LLM receives resource JSON with all Azure properties
3. LLM analyzes:
   - NIC attachments (by parsing properties.virtualMachine.id)
   - Storage account references (Function App → Storage)
   - VNet integrations (subnet associations)
   - Naming conventions (vm-hub-jumpbox-nic → vm-hub-jumpbox)
4. Returns structured JSON with connections and confidence levels
5. **azdoc automatically renders high-confidence connections in diagrams**
6. Falls back to rule-based connections if AI is disabled or unavailable

### Layout Optimization Flow:
1. LLM receives resources + subnet information
2. Analyzes resource types and relationships
3. Calculates optimal x,y coordinates
4. Groups related resources (compute in subnets, network resources outside)
5. Returns positioning suggestions with reasoning
6. **[Coming Soon]** azdoc applies layout while respecting diagram constraints

### Current Implementation Status:
✅ **Connection Discovery** - Fully implemented and active
- AI analyzes all Azure resource properties
- Returns connection suggestions with confidence scores
- High-confidence connections are automatically rendered
- Visual feedback shows which connections were AI-discovered

🚧 **Layout Optimization** - Prepared but not yet integrated
- LLM client method implemented (`OptimizeLayout`)
- Will be integrated in next iteration
- Current layout uses rule-based positioning

## Cost Estimate

**Per Diagram Generation:**
- Connection Discovery: ~500-1000 tokens (~$0.01-0.02 with GPT-4 Turbo)
- Layout Optimization: ~1000-2000 tokens (~$0.02-0.04 with GPT-4 Turbo)
- **Total per diagram: ~$0.03-0.06**

**For typical environment (2 VNets + Overview):**
- 3 diagrams × $0.05 = **~$0.15 per documentation build**

## Benefits

### Without AI:
- Manual rule-based positioning
- Simple naming pattern matching for connections
- Fixed layouts
- May miss complex relationships
- Limited to hardcoded connection patterns

### With AI (Current):
- ✅ Intelligent understanding of Azure resource properties
- ✅ Discovers all relationships automatically by analyzing Azure metadata
- ✅ Adapts to different resource types and configurations
- ✅ Natural language explanations for each connection
- ✅ Confidence scoring helps filter noise
- ✅ Console output shows exactly which AI connections were applied
- ✅ Graceful fallback if API unavailable

## Example Output

### Console Output When Using AI:
```
Generating Draw.io diagrams...
🤖 AI-enhanced diagram generation enabled
  🤖 Analyzing resources with AI for intelligent connections...
  ✅ AI discovered 8 connections
    🤖 Applying 8 AI-discovered connections to vnet-hub-eastus...
      ✅ vm-hub-jumpbox-nic → vm-hub-jumpbox (attached to)
      ✅ func-spoke-app-18709 → stfunc13472 (uses)
      ✅ subnet-management → nsg-management (protected by)
    ✨ Applied 3 high-confidence AI connections
```

**AI-Discovered Connections JSON:**
```json
[
  {
    "source_resource": "vm-hub-jumpbox-nic",
    "target_resource": "vm-hub-jumpbox",
    "connection_type": "association",
    "label": "attached to",
    "confidence": "high",
    "reason": "NIC properties.virtualMachine.id references this VM"
  },
  {
    "source_resource": "func-spoke-app-18709",
    "target_resource": "stfunc13472",
    "connection_type": "association",
    "label": "uses",
    "confidence": "high",
    "reason": "Function App requires storage account for runtime"
  }
]
```

**AI-Suggested Layout:**
```json
[
  {
    "resource_name": "vm-hub-jumpbox",
    "x": 120.0,
    "y": 160.0,
    "subnet_name": "subnet-management",
    "grouping": "subnet",
    "reason": "VM should be inside management subnet as it's a jumpbox for administrative access"
  }
]
```

## Privacy & Security

- Resource data is sent to OpenAI API (review their data usage policy)
- Only resource metadata is sent (no secrets/keys)
- API calls are over HTTPS
- For sensitive environments, consider using local LLM (future feature)

## Troubleshooting

**API Key Error:**
```
Warning: OpenAI API key not found. AI features disabled.
```
→ Set OPENAI_API_KEY environment variable

**Rate Limit Error:**
```
OpenAI API error: rate limit exceeded
```
→ Wait a few moments and retry, or upgrade your OpenAI plan

**Invalid JSON Response:**
```
failed to parse LLM response
```
→ Rare issue with LLM output formatting, retry the command

## Future Enhancements

- [ ] Support for local LLMs (Ollama integration)
- [ ] Anthropic Claude API support
- [ ] Caching LLM responses for faster rebuilds
- [ ] Security analysis powered by LLM
- [ ] Natural language documentation generation
