//go:build !account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_CurrentAccount(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	account, err := client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, account)
}

func TestInt_CurrentAccountName(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	accountName, err := client.ContextFunctions.CurrentAccountName(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, accountName)
}

func TestInt_CurrentOrganizationName(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	organizationName, err := client.ContextFunctions.CurrentOrganizationName(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, organizationName)
}

func TestInt_CurrentRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	role, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, role.Name())
}

func TestInt_CurrentRegion(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	region, err := client.ContextFunctions.CurrentRegion(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, region)
}

func TestInt_CurrentSession(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	session, err := client.ContextFunctions.CurrentSession(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, session)
}

func TestInt_CurrentUser(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	user, err := client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, user.Name())
}

func TestInt_CurrentSessionDetails(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	account, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	require.NoError(t, err)
	assert.NotNil(t, account)
	assert.NotEmpty(t, account.Account)
	assert.NotEmpty(t, account.Role)
	assert.NotEmpty(t, account.Region)
	assert.NotEmpty(t, account.Session)
	assert.NotEmpty(t, account.User)
}

func TestInt_CurrentDatabase(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)
	err := client.Sessions.UseDatabase(ctx, databaseTest.ID())
	require.NoError(t, err)
	db, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, db)
}

func TestInt_CurrentSchema(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// new database and schema created on purpose
	databaseTest, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, databaseTest.ID())
	t.Cleanup(schemaCleanup)
	err := client.Sessions.UseSchema(ctx, schemaTest.ID())
	require.NoError(t, err)
	schema, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, schema)
}

func TestInt_CurrentWarehouse(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// new warehouse created on purpose
	warehouseTest, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	err := client.Sessions.UseWarehouse(ctx, warehouseTest.ID())
	require.NoError(t, err)
	warehouse, err := client.ContextFunctions.CurrentWarehouse(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, warehouse)
}

func TestInt_IsRoleInSession(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	currentRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)
	role, err := client.ContextFunctions.IsRoleInSession(ctx, currentRole)
	require.NoError(t, err)
	assert.True(t, role)
}

func TestInt_RolesUse(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	currentRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	role, cleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(cleanup)
	require.NotEqual(t, currentRole.Name(), role.Name)

	err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(role.ID(), sdk.GrantRole{Role: &currentRole}))
	require.NoError(t, err)

	err = client.Sessions.UseRole(ctx, role.ID())
	require.NoError(t, err)

	activeRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	assert.Equal(t, activeRole.Name(), role.Name)

	err = client.Sessions.UseRole(ctx, currentRole)
	require.NoError(t, err)
}

func TestInt_RolesUseSecondaryRoles(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	currentRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	role, cleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(cleanup)
	require.NotEqual(t, currentRole.Name(), role.Name)

	user, err := client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)

	err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(role.ID(), sdk.GrantRole{User: &user}))
	require.NoError(t, err)

	err = client.Sessions.UseRole(ctx, role.ID())
	require.NoError(t, err)

	err = client.Sessions.UseSecondaryRoles(ctx, sdk.SecondaryRolesAll)
	require.NoError(t, err)

	r, err := client.ContextFunctions.CurrentSecondaryRoles(ctx)
	require.NoError(t, err)

	names := make([]string, len(r.Roles))
	for i, v := range r.Roles {
		names[i] = v.Name()
	}
	assert.Equal(t, sdk.SecondaryRolesAll, r.Value)
	assert.Contains(t, names, currentRole.Name())

	err = client.Sessions.UseSecondaryRoles(ctx, sdk.SecondaryRolesNone)
	require.NoError(t, err)

	secondaryRolesAfter, err := client.ContextFunctions.CurrentSecondaryRoles(ctx)
	require.NoError(t, err)

	assert.Equal(t, sdk.SecondaryRolesNone, secondaryRolesAfter.Value)
	assert.Empty(t, secondaryRolesAfter.Roles)

	t.Cleanup(func() {
		err = client.Sessions.UseRole(ctx, currentRole)
		require.NoError(t, err)
	})
}

func TestInt_LastQueryId(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	lastQueryId, err := client.ContextFunctions.LastQueryId(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, lastQueryId)
}
