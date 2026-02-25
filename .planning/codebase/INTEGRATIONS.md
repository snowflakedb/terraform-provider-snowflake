# External Integrations

**Analysis Date:** 2026-02-25

## APIs & External Services

**Snowflake Data Warehouse:**
- Snowflake REST API and SQL - Primary integration for all infrastructure management
  - SDK/Client: `github.com/snowflakedb/gosnowflake` (v1.18.1)
  - Auth: Multiple methods supported via `SNOWFLAKE_AUTHENTICATOR`:
    - Basic: `SNOWFLAKE_USER`, `SNOWFLAKE_PASSWORD`
    - Certificate-based: `SNOWFLAKE_PRIVATE_KEY_PATH`, `SNOWFLAKE_PRIVATE_KEY`, `SNOWFLAKE_PRIVATE_KEY_PASSPHRASE`
    - OAuth: `SNOWFLAKE_TOKEN`, `SNOWFLAKE_AUTHENTICATOR=OAUTH`
    - JWT: `SNOWFLAKE_PRIVATE_KEY`, `SNOWFLAKE_AUTHENTICATOR=JWT`
    - SAML: `SNOWFLAKE_AUTHENTICATOR=SAML`
    - External Browser: `SNOWFLAKE_AUTHENTICATOR=EXTERNALBROWSER`
    - Okta: `SNOWFLAKE_OKTA_URL` + `SNOWFLAKE_AUTHENTICATOR=OKTA`

**Okta (Optional):**
- Used for SAML-based authentication if configured
- Connection: `SNOWFLAKE_OKTA_URL` environment variable
- SDK: gosnowflake (native support)
- Tests: `TEST_SF_TF_SKIP_SAML_INTEGRATION_TEST` skips these when not configured

## Data Storage

**Databases:**
- Snowflake Cloud Data Warehouse - Primary data store
  - Connection: Via `SNOWFLAKE_ACCOUNT_NAME`, `SNOWFLAKE_ORGANIZATION_NAME`, `SNOWFLAKE_WAREHOUSE`, `SNOWFLAKE_ROLE`
  - Client: `github.com/snowflakedb/gosnowflake` (v1.18.1)
  - Protocol: HTTP/HTTPS with configurable endpoints

**File Storage:**
- Snowflake Internal Stages - Terraform state and object configurations stored within Snowflake
- Cloud Storage (via Snowflake):
  - AWS S3 - Supported through Snowflake External Stages
  - Azure Blob Storage - Supported through Snowflake External Stages
  - Google Cloud Storage - Supported through Snowflake External Stages
- Local filesystem - For provider configuration files (`~/.snowflake/config.toml`)

**Caching:**
- None - Direct connection to Snowflake for each operation
- Session management: `SNOWFLAKE_KEEP_SESSION_ALIVE` controls session persistence

## Authentication & Identity

**Auth Provider:**
- Snowflake Native - Custom authentication handling via gosnowflake driver
  - Implementation: Multiple auth methods (user/password, JWT, OAuth, SAML, external browser, certificate)
  - Profile support: TOML-based configuration at `~/.snowflake/config.toml`
  - Token accessor pattern: Optional OAuth token refresh flow via `token_accessor` block

**Multi-Factor Authentication:**
- Duo - Supported via Okta integration if configured
  - Configuration: `SNOWFLAKE_PASSCODE` or `SNOWFLAKE_PASSCODE_IN_PASSWORD`

## Monitoring & Observability

**Error Tracking:**
- None - Errors surfaced directly to Terraform execution context

**Logs:**
- Approach: Terraform Plugin Logging (via terraform-plugin-log)
- Provider version tracking: `internal/tracking/version.go` tracks provider version for debugging
- Query tracking: `internal/tracking/query.go` tracks executed queries
- Driver-level tracing: `SNOWFLAKE_DRIVER_TRACING=debug` enables gosnowflake debug output

## CI/CD & Deployment

**Hosting:**
- Terraform Registry (official distribution)
- GitHub Releases (via GoReleaser)
- Local plugin directory: `~/.terraform.d/plugins/`

**CI Pipeline:**
- GitHub Actions (implied by `.github/` structure and GoReleaser GitHub release configuration)
- Docker Compose (for containerized test execution):
  - `packaging/docker-compose.yml` orchestrates test environment
  - Supports pre-prod and gov environments via `TEST_SF_TF_SNOWFLAKE_TESTING_ENVIRONMENT`

**Release Process:**
- GoReleaser - Builds and releases provider binaries
  - Multi-platform: Windows, Linux, Darwin, FreeBSD
  - Multi-arch: amd64, arm64, 386
  - Signing: GPG signing with environment-based key management
  - Registry manifest generation for Terraform Registry

## Environment Configuration

