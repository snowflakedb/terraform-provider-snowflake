//go:build !account_level_tests

package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func TestInt_SafeShowByIdOnAccountObjectIdentifier(t *testing.T) {
	networkPolicyShowById := func(ctx context.Context, id sdk.AccountObjectIdentifier) (*sdk.NetworkPolicy, error) {
		return testClientHelper().NetworkPolicy.Show(t, id)
	}

	networkPolicy, cleanupNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(cleanupNetworkPolicy)

	value, err := sdk.SafeShowById(testClient(t), networkPolicyShowById, testContext(t), networkPolicy.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupNetworkPolicy()

	_, err = sdk.SafeShowById(testClient(t), networkPolicyShowById, testContext(t), networkPolicy.ID())
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
}

func TestInt_SafeShowByIdOnDatabaseObjectIdentifier(t *testing.T) {
	databaseRoleShowById := func(ctx context.Context, id sdk.DatabaseObjectIdentifier) (*sdk.DatabaseRole, error) {
		return testClientHelper().DatabaseRole.Show(t, id)
	}

	databaseRole, cleanupDatabaseRole := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(cleanupDatabaseRole)

	value, err := sdk.SafeShowById(testClient(t), databaseRoleShowById, testContext(t), databaseRole.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupDatabaseRole()

	_, err = sdk.SafeShowById(testClient(t), databaseRoleShowById, testContext(t), databaseRole.ID())
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

	invalidDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	invalidDatabaseRoleId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(invalidDatabaseId)

	_, err = sdk.SafeShowById(testClient(t), databaseRoleShowById, testContext(t), invalidDatabaseRoleId)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
}

func TestInt_SafeShowByIdOnSchemaObjectIdentifier(t *testing.T) {
	tableShowById := func(ctx context.Context, id sdk.SchemaObjectIdentifier) (*sdk.Table, error) {
		return testClientHelper().Table.Show(t, id)
	}

	table, cleanupTable := testClientHelper().Table.Create(t)
	t.Cleanup(cleanupTable)

	value, err := sdk.SafeShowById(testClient(t), tableShowById, testContext(t), table.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupTable()

	_, err = sdk.SafeShowById(testClient(t), tableShowById, testContext(t), table.ID())
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

	invalidSchemaIdOnValidDatabase := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
	invalidTableIdOnValidDatabase := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(invalidSchemaIdOnValidDatabase)

	_, err = sdk.SafeShowById(testClient(t), tableShowById, testContext(t), invalidTableIdOnValidDatabase)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)

	invalidDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	invalidSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(invalidDatabaseId)
	invalidTableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(invalidSchemaId)

	_, err = sdk.SafeShowById(testClient(t), tableShowById, testContext(t), invalidTableId)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)
}

func TestInt_SafeShowByIdOnSchemaObjectIdentifierWithArguments(t *testing.T) {
	procedureShowById := func(ctx context.Context, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.Procedure, error) {
		return testClientHelper().Procedure.Show(t, id)
	}

	procedure, cleanupProcedure := testClientHelper().Procedure.Create(t, sdk.DataTypeInt)
	t.Cleanup(cleanupProcedure)

	value, err := sdk.SafeShowById(testClient(t), procedureShowById, testContext(t), procedure.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupProcedure()

	_, err = sdk.SafeShowById(testClient(t), procedureShowById, testContext(t), procedure.ID())
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

	invalidSchemaIdOnValidDatabase := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
	invalidProcedureIdOnValidDatabase := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchema(invalidSchemaIdOnValidDatabase)

	_, err = sdk.SafeShowById(testClient(t), procedureShowById, testContext(t), invalidProcedureIdOnValidDatabase)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

	invalidDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	invalidSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(invalidDatabaseId)
	invalidProcedureId := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchema(invalidSchemaId)

	_, err = sdk.SafeShowById(testClient(t), procedureShowById, testContext(t), invalidProcedureId)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
}
