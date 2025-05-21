//go:build !account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tags(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	tagModel := model.TagBase("test", id).
		WithComment(comment).
		WithAllowedValuesValue(tfconfig.ListVariable(tfconfig.StringVariable("foo"), tfconfig.StringVariable(""), tfconfig.StringVariable("bar")))
	tagsModel := datasourcemodel.Tags("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(tagModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tagModel, tagsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(tagsModel.DatasourceReference(), "tags.#", "1")),

					resourceshowoutputassert.TagsDatasourceShowOutput(t, "snowflake_tags.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment(comment).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttr(tagsModel.DatasourceReference(), "tags.0.show_output.0.allowed_values.#", "3")),
					assert.Check(resource.TestCheckTypeSetElemAttr(tagsModel.DatasourceReference(), "tags.0.show_output.0.allowed_values.*", "foo")),
					assert.Check(resource.TestCheckTypeSetElemAttr(tagsModel.DatasourceReference(), "tags.0.show_output.0.allowed_values.*", "")),
					assert.Check(resource.TestCheckTypeSetElemAttr(tagsModel.DatasourceReference(), "tags.0.show_output.0.allowed_values.*", "bar")),
				),
			},
		},
	})
}

func TestAcc_Tags_Filtering(t *testing.T) {
	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifier()

	model1 := model.TagBase("test1", id1)
	model2 := model.TagBase("test2", id2)
	model3 := model.TagBase("test3", id3)
	tagsModelLikeFirstOne := datasourcemodel.Tags("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	tagsModelLikePrefix := datasourcemodel.Tags("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, tagsModelLikeFirstOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tagsModelLikeFirstOne.DatasourceReference(), "tags.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, tagsModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tagsModelLikePrefix.DatasourceReference(), "tags.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Tags_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
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
