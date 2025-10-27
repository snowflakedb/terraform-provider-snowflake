//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Databases_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomAccountObjectIdentifier()

	databaseModel1 := model.DatabaseWithParametersSet("test", idOne.Name())
	databaseModel2 := model.DatabaseWithParametersSet("test1", idTwo.Name())
	databaseModel3 := model.DatabaseWithParametersSet("test2", idThree.Name())
	databasesWithLikeModel := datasourcemodel.Databases("test").
		WithLike(idOne.Name()).
		WithDependsOn(databaseModel1.ResourceReference(), databaseModel2.ResourceReference(), databaseModel3.ResourceReference())
	databasesWithStartsWithModel := datasourcemodel.Databases("test").
		WithStartsWith(prefix).
		WithDependsOn(databaseModel1.ResourceReference(), databaseModel2.ResourceReference(), databaseModel3.ResourceReference())
	databasesWithLimitModel := datasourcemodel.Databases("test").
		WithRowsAndFrom(1, prefix).
		WithDependsOn(databaseModel1.ResourceReference(), databaseModel2.ResourceReference(), databaseModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseModel1, databaseModel2, databaseModel3, databasesWithLikeModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(databasesWithLikeModel.DatasourceReference(), "databases.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, databaseModel1, databaseModel2, databaseModel3, databasesWithStartsWithModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(databasesWithLikeModel.DatasourceReference(), "databases.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, databaseModel1, databaseModel2, databaseModel3, databasesWithLimitModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(databasesWithLikeModel.DatasourceReference(), "databases.#", "1"),
				),
			},
		},
	})
}

func TestAcc_Databases_CompleteUseCase(t *testing.T) {
	databaseId := testClient().Ids.RandomAccountObjectIdentifier()
	databaseName := databaseId.Name()
	comment := random.Comment()
	secondaryAccountId := secondaryTestClient().Account.GetAccountIdentifier(t)

	databaseModel := model.Database("test", databaseName).
		WithComment(comment).
		WithReplication(secondaryAccountId, true, true)
	databasesModel := datasourcemodel.Databases("test").
		WithLike(databaseName).
		WithStartsWith(databaseName).
		WithLimit(1).
		WithDependsOn(databaseModel.ResourceReference())
	databasesWithoutOptionalsModel := datasourcemodel.Databases("test").
		WithLike(databaseName).
		WithStartsWith(databaseName).
		WithLimit(1).
		WithDependsOn(databaseModel.ResourceReference()).
		WithWithDescribe(false).
		WithWithParameters(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, databaseModel, databasesModel),
				Check: assertThat(t,
					resourceshowoutputassert.DatabasesDatasourceShowOutput(t, databasesModel.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(databaseName).
						HasKind("STANDARD").
						HasTransient(false).
						HasIsDefault(false).
						HasIsCurrent(true).
						HasOriginEmpty().
						HasOwnerNotEmpty().
						HasComment(comment).
						HasOptions("").
						HasRetentionTimeNotEmpty().
						HasResourceGroup("").
						HasOwnerRoleTypeNotEmpty(),

					assert.Check(resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.describe_output.#", "2")),
					assert.Check(resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.describe_output.0.created_on")),

					assert.Check(resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.describe_output.0.name", "INFORMATION_SCHEMA")),
					assert.Check(resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.describe_output.0.kind", "SCHEMA")),

					assert.Check(resource.TestCheckResourceAttrSet(databasesModel.DatasourceReference(), "databases.0.describe_output.1.created_on")),
					assert.Check(resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.describe_output.1.name", "PUBLIC")),
					assert.Check(resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.describe_output.1.kind", "SCHEMA")),

					resourceparametersassert.DatabasesDatasourceParameters(t, databasesModel.DatasourceReference()).
						HasAllDefaultParameters(),
				),
			},
			{
				Config: accconfig.FromModels(t, databaseModel, databasesWithoutOptionalsModel),
				Check: assertThat(t,
					resourceshowoutputassert.DatabasesDatasourceShowOutput(t, databasesModel.DatasourceReference()).
						HasCreatedOnNotEmpty().
						HasName(databaseName).
						HasKind("STANDARD").
						HasTransient(false).
						HasIsDefault(false).
						HasIsCurrent(false).
						HasOriginEmpty().
						HasOwnerNotEmpty().
						HasComment(comment).
						HasOptions("").
						HasRetentionTimeNotEmpty().
						HasResourceGroup("").
						HasOwnerRoleTypeNotEmpty(),

					assert.Check(resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(databasesModel.DatasourceReference(), "databases.0.parameters.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Databases_DatabaseNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      databasesWithPostcondition(),
				ExpectError: regexp.MustCompile("there should be at least one database"),
			},
		},
	})
}

func databasesWithPostcondition() string {
	return `
data "snowflake_databases" "test" {
  like = "non-existing-database"

  lifecycle {
    postcondition {
      condition     = length(self.databases) > 0
      error_message = "there should be at least one database"
    }
  }
}
`
}
