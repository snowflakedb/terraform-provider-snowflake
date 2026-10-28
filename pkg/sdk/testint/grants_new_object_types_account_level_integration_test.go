//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// Snowflake Intelligence is an account-level object. The ALL grant cases create one,
// so this test runs with the account_level_tests build tag.
func TestInt_GrantPrivileges_OnFutureAndAll_UnsupportedObjectTypes_SnowflakeIntelligences(t *testing.T) {
	testGrantPrivilegesOnFutureAndAllUnsupportedObjectTypes(t, []unsupportedFutureAndAllGrantTestCase{
		{
			objectTypePlural: sdk.PluralObjectTypeSnowflakeIntelligences,
			privilege:        sdk.SchemaObjectPrivilegeUsage,
			createObject: func(t *testing.T) (sdk.ObjectIdentifier, func()) {
				t.Helper()
				return testClientHelper().SnowflakeIntelligence.Create(t)
			},
			expectedFutureError: "unexpected 'SNOWFLAKE'.",
			expectedAllError:    "unexpected 'INTELLIGENCES'.",
		},
	})
}
