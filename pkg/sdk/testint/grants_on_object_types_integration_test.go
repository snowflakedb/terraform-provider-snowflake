//go:build non_account_level_tests

package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests are the source of truth for which object types Snowflake accepts for each grant operation.
// The valid/invalid lists are validated by running against Snowflake and can later be adopted by grants_validations.go.

// accountObjectForType maps an ObjectType to the appropriate field in GrantOnAccountObject.
func accountObjectForType(objectType sdk.ObjectType, id sdk.AccountObjectIdentifier) *sdk.GrantOnAccountObject {
	obj := &sdk.GrantOnAccountObject{}
	switch objectType {
	case sdk.ObjectTypeUser:
		obj.User = &id
	case sdk.ObjectTypeResourceMonitor:
		obj.ResourceMonitor = &id
	case sdk.ObjectTypeWarehouse:
		obj.Warehouse = &id
	case sdk.ObjectTypeComputePool:
		obj.ComputePool = &id
	case sdk.ObjectTypeDatabase:
		obj.Database = &id
	case sdk.ObjectTypeIntegration:
		obj.Integration = &id
	case sdk.ObjectTypeFailoverGroup:
		obj.FailoverGroup = &id
	case sdk.ObjectTypeReplicationGroup:
		obj.ReplicationGroup = &id
	case sdk.ObjectTypeExternalVolume:
		obj.ExternalVolume = &id
	}
	return obj
}

// =====================================================
// Tests 1-2: Direct grant on schema object
// =====================================================

func TestInt_GrantOnSchemaObject_ToAccountRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	type test struct {
		objectType sdk.ObjectType
	}

	// Schema-level object types valid for GRANT ... ON <type> <name> TO ROLE.
	// Based on validGrantToSchemaObjectTypes in grants_validations.go.
	valid := []test{
		{objectType: sdk.ObjectTypeAgent},
		{objectType: sdk.ObjectTypeAggregationPolicy},
		{objectType: sdk.ObjectTypeAlert},
		{objectType: sdk.ObjectTypeAuthenticationPolicy},
		{objectType: sdk.ObjectTypeCortexSearchService},
		{objectType: sdk.ObjectTypeDataMetricFunction},
		{objectType: sdk.ObjectTypeDataset},
		{objectType: sdk.ObjectTypeDbtProject},
		{objectType: sdk.ObjectTypeDynamicTable},
		{objectType: sdk.ObjectTypeEventTable},
		{objectType: sdk.ObjectTypeExperiment},
		{objectType: sdk.ObjectTypeExternalTable},
		{objectType: sdk.ObjectTypeFileFormat},
		{objectType: sdk.ObjectTypeFunction},
		{objectType: sdk.ObjectTypeGateway},
		{objectType: sdk.ObjectTypeGitRepository},
		{objectType: sdk.ObjectTypeHybridTable},
		{objectType: sdk.ObjectTypeImageRepository},
		{objectType: sdk.ObjectTypeIcebergTable},
		{objectType: sdk.ObjectTypeJoinPolicy},
		{objectType: sdk.ObjectTypeMaskingPolicy},
		{objectType: sdk.ObjectTypeMaterializedView},
		{objectType: sdk.ObjectTypeMcpServer},
		{objectType: sdk.ObjectTypeModel},
		{objectType: sdk.ObjectTypeModelMonitor},
		{objectType: sdk.ObjectTypeNetworkRule},
		{objectType: sdk.ObjectTypeNotebook},
		{objectType: sdk.ObjectTypeNotebookProject},
		{objectType: sdk.ObjectTypeOnlineFeatureTable},
		{objectType: sdk.ObjectTypePackagesPolicy},
		{objectType: sdk.ObjectTypePasswordPolicy},
		{objectType: sdk.ObjectTypePipe},
		{objectType: sdk.ObjectTypePrivacyPolicy},
		{objectType: sdk.ObjectTypeProcedure},
		{objectType: sdk.ObjectTypeProjectionPolicy},
		{objectType: sdk.ObjectTypeRowAccessPolicy},
		{objectType: sdk.ObjectTypeSecret},
		{objectType: sdk.ObjectTypeSemanticView},
		{objectType: sdk.ObjectTypeService},
		{objectType: sdk.ObjectTypeSessionPolicy},
		{objectType: sdk.ObjectTypeSequence},
		{objectType: sdk.ObjectTypeSnapshot},
		{objectType: sdk.ObjectTypeSnapshotPolicy},
		{objectType: sdk.ObjectTypeSnapshotSet},
		{objectType: sdk.ObjectTypeStage},
		{objectType: sdk.ObjectTypeStorageLifecyclePolicy},
		{objectType: sdk.ObjectTypeStream},
		{objectType: sdk.ObjectTypeStreamlit},
		{objectType: sdk.ObjectTypeTable},
		{objectType: sdk.ObjectTypeTag},
		{objectType: sdk.ObjectTypeTask},
		{objectType: sdk.ObjectTypeView},
		{objectType: sdk.ObjectTypeWorkspace},
	}

	// Schema-level object types NOT valid for GRANT ... ON <type> <name> TO ROLE.
	// When a test for an "invalid" type starts returning a "not exist" error (instead of a type error),
	// it means Snowflake now supports it — move it to the valid list.
	invalid := []test{
		{objectType: sdk.ObjectTypeBudget},
		{objectType: sdk.ObjectTypeClassification},
		{objectType: sdk.ObjectTypeExternalFunction},
	}

	for _, tc := range valid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			err := client.Grants.GrantPrivilegesToAccountRole(
				ctx,
				&sdk.AccountRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)},
				&sdk.AccountRoleGrantOn{
					SchemaObject: &sdk.GrantOnSchemaObject{
						SchemaObject: &sdk.Object{
							ObjectType: tc.objectType,
							Name:       NonExistingSchemaObjectIdentifier,
						},
					},
				},
				role.ID(),
				nil,
			)
			assert.Error(t, err)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			err := client.Grants.GrantPrivilegesToAccountRole(
				ctx,
				&sdk.AccountRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)},
				&sdk.AccountRoleGrantOn{
					SchemaObject: &sdk.GrantOnSchemaObject{
						SchemaObject: &sdk.Object{
							ObjectType: tc.objectType,
							Name:       NonExistingSchemaObjectIdentifier,
						},
					},
				},
				role.ID(),
				nil,
			)
			assert.Error(t, err)
		})
	}
}

