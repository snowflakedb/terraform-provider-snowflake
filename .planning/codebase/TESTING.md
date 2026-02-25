# Testing Patterns

**Analysis Date:** 2025-02-25

## Test Framework

**Runner:**
- Go's built-in `testing` package
- Terraform Plugin Testing Framework: `github.com/hashicorp/terraform-plugin-testing`
- Config: Makefile targets for different test types

**Assertion Library:**
- `github.com/stretchr/testify/assert` - Soft assertions (continue on failure)
- `github.com/stretchr/testify/require` - Hard assertions (fail immediately)

**Run Commands:**
```bash
make test-unit              # Run unit tests
make test-acceptance        # Run acceptance tests
make test-account-level-features  # Run account-level acceptance tests
make test-integration       # Run SDK integration tests
make test-functional        # Run functional tests for terraform libraries
make test-architecture      # Check architecture constraints
make test-acceptance-<Resource>  # Run acceptance tests for specific resource
make test-main-terraform-use-cases  # Run main use case tests
```

## Test File Organization

**Location:**
- Unit tests co-located with source: `*.go` and `*_test.go` in same directory
- Acceptance tests: `pkg/testacc/resource_<name>_acceptance_test.go` or `pkg/testacc/data_source_<name>_acceptance_test.go`
- Integration tests (SDK): `pkg/sdk/testint/` directory
- Functional tests: `pkg/testfunctional/` directory

**Naming:**
- Unit test files: `<source_file>_test.go`
- Acceptance test files: `resource_<resource_name>_acceptance_test.go`, `data_source_<datasource_name>_acceptance_test.go`
- Test functions: `Test<Function>`, `TestAcc_<ResourceName>_<UseCase>`, `TestInt_<Description>`

**Packages and Build Tags:**
```go
//go:build non_account_level_tests

package testacc
```
- Build tags separate account-level tests from non-account-level tests
- Acceptance tests tagged: `non_account_level_tests` or `account_level_tests`

## Test Structure

**Unit Test Pattern:**
```go
package resources

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func Test_suppressIdentifierQuoting(t *testing.T) {
	t.Run("old identifier with too many parts", func(t *testing.T) {
		result := suppressIdentifierQuoting("", incorrectId, firstId, nil)
		require.False(t, result)
	})

	t.Run("identifiers the same (but different quoting)", func(t *testing.T) {
		result := suppressIdentifierQuoting("", firstId, firstIdQuoted, nil)
		require.True(t, result)
	})
}
```

**Acceptance Test Pattern (using Config Models):**
```go
func TestAcc_<ResourceName>_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	basic := model.<ResourceName>("test", id.DatabaseName(), id.SchemaName(), id.Name())
	complete := model.<ResourceName>("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.<ResourceName>),
		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      basic.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update - set optionals
			{
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}
```

**Patterns - Test Suites:**
- Use `t.Run()` for sub-tests and better organization
- Parallel execution via `t.Parallel()` when tests don't share state
- One-level nesting for simple unit tests, deeper for complex scenarios

**Setup/Teardown:**
- Setup: Initialize test data before `t.Run()` or in PreConfig
- Teardown: Use `t.Cleanup()` for resource cleanup
- Example: `t.Cleanup(apiIntegrationCleanup)`

**Test Data Lifecycle:**
```go
// Setup
apiIntegration, apiIntegrationCleanup := testClient().SecurityIntegration.CreateApiAuthenticationWithClientCredentialsFlowWithEnabled(t, true)
t.Cleanup(apiIntegrationCleanup)  // Teardown

// Test uses apiIntegration
```

## Mocking

**Framework:**
- No external mocking framework used
- Manual struct embedding with stubbed methods for unit tests
- Example from tasks_test.go:
```go
type testTasks struct {
	tasks
	stubbedTasks map[string]*Task
}

func (v *testTasks) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Task, error) {
	t, ok := v.stubbedTasks[id.Name()]
	if !ok {
		return nil, errors.New("no task configured, check test config")
	}
	return t, nil
}
```

**What to Mock:**
- External dependencies (SDK methods, Snowflake API calls)
- Complex objects that are expensive to create
- Methods with side effects you want to control

**What NOT to Mock:**
- SDK identifiers and helper functions (test with real instances)
- Schema definitions and validation functions
- Resource CRUD operation handlers (use full integration tests instead)

**Acceptance Test Fixtures:**
- Use `testClient()` helper for accessing test context
- Use `model.<ResourceName>()` builders for resource creation
- Use `datasourcemodel.<DataSourceName>()` builders for datasource creation
- All test data generated with prefixes to avoid conflicts

## Fixtures and Factories

**Test Data Builders:**
Located in `pkg/acceptance/bettertestspoc/config/model/` and generated via `make generate-resource-model-builders`

```go
basic := model.Secret("test", id.DatabaseName(), id.SchemaName(), id.Name(), secretString)
complete := model.Secret("test", id.DatabaseName(), id.SchemaName(), id.Name(), secretString).
	WithComment(comment)
```

**Builder Methods:**
- Fluent API: `WithComment()`, `WithDescription()`, `WithEnabled()`, etc.
- Always return self for chaining
- Factories handle generation of test resources

**Identifier Generation:**
```go
id := testClient().Ids.RandomSchemaObjectIdentifier()
id := testClient().Ids.RandomAccountObjectIdentifier()
idWithPrefix := testClient().Ids.RandomAccountObjectIdentifierWithPrefix("prefix")
```

