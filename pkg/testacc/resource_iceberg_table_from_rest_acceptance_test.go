//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_IcebergTableFromRest_BasicUseCase(t *testing.T) {
	// TODO(SNOW-3725859): Provide the external volume and catalog integration dynamically. Unskip and
	// fold these tests into the main suite.
	// Also, add tests for detecting external volume and catalog changes.
	// t.Skip("Iceberg REST tests require preconfigured external catalog integrations and are not run by default")
	const (
		icebergRestCatalogName = "REST_CATALOG_INTEGRATION"
		icebergRestVolumeName  = "GLUE_EXTERNAL_VOLUME"
		// Values that must match the manually preconfigured REST catalog contents.
		icebergRestCatalogTableName = "TEST"
		icebergRestCatalogNamespace = "glue_iceberg_schema"
	)
	externalVolumeId := sdk.NewAccountObjectIdentifier(icebergRestVolumeName)
	catalogId := sdk.NewAccountObjectIdentifier(icebergRestCatalogName)

	// Create a dedicated database with external_volume and catalog set at db level so the table
	// can be created without specifying them explicitly (matching the "required fields only" test case).
	dbForIcebergRest, dbCleanup := testClient().Database.CreateDatabaseWithRequest(t, sdk.NewCreateDatabaseRequest(testClient().Ids.RandomAccountObjectIdentifier()).WithCatalog(catalogId).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbCleanup)
	schemaIdForIcebergRest := sdk.NewDatabaseObjectIdentifier(dbForIcebergRest.ID().Name(), "PUBLIC")

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForIcebergRest)
	comment := random.Comment()
	externalComment := random.Comment()

	// modelBasic relies on db-level external_volume and catalog defaults — no explicit values.
	modelBasic := model.IcebergTableFromRestWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		icebergRestCatalogTableName,
	)

	// modelWithAllOptional only sets alterable fields so the transition from modelBasic is an update
	// (not a force-new recreate). The force-new fields are covered by the complete use case.
	modelWithAllOptional := model.IcebergTableFromRestWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		icebergRestCatalogTableName,
	).WithComment(comment).
		WithReplaceInvalidCharacters(true).
		WithAutoRefresh("true").
		WithTargetFileSize(string(sdk.IcebergTableTargetFileSize64mb)).
		WithEnableIcebergMergeOnRead(true)

	ref := modelBasic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableFromRestResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasCatalogTableNameString(icebergRestCatalogTableName).
			HasNoCatalogNamespace().
			HasNoPathLayout().
			HasExternalVolumeString(externalVolumeId.Name()).
			HasCatalogString(catalogId.Name()).
			HasAutoRefreshString(r.BooleanDefault).
			HasTargetFileSizeString(string(sdk.IcebergTableTargetFileSizeAuto)).
			HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyOptimized)).
			HasIcebergMergeOnReadBehaviorString("auto").
			HasEnableIcebergMergeOnRead(true).
			HasCommentEmpty().
			HasReplaceInvalidCharacters(false).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.IcebergTableShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasExternalVolumeName(externalVolumeId).
			HasCatalogName(catalogId).
			HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
			HasCatalogTableName(icebergRestCatalogTableName).
			HasCatalogNamespace(icebergRestCatalogNamespace).
			HasCanWriteMetadata(true).
			HasComment("").
			HasNameMapping("").
			HasOwnerRoleType("ROLE").
			HasCatalogSyncName("").
			HasAutoRefreshStatusEmpty(),
		resourceparametersassert.IcebergTableResourceParameters(t, ref).
			HasExternalVolume(externalVolumeId.Name()).
			HasExternalVolumeLevel(sdk.ParameterTypeDatabase).
			HasCatalog(catalogId.Name()).
			HasCatalogLevel(sdk.ParameterTypeDatabase).
			HasReplaceInvalidCharacters(false).
			HasReplaceInvalidCharactersLevel(sdk.ParameterTypeSnowflakeDefault).
			HasTargetFileSize(sdk.IcebergTableTargetFileSizeAuto).
			HasTargetFileSizeLevel(sdk.ParameterTypeSnowflakeDefault).
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			HasStorageSerializationPolicyLevel(sdk.ParameterTypeSnowflakeDefault).
			HasEnableIcebergMergeOnRead(true).
			HasEnableIcebergMergeOnReadLevel(sdk.ParameterTypeSnowflakeDefault).
			HasIcebergMergeOnReadBehavior("auto").
			HasIcebergMergeOnReadBehaviorLevel(sdk.ParameterTypeSnowflakeDefault),
	}

	allOptionalAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableFromRestResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasCatalogTableNameString(icebergRestCatalogTableName).
			HasNoCatalogNamespace().
			HasNoPathLayout().
			HasExternalVolumeString(externalVolumeId.Name()).
			HasCatalogString(catalogId.Name()).
			HasComment(comment).
			HasReplaceInvalidCharacters(true).
			HasAutoRefreshString("true").
			HasTargetFileSizeString(string(sdk.IcebergTableTargetFileSize64mb)).
			HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyOptimized)).
			HasIcebergMergeOnReadBehaviorString("auto").
			HasEnableIcebergMergeOnRead(true).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.IcebergTableShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasExternalVolumeName(externalVolumeId).
			HasCatalogName(catalogId).
			HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
			HasCatalogTableName(icebergRestCatalogTableName).
			HasCatalogNamespace(icebergRestCatalogNamespace).
			HasCanWriteMetadata(true).
			HasComment(comment).
			HasNameMapping("").
			HasOwnerRoleType("ROLE").
			HasCatalogSyncName("").
			HasAutoRefreshStatusNotEmpty(),
		resourceparametersassert.IcebergTableResourceParameters(t, ref).
			HasExternalVolume(externalVolumeId.Name()).
			HasExternalVolumeLevel(sdk.ParameterTypeDatabase).
			HasCatalog(catalogId.Name()).
			HasCatalogLevel(sdk.ParameterTypeDatabase).
			HasReplaceInvalidCharacters(true).
			HasReplaceInvalidCharactersLevel(sdk.ParameterTypeTable).
			HasTargetFileSize(sdk.IcebergTableTargetFileSize64mb).
			HasTargetFileSizeLevel(sdk.ParameterTypeTable).
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			HasStorageSerializationPolicyLevel(sdk.ParameterTypeSnowflakeDefault).
			HasEnableIcebergMergeOnRead(true).
			HasEnableIcebergMergeOnReadLevel(sdk.ParameterTypeTable).
			HasIcebergMergeOnReadBehavior("auto").
			HasIcebergMergeOnReadBehaviorLevel(sdk.ParameterTypeSnowflakeDefault),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromRest),
		Steps: []resource.TestStep{
			// Create with only required fields (external_volume and catalog come from db defaults)
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:            accconfig.FromModels(t, modelBasic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"path_layout",
					"auto_refresh",
				},
			},
			// Set the alterable optional fields
			{
				Config: accconfig.FromModels(t, modelWithAllOptional),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, allOptionalAssertions...),
			},
			// Import
			{
				Config:            accconfig.FromModels(t, modelWithAllOptional),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"path_layout",
					"auto_refresh",
				},
			},
			// Change alterable fields externally and detect drift
			{
				PreConfig: func() {
					testClient().IcebergTable.Alter(t, sdk.NewAlterIcebergTableRequest(id).WithSet(
						*sdk.NewIcebergTableSetPropertiesRequest().
							WithComment(externalComment),
					))
					testClient().IcebergTable.Alter(t, sdk.NewAlterIcebergTableRequest(id).WithSet(
						*sdk.NewIcebergTableSetPropertiesRequest().
							WithReplaceInvalidCharacters(false),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", new(comment), new(externalComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, new(externalComment), new(comment)),
						planchecks.ExpectDrift(ref, "replace_invalid_characters", new("true"), new("false")),
						planchecks.ExpectChange(ref, "replace_invalid_characters", tfjson.ActionUpdate, new("false"), new("true")),
					},
				},
				Config: accconfig.FromModels(t, modelWithAllOptional),
				Check:  assertThat(t, allOptionalAssertions...),
			},
			// Bring back only required fields
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, basicAssertions...),
			},
		},
	})
}

