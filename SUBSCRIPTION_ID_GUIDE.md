# Azure Subscription ID Guide

## What is a Subscription ID?

An Azure Subscription ID is a unique GUID (Globally Unique Identifier) that identifies your Azure subscription. It looks like this:

```
12345678-1234-1234-1234-123456789abc
```

## Why Do You Need It?

azdoc needs to know which Azure subscription to scan and document. You must provide this ID when running commands like:
- `azdoc scan`
- `azdoc all`
- `azdoc doctor`

## How to Find Your Subscription ID

### Method 1: Quick Helper Script (Easiest)

```bash
cd examples
./get-subscription-id.sh
```

This will:
- Show your current subscription
- List all available subscriptions
- Provide copy-paste ready commands
- Give you quick one-liners

### Method 2: Azure CLI

```bash
# Show current subscription ID
az account show --query id -o tsv

# Show current subscription details
az account show

# List all your subscriptions
az account list --output table

# List with just names and IDs
az account list --query "[].{Name:name, ID:id}" --output table
```

### Method 3: Azure Portal

1. Go to https://portal.azure.com
2. Search for "Subscriptions" in the top search bar
3. Click on "Subscriptions"
4. Your subscription ID is shown in the "Subscription ID" column

### Method 4: Check Environment Variable

```bash
# If you've already set it
echo $AZURE_SUBSCRIPTION_ID

# Or
echo $SUBSCRIPTION_ID
```

## How to Use It

### Option 1: Command-Line Flag (Recommended)

```bash
# Direct usage
./azdoc scan --subscription-id "12345678-1234-1234-1234-123456789abc"

# Using command substitution (auto-detects current subscription)
./azdoc scan --subscription-id $(az account show --query id -o tsv)
```

### Option 2: Environment Variable

```bash
# Set the environment variable
export AZURE_SUBSCRIPTION_ID="12345678-1234-1234-1234-123456789abc"

# Or use the SUBSCRIPTION_ID variable (for examples)
export SUBSCRIPTION_ID="12345678-1234-1234-1234-123456789abc"

# Then run without the flag
./azdoc scan
```

### Option 3: Configuration File

Create `azdoc.yaml`:

```yaml
subscription-id: "12345678-1234-1234-1234-123456789abc"
```

Then run:

```bash
./azdoc scan --config azdoc.yaml
```

## Multiple Subscriptions

If you have multiple Azure subscriptions:

### List All Subscriptions

```bash
az account list --output table
```

### Switch Active Subscription

```bash
# By name
az account set --subscription "My Production Subscription"

# By ID
az account set --subscription "12345678-1234-1234-1234-123456789abc"

# Verify
az account show
```

### Document Multiple Subscriptions

```bash
# Subscription 1
./azdoc all --subscription-id "12345678-1234-1234-1234-111111111111" \
  --json-out ./data/sub1 \
  --out ./docs/sub1

# Subscription 2
./azdoc all --subscription-id "12345678-1234-1234-1234-222222222222" \
  --json-out ./data/sub2 \
  --out ./docs/sub2
```

## Quick Reference Commands

### Get Subscription ID

```bash
# Current subscription ID only
az account show --query id -o tsv

# Current subscription with details
az account show --query "{Name:name, ID:id, State:state}" -o table

# All subscriptions
az account list --query "[].{Name:name, ID:id, State:state}" -o table
```

### Set Subscription

```bash
# By name
az account set --subscription "My Subscription"

# By ID
az account set --subscription "12345678-1234-1234-1234-123456789abc"

# By environment variable
az account set --subscription "$SUBSCRIPTION_ID"
```

### Verify Current Subscription

```bash
# Full details
az account show

# Just the ID
az account show --query id -o tsv

# Just the name
az account show --query name -o tsv
```

## Common Issues

### Issue: "No subscriptions found"

**Problem:** You're logged in but have no subscriptions.

**Solutions:**
1. Sign up for Azure Free Account: https://azure.microsoft.com/free/
2. Ask your organization admin to add you to a subscription
3. Check if you're in the right Azure AD tenant: `az account list --all`

### Issue: "Subscription ... doesn't exist"

