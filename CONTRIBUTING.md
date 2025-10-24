# Contributing to azdoc

Thank you for considering contributing to azdoc! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Azure subscription (for testing)
- Azure CLI installed and configured
- Git

### Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/azdoc.git
   cd azdoc
   ```

3. Install dependencies:
   ```bash
   make deps
   ```

4. Build the project:
   ```bash
   make build
   ```

5. Run tests:
   ```bash
   make test
   ```

## Project Structure

```
azdoc/
├── cmd/azdoc/              # CLI entry point
│   ├── main.go             # Main entry point
│   └── commands/           # Cobra command implementations
│       ├── root.go         # Root command and global flags
│       ├── scan.go         # Scan command
│       ├── build.go        # Build command
│       ├── explain.go      # LLM explanation command
│       ├── all.go          # All-in-one command
│       ├── doctor.go       # Diagnostics command
│       └── version.go      # Version command
├── pkg/                    # Public packages
│   ├── auth/               # Azure authentication
│   ├── cache/              # Local caching layer
│   ├── config/             # Configuration management
│   ├── discovery/          # Resource discovery
│   ├── graph/              # Topology graph builder
│   ├── llm/                # LLM integration
│   ├── models/             # Data models
│   │   ├── subscription.go # Subscription models
│   │   ├── network.go      # Network resource models
│   │   ├── security.go     # Security resource models
│   │   ├── routing.go      # Routing models
│   │   └── loadbalancer.go # Load balancer models
│   └── renderer/           # Output renderers
│       ├── markdown.go     # Markdown renderer
│       └── diagram.go      # Draw.io renderer
├── internal/               # Private packages
│   └── utils/              # Utility functions
├── tests/                  # Test files
│   ├── integration/        # Integration tests
│   ├── golden/             # Golden file tests
│   └── fixtures/           # Test fixtures
├── Makefile                # Build automation
├── README.md               # Project documentation
├── CONTRIBUTING.md         # This file
├── LICENSE                 # MIT License
└── azdoc.yaml.example      # Example configuration
```

## Development Workflow

### Making Changes

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and commit:
   ```bash
   git add .
   git commit -m "Add feature: description"
   ```

3. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

4. Create a pull request

### Commit Messages

Follow conventional commits:
- `feat: add new feature`
- `fix: fix bug in component`
- `docs: update documentation`
- `test: add tests for feature`
- `refactor: refactor code`
- `chore: update dependencies`

### Code Style

- Run `make fmt` before committing
- Follow Go standard naming conventions
- Add comments for exported functions and types
- Keep functions small and focused

### Testing

- Write unit tests for all new code
- Run `make test` to ensure all tests pass
- Add integration tests for end-to-end features
- Update golden files when output changes

## Architecture Guidelines

### Data Flow

1. **Discovery**: Azure Resource Graph → Raw JSON cache
2. **Normalization**: Raw JSON → Typed models
3. **Graph Building**: Typed models → Topology graph
4. **Rendering**: Topology graph → Markdown + Diagrams
5. **LLM (optional)**: Normalized data → Explanations

### Key Principles

- **Deterministic**: Same inputs → same outputs
- **Cacheable**: Support offline builds from cached data
- **Extensible**: Easy to add new resource types
- **Safe**: Never log secrets or PII
- **Fast**: Concurrent API calls, efficient caching

### Adding New Resource Types

1. Add model to `pkg/models/`
2. Add discovery logic to `pkg/discovery/`
3. Update graph builder in `pkg/graph/`
4. Update renderers in `pkg/renderer/`
5. Add tests

Example:

```go
// pkg/models/newresource.go
type NewResource struct {
    ID            string `json:"id"`
    Name          string `json:"name"`
    // ... other fields
}

// pkg/discovery/newresource.go
func (c *Client) DiscoverNewResources(ctx context.Context) ([]models.NewResource, error) {
    // Discovery logic
}

// pkg/graph/builder.go
func (b *Builder) addNewResourceNodes() error {
    // Add to graph
}

// pkg/renderer/markdown.go
func (r *MarkdownRenderer) renderNewResources() string {
    // Render in Markdown
}
```

## Testing Guidelines

### Unit Tests

```go
// pkg/graph/builder_test.go
func TestBuilder_Build(t *testing.T) {
    builder := NewBuilder(testData)
    topology, err := builder.Build()

    assert.NoError(t, err)
    assert.Equal(t, expectedNodeCount, len(topology.Nodes))
}
```

### Integration Tests

```go
// tests/integration/scan_test.go
func TestScan_RealSubscription(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Test against real Azure subscription
}
```

Run integration tests:
```bash
go test ./tests/integration/... -v
```

### Golden File Tests

For testing output formatting:

```go
// tests/golden/markdown_test.go
func TestMarkdownRenderer_Golden(t *testing.T) {
    output := renderer.Render(testTopology)
    golden.Assert(t, "expected.md", output)
}
```

## Pull Request Process

1. Update documentation if needed
2. Add tests for new features
3. Ensure all tests pass
4. Update CHANGELOG.md
5. Request review from maintainers

### PR Checklist

- [ ] Code follows project style
- [ ] Tests added and passing
- [ ] Documentation updated
- [ ] Commit messages follow conventions
- [ ] No secrets or PII in code

## Release Process

Maintainers will:

1. Update version in `Makefile`
2. Update `CHANGELOG.md`
3. Create git tag: `git tag v1.0.0`
4. Push tag: `git push origin v1.0.0`
5. GitHub Actions will build and release

## Getting Help

- Open an issue for bugs or feature requests
- Join discussions for questions
- Tag maintainers for urgent issues

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what's best for the project

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
