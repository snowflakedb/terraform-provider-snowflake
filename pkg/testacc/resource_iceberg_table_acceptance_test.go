//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/customassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_IcebergTable_BasicUseCase(t *testing.T) {
	// Snowflake-managed Iceberg tables require an external volume, which needs preconfigured AWS storage.
	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	baseLocation := "iceberg_table_test_table"

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	// Create a dedicated database with the external volume set at db level so the Snowflake-managed table
	// can be created without specifying the external volume explicitly on the resource.
	db, dbCleanup := testClient().Database.CreateDatabaseWithRequest(t, testClient().Database.TestParametersSet(testClient().Ids.RandomAccountObjectIdentifier()).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbCleanup)
	schemaId := sdk.NewDatabaseObjectIdentifier(db.ID().Name(), "PUBLIC")

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0},
		{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg},
	}
	baseLocationChanged := "iceberg_table_test_table_changed"
	comment := random.Comment()
	commentChanged := random.Comment()
	externalComment := random.Comment()

	rowAccessPolicy, rowAccessPolicyCleanup := testClient().RowAccessPolicy.CreateRowAccessPolicyWithDataType(t, testdatatypes.DataTypeNumber)
	t.Cleanup(rowAccessPolicyCleanup)

	aggregationPolicy, aggregationPolicyCleanup := testClient().AggregationPolicy.CreateAggregationPolicy(t)
	t.Cleanup(aggregationPolicyCleanup)

	// modelBasic relies on defaults for all the alterable optional fields.
	modelBasic := model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns)

	// modelWithAlterableOptional sets only the fields that can be altered in place (no ForceNew fields).
	modelWithAlterableOptional := model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithComment(comment).
		WithErrorLogging("true").
		WithTargetFileSize(string(sdk.IcebergTableTargetFileSize64mb)).
		WithDataRetentionTimeInDays(5).
		WithMaxDataExtensionTimeInDays(10).
		WithEnableDataCompaction(false).
		WithEnableIcebergMergeOnRead(false).
		WithRowAccessPolicy(rowAccessPolicy.ID(), "ID").
		WithAggregationPolicy(aggregationPolicy, "ID")

	// modelWithAllOptional sets every optional field explicitly, including the ForceNew ones.
	modelWithAllOptional := model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithCatalog("SNOWFLAKE").
		WithBaseLocation(baseLocation).
		WithExternalVolume(externalVolumeId.Name()).
		WithComment(comment).
		WithChangeTracking("true").
		WithIcebergVersion(2).
		WithPathLayout(string(sdk.IcebergTablePathLayoutFlat)).
		WithErrorLogging("true").
		WithTargetFileSize(string(sdk.IcebergTableTargetFileSize64mb)).
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyOptimized)).
		WithDataRetentionTimeInDays(5).
		WithMaxDataExtensionTimeInDays(10).
		WithEnableDataCompaction(false).
		WithEnableIcebergMergeOnRead(false).
		WithRowAccessPolicy(rowAccessPolicy.ID(), "ID").
		WithAggregationPolicy(aggregationPolicy, "ID").
		WithPartitionBy(
			model.IcebergTablePartitionByIdentity("NAME"),
			model.IcebergTablePartitionByBucket(4, "ID"),
		)

	// modelWithAllOptionalChanged sets every optional field to a value different from modelWithAllOptional,
	// so that reapplying it always forces a new resource (ForceNew fields changed).
	modelWithAllOptionalChanged := model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithCatalog("SNOWFLAKE").
		WithBaseLocation(baseLocationChanged).
		WithExternalVolume(externalVolumeId.Name()).
		WithComment(commentChanged).
		WithChangeTracking("true").
		WithIcebergVersion(3).
		WithPathLayout(string(sdk.IcebergTablePathLayoutHierarchical)).
		WithErrorLogging("false").
		WithTargetFileSize(string(sdk.IcebergTableTargetFileSize128mb)).
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyCompatible)).
		WithDataRetentionTimeInDays(3).
		WithMaxDataExtensionTimeInDays(7).
		WithEnableDataCompaction(true).
		WithEnableIcebergMergeOnRead(true).
		WithRowAccessPolicy(rowAccessPolicy.ID(), "ID").
		WithAggregationPolicy(aggregationPolicy, "ID")

	// modelWithAllOptionalUnset starts from modelWithAllOptionalChanged's ForceNew field values, but omits every alterable optional field from the config to exercise the
	// UNSET code path in UpdateIcebergTable.
	modelWithAllOptionalUnset := model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithCatalog("SNOWFLAKE").
		WithBaseLocation(baseLocationChanged).
		WithExternalVolume(externalVolumeId.Name()).
		WithChangeTracking("true").
		WithIcebergVersion(3).
		WithPathLayout(string(sdk.IcebergTablePathLayoutHierarchical)).
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyCompatible))

	// modelWithClusterBy mirrors modelWithAllOptionalUnset (so no replace is triggered) but sets cluster_by
	// instead of partition_by (the two are mutually exclusive).
	modelWithClusterBy := model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithCatalog("SNOWFLAKE").
		WithBaseLocation(baseLocationChanged).
		WithExternalVolume(externalVolumeId.Name()).
		WithChangeTracking("true").
		WithIcebergVersion(3).
		WithPathLayout(string(sdk.IcebergTablePathLayoutHierarchical)).
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyCompatible)).
		WithClusterBy("ID", "NAME")

	// modelWithAlteredClusterBy mirrors modelWithClusterBy but changes the cluster_by columns, to prove that
	// changing cluster_by triggers an in-place update (ALTER ... CLUSTER BY) rather than a replace.
	modelWithAlteredClusterBy := model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithCatalog("SNOWFLAKE").
		WithBaseLocation(baseLocationChanged).
		WithExternalVolume(externalVolumeId.Name()).
		WithChangeTracking("true").
		WithIcebergVersion(3).
		WithPathLayout(string(sdk.IcebergTablePathLayoutHierarchical)).
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyCompatible)).
		WithClusterBy("NAME")

	ref := modelBasic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasCommentEmpty().
			HasChangeTrackingString("default").
			HasErrorLoggingString("default").
			HasStorageSerializationPolicy(string(sdk.StorageSerializationPolicyOptimized)).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasNoRowAccessPolicy().
			HasNoAggregationPolicy(),
		resourceshowoutputassert.IcebergTableShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasIcebergTableType(sdk.IcebergTableTypeManaged).
			HasComment(""),
		resourceparametersassert.IcebergTableResourceParameters(t, ref).
			HasCatalog("SNOWFLAKE").
			HasCatalogLevel(sdk.ParameterTypeDatabase).
			HasExternalVolume(externalVolumeId.Name()).
			HasExternalVolumeLevel(sdk.ParameterTypeDatabase).
			HasTargetFileSize(sdk.IcebergTableTargetFileSizeAuto).
			HasTargetFileSizeLevel(sdk.ParameterTypeSnowflakeDefault).
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			HasStorageSerializationPolicyLevel(sdk.ParameterTypeTable).
			HasDataRetentionTimeInDays(1).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeDatabase).
			HasMaxDataExtensionTimeInDays(1).
			HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeDatabase).
			HasEnableDataCompaction(true).
			HasEnableDataCompactionLevel(sdk.ParameterTypeSnowflakeDefault).
			HasEnableIcebergMergeOnRead(true).
			HasEnableIcebergMergeOnReadLevel(sdk.ParameterTypeSnowflakeDefault),
	}

	alterableOptionalAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasComment(comment).
			HasErrorLoggingString("true").
			HasTargetFileSize(string(sdk.IcebergTableTargetFileSize64mb)).
			HasDataRetentionTimeInDays(5).
			HasMaxDataExtensionTimeInDays(10).
			HasEnableDataCompaction(false).
			HasEnableIcebergMergeOnRead(false).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasRowAccessPolicy(rowAccessPolicy.ID(), "ID").
			HasAggregationPolicy(aggregationPolicy, "ID"),
		resourceshowoutputassert.IcebergTableShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasIcebergTableType(sdk.IcebergTableTypeManaged).
			HasComment(comment),
		resourceparametersassert.IcebergTableResourceParameters(t, ref).
			HasCatalog("SNOWFLAKE").
			HasCatalogLevel(sdk.ParameterTypeDatabase).
			HasExternalVolume(externalVolumeId.Name()).
			HasExternalVolumeLevel(sdk.ParameterTypeDatabase).
			HasTargetFileSize(sdk.IcebergTableTargetFileSize64mb).
			HasTargetFileSizeLevel(sdk.ParameterTypeTable).
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			HasStorageSerializationPolicyLevel(sdk.ParameterTypeTable).
			HasDataRetentionTimeInDays(5).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeTable).
			HasMaxDataExtensionTimeInDays(10).
			HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeTable).
			HasEnableDataCompaction(false).
			HasEnableDataCompactionLevel(sdk.ParameterTypeTable).
			HasEnableIcebergMergeOnRead(false).
			HasEnableIcebergMergeOnReadLevel(sdk.ParameterTypeTable),
	}

	allOptionalAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasComment(comment).
			HasChangeTrackingString("true").
			HasIcebergVersion(2).
			HasPathLayoutString(string(sdk.IcebergTablePathLayoutFlat)).
			HasErrorLoggingString("true").
			HasTargetFileSize(string(sdk.IcebergTableTargetFileSize64mb)).
			HasStorageSerializationPolicy(string(sdk.StorageSerializationPolicyOptimized)).
			HasDataRetentionTimeInDays(5).
			HasMaxDataExtensionTimeInDays(10).
			HasEnableDataCompaction(false).
			HasEnableIcebergMergeOnRead(false).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasRowAccessPolicy(rowAccessPolicy.ID(), "ID").
			HasAggregationPolicy(aggregationPolicy, "ID").
			HasPartitionByLength(2).
			HasPartitionByIdentity(0, "NAME").
			HasPartitionByBucket(1, 4, "ID"),
		resourceshowoutputassert.IcebergTableShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasExternalVolumeName(externalVolumeId).
			HasIcebergTableType(sdk.IcebergTableTypeManaged).
			HasComment(comment).
			HasIcebergTableFormatVersion(2),
		resourceparametersassert.IcebergTableResourceParameters(t, ref).
			HasCatalog("SNOWFLAKE").
			HasCatalogLevel(sdk.ParameterTypeTable).
			HasExternalVolume(externalVolumeId.FullyQualifiedName()).
			HasExternalVolumeLevel(sdk.ParameterTypeTable).
			HasTargetFileSize(sdk.IcebergTableTargetFileSize64mb).
			HasTargetFileSizeLevel(sdk.ParameterTypeTable).
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			HasStorageSerializationPolicyLevel(sdk.ParameterTypeTable).
			HasDataRetentionTimeInDays(5).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeTable).
			HasMaxDataExtensionTimeInDays(10).
			HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeTable).
			HasEnableDataCompaction(false).
			HasEnableDataCompactionLevel(sdk.ParameterTypeTable).
			HasEnableIcebergMergeOnRead(false).
			HasEnableIcebergMergeOnReadLevel(sdk.ParameterTypeTable),
		assert.Check(resource.TestCheckResourceAttrWith(ref, "base_location", customassert.HasPrefixFunc(baseLocation))),
	}

	changedOptionalAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasComment(commentChanged).
			HasChangeTrackingString("true").
			HasIcebergVersion(3).
			HasPathLayoutString(string(sdk.IcebergTablePathLayoutHierarchical)).
			HasErrorLoggingString("false").
			HasTargetFileSize(string(sdk.IcebergTableTargetFileSize128mb)).
			HasStorageSerializationPolicy(string(sdk.StorageSerializationPolicyCompatible)).
			HasDataRetentionTimeInDays(3).
			HasMaxDataExtensionTimeInDays(7).
			HasEnableDataCompaction(true).
			HasEnableIcebergMergeOnRead(true).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasRowAccessPolicy(rowAccessPolicy.ID(), "ID").
			HasAggregationPolicy(aggregationPolicy, "ID"),
		resourceshowoutputassert.IcebergTableShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasExternalVolumeName(externalVolumeId).
			HasIcebergTableType(sdk.IcebergTableTypeManaged).
			HasComment(commentChanged).
			HasIcebergTableFormatVersion(3),
		resourceparametersassert.IcebergTableResourceParameters(t, ref).
			HasCatalog("SNOWFLAKE").
			HasCatalogLevel(sdk.ParameterTypeTable).
			HasExternalVolume(externalVolumeId.FullyQualifiedName()).
			HasExternalVolumeLevel(sdk.ParameterTypeTable).
			HasTargetFileSize(sdk.IcebergTableTargetFileSize128mb).
			HasTargetFileSizeLevel(sdk.ParameterTypeTable).
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyCompatible).
			HasStorageSerializationPolicyLevel(sdk.ParameterTypeTable).
			HasDataRetentionTimeInDays(3).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeTable).
			HasMaxDataExtensionTimeInDays(7).
			HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeTable).
			HasEnableDataCompaction(true).
			HasEnableDataCompactionLevel(sdk.ParameterTypeTable).
			HasEnableIcebergMergeOnRead(true).
			HasEnableIcebergMergeOnReadLevel(sdk.ParameterTypeTable),
		assert.Check(resource.TestCheckResourceAttrWith(ref, "base_location", customassert.HasPrefixFunc(baseLocationChanged))),
	}

	// unsetOptionalAssertions mirrors basicAssertions for the alterable optional fields (back to Snowflake
	// defaults after being unset), while keeping modelWithAllOptionalChanged's ForceNew field values.
	unsetOptionalAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasCommentEmpty().
			HasChangeTrackingString("true").
			HasIcebergVersion(3).
			HasPathLayoutString(string(sdk.IcebergTablePathLayoutHierarchical)).
			HasErrorLoggingString("default").
			HasTargetFileSize(string(sdk.IcebergTableTargetFileSizeAuto)).
			HasStorageSerializationPolicy(string(sdk.StorageSerializationPolicyCompatible)).
			HasFullyQualifiedNameString(id.FullyQualifiedName()).
			HasNoRowAccessPolicy().
			HasNoAggregationPolicy(),
		resourceshowoutputassert.IcebergTableShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasExternalVolumeName(externalVolumeId).
			HasIcebergTableType(sdk.IcebergTableTypeManaged).
			HasComment("").
			HasIcebergTableFormatVersion(3),
		resourceparametersassert.IcebergTableResourceParameters(t, ref).
			HasCatalog("SNOWFLAKE").
			HasCatalogLevel(sdk.ParameterTypeTable).
			HasExternalVolume(externalVolumeId.FullyQualifiedName()).
			HasExternalVolumeLevel(sdk.ParameterTypeTable).
			HasTargetFileSize(sdk.IcebergTableTargetFileSizeAuto).
			HasTargetFileSizeLevel(sdk.ParameterTypeSnowflakeDefault).
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyCompatible).
			HasStorageSerializationPolicyLevel(sdk.ParameterTypeTable).
			HasDataRetentionTimeInDays(1).
			HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeDatabase).
			HasMaxDataExtensionTimeInDays(1).
			HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeDatabase).
			HasEnableDataCompaction(true).
			HasEnableDataCompactionLevel(sdk.ParameterTypeSnowflakeDefault).
			HasEnableIcebergMergeOnRead(true).
			HasEnableIcebergMergeOnReadLevel(sdk.ParameterTypeSnowflakeDefault),
		assert.Check(resource.TestCheckResourceAttrWith(ref, "base_location", customassert.HasPrefixFunc(baseLocationChanged))),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTable),
		Steps: []resource.TestStep{
			// Create with only required fields
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:                  accconfig.FromModels(t, modelBasic),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"change_tracking", "error_logging", "path_layout", "iceberg_version", "base_location"},
			},
			// Change only alterable fields - expect an in-place update
			{
				Config: accconfig.FromModels(t, modelWithAlterableOptional),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, alterableOptionalAssertions...),
			},
			// Set all possible fields, including ForceNew ones - expect destroy before create
			{
				Config: accconfig.FromModels(t, modelWithAllOptional),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t, allOptionalAssertions...),
			},
			// Import with all fields set
			{
				Config:                  accconfig.FromModels(t, modelWithAllOptional),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"change_tracking", "error_logging", "path_layout", "iceberg_version", "base_location"},
			},
			// The partition spec changes externally (CREATE OR REPLACE with a different partition_by) while
			// the config keeps requesting the original partitioning.
			{
				PreConfig: func() {
					req := sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
						Columns: []sdk.IcebergTableColumnRequest{
							*sdk.NewIcebergTableColumnRequest(columns[0].Name, columns[0].Type),
							*sdk.NewIcebergTableColumnRequest(columns[1].Name, columns[1].Type),
						},
					}).
						WithOrReplace(true).
						WithExternalVolume(externalVolumeId).
						WithCatalog(sdk.IcebergTableCatalogSnowflake).
						WithBaseLocation(baseLocation).
						WithComment(comment).
						WithChangeTracking(true).
						WithIcebergVersion(2).
						WithPathLayout(sdk.IcebergTablePathLayoutFlat).
						WithErrorLogging(true).
						WithTargetFileSize(sdk.IcebergTableTargetFileSize64mb).
						WithStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
						WithDataRetentionTimeInDays(5).
						WithMaxDataExtensionTimeInDays(10).
						WithEnableDataCompaction(false).
						WithEnableIcebergMergeOnRead(false).
						WithPartitionBy([]sdk.IcebergTablePartitionExpressionRequest{
							{Bucket: &sdk.IcebergTablePartitionBucketRequest{Args: sdk.IcebergTablePartitionBucketArgsRequest{NumBuckets: 4, Column: "NAME"}}},
						})
					testClient().IcebergTable.CreateWithRequest(t, req)
				},
				Config: accconfig.FromModels(t, modelWithAllOptional),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t, allOptionalAssertions...),
			},
			// The underlying table gets recreated externally (CREATE OR REPLACE) and every field is changed
			// in the config at the same time - still expect Terraform to destroy and recreate the resource.
			{
				PreConfig: func() {
					req := sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{
						Columns: []sdk.IcebergTableColumnRequest{
							*sdk.NewIcebergTableColumnRequest(columns[0].Name, columns[0].Type),
							*sdk.NewIcebergTableColumnRequest(columns[1].Name, columns[1].Type),
						},
					}).
						WithOrReplace(true).
						WithExternalVolume(externalVolumeId).
						WithCatalog(sdk.IcebergTableCatalogSnowflake).
						WithBaseLocation(baseLocationChanged).
						WithComment(commentChanged).
						WithChangeTracking(true).
						WithIcebergVersion(3).
						WithPathLayout(sdk.IcebergTablePathLayoutHierarchical).
						WithErrorLogging(false).
						WithTargetFileSize(sdk.IcebergTableTargetFileSize128mb).
						WithStorageSerializationPolicy(sdk.StorageSerializationPolicyCompatible).
						WithDataRetentionTimeInDays(3).
						WithMaxDataExtensionTimeInDays(7).
						WithEnableDataCompaction(true).
						WithEnableIcebergMergeOnRead(true)
					testClient().IcebergTable.CreateWithRequest(t, req)
				},
				Config: accconfig.FromModels(t, modelWithAllOptionalChanged),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assertThat(t, changedOptionalAssertions...),
			},
			// Change the alterable fields externally and detect drift - expect Terraform to fix it with an update
			{
				PreConfig: func() {
					testClient().IcebergTable.Alter(t, sdk.NewAlterIcebergTableRequest(id).WithSet(
						*sdk.NewIcebergTableSetPropertiesRequest().
							WithComment(externalComment).
							WithErrorLogging(true).
							WithTargetFileSize(sdk.IcebergTableTargetFileSize64mb).
							WithDataRetentionTimeInDays(1).
							WithMaxDataExtensionTimeInDays(1).
							WithEnableDataCompaction(false).
							WithEnableIcebergMergeOnRead(false),
					))
					testClient().IcebergTable.Alter(t, sdk.NewAlterIcebergTableRequest(id).
						WithDropRowAccessPolicy(*sdk.NewViewDropRowAccessPolicyRequest(rowAccessPolicy.ID())))
					testClient().IcebergTable.Alter(t, sdk.NewAlterIcebergTableRequest(id).
						WithUnsetAggregationPolicy(*sdk.NewViewUnsetAggregationPolicyRequest()))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelWithAllOptionalChanged),
				Check:  assertThat(t, changedOptionalAssertions...),
			},
			// Unset all alterable optional fields (keeping the ForceNew ones) - expect Terraform to issue
			// UNSET for the removed fields via an in-place update, not a replace.
			{
				Config: accconfig.FromModels(t, modelWithAllOptionalUnset),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, unsetOptionalAssertions...),
			},
			// Switch from partition_by to cluster_by - expect an in-place update (cluster_by is not ForceNew)
			{
				Config: accconfig.FromModels(t, modelWithClusterBy),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.IcebergTableResource(t, ref).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasClusterBy("ID", "NAME").
						HasPartitionByEmpty(),
				),
			},
			// Change cluster_by to a different set of columns - expect an in-place update (ALTER ... CLUSTER BY)
			{
				Config: accconfig.FromModels(t, modelWithAlteredClusterBy),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.IcebergTableResource(t, ref).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasClusterBy("NAME").
						HasPartitionByEmpty(),
				),
			},
			// Unset cluster_by - expect an in-place update (ALTER ... DROP CLUSTERING KEY) rather than a replace
			{
				Config: accconfig.FromModels(t, modelWithAllOptionalUnset),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.IcebergTableResource(t, ref).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasClusterByEmpty().
						HasPartitionByEmpty(),
				),
			},
		},
	})
}

