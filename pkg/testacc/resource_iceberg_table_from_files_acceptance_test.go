//go:build non_account_level_tests

package testacc

import (
	"strings"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_IcebergTableFromFiles_BasicUseCase(t *testing.T) {
	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	baseLocationPrefix := "iceberg_test_table"
	metadataFilePath := baseLocationPrefix + "/metadata/v1.metadata.json"

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	externalVolumeId2, externalVolumeCleanup2 := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup2)

	// Create a dedicated database with external_volume and catalog set at db level so the table
	// can be created without specifying them explicitly (matching the "required fields only" test case).
	dbForIcebergFiles, dbCleanup := testClient().Database.CreateDatabaseWithRequest(t, sdk.NewCreateDatabaseRequest(testClient().Ids.RandomAccountObjectIdentifier()).WithCatalog(catalogId).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbCleanup)
	schemaIdForIcebergFiles := sdk.NewDatabaseObjectIdentifier(dbForIcebergFiles.ID().Name(), "PUBLIC")

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForIcebergFiles)
	comment := random.Comment()
	externalComment := random.Comment()

	// modelBasic relies on db-level external_volume and catalog defaults — no explicit values.
	modelBasic := model.IcebergTableFromFilesWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		metadataFilePath,
	)

	// modelWithAllOptional also relies on db-level defaults for external_volume/catalog.
	modelWithAllOptional := model.IcebergTableFromFilesWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		metadataFilePath,
	).WithComment(comment).
		WithReplaceInvalidCharacters(true)

	ref := modelBasic.ResourceReference()

	// external_volume and catalog are inherited from the database (DATABASE level).
	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableFromFilesResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasMetadataFilePathString(metadataFilePath).
			HasExternalVolume(externalVolumeId.Name()).
			HasCatalog(catalogId.Name()).
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
			HasCatalogTableName("").
			HasCatalogNamespace("").
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
			HasReplaceInvalidCharactersLevel(sdk.ParameterTypeSnowflakeDefault),
	}

	// replace_invalid_characters is set explicitly at table level; external_volume/catalog stay at db level.
	allOptionalAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableFromFilesResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasMetadataFilePathString(metadataFilePath).
			HasExternalVolume(externalVolumeId.Name()).
			HasCatalog(catalogId.Name()).
			HasComment(comment).
			HasReplaceInvalidCharacters(true).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.IcebergTableShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasExternalVolumeName(externalVolumeId).
			HasCatalogName(catalogId).
			HasIcebergTableType(sdk.IcebergTableTypeUnmanaged).
			HasCatalogTableName("").
			HasCatalogNamespace("").
			HasCanWriteMetadata(true).
			HasComment(comment).
			HasNameMapping("").
			HasOwnerRoleType("ROLE").
			HasCatalogSyncName("").
			HasAutoRefreshStatusEmpty(),
		resourceparametersassert.IcebergTableResourceParameters(t, ref).
			HasExternalVolume(externalVolumeId.Name()).
			HasExternalVolumeLevel(sdk.ParameterTypeDatabase).
			HasCatalog(catalogId.Name()).
			HasCatalogLevel(sdk.ParameterTypeDatabase).
			HasReplaceInvalidCharacters(true).
			HasReplaceInvalidCharactersLevel(sdk.ParameterTypeTable),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromFiles),
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
					"metadata_file_path",
				},
			},
			// Set all optional fields
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
					"metadata_file_path",
				},
			},
			// Change alterable fields externally and detect drift
			{
				PreConfig: func() {
					testClient().IcebergTable.Alter(t, sdk.NewAlterIcebergTableRequest(id).WithSet(
						*sdk.NewIcebergTableSetPropertiesRequest().
							WithReplaceInvalidCharacters(false).
							WithComment(externalComment),
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
			// Change force new fields externally and detect drift
			{
				PreConfig: func() {
					testClient().IcebergTable.CreateFromIcebergFiles(
						t, id,
						sdk.NewCreateFromIcebergFilesIcebergTableRequest(id, metadataFilePath).
							WithOrReplace(true).
							WithComment(externalComment).
							WithExternalVolume(externalVolumeId2),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
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

func TestAcc_IcebergTableFromFiles_CompleteUseCase(t *testing.T) {
	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	baseLocationPrefix := "iceberg_test_table"
	metadataFilePath := baseLocationPrefix + "/metadata/v1.metadata.json"

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.IcebergTableFromFilesWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		metadataFilePath,
	).WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithComment(comment).
		WithReplaceInvalidCharacters(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromFiles),
		Steps: []resource.TestStep{
			// Create with all fields
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(
					t,
					resourceassert.IcebergTableFromFilesResource(t, modelComplete.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasMetadataFilePathString(metadataFilePath).
						HasExternalVolume(externalVolumeId.FullyQualifiedName()).
						HasCatalog(catalogId.FullyQualifiedName()).
						HasComment(comment).
						HasReplaceInvalidCharacters(true).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.IcebergTableShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()),
					resourceparametersassert.IcebergTableResourceParameters(t, modelComplete.ResourceReference()).
						HasExternalVolume(externalVolumeId.FullyQualifiedName()).
						HasCatalog(catalogId.FullyQualifiedName()).
						HasReplaceInvalidCharacters(true),
				),
			},
			// Import - complete
			{
				Config:            accconfig.FromModels(t, modelComplete),
				ResourceName:      modelComplete.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"metadata_file_path",
				},
			},
		},
	})
}

func TestAcc_IcebergTableFromFiles_Import(t *testing.T) {
	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	baseLocationPrefix := "iceberg_test_table"
	metadataFilePath := baseLocationPrefix + "/metadata/v1.metadata.json"

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := testClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.IcebergTableFromFilesWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		metadataFilePath,
	).WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithComment(comment).
		WithReplaceInvalidCharacters(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromFiles),
		Steps: []resource.TestStep{
			// Import the externally created resource
			{
				PreConfig: func() {
					testClient().IcebergTable.CreateFromIcebergFiles(t, id,
						sdk.NewCreateFromIcebergFilesIcebergTableRequest(id, metadataFilePath).
							WithComment(comment).
							WithExternalVolume(externalVolumeId).
							WithCatalog(catalogId).
							WithReplaceInvalidCharacters(true))
				},
				Config:             config.FromModels(t, modelComplete),
				ResourceName:       modelComplete.ResourceReference(),
				ImportState:        true,
				ImportStateId:      id.FullyQualifiedName(),
				ImportStatePersist: true,
			},
			{
				Config: accconfig.FromModels(t, modelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