func TestAcc_IcebergTableFromRest_CompleteUseCase(t *testing.T) {
	// TODO(SNOW-3725859): Provide the external volume and catalog integration dynamically. Unskip and
	// fold these tests into the main suite.
	// t.Skip("Iceberg REST tests require preconfigured external catalog integrations and are not run by default")
	const (
		icebergRestCatalogName = "REST_CATALOG_INTEGRATION"
		icebergRestVolumeName  = "GLUE_EXTERNAL_VOLUME"
		// Values that must match the manually preconfigured REST catalog contents.
		icebergRestCatalogTableName = "TEST"
		icebergRestCatalogNamespace = "glue_iceberg_schema"
	)
	externalVolumeId := sdk.NewAccountObjectIdentifier(icebergRestVolumeName)
	catalogId := sdk.NewAccountObjectIdentifier(icebergRestCatalogName)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.IcebergTableFromRestWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		icebergRestCatalogTableName,
	).WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithCatalogNamespace(icebergRestCatalogNamespace).
		WithPathLayout(string(sdk.IcebergTablePathLayoutHierarchical)).
		WithAutoRefresh("true").
		WithTargetFileSize(string(sdk.IcebergTableTargetFileSize128mb)).
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyOptimized)).
		WithIcebergMergeOnReadBehavior(string(sdk.IcebergTableIcebergMergeOnReadBehaviorEnabled)).
		WithEnableIcebergMergeOnRead(true).
		WithReplaceInvalidCharacters(true).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromRest),
		Steps: []resource.TestStep{
			// Create with all fields
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(
					t,
					resourceassert.IcebergTableFromRestResource(t, modelComplete.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasCatalogTableNameString(icebergRestCatalogTableName).
						HasCatalogNamespaceString(icebergRestCatalogNamespace).
						HasPathLayoutString(string(sdk.IcebergTablePathLayoutHierarchical)).
						HasExternalVolumeString(externalVolumeId.FullyQualifiedName()).
						HasCatalogString(catalogId.FullyQualifiedName()).
						HasAutoRefreshString("true").
						HasTargetFileSizeString(string(sdk.IcebergTableTargetFileSize128mb)).
						HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyOptimized)).
						HasIcebergMergeOnReadBehaviorString(string(sdk.IcebergTableIcebergMergeOnReadBehaviorEnabled)).
						HasEnableIcebergMergeOnRead(true).
						HasComment(comment).
						HasReplaceInvalidCharacters(true).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.IcebergTableShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasExternalVolumeName(externalVolumeId).
						HasCatalogName(catalogId).
						HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
						HasCatalogTableName(icebergRestCatalogTableName).
						HasCatalogNamespace(icebergRestCatalogNamespace).
						HasCanWriteMetadata(true).
						HasComment(comment).
						HasNameMapping("").
						HasOwnerRoleType("ROLE").
						HasCatalogSyncName("").
						HasAutoRefreshStatusNotEmpty(),
					resourceparametersassert.IcebergTableResourceParameters(t, modelComplete.ResourceReference()).
						HasExternalVolume(externalVolumeId.FullyQualifiedName()).
						HasExternalVolumeLevel(sdk.ParameterTypeTable).
						HasCatalog(catalogId.FullyQualifiedName()).
						HasCatalogLevel(sdk.ParameterTypeTable).
						HasReplaceInvalidCharacters(true).
						HasReplaceInvalidCharactersLevel(sdk.ParameterTypeTable).
						HasTargetFileSize(sdk.IcebergTableTargetFileSize128mb).
						HasTargetFileSizeLevel(sdk.ParameterTypeTable).
						HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
						HasStorageSerializationPolicyLevel(sdk.ParameterTypeTable).
						HasEnableIcebergMergeOnRead(true).
						HasEnableIcebergMergeOnReadLevel(sdk.ParameterTypeTable).
						HasIcebergMergeOnReadBehavior(string(sdk.IcebergTableIcebergMergeOnReadBehaviorEnabled)).
						HasIcebergMergeOnReadBehaviorLevel(sdk.ParameterTypeTable),
				),
			},
			// Import - complete (path_layout is not returned by SHOW/DESCRIBE)
			{
				Config:            accconfig.FromModels(t, modelComplete),
				ResourceName:      modelComplete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"path_layout",
					"auto_refresh",
					"catalog_namespace",
				},
			},
		},
	})
}

