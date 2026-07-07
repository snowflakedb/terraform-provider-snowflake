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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_IcebergTableFromAwsGlue_BasicUseCase(t *testing.T) {
	// TODO(SNOW-3725859): Provide the external volume and catalog integration dynamically. Unskip and
	// fold these tests into the main suite.
	t.Skip("Iceberg AWS Glue tests require preconfigured external catalog integrations and are not run by default")
	const (
		glueCatalogName = "GLUE_CATALOG_INTEGRATION"
		glueVolumeName  = "GLUE_EXTERNAL_VOLUME"
		// Values that must match the manually preconfigured AWS Glue catalog contents.
		glueCatalogTableName = "TEST"
		glueCatalogNamespace = "glue_iceberg_schema"
	)
	externalVolumeId := sdk.NewAccountObjectIdentifier(glueVolumeName)
	catalogId := sdk.NewAccountObjectIdentifier(glueCatalogName)

	// Create a dedicated database with external_volume and catalog set at db level so the table
	// can be created without specifying them explicitly (matching the "required fields only" test case).
	dbForIcebergGlue, dbCleanup := testClient().Database.CreateDatabaseWithRequest(t, testClient().Database.TestParametersSet(testClient().Ids.RandomAccountObjectIdentifier()).WithCatalog(catalogId).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbCleanup)
	schemaIdForIcebergGlue := sdk.NewDatabaseObjectIdentifier(dbForIcebergGlue.ID().Name(), "PUBLIC")

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaIdForIcebergGlue)
	comment := random.Comment()
	externalComment := random.Comment()

	// modelBasic relies on db-level external_volume and catalog defaults — no explicit values.
	modelBasic := model.IcebergTableFromAwsGlueWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		glueCatalogTableName,
	)

	// modelWithAllOptional only sets alterable fields so the transition from modelBasic is an update
	// (not a force-new recreate). The force-new fields are covered by the complete use case.
	modelWithAllOptional := model.IcebergTableFromAwsGlueWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		glueCatalogTableName,
	).WithComment(comment).
		WithReplaceInvalidCharacters(true).
		WithAutoRefresh("true")

	ref := modelBasic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableFromAwsGlueResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasCatalogTableNameString(glueCatalogTableName).
			HasExternalVolumeString(externalVolumeId.Name()).
			HasCatalogString(catalogId.Name()).
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
			HasCatalogTableName(glueCatalogTableName).
			HasCatalogNamespace(glueCatalogNamespace).
			HasComment(""),
		resourceparametersassert.IcebergTableResourceParameters(t, ref).
			HasExternalVolume(externalVolumeId.Name()).
			HasExternalVolumeLevel(sdk.ParameterTypeDatabase).
			HasCatalog(catalogId.Name()).
			HasCatalogLevel(sdk.ParameterTypeDatabase).
			HasReplaceInvalidCharacters(false),
	}

	allOptionalAssertions := []assert.TestCheckFuncProvider{
		resourceassert.IcebergTableFromAwsGlueResource(t, ref).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNameString(id.Name()).
			HasCatalogTableNameString(glueCatalogTableName).
			HasExternalVolumeString(externalVolumeId.Name()).
			HasCatalogString(catalogId.Name()).
			HasComment(comment).
			HasReplaceInvalidCharacters(true).
			HasAutoRefreshString("true").
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.IcebergTableShowOutput(t, ref).
			HasAutoRefreshStatus(sdk.IcebergTableAutoRefreshStatus{
				CurrentSnapshotId:    0,
				PendingSnapshotCount: 0,
				ExecutionState:       "RUNNING",
			}),
		resourceparametersassert.IcebergTableResourceParameters(t, ref).
			HasExternalVolume(externalVolumeId.Name()).
			HasCatalog(catalogId.Name()).
			HasReplaceInvalidCharacters(true).
			HasReplaceInvalidCharactersLevel(sdk.ParameterTypeTable),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromAwsGlue),
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
					"catalog_namespace",
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
			// Change alterable fields externally and detect drift
			{
				PreConfig: func() {
					testClient().IcebergTable.Alter(t, sdk.NewAlterIcebergTableRequest(id).WithSet(
						*sdk.NewIcebergTableSetPropertiesRequest().
							WithReplaceInvalidCharacters(false).
							WithAutoRefresh(false).
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
						planchecks.ExpectDrift(ref, "auto_refresh", new("true"), new("false")),
						planchecks.ExpectChange(ref, "auto_refresh", tfjson.ActionUpdate, new("false"), new("true")),
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

func TestAcc_IcebergTableFromAwsGlue_CompleteUseCase(t *testing.T) {
	// TODO(SNOW-3725859): Provide the external volume and catalog integration dynamically. Unskip and
	// fold these tests into the main suite.
	t.Skip("Iceberg AWS Glue tests require preconfigured external catalog integrations and are not run by default")
	const (
		glueCatalogName = "GLUE_CATALOG_INTEGRATION"
		glueVolumeName  = "GLUE_EXTERNAL_VOLUME"
		// Values that must match the manually preconfigured AWS Glue catalog contents.
		glueCatalogTableName = "TEST"
		glueCatalogNamespace = "glue_iceberg_schema"
	)
	externalVolumeId := sdk.NewAccountObjectIdentifier(glueVolumeName)
	catalogId := sdk.NewAccountObjectIdentifier(glueCatalogName)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	modelComplete := model.IcebergTableFromAwsGlueWithDefaultMeta(
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		glueCatalogTableName,
	).WithExternalVolume(externalVolumeId.Name()).
		WithCatalog(catalogId.Name()).
		WithCatalogNamespace(glueCatalogNamespace).
		WithAutoRefresh("true").
		WithReplaceInvalidCharacters(true).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromAwsGlue),
		Steps: []resource.TestStep{
			// Create with all fields
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(
					t,
					resourceassert.IcebergTableFromAwsGlueResource(t, modelComplete.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasCatalogTableNameString(glueCatalogTableName).
						HasCatalogNamespaceString(glueCatalogNamespace).
						HasExternalVolumeString(externalVolumeId.FullyQualifiedName()).
						HasCatalogString(catalogId.FullyQualifiedName()).
						HasAutoRefreshString("true").
						HasComment(comment).
						HasReplaceInvalidCharacters(true).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.IcebergTableShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCatalogTableName(glueCatalogTableName).
						HasCatalogNamespace(glueCatalogNamespace),
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

func TestAcc_IcebergTableFromAwsGlue_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	baseModel := func() *model.IcebergTableFromAwsGlueModel {
		return model.IcebergTableFromAwsGlueWithDefaultMeta(
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
		CheckDestroy: CheckDestroy(t, resources.IcebergTableFromAwsGlue),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, baseModel().WithCatalogTableName("")),
				ExpectError: regexp.MustCompile(`expected "catalog_table_name" to not be an empty string`),
			},
			{
				Config:      accconfig.FromModels(t, baseModel().WithAutoRefresh("INVALID")),
				ExpectError: regexp.MustCompile(`expected .* to be one of .*, got INVALID`),
			},
		},
	})
}