func TestAcc_IcebergTable_BasicUseCase_Columns(t *testing.T) {
	// Snowflake-managed Iceberg tables require an external volume, which needs preconfigured AWS storage.
	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	// Create a dedicated database with the external volume set at db level so the Snowflake-managed table
	// can be created without specifying the external volume explicitly on the resource.
	db, dbCleanup := testClient().Database.CreateDatabaseWithRequest(t, testClient().Database.TestParametersSet(testClient().Ids.RandomAccountObjectIdentifier()).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbCleanup)
	schemaId := sdk.NewDatabaseObjectIdentifier(db.ID().Name(), "PUBLIC")

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	maskingPolicy, maskingPolicyCleanup := testClient().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarcharIceberg)
	t.Cleanup(maskingPolicyCleanup)
	maskingPolicyId := maskingPolicy.ID()

	conditionalMaskingPolicy, conditionalMaskingPolicyCleanup := testClient().MaskingPolicy.CreateMaskingPolicy(t)
	t.Cleanup(conditionalMaskingPolicyCleanup)
	conditionalMaskingPolicyId := conditionalMaskingPolicy.ID()

	projectionPolicyId, projectionPolicyCleanup := testClient().ProjectionPolicy.CreateProjectionPolicy(t)
	t.Cleanup(projectionPolicyCleanup)

	// FK constraints must reference a column backed by a UNIQUE or PRIMARY KEY constraint on the target table.
	fkRefTable, fkRefTableCleanup := testClient().Table.CreateWithPredefinedColumnsForIcebergTable(t)
	t.Cleanup(fkRefTableCleanup)

	idComment := random.Comment()
	nameComment := random.Comment()
	statusDefault := "'active'"
	emptyDefault := "''"

	pkConstraintName := random.AlphaN(6)
	uniqueConstraintName := random.AlphaN(6)
	fkConstraintName := random.AlphaN(6)
	checkConstraintName := random.AlphaN(6)

	columns := []model.IcebergTableColumnRequest{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0, NotNull: new("true"), Comment: idComment},
		{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg, Comment: nameComment, MaskingPolicy: &maskingPolicyId},
		{Name: "REGION", Type: testdatatypes.DataTypeVarcharIceberg, ProjectionPolicy: &projectionPolicyId},
		{Name: "STATUS", Type: testdatatypes.DataTypeVarcharIceberg, DefaultExpression: statusDefault},
		{Name: "NOTES", Type: testdatatypes.DataTypeVarcharIceberg, DefaultExpression: emptyDefault},
		{Name: "CATEGORY", Type: testdatatypes.DataTypeVarcharIceberg, MaskingPolicy: &conditionalMaskingPolicyId, MaskingPolicyUsing: []string{"CATEGORY", "STATUS"}},
		{Name: "REF_ID", Type: testdatatypes.DataTypeNumber_38_0},
	}

	primaryKeyConstraint := sdk.TableOutOfLineUniquePKRequest{
		Name:               new(pkConstraintName),
		PrimaryKey:         new(true),
		Columns:            []sdk.Column{{Value: "ID"}},
		NotEnforced:        new(true),
		NotDeferrable:      new(true),
		InitiallyImmediate: new(true),
		Disable:            new(true),
		Novalidate:         new(true),
		Rely:               new(true),
	}
	uniqueConstraint := sdk.TableOutOfLineUniquePKRequest{
		Name:    new(uniqueConstraintName),
		Unique:  new(true),
		Columns: []sdk.Column{{Value: "NAME"}},
	}
	foreignKeyConstraint := sdk.TableOutOfLineFKRequest{
		Name:       new(fkConstraintName),
		Columns:    []sdk.Column{{Value: "REF_ID"}},
		References: fkRefTable.ID(),
		RefColumns: []sdk.Column{{Value: "id"}},
		Match:      new(sdk.SimpleMatchType),
		On: &sdk.ForeignKeyOnAction{
			OnUpdate: new(sdk.ForeignKeyCascadeAction),
			OnDelete: new(sdk.ForeignKeySetNullAction),
		},
	}
	checkConstraint := sdk.TableOutOfLineCHRequest{
		Name:           new(checkConstraintName),
		Expression:     "ID > 0",
		EnableValidate: new(true),
	}

	modelWithColumns := model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), nil).
		WithColumns(columns...).
		WithPrimaryKeyConstraints(primaryKeyConstraint).
		WithUniqueConstraints(uniqueConstraint).
		WithForeignKeyConstraints(foreignKeyConstraint).
		WithCheckConstraints(checkConstraint)

	ref := modelWithColumns.ResourceReference()

	columnAssertions := resourceassert.IcebergTableResource(t, ref).
		HasDatabaseString(id.DatabaseName()).
		HasSchemaString(id.SchemaName()).
		HasNameString(id.Name()).
		HasColumns(
			resourceassert.ExpectedColumn{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0.ToSql(), NotNull: true, Comment: idComment},
			resourceassert.ExpectedColumn{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), Comment: nameComment, MaskingPolicy: &maskingPolicyId, MaskingPolicyUsing: []string{"NAME"}},
			resourceassert.ExpectedColumn{Name: "REGION", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), ProjectionPolicy: &projectionPolicyId},
			resourceassert.ExpectedColumn{Name: "STATUS", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), DefaultExpression: statusDefault},
			resourceassert.ExpectedColumn{Name: "NOTES", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), DefaultExpression: emptyDefault},
			resourceassert.ExpectedColumn{Name: "CATEGORY", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), MaskingPolicy: &conditionalMaskingPolicyId, MaskingPolicyUsing: []string{"CATEGORY", "STATUS"}},
			resourceassert.ExpectedColumn{Name: "REF_ID", Type: testdatatypes.DataTypeNumber_38_0.ToSql()},
		).
		HasPrimaryKeyConstraints(primaryKeyConstraint).
		HasUniqueConstraints(uniqueConstraint).
		HasForeignKeyConstraints(foreignKeyConstraint).
		HasCheckConstraints(checkConstraint)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTable),
		Steps: []resource.TestStep{
			// Create with columns exercising not_null, comment, default, masking_policy (plain and conditional
			// with using), projection_policy, and all constraint kinds: primary key, unique, foreign key, check.
			{
				Config: accconfig.FromModels(t, modelWithColumns),
				Check:  assertThat(t, columnAssertions),
			},
			// Import
			{
				Config:            accconfig.FromModels(t, modelWithColumns),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// Out-of-line constraints are ForceNew-only and are not read back by the resource, so they
				// cannot be verified against imported state.
				ImportStateVerifyIgnore: []string{
					"change_tracking", "error_logging", "path_layout", "iceberg_version", "base_location",
					"primary_key_constraint", "unique_constraint", "foreign_key_constraint", "check_constraint",
				},
			},
		},
	})
}