func TestInt_GrantOnSchemaObject_ToDatabaseRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	type test struct {
		objectType sdk.ObjectType
	}

	// Schema-level object types valid for GRANT ... ON <type> <name> TO DATABASE ROLE.
	// Based on validGrantToSchemaObjectTypes in grants_validations.go.
	valid := []test{
		{objectType: sdk.ObjectTypeAgent},
		{objectType: sdk.ObjectTypeAggregationPolicy},
		{objectType: sdk.ObjectTypeAlert},
		{objectType: sdk.ObjectTypeAuthenticationPolicy},
		{objectType: sdk.ObjectTypeCortexSearchService},
		{objectType: sdk.ObjectTypeDataMetricFunction},
		{objectType: sdk.ObjectTypeDataset},
		{objectType: sdk.ObjectTypeDbtProject},
		{objectType: sdk.ObjectTypeDynamicTable},
		{objectType: sdk.ObjectTypeEventTable},
		{objectType: sdk.ObjectTypeExperiment},
		{objectType: sdk.ObjectTypeExternalTable},
		{objectType: sdk.ObjectTypeFileFormat},
		{objectType: sdk.ObjectTypeFunction},
		{objectType: sdk.ObjectTypeGateway},
		{objectType: sdk.ObjectTypeGitRepository},
		{objectType: sdk.ObjectTypeHybridTable},
		{objectType: sdk.ObjectTypeImageRepository},
		{objectType: sdk.ObjectTypeIcebergTable},
		{objectType: sdk.ObjectTypeJoinPolicy},
		{objectType: sdk.ObjectTypeMaskingPolicy},
		{objectType: sdk.ObjectTypeMaterializedView},
		{objectType: sdk.ObjectTypeMcpServer},
		{objectType: sdk.ObjectTypeModel},
		{objectType: sdk.ObjectTypeModelMonitor},
		{objectType: sdk.ObjectTypeNetworkRule},
		{objectType: sdk.ObjectTypeNotebook},
		{objectType: sdk.ObjectTypeNotebookProject},
		{objectType: sdk.ObjectTypeOnlineFeatureTable},
		{objectType: sdk.ObjectTypePackagesPolicy},
		{objectType: sdk.ObjectTypePasswordPolicy},
		{objectType: sdk.ObjectTypePipe},
		{objectType: sdk.ObjectTypePrivacyPolicy},
		{objectType: sdk.ObjectTypeProcedure},
		{objectType: sdk.ObjectTypeProjectionPolicy},
		{objectType: sdk.ObjectTypeRowAccessPolicy},
		{objectType: sdk.ObjectTypeSecret},
		{objectType: sdk.ObjectTypeSemanticView},
		{objectType: sdk.ObjectTypeService},
		{objectType: sdk.ObjectTypeSessionPolicy},
		{objectType: sdk.ObjectTypeSequence},
		{objectType: sdk.ObjectTypeSnapshot},
		{objectType: sdk.ObjectTypeSnapshotPolicy},
		{objectType: sdk.ObjectTypeSnapshotSet},
		{objectType: sdk.ObjectTypeStage},
		{objectType: sdk.ObjectTypeStorageLifecyclePolicy},
		{objectType: sdk.ObjectTypeStream},
		{objectType: sdk.ObjectTypeStreamlit},
		{objectType: sdk.ObjectTypeTable},
		{objectType: sdk.ObjectTypeTag},
		{objectType: sdk.ObjectTypeTask},
		{objectType: sdk.ObjectTypeView},
		{objectType: sdk.ObjectTypeWorkspace},
	}

	// Schema-level object types NOT valid for GRANT ... ON <type> <name> TO DATABASE ROLE.
	invalid := []test{
		{objectType: sdk.ObjectTypeBudget},
		{objectType: sdk.ObjectTypeClassification},
		{objectType: sdk.ObjectTypeExternalFunction},
	}

	for _, tc := range valid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			err := client.Grants.GrantPrivilegesToDatabaseRole(
				ctx,
				&sdk.DatabaseRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)},
				&sdk.DatabaseRoleGrantOn{
					SchemaObject: &sdk.GrantOnSchemaObject{
						SchemaObject: &sdk.Object{
							ObjectType: tc.objectType,
							Name:       NonExistingSchemaObjectIdentifier,
						},
					},
				},
				databaseRole.ID(),
				nil,
			)
			assert.Error(t, err)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			err := client.Grants.GrantPrivilegesToDatabaseRole(
				ctx,
				&sdk.DatabaseRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)},
				&sdk.DatabaseRoleGrantOn{
					SchemaObject: &sdk.GrantOnSchemaObject{
						SchemaObject: &sdk.Object{
							ObjectType: tc.objectType,
							Name:       NonExistingSchemaObjectIdentifier,
						},
					},
				},
				databaseRole.ID(),
				nil,
			)
			assert.Error(t, err)
		})
	}
}

// =====================================================
// Tests 3-4: Grant on ALL in database
// =====================================================

