//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// The data source covers both deployment types - SHOW OPENFLOW DEPLOYMENTS does not split them - but this
// test uses BYOC only, since a Snowflake-managed deployment takes six to seven minutes to reach ACTIVE while
// BYOC settles at INACTIVE in about ninety seconds. The type column is asserted so the distinction is still
// covered.
//
// SHOW OPENFLOW DEPLOYMENTS is account-scoped, so there is no `in` filter to exercise. The other three
// filters it does accept are covered below.
func TestAcc_OpenflowDeployments_BasicUseCase_DifferentFiltering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	currentRole := testClient().Context.CurrentRole(t)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	displayName := random.AlphaN(12)

	deploymentModel := model.OpenflowDeploymentByoc("t", id.Name(), string(sdk.OpenflowVpcTypeManaged)).
		WithDisplayName(displayName).
		WithComment(comment)

	dsLikeExact := datasourcemodel.OpenflowDeployments("test").
		WithLike(id.Name()).
		WithDependsOn(deploymentModel.ResourceReference())
	dsLikeMatchingNothing := datasourcemodel.OpenflowDeployments("test").
		WithLike("non_existing_deployment").
		WithDependsOn(deploymentModel.ResourceReference())
	dsWithoutDescribe := datasourcemodel.OpenflowDeployments("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(deploymentModel.ResourceReference())
	dsStartsWith := datasourcemodel.OpenflowDeployments("test").
		WithStartsWith(id.Name()).
		WithDependsOn(deploymentModel.ResourceReference())
	// The account holds other deployments, so LIMIT 1 is asserted as returning exactly one row rather than
	// as returning this deployment.
	dsLimit := datasourcemodel.OpenflowDeployments("test").
		WithLimit(1).
		WithDependsOn(deploymentModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		// The deployment here is only a fixture for the data source to read; its own resource tests cover
		// destroy.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// 1. like on an exact name, describe on by default - assert every exposed field
			{
				Config: config.FromModels(t, deploymentModel, dsLikeExact),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLikeExact.DatasourceReference(), "openflow_deployments.#", "1")),
					resourceshowoutputassert.OpenflowDeploymentsDatasourceShowOutput(t, dsLikeExact.DatasourceReference()).
						HasName(id.Name()).
						HasType(sdk.OpenflowDeploymentTypeByoc).
						// BYOC creation is complete at INACTIVE: the deployment only becomes ACTIVE once the
						// customer provisions it in their own cloud account.
						HasStatus(sdk.OpenflowDeploymentStatusInactive).
						HasVpcType(sdk.OpenflowVpcTypeManaged).
						HasUsePrivateLink(false).
						HasUseUserAuthOverPrivateLink(false).
						HasDisplayName(displayName).
						HasOwner(currentRole.Name()).
						HasComment(comment).
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty(),
					// describe_output is the SHOW set minus created_on and updated_on, plus the key that names
					// the customer's cloud resources.
					resourceshowoutputassert.OpenflowDeploymentsDatasourceDescribeOutput(t, dsLikeExact.DatasourceReference()).
						HasName(id.Name()).
						HasType(sdk.OpenflowDeploymentTypeByoc).
						HasVpcType(sdk.OpenflowVpcTypeManaged).
						HasOwner(currentRole.Name()).
						HasComment(comment).
						HasKeyNotEmpty(),
				),
			},
			// 2. like matching nothing
			{
				Config: config.FromModels(t, deploymentModel, dsLikeMatchingNothing),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsLikeMatchingNothing.DatasourceReference(), "openflow_deployments.#", "0"),
				),
			},
			// 3. with_describe = false suppresses the DESCRIBE call, so the block must be empty
			{
				Config: config.FromModels(t, deploymentModel, dsWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsWithoutDescribe.DatasourceReference(), "openflow_deployments.#", "1")),
					resourceshowoutputassert.OpenflowDeploymentsDatasourceShowOutput(t, dsWithoutDescribe.DatasourceReference()).
						HasName(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(dsWithoutDescribe.DatasourceReference(), "openflow_deployments.0.describe_output.#", "0")),
				),
			},
			// 4. starts_with on the full name, which is unique, so it returns exactly this deployment
			{
				Config: config.FromModels(t, deploymentModel, dsStartsWith),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsStartsWith.DatasourceReference(), "openflow_deployments.#", "1")),
					resourceshowoutputassert.OpenflowDeploymentsDatasourceShowOutput(t, dsStartsWith.DatasourceReference()).
						HasName(id.Name()),
				),
			},
			// 5. limit caps the result set
			{
				Config: config.FromModels(t, deploymentModel, dsLimit),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLimit.DatasourceReference(), "openflow_deployments.#", "1")),
				),
			},
		},
	})
}
