//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"strings"
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
)

const aggregationPolicyBodyMinGroupSize = "AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 5)"

func TestAcc_AggregationPolicies_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	policyModel := model.AggregationPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), aggregationPolicyBodyMinGroupSize).
		WithComment(comment)

	aggregationPoliciesModel := datasourcemodel.AggregationPolicies("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(policyModel.ResourceReference())

	aggregationPoliciesModelWithoutDescribe := datasourcemodel.AggregationPolicies("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(policyModel.ResourceReference())

	showOutputAssertions := resourceshowoutputassert.AggregationPoliciesDatasourceShowOutput(t, aggregationPoliciesModel.DatasourceReference()).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasKind("AGGREGATION_POLICY").
		HasOwner(snowflakeroles.Accountadmin.Name()).
		HasComment(comment).
		HasOwnerRoleType("ROLE").
		HasOptions("")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, policyModel, aggregationPoliciesModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(aggregationPoliciesModel.DatasourceReference(), "aggregation_policies.#", "1")),
					showOutputAssertions,
					resourceshowoutputassert.AggregationPoliciesDatasourceDescribeOutput(t, aggregationPoliciesModel.DatasourceReference()).
						HasName(id.Name()).
						HasSignature("()").
						HasReturnType("AGGREGATION_CONSTRAINT").
						HasBody(aggregationPolicyBodyMinGroupSize),
				),
			},
			{
				Config: accconfig.FromModels(t, policyModel, aggregationPoliciesModelWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(aggregationPoliciesModelWithoutDescribe.DatasourceReference(), "aggregation_policies.#", "1")),
					showOutputAssertions,
					assert.Check(resource.TestCheckResourceAttr(aggregationPoliciesModelWithoutDescribe.DatasourceReference(), "aggregation_policies.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_AggregationPolicies_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, secondSchema.ID())

	model1 := model.AggregationPolicy("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), aggregationPolicyBodyMinGroupSize)
	model2 := model.AggregationPolicy("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), aggregationPolicyBodyMinGroupSize)
	model3 := model.AggregationPolicy("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name(), aggregationPolicyBodyMinGroupSize)

	likeFirst := datasourcemodel.AggregationPolicies("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	likePrefix := datasourcemodel.AggregationPolicies("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	inSchema := datasourcemodel.AggregationPolicies("test").
		WithInSchema(id1.SchemaId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	inDatabase := datasourcemodel.AggregationPolicies("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	withLimit := datasourcemodel.AggregationPolicies("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithLimit(1).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	otherCasePrefix := strings.ToUpper(prefix)
	if otherCasePrefix == prefix {
		otherCasePrefix = strings.ToLower(prefix)
	}
	likeOtherCase := datasourcemodel.AggregationPolicies("test").
		WithLike(otherCasePrefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	inAccount := datasourcemodel.AggregationPolicies("test").
		WithLike(prefix+"%").
		WithInAccount().
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	withLimitFrom := datasourcemodel.AggregationPolicies("test").
		WithInDatabase(id1.DatabaseId()).
		WithRowsAndFrom(2, prefix).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	likePrefixWithoutDescribe := datasourcemodel.AggregationPolicies("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithWithDescribe(false).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likeFirst),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likeFirst.DatasourceReference(), "aggregation_policies.#", "1")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likePrefix),
				// With more than one result, every item holds its own show and describe output.
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(likePrefix.DatasourceReference(), "id", "aggregation_policies_read"),
					resource.TestCheckResourceAttr(likePrefix.DatasourceReference(), "aggregation_policies.#", "3"),
					resource.TestCheckResourceAttr(likePrefix.DatasourceReference(), "aggregation_policies.0.describe_output.#", "1"),
					resource.TestCheckResourceAttrPair(
						likePrefix.DatasourceReference(), "aggregation_policies.0.show_output.0.name",
						likePrefix.DatasourceReference(), "aggregation_policies.0.describe_output.0.name",
					),
					resource.TestCheckResourceAttr(likePrefix.DatasourceReference(), "aggregation_policies.1.describe_output.#", "1"),
					resource.TestCheckResourceAttrPair(
						likePrefix.DatasourceReference(), "aggregation_policies.1.show_output.0.name",
						likePrefix.DatasourceReference(), "aggregation_policies.1.describe_output.0.name",
					),
					resource.TestCheckResourceAttr(likePrefix.DatasourceReference(), "aggregation_policies.2.describe_output.#", "1"),
					resource.TestCheckResourceAttrPair(
						likePrefix.DatasourceReference(), "aggregation_policies.2.show_output.0.name",
						likePrefix.DatasourceReference(), "aggregation_policies.2.describe_output.0.name",
					),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likePrefixWithoutDescribe),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(likePrefixWithoutDescribe.DatasourceReference(), "aggregation_policies.#", "3"),
					resource.TestCheckResourceAttr(likePrefixWithoutDescribe.DatasourceReference(), "aggregation_policies.0.describe_output.#", "0"),
					resource.TestCheckResourceAttr(likePrefixWithoutDescribe.DatasourceReference(), "aggregation_policies.1.describe_output.#", "0"),
					resource.TestCheckResourceAttr(likePrefixWithoutDescribe.DatasourceReference(), "aggregation_policies.2.describe_output.#", "0"),
				),
			},
			{
				// LIKE is case-insensitive.
				Config: accconfig.FromModels(t, model1, model2, model3, likeOtherCase),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likeOtherCase.DatasourceReference(), "aggregation_policies.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inAccount),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inAccount.DatasourceReference(), "aggregation_policies.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, withLimitFrom),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(withLimitFrom.DatasourceReference(), "aggregation_policies.#", "2")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inSchema),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inSchema.DatasourceReference(), "aggregation_policies.#", "2")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inDatabase),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inDatabase.DatasourceReference(), "aggregation_policies.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, withLimit),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(withLimit.DatasourceReference(), "aggregation_policies.#", "1")),
			},
		},
	})
}

func TestAcc_AggregationPolicies_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.AggregationPolicies("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_AggregationPolicies_EmptyResult(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	aggregationPoliciesModel := datasourcemodel.AggregationPolicies("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, aggregationPoliciesModel),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(aggregationPoliciesModel.DatasourceReference(), "aggregation_policies.#", "0")),
			},
		},
	})
}
