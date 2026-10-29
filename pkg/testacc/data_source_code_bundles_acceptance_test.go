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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CodeBundles_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	stage, cleanupStage := testClient().Stage.CreateStage(t)
	t.Cleanup(cleanupStage)
	testClient().Stage.PutOnStageWithContent(t, stage.ID(), "main.py", "print('code bundle')\n")
	from := stage.Location()

	bundleModel := model.CodeBundle("test", id.DatabaseName(), id.SchemaName(), id.Name(), from).
		WithComment(comment)

	codeBundlesModel := datasourcemodel.CodeBundles("test").
		WithInSchema(id.SchemaId()).
		WithDependsOn(bundleModel.ResourceReference())
	codeBundlesModelWithoutDescribe := datasourcemodel.CodeBundles("test").
		WithWithDescribe(false).
		WithInSchema(id.SchemaId()).
		WithDependsOn(bundleModel.ResourceReference())

	showOutputAssertions := resourceshowoutputassert.CodeBundlesDatasourceShowOutput(t, codeBundlesModel.DatasourceReference()).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwner("").
		HasComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, bundleModel, codeBundlesModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(codeBundlesModel.DatasourceReference(), "code_bundles.#", "1")),
					showOutputAssertions,
					resourceshowoutputassert.CodeBundlesDatasourceDescribeOutput(t, codeBundlesModel.DatasourceReference()).
						HasName(id.Name()).
						HasComment(comment).
						HasDefaultVersion("LAST").
						HasContainerType("Stage").
						HasContainerName(stage.ID().Name()).
						HasOwnerRole("ACCOUNTADMIN").
						HasOwnerRoleType("ROLE"),
				),
			},
			{
				Config: accconfig.FromModels(t, bundleModel, codeBundlesModelWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(codeBundlesModelWithoutDescribe.DatasourceReference(), "code_bundles.#", "1")),
					showOutputAssertions,
					assert.Check(resource.TestCheckResourceAttr(codeBundlesModelWithoutDescribe.DatasourceReference(), "code_bundles.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_CodeBundles_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	stage, cleanupStage := testClient().Stage.CreateStage(t)
	t.Cleanup(cleanupStage)
	testClient().Stage.PutOnStageWithContent(t, stage.ID(), "main.py", "print('code bundle')\n")
	from := stage.Location()

	id1 := testClient().Ids.RandomSchemaObjectIdentifier()
	id2 := testClient().Ids.RandomSchemaObjectIdentifier()
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchema(secondSchema.ID())

	model1 := model.CodeBundle("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), from)
	model2 := model.CodeBundle("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), from)
	model3 := model.CodeBundle("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name(), from)

	inFirstSchema := datasourcemodel.CodeBundles("test").
		WithInSchema(id1.SchemaId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	inSecondSchema := datasourcemodel.CodeBundles("test").
		WithInSchema(id3.SchemaId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	inDatabase := datasourcemodel.CodeBundles("test").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inFirstSchema),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inFirstSchema.DatasourceReference(), "code_bundles.#", "2")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inSecondSchema),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inSecondSchema.DatasourceReference(), "code_bundles.#", "1")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inDatabase),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inDatabase.DatasourceReference(), "code_bundles.#", "3")),
			},
		},
	})
}

func TestAcc_CodeBundles_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.CodeBundles("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_CodeBundles_EmptyResult(t *testing.T) {
	emptySchema, cleanupSchema := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(cleanupSchema)

	codeBundlesModel := datasourcemodel.CodeBundles("test").
		WithInSchema(emptySchema.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, codeBundlesModel),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(codeBundlesModel.DatasourceReference(), "code_bundles.#", "0")),
			},
		},
	})
}
