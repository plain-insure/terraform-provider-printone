# Terraform Provider PrintOne

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

The terraform-provider-printone is a Terraform Provider built with the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). This provider enables management of PrintOne webhook resources through Terraform.

## Required Dependencies

**CRITICAL**: Ensure these exact dependency versions before working:
- Go >= 1.24.0 (currently tested with Go 1.24.7)
- Terraform >= 1.0 (optional, only needed for documentation generation)
- golangci-lint for linting (install separately if needed)

Install golangci-lint if required for linting:
```bash
# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.62.2
export PATH=$PATH:$(go env GOPATH)/bin
```

## Working Effectively

**CRITICAL BUILD TIMING**: Follow exact sequence and timeouts:

1. **Bootstrap dependencies** (takes ~40 seconds):
   ```bash
   go mod download
   ```

2. **Build the provider** (takes ~25 seconds):
   ```bash
   go build -v ./...
   ```
   OR use make targets:
   ```bash
   make build    # Same as go build -v ./...
   make install  # Build + install binary (takes ~2 seconds)
   ```

3. **Format code** (takes <1 second):
   ```bash
   make fmt
   ```

4. **Run tests** (takes ~6 seconds):
   ```bash
   make test
   ```

5. **Run acceptance tests** (takes <1 second - no real tests exist yet):
   ```bash
   make testacc
   ```
   **WARNING**: Acceptance tests can create real resources and cost money. Currently no tests exist.

## Critical Limitations and Workarounds

**LINTING ISSUES**: 
- `make lint` FAILS if golangci-lint is not installed
- **WORKAROUND**: Install golangci-lint first or skip linting during development
- The .golangci.yml config is properly configured and functional when golangci-lint is available
- If lint is required, try: `golangci-lint run --timeout=5m --concurrency=1`

**DOCUMENTATION GENERATION ISSUES**:
- `make generate` FAILS if Terraform CLI is not installed in PATH
- The provider name is correctly set to "printone" in tools/tools.go
- **WORKAROUND**: Install Terraform CLI or skip documentation generation (docs are already generated)
- Documentation generation works correctly when Terraform is available

## Validation Scenarios

After making changes, always:

1. **Build validation**: Run `go build -v ./...` to ensure code compiles
2. **Format validation**: Run `make fmt` to fix formatting
3. **Test validation**: Run `make test` to check unit tests pass
4. **Provider functionality**: Test the built provider binary:
   ```bash
   go install .
   # Binary installed at: /home/runner/go/bin/terraform-provider-printone
   /home/runner/go/bin/terraform-provider-printone --help
   ```

## Repository Structure

Key directories and files:
```
.
├── README.md                    # Main documentation
├── go.mod                       # Go module definition (Go 1.24.0)
├── GNUmakefile                  # Build targets: fmt, lint, test, testacc, build, install, generate
├── main.go                      # Provider entry point
├── internal/provider/           # Provider implementation
│   ├── provider.go             # Main provider configuration
│   ├── webhook_resource.go     # Webhook resource implementation
│   ├── webhook_data_source.go  # Webhook data source implementation
│   ├── webhooksecret_data_source.go # Webhook secret data source
│   ├── webhook_helpers.go      # Webhook helper functions
│   ├── resource_webhook/       # Webhook resource module
│   ├── datasource_webhook/     # Webhook data source module
│   ├── datasource_webhooksecret/ # Webhook secret data source module
│   └── provider_printone/      # Provider module
├── examples/                    # Terraform examples for documentation
├── docs/                        # Generated documentation
├── test/                        # Test configuration files
├── tools/                       # Documentation generation tools
└── .github/workflows/           # CI/CD workflows
```

## Common Tasks and Commands

**Development workflow**:
```bash
# 1. Install dependencies (40 seconds)
go mod download

# 2. Build and install (30 seconds total)
make install

# 3. Format code (<1 second)
make fmt

# 4. Run tests (6 seconds)
make test

# 5. Test provider binary
go install .
$GOPATH/bin/terraform-provider-printone --help
```

