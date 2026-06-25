//go:build non_account_level_tests

package testacc

import (
	"strings"
	"testing"

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

func deltaLakeCatalog(t *testing.T) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	return testClient().CatalogIntegration.CreateFunc(t,
		sdk.NewCreateCatalogIntegrationRequest(testClient().Ids.RandomAccountObjectIdentifier(), true).
			WithObjectStorageCatalogSourceParams(*sdk.NewObjectStorageParamsRequest(sdk.CatalogIntegrationTableFormatDelta)),
	)
}

func TestAcc_IcebergTableFromDeltaFiles_BasicUseCase(t *testing.T) {
	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	baseLocation := "delta_lake_test_table/"

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := deltaLakeCatalog(t)
	t.Cleanup(catalogCleanup)

	externalVolumeId2, externalVolumeCleanup2 := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup2)

	catalogId2, catalogCleanup2 := deltaLakeCatalog(t)
	t.Cleanup(catalogCleanup2)

	// Create a dedicated database with external_volume and catalog set at db level so the table
	// can be created without specifying them explicitly (matching the "required fields only" test case).
	dbForDeltaFiles, dbCleanup := testClient().Database.CreateDatabaseWithRequest(t, sdk.NewCreateDatabaseRequest(testClient().Ids.RandomAccountObjectIdentifier()).
		WithCatalog(catalogId).
		WithExternalVolume(externalVolumeId))
	t.Cleanup(dbCleanup)
	schemaIdForDeltaFiles := sdk.NewDatabaseObjectIdentifier(dbForDeltaFiles.ID().Name(), "PUBLIC")

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForDeltaFiles)
	comment := random.Comment()
	externalComment := random.Comment()

	// modelBasic relies on db-level external_volume and catalog defaults — no explicit values.
	modelBasic := model.IcebergTableFromDeltaFilesWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		baseLocation,
	)

	// modelWithAllOptional also relies on db-level defaults for external_volume/catalog.
	modelWithAllOptional := model.IcebergTableFromDeltaFilesWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		baseLocation,
	).WithComment(comment).
		WithReplaceInvalidCharacters(true).
		WithAutoRefresh("true")

	ref := modelBasic.ResourceReference()

	// external_volume and catalog are inherited from the database (DATABASE level).
	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableFromDeltaFilesResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasBaseLocationString(baseLocation).
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

	// replace_invalid_characters and auto_refresh are set explicitly at table level; external_volume/catalog stay at db level.
	allOptionalAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableFromDeltaFilesResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasBaseLocationString(baseLocation).
			HasExternalVolume(externalVolumeId.Name()).
			HasCatalog(catalogId.Name()).
			HasComment(comment).
			HasReplaceInvalidCharacters(true).
			HasAutoRefreshString("true").
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
			HasAutoRefreshStatus(sdk.IcebergTableAutoRefreshStatus{
				CurrentSnapshotId:    0,
				PendingSnapshotCount: 0,
				ExecutionState:       "RUNNING",
			}),
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
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromDeltaFiles),
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
					"auto_refresh",
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
					"auto_refresh",
				},
			},
			// Change alterable fields externally and detect drift
			{
				PreConfig: func() {
					testClient().IcebergTable.CreateFromDeltaLake(t, id,
						sdk.NewCreateFromDeltaLakeIcebergTableRequest(id, baseLocation).
							WithOrReplace(true).
							WithComment(externalComment).
							WithExternalVolume(externalVolumeId2).
							WithCatalog(catalogId2).
							WithReplaceInvalidCharacters(false).
							WithAutoRefresh(false),
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

func TestAcc_IcebergTableFromDeltaFiles_CompleteUseCase(t *testing.T) {
	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	// Here, the trailing slash is missing on purpose to test the diff suppression.
	baseLocation := "delta_lake_test_table"

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := deltaLakeCatalog(t)
	t.Cleanup(catalogCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.IcebergTableFromDeltaFilesWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		baseLocation,
	).WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithComment(comment).
		WithReplaceInvalidCharacters(true).
		WithAutoRefresh("false")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromDeltaFiles),
		Steps: []resource.TestStep{
			// Create with all fields
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.IcebergTableFromDeltaFilesResource(t, modelComplete.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasBaseLocationString(baseLocation+"/").
						HasExternalVolume(externalVolumeId.FullyQualifiedName()).
						HasCatalog(catalogId.FullyQualifiedName()).
						HasComment(comment).
						HasReplaceInvalidCharacters(true).
						HasAutoRefreshString("false").
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
					"auto_refresh",
				},
			},
		},
	})
}

func TestAcc_IcebergTableFromDeltaFiles_Import(t *testing.T) {
	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"
	baseLocation := "delta_lake_test_table"

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := deltaLakeCatalog(t)
	t.Cleanup(catalogCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.IcebergTableFromDeltaFilesWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		baseLocation,
	).WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithComment(comment).
		WithAutoRefresh("false").
		WithReplaceInvalidCharacters(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromDeltaFiles),
		Steps: []resource.TestStep{
			// Import the externally created resource
			{
				PreConfig: func() {
					testClient().IcebergTable.CreateFromDeltaLake(t, id,
						sdk.NewCreateFromDeltaLakeIcebergTableRequest(id, baseLocation).
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
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionNoop),
						planchecks.PrintPlanDetails(modelComplete.ResourceReference(), "auto_refresh", "replace_invalid_characters", "comment"),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
