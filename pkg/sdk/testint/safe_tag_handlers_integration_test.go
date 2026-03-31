//go:build non_account_level_tests

package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SafeUnsetTagFromTable(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	tags := []sdk.TagAssociation{
		{
			Name:  tag.ID(),
			Value: "abc",
		},
	}
	unsetTags := []sdk.ObjectIdentifier{tag.ID()}

	err := client.Tags.Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeTable, table.ID()).WithSetTags(tags))
	require.NoError(t, err)

	t.Run("unset existing tag", func(t *testing.T) {
		err := client.Tags.UnsetSafely(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeTable, table.ID()).WithUnsetTags(unsetTags))
		assert.NoError(t, err)
	})

	t.Run("unset already-unset tag", func(t *testing.T) {
		// Snowflake returns success when unsetting a tag that was never set.
		err := client.Tags.UnsetSafely(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeTable, table.ID()).WithUnsetTags(unsetTags))
		assert.NoError(t, err)
	})

	t.Run("non-existing tag on existing table", func(t *testing.T) {
		nonExistingTagIds := []sdk.ObjectIdentifier{NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema}
		err := client.Tags.UnsetSafely(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeTable, table.ID()).WithUnsetTags(nonExistingTagIds))
		assert.NoError(t, err)
	})

	t.Run("non-existing table", func(t *testing.T) {
		err := client.Tags.UnsetSafely(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeTable, NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema).WithUnsetTags(unsetTags))
		assert.NoError(t, err)
	})
}

func TestInt_SafeUnsetTagFromColumn(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	columnId := sdk.NewTableColumnIdentifier(table.ID().DatabaseName(), table.ID().SchemaName(), table.ID().Name(), "ID")

	tags := []sdk.TagAssociation{
		{
			Name:  tag.ID(),
			Value: "abc",
		},
	}
	unsetTags := []sdk.ObjectIdentifier{tag.ID()}

	err := client.Tags.Set(ctx, sdk.NewSetTagRequest(sdk.ObjectTypeColumn, columnId).WithSetTags(tags))
	require.NoError(t, err)

	t.Run("unset existing tag on column", func(t *testing.T) {
		err := client.Tags.UnsetSafely(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeColumn, columnId).WithUnsetTags(unsetTags))
		assert.NoError(t, err)
	})

	t.Run("non-existing tag on existing column", func(t *testing.T) {
		nonExistingTagIds := []sdk.ObjectIdentifier{NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema}
		err := client.Tags.UnsetSafely(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeColumn, columnId).WithUnsetTags(nonExistingTagIds))
		assert.NoError(t, err)
	})

	t.Run("non-existing column", func(t *testing.T) {
		nonExistingColumnId := sdk.NewTableColumnIdentifier("does_not_exist", "does_not_exist", "does_not_exist", "does_not_exist")
		err := client.Tags.UnsetSafely(ctx, sdk.NewUnsetTagRequest(sdk.ObjectTypeColumn, nonExistingColumnId).WithUnsetTags(unsetTags))
		assert.NoError(t, err)
	})
}

func TestInt_SafeUnsetTagOnNonExistingSchemaObject(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	unsetTags := []sdk.ObjectIdentifier{NonExistingSchemaObjectIdentifier}

	testCases := []struct {
		ObjectType sdk.ObjectType
	}{
		{ObjectType: sdk.ObjectTypeTable},
		{ObjectType: sdk.ObjectTypeExternalTable},
		{ObjectType: sdk.ObjectTypeEventTable},
		{ObjectType: sdk.ObjectTypeView},
		{ObjectType: sdk.ObjectTypeMaterializedView},
		{ObjectType: sdk.ObjectTypePipe},
		{ObjectType: sdk.ObjectTypeStage},
		{ObjectType: sdk.ObjectTypeStream},
		{ObjectType: sdk.ObjectTypeTask},
		{ObjectType: sdk.ObjectTypeAlert},
		{ObjectType: sdk.ObjectTypeMaskingPolicy},
		{ObjectType: sdk.ObjectTypeRowAccessPolicy},
		{ObjectType: sdk.ObjectTypeImageRepository},
		{ObjectType: sdk.ObjectTypeGitRepository},
		{ObjectType: sdk.ObjectTypeService},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			err := client.Tags.UnsetSafely(ctx, sdk.NewUnsetTagRequest(tt.ObjectType, NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema).WithUnsetTags(unsetTags))
			assert.NoError(t, err)
		})
	}
}

func TestInt_SafeUnsetTagOnNonExistingSchemaObjectWithArguments(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	unsetTags := []sdk.ObjectIdentifier{NonExistingSchemaObjectIdentifier}

	testCases := []struct {
		ObjectType sdk.ObjectType
	}{
		{ObjectType: sdk.ObjectTypeFunction},
		{ObjectType: sdk.ObjectTypeProcedure},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			err := client.Tags.UnsetSafely(ctx, sdk.NewUnsetTagRequest(tt.ObjectType, NonExistingSchemaObjectIdentifierWithArgumentsWithNonExistingDatabaseAndSchema).WithUnsetTags(unsetTags))
			assert.NoError(t, err)
		})
	}
}

func TestInt_SafeUnsetTagOnNonExistingAccountObject(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	unsetTags := []sdk.ObjectIdentifier{NonExistingSchemaObjectIdentifier}

	testCases := []struct {
		ObjectType sdk.ObjectType
	}{
		{ObjectType: sdk.ObjectTypeDatabase},
		{ObjectType: sdk.ObjectTypeWarehouse},
		{ObjectType: sdk.ObjectTypeComputePool},
		{ObjectType: sdk.ObjectTypeRole},
		{ObjectType: sdk.ObjectTypeUser},
		{ObjectType: sdk.ObjectTypeIntegration},
		{ObjectType: sdk.ObjectTypeNetworkPolicy},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			err := client.Tags.UnsetSafely(ctx, sdk.NewUnsetTagRequest(tt.ObjectType, NonExistingAccountObjectIdentifier).WithUnsetTags(unsetTags))
			assert.NoError(t, err)
		})
	}
}

func TestInt_SafeUnsetTagOnNonExistingDatabaseObject(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	unsetTags := []sdk.ObjectIdentifier{NonExistingSchemaObjectIdentifier}

	testCases := []struct {
		ObjectType sdk.ObjectType
	}{
		{ObjectType: sdk.ObjectTypeSchema},
		{ObjectType: sdk.ObjectTypeDatabaseRole},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			err := client.Tags.UnsetSafely(ctx, sdk.NewUnsetTagRequest(tt.ObjectType, NonExistingDatabaseObjectIdentifierWithNonExistingDatabase).WithUnsetTags(unsetTags))
			assert.NoError(t, err)
		})
	}
}
