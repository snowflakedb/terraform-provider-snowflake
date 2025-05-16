//go:build !account_level_tests

package datasources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO(SNOW-1548063): 090105 (22000): Cannot perform operation. This session does not have a current database. Call 'USE DATABASE', or use a qualified name.
func TestAcc_Streamlits(t *testing.T) {
	t.Skip("Skipping because of the error: 090105 (22000): Cannot perform operation. This session does not have a current database. Call 'USE DATABASE', or use a qualified name.")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	// warehouse is needed because default warehouse uses lowercase, and it fails in snowflake.
	// TODO(SNOW-1541938): use a default warehouse after fix on snowflake side
	warehouse, warehouseCleanup := acc.TestClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	networkRule, networkRuleCleanup := acc.TestClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)
	externalAccessIntegrationId, externalAccessIntegrationCleanup := acc.TestClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	databaseId := acc.TestClient().Ids.DatabaseId()
	schemaId := acc.TestClient().Ids.SchemaId()
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
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
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, streamlitModel, streamlitsModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.#", "1"),

					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.#", "1"),
					resource.TestCheckResourceAttrSet(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.created_on"),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.database_name", databaseId.Name()),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.schema_name", schemaId.Name()),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.title", title),
					resource.TestCheckResourceAttr(streamlitsModel.DatasourceReference(), "streamlits.0.show_output.0.owner", acc.TestClient().Context.CurrentRole(t).Name()),
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
				Config: accconfig.FromModels(t, streamlitModel, streamlitsModelWithoutDescribe),
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

func TestAcc_Streamlits_Filtering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stage, stageCleanup := acc.TestClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	mainFile := random.AlphaN(4)

	streamlitModel1 := model.StreamlitWithIds("test1", idOne, mainFile, stage.ID())
	streamlitModel2 := model.StreamlitWithIds("test2", idTwo, mainFile, stage.ID())
	streamlitModel3 := model.StreamlitWithIds("test3", idThree, mainFile, stage.ID())
	streamlitsModelLikeFirst := datasourcemodel.Streamlits("test").
		WithLike(idOne.Name()).
		WithInDatabase(idOne.DatabaseId()).
		WithDependsOn(streamlitModel1.ResourceReference(), streamlitModel2.ResourceReference(), streamlitModel3.ResourceReference())
	streamlitsModelLikePrefix := datasourcemodel.Streamlits("test").
		WithLike(prefix+"%").
		WithInDatabase(idOne.DatabaseId()).
		WithDependsOn(streamlitModel1.ResourceReference(), streamlitModel2.ResourceReference(), streamlitModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Streamlit),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, streamlitModel1, streamlitModel2, streamlitModel3, streamlitsModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitsModelLikeFirst.DatasourceReference(), "streamlits.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, streamlitModel1, streamlitModel2, streamlitModel3, streamlitsModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(streamlitsModelLikePrefix.DatasourceReference(), "streamlits.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Streamlits_badCombination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
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
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streamlits/non_existing"),
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
`, acc.TestDatabaseName, acc.TestSchemaName)
}

func streamlitsDatasourceEmptyIn() string {
	return `
data "snowflake_streamlits" "test" {
  in {
  }
}
`
}
