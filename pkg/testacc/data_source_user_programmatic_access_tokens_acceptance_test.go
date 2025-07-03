//go:build !account_level_tests

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

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	modelComplete := model.UserProgrammaticAccessToken("test", id.Name(), user.ID().Name()).
		WithRoleRestriction(snowflakeroles.Public.Name()).
		WithDaysToExpiry(10).
		WithMinsToBypassNetworkPolicyRequirement(10).
		WithDisabled("true").
		WithComment(comment)

	datasourceModel := datasourcemodel.UserProgrammaticAccessTokens("test", user.ID().Name()).
		WithDependsOn(modelComplete.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelComplete, datasourceModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(datasourceModel.DatasourceReference(), "user_programmatic_access_tokens.#", "1")),

					resourceshowoutputassert.ProgrammaticAccessTokensDatasourceShowOutput(t, datasourceModel.DatasourceReference()).
						HasName(id.Name()).
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
		},
	})
}
