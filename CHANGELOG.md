# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- Complete CLI framework with Cobra + Viper
- Azure authentication (Azure CLI + Service Principal)
- Comprehensive data models for Azure network resources
- Examples directory with deployment scripts
- Test environment deployment (Hub-and-Spoke topology)
- Helper script to find Azure subscription ID
- Cost estimation and checking scripts
- Comprehensive documentation (4,500+ lines)

### Fixed
- **2024-10-24**: Fixed `--subscription-id` flag not working in `all` command
  - Issue: Circular dependency when inheriting flags from subcommands
  - Solution: Explicitly define all flags in `all` command
  - Commands affected: `azdoc all`

## [0.1.0] - 2024-10-24

### Initial Release
- Foundation complete
- CLI skeleton working
- All commands functional (scan, build, explain, all, doctor, version)
- Test environment ready
- Documentation complete

---

## How to Use This Changelog

- **Added**: New features
- **Changed**: Changes in existing functionality
- **Deprecated**: Soon-to-be removed features
- **Removed**: Removed features
- **Fixed**: Bug fixes
- **Security**: Security fixes