**Problem:** The subscription ID is wrong or you don't have access.

**Solutions:**
1. Double-check the subscription ID (copy-paste carefully)
2. Verify you have access: `az account list`
3. Ask admin for Reader or Contributor role
4. Make sure you're in the correct tenant

### Issue: "Please run 'az login'"

**Problem:** Not authenticated to Azure.

**Solution:**
```bash
az login
az account show
```

### Issue: Wrong subscription being used

**Problem:** Azure CLI is using a different subscription than you expected.

**Solution:**
```bash
# Check current subscription
az account show --query name -o tsv

# Set the correct one
az account set --subscription "Correct Subscription Name"

# Or use explicit --subscription-id flag
./azdoc scan --subscription-id "12345678-1234-1234-1234-123456789abc"
```

## Security Best Practices

### Don't Hardcode Subscription IDs

❌ **Bad:**
```bash
# Hardcoded in script
./azdoc scan --subscription-id "12345678-1234-1234-1234-123456789abc"
```

✅ **Good:**
```bash
# Use environment variable
export SUBSCRIPTION_ID=$(az account show --query id -o tsv)
./azdoc scan --subscription-id "$SUBSCRIPTION_ID"

# Or use command substitution
./azdoc scan --subscription-id $(az account show --query id -o tsv)
```

### Don't Commit Subscription IDs to Git

Add to `.gitignore`:
```
azdoc.yaml
.env
*.local
```

Use example files instead:
```bash
# Checked into git
azdoc.yaml.example

# Local file (gitignored)
azdoc.yaml
```

## Examples

### Example 1: First Time Setup

```bash
# 1. Login
az login

# 2. Get subscription ID
SUBSCRIPTION_ID=$(az account show --query id -o tsv)
echo "Using subscription: $SUBSCRIPTION_ID"

# 3. Verify access
./azdoc doctor --subscription-id "$SUBSCRIPTION_ID"

# 4. Run scan
./azdoc scan --subscription-id "$SUBSCRIPTION_ID"
```

### Example 2: Using Helper Script

```bash
# Run the helper
cd examples
./get-subscription-id.sh

# It will show you copy-paste ready commands
# Just copy and run them!
```

### Example 3: Automated Script

```bash
#!/bin/bash

# Get current subscription automatically
SUBSCRIPTION_ID=$(az account show --query id -o tsv)

if [ -z "$SUBSCRIPTION_ID" ]; then
    echo "Error: Not logged in to Azure"
    echo "Run: az login"
    exit 1
fi

echo "Documenting subscription: $SUBSCRIPTION_ID"

# Run azdoc
./azdoc all --subscription-id "$SUBSCRIPTION_ID"
```

### Example 4: CI/CD Pipeline

```bash
#!/bin/bash
# For GitHub Actions or Azure DevOps

# Subscription ID from secret/variable
SUBSCRIPTION_ID="${AZURE_SUBSCRIPTION_ID}"

# Or from service principal login
SUBSCRIPTION_ID=$(az account show --query id -o tsv)

# Run azdoc
./azdoc all \
  --subscription-id "$SUBSCRIPTION_ID" \
  --no-progress \
  --quiet
```

## Quick Start (Copy-Paste)

```bash
# 1. Login to Azure
az login

# 2. Get your subscription ID
az account show --query id -o tsv

# 3. Run azdoc with your subscription ID
./azdoc scan --subscription-id $(az account show --query id -o tsv)

# 4. Build documentation
./azdoc build

# Or do it all at once
./azdoc all --subscription-id $(az account show --query id -o tsv)
```

## Need Help?

Run the helper script:
```bash
cd examples
./get-subscription-id.sh
```

Or check:
- [Azure CLI Documentation](https://docs.microsoft.com/cli/azure/)
- [Azure Subscriptions Overview](https://docs.microsoft.com/azure/cost-management-billing/manage/create-subscription)
- [azdoc README](../README.md)

## Summary

**Simplest way to use azdoc:**

```bash
./azdoc all --subscription-id $(az account show --query id -o tsv)
```

This automatically uses your current Azure CLI subscription!
