//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_UserProgrammaticAccessTokens(t *testing.T) {
	currentUser := testClient().Context.CurrentUser(t)

	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	user2, user2Cleanup := testClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)

	id1 := testClient().Ids.RandomAccountObjectIdentifier()
	id2 := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	modelComplete1 := model.UserProgrammaticAccessToken("test", id1.Name(), user.ID().Name()).
		WithRoleRestriction(snowflakeroles.Public.Name()).
		WithDaysToExpiry(10).
		WithMinsToBypassNetworkPolicyRequirement(10).
		WithDisabled("true").
		WithComment(comment)

	modelComplete2 := model.UserProgrammaticAccessToken("test2", id2.Name(), user.ID().Name()).
		WithRoleRestriction(snowflakeroles.Public.Name()).
		WithDaysToExpiry(10).
		WithMinsToBypassNetworkPolicyRequirement(10).
		WithDisabled("true").
		WithComment(comment)

	modelWithDifferentUser := model.UserProgrammaticAccessToken("test3", id1.Name(), user2.ID().Name())

	datasourceModelWithOneToken := datasourcemodel.UserProgrammaticAccessTokens("test", user.ID().Name()).
		WithDependsOn(modelComplete1.ResourceReference())

	datasourceModelWithTwoTokens := datasourcemodel.UserProgrammaticAccessTokens("test", user.ID().Name()).
		WithDependsOn(modelComplete1.ResourceReference(), modelComplete2.ResourceReference())

	datasourceModelWithDifferentUser := datasourcemodel.UserProgrammaticAccessTokens("test", user2.ID().Name()).
		WithDependsOn(modelWithDifferentUser.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete1, datasourceModelWithOneToken),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithOneToken.DatasourceReference(), "user_programmatic_access_tokens.#", "1")),

					resourceshowoutputassert.ProgrammaticAccessTokensDatasourceShowOutput(t, datasourceModelWithOneToken.DatasourceReference()).
						HasName(id1.Name()).
						HasUserName(user.ID()).
						HasRoleRestriction(snowflakeroles.Public).
						HasExpiresAtNotEmpty().
						HasStatus(sdk.ProgrammaticAccessTokenStatusDisabled).
						HasComment(comment).
						HasCreatedOnNotEmpty().
						HasCreatedBy(currentUser.Name()).
						HasMinsToBypassNetworkPolicyRequirementNotEmpty().
						HasRotatedTo(""),
				),
			},
			{
				Config: accconfig.FromModels(t, modelComplete1, modelComplete2, datasourceModelWithTwoTokens),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithTwoTokens.DatasourceReference(), "user_programmatic_access_tokens.#", "2")),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithDifferentUser, datasourceModelWithDifferentUser),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithDifferentUser.DatasourceReference(), "user_programmatic_access_tokens.#", "1")),
				),
			},
		},
	})
}
