//go:build non_account_level_tests

package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func schemaObjectShowByIDWrapper[T any](showByIdFn func(context.Context, sdk.SchemaObjectIdentifier) (*T, error)) func(context.Context, sdk.SchemaObjectIdentifier) error {
	return func(ctx context.Context, id sdk.SchemaObjectIdentifier) error {
		_, err := showByIdFn(ctx, id)
		return err
	}
}

func schemaObjectDropWrapper[T any](dropFn func(ctx context.Context, req T) error, req T) func(context.Context) error {
	return func(ctx context.Context) error {
		return dropFn(ctx, req)
	}
}

func schemaObjectWithArgumentsShowByIDWrapper[T any](showByIdFn func(context.Context, sdk.SchemaObjectIdentifierWithArguments) (*T, error)) func(context.Context, sdk.SchemaObjectIdentifierWithArguments) error {
	return func(ctx context.Context, id sdk.SchemaObjectIdentifierWithArguments) error {
		_, err := showByIdFn(ctx, id)
		return err
	}
}

func TestInt_ShowSchemaObjectInNonExistingDatabase(t *testing.T) {
	testCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr error
		ShowFn      func(context.Context, sdk.SchemaObjectIdentifier) error
	}{
		// Only object types that use IN SCHEMA in their ShowByID implementation
		{ObjectType: sdk.ObjectTypeTable, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tables.ShowByID)},
		{ObjectType: sdk.ObjectTypeDynamicTable, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).DynamicTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeCortexSearchService, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).CortexSearchServices.ShowByID)},
		{ObjectType: sdk.ObjectTypeExternalTable, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).ExternalTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeEventTable, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).EventTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeView, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Views.ShowByID)},
		{ObjectType: sdk.ObjectTypeMaterializedView, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).MaterializedViews.ShowByID)},
		{ObjectType: sdk.ObjectTypeSequence, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Sequences.ShowByID)},
		{ObjectType: sdk.ObjectTypeStream, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Streams.ShowByID)},
		{ObjectType: sdk.ObjectTypeTask, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tasks.ShowByID)},
		{ObjectType: sdk.ObjectTypeMaskingPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).MaskingPolicies.ShowByID)},
		{ObjectType: sdk.ObjectTypeRowAccessPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).RowAccessPolicies.ShowByID)},
		{ObjectType: sdk.ObjectTypeTag, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tags.ShowByID)},
		{ObjectType: sdk.ObjectTypeSecret, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Secrets.ShowByID)},
		{ObjectType: sdk.ObjectTypeStage, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Stages.ShowByID)},
		{ObjectType: sdk.ObjectTypeFileFormat, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).FileFormats.ShowByID)},
		{ObjectType: sdk.ObjectTypePipe, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Pipes.ShowByID)},
		{ObjectType: sdk.ObjectTypeAlert, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Alerts.ShowByID)},
		{ObjectType: sdk.ObjectTypeStreamlit, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Streamlits.ShowByID)},
		{ObjectType: sdk.ObjectTypeNetworkRule, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).NetworkRules.ShowByID)},
		{ObjectType: sdk.ObjectTypeAuthenticationPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).AuthenticationPolicies.ShowByID)},
		{ObjectType: sdk.ObjectTypeImageRepository, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).ImageRepositories.ShowByID)},
		{ObjectType: sdk.ObjectTypeService, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Services.ShowByID)},
		{ObjectType: sdk.ObjectTypeGitRepository, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).GitRepositories.ShowByID)},
		{ObjectType: sdk.ObjectTypeNotebook, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Notebooks.ShowByID)},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.ShowFn(ctx, sdk.NewSchemaObjectIdentifier("non-existing-database", "non-existing-schema", "non-existing-schema-object"))
			assert.ErrorIs(t, err, tt.ExpectedErr)
		})
	}

	schemaObjectWithArgumentsTestCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr error
		ShowFn      func(context.Context, sdk.SchemaObjectIdentifierWithArguments) error
	}{
		{ObjectType: sdk.ObjectTypeFunction, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).Functions.ShowByID)},
		{ObjectType: sdk.ObjectTypeExternalFunction, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).ExternalFunctions.ShowByID)},
		{ObjectType: sdk.ObjectTypeProcedure, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).Procedures.ShowByID)},
	}

	for _, tt := range schemaObjectWithArgumentsTestCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.ShowFn(ctx, sdk.NewSchemaObjectIdentifierWithArguments("non-existing-database", "non-existing-schema", "non-existing-schema-object"))
			assert.ErrorIs(t, err, tt.ExpectedErr)
		})
	}
}

