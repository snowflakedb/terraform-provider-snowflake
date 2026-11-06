//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Pipes_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := testClient().PipeCopyStatement(t)
	comment := random.Comment()

	pipeModel := model.PipeWithId("test", id, statement).
		WithComment(comment).
		WithLogEventLevelSdk(sdk.LogLevelInfo)

	pipesModel := datasourcemodel.Pipes("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(pipeModel.ResourceReference())

	pipesModelWithoutAdditional := datasourcemodel.Pipes("test").
		WithWithDescribe(false).
		WithWithParameters(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(pipeModel.ResourceReference())

	showOutputAssertions := resourceshowoutputassert.PipesDatasourceShowOutput(t, pipesModel.DatasourceReference()).
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasDefinition(statement).
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasComment(comment).
		HasIntegration("").
		HasErrorIntegration("").
		HasNotificationChannel("").
		HasOwnerRoleType("ROLE").
		HasKind("STAGE")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, pipeModel, pipesModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(pipesModel.DatasourceReference(), "pipes.#", "1")),
					showOutputAssertions,
					resourceshowoutputassert.PipesDatasourceDescribeOutput(t, pipesModel.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasDefinition(statement).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment(comment),
					resourceparametersassert.PipesDatasourceParameters(t, pipesModel.DatasourceReference()).
						HasLogEventLevel(sdk.LogLevelInfo).
						HasLogEventLevelLevel(sdk.ParameterTypePipe).
						HasPipeExecutionPaused(false),
				),
			},
			{
				Config: accconfig.FromModels(t, pipeModel, pipesModelWithoutAdditional),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(pipesModelWithoutAdditional.DatasourceReference(), "pipes.#", "1")),
					showOutputAssertions,
					assert.Check(resource.TestCheckResourceAttr(pipesModelWithoutAdditional.DatasourceReference(), "pipes.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(pipesModelWithoutAdditional.DatasourceReference(), "pipes.0.parameters.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Pipes_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)
	statement := testClient().PipeCopyStatement(t)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, secondSchema.ID())

	model1 := model.PipeWithId("test1", id1, statement)
	model2 := model.PipeWithId("test2", id2, statement)
	model3 := model.PipeWithId("test3", id3, statement)
	dependsOn := []string{model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference()}

	likeFirst := datasourcemodel.Pipes("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(dependsOn...)
	likePrefixInSchema := datasourcemodel.Pipes("test").
		WithLike(prefix + "%").
		WithInSchema(id1.SchemaId()).
		WithDependsOn(dependsOn...)
	likePrefixInDatabase := datasourcemodel.Pipes("test").
		WithLike(prefix + "%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(dependsOn...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likeFirst),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likeFirst.DatasourceReference(), "pipes.#", "1")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likePrefixInSchema),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likePrefixInSchema.DatasourceReference(), "pipes.#", "2")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likePrefixInDatabase),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likePrefixInDatabase.DatasourceReference(), "pipes.#", "3")),
			},
		},
	})
}

func TestAcc_Pipes_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.Pipes("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Pipes_InApplication(t *testing.T) {
	applicationPackage, cleanup := testClient().ApplicationPackage.CreateApplicationPackage(t)
	t.Cleanup(cleanup)

	inApplicationPackage := datasourcemodel.Pipes("test").
		WithInApplicationPackage(applicationPackage.ID())
	inMissingApplication := datasourcemodel.Pipes("test").
		WithInApplication(testClient().Ids.RandomAccountObjectIdentifier())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, inApplicationPackage),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inApplicationPackage.DatasourceReference(), "pipes.#", "0")),
			},
			{
				Config:      accconfig.FromModels(t, inMissingApplication),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_Pipes_EmptyResult(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	pipesModel := datasourcemodel.Pipes("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, pipesModel),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(pipesModel.DatasourceReference(), "pipes.#", "0")),
			},
		},
	})
}