// TestAcc_IcebergTable_BasicUseCase_ColumnAlters proves that changing the "column" list in place
// (adding columns, dropping columns, renaming columns, and altering not_null/comment/masking_policy/
// projection_policy on existing columns) is always applied via ALTER ICEBERG TABLE - never by
// destroying and recreating the resource. Every step below asserts plancheck.ResourceActionUpdate
// (never ResourceActionDestroyBeforeCreate/ResourceActionReplace) to make that guarantee explicit,
// and cross-checks the applied state against a DESCRIBE-backed read using resourceassert.HasColumns.
func TestAcc_IcebergTable_BasicUseCase_ColumnAlters(t *testing.T) {
	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	// Create a dedicated database with the external volume set at db level so the Snowflake-managed table
	// can be created without specifying the external volume explicitly on the resource.
	db, dbCleanup := testClient().Database.CreateDatabaseWithRequest(t, testClient().Database.TestParametersSet(testClient().Ids.RandomAccountObjectIdentifier()).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbCleanup)
	schemaId := sdk.NewDatabaseObjectIdentifier(db.ID().Name(), "PUBLIC")

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	// Two arity-1 identity masking policies and two projection policies, so masking_policy/projection_policy
	// can be swapped from one policy to another without changing the column's data type.
	maskingPolicy1, maskingPolicy1Cleanup := testClient().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarcharIceberg)
	t.Cleanup(maskingPolicy1Cleanup)
	maskingPolicy1Id := maskingPolicy1.ID()

	maskingPolicy2, maskingPolicy2Cleanup := testClient().MaskingPolicy.CreateMaskingPolicyIdentity(t, testdatatypes.DataTypeVarcharIceberg)
	t.Cleanup(maskingPolicy2Cleanup)
	maskingPolicy2Id := maskingPolicy2.ID()

	projectionPolicy1Id, projectionPolicy1Cleanup := testClient().ProjectionPolicy.CreateProjectionPolicy(t)
	t.Cleanup(projectionPolicy1Cleanup)

	projectionPolicy2Id, projectionPolicy2Cleanup := testClient().ProjectionPolicy.CreateProjectionPolicy(t)
	t.Cleanup(projectionPolicy2Cleanup)

	idComment := random.Comment()
	idCommentChanged := random.Comment()
	statusComment := random.Comment()

	ref := model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), nil).ResourceReference()

	newModel := func(columns ...model.IcebergTableColumnRequest) *model.IcebergTableModel {
		return model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), nil).WithColumns(columns...)
	}

	updateStep := func(config *model.IcebergTableModel, columnAssertions resourceassert.IcebergTableResourceAssert, expectedChanges ...plancheck.PlanCheck) resource.TestStep {
		return resource.TestStep{
			Config: accconfig.FromModels(t, config),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: append([]plancheck.PlanCheck{
					plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
				}, expectedChanges...),
			},
			Check: assertThat(t, &columnAssertions),
		}
	}

	// v1 (create): three columns - ID (not_null + comment), NAME (plain), REGION (masking_policy1).
	v1 := newModel(
		model.IcebergTableColumnRequest{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0, NotNull: new("true"), Comment: idComment},
		model.IcebergTableColumnRequest{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg},
		model.IcebergTableColumnRequest{Name: "REGION", Type: testdatatypes.DataTypeVarcharIceberg, MaskingPolicy: &maskingPolicy1Id},
	)
	v1Assertions := *resourceassert.IcebergTableResource(t, ref).HasColumns(
		resourceassert.ExpectedColumn{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0.ToSql(), NotNull: true, Comment: idComment},
		resourceassert.ExpectedColumn{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg.ToSql()},
		resourceassert.ExpectedColumn{Name: "REGION", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), MaskingPolicy: &maskingPolicy1Id, MaskingPolicyUsing: []string{"REGION"}},
	)

	// v2 (add a column at the end + alter existing columns in place): adds STATUS (not_null + comment set
	// directly at ADD COLUMN time), sets NAME.not_null, changes ID's comment, and swaps REGION's masking
	// policy - all in a single apply.
	v2 := newModel(
		model.IcebergTableColumnRequest{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0, NotNull: new("true"), Comment: idCommentChanged},
		model.IcebergTableColumnRequest{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg, NotNull: new("true")},
		model.IcebergTableColumnRequest{Name: "REGION", Type: testdatatypes.DataTypeVarcharIceberg, MaskingPolicy: &maskingPolicy2Id},
		model.IcebergTableColumnRequest{Name: "STATUS", Type: testdatatypes.DataTypeVarcharIceberg, NotNull: new("true"), Comment: statusComment},
	)
	v2Assertions := *resourceassert.IcebergTableResource(t, ref).HasColumns(
		resourceassert.ExpectedColumn{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0.ToSql(), NotNull: true, Comment: idCommentChanged},
		resourceassert.ExpectedColumn{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), NotNull: true},
		resourceassert.ExpectedColumn{Name: "REGION", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), MaskingPolicy: &maskingPolicy2Id, MaskingPolicyUsing: []string{"REGION"}},
		resourceassert.ExpectedColumn{Name: "STATUS", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), NotNull: true, Comment: statusComment},
	)

	// v3 (add another column with projection_policy set at ADD COLUMN time, rename REGION -> REGION_CODE,
	// drop NAME's not_null, and unset ID's comment) - again all in a single apply.
	v3 := newModel(
		model.IcebergTableColumnRequest{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0, NotNull: new("true")},
		model.IcebergTableColumnRequest{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg, NotNull: new("false")},
		model.IcebergTableColumnRequest{Name: "REGION_CODE", Type: testdatatypes.DataTypeVarcharIceberg, MaskingPolicy: &maskingPolicy2Id},
		model.IcebergTableColumnRequest{Name: "STATUS", Type: testdatatypes.DataTypeVarcharIceberg, NotNull: new("true"), Comment: statusComment},
		model.IcebergTableColumnRequest{Name: "CATEGORY", Type: testdatatypes.DataTypeVarcharIceberg, ProjectionPolicy: &projectionPolicy1Id},
	)
	v3Assertions := *resourceassert.IcebergTableResource(t, ref).HasColumns(
		resourceassert.ExpectedColumn{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0.ToSql(), NotNull: true},
		resourceassert.ExpectedColumn{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg.ToSql()},
		resourceassert.ExpectedColumn{Name: "REGION_CODE", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), MaskingPolicy: &maskingPolicy2Id, MaskingPolicyUsing: []string{"REGION_CODE"}},
		resourceassert.ExpectedColumn{Name: "STATUS", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), NotNull: true, Comment: statusComment},
		resourceassert.ExpectedColumn{Name: "CATEGORY", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), ProjectionPolicy: &projectionPolicy1Id},
	)

	// v3b (swap CATEGORY's projection policy to a different one) - proves SetProjectionPolicyOnColumn
	// correctly overwrites an existing projection policy rather than requiring an unset first.
	v3b := newModel(
		model.IcebergTableColumnRequest{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0, NotNull: new("true")},
		model.IcebergTableColumnRequest{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg, NotNull: new("false")},
		model.IcebergTableColumnRequest{Name: "REGION_CODE", Type: testdatatypes.DataTypeVarcharIceberg, MaskingPolicy: &maskingPolicy2Id},
		model.IcebergTableColumnRequest{Name: "STATUS", Type: testdatatypes.DataTypeVarcharIceberg, NotNull: new("true"), Comment: statusComment},
		model.IcebergTableColumnRequest{Name: "CATEGORY", Type: testdatatypes.DataTypeVarcharIceberg, ProjectionPolicy: &projectionPolicy2Id},
	)
	v3bAssertions := *resourceassert.IcebergTableResource(t, ref).HasColumns(
		resourceassert.ExpectedColumn{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0.ToSql(), NotNull: true},
		resourceassert.ExpectedColumn{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg.ToSql()},
		resourceassert.ExpectedColumn{Name: "REGION_CODE", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), MaskingPolicy: &maskingPolicy2Id, MaskingPolicyUsing: []string{"REGION_CODE"}},
		resourceassert.ExpectedColumn{Name: "STATUS", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), NotNull: true, Comment: statusComment},
		resourceassert.ExpectedColumn{Name: "CATEGORY", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), ProjectionPolicy: &projectionPolicy2Id},
	)

	// v4 (unset REGION_CODE's masking policy and drop the trailing CATEGORY column - exercising DROP
	// COLUMN together with an in-place alter in the same apply).
	v4 := newModel(
		model.IcebergTableColumnRequest{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0, NotNull: new("true")},
		model.IcebergTableColumnRequest{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg, NotNull: new("false")},
		model.IcebergTableColumnRequest{Name: "REGION_CODE", Type: testdatatypes.DataTypeVarcharIceberg},
		model.IcebergTableColumnRequest{Name: "STATUS", Type: testdatatypes.DataTypeVarcharIceberg, NotNull: new("true"), Comment: statusComment},
	)
	v4Assertions := *resourceassert.IcebergTableResource(t, ref).HasColumns(
		resourceassert.ExpectedColumn{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0.ToSql(), NotNull: true},
		resourceassert.ExpectedColumn{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg.ToSql()},
		resourceassert.ExpectedColumn{Name: "REGION_CODE", Type: testdatatypes.DataTypeVarcharIceberg.ToSql()},
		resourceassert.ExpectedColumn{Name: "STATUS", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), NotNull: true, Comment: statusComment},
	)

	// v5 (drop two trailing columns - STATUS and REGION_CODE - in a single apply, leaving only ID and NAME).
	v5 := newModel(
		model.IcebergTableColumnRequest{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0, NotNull: new("true")},
		model.IcebergTableColumnRequest{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg, NotNull: new("false")},
	)
	v5Assertions := *resourceassert.IcebergTableResource(t, ref).HasColumns(
		resourceassert.ExpectedColumn{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0.ToSql(), NotNull: true},
		resourceassert.ExpectedColumn{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg.ToSql()},
	)

	// v6 (add two new columns - REGION2 and EXTRA - in a single apply, proving multi-column ADD works too).
	v6 := newModel(
		model.IcebergTableColumnRequest{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0, NotNull: new("true")},
		model.IcebergTableColumnRequest{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg, NotNull: new("false")},
		model.IcebergTableColumnRequest{Name: "REGION2", Type: testdatatypes.DataTypeVarcharIceberg, Comment: statusComment},
		model.IcebergTableColumnRequest{Name: "EXTRA", Type: testdatatypes.DataTypeVarcharIceberg, NotNull: new("true")},
	)
	v6Assertions := *resourceassert.IcebergTableResource(t, ref).HasColumns(
		resourceassert.ExpectedColumn{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0.ToSql(), NotNull: true},
		resourceassert.ExpectedColumn{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg.ToSql()},
		resourceassert.ExpectedColumn{Name: "REGION2", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), Comment: statusComment},
		resourceassert.ExpectedColumn{Name: "EXTRA", Type: testdatatypes.DataTypeVarcharIceberg.ToSql(), NotNull: true},
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTable),
		Steps: []resource.TestStep{
			// Create with the initial three columns.
			{
				Config: accconfig.FromModels(t, v1),
				Check:  assertThat(t, &v1Assertions),
			},
			// Add a column at the end (with not_null + comment set at ADD COLUMN time) while altering
			// not_null/comment/masking_policy on existing columns - expect an in-place update.
			updateStep(
				v2, v2Assertions,
				planchecks.ExpectChange(ref, "column.0.comment", tfjson.ActionUpdate, sdk.String(idComment), sdk.String(idCommentChanged)),
				planchecks.ExpectChange(ref, "column.1.not_null", tfjson.ActionUpdate, sdk.String("false"), sdk.String("true")),
				planchecks.ExpectChange(ref, "column.2.masking_policy.0.policy_name", tfjson.ActionUpdate, sdk.String(maskingPolicy1Id.FullyQualifiedName()), sdk.String(maskingPolicy2Id.FullyQualifiedName())),
			),
			// Add another column (with projection_policy set at ADD COLUMN time), rename a column, and
			// alter not_null/comment on existing columns - expect an in-place update.
			updateStep(
				v3, v3Assertions,
				planchecks.ExpectChange(ref, "column.0.comment", tfjson.ActionUpdate, sdk.String(idCommentChanged), nil),
				planchecks.ExpectChange(ref, "column.1.not_null", tfjson.ActionUpdate, sdk.String("true"), sdk.String("false")),
				planchecks.ExpectChange(ref, "column.2.name", tfjson.ActionUpdate, sdk.String("REGION"), sdk.String("REGION_CODE")),
			),
			// Swap the projection policy on an existing column to a different policy - expect an
			// in-place update.
			updateStep(
				v3b, v3bAssertions,
				planchecks.ExpectChange(ref, "column.4.projection_policy.0.policy_name", tfjson.ActionUpdate, sdk.String(projectionPolicy1Id.FullyQualifiedName()), sdk.String(projectionPolicy2Id.FullyQualifiedName())),
			),
			// Unset a masking policy and drop the trailing column - expect an in-place update.
			// No planchecks.ExpectChange here: unsetting "masking_policy" plans it as an empty list
			// (not null), and the dropped column has no counterpart index in the new list - both would
			// make ExpectChange index out of range into the JSON plan representation.
			updateStep(v4, v4Assertions),
			// Drop two trailing columns in one apply - expect an in-place update.
			// No planchecks.ExpectChange here: the surviving columns (ID, NAME) are unchanged, and the
			// dropped ones have no counterpart index in the new list to diff against.
			updateStep(v5, v5Assertions),
			// Add two columns back in one apply - expect an in-place update.
			// No planchecks.ExpectChange here: the surviving columns (ID, NAME) are unchanged, and the
			// added ones have no counterpart index in the old list to diff against.
			updateStep(v6, v6Assertions),
			// Import at the final state to prove the full column list (including masking/projection
			// policies picked up along the way) round-trips through Read.
			{
				Config:                  accconfig.FromModels(t, v6),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"change_tracking", "error_logging", "path_layout", "iceberg_version", "base_location"},
			},
		},
	})
}

