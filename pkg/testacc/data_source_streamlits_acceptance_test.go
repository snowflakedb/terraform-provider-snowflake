//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Streamlits_BasicUseCase_DifferentFiltering(t *testing.T) {
	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	prefix := random.AlphaN(4)
	streamlitId1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	streamlitId2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	streamlitId3 := testClient().Ids.RandomSchemaObjectIdentifier()

	mainFile := random.AlphaN(4)

	streamlitModel1 := model.StreamlitWithIds("test", streamlitId1, mainFile, stage.ID())
	streamlitModel2 := model.StreamlitWithIds("test1", streamlitId2, mainFile, stage.ID())
	streamlitModel3 := model.StreamlitWithIds("test2", streamlitId3, mainFile, stage.ID())

	datasourceModelLikeExact := datasourcemodel.Streamlits("test").
		WithLike(streamlitId1.Name()).
		WithDependsOn(streamlitModel1.ResourceReference(), streamlitModel2.ResourceReference(), streamlitModel3.ResourceReference())

	datasourceModelLikePrefix := datasourcemodel.Streamlits("test").
		WithLike(prefix+"%").
		WithDependsOn(streamlitModel1.ResourceReference(), streamlitModel2.ResourceReference(), streamlitModel3.ResourceReference())

	datasourceModelInDatabase := datasourcemodel.Streamlits("test").
		WithInDatabase(streamlitId1.DatabaseId()).
		WithDependsOn(streamlitModel1.ResourceReference(), streamlitModel2.ResourceReference(), streamlitModel3.ResourceReference())

	datasourceModelInSchema := datasourcemodel.Streamlits("test").
		WithInSchema(streamlitId1.SchemaId()).
		WithDependsOn(streamlitModel1.ResourceReference(), streamlitModel2.ResourceReference(), streamlitModel3.ResourceReference())

	datasourceModelLikeAndInDatabase := datasourcemodel.Streamlits("test").
		WithLike(prefix+"%").
		WithInDatabase(streamlitId1.DatabaseId()).
		WithDependsOn(streamlitModel1.ResourceReference(), streamlitModel2.ResourceReference(), streamlitModel3.ResourceReference())

	datasourceModelLikeAndInSchema := datasourcemodel.Streamlits("test").
		WithLike(prefix+"%").
		WithInSchema(streamlitId1.SchemaId()).
		WithDependsOn(streamlitModel1.ResourceReference(), streamlitModel2.ResourceReference(), streamlitModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, streamlitModel1, streamlitModel2, streamlitModel3, datasourceModelLikeExact),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "streamlits.#", "1"),
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "streamlits.0.show_output.0.name", streamlitId1.Name()),
				),
			},
			{
				Config: config.FromModels(t, streamlitModel1, streamlitModel2, streamlitModel3, datasourceModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikePrefix.DatasourceReference(), "streamlits.#", "2"),
				),
			},
			{
				Config: config.FromModels(t, streamlitModel1, streamlitModel2, streamlitModel3, datasourceModelInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelInDatabase.DatasourceReference(), "streamlits.#", "3"),
				),
			},
			{
				Config: config.FromModels(t, streamlitModel1, streamlitModel2, streamlitModel3, datasourceModelInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelInSchema.DatasourceReference(), "streamlits.#", "3"),
				),
			},
			{
				Config: config.FromModels(t, streamlitModel1, streamlitModel2, streamlitModel3, datasourceModelLikeAndInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeAndInDatabase.DatasourceReference(), "streamlits.#", "2"),
				),
			},
			{
				Config: config.FromModels(t, streamlitModel1, streamlitModel2, streamlitModel3, datasourceModelLikeAndInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeAndInSchema.DatasourceReference(), "streamlits.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Streamlits_CompleteUseCase(t *testing.T) {
	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	streamlitId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	title := random.AlphaN(4)
	mainFile := random.AlphaN(4)
	directoryLocation := random.AlphaN(4)

	streamlitModel := model.StreamlitWithIds("test", streamlitId, mainFile, stage.ID()).
		WithComment(comment).
		WithTitle(title).
		WithQueryWarehouse(warehouse.ID().Name()).
		WithDirectoryLocation(directoryLocation)

	datasourceModelWithoutDescribe := datasourcemodel.Streamlits("test").
		WithLike(streamlitId.Name()).
		WithInDatabase(streamlitId.DatabaseId()).
		WithWithDescribe(false).
		WithDependsOn(streamlitModel.ResourceReference())

	datasourceModelWithDescribe := datasourcemodel.Streamlits("test").
		WithLike(streamlitId.Name()).
		WithInDatabase(streamlitId.DatabaseId()).
		WithWithDescribe(true).
		WithDependsOn(streamlitModel.ResourceReference())

	commonShowOutputAssert := func(t *testing.T, datasourceReference string) *resourceshowoutputassert.StreamlitShowOutputAssert {
		return resourceshowoutputassert.StreamlitsDatasourceShowOutput(t, datasourceReference).
			HasCreatedOnNotEmpty().
			HasName(streamlitId.Name()).
			HasDatabaseName(streamlitId.DatabaseName()).
			HasSchemaName(streamlitId.SchemaName()).
			HasTitle(title).
			HasComment(comment).
			HasQueryWarehouse(warehouse.ID().Name()).
			HasUrlIdNotEmpty().
			HasOwnerNotEmpty().
			HasOwnerRoleType("ROLE")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, streamlitModel, datasourceModelWithoutDescribe),
				Check: assertThat(t,
					commonShowOutputAssert(t, datasourceModelWithoutDescribe.DatasourceReference()),

					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutDescribe.DatasourceReference(), "streamlits.0.describe_output.#", "0")),
				),
			},
			{
				Config: config.FromModels(t, streamlitModel, datasourceModelWithDescribe),
				Check: assertThat(t,
					commonShowOutputAssert(t, datasourceModelWithDescribe.DatasourceReference()),

					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "streamlits.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "streamlits.0.describe_output.0.name", streamlitId.Name())),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "streamlits.0.describe_output.0.title", title)),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "streamlits.0.describe_output.0.main_file", mainFile)),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDescribe.DatasourceReference(), "streamlits.0.describe_output.0.query_warehouse", warehouse.ID().Name())),
					assert.Check(resource.TestCheckResourceAttrSet(datasourceModelWithDescribe.DatasourceReference(), "streamlits.0.describe_output.0.url_id")),
					assert.Check(resource.TestCheckResourceAttrSet(datasourceModelWithDescribe.DatasourceReference(), "streamlits.0.describe_output.0.root_location")),
					assert.Check(resource.TestCheckResourceAttrSet(datasourceModelWithDescribe.DatasourceReference(), "streamlits.0.describe_output.0.default_packages")),
					assert.Check(resource.TestCheckResourceAttrSet(datasourceModelWithDescribe.DatasourceReference(), "streamlits.0.describe_output.0.user_packages.#")),
					assert.Check(resource.TestCheckResourceAttrSet(datasourceModelWithDescribe.DatasourceReference(), "streamlits.0.describe_output.0.import_urls.#")),
				),
			},
		},
	})
}

