//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// A Snowflake-managed deployment takes six to seven minutes to reach ACTIVE, with a comparable terminate, so
// this covers the lifecycle in six steps rather than the nine the cheaper BYOC resource uses.
//
// There is no ForceNew attribute here and so no third model: this resource takes none of BYOC's networking
// options, which is why the two deployment types are separate resources. Rename is covered on BYOC.
func TestAcc_OpenflowDeploymentSnowflakeManaged_BasicUseCase(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	currentRole := testClient().Context.CurrentRole(t)
	// event_table is a parameter, so with none set on the deployment it reports what the account has.
	accountEventTable := helpers.FindParameter(t, testClient().Parameter.ShowAccountParameters(t), sdk.AccountParameterEventTable).Value

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	displayName := random.AlphaN(12)
	externallyChangedComment := random.Comment()
	externallyChangedDisplayName := random.AlphaN(12)

	basic := model.OpenflowDeploymentSnowflakeManaged("t", id.Name())
	withOptionals := model.OpenflowDeploymentSnowflakeManaged("t", id.Name()).
		WithDisplayName(displayName).
		WithComment(comment)
	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowDeploymentSnowflakeManagedResource(t, ref).
			HasNameString(id.Name()).
			HasTypeString(string(sdk.OpenflowDeploymentTypeSnowflake)).
			HasNoDisplayName().
			HasNoComment().
			HasEventTableString(accountEventTable).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowDeploymentShowOutput(t, ref).
			HasName(id.Name()).
			HasType(sdk.OpenflowDeploymentTypeSnowflake).
			HasStatus(sdk.OpenflowDeploymentStatusActive).
			HasOwner(currentRole.Name()),
		resourceshowoutputassert.OpenflowDeploymentDescribeOutput(t, ref).
			HasName(id.Name()).
			HasType(sdk.OpenflowDeploymentTypeSnowflake).
			HasOwner(currentRole.Name()),
		objectassert.OpenflowDeployment(t, id).
			HasName(id.Name()).
			HasType(sdk.OpenflowDeploymentTypeSnowflake).
			// Unlike BYOC, a Snowflake-managed deployment reaches ACTIVE on its own, so create is only
			// complete once it gets there.
			HasStatus(sdk.OpenflowDeploymentStatusActive).
			HasOwner(currentRole.Name()).
			HasNoDisplayName().
			HasNoComment().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
	}

	assertOptionalsUnset := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowDeploymentSnowflakeManagedResource(t, ref).
			HasNameString(id.Name()).
			HasTypeString(string(sdk.OpenflowDeploymentTypeSnowflake)).
			HasDisplayNameString("").
			HasCommentString("").
			HasEventTableString(accountEventTable).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		objectassert.OpenflowDeployment(t, id).
			HasName(id.Name()).
			HasNoDisplayName().
			HasNoComment(),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowDeploymentSnowflakeManagedResource(t, ref).
			HasNameString(id.Name()).
			HasTypeString(string(sdk.OpenflowDeploymentTypeSnowflake)).
			HasDisplayNameString(displayName).
			HasCommentString(comment),
		resourceshowoutputassert.OpenflowDeploymentShowOutput(t, ref).
			HasName(id.Name()).
			HasDisplayName(displayName).
			HasComment(comment),
		objectassert.OpenflowDeployment(t, id).
			HasDisplayName(displayName).
			HasComment(comment).
			HasStatus(sdk.OpenflowDeploymentStatusActive),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowDeploymentSnowflakeManaged),
		Steps: []resource.TestStep{
			// 1. Create - required fields only
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// 2. Import - required fields only
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// 3. Update - set the optionals. Both are altered in place; recreating a deployment would cost
			// minutes, so a plan that replaces here is a bug.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// 4. Import - with the optionals set
			{
				Config:            config.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// 5. External change - both configured fields are changed outside Terraform and must be reconciled
			{
				PreConfig: func() {
					testClient().OpenflowDeployment.Alter(t, sdk.NewAlterOpenflowDeploymentRequest(id).
						WithSet(*sdk.NewOpenflowDeploymentSetRequest().
							WithComment(externallyChangedComment).
							WithDisplayName(externallyChangedDisplayName)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externallyChangedComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externallyChangedComment), sdk.String(comment)),
						planchecks.ExpectDrift(ref, "display_name", sdk.String(displayName), sdk.String(externallyChangedDisplayName)),
						planchecks.ExpectChange(ref, "display_name", tfjson.ActionUpdate, sdk.String(externallyChangedDisplayName), sdk.String(displayName)),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// 6. Update - unset the optionals again
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertOptionalsUnset...),
			},
		},
	})
}