func TestAcc_IcebergTable_CompleteUseCase(t *testing.T) {
	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	baseLocation := random.String()

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0},
		{Name: "NAME", Type: testdatatypes.DataTypeVarcharIceberg},
	}
	comment := random.Comment()

	rowAccessPolicy, rowAccessPolicyCleanup := testClient().RowAccessPolicy.CreateRowAccessPolicyWithDataType(t, testdatatypes.DataTypeNumber)
	t.Cleanup(rowAccessPolicyCleanup)

	aggregationPolicy, aggregationPolicyCleanup := testClient().AggregationPolicy.CreateAggregationPolicy(t)
	t.Cleanup(aggregationPolicyCleanup)

	modelComplete := model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns).
		WithCatalog("SNOWFLAKE").
		WithBaseLocation(baseLocation).
		WithExternalVolume(externalVolumeId.Name()).
		WithComment(comment).
		WithChangeTracking("true").
		WithIcebergVersion(2).
		WithPathLayout(string(sdk.IcebergTablePathLayoutFlat)).
		WithErrorLogging("true").
		WithTargetFileSize(string(sdk.IcebergTableTargetFileSize64mb)).
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyOptimized)).
		WithDataRetentionTimeInDays(5).
		WithMaxDataExtensionTimeInDays(10).
		WithEnableDataCompaction(false).
		WithEnableIcebergMergeOnRead(false).
		WithRowAccessPolicy(rowAccessPolicy.ID(), "ID").
		WithAggregationPolicy(aggregationPolicy, "ID").
		WithPartitionBy(
			model.IcebergTablePartitionByBucket(4, "ID"),
			model.IcebergTablePartitionByTruncate(10, "NAME"),
		)

	ref := modelComplete.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTable),
		Steps: []resource.TestStep{
			// Create with all fields
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(
					t,
					resourceassert.IcebergTableResource(t, ref).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasExternalVolume(externalVolumeId.FullyQualifiedName()).
						HasComment(comment).
						HasChangeTrackingString("true").
						HasIcebergVersion(2).
						HasPathLayoutString(string(sdk.IcebergTablePathLayoutFlat)).
						HasErrorLoggingString("true").
						HasTargetFileSize(string(sdk.IcebergTableTargetFileSize64mb)).
						HasStorageSerializationPolicy(string(sdk.StorageSerializationPolicyOptimized)).
						HasDataRetentionTimeInDays(5).
						HasMaxDataExtensionTimeInDays(10).
						HasEnableDataCompaction(false).
						HasEnableIcebergMergeOnRead(false).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasRowAccessPolicy(rowAccessPolicy.ID(), "ID").
						HasAggregationPolicy(aggregationPolicy, "ID").
						HasPartitionByLength(2).
						HasPartitionByBucket(0, 4, "ID").
						HasPartitionByTruncate(1, 10, "NAME"),
					resourceshowoutputassert.IcebergTableShowOutput(t, ref).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasExternalVolumeName(externalVolumeId).
						HasIcebergTableType(sdk.IcebergTableTypeManaged).
						HasComment(comment).
						HasIcebergTableFormatVersion(2),
					resourceparametersassert.IcebergTableResourceParameters(t, ref).
						HasCatalog("SNOWFLAKE").
						HasCatalogLevel(sdk.ParameterTypeTable).
						HasExternalVolume(externalVolumeId.FullyQualifiedName()).
						HasExternalVolumeLevel(sdk.ParameterTypeTable).
						HasTargetFileSize(sdk.IcebergTableTargetFileSize64mb).
						HasTargetFileSizeLevel(sdk.ParameterTypeTable).
						HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
						HasStorageSerializationPolicyLevel(sdk.ParameterTypeTable).
						HasDataRetentionTimeInDays(5).
						HasDataRetentionTimeInDaysLevel(sdk.ParameterTypeTable).
						HasMaxDataExtensionTimeInDays(10).
						HasMaxDataExtensionTimeInDaysLevel(sdk.ParameterTypeTable).
						HasEnableDataCompaction(false).
						HasEnableDataCompactionLevel(sdk.ParameterTypeTable).
						HasEnableIcebergMergeOnRead(false).
						HasEnableIcebergMergeOnReadLevel(sdk.ParameterTypeTable),
					assert.Check(resource.TestCheckResourceAttrWith(ref, "base_location", customassert.HasPrefixFunc(baseLocation))),
				),
			},
			// Import
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"change_tracking", "error_logging", "path_layout", "iceberg_version", "base_location"},
			},
		},
	})
}

