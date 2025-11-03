//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tags_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	tagId1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	tagId2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	tagId3 := testClient().Ids.RandomSchemaObjectIdentifier()

	tagModel1 := model.TagBase("test", tagId1)
	tagModel2 := model.TagBase("test1", tagId2)
	tagModel3 := model.TagBase("test2", tagId3)

	// like only (exact match)
	datasourceModelLikeExact := datasourcemodel.Tags("test").
		WithLike(tagId1.Name()).
		WithDependsOn(tagModel1.ResourceReference(), tagModel2.ResourceReference(), tagModel3.ResourceReference())

	// like only (prefix pattern)
	datasourceModelLikePrefix := datasourcemodel.Tags("test").
		WithLike(prefix+"%").
		WithDependsOn(tagModel1.ResourceReference(), tagModel2.ResourceReference(), tagModel3.ResourceReference())

	// in database only
	datasourceModelInDatabase := datasourcemodel.Tags("test").
		WithInDatabase(tagId1.DatabaseId()).
		WithDependsOn(tagModel1.ResourceReference(), tagModel2.ResourceReference(), tagModel3.ResourceReference())

	// in schema only
	datasourceModelInSchema := datasourcemodel.Tags("test").
		WithInSchema(tagId1.SchemaId()).
		WithDependsOn(tagModel1.ResourceReference(), tagModel2.ResourceReference(), tagModel3.ResourceReference())

	// like + in database (combined filters)
	datasourceModelLikeAndInDatabase := datasourcemodel.Tags("test").
		WithLike(prefix+"%").
		WithInDatabase(tagId1.DatabaseId()).
		WithDependsOn(tagModel1.ResourceReference(), tagModel2.ResourceReference(), tagModel3.ResourceReference())

	// like + in schema (combined filters)
	datasourceModelLikeAndInSchema := datasourcemodel.Tags("test").
		WithLike(prefix+"%").
		WithInSchema(tagId1.SchemaId()).
		WithDependsOn(tagModel1.ResourceReference(), tagModel2.ResourceReference(), tagModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, tagModel1, tagModel2, tagModel3, datasourceModelLikeExact),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "tags.#", "1"),
					resource.TestCheckResourceAttr(datasourceModelLikeExact.DatasourceReference(), "tags.0.show_output.0.name", tagId1.Name()),
				),
			},
			{
				Config: config.FromModels(t, tagModel1, tagModel2, tagModel3, datasourceModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikePrefix.DatasourceReference(), "tags.#", "2"),
				),
			},
			{
				Config: config.FromModels(t, tagModel1, tagModel2, tagModel3, datasourceModelInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelInDatabase.DatasourceReference(), "tags.#", "3"),
				),
			},
			{
				Config: config.FromModels(t, tagModel1, tagModel2, tagModel3, datasourceModelInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelInSchema.DatasourceReference(), "tags.#", "3"),
				),
			},
			{
				Config: config.FromModels(t, tagModel1, tagModel2, tagModel3, datasourceModelLikeAndInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeAndInDatabase.DatasourceReference(), "tags.#", "2"),
				),
			},
			{
				Config: config.FromModels(t, tagModel1, tagModel2, tagModel3, datasourceModelLikeAndInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModelLikeAndInSchema.DatasourceReference(), "tags.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Tags_CompleteUseCase(t *testing.T) {
	tagId := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	tagModel := model.TagBase("test", tagId).
		WithComment(comment).
		WithAllowedValuesValue(tfconfig.ListVariable(tfconfig.StringVariable("foo"), tfconfig.StringVariable("bar")))

	datasourceModel := datasourcemodel.Tags("test").
		WithLike(tagId.Name()).
		WithInDatabase(tagId.DatabaseId()).
		WithDependsOn(tagModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, tagModel, datasourceModel),
				Check: assertThat(t,
					resourceshowoutputassert.TagsDatasourceShowOutput(t, "test").
						HasCreatedOnNotEmpty().
						HasName(tagId.Name()).
						HasDatabaseName(tagId.DatabaseName()).
						HasSchemaName(tagId.SchemaName()).
						HasComment(comment).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasAllowedValuesUnordered("foo", "bar"),
				),
			},
		},
	})
}

func TestAcc_Tags_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      tagsDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func tagsDatasourceEmptyIn() string {
	return `
data "snowflake_tags" "test" {
  in {
  }
}
`
}

func TestAcc_Tags_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Tags/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one tag"),
			},
		},
	})
}