**Random Data:**
```go
import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"

comment := random.Comment()
secretString := random.String()
alphaString := random.AlphaN(4)
```

**Location:**
- Builders: Generated in `pkg/acceptance/bettertestspoc/config/model/`
- Random helpers: `pkg/acceptance/helpers/random/`
- Test client: `pkg/testacc/helpers.go` (provides `testClient()`)

## Coverage

**Requirements:**
- No enforced minimum
- Coverage measured but not gated

**View Coverage:**
```bash
go test -cover ./pkg/resources
go test -cover ./pkg/sdk
```

**Coverage Tools:**
- Standard Go coverage: `go test -cover`
- Per-file coverage available via `-coverprofile=coverage.out`

## Test Types

**Unit Tests:**
- Scope: Individual functions and helper utilities
- Approach: Test with table-driven tests, verify business logic
- Location: `*_test.go` co-located with source
- Example files: `pkg/resources/common_test.go`, `pkg/helpers/helpers_test.go`, `pkg/sdk/tasks_test.go`

**Integration Tests (SDK):**
- Scope: SDK client methods with real Snowflake backend
- Approach: Full round-trip tests (create, read, update, delete)
- Location: `pkg/sdk/testint/`
- Requires: Snowflake test account and credentials
- Run: `make test-integration`

**Acceptance Tests (Resources/DataSources):**
- Scope: Full provider lifecycle (plan, apply, import, destroy)
- Approach: Multi-step tests with Terraform configurations
- Location: `pkg/testacc/resource_*_acceptance_test.go` or `pkg/testacc/data_source_*_acceptance_test.go`
- Requires: Terraform binary and Snowflake test account
- Run: `make test-acceptance`
- Resource tests must include: `BasicUseCase` and `CompleteUseCase` tests

**E2E Tests:**
- Not formally defined; acceptance tests serve as E2E verification
- Main use cases tested via `make test-main-terraform-use-cases`

## Common Patterns

**Async Testing:**
- Terraform Plugin Testing handles terraform execution asynchronously
- PreConfig/Check functions run before/after apply
- No explicit async/await needed

**Error Testing:**
```go
func Test_GetRootTasks_Error(t *testing.T) {
	_, err := GetRootTasks(client, ctx, invalidId)
	require.Error(t, err)
	require.ErrorContains(t, err, "expected error message")
}
```

**Parametric/Table-Driven Tests:**
```go
tests := []struct {
	name    string
	input   string
	expected string
}{
	{"case 1", "input1", "expected1"},
	{"case 2", "input2", "expected2"},
}

for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		result := Process(tt.input)
		assert.Equal(t, tt.expected, result)
	})
}
```

**Assertion Helpers (Acceptance Tests):**
Located in `pkg/acceptance/bettertestspoc/assert/`

```go
objectassert.Secret(t, id).
	HasName(id.Name()).
	HasDatabaseName(id.DatabaseName()).
	HasComment(comment)

resourceassert.SecretResource(t, basicModel.ResourceReference()).
	HasNameString(id.Name()).
	HasCommentString(comment)

resourceshowoutputassert.SecretShowOutput(t, model.ResourceReference()).
	HasCreatedOnNotEmpty().
	HasName(id.Name()).
	HasComment(comment)
```

**Assertion Patterns:**
- Object assertions verify Snowflake objects via DESCRIBE/SHOW
- Resource assertions verify terraform state
- Show output assertions verify computed show_output field
- Chained with `.HasX()` methods for fluent API
- Fail immediately on first assertion failure

## Acceptance Test Structure

**Mandatory Test Cases per Resource:**

1. **TestAcc_<ResourceName>_BasicUseCase:**
   - Create with required fields only
   - Import the resource
   - Update with optional fields set
   - Import with optionals
   - Update to unset optionals
   - Verify external changes are detected
   - Destroy and verify

2. **TestAcc_<ResourceName>_CompleteUseCase:**
   - Create with all fields set
   - Import with all fields

**Optional Test Cases:**
- Migration tests: `TestAcc_<ResourceName>_migrateFromVersion_<VERSION>`
- Specific feature tests: `TestAcc_<ResourceName>_<FeatureName>`
- Filtering tests: `TestAcc_<DataSourceName>_<FilterType>`

**State Upgraders:**
- Must be tested if implemented
- Named: `TestAcc_<ResourceName>_migrateFromVersion_2_13_0` (use latest version)
- Covers: Old state format → New state format migration

**External Changes:**
```go
// Detect external changes made via SDK/API
{
	PreConfig: func() {
		testClient().Secret.Alter(t, sdk.NewAlterSecretRequest(id).
			WithSet(*sdk.NewSecretSetRequest().WithComment(comment)))
	},
	ConfigPlanChecks: resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(model.ResourceReference(), plancheck.ResourceActionUpdate),
		},
	},
	Config: config.FromModels(t, model),
	Check:  assertThat(t, assertions...),
}
```

## Coverage Gaps

**Untested Areas:**
- Legacy resources using old patterns may have lower coverage
- Generated code from `make generate-*` targets is auto-generated and not individually tested
- Scripts in `scripts/` directory not covered

**Priority:** Medium
- Focus on new resource implementations
- Integration with acceptance tests provides good coverage for main paths
- Unit tests cover validation and helper logic

---

*Testing analysis: 2025-02-25*
