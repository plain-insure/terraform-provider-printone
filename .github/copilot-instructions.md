# Terraform Provider PrintOne

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

The terraform-provider-printone is a Terraform Provider built with the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). This provider enables management of PrintOne webhook resources through Terraform.

## Required Dependencies

**CRITICAL**: Ensure these exact dependency versions before working:
- Go >= 1.23 (currently tested with Go 1.24.5)
- Terraform >= 1.0 (currently tested with Terraform 1.12.2)
- golangci-lint for linting (install separately)

Install Go and Terraform first, then golangci-lint:
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
- `make lint` FAILS due to resource constraints in most environments
- golangci-lint gets killed due to memory limitations
- **WORKAROUND**: Use lighter linting or skip during development
- The .golangci.yml config was fixed to remove non-existent 'usetesting' linter
- If lint is required, try: `golangci-lint run --timeout=5m --concurrency=1`

**DOCUMENTATION GENERATION ISSUES**:
- `make generate` FAILS due to template configuration issues
- The provider name "scaffolding" in tools/tools.go needs to be updated to "printone"
- **WORKAROUND**: Documentation generation currently broken, manual docs updates required

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
├── go.mod                       # Go module definition (Go 1.23.7)
├── GNUmakefile                  # Build targets: fmt, lint, test, testacc, build, install, generate
├── main.go                      # Provider entry point
├── internal/provider/           # Provider implementation
│   ├── provider.go             # Main provider configuration
│   ├── webhook_resource.go     # Webhook resource implementation
│   ├── webhook_data_source.go  # Webhook data source implementation
│   └── provider_test.go        # Provider tests
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

1. **Linting**: `make lint` fails due to memory constraints
2. **Documentation**: `make generate` fails due to provider name mismatch
3. **No real tests**: Project contains mostly scaffolding, minimal test coverage
4. **Provider name**: Some files still reference "scaffolding" instead of "printone"

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