# Technology Stack

**Analysis Date:** 2026-02-25

## Languages

**Primary:**
- Go 1.24.9 - Main provider implementation and SDK code
- HCL (HashiCorp Configuration Language) - Terraform configuration examples

**Secondary:**
- Bash - Build scripts and test orchestration
- TOML - Configuration file format for Snowflake connection profiles

## Runtime

**Environment:**
- Go 1.24.9 - Required runtime for compilation and execution

**Package Manager:**
- Go Modules - Dependency management via `go mod`
- Lockfile: Present (`go.sum`)

## Frameworks

**Core:**
- HashiCorp Terraform Plugin SDK v2 (v2.38.1) - Legacy provider framework (SDKv2)
- HashiCorp Terraform Plugin Framework (v1.17.0) - Modern provider framework
- HashiCorp Terraform Plugin Go (v0.29.0) - gRPC server implementation for provider protocol
- HashiCorp Terraform Plugin Mux (v0.21.0) - Multiplexing provider servers (SDKv2 + Framework)
- HashiCorp Terraform Plugin Testing (v1.14.0) - Integration testing framework

**Database:**
- Snowflake Go Driver (gosnowflake v1.18.1) - Native Go driver for Snowflake connections via SQL
- sqlx (v1.4.0) - Extension to database/sql with reflection and advanced scanning

**Build/Dev:**
- GoReleaser - Binary building and release management (`.goreleaser.yml`)
- golangci-lint (v2.6.1) - Go linter and formatter orchestration
- gofumpt (v0.9.2) - Opinionated Go formatter
- tfplugindocs (v0.24.0) - Terraform provider documentation generator

**Code Generation:**
- Terraform Plugin Framework generator tools - For SDK scaffolding and code generation
- Custom code generation scripts in `Makefile` for:
  - SDK objects (`generate-sdk`)
  - Show output schemas (`generate-show-output-schemas`)
  - Resource and datasource assertions
  - Test configuration model builders

## Key Dependencies

**Critical:**
- `github.com/snowflakedb/gosnowflake` (v1.18.1) - Native driver for Snowflake REST API communication, handles authentication, connection management, and SQL execution
- `github.com/hashicorp/terraform-plugin-sdk/v2` (v2.38.1) - Provider lifecycle management, schema definition, CRUD operations
- `github.com/hashicorp/terraform-plugin-framework` (v1.17.0) - Modern framework for future provider development (currently experimental/PoC)

**Infrastructure:**
- `github.com/hashicorp/terraform-plugin-mux` (v0.21.0) - Combines SDKv2 and Plugin Framework servers in single provider
- `github.com/hashicorp/terraform-plugin-testing` (v1.14.0) - Acceptance test utilities
- `github.com/jmoiron/sqlx` (v1.4.0) - Enhanced SQL database operations
- `github.com/hashicorp/hcl` (v1.0.0) - HCL parsing for configuration files

**Cloud SDKs:**
- `github.com/Azure/azure-sdk-for-go/sdk/storage/azblob` (v1.6.3) - Azure Blob Storage support (transitive via gosnowflake)
- `github.com/aws/aws-sdk-go-v2/*` - AWS SDK components (transitive via gosnowflake for S3, STS, SSO, OIDC)

**Cryptography & Security:**
- `github.com/youmark/pkcs8` (v0.0.0-20240726163527-a2c0da244d78) - PKCS#8 private key support for certificate-based auth
- `golang.org/x/crypto` (v0.45.0) - Cryptographic functions for TLS/SSL operations

**Testing & Utilities:**
- `github.com/stretchr/testify` (v1.11.1) - Test assertions and mocking
- `github.com/brianvoe/gofakeit/v6` (v6.28.0) - Fake data generation for tests

**Serialization:**
- `github.com/pelletier/go-toml/v2` (v2.2.4) - TOML parsing for connection profile files
- `github.com/vmihailenco/msgpack` (v5.4.1) - Message pack serialization (transitive)

## Configuration

**Environment:**
Authentication and connection configured via:
- Environment variables with `SNOWFLAKE_` prefix (see `pkg/internal/snowflakeenvs/`):
  - `SNOWFLAKE_ACCOUNT_NAME` - Snowflake account identifier
  - `SNOWFLAKE_ORGANIZATION_NAME` - Snowflake organization (for new account identifiers)
  - `SNOWFLAKE_USER` - Database user
  - `SNOWFLAKE_PASSWORD` - User password or PAT token
  - `SNOWFLAKE_WAREHOUSE` - Default warehouse
  - `SNOWFLAKE_ROLE` - Default role
  - `SNOWFLAKE_AUTHENTICATOR` - Authentication type (USER_PASSWORD, JWT, OAUTH, SAML, EXTERNALBROWSER, etc.)
  - `SNOWFLAKE_PRIVATE_KEY_PATH` / `SNOWFLAKE_PRIVATE_KEY` - For key-pair authentication
  - `SNOWFLAKE_TOKEN` - OAuth token
  - Connection tuning: `SNOWFLAKE_LOGIN_TIMEOUT`, `SNOWFLAKE_REQUEST_TIMEOUT`, `SNOWFLAKE_CLIENT_TIMEOUT`, etc.
  - Security: `SNOWFLAKE_INSECURE_MODE`, `SNOWFLAKE_OCSP_FAIL_OPEN`, `SNOWFLAKE_DISABLE_OCSP_CHECKS`

- TOML profile file: `~/.snowflake/config.toml` (optional, legacy format also supported)
  - Stored in `pkg/sdk/config_dto.go` and related DTO builders

**Build:**
- `.goreleaser.yml` - Multi-platform binary builds (Windows, Linux, Darwin, FreeBSD on amd64, arm64, 386)
- `Makefile` - Primary orchestration for build, test, lint, code generation
- `go.mod` / `go.sum` - Dependency management with version pinning

## Platform Requirements

**Development:**
- Go 1.24.9 or higher
- golangci-lint (installed via `make dev-setup`)
- tfplugindocs (installed via `make dev-setup`)
- gofumpt (installed via `make dev-setup`)
- Make or compatible build tool
- Git (for version tracking and hooks)

**Testing:**
- Active Snowflake account with appropriate permissions
- Terraform (TF_ACC=1 enables acceptance tests)
- Docker/Docker Compose (for containerized test execution)
- Test environment variables (see Makefile test targets):
  - `TF_ACC=1` - Enables acceptance tests
  - `TEST_SF_TF_REQUIRE_TEST_OBJECT_SUFFIX=1` - Ensures test object isolation
  - `TEST_SF_TF_REQUIRE_GENERATED_RANDOM_VALUE=1` - Randomization for test isolation
  - `SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES=true` - Enables preview features during testing
  - `SNOWFLAKE_DRIVER_TRACING=debug` - Driver-level tracing

**Production:**
- Binary deployment to Terraform Registry or local plugin cache
- Terraform CLI 1.0+ (supports gRPC protocol v6)
- Snowflake account with network connectivity
- Credentials available via environment variables or config file

---

*Stack analysis: 2026-02-25*
