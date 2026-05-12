//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ImageRepositories(t *testing.T) {
	// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	comment := random.Comment()

	imageRepositoryModel := model.ImageRepository("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithComment(comment)

	dataSourceModel := datasourcemodel.ImageRepositories("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(imageRepositoryModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, imageRepositoryModel, dataSourceModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "image_repositories.#", "1")),

					resourceshowoutputassert.ImageRepositoriesDatasourceShowOutput(t, "snowflake_image_repositories.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasRepositoryUrlNotEmpty().
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment).
						HasPrivatelinkRepositoryUrl(""),
				),
			},
		},
	})
}

func TestAcc_ImageRepositories_Filtering(t *testing.T) {
	// TODO(SNOW-2070746): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	prefix := random.AlphaUpperN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
	id2 := testClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	model1 := model.ImageRepository("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name())
	model2 := model.ImageRepository("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name())
	model3 := model.ImageRepository("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name())

	dataSourceModelLikeFirstOne := datasourcemodel.ImageRepositories("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	dataSourceModelLikePrefix := datasourcemodel.ImageRepositories("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, dataSourceModelLikeFirstOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelLikeFirstOne.DatasourceReference(), "image_repositories.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, dataSourceModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelLikePrefix.DatasourceReference(), "image_repositories.#", "2"),
				),
			},
		},
	})
}

func TestAcc_ImageRepositories_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.ImageRepositories("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_ImageRepositories_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_ImageRepositories/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one image repository"),
			},
		},
	})
}