func TestInt_GrantOnAllInDatabase_ToAccountRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	databaseId := testClientHelper().Ids.DatabaseId()

	type test struct {
		objectType sdk.ObjectType
	}

	// Schema-level object types valid for GRANT ... ON ALL <plural> IN DATABASE.
	// Based on validGrantToSchemaObjectTypes minus invalidGrantToAllObjectTypes in grants_validations.go.
	valid := []test{
		{objectType: sdk.ObjectTypeAgent},
		{objectType: sdk.ObjectTypeAlert},
		{objectType: sdk.ObjectTypeCortexSearchService},
		{objectType: sdk.ObjectTypeDataMetricFunction},
		{objectType: sdk.ObjectTypeDataset},
		{objectType: sdk.ObjectTypeDbtProject},
		{objectType: sdk.ObjectTypeDynamicTable},
		{objectType: sdk.ObjectTypeEventTable},
		{objectType: sdk.ObjectTypeExternalTable},
		{objectType: sdk.ObjectTypeFileFormat},
		{objectType: sdk.ObjectTypeFunction},
		{objectType: sdk.ObjectTypeGitRepository},
		{objectType: sdk.ObjectTypeImageRepository},
		{objectType: sdk.ObjectTypeIcebergTable},
		{objectType: sdk.ObjectTypeMaterializedView},
		{objectType: sdk.ObjectTypeMcpServer},
		{objectType: sdk.ObjectTypeModel},
		{objectType: sdk.ObjectTypeModelMonitor},
		{objectType: sdk.ObjectTypeNetworkRule},
		{objectType: sdk.ObjectTypeOnlineFeatureTable},
		{objectType: sdk.ObjectTypePipe},
		{objectType: sdk.ObjectTypeProcedure},
		{objectType: sdk.ObjectTypeSecret},
		{objectType: sdk.ObjectTypeSemanticView},
		{objectType: sdk.ObjectTypeService},
		{objectType: sdk.ObjectTypeSequence},
		{objectType: sdk.ObjectTypeSnapshot},
		{objectType: sdk.ObjectTypeSnapshotPolicy},
		{objectType: sdk.ObjectTypeStage},
		{objectType: sdk.ObjectTypeStream},
		{objectType: sdk.ObjectTypeStreamlit},
		{objectType: sdk.ObjectTypeTable},
		{objectType: sdk.ObjectTypeTask},
		{objectType: sdk.ObjectTypeView},
		{objectType: sdk.ObjectTypeWorkspace},
	}

	invalid := []test{
		{objectType: sdk.ObjectTypeAggregationPolicy},
		{objectType: sdk.ObjectTypeAuthenticationPolicy},
		{objectType: sdk.ObjectTypeBudget},
		{objectType: sdk.ObjectTypeClassification},
		{objectType: sdk.ObjectTypeExperiment},
		{objectType: sdk.ObjectTypeExternalFunction},
		{objectType: sdk.ObjectTypeGateway},
		{objectType: sdk.ObjectTypeHybridTable},
		{objectType: sdk.ObjectTypeJoinPolicy},
		{objectType: sdk.ObjectTypeMaskingPolicy},
		{objectType: sdk.ObjectTypeNotebook},
		{objectType: sdk.ObjectTypeNotebookProject},
		{objectType: sdk.ObjectTypePackagesPolicy},
		{objectType: sdk.ObjectTypePasswordPolicy},
		{objectType: sdk.ObjectTypePrivacyPolicy},
		{objectType: sdk.ObjectTypeProjectionPolicy},
		{objectType: sdk.ObjectTypeRowAccessPolicy},
		{objectType: sdk.ObjectTypeSessionPolicy},
		{objectType: sdk.ObjectTypeSnapshotSet},
		{objectType: sdk.ObjectTypeStorageLifecyclePolicy},
		{objectType: sdk.ObjectTypeTag},
	}

	for _, tc := range valid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.AccountRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InDatabase:       sdk.Pointer(databaseId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = client.Grants.RevokePrivilegesFromAccountRole(context.Background(), privileges, on, role.ID(), nil)
			})
		})
	}

	for _, tc := range invalid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.AccountRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InDatabase:       sdk.Pointer(databaseId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
			assert.Error(t, err)
		})
	}
}