func TestInt_ShowSchemaObjectInNonExistingSchema(t *testing.T) {
	testCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr error
		ShowFn      func(context.Context, sdk.SchemaObjectIdentifier) error
	}{
		// Only object types that use IN SCHEMA in their ShowByID implementation
		{ObjectType: sdk.ObjectTypeTable, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tables.ShowByID)},
		{ObjectType: sdk.ObjectTypeDynamicTable, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).DynamicTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeCortexSearchService, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).CortexSearchServices.ShowByID)},
		{ObjectType: sdk.ObjectTypeExternalTable, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).ExternalTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeEventTable, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).EventTables.ShowByID)},
		{ObjectType: sdk.ObjectTypeView, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Views.ShowByID)},
		{ObjectType: sdk.ObjectTypeMaterializedView, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).MaterializedViews.ShowByID)},
		{ObjectType: sdk.ObjectTypeSequence, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Sequences.ShowByID)},
		{ObjectType: sdk.ObjectTypeStream, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Streams.ShowByID)},
		{ObjectType: sdk.ObjectTypeTask, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tasks.ShowByID)},
		{ObjectType: sdk.ObjectTypeMaskingPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).MaskingPolicies.ShowByID)},
		{ObjectType: sdk.ObjectTypeRowAccessPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).RowAccessPolicies.ShowByID)},
		{ObjectType: sdk.ObjectTypeTag, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Tags.ShowByID)},
		{ObjectType: sdk.ObjectTypeSecret, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Secrets.ShowByID)},
		{ObjectType: sdk.ObjectTypeStage, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Stages.ShowByID)},
		{ObjectType: sdk.ObjectTypeFileFormat, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).FileFormats.ShowByID)},
		{ObjectType: sdk.ObjectTypePipe, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Pipes.ShowByID)},
		{ObjectType: sdk.ObjectTypeAlert, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Alerts.ShowByID)},
		{ObjectType: sdk.ObjectTypeStreamlit, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Streamlits.ShowByID)},
		{ObjectType: sdk.ObjectTypeNetworkRule, ExpectedErr: sdk.ErrDoesNotExistOrOperationCannotBePerformed, ShowFn: schemaObjectShowByIDWrapper(testClient(t).NetworkRules.ShowByID)},
		{ObjectType: sdk.ObjectTypeAuthenticationPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).AuthenticationPolicies.ShowByID)},
		{ObjectType: sdk.ObjectTypeImageRepository, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).ImageRepositories.ShowByID)},
		{ObjectType: sdk.ObjectTypeService, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Services.ShowByID)},
		{ObjectType: sdk.ObjectTypeGitRepository, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).GitRepositories.ShowByID)},
		{ObjectType: sdk.ObjectTypeNotebook, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectShowByIDWrapper(testClient(t).Notebooks.ShowByID)},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.ShowFn(ctx, sdk.NewSchemaObjectIdentifier(testClientHelper().Ids.DatabaseId().Name(), "non-existing-schema", "non-existing-schema-object"))
			assert.ErrorIs(t, err, tt.ExpectedErr)
		})
	}

	schemaObjectWithArgumentsTestCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr error
		ShowFn      func(context.Context, sdk.SchemaObjectIdentifierWithArguments) error
	}{
		{ObjectType: sdk.ObjectTypeFunction, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).Functions.ShowByID)},
		{ObjectType: sdk.ObjectTypeExternalFunction, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).ExternalFunctions.ShowByID)},
		{ObjectType: sdk.ObjectTypeProcedure, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, ShowFn: schemaObjectWithArgumentsShowByIDWrapper(testClient(t).Procedures.ShowByID)},
	}

	for _, tt := range schemaObjectWithArgumentsTestCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.ShowFn(ctx, sdk.NewSchemaObjectIdentifierWithArguments(testClientHelper().Ids.DatabaseId().Name(), "non-existing-schema", "non-existing-schema-object"))
			assert.ErrorIs(t, err, tt.ExpectedErr)
		})
	}
}

