# Logo & Images Setup Instructions

The azdoc logo and example diagram have been integrated into the project, but the actual image files need to be placed in the correct location.

## Required Actions

### 1. Logo Image

Save the azdoc logo image to:
```
assets/azdoc-logo.png
```

**Specifications:**
- Format: PNG
- Recommended size: 512x512 pixels or larger (will be scaled in display)
- Transparent background preferred

### 2. Example Network Diagram

Save the example network diagram (vnet-spoke1 screenshot) to:
```
assets/example-diagram.png
```

**Specifications:**
- Format: PNG
- The diagram showing "Virtual Network: vnet-spoke1, Resource Group: azdoc-test-rg"
- Will be displayed at 800px width in README

## Where the Logo Appears

### README.md
- Displays at the top of the README
- Size: 300px width
- Centered alignment

### Generated Documentation Reports
- Appears at the top of each generated `SUBSCRIPTION.md` file
- Size: 200px width
- Centered alignment
- Relative path: `../assets/azdoc-logo.png` (from docs/ directory)

## File Structure

```
azure-docs/
├── assets/
│   └── azdoc-logo.png          ← Place logo here
├── README.md                    ← Logo referenced here
└── docs/
    └── SUBSCRIPTION.md          ← Logo referenced here (generated)
```

## Next Steps

Once you place the logo file at `assets/azdoc-logo.png`, both the README and generated documentation will display it correctly.
