//go:build non_account_level_tests

package testint

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_IdentifiersForOnePartIdentifierAsNameAndReference(t *testing.T) {
	identifier := func(prefix string) string {
		return testClientHelper().Ids.WithTestObjectSuffix(prefix)
	}

	identifierLowercase := func(prefix string) string {
		return strings.ToLower(identifier(prefix))
	}

	wrapInDoubleQuotes := func(text string) string {
		return `"` + text + `"`
	}

	testCases := []struct {
		Name     string
		ShowName string
		Error    string
	}{
		// special cases
		{Name: ``, Error: "invalid object identifier"},
		{Name: `"`, Error: "invalid object identifier"},
		// This is a valid identifier, but because in NewXIdentifier functions we're trimming double quotes it won't work
		{Name: `""`, Error: "invalid object identifier"},
		// This is a valid identifier, but because in NewXIdentifier functions we're trimming double quotes it won't work
		{Name: `""""`, Error: "invalid object identifier"},
		// This name is hardcoded on purpose, without test object suffix as we want to check such special case.
		{Name: `"."`, ShowName: `.`},

		// lower case
		{Name: identifierLowercase(`abc`), ShowName: identifierLowercase(`abc`)},
		{Name: identifierLowercase(`ab.c`), ShowName: identifierLowercase(`ab.c`)},
		{Name: identifierLowercase(`a"bc`), Error: `unexpected '"`},
		{Name: wrapInDoubleQuotes(identifierLowercase(`a""bc`)), ShowName: identifierLowercase(`a"bc`)},

		// upper case
		{Name: identifier(`ABC`), ShowName: identifier(`ABC`)},
		{Name: identifier(`AB.C`), ShowName: identifier(`AB.C`)},
		{Name: identifier(`A"BC`), Error: `unexpected '"`},
		{Name: wrapInDoubleQuotes(identifier(`A""BC`)), ShowName: identifier(`A"BC`)},

		// mixed case
		{Name: identifier(`AbC`), ShowName: identifier(`AbC`)},
		{Name: identifier(`Ab.C`), ShowName: identifier(`Ab.C`)},
		{Name: identifier(`A"bC`), Error: `unexpected '"`},
		{Name: wrapInDoubleQuotes(identifier(`A""bC`)), ShowName: identifier(`A"bC`)},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("one part identifier name and reference for input: %s", testCase.Name), func(t *testing.T) {
			ctx := context.Background()

			id := sdk.NewAccountObjectIdentifier(testCase.Name)
			err := testClient(t).ResourceMonitors.Create(ctx, id, new(sdk.CreateResourceMonitorOptions))
			if err == nil {
				t.Cleanup(testClientHelper().ResourceMonitor.DropResourceMonitorFunc(t, id))
			}
			if testCase.Error != "" {
				require.ErrorContains(t, err, testCase.Error)
			} else {
				require.NoError(t, err)
			}

			err = testClient(t).Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{
				ResourceMonitor: &id,
			})
			if err == nil {
				t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))
			}
			if testCase.Error != "" {
				require.ErrorContains(t, err, testCase.Error)
			} else {
				require.NoError(t, err)
				var result struct {
					Name            string `db:"name"`
					ResourceMonitor string `db:"resource_monitor"`
				}
				err = testClient(t).QueryOneForTests(ctx, &result, fmt.Sprintf("SHOW WAREHOUSES LIKE '%s'", testCase.ShowName))
				require.NoError(t, err)

				// For one part identifiers, we expect Snowflake to return unescaped identifiers (just like the ones we used for SHOW)
				assert.Equal(t, testCase.ShowName, result.Name)
				assert.Equal(t, testCase.ShowName, result.ResourceMonitor)
			}
		})
	}
}

