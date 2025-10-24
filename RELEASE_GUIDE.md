# Release Guide - Creating GitHub Releases for azdoc

## 📦 How the Download Command Works

The command from the README:
```bash
curl -LO https://github.com/automationpi/azdocs/releases/latest/download/azdoc-linux-amd64
```

This command works by:
1. GitHub hosts binary artifacts attached to releases
2. `/releases/latest/download/` is a special GitHub URL that automatically redirects to the latest release
3. The binary name (`azdoc-linux-amd64`) must match the filename uploaded to the release

---

## 🚀 How to Create Your First Release

### Option 1: Automated Release (Recommended) ✅

We've set up GitHub Actions to automatically build and release binaries when you push a version tag.

#### Step 1: Push Code to GitHub

```bash
# Make sure all changes are committed
git status

# Push to GitHub
git push -u origin master
```

#### Step 2: Create and Push a Version Tag

```bash
# Create a version tag
git tag -a v1.0.0 -m "Release v1.0.0 - AI-Powered Azure Analysis"

# Push the tag to GitHub (this triggers the release workflow)
git push origin v1.0.0
```

#### Step 3: Wait for GitHub Actions

- Go to: https://github.com/automationpi/azdocs/actions
- Watch the "Release" workflow run
- It will automatically:
  - ✅ Build binaries for all platforms
  - ✅ Generate checksums
  - ✅ Create GitHub release
  - ✅ Upload all artifacts

#### Step 4: Verify the Release

Go to: https://github.com/automationpi/azdocs/releases

You should see release `v1.0.0` with the following downloads:
- `azdoc-linux-amd64`
- `azdoc-linux-arm64`
- `azdoc-darwin-amd64`
- `azdoc-darwin-arm64`
- `azdoc-windows-amd64.exe`
- `checksums.txt`

Now the curl command will work! ✅

---

### Option 2: Manual Release (Alternative)

If you prefer to create releases manually:

#### Step 1: Build Binaries Locally

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o azdoc-linux-amd64 ./cmd/azdoc

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o azdoc-linux-arm64 ./cmd/azdoc

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o azdoc-darwin-amd64 ./cmd/azdoc

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o azdoc-darwin-arm64 ./cmd/azdoc

# Windows
GOOS=windows GOARCH=amd64 go build -o azdoc-windows-amd64.exe ./cmd/azdoc

# Generate checksums
sha256sum azdoc-* > checksums.txt
```

#### Step 2: Create Tag

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

#### Step 3: Create Release on GitHub

1. Go to: https://github.com/automationpi/azdocs/releases/new
2. Select tag: `v1.0.0`
3. Release title: `azdoc v1.0.0 - AI-Powered Azure Analysis`
4. Description: Copy from release template below
5. Attach binaries:
   - Drag and drop all `azdoc-*` files
   - Drag and drop `checksums.txt`
6. Click "Publish release"

---

## 📝 Release Description Template

```markdown
## azdoc v1.0.0

### 🚀 AI-Powered Azure Documentation & Analysis

Harness the power of GenAI to analyze your Azure infrastructure, identify security risks, optimize costs, and generate actionable insights—all automatically.

### ✨ Features

- 🤖 **AI Security Insights** - Strategic security analysis with prioritized remediation
- 💰 **Cost Optimization** - Identify savings opportunities with effort assessment
- 🏷️ **Tagging Compliance** - Track governance policy adherence
- 🔄 **DR & Monitoring** - Assess backup coverage and monitoring gaps
- 📊 **Executive Dashboard** - At-a-glance infrastructure health scores
- 🎨 **Smart Diagrams** - AI-generated network topology with intelligent connections

### 📥 Installation

**Linux (AMD64):**
```bash
curl -LO https://github.com/automationpi/azdocs/releases/download/v1.0.0/azdoc-linux-amd64
chmod +x azdoc-linux-amd64
sudo mv azdoc-linux-amd64 /usr/local/bin/azdoc
```

**macOS (Intel):**
```bash
curl -LO https://github.com/automationpi/azdocs/releases/download/v1.0.0/azdoc-darwin-amd64
chmod +x azdoc-darwin-amd64
sudo mv azdoc-darwin-amd64 /usr/local/bin/azdoc
```

**macOS (Apple Silicon M1/M2):**
```bash
curl -LO https://github.com/automationpi/azdocs/releases/download/v1.0.0/azdoc-darwin-arm64
chmod +x azdoc-darwin-arm64
sudo mv azdoc-darwin-arm64 /usr/local/bin/azdoc
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://github.com/automationpi/azdocs/releases/download/v1.0.0/azdoc-windows-amd64.exe" -OutFile "azdoc.exe"
```

### 🤖 Quick Start

```bash
# With AI (recommended)
export OPENAI_API_KEY="sk-..."
azdoc all --subscription-id <sub-id> --with-diagrams --enable-ai

# Without AI
azdoc all --subscription-id <sub-id> --with-diagrams
```

### 📄 Documentation

- [Full Documentation](https://github.com/automationpi/azdocs#readme)
- [Example Report](https://github.com/automationpi/azdocs/blob/master/examples/EXAMPLE_REPORT.md)
- [AI Features Guide](https://github.com/automationpi/azdocs/blob/master/AI_FEATURES.md)

### 🔒 Verify Checksums

```bash
curl -LO https://github.com/automationpi/azdocs/releases/download/v1.0.0/checksums.txt
sha256sum -c checksums.txt
```

### 💡 What's New in v1.0.0

- Initial release with full platform team features
- AI-powered security and cost analysis
- Executive summary dashboard
- Comprehensive Azure documentation generation
- Smart network diagram generation
```

---

## 🔄 Future Releases

For subsequent releases:

```bash
# Update version
git tag -a v1.1.0 -m "Release v1.1.0 - New Features"
git push origin v1.1.0

# GitHub Actions will automatically:
# - Build new binaries
# - Create new release
# - Update "latest" link
```

The `/releases/latest/download/` URL will always point to the newest release!

---

## ✅ Checklist for First Release

- [ ] Code pushed to GitHub
- [ ] GitHub Actions workflow committed (`.github/workflows/release.yml`)
- [ ] Tag created (`git tag -a v1.0.0 -m "Release v1.0.0"`)
- [ ] Tag pushed (`git push origin v1.0.0`)
- [ ] GitHub Actions workflow completed successfully
- [ ] Release visible at: https://github.com/automationpi/azdocs/releases
- [ ] Binaries downloadable via curl
- [ ] Test download and execute binary

---

## 🧪 Testing the Release

After creating the release, test the download command:

```bash
# Download
curl -LO https://github.com/automationpi/azdocs/releases/latest/download/azdoc-linux-amd64

# Make executable
chmod +x azdoc-linux-amd64

# Test
./azdoc-linux-amd64 version
```

Expected output:
```
azdoc version v1.0.0 (commit: abc123, built: 2025-01-24T10:30:00Z)
```

---

## 🎯 Current Status

✅ GitHub Actions workflow created
✅ Automated release process configured
⏳ Waiting for: First push to GitHub and tag creation

**Next steps:**
1. Push code: `git push -u origin master`
2. Create tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
3. Push tag: `git push origin v1.0.0`
4. Watch GitHub Actions build and release!