**Required env vars:**
- `SNOWFLAKE_ACCOUNT_NAME` or `SNOWFLAKE_ORGANIZATION_NAME` + `SNOWFLAKE_ACCOUNT_NAME` - Account identification
- At least one auth method:
  - `SNOWFLAKE_USER` + `SNOWFLAKE_PASSWORD`
  - `SNOWFLAKE_PRIVATE_KEY_PATH`/`SNOWFLAKE_PRIVATE_KEY`
  - `SNOWFLAKE_TOKEN` (with `SNOWFLAKE_AUTHENTICATOR=OAUTH`)

**Optional connection tuning:**
- `SNOWFLAKE_WAREHOUSE` - Default warehouse for queries
- `SNOWFLAKE_ROLE` - Default role
- `SNOWFLAKE_PROTOCOL` - HTTP or HTTPS (default: HTTPS)
- `SNOWFLAKE_HOST` - Custom host for PrivateLink
- `SNOWFLAKE_PORT` - Custom port
- `SNOWFLAKE_CLIENT_IP` - For network policy checks

**Secrets location:**
- Environment variables (recommended for CI/CD)
- TOML profile file: `~/.snowflake/config.toml` (local development)
- Note: Never commit `.env` or credential files - `.gitignore` excludes sensitive files

**Test-specific env vars:**
- `TF_ACC=1` - Enable acceptance tests
- `TEST_SF_TF_REQUIRE_TEST_OBJECT_SUFFIX=1` - Test isolation
- `TEST_SF_TF_REQUIRE_GENERATED_RANDOM_VALUE=1` - Randomization
- `SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES=true` - Preview features
- `TEST_SF_TF_SKIP_SAML_INTEGRATION_TEST=true` - Skip SAML tests
- `TEST_SF_TF_SKIP_MANAGED_ACCOUNT_TEST=true` - Skip managed account tests
- `TEST_SF_TF_SNOWFLAKE_TESTING_ENVIRONMENT` - Environment selection (PRE_PROD_GOV, etc.)
- `SNOWFLAKE_DRIVER_TRACING=debug` - Driver debug logging

## Webhooks & Callbacks

**Incoming:**
- None - Provider is unidirectional (Terraform → Snowflake)

**Outgoing:**
- Snowflake Event Notifications (for resources that support it)
- Task notifications and alert webhooks (configured within Snowflake resources)
- No direct provider-initiated callbacks

## Network Configuration

**Connectivity:**
- Snowflake REST API - HTTPS (default) or HTTP (custom)
- Custom endpoints: `SNOWFLAKE_HOST` + `SNOWFLAKE_PORT`
- PrivateLink support: Via custom host configuration
- OCSP certificate validation: `SNOWFLAKE_OCSP_FAIL_OPEN`, `SNOWFLAKE_DISABLE_OCSP_CHECKS`

**Connection Pooling:**
- Managed by gosnowflake driver
- Session persistence: `SNOWFLAKE_KEEP_SESSION_ALIVE`
- Timeout configuration:
  - `SNOWFLAKE_LOGIN_TIMEOUT` - Login retry timeout
  - `SNOWFLAKE_REQUEST_TIMEOUT` - Request timeout
  - `SNOWFLAKE_CLIENT_TIMEOUT` - Client authentication timeout
  - `SNOWFLAKE_JWT_CLIENT_TIMEOUT` - JWT authentication timeout
  - `SNOWFLAKE_EXTERNAL_BROWSER_TIMEOUT` - External browser auth timeout
  - `SNOWFLAKE_JWT_EXPIRE_TIMEOUT` - JWT token expiration

## OAuth Token Refresh

**Token Accessor Pattern:**
- Optional OAuth token refresh mechanism via `token_accessor` block in provider config
- Fields:
  - `token_endpoint` - OAuth provider token endpoint (env: `SNOWFLAKE_TOKEN_ACCESSOR_TOKEN_ENDPOINT`)
  - `refresh_token` - OAuth refresh token (env: `SNOWFLAKE_TOKEN_ACCESSOR_REFRESH_TOKEN`)
  - `client_id` - OAuth client ID (env: `SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_ID`)
  - `client_secret` - OAuth client secret (env: `SNOWFLAKE_TOKEN_ACCESSOR_CLIENT_SECRET`)
  - `redirect_uri` - OAuth redirect URI (env: `SNOWFLAKE_TOKEN_ACCESSOR_REDIRECT_URI`)

## Terraform Integration Points

**Provider Registry:**
- Namespace: `snowflakedb/snowflake`
- Deployed to: Terraform Registry (registry.terraform.io)
- Protocol: gRPC v6 (tfprotov6)

**State Management:**
- Terraform state format - Standard Terraform JSON
- State encryption: Handled by Terraform (provider doesn't manage)
- State import: Supported for most resources via `terraform import`

---

*Integration audit: 2026-02-25*
