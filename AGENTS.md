## Cursor Cloud specific instructions

This is the **Snowflake Terraform Provider** — a Go-based Terraform provider that manages Snowflake cloud data warehouse resources via Infrastructure as Code. There are no local databases or services; all integration/acceptance tests require a live Snowflake account.

### Key Commands

All standard development commands are in the `Makefile`. The most relevant for daily development:

| Task | Command |
|------|---------|
| Build binary | `make build-local` |
| Lint | `make lint` |
| Unit tests | `make test-unit` |
| Architecture tests | `make test-architecture` |
| Format code | `make fmt` |
| Pre-push checks | `make pre-push-check` |
| Terraform fmt check | `make terraform-fmt-check` |

### Gotchas

- **`pkg/sdk/client_integration_test.go`** contains `TestClient_*` tests that have NO build tag and require a live Snowflake connection. When running `make test-unit`, these will fail without `~/.snowflake/config`. To run SDK unit tests only, exclude `TestClient_*` tests: `go test -run "^Test[^C]|^TestC[^l]" ./pkg/sdk/`.
- **Integration and acceptance tests** (`make test-integration`, `make test-acceptance`) all require a live Snowflake account configured in `~/.snowflake/config` with profiles `[default]` and `[secondary_test_account]`.
- **Dev tools** (`golangci-lint`, `tfplugindocs`, `gofumpt`) are installed locally to `./bin/` and `./tools/bin/` via `make dev-setup`, not globally.
- **Terraform CLI** is needed for `make terraform-fmt-check` and `make terraform-fmt` (formats HCL in `./examples/` and `./pkg/testacc/testdata/`).
- The Go version required is specified in `go.mod` (currently Go 1.25.7).