func TestInt_GrantOnAllInDatabase_ToDatabaseRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)
	databaseId := testClientHelper().Ids.DatabaseId()

	type test struct {
		objectType sdk.ObjectType
	}

	valid := []test{
		{objectType: sdk.ObjectTypeAgent},
		{objectType: sdk.ObjectTypeAlert},
		{objectType: sdk.ObjectTypeCortexSearchService},
		{objectType: sdk.ObjectTypeDataMetricFunction},
		{objectType: sdk.ObjectTypeDataset},
		{objectType: sdk.ObjectTypeDbtProject},
		{objectType: sdk.ObjectTypeDynamicTable},
		{objectType: sdk.ObjectTypeEventTable},
		{objectType: sdk.ObjectTypeExternalTable},
		{objectType: sdk.ObjectTypeFileFormat},
		{objectType: sdk.ObjectTypeFunction},
		{objectType: sdk.ObjectTypeGitRepository},
		{objectType: sdk.ObjectTypeImageRepository},
		{objectType: sdk.ObjectTypeIcebergTable},
		{objectType: sdk.ObjectTypeMaterializedView},
		{objectType: sdk.ObjectTypeMcpServer},
		{objectType: sdk.ObjectTypeModel},
		{objectType: sdk.ObjectTypeModelMonitor},
		{objectType: sdk.ObjectTypeNetworkRule},
		{objectType: sdk.ObjectTypeOnlineFeatureTable},
		{objectType: sdk.ObjectTypePipe},
		{objectType: sdk.ObjectTypeProcedure},
		{objectType: sdk.ObjectTypeSecret},
		{objectType: sdk.ObjectTypeSemanticView},
		{objectType: sdk.ObjectTypeService},
		{objectType: sdk.ObjectTypeSequence},
		{objectType: sdk.ObjectTypeSnapshot},
		{objectType: sdk.ObjectTypeSnapshotPolicy},
		{objectType: sdk.ObjectTypeStage},
		{objectType: sdk.ObjectTypeStream},
		{objectType: sdk.ObjectTypeStreamlit},
		{objectType: sdk.ObjectTypeTable},
		{objectType: sdk.ObjectTypeTask},
		{objectType: sdk.ObjectTypeView},
		{objectType: sdk.ObjectTypeWorkspace},
	}

	invalid := []test{
		{objectType: sdk.ObjectTypeAggregationPolicy},
		{objectType: sdk.ObjectTypeAuthenticationPolicy},
		{objectType: sdk.ObjectTypeBudget},
		{objectType: sdk.ObjectTypeClassification},
		{objectType: sdk.ObjectTypeExperiment},
		{objectType: sdk.ObjectTypeExternalFunction},
		{objectType: sdk.ObjectTypeGateway},
		{objectType: sdk.ObjectTypeHybridTable},
		{objectType: sdk.ObjectTypeJoinPolicy},
		{objectType: sdk.ObjectTypeMaskingPolicy},
		{objectType: sdk.ObjectTypeNotebook},
		{objectType: sdk.ObjectTypeNotebookProject},
		{objectType: sdk.ObjectTypePackagesPolicy},
		{objectType: sdk.ObjectTypePasswordPolicy},
		{objectType: sdk.ObjectTypePrivacyPolicy},
		{objectType: sdk.ObjectTypeProjectionPolicy},
		{objectType: sdk.ObjectTypeRowAccessPolicy},
		{objectType: sdk.ObjectTypeSessionPolicy},
		{objectType: sdk.ObjectTypeSnapshotSet},
		{objectType: sdk.ObjectTypeStorageLifecyclePolicy},
		{objectType: sdk.ObjectTypeTag},
	}

	for _, tc := range valid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.DatabaseRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InDatabase:       sdk.Pointer(databaseId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRole.ID(), nil)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = client.Grants.RevokePrivilegesFromDatabaseRole(context.Background(), privileges, on, databaseRole.ID(), nil)
			})
		})
	}

	for _, tc := range invalid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.DatabaseRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InDatabase:       sdk.Pointer(databaseId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRole.ID(), nil)
			assert.Error(t, err)
		})
	}
}

// =====================================================
// Tests 5-6: Grant on ALL in schema
// =====================================================

func TestInt_GrantOnAllInSchema_ToAccountRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	schemaId := testClientHelper().Ids.SchemaId()

	type test struct {
		objectType sdk.ObjectType
	}

	valid := []test{
		{objectType: sdk.ObjectTypeAgent},
		{objectType: sdk.ObjectTypeAlert},
		{objectType: sdk.ObjectTypeCortexSearchService},
		{objectType: sdk.ObjectTypeDataMetricFunction},
		{objectType: sdk.ObjectTypeDataset},
		{objectType: sdk.ObjectTypeDbtProject},
		{objectType: sdk.ObjectTypeDynamicTable},
		{objectType: sdk.ObjectTypeEventTable},
		{objectType: sdk.ObjectTypeExternalTable},
		{objectType: sdk.ObjectTypeFileFormat},
		{objectType: sdk.ObjectTypeFunction},
		{objectType: sdk.ObjectTypeGitRepository},
		{objectType: sdk.ObjectTypeImageRepository},
		{objectType: sdk.ObjectTypeIcebergTable},
		{objectType: sdk.ObjectTypeMaterializedView},
		{objectType: sdk.ObjectTypeMcpServer},
		{objectType: sdk.ObjectTypeModel},
		{objectType: sdk.ObjectTypeModelMonitor},
		{objectType: sdk.ObjectTypeNetworkRule},
		{objectType: sdk.ObjectTypeOnlineFeatureTable},
		{objectType: sdk.ObjectTypePipe},
		{objectType: sdk.ObjectTypeProcedure},
		{objectType: sdk.ObjectTypeSecret},
		{objectType: sdk.ObjectTypeSemanticView},
		{objectType: sdk.ObjectTypeService},
		{objectType: sdk.ObjectTypeSequence},
		{objectType: sdk.ObjectTypeSnapshot},
		{objectType: sdk.ObjectTypeSnapshotPolicy},
		{objectType: sdk.ObjectTypeStage},
		{objectType: sdk.ObjectTypeStream},
		{objectType: sdk.ObjectTypeStreamlit},
		{objectType: sdk.ObjectTypeTable},
		{objectType: sdk.ObjectTypeTask},
		{objectType: sdk.ObjectTypeView},
		{objectType: sdk.ObjectTypeWorkspace},
	}

	invalid := []test{
		{objectType: sdk.ObjectTypeAggregationPolicy},
		{objectType: sdk.ObjectTypeAuthenticationPolicy},
		{objectType: sdk.ObjectTypeBudget},
		{objectType: sdk.ObjectTypeClassification},
		{objectType: sdk.ObjectTypeExperiment},
		{objectType: sdk.ObjectTypeExternalFunction},
		{objectType: sdk.ObjectTypeGateway},
		{objectType: sdk.ObjectTypeHybridTable},
		{objectType: sdk.ObjectTypeJoinPolicy},
		{objectType: sdk.ObjectTypeMaskingPolicy},
		{objectType: sdk.ObjectTypeNotebook},
		{objectType: sdk.ObjectTypeNotebookProject},
		{objectType: sdk.ObjectTypePackagesPolicy},
		{objectType: sdk.ObjectTypePasswordPolicy},
		{objectType: sdk.ObjectTypePrivacyPolicy},
		{objectType: sdk.ObjectTypeProjectionPolicy},
		{objectType: sdk.ObjectTypeRowAccessPolicy},
		{objectType: sdk.ObjectTypeSessionPolicy},
		{objectType: sdk.ObjectTypeSnapshotSet},
		{objectType: sdk.ObjectTypeStorageLifecyclePolicy},
		{objectType: sdk.ObjectTypeTag},
	}

	for _, tc := range valid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.AccountRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InSchema:         sdk.Pointer(schemaId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = client.Grants.RevokePrivilegesFromAccountRole(context.Background(), privileges, on, role.ID(), nil)
			})
		})
	}

	for _, tc := range invalid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.AccountRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InSchema:         sdk.Pointer(schemaId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
			assert.Error(t, err)
		})
	}
}