func TestAcc_IcebergTable_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	columns := []sdk.TableColumnSignature{
		{Name: "ID", Type: testdatatypes.DataTypeNumber_38_0},
	}

	baseModel := func() *model.IcebergTableModel {
		return model.IcebergTableWithDefaultMeta(id.DatabaseName(), id.SchemaName(), id.Name(), columns)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTable),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, baseModel().WithPathLayout("INVALID")),
				ExpectError: regexp.MustCompile(`expected .*path_layout.* to be one of .*, got INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, baseModel().WithErrorLogging("INVALID")),
				ExpectError: regexp.MustCompile(`expected .*error_logging.* to be one of .*, got INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, baseModel().WithChangeTracking("INVALID")),
				ExpectError: regexp.MustCompile(`expected .*change_tracking.* to be one of .*, got INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, baseModel().WithTargetFileSize("INVALID")),
				ExpectError: regexp.MustCompile(`expected .*target_file_size.* to be one of .*, got INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, baseModel().WithStorageSerializationPolicy("INVALID")),
				ExpectError: regexp.MustCompile(`expected .*storage_serialization_policy.* to be one of .*, got INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, baseModel().WithIcebergVersion(0)),
				ExpectError: regexp.MustCompile(`expected .*iceberg_version.* to be at least \(1\), got 0`),
			},
			{
				Config: accconfig.FromModels(
					t, baseModel().
						WithPartitionBy(model.IcebergTablePartitionByIdentity("ID")).
						WithClusterBy("ID"),
				),
				ExpectError: regexp.MustCompile(`"cluster_by": conflicts with partition_by`),
			},
		},
	})
}
