package customassert

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var _ assert.TestCheckFuncProvider = (*FutureGrantsAssert)(nil)

type FutureGrantsAssert struct {
	queryName  string
	roleId     sdk.AccountObjectIdentifier
	query      func(t *testing.T, testClient *helpers.TestClient) ([]sdk.Grant, error)
	assertions []func(t *testing.T, grants []sdk.Grant) error
}

func FutureGrantsInDatabaseToRole(t *testing.T, databaseId sdk.AccountObjectIdentifier, roleId sdk.AccountObjectIdentifier) *FutureGrantsAssert {
	t.Helper()
	return &FutureGrantsAssert{
		queryName: fmt.Sprintf("FUTURE_GRANTS_IN_DATABASE[%s]_TO_ROLE[%s]", databaseId.FullyQualifiedName(), roleId.FullyQualifiedName()),
		roleId:    roleId,
		query: func(t *testing.T, testClient *helpers.TestClient) ([]sdk.Grant, error) {
			t.Helper()
			return testClient.Grant.ShowFutureGrantsInDatabase(t, databaseId)
		},
		assertions: make([]func(t *testing.T, grants []sdk.Grant) error, 0),
	}
}

func FutureGrantsInSchemaToRole(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, roleId sdk.AccountObjectIdentifier) *FutureGrantsAssert {
	t.Helper()
	return &FutureGrantsAssert{
		queryName: fmt.Sprintf("FUTURE_GRANTS_IN_SCHEMA[%s]_TO_ROLE[%s]", schemaId.FullyQualifiedName(), roleId.FullyQualifiedName()),
		roleId:    roleId,
		query: func(t *testing.T, testClient *helpers.TestClient) ([]sdk.Grant, error) {
			t.Helper()
			return testClient.Grant.ShowFutureGrantsInSchema(t, schemaId)
		},
		assertions: make([]func(t *testing.T, grants []sdk.Grant) error, 0),
	}
}

func (a *FutureGrantsAssert) HasPrivilegesOnObjectTypeEqualTo(objectType sdk.ObjectType, expectedPrivileges ...string) *FutureGrantsAssert {
	a.assertions = append(a.assertions, func(t *testing.T, grants []sdk.Grant) error {
		t.Helper()

		filteredGrants := slices.DeleteFunc(grants, func(grant sdk.Grant) bool {
			return grant.GrantOn != objectType
		})

		actual := extractFuturePrivilegesForRole(filteredGrants, a.roleId)
		slices.Sort(actual)
		actual = slices.Compact(actual)

		expected := slices.Clone(expectedPrivileges)
		slices.Sort(expected)
		expected = slices.Compact(expected)

		if !slices.Equal(actual, expected) {
			return fmt.Errorf("expected future privileges: %v; got: %v", expected, actual)
		}

		return nil
	})
	return a
}

// ToTerraformTestCheckFunc TODO(SNOW-1501905): If possible, unify with runSnowflakeObjectsAssertions (SnowflakeObjectAssert)
func (a *FutureGrantsAssert) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		if testClient == nil {
			return fmt.Errorf("%s: testClient must not be nil", a.queryName)
		}

		grants, err := a.query(t, testClient)
		if err != nil {
			return fmt.Errorf("%s: query failed: %w", a.queryName, err)
		}

		var result []error
		for i, assertion := range a.assertions {
			if err := assertion(t, grants); err != nil {
				result = append(result, fmt.Errorf("%s assertion [%d/%d]: %w", a.queryName, i+1, len(a.assertions), err))
			}
		}

		return errors.Join(result...)
	}
}

func extractFuturePrivilegesForRole(grants []sdk.Grant, roleId sdk.AccountObjectIdentifier) []string {
	privileges := make([]string, 0)
	for _, grant := range grants {
		if grant.GrantTo != sdk.ObjectTypeRole || grant.GranteeName.Name() != roleId.Name() {
			continue
		}

		privileges = append(privileges, grant.Privilege)
	}

	return privileges
}