func TestInt_IdentifiersForTwoPartIdentifierAsReference(t *testing.T) {
	identifier := func(prefix string) string {
		return testClientHelper().Ids.WithTestObjectSuffix(prefix)
	}

	identifierLowercase := func(prefix string) string {
		return strings.ToLower(identifier(prefix))
	}

	wrapInDoubleQuotes := func(text string) string {
		return `"` + text + `"`
	}

	type RawGrantOutput struct {
		Name      string `db:"name"`
		Privilege string `db:"privilege"`
	}

	testCases := []struct {
		Name                            string
		OverrideExpectedSnowflakeOutput string
		Error                           string
	}{
		// special cases
		{Name: ``, Error: "invalid object identifier"},
		{Name: `"`, Error: "invalid object identifier"},
		// This is a valid identifier, but because in NewXIdentifier functions we're trimming double quotes it won't work
		{Name: `""`, Error: "invalid object identifier"},
		// This is a valid identifier, but because in NewXIdentifier functions we're trimming double quotes it won't work
		{Name: `""""`, Error: "invalid object identifier"},
		// This name is hardcoded on purpose, without test object suffix as we want to check such special case.
		{Name: `"."`},

		// lower case
		{Name: identifierLowercase(`abc`)},
		{Name: identifierLowercase(`ab.c`)},
		{Name: identifierLowercase(`a"bc`), Error: `unexpected '"`},
		{Name: wrapInDoubleQuotes(identifierLowercase(`a""bc`))},

		// upper case
		{Name: identifier(`ABC`), OverrideExpectedSnowflakeOutput: identifier(`ABC`) + "." + identifier(`ABC`)},
		{Name: identifier(`AB.C`)},
		{Name: identifier(`A"BC`), Error: `unexpected '"`},
		{Name: wrapInDoubleQuotes(identifier(`A""BC`))},

		// mixed case
		{Name: identifier(`AbC`)},
		{Name: identifier(`Ab.C`)},
		{Name: identifier(`A"bC`), Error: `unexpected '"`},
		{Name: wrapInDoubleQuotes(identifier(`A""bC`))},
	}

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("two part identifier reference for input: %s", testCase.Name), func(t *testing.T) {
			ctx := context.Background()

			id := sdk.NewDatabaseObjectIdentifier(testCase.Name, testCase.Name)
			err := testClient(t).Databases.Create(ctx, id.DatabaseId(), new(sdk.CreateDatabaseOptions))
			if err == nil {
				t.Cleanup(testClientHelper().Database.DropDatabaseFunc(t, id.DatabaseId()))
			}
			if testCase.Error != "" {
				require.ErrorContains(t, err, testCase.Error)
			} else {
				require.NoError(t, err)
			}

			err = testClient(t).Schemas.Create(ctx, id, new(sdk.CreateSchemaOptions))
			if err == nil {
				t.Cleanup(testClientHelper().Schema.DropSchemaFunc(t, id))
			}
			if testCase.Error != "" {
				require.ErrorContains(t, err, testCase.Error)
			} else {
				require.NoError(t, err)

				testClientHelper().Grant.GrantOnSchemaToAccountRole(t, id, role.ID(), sdk.SchemaPrivilegeCreateTable)

				var grants []RawGrantOutput
				err = testClient(t).QueryForTests(ctx, &grants, fmt.Sprintf("SHOW GRANTS ON SCHEMA %s", id.FullyQualifiedName()))
				require.NoError(t, err)

				createTableGrant, err := collections.FindFirst(grants, func(output RawGrantOutput) bool { return output.Privilege == sdk.SchemaPrivilegeCreateTable.String() })
				require.NoError(t, err)

				// For two part identifiers, we expect Snowflake to return escaped identifiers with exception
				// to identifiers that don't have any lowercase character and special symbol in it.
				if testCase.OverrideExpectedSnowflakeOutput != "" {
					assert.Equal(t, testCase.OverrideExpectedSnowflakeOutput, createTableGrant.Name)
				} else {
					assert.Equal(t, id.FullyQualifiedName(), createTableGrant.Name)
				}
			}
		})
	}
}
