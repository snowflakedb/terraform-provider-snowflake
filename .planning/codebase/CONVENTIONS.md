# Coding Conventions

**Analysis Date:** 2025-02-25

## Naming Patterns

**Files:**
- Snake case for resource/datasource files: `account.go`, `account_parameter.go`
- Acceptance test files: `resource_<name>_acceptance_test.go` or `data_source_<name>_acceptance_test.go`
- Helper files: `resource_helpers_<purpose>.go` (e.g., `resource_helpers_create.go`, `resource_helpers_update.go`)
- Test files co-located with source: `*_test.go` in same package
- Generate directives: `generate.go` for codegen targets

**Functions:**
- PascalCase for exported functions (Go convention): `CreateAccount`, `ImportAccount`, `DeleteAccount`
- camelCase for unexported functions: `accountAuthenticationPolicyAttachmentConfig`, `stringAttributeCreate`
- CRUD operations: `Create<Resource>`, `Read<Resource>`, `Update<Resource>`, `Delete<Resource>`
- Import functions: `Import<Resource>`
- Helper functions: verbose and descriptive: `suppressIdentifierQuoting`, `IgnoreAfterCreation`, `NormalizeAndCompare`

**Variables:**
- camelCase throughout: `accountId`, `currentRole`, `oauthRefreshToken`
- Pointer variables and config fields: `opts`, `createResponse`, `d` (terraform schema.ResourceData)
- Test setup variables: `basic`, `complete`, `assertBasic`, `assertComplete`
- Error variables: always `err`

**Types/Constants:**
- PascalCase for types: `Task`, `Secret`, `SchemaObjectIdentifier`
- UPPER_CASE for constants: `BooleanDefault`, `BooleanTrue`, `IntDefault`
- Map schema variables: snake_case with `Schema` suffix: `accountSchema`, `alertsSchema`

**SDK Identifiers:**
- Fully qualified names use dot separator: `"db"."schema"."table"` (quoting per part)
- Unquoted parts use dot separator: `db.schema.table`
- Identifier types use abbreviations: `NewSchemaObjectIdentifier`, `NewAccountObjectIdentifier`, `NewDatabaseObjectIdentifier`

## Code Style

**Formatting:**
- Tool: `gofumpt` (enforced via `make fmt`)
- Also uses `goimports` for import organization
- Linted with `golangci-lint`

**Linting:**
- Config: `.golangci.yml`
- Enabled linters: errcheck, errname, errorlint, gocritic, gosec, govet, ineffassign, makezero, misspell, prealloc, revive, staticcheck, testifylint, thelper, unconvert, wastedassign, whitespace
- Disabled testifylint rules: `require-error` (allow assert), `empty` (allow Equal with empty string)

**Line Length:**
- Implicit limit based on gofumpt defaults (typically 80-100 chars)
- Long strings and function calls broken into multiple lines

**Indentation:**
- Tabs (Go standard)

## Import Organization

**Order:**
1. Standard library: `context`, `errors`, `fmt`, `log`, `strings`, `time`, `reflect`
2. External packages: `github.com/hashicorp/*`, `github.com/stretchr/*`
3. Internal packages: `github.com/Snowflake-Labs/terraform-provider-snowflake/*`

**Path Aliases:**
- No aliases used; full import paths
- Internal imports use full path: `"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"`

**Example from account.go:**
```go
import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)
```

## Error Handling

**Pattern - Simple errors:**
```go
if err != nil {
    return diag.FromErr(err)  // For diagnostic returns
}
```

**Pattern - Wrapped errors with context:**
```go
if err != nil {
    return diag.FromErr(fmt.Errorf("failed to query account (%s) after creation, err: %w", id.FullyQualifiedName(), err))
}
```

**Pattern - Error checking with specific types:**
```go
if errors.Is(err, sdk.ErrObjectNotFound) {
    return nil, false  // Retryable
} else {
    return err, true
}
```

**Pattern - Multiple operations with errors.Join:**
```go
if err := errors.Join(
    d.Set("edition", string(*account.Edition)),
    d.Set("region", account.SnowflakeRegion),
    d.Set("comment", comment),
); err != nil {
    return nil, err
}
```

**Pattern - Retry logic with error handling:**
```go
if err := util.Retry(5, 3*time.Second, func() (error, bool) {
    _, err = client.Accounts.ShowByID(ctx, id)
    if err != nil {
        log.Printf("[DEBUG] retryable operation resulted in error: %v", err)
        if errors.Is(err, sdk.ErrObjectNotFound) {
            return nil, false  // Retryable
        } else {
            return err, true   // Not retryable
        }
    }
    return nil, true
}); err != nil {
    return diag.FromErr(err)
}
```

