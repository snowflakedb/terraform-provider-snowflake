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
)

const projectionPolicyBodyDeny = "PROJECTION_CONSTRAINT(ALLOW => false)"

func TestAcc_ProjectionPolicies_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	policyModel := model.ProjectionPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name(), projectionPolicyBodyDeny).
		WithComment(comment)

	projectionPoliciesModel := datasourcemodel.ProjectionPolicies("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(policyModel.ResourceReference())

	projectionPoliciesModelWithoutDescribe := datasourcemodel.ProjectionPolicies("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(policyModel.ResourceReference())

	showOutputAssertions := resourceshowoutputassert.ProjectionPoliciesDatasourceShowOutput(t, projectionPoliciesModel.DatasourceReference()).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasKind("PROJECTION_POLICY").
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
				Config: accconfig.FromModels(t, policyModel, projectionPoliciesModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(projectionPoliciesModel.DatasourceReference(), "projection_policies.#", "1")),
					showOutputAssertions,
					resourceshowoutputassert.ProjectionPoliciesDatasourceDescribeOutput(t, projectionPoliciesModel.DatasourceReference()).
						HasName(id.Name()).
						HasSignature("()").
						HasReturnType("PROJECTION_CONSTRAINT").
						HasBody(projectionPolicyBodyDeny),
				),
			},
			{
				Config: accconfig.FromModels(t, policyModel, projectionPoliciesModelWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(projectionPoliciesModelWithoutDescribe.DatasourceReference(), "projection_policies.#", "1")),
					showOutputAssertions,
					assert.Check(resource.TestCheckResourceAttr(projectionPoliciesModelWithoutDescribe.DatasourceReference(), "projection_policies.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_ProjectionPolicies_Filtering(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchemaInDatabase(t, sdk.NewAccountObjectIdentifier(TestDatabaseName))
	t.Cleanup(secondSchemaCleanup)

	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, secondSchema.ID())

	model1 := model.ProjectionPolicy("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), projectionPolicyBodyDeny)
	model2 := model.ProjectionPolicy("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), projectionPolicyBodyDeny)
	model3 := model.ProjectionPolicy("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name(), projectionPolicyBodyDeny)

	likeFirst := datasourcemodel.ProjectionPolicies("test").
		WithLike(id1.Name()).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	likePrefix := datasourcemodel.ProjectionPolicies("test").
		WithLike(prefix+"%").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	inSchema := datasourcemodel.ProjectionPolicies("test").
		WithInSchema(id1.SchemaId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	inDatabase := datasourcemodel.ProjectionPolicies("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	withLimit := datasourcemodel.ProjectionPolicies("test").
		WithStartsWith(prefix).
		WithInDatabase(id1.DatabaseId()).
		WithLimit(1).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likeFirst),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likeFirst.DatasourceReference(), "projection_policies.#", "1")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, likePrefix),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(likePrefix.DatasourceReference(), "projection_policies.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inSchema),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inSchema.DatasourceReference(), "projection_policies.#", "2")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, inDatabase),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(inDatabase.DatasourceReference(), "projection_policies.#", "3")),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, withLimit),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(withLimit.DatasourceReference(), "projection_policies.#", "1")),
			},
		},
	})
}

func TestAcc_ProjectionPolicies_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.ProjectionPolicies("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_ProjectionPolicies_EmptyResult(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	projectionPoliciesModel := datasourcemodel.ProjectionPolicies("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, projectionPoliciesModel),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(projectionPoliciesModel.DatasourceReference(), "projection_policies.#", "0")),
			},
		},
	})
}