// TODO(SNOW-1548063): 090105 (22000): Cannot perform operation. This session does not have a current database. Call 'USE DATABASE', or use a qualified name.
func TestAcc_Streamlits(t *testing.T) {
	t.Skip("Skipping because of the error: 090105 (22000): Cannot perform operation. This session does not have a current database. Call 'USE DATABASE', or use a qualified name.")

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	// warehouse is needed because default warehouse uses lowercase, and it fails in snowflake.
	// TODO(SNOW-1541938): use a default warehouse after fix on snowflake side
	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	networkRule, networkRuleCleanup := testClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)
	externalAccessIntegrationId, externalAccessIntegrationCleanup := testClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	databaseId := testClient().Ids.DatabaseId()
	schemaId := testClient().Ids.SchemaId()
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	mainFile := random.AlphaN(4)
	comment := random.Comment()
	title := random.AlphaN(4)
	directoryLocation := random.AlphaN(4)
	rootLocation := fmt.Sprintf("@%s/%s", stage.ID().FullyQualifiedName(), directoryLocation)

	streamlitModel := model.StreamlitWithIds("test", id, mainFile, stage.ID()).
		WithComment(comment).
		WithTitle(title).
		WithDirectoryLocation(directoryLocation).
		WithQueryWarehouse(warehouse.ID().Name()).
		WithExternalAccessIntegrations(externalAccessIntegrationId)
	streamlitsModel := datasourcemodel.Streamlits("test").
		WithLike(id.Name()).
		WithDependsOn(streamlitModel.ResourceReference())
	streamlitsModelWithoutDescribe := datasourcemodel.Streamlits("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(streamlitModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, streamlitModel, streamlitsModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.#", "1"),

					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.#", "1"),
					resource.TestCheckResourceAttrSet(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.schema_name", schemaId.Name()),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.title", title),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.owner", testClient().Context.CurrentRole(t).Name()),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.url_id"),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.owner_role_type", "ROLE"),

					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.title", title),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.root_location", rootLocation),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.main_file", mainFile),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.url_id"),
					resource.TestCheckResourceAttrSet(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.default_packages"),
					resource.TestCheckResourceAttrSet(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.user_packages.#"),
					resource.TestCheckResourceAttrSet(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.import_urls.#"),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.external_access_integrations.#", "1"),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.external_access_integrations.0", externalAccessIntegrationId.Name()),
					resource.TestCheckResourceAttrSet(streamlitsModel.DatasourceReference(), "streamlits.0.describe_output.0.external_access_secrets"),
				),
			},
			{
				Config: config.FromModels(t, streamlitModel, streamlitsModelWithoutDescribe),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.#", "1"),

					resource.TestCheckResourceAttr(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.show_output.#", "1"),
					resource.TestCheckResourceAttrSet(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttr(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.show_output.0.schema_name", schemaId.Name()),
					resource.TestCheckResourceAttr(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.show_output.0.title", title),
					resource.TestCheckResourceAttr(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.show_output.0.owner", snowflakeroles.Accountadmin.Name()),
					resource.TestCheckResourceAttr(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.show_output.0.comment", comment),
					resource.TestCheckResourceAttr(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.show_output.0.query_warehouse", warehouse.ID().Name()),
					resource.TestCheckResourceAttrSet(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.show_output.0.url_id"),
					resource.TestCheckResourceAttr(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.show_output.0.owner_role_type", "ROLE"),

					resource.TestCheckResourceAttr(streamlitsModelWithoutDescribe.DatasourceReference(), "streamlits.0.describe_output.#", "0"),
				),
			},
		},
	})
}

func TestAcc_Streamlits_badCombination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      streamlitsDatasourceConfigDbAndSchema(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Streamlits_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      streamlitsDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Streamlits_StreamlitNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Streamlits/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one streamlit"),
			},
		},
	})
}

func streamlitsDatasourceConfigDbAndSchema() string {
	return fmt.Sprintf(`
data "snowflake_streamlits" "test" {
  in {
    database = "%s"
    schema   = "%s"
  }
}
`, TestDatabaseName, TestSchemaName)
}

func streamlitsDatasourceEmptyIn() string {
	return `
data "snowflake_streamlits" "test" {
  in {
  }
}
`
}