func TestAcc_IcebergTableFromRest_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	baseModel := func() *model.IcebergTableFromRestModel {
		return model.IcebergTableFromRestWithDefaultMeta(
			id.DatabaseName(),
			id.SchemaName(),
			id.Name(),
			"TEST",
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromRest),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, baseModel().WithCatalogTableName("")),
				ExpectError: regexp.MustCompile(`expected "catalog_table_name" to not be an empty string`),
			},
			{
				Config:      accconfig.FromModels(t, baseModel().WithPathLayout("INVALID")),
				ExpectError: regexp.MustCompile(`expected .* to be one of .*, got INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, baseModel().WithAutoRefresh("INVALID")),
				ExpectError: regexp.MustCompile(`expected .* to be one of .*, got INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, baseModel().WithTargetFileSize("INVALID")),
				ExpectError: regexp.MustCompile(`expected .* to be one of .*, got INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, baseModel().WithStorageSerializationPolicy("INVALID")),
				ExpectError: regexp.MustCompile(`expected .* to be one of .*, got INVALID`),
			},
			{
				Config:      accconfig.FromModels(t, baseModel().WithIcebergMergeOnReadBehavior("INVALID")),
				ExpectError: regexp.MustCompile(`expected .* to be one of .*, got INVALID`),
			},
		},
	})
}