func TestInt_GrantOnAllInSchema_ToDatabaseRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)
	schemaId := testClientHelper().Ids.SchemaId()

	type test struct {
		objectType sdk.ObjectType
	}

	valid := []test{
		{objectType: sdk.ObjectTypeAgent},
		{objectType: sdk.ObjectTypeAlert},
		{objectType: sdk.ObjectTypeCortexSearchService},
		{objectType: sdk.ObjectTypeDataMetricFunction},
		{objectType: sdk.ObjectTypeDataset},
		{objectType: sdk.ObjectTypeDbtProject},
		{objectType: sdk.ObjectTypeDynamicTable},
		{objectType: sdk.ObjectTypeEventTable},
		{objectType: sdk.ObjectTypeExternalTable},
		{objectType: sdk.ObjectTypeFileFormat},
		{objectType: sdk.ObjectTypeFunction},
		{objectType: sdk.ObjectTypeGitRepository},
		{objectType: sdk.ObjectTypeImageRepository},
		{objectType: sdk.ObjectTypeIcebergTable},
		{objectType: sdk.ObjectTypeMaterializedView},
		{objectType: sdk.ObjectTypeMcpServer},
		{objectType: sdk.ObjectTypeModel},
		{objectType: sdk.ObjectTypeModelMonitor},
		{objectType: sdk.ObjectTypeNetworkRule},
		{objectType: sdk.ObjectTypeOnlineFeatureTable},
		{objectType: sdk.ObjectTypePipe},
		{objectType: sdk.ObjectTypeProcedure},
		{objectType: sdk.ObjectTypeSecret},
		{objectType: sdk.ObjectTypeSemanticView},
		{objectType: sdk.ObjectTypeService},
		{objectType: sdk.ObjectTypeSequence},
		{objectType: sdk.ObjectTypeSnapshot},
		{objectType: sdk.ObjectTypeSnapshotPolicy},
		{objectType: sdk.ObjectTypeStage},
		{objectType: sdk.ObjectTypeStream},
		{objectType: sdk.ObjectTypeStreamlit},
		{objectType: sdk.ObjectTypeTable},
		{objectType: sdk.ObjectTypeTask},
		{objectType: sdk.ObjectTypeView},
		{objectType: sdk.ObjectTypeWorkspace},
	}

	invalid := []test{
		{objectType: sdk.ObjectTypeAggregationPolicy},
		{objectType: sdk.ObjectTypeAuthenticationPolicy},
		{objectType: sdk.ObjectTypeBudget},
		{objectType: sdk.ObjectTypeClassification},
		{objectType: sdk.ObjectTypeExperiment},
		{objectType: sdk.ObjectTypeExternalFunction},
		{objectType: sdk.ObjectTypeGateway},
		{objectType: sdk.ObjectTypeHybridTable},
		{objectType: sdk.ObjectTypeJoinPolicy},
		{objectType: sdk.ObjectTypeMaskingPolicy},
		{objectType: sdk.ObjectTypeNotebook},
		{objectType: sdk.ObjectTypeNotebookProject},
		{objectType: sdk.ObjectTypePackagesPolicy},
		{objectType: sdk.ObjectTypePasswordPolicy},
		{objectType: sdk.ObjectTypePrivacyPolicy},
		{objectType: sdk.ObjectTypeProjectionPolicy},
		{objectType: sdk.ObjectTypeRowAccessPolicy},
		{objectType: sdk.ObjectTypeSessionPolicy},
		{objectType: sdk.ObjectTypeSnapshotSet},
		{objectType: sdk.ObjectTypeStorageLifecyclePolicy},
		{objectType: sdk.ObjectTypeTag},
	}

	for _, tc := range valid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.DatabaseRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InSchema:         sdk.Pointer(schemaId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRole.ID(), nil)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = client.Grants.RevokePrivilegesFromDatabaseRole(context.Background(), privileges, on, databaseRole.ID(), nil)
			})
		})
	}

	for _, tc := range invalid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.DatabaseRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InSchema:         sdk.Pointer(schemaId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRole.ID(), nil)
			assert.Error(t, err)
		})
	}
}

// =====================================================
// Tests 7-8: Grant on FUTURE in database
// =====================================================

