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

func TestAcc_Quotas_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	quotaModel := model.QuotaFromId("test", id).
		WithComment(comment).
		WithSharedResources(model.QuotaSharedResourceModel{Domain: sdk.QuotaSharedResourceDomainAiFunction}).
		WithPerUserLimitMonthly(1000)

	quotasModel := datasourcemodel.Quotas("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(quotaModel.ResourceReference())

	quotasModelWithoutDescribe := datasourcemodel.Quotas("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(quotaModel.ResourceReference())

	showOutputAssertions := resourceshowoutputassert.QuotasDatasourceShowOutput(t, quotasModel.DatasourceReference()).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasOwnerNotEmpty().
		HasComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, quotaModel, quotasModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(quotasModel.DatasourceReference(), "quotas.#", "1")),
					showOutputAssertions,
					assert.Check(resource.TestCheckResourceAttr(quotasModel.DatasourceReference(), "quotas.0.describe_output.0.scope.0.shared_resources.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(quotasModel.DatasourceReference(), "quotas.0.describe_output.0.scope.0.shared_resources.0.domain", string(sdk.QuotaSharedResourceDomainAiFunction))),
					assert.Check(resource.TestCheckResourceAttr(quotasModel.DatasourceReference(), "quotas.0.describe_output.0.config.0.per_user_limit", "1000")),
				),
			},
			{
				Config: accconfig.FromModels(t, quotaModel, quotasModelWithoutDescribe),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(quotasModelWithoutDescribe.DatasourceReference(), "quotas.#", "1")),
					showOutputAssertions,
					assert.Check(resource.TestCheckResourceAttr(quotasModelWithoutDescribe.DatasourceReference(), "quotas.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Quotas_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchema(secondSchema.ID())

	model1 := model.QuotaFromId("test1", id1)
	model2 := model.QuotaFromId("test2", id2)
	model3 := model.QuotaFromId("test3", id3)
	dependsOn := []string{model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference()}

	quotasModelLikeFirst := datasourcemodel.Quotas("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(dependsOn...)

	quotasModelLikePrefix := datasourcemodel.Quotas("test").
		WithLike(prefix + "%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(dependsOn...)

	quotasModelInSchema := datasourcemodel.Quotas("test").
		WithInSchema(id1.SchemaId()).
		WithDependsOn(dependsOn...)

	quotasModelLimit := datasourcemodel.Quotas("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithRowsAndFrom(1, prefix).
		WithDependsOn(dependsOn...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, quotasModelLikeFirst),
				Check:  resource.TestCheckResourceAttr(quotasModelLikeFirst.DatasourceReference(), "quotas.#", "1"),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, quotasModelLikePrefix),
				Check:  resource.TestCheckResourceAttr(quotasModelLikePrefix.DatasourceReference(), "quotas.#", "2"),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, quotasModelInSchema),
				Check:  resource.TestCheckResourceAttr(quotasModelInSchema.DatasourceReference(), "quotas.#", "2"),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, quotasModelLimit),
				Check:  resource.TestCheckResourceAttr(quotasModelLimit.DatasourceReference(), "quotas.#", "1"),
			},
		},
	})
}

func TestAcc_Quotas_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.Quotas("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}
