//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_CortexAgents_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	response := "You are a helpful assistant"
	hclSpec := model.SampleSpecAsYamlencodeHCL(response)
	normalizedSpec, err := sdk.NormalizeCortexAgentSpecification(
		testClient().CortexAgent.SampleSpecWithResponse(t, response))
	require.NoError(t, err)
	comment := random.Comment()
	completeProfile := sdk.CortexAgentProfile{
		DisplayName: sdk.String("My Helpful Assistant"),
		Avatar:      sdk.String("business-icon.png"),
		Color:       sdk.String("red"),
	}

	completeModel := model.CortexAgentWithSpecification("test", id.DatabaseName(), id.SchemaName(), id.Name(), hclSpec).
		WithComment(comment).
		WithProfile(completeProfile)

	cortexAgentsModel := datasourcemodel.CortexAgents("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(completeModel.ResourceReference())

	cortexAgentsModelWithoutDescribe := datasourcemodel.CortexAgents("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(completeModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, completeModel, cortexAgentsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(cortexAgentsModel.DatasourceReference(), "cortex_agents.#", "1")),
					resourceshowoutputassert.CortexAgentsDatasourceShowOutput(t, "snowflake_cortex_agents.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasProfile(completeProfile),
					resourceshowoutputassert.CortexAgentsDatasourceDescribeOutput(t, "snowflake_cortex_agents.test").
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasProfile(completeProfile).
						HasAgentSpec(normalizedSpec).
						HasCreatedOnNotEmpty().
						HasDefaultVersionName("LAST").
						HasVersions(`["VERSION$1"]`).
						HasAliases(`{"DEFAULT":"VERSION$1","FIRST":"VERSION$1","LAST":"VERSION$1"}`),
				),
			},
			{
				Config: accconfig.FromModels(t, completeModel, cortexAgentsModelWithoutDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(cortexAgentsModelWithoutDescribe.DatasourceReference(), "cortex_agents.#", "1")),
					resourceshowoutputassert.CortexAgentsDatasourceShowOutput(t, "snowflake_cortex_agents.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment).
						HasProfile(completeProfile),
					assert.Check(resource.TestCheckResourceAttr(cortexAgentsModelWithoutDescribe.DatasourceReference(), "cortex_agents.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_CortexAgents_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchema(secondSchema.ID())

	response := "You are a helpful assistant"
	hclSpec := model.SampleSpecAsYamlencodeHCL(response)

	model1 := model.CortexAgentWithSpecification("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), hclSpec)
	model2 := model.CortexAgentWithSpecification("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), hclSpec)
	model3 := model.CortexAgentWithSpecification("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name(), hclSpec)

	cortexAgentsModelLikeFirst := datasourcemodel.CortexAgents("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	cortexAgentsModelLikePrefix := datasourcemodel.CortexAgents("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	cortexAgentsModelStartsWith := datasourcemodel.CortexAgents("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	cortexAgentsModelLimit := datasourcemodel.CortexAgents("test").
		WithRowsAndFrom(1, prefix).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	cortexAgentsModelInSchema := datasourcemodel.CortexAgents("test").
		WithInSchema(id1.SchemaId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, cortexAgentsModelLikeFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cortexAgentsModelLikeFirst.DatasourceReference(), "cortex_agents.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, cortexAgentsModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cortexAgentsModelLikePrefix.DatasourceReference(), "cortex_agents.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, cortexAgentsModelStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cortexAgentsModelStartsWith.DatasourceReference(), "cortex_agents.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, cortexAgentsModelLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cortexAgentsModelLimit.DatasourceReference(), "cortex_agents.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, cortexAgentsModelInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cortexAgentsModelInSchema.DatasourceReference(), "cortex_agents.#", "2"),
				),
			},
		},
	})
}

func TestAcc_CortexAgents_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.CortexAgents("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}