**CI Pipeline compatibility**:
- `.github/workflows/test.yml` runs build, lint, generate, and acceptance tests
- Uses Go matrix testing with multiple Terraform versions (1.0.* through 1.4.*)
- Build timeout: 5 minutes, Test timeout: 15 minutes
- **NEVER CANCEL** builds or tests - they complete quickly for this project

## Known Issues

1. **Linting**: `make lint` fails if golangci-lint is not installed
2. **Documentation**: `make generate` fails if Terraform CLI is not installed
3. **No real tests**: Project contains mostly scaffolding, minimal test coverage
4. **Provider functionality**: Limited to webhook management resources currently

## Provider Configuration

The provider address is: `registry.terraform.io/plain-insure/printone`

Resources available:
- `printone_webhook` (resource)
- `printone_webhook` (data source)

## Manual Testing Approach

Since automated tests are minimal:
1. Build the provider: `go install .`
2. Create a test Terraform configuration using the provider
3. Run `terraform init` and `terraform plan` to verify provider loads
4. Test resource operations with real or mock endpoints

**CRITICAL**: Always build and validate your changes work before committing.

## Code Style and Conventions

**Go Code Style**:
- Follow standard Go conventions and formatting (use `make fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and small

**Terraform Provider Patterns**:
- Resource implementations follow the Terraform Plugin Framework patterns
- Use proper validation and error handling
- Implement Create, Read, Update, Delete (CRUD) operations consistently
- Use appropriate schema definitions with proper types and validation

**File Organization**:
- Provider resources: `internal/provider/resource_*.go`
- Data sources: `internal/provider/datasource_*.go`
- Client logic: `internal/client/`
- Examples: `examples/` (for documentation generation)
- Tests: `*_test.go` files alongside implementation

## Security Considerations

- API keys should be handled securely (prefer environment variables)
- Sensitive data should be marked as sensitive in schema
- Validate all inputs and sanitize outputs
- Follow least privilege principles for API access

## Troubleshooting Common Issues

**Build Failures**:
- Run `go mod download` first to ensure dependencies are available
- Check Go version compatibility (requires Go 1.24.0+)
- Verify all imports are available and correct

**Test Failures**:
- Most tests are currently minimal/scaffolding
- Focus on testing business logic and error cases
- Use testify/assert for cleaner test assertions when available

**Provider Runtime Issues**:
- Check API connectivity and authentication
- Verify Terraform version compatibility
- Review provider configuration and required fields

## Development Workflow Best Practices

**Before Making Changes**:
1. Pull latest changes and ensure clean build state
2. Run `go mod download` to ensure dependencies are current
3. Test basic functionality with `make test`

**During Development**:
1. Make small, incremental changes
2. Run `make fmt` frequently to maintain code style
3. Test changes early and often with `make build && make test`
4. Use meaningful commit messages describing changes

**Before Committing**:
1. Run full validation: `make build && make fmt && make test`
2. Test provider binary functionality if applicable
3. Update documentation if public interfaces changed
4. Verify no unintended files are being committed

## Working with External APIs

**PrintOne API Integration**:
- Base API endpoint: `https://api.print.one`
- Authentication via `x-api-key` header
- Support for custom endpoints via provider configuration
- Client implementation in `internal/client/`

**Testing API Integration**:
- Use mock responses for unit tests when possible
- Consider API rate limits during development
- Validate error handling for API failures
- Test with invalid/expired API keys

## Environment Setup Tips

**Go Environment**:
- Ensure `GOPATH` and `GOROOT` are properly set
- Use Go modules (this project uses go.mod)
- Recommended to use latest stable Go version (1.24.7+)

**Development Tools**:
- IDE: VS Code with Go extension recommended
- Formatting: Built into make targets (`make fmt`)
- Testing: Standard go test tools
- Debugging: Use `go install .` for local provider testing