**Special case - Parsing/Type conversion:**
```go
userType, err := sdk.ToUserType(v.(string))
if err != nil {
    return diag.FromErr(err)
}
```

**Pattern - Two-phase error handling:**
```go
if _, err := ImportName[sdk.AccountIdentifier](context.Background(), d, nil); err != nil {
    return diag.FromErr(err)
}
```

## Logging

**Framework:** Standard Go `log` package (output to `[DEBUG]` prefix)

**Patterns:**
- Debug logging for operations: `log.Printf("[DEBUG] retryable operation resulted in error: %v", err)`
- Used primarily for error context and retry logic
- Rarely used for success paths
- Debug level indicated by `[DEBUG]` prefix string

## Comments

**When to Comment:**
- Complex logic or non-obvious behavior
- Business logic that isn't self-evident
- TODO items with ticket references: `// TODO [SNOW-1763613]: unskip`
- Explanatory comments for test setup or configurations

**No JSDoc/TSDoc:**
- Go uses plain comments with Go conventions
- Export comments for public functions/types: `// CreateAccount creates a new account` (above the function)
- Often omitted for obvious CRUD operations

## Function Design

**Size:**
- Typical resource functions: 20-50 lines (Create, Read, Update)
- Helper functions: 5-15 lines (composition over monolithic functions)
- Schema definition functions: 50-200 lines (maps are large)

**Parameters:**
- Context passed as first parameter: `ctx context.Context`
- Data passed as pointer: `d *schema.ResourceData` (terraform convention)
- Meta cast to provider context: `meta.(*provider.Context)`
- Options structs for create/update operations: `opts := &sdk.CreateAccountOptions{}`

**Return Values:**
- CRUD operations: `diag.Diagnostics` (terraform diagnostic type)
- Read/Import: `[]*schema.ResourceData, error`
- Helper functions: `error` or typed results with error
- Test helpers: `func(*testing.T)` with assertion via assertions or return values

## Module Design

**Exports:**
- Resources exported via init functions: `func init() { resource.AddProvider(...) }`
- Datasources exported via init functions
- Each resource/datasource file defines its CRUD functions and schema
- No barrel files used; direct imports

**Resource File Pattern in `pkg/resources/`:**
```go
// account.go
var accountSchema = map[string]*schema.Schema{ ... }

func CreateAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics { ... }
func ReadAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics { ... }
func UpdateAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics { ... }
func DeleteAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics { ... }
func ImportAccount(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) { ... }
```

**Generated Code:**
- Marked with `// Code generated by...` comment at top
- Located in `pkg/schemas/` and `pkg/acceptance/` directories
- Minor adjustments allowed if annotated with `// Adjusted manually` comment
- Not to be modified directly; regenerate via `make` targets

## SDK Conventions

**Pointer wrapping:**
- Use `sdk.String()`, `sdk.Bool()`, `sdk.Int()` for pointer values
- Use `sdk.Pointer()` for non-primitive types: `sdk.Pointer(objectIdentifier)`

**Identifier handling:**
- Parse with SDK parsers: `sdk.ParseAccountIdentifier()`, `sdk.ParseSchemaObjectIdentifier()`
- Create with SDK constructors: `sdk.NewAccountObjectIdentifier()`
- Extract with methods: `id.FullyQualifiedName()`, `id.DatabaseName()`, `id.SchemaName()`

**Enum conversion:**
- Use SDK converters: `sdk.ToUserType()`, `sdk.ToAccountEdition()`
- Validate with converters (they return error for invalid values)

## Resource ID Encoding

**Current Pattern (dot separator):**
- Use `helpers.EncodeResourceIdentifier()` for setting IDs
- Use `sdk.Parse<Type>ObjectIdentifier()` for reading IDs
- Example: `helpers.EncodeResourceIdentifier(sdk.NewAccountIdentifier(orgName, accountName))`

**Old Pattern (DEPRECATED - pipe separator):**
- No longer used; being migrated away
- Don't add new resources with pipe separator

## Diff Suppression & Validation

**Diff Suppress Functions:**
- Used in schema for fields that should ignore changes in some cases
- Common examples: `IgnoreAfterCreation`, `NormalizeAndCompare`, `suppressIdentifierQuoting`
- Located in `common.go` or resource-specific files
- Composed with `SuppressIfAny()` for multiple conditions

**Validation Functions:**
- Use SDK validators: `sdkValidation(sdk.ToAccountEdition)`
- Converts enum validators to terraform ValidateDiagFunc
- Compose with `ValidateDiagFunc` field in schema

---

*Convention analysis: 2025-02-25*
