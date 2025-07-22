//go:build !account_level_tests

package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func TestInt_SafeDropOnAccountObjectIdentifier(t *testing.T) {
	networkPolicy, cleanupNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(cleanupNetworkPolicy)

	ctx := context.Background()
	networkPolicyDrop := func() error {
		return testClient(t).NetworkPolicies.Drop(ctx, sdk.NewDropNetworkPolicyRequest(networkPolicy.ID()).WithIfExists(true))
	}

	err := sdk.SafeDrop(testClient(t), networkPolicyDrop, ctx, networkPolicy.ID())
	assert.NoError(t, err)

	err = sdk.SafeDrop(testClient(t), networkPolicyDrop, ctx, networkPolicy.ID())
	assert.NoError(t, err)
}

func TestInt_SafeDropOnDatabaseObjectIdentifier(t *testing.T) {
	databaseRole, cleanupDatabaseRole := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(cleanupDatabaseRole)

	ctx := context.Background()
	databaseRoleDrop := func(id sdk.DatabaseObjectIdentifier) func() error {
		return func() error {
			return testClient(t).DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id).WithIfExists(true))
		}
	}

	err := sdk.SafeDrop(testClient(t), databaseRoleDrop(databaseRole.ID()), testContext(t), databaseRole.ID())
	assert.NoError(t, err)

	invalidDatabaseRoleId := NonExistingDatabaseObjectIdentifierWithNonExistingDatabase
	err = sdk.SafeDrop(testClient(t), databaseRoleDrop(invalidDatabaseRoleId), testContext(t), invalidDatabaseRoleId)
	assert.NoError(t, err)
}

func TestInt_SafeDropOnSchemaObjectIdentifier(t *testing.T) {
	table, cleanupTable := testClientHelper().Table.Create(t)
	t.Cleanup(cleanupTable)

	ctx := context.Background()
	tableDrop := func(id sdk.SchemaObjectIdentifier) func() error {
		return func() error {
			return testClient(t).Tables.Drop(ctx, sdk.NewDropTableRequest(id).WithIfExists(sdk.Bool(true)))
		}
	}

	err := sdk.SafeDrop(testClient(t), tableDrop(table.ID()), ctx, table.ID())
	assert.NoError(t, err)

	invalidTableIdOnValidDatabase := NonExistingSchemaObjectIdentifierWithNonExistingSchema
	err = sdk.SafeDrop(testClient(t), tableDrop(invalidTableIdOnValidDatabase), ctx, invalidTableIdOnValidDatabase)
	assert.NoError(t, err)

	invalidTableId := NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema
	err = sdk.SafeDrop(testClient(t), tableDrop(invalidTableId), ctx, invalidTableId)
	assert.NoError(t, err)
}

func TestInt_SafeDropOnSchemaObjectIdentifierWithArguments(t *testing.T) {
	procedure, procedureCleanup := testClientHelper().Procedure.Create(t, sdk.DataTypeInt)
	t.Cleanup(procedureCleanup)

	ctx := context.Background()
	procedureDrop := func(id sdk.SchemaObjectIdentifierWithArguments) func() error {
		return func() error {
			return testClient(t).Procedures.Drop(ctx, sdk.NewDropProcedureRequest(id).WithIfExists(true))
		}
	}

	err := sdk.SafeDrop(testClient(t), procedureDrop(procedure.ID()), ctx, procedure.ID())
	assert.NoError(t, err)

	invalidProcedureIdOnValidDatabase := NonExistingSchemaObjectIdentifierWithArgumentsWithNonExistingSchema
	err = sdk.SafeDrop(testClient(t), procedureDrop(invalidProcedureIdOnValidDatabase), ctx, invalidProcedureIdOnValidDatabase)
	assert.NoError(t, err)

	invalidProcedureId := NonExistingSchemaObjectIdentifierWithArgumentsWithNonExistingDatabaseAndSchema
	err = sdk.SafeDrop(testClient(t), procedureDrop(invalidProcedureId), ctx, invalidProcedureId)
	assert.NoError(t, err)
}

func TestInt_SafeRemoveProgrammaticAccessToken(t *testing.T) {
	user, cleanupUser := testClientHelper().User.CreateUser(t)
	t.Cleanup(cleanupUser)

	token, cleanupToken := testClientHelper().User.AddProgrammaticAccessToken(t, user.ID())
	t.Cleanup(cleanupToken)

	ctx := context.Background()
	err := sdk.SafeRemoveProgrammaticAccessToken(testClient(t), ctx, sdk.NewRemoveUserProgrammaticAccessTokenRequest(user.ID(), token.ID()))
	assert.NoError(t, err)

	err = sdk.SafeRemoveProgrammaticAccessToken(testClient(t), ctx, sdk.NewRemoveUserProgrammaticAccessTokenRequest(user.ID(), token.ID()))
	assert.NoError(t, err)

	invalidUserId := NonExistingAccountObjectIdentifier
	err = sdk.SafeRemoveProgrammaticAccessToken(testClient(t), ctx, sdk.NewRemoveUserProgrammaticAccessTokenRequest(invalidUserId, token.ID()))
	assert.NoError(t, err)
}
