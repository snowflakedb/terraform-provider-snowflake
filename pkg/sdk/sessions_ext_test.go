package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	sessionsTests.Alter.
		withModifyAndExpectedSqlf(
			case_Sessions_sql_Alter_Set,
			func(opts *AlterSessionOptions) {
				opts.Set = &SessionSet{SessionParameters: &SessionParameters{
					AbortDetachedQuery:    new(true),
					ClientPrefetchThreads: new(5),
				}}
			},
			`ALTER SESSION SET ABORT_DETACHED_QUERY = true, CLIENT_PREFETCH_THREADS = 5`,
		).
		withModifyAndExpectedSqlf(
			case_Sessions_sql_Alter_Unset,
			func(opts *AlterSessionOptions) {
				opts.Unset = &SessionUnset{SessionParametersUnset: &SessionParametersUnset{
					AbortDetachedQuery:    new(true),
					ClientPrefetchThreads: new(true),
				}}
			},
			`ALTER SESSION UNSET ABORT_DETACHED_QUERY, CLIENT_PREFETCH_THREADS`,
		)

	warehouseId := sessionsTestIdAccountObjectIdentifier
	databaseId := sessionsTestIdAccountObjectIdentifier
	schemaId := sessionsTestIdDatabaseObjectIdentifier
	roleId := sessionsTestIdAccountObjectIdentifier

	sessionsTests.UseWarehouse.
		withExpectedSqlf(case_Sessions_sql_UseWarehouse_basic, `USE WAREHOUSE %s`, warehouseId.FullyQualifiedName())
	sessionsTests.UseDatabase.
		withExpectedSqlf(case_Sessions_sql_UseDatabase_basic, `USE DATABASE %s`, databaseId.FullyQualifiedName())
	sessionsTests.UseSchema.
		withExpectedSqlf(case_Sessions_sql_UseSchema_basic, `USE SCHEMA %s`, schemaId.FullyQualifiedName())
	sessionsTests.UseRole.
		withExpectedSqlf(case_Sessions_sql_UseRole_basic, `USE ROLE %s`, roleId.FullyQualifiedName())
	sessionsTests.UseSecondaryRoles.
		withModifyAndExpectedSqlf(
			case_Sessions_sql_UseSecondaryRoles_basic,
			func(opts *UseSecondaryRolesSessionOptions) { opts.Option = SecondaryRoleOptionAll },
			`USE SECONDARY ROLES ALL`,
		)
}

func TestSessionParametersAdditionalValidations(t *testing.T) {
	err := (&SessionParameters{ClientPrefetchThreads: new(-1)}).additionalValidations()
	require.ErrorContains(t, err, "SessionParameters field: ClientPrefetchThreads must be greater than or equal to 0")
}
