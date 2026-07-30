//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_McpServers_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	cleanup := testClient().McpServer.Create(t, id)
	t.Cleanup(cleanup)

	datasourceModel := datasourcemodel.McpServers("test").
		WithLike(id.Name()).
		WithInSchema(id.SchemaId())
	datasourceModelWithoutDescribe := datasourcemodel.McpServers("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInSchema(id.SchemaId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, datasourceModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModel.DatasourceReference(), "mcp_servers.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModel.DatasourceReference(), "mcp_servers.0.describe_output.#", "1")),
					resourceshowoutputassert.McpServersDatasourceShowOutput(t, datasourceModel.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(""),
					resourceshowoutputassert.McpServersDatasourceDescribeOutput(t, datasourceModel.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(""),
				),
			},
			{
				Config: accconfig.FromModels(t, datasourceModelWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutDescribe.DatasourceReference(), "mcp_servers.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutDescribe.DatasourceReference(), "mcp_servers.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_McpServers_Filtering(t *testing.T) {
	prefix := testClient().Ids.Alpha()
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)

	cleanup1 := testClient().McpServer.Create(t, id1)
	t.Cleanup(cleanup1)
	cleanup2 := testClient().McpServer.Create(t, id2)
	t.Cleanup(cleanup2)

	modelByPrefix := datasourcemodel.McpServers("by_prefix").
		WithLike(prefix + "%").
		WithInSchema(id1.SchemaId())
	modelById1 := datasourcemodel.McpServers("by_id1").
		WithLike(id1.Name()).
		WithInSchema(id1.SchemaId())
	modelBySchema := datasourcemodel.McpServers("by_schema").
		WithInSchema(id1.SchemaId()).
		WithLike(prefix + "%")
	modelInDatabase := datasourcemodel.McpServers("in_database").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId())
	modelNoResults := datasourcemodel.McpServers("no_results").
		WithLike("non_existent_%").
		WithInSchema(id1.SchemaId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelByPrefix, modelById1, modelBySchema, modelInDatabase, modelNoResults),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(modelByPrefix.DatasourceReference(), "mcp_servers.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(modelById1.DatasourceReference(), "mcp_servers.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelBySchema.DatasourceReference(), "mcp_servers.#", "2")),
					assert.Check(resource.TestCheckResourceAttr(modelInDatabase.DatasourceReference(), "mcp_servers.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(modelNoResults.DatasourceReference(), "mcp_servers.#", "0")),
				),
			},
		},
	})
}

func TestAcc_McpServers_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.McpServers("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}
