//go:build non_account_level_tests

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
	assert.NotEmpty(t, account.Value)
}

func TestInt_CurrentAccountName(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	accountName, err := client.ContextFunctions.CurrentAccountName(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, accountName.Value)
}

func TestInt_CurrentOrganizationName(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	organizationName, err := client.ContextFunctions.CurrentOrganizationName(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, organizationName.Value)
}

func TestInt_CurrentRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	role, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, role.Value.Name())
}

func TestInt_CurrentRegion(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	region, err := client.ContextFunctions.CurrentRegion(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, region.Value)
}

func TestInt_CurrentSession(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	session, err := client.ContextFunctions.CurrentSession(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, session.Value)
}

func TestInt_CurrentUser(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	user, err := client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, user.Value.Name())
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
	err := client.Sessions.UseDatabase(ctx, sdk.NewUseDatabaseSessionRequest(databaseTest.ID()))
	require.NoError(t, err)
	db, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, db.Value)
}

func TestInt_CurrentSchema(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// new database and schema created on purpose
	databaseTest, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, databaseTest.ID())
	t.Cleanup(schemaCleanup)
	err := client.Sessions.UseSchema(ctx, sdk.NewUseSchemaSessionRequest(schemaTest.ID()))
	require.NoError(t, err)
	schema, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, schema.Value)
}

func TestInt_CurrentWarehouse(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// new warehouse created on purpose
	warehouseTest, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	err := client.Sessions.UseWarehouse(ctx, sdk.NewUseWarehouseSessionRequest(warehouseTest.ID()))
	require.NoError(t, err)
	warehouse, err := client.ContextFunctions.CurrentWarehouse(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, warehouse.Value)
}

func TestInt_IsRoleInSession(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	currentRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)
	role, err := client.ContextFunctions.IsRoleInSession(ctx, sdk.NewIsRoleInSessionRequest(*sdk.NewIsRoleInSessionArgumentsRequest(currentRole.Value)))
	require.NoError(t, err)
	assert.True(t, role.Value)
}

func TestInt_RolesUse(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	currentRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	role, cleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(cleanup)
	require.NotEqual(t, currentRole.Value.Name(), role.Name)

	err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(role.ID(), *sdk.NewGrantRoleToRequest().WithRole(currentRole.Value)))
	require.NoError(t, err)

	err = client.Sessions.UseRole(ctx, sdk.NewUseRoleSessionRequest(role.ID()))
	require.NoError(t, err)

	activeRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	assert.Equal(t, activeRole.Value.Name(), role.Name)

	err = client.Sessions.UseRole(ctx, sdk.NewUseRoleSessionRequest(currentRole.Value))
	require.NoError(t, err)
}

func TestInt_RolesUseSecondaryRoles(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	currentRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	role, cleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(cleanup)
	require.NotEqual(t, currentRole.Value.Name(), role.Name)

	user, err := client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)

	err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(role.ID(), *sdk.NewGrantRoleToRequest().WithUser(user.Value)))
	require.NoError(t, err)

	err = client.Sessions.UseRole(ctx, sdk.NewUseRoleSessionRequest(role.ID()))
	require.NoError(t, err)

	err = client.Sessions.UseSecondaryRoles(ctx, sdk.NewUseSecondaryRolesSessionRequest(sdk.SecondaryRoleOptionAll))
	require.NoError(t, err)

	r, err := client.ContextFunctions.CurrentSecondaryRoles(ctx)
	require.NoError(t, err)

	names := make([]string, len(r.Roles))
	for i, v := range r.Roles {
		names[i] = v.Name()
	}
	assert.Equal(t, sdk.SecondaryRoleOptionAll, r.Value)
	assert.Contains(t, names, currentRole.Value.Name())

	err = client.Sessions.UseSecondaryRoles(ctx, sdk.NewUseSecondaryRolesSessionRequest(sdk.SecondaryRoleOptionNone))
	require.NoError(t, err)

	secondaryRolesAfter, err := client.ContextFunctions.CurrentSecondaryRoles(ctx)
	require.NoError(t, err)

	assert.Equal(t, sdk.SecondaryRoleOptionNone, secondaryRolesAfter.Value)
	assert.Empty(t, secondaryRolesAfter.Roles)

	t.Cleanup(func() {
		err = client.Sessions.UseRole(ctx, sdk.NewUseRoleSessionRequest(currentRole.Value))
		require.NoError(t, err)
	})
}

func TestInt_LastQueryId(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	lastQueryId, err := client.ContextFunctions.LastQueryId(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, lastQueryId.Value)
}