func TestInt_DropSchemaObjectInNonExistingDatabase(t *testing.T) {
	id := sdk.NewSchemaObjectIdentifier("non-existing-database", "non-existing-schema", "non-existing-schema-object")
	idWithArguments := sdk.NewSchemaObjectIdentifierWithArguments("non-existing-database", "non-existing-schema", "non-existing-schema-object", sdk.DataTypeInt)

	testCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr error
		DropFn      func(context.Context) error
	}{
		// Only object types that use IN SCHEMA in their ShowByID implementation
		{ObjectType: sdk.ObjectTypeTable, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tables.Drop, sdk.NewDropTableRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeDynamicTable, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).DynamicTables.Drop, sdk.NewDropDynamicTableRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeCortexSearchService, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).CortexSearchServices.Drop, sdk.NewDropCortexSearchServiceRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeExternalTable, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).ExternalTables.Drop, sdk.NewDropExternalTableRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeEventTable, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).EventTables.Drop, sdk.NewDropEventTableRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeView, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Views.Drop, sdk.NewDropViewRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeMaterializedView, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).MaterializedViews.Drop, sdk.NewDropMaterializedViewRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeSequence, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Sequences.Drop, sdk.NewDropSequenceRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeStream, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Streams.Drop, sdk.NewDropStreamRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeTask, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tasks.Drop, sdk.NewDropTaskRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeMaskingPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: func(ctx context.Context) error {
			return testClient(t).MaskingPolicies.Drop(ctx, id, &sdk.DropMaskingPolicyOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeRowAccessPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).RowAccessPolicies.Drop, sdk.NewDropRowAccessPolicyRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeTag, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tags.Drop, sdk.NewDropTagRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeSecret, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Secrets.Drop, sdk.NewDropSecretRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeStage, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Stages.Drop, sdk.NewDropStageRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeFileFormat, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: func(ctx context.Context) error {
			return testClient(t).FileFormats.Drop(ctx, id, &sdk.DropFileFormatOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypePipe, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: func(ctx context.Context) error {
			return testClient(t).Pipes.Drop(ctx, id, &sdk.DropPipeOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeAlert, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: func(ctx context.Context) error {
			return testClient(t).Alerts.Drop(ctx, id, &sdk.DropAlertOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeStreamlit, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Streamlits.Drop, sdk.NewDropStreamlitRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeNetworkRule, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).NetworkRules.Drop, sdk.NewDropNetworkRuleRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeAuthenticationPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).AuthenticationPolicies.Drop, sdk.NewDropAuthenticationPolicyRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeImageRepository, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).ImageRepositories.Drop, sdk.NewDropImageRepositoryRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeService, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Services.Drop, sdk.NewDropServiceRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeGitRepository, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).GitRepositories.Drop, sdk.NewDropGitRepositoryRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeNotebook, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Notebooks.Drop, sdk.NewDropNotebookRequest(id).WithIfExists(true))},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.DropFn(ctx)
			assert.ErrorIs(t, err, tt.ExpectedErr)
		})
	}

	schemaObjectWithArgumentsTestCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr error
		DropFn      func(ctx context.Context) error
	}{
		{ObjectType: sdk.ObjectTypeFunction, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Functions.Drop, sdk.NewDropFunctionRequest(idWithArguments).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeProcedure, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Procedures.Drop, sdk.NewDropProcedureRequest(idWithArguments).WithIfExists(true))},
	}

	for _, tt := range schemaObjectWithArgumentsTestCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.DropFn(ctx)
			assert.ErrorIs(t, err, tt.ExpectedErr)
		})
	}
}

func TestInt_DropSchemaInNonExistingDatabase(t *testing.T) {
	ctx := context.Background()
	err := testClient(t).Schemas.Drop(ctx, sdk.NewDatabaseObjectIdentifier("non-existing-database", "non-existing-schema"), &sdk.DropSchemaOptions{IfExists: sdk.Bool(true)})
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
}