func TestInt_GrantOnFutureInDatabase_ToAccountRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	databaseId := testClientHelper().Ids.DatabaseId()

	type test struct {
		objectType sdk.ObjectType
	}

	// Schema-level object types valid for GRANT ... ON FUTURE <plural> IN DATABASE.
	// Based on validGrantToSchemaObjectTypes minus invalidGrantToFutureObjectTypes in grants_validations.go.
	valid := []test{
		{objectType: sdk.ObjectTypeAgent},
		{objectType: sdk.ObjectTypeAlert},
		{objectType: sdk.ObjectTypeCortexSearchService},
		{objectType: sdk.ObjectTypeDataMetricFunction},
		{objectType: sdk.ObjectTypeDataset},
		{objectType: sdk.ObjectTypeDbtProject},
		{objectType: sdk.ObjectTypeDynamicTable},
		{objectType: sdk.ObjectTypeEventTable},
		{objectType: sdk.ObjectTypeExternalTable},
		{objectType: sdk.ObjectTypeFileFormat},
		{objectType: sdk.ObjectTypeFunction},
		{objectType: sdk.ObjectTypeGitRepository},
		{objectType: sdk.ObjectTypeImageRepository},
		{objectType: sdk.ObjectTypeIcebergTable},
		{objectType: sdk.ObjectTypeMaterializedView},
		{objectType: sdk.ObjectTypeMcpServer},
		{objectType: sdk.ObjectTypeModel},
		{objectType: sdk.ObjectTypeModelMonitor},
		{objectType: sdk.ObjectTypeNetworkRule},
		{objectType: sdk.ObjectTypeOnlineFeatureTable},
		{objectType: sdk.ObjectTypePipe},
		{objectType: sdk.ObjectTypeProcedure},
		{objectType: sdk.ObjectTypeSecret},
		{objectType: sdk.ObjectTypeSemanticView},
		{objectType: sdk.ObjectTypeService},
		{objectType: sdk.ObjectTypeSequence},
		{objectType: sdk.ObjectTypeStage},
		{objectType: sdk.ObjectTypeStream},
		{objectType: sdk.ObjectTypeStreamlit},
		{objectType: sdk.ObjectTypeTable},
		{objectType: sdk.ObjectTypeTask},
		{objectType: sdk.ObjectTypeView},
		{objectType: sdk.ObjectTypeWorkspace},
	}

	invalid := []test{
		{objectType: sdk.ObjectTypeAggregationPolicy},
		{objectType: sdk.ObjectTypeAuthenticationPolicy},
		{objectType: sdk.ObjectTypeBudget},
		{objectType: sdk.ObjectTypeClassification},
		{objectType: sdk.ObjectTypeExperiment},
		{objectType: sdk.ObjectTypeExternalFunction},
		{objectType: sdk.ObjectTypeGateway},
		{objectType: sdk.ObjectTypeHybridTable},
		{objectType: sdk.ObjectTypeJoinPolicy},
		{objectType: sdk.ObjectTypeMaskingPolicy},
		{objectType: sdk.ObjectTypeNotebook},
		{objectType: sdk.ObjectTypeNotebookProject},
		{objectType: sdk.ObjectTypePackagesPolicy},
		{objectType: sdk.ObjectTypePasswordPolicy},
		{objectType: sdk.ObjectTypePrivacyPolicy},
		{objectType: sdk.ObjectTypeProjectionPolicy},
		{objectType: sdk.ObjectTypeRowAccessPolicy},
		{objectType: sdk.ObjectTypeSessionPolicy},
		{objectType: sdk.ObjectTypeSnapshot},
		{objectType: sdk.ObjectTypeSnapshotPolicy},
		{objectType: sdk.ObjectTypeSnapshotSet},
		{objectType: sdk.ObjectTypeStorageLifecyclePolicy},
		{objectType: sdk.ObjectTypeTag},
	}

	for _, tc := range valid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.AccountRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					Future: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InDatabase:       sdk.Pointer(databaseId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = client.Grants.RevokePrivilegesFromAccountRole(context.Background(), privileges, on, role.ID(), nil)
			})
		})
	}

	for _, tc := range invalid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.AccountRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					Future: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InDatabase:       sdk.Pointer(databaseId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
			assert.Error(t, err)
		})
	}
}

func TestInt_GrantOnFutureInDatabase_ToDatabaseRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)
	databaseId := testClientHelper().Ids.DatabaseId()

	type test struct {
		objectType sdk.ObjectType
	}

	valid := []test{
		{objectType: sdk.ObjectTypeAgent},
		{objectType: sdk.ObjectTypeAlert},
		{objectType: sdk.ObjectTypeCortexSearchService},
		{objectType: sdk.ObjectTypeDataMetricFunction},
		{objectType: sdk.ObjectTypeDataset},
		{objectType: sdk.ObjectTypeDbtProject},
		{objectType: sdk.ObjectTypeDynamicTable},
		{objectType: sdk.ObjectTypeEventTable},
		{objectType: sdk.ObjectTypeExternalTable},
		{objectType: sdk.ObjectTypeFileFormat},
		{objectType: sdk.ObjectTypeFunction},
		{objectType: sdk.ObjectTypeGitRepository},
		{objectType: sdk.ObjectTypeImageRepository},
		{objectType: sdk.ObjectTypeIcebergTable},
		{objectType: sdk.ObjectTypeMaterializedView},
		{objectType: sdk.ObjectTypeMcpServer},
		{objectType: sdk.ObjectTypeModel},
		{objectType: sdk.ObjectTypeModelMonitor},
		{objectType: sdk.ObjectTypeNetworkRule},
		{objectType: sdk.ObjectTypeOnlineFeatureTable},
		{objectType: sdk.ObjectTypePipe},
		{objectType: sdk.ObjectTypeProcedure},
		{objectType: sdk.ObjectTypeSecret},
		{objectType: sdk.ObjectTypeSemanticView},
		{objectType: sdk.ObjectTypeService},
		{objectType: sdk.ObjectTypeSequence},
		{objectType: sdk.ObjectTypeStage},
		{objectType: sdk.ObjectTypeStream},
		{objectType: sdk.ObjectTypeStreamlit},
		{objectType: sdk.ObjectTypeTable},
		{objectType: sdk.ObjectTypeTask},
		{objectType: sdk.ObjectTypeView},
		{objectType: sdk.ObjectTypeWorkspace},
	}

	invalid := []test{
		{objectType: sdk.ObjectTypeAggregationPolicy},
		{objectType: sdk.ObjectTypeAuthenticationPolicy},
		{objectType: sdk.ObjectTypeBudget},
		{objectType: sdk.ObjectTypeClassification},
		{objectType: sdk.ObjectTypeExperiment},
		{objectType: sdk.ObjectTypeExternalFunction},
		{objectType: sdk.ObjectTypeGateway},
		{objectType: sdk.ObjectTypeHybridTable},
		{objectType: sdk.ObjectTypeJoinPolicy},
		{objectType: sdk.ObjectTypeMaskingPolicy},
		{objectType: sdk.ObjectTypeNotebook},
		{objectType: sdk.ObjectTypeNotebookProject},
		{objectType: sdk.ObjectTypePackagesPolicy},
		{objectType: sdk.ObjectTypePasswordPolicy},
		{objectType: sdk.ObjectTypePrivacyPolicy},
		{objectType: sdk.ObjectTypeProjectionPolicy},
		{objectType: sdk.ObjectTypeRowAccessPolicy},
		{objectType: sdk.ObjectTypeSessionPolicy},
		{objectType: sdk.ObjectTypeSnapshot},
		{objectType: sdk.ObjectTypeSnapshotPolicy},
		{objectType: sdk.ObjectTypeSnapshotSet},
		{objectType: sdk.ObjectTypeStorageLifecyclePolicy},
		{objectType: sdk.ObjectTypeTag},
	}

	for _, tc := range valid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.DatabaseRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					Future: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InDatabase:       sdk.Pointer(databaseId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRole.ID(), nil)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = client.Grants.RevokePrivilegesFromDatabaseRole(context.Background(), privileges, on, databaseRole.ID(), nil)
			})
		})
	}

	for _, tc := range invalid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.DatabaseRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					Future: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InDatabase:       sdk.Pointer(databaseId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRole.ID(), nil)
			assert.Error(t, err)
		})
	}
}