func TestInt_DropSchemaObjectInNonExistingSchema(t *testing.T) {
	id := sdk.NewSchemaObjectIdentifier(testClientHelper().Ids.DatabaseId().Name(), "non-existing-schema", "non-existing-schema-object")
	idWithArguments := sdk.NewSchemaObjectIdentifierWithArguments(testClientHelper().Ids.DatabaseId().Name(), "non-existing-schema", "non-existing-schema-object", sdk.DataTypeInt)

	testCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr error
		DropFn      func(ctx context.Context) error
	}{
		// Only object types that use IN SCHEMA in their ShowByID implementation
		{ObjectType: sdk.ObjectTypeTable, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tables.Drop, sdk.NewDropTableRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeDynamicTable, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).DynamicTables.Drop, sdk.NewDropDynamicTableRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeCortexSearchService, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).CortexSearchServices.Drop, sdk.NewDropCortexSearchServiceRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeExternalTable, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).ExternalTables.Drop, sdk.NewDropExternalTableRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeEventTable, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).EventTables.Drop, sdk.NewDropEventTableRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeView, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Views.Drop, sdk.NewDropViewRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeMaterializedView, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).MaterializedViews.Drop, sdk.NewDropMaterializedViewRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeSequence, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Sequences.Drop, sdk.NewDropSequenceRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeStream, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Streams.Drop, sdk.NewDropStreamRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeTask, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tasks.Drop, sdk.NewDropTaskRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeMaskingPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: func(ctx context.Context) error {
			return testClient(t).MaskingPolicies.Drop(ctx, id, &sdk.DropMaskingPolicyOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeRowAccessPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).RowAccessPolicies.Drop, sdk.NewDropRowAccessPolicyRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeTag, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Tags.Drop, sdk.NewDropTagRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeSecret, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Secrets.Drop, sdk.NewDropSecretRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeStage, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Stages.Drop, sdk.NewDropStageRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeFileFormat, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: func(ctx context.Context) error {
			return testClient(t).FileFormats.Drop(ctx, id, &sdk.DropFileFormatOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypePipe, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: func(ctx context.Context) error {
			return testClient(t).Pipes.Drop(ctx, id, &sdk.DropPipeOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeAlert, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: func(ctx context.Context) error {
			return testClient(t).Alerts.Drop(ctx, id, &sdk.DropAlertOptions{IfExists: sdk.Bool(true)})
		}},
		{ObjectType: sdk.ObjectTypeStreamlit, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Streamlits.Drop, sdk.NewDropStreamlitRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeNetworkRule, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).NetworkRules.Drop, sdk.NewDropNetworkRuleRequest(id).WithIfExists(sdk.Bool(true)))},
		{ObjectType: sdk.ObjectTypeAuthenticationPolicy, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).AuthenticationPolicies.Drop, sdk.NewDropAuthenticationPolicyRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeImageRepository, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).ImageRepositories.Drop, sdk.NewDropImageRepositoryRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeService, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Services.Drop, sdk.NewDropServiceRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeGitRepository, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).GitRepositories.Drop, sdk.NewDropGitRepositoryRequest(id).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeNotebook, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Notebooks.Drop, sdk.NewDropNotebookRequest(id).WithIfExists(true))},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.DropFn(ctx)
			assert.ErrorIs(t, err, tt.ExpectedErr)
		})
	}

	schemaObjectWithArgumentsTestCases := []struct {
		ObjectType  sdk.ObjectType
		ExpectedErr error
		DropFn      func(ctx context.Context) error
	}{
		{ObjectType: sdk.ObjectTypeFunction, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Functions.Drop, sdk.NewDropFunctionRequest(idWithArguments).WithIfExists(true))},
		{ObjectType: sdk.ObjectTypeProcedure, ExpectedErr: sdk.ErrObjectNotExistOrAuthorized, DropFn: schemaObjectDropWrapper(testClient(t).Procedures.Drop, sdk.NewDropProcedureRequest(idWithArguments).WithIfExists(true))},
	}

	for _, tt := range schemaObjectWithArgumentsTestCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			ctx := context.Background()
			err := tt.DropFn(ctx)
			assert.ErrorIs(t, err, tt.ExpectedErr)
		})
	}
}

func TestInt_SnowflakeReturnsACustomErrorOnInvalidUserProgrammaticAccessToken(t *testing.T) {
	ctx := context.Background()

	user, userCleanup := testClientHelper().User.CreateUser(t)
	t.Cleanup(userCleanup)

	invalidUserId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	invalidTokenName := testClientHelper().Ids.RandomAccountObjectIdentifier()

	err := testClient(t).Users.RemoveProgrammaticAccessToken(ctx, sdk.NewRemoveUserProgrammaticAccessTokenRequest(user.ID(), invalidTokenName))
	assert.ErrorIs(t, err, sdk.ErrPatNotFound)
	assert.NotErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	assert.NotErrorIs(t, err, sdk.ErrObjectNotFound)

	err = testClient(t).Users.RemoveProgrammaticAccessToken(ctx, sdk.NewRemoveUserProgrammaticAccessTokenRequest(invalidUserId, invalidTokenName))
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
}