// =====================================================
// Tests 9-10: Grant on FUTURE in schema
// =====================================================

func TestInt_GrantOnFutureInSchema_ToAccountRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	schemaId := testClientHelper().Ids.SchemaId()

	type test struct {
		objectType sdk.ObjectType
	}

	valid := []test{
		{objectType: sdk.ObjectTypeAgent},
		{objectType: sdk.ObjectTypeAlert},
		{objectType: sdk.ObjectTypeCortexSearchService},
		{objectType: sdk.ObjectTypeDataMetricFunction},
		{objectType: sdk.ObjectTypeDataset},
		{objectType: sdk.ObjectTypeDbtProject},
		{objectType: sdk.ObjectTypeDynamicTable},
		{objectType: sdk.ObjectTypeEventTable},
		{objectType: sdk.ObjectTypeExternalTable},
		{objectType: sdk.ObjectTypeFileFormat},
		{objectType: sdk.ObjectTypeFunction},
		{objectType: sdk.ObjectTypeGitRepository},
		{objectType: sdk.ObjectTypeImageRepository},
		{objectType: sdk.ObjectTypeIcebergTable},
		{objectType: sdk.ObjectTypeMaterializedView},
		{objectType: sdk.ObjectTypeMcpServer},
		{objectType: sdk.ObjectTypeModel},
		{objectType: sdk.ObjectTypeModelMonitor},
		{objectType: sdk.ObjectTypeNetworkRule},
		{objectType: sdk.ObjectTypeOnlineFeatureTable},
		{objectType: sdk.ObjectTypePipe},
		{objectType: sdk.ObjectTypeProcedure},
		{objectType: sdk.ObjectTypeSecret},
		{objectType: sdk.ObjectTypeSemanticView},
		{objectType: sdk.ObjectTypeService},
		{objectType: sdk.ObjectTypeSequence},
		{objectType: sdk.ObjectTypeStage},
		{objectType: sdk.ObjectTypeStream},
		{objectType: sdk.ObjectTypeStreamlit},
		{objectType: sdk.ObjectTypeTable},
		{objectType: sdk.ObjectTypeTask},
		{objectType: sdk.ObjectTypeView},
		{objectType: sdk.ObjectTypeWorkspace},
	}

	invalid := []test{
		{objectType: sdk.ObjectTypeAggregationPolicy},
		{objectType: sdk.ObjectTypeAuthenticationPolicy},
		{objectType: sdk.ObjectTypeBudget},
		{objectType: sdk.ObjectTypeClassification},
		{objectType: sdk.ObjectTypeExperiment},
		{objectType: sdk.ObjectTypeExternalFunction},
		{objectType: sdk.ObjectTypeGateway},
		{objectType: sdk.ObjectTypeHybridTable},
		{objectType: sdk.ObjectTypeJoinPolicy},
		{objectType: sdk.ObjectTypeMaskingPolicy},
		{objectType: sdk.ObjectTypeNotebook},
		{objectType: sdk.ObjectTypeNotebookProject},
		{objectType: sdk.ObjectTypePackagesPolicy},
		{objectType: sdk.ObjectTypePasswordPolicy},
		{objectType: sdk.ObjectTypePrivacyPolicy},
		{objectType: sdk.ObjectTypeProjectionPolicy},
		{objectType: sdk.ObjectTypeRowAccessPolicy},
		{objectType: sdk.ObjectTypeSessionPolicy},
		{objectType: sdk.ObjectTypeSnapshot},
		{objectType: sdk.ObjectTypeSnapshotPolicy},
		{objectType: sdk.ObjectTypeSnapshotSet},
		{objectType: sdk.ObjectTypeStorageLifecyclePolicy},
		{objectType: sdk.ObjectTypeTag},
	}

	for _, tc := range valid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.AccountRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					Future: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InSchema:         sdk.Pointer(schemaId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = client.Grants.RevokePrivilegesFromAccountRole(context.Background(), privileges, on, role.ID(), nil)
			})
		})
	}

	for _, tc := range invalid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.AccountRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					Future: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InSchema:         sdk.Pointer(schemaId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
			assert.Error(t, err)
		})
	}
}

func TestInt_GrantOnFutureInSchema_ToDatabaseRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)
	schemaId := testClientHelper().Ids.SchemaId()

	type test struct {
		objectType sdk.ObjectType
	}

	valid := []test{
		{objectType: sdk.ObjectTypeAgent},
		{objectType: sdk.ObjectTypeAlert},
		{objectType: sdk.ObjectTypeCortexSearchService},
		{objectType: sdk.ObjectTypeDataMetricFunction},
		{objectType: sdk.ObjectTypeDataset},
		{objectType: sdk.ObjectTypeDbtProject},
		{objectType: sdk.ObjectTypeDynamicTable},
		{objectType: sdk.ObjectTypeEventTable},
		{objectType: sdk.ObjectTypeExternalTable},
		{objectType: sdk.ObjectTypeFileFormat},
		{objectType: sdk.ObjectTypeFunction},
		{objectType: sdk.ObjectTypeGitRepository},
		{objectType: sdk.ObjectTypeImageRepository},
		{objectType: sdk.ObjectTypeIcebergTable},
		{objectType: sdk.ObjectTypeMaterializedView},
		{objectType: sdk.ObjectTypeMcpServer},
		{objectType: sdk.ObjectTypeModel},
		{objectType: sdk.ObjectTypeModelMonitor},
		{objectType: sdk.ObjectTypeNetworkRule},
		{objectType: sdk.ObjectTypeOnlineFeatureTable},
		{objectType: sdk.ObjectTypePipe},
		{objectType: sdk.ObjectTypeProcedure},
		{objectType: sdk.ObjectTypeSecret},
		{objectType: sdk.ObjectTypeSemanticView},
		{objectType: sdk.ObjectTypeService},
		{objectType: sdk.ObjectTypeSequence},
		{objectType: sdk.ObjectTypeStage},
		{objectType: sdk.ObjectTypeStream},
		{objectType: sdk.ObjectTypeStreamlit},
		{objectType: sdk.ObjectTypeTable},
		{objectType: sdk.ObjectTypeTask},
		{objectType: sdk.ObjectTypeView},
		{objectType: sdk.ObjectTypeWorkspace},
	}

	invalid := []test{
		{objectType: sdk.ObjectTypeAggregationPolicy},
		{objectType: sdk.ObjectTypeAuthenticationPolicy},
		{objectType: sdk.ObjectTypeBudget},
		{objectType: sdk.ObjectTypeClassification},
		{objectType: sdk.ObjectTypeExperiment},
		{objectType: sdk.ObjectTypeExternalFunction},
		{objectType: sdk.ObjectTypeGateway},
		{objectType: sdk.ObjectTypeHybridTable},
		{objectType: sdk.ObjectTypeJoinPolicy},
		{objectType: sdk.ObjectTypeMaskingPolicy},
		{objectType: sdk.ObjectTypeNotebook},
		{objectType: sdk.ObjectTypeNotebookProject},
		{objectType: sdk.ObjectTypePackagesPolicy},
		{objectType: sdk.ObjectTypePasswordPolicy},
		{objectType: sdk.ObjectTypePrivacyPolicy},
		{objectType: sdk.ObjectTypeProjectionPolicy},
		{objectType: sdk.ObjectTypeRowAccessPolicy},
		{objectType: sdk.ObjectTypeSessionPolicy},
		{objectType: sdk.ObjectTypeSnapshot},
		{objectType: sdk.ObjectTypeSnapshotPolicy},
		{objectType: sdk.ObjectTypeSnapshotSet},
		{objectType: sdk.ObjectTypeStorageLifecyclePolicy},
		{objectType: sdk.ObjectTypeTag},
	}

	for _, tc := range valid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.DatabaseRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					Future: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InSchema:         sdk.Pointer(schemaId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRole.ID(), nil)
			require.NoError(t, err)
			t.Cleanup(func() {
				_ = client.Grants.RevokePrivilegesFromDatabaseRole(context.Background(), privileges, on, databaseRole.ID(), nil)
			})
		})
	}

	for _, tc := range invalid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			privileges := &sdk.DatabaseRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)}
			on := &sdk.DatabaseRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					Future: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectType.Plural(),
						InSchema:         sdk.Pointer(schemaId),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRole.ID(), nil)
			assert.Error(t, err)
		})
	}
}

// =====================================================
// Test 11: Grant on account object
// =====================================================

func TestInt_GrantOnAccountObject_ToAccountRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	nonExistingId := testClientHelper().Ids.RandomAccountObjectIdentifier()

	type test struct {
		objectType sdk.ObjectType
	}

	// Account-level object types valid for GRANT ... ON <type> <name> TO ROLE.
	// Based on validGrantToAccountObjectTypes in grants_validations.go.
	valid := []test{
		{objectType: sdk.ObjectTypeUser},
		{objectType: sdk.ObjectTypeResourceMonitor},
		{objectType: sdk.ObjectTypeWarehouse},
		{objectType: sdk.ObjectTypeComputePool},
		{objectType: sdk.ObjectTypeDatabase},
		{objectType: sdk.ObjectTypeIntegration},
		{objectType: sdk.ObjectTypeFailoverGroup},
		{objectType: sdk.ObjectTypeReplicationGroup},
		{objectType: sdk.ObjectTypeExternalVolume},
	}

	// Account-level object types NOT valid for GRANT ... ON <type> <name> TO ROLE.
	// These are account-level objects that are either not grantable or use a different grant mechanism.
	//
	// Not testable through GrantOnAccountObject (no struct field):
	//   - ObjectTypeConnection (valid per grants_validations.go, but no SDK struct field; TODO: SNOW-2370066)
	//   - ObjectTypeNetworkPolicy
	//   - ObjectTypeRole
	//   - ObjectTypeShare (uses separate GRANT ... TO SHARE mechanism)
	//   - ObjectTypeApplication
	//   - ObjectTypeApplicationPackage
	invalid := []test{}
	_ = invalid

	for _, tc := range valid {
		t.Run(tc.objectType.String(), func(t *testing.T) {
			t.Parallel()
			err := client.Grants.GrantPrivilegesToAccountRole(
				ctx,
				&sdk.AccountRoleGrantPrivileges{AllPrivileges: sdk.Bool(true)},
				&sdk.AccountRoleGrantOn{
					AccountObject: accountObjectForType(tc.objectType, nonExistingId),
				},
				role.ID(),
				nil,
			)
			assert.Error(t, err)
		})
	}
}
