//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Provider_OauthWithClientCredentials(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	oauthClientId := testenvs.GetOrSkipTest(t, testenvs.OauthClientId)
	oauthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OauthClientSecret)
	oauthTokenRequestURL := testenvs.GetOrSkipTest(t, testenvs.OauthTokenRequestUrl)

	// user := testClient().SetUpTemporaryUserWithOauthClientCredentials(t, oauthClientId)
	user := helpers.TmpUser{
		UserId:    sdk.NewAccountObjectIdentifier(oauthClientId),
		AccountId: testClient().Context.CurrentAccountId(t),
		RoleId:    snowflakeroles.Public,
	}
	userConfig := testClient().TempTomlConfigForServiceUserWithOauthClientCredentials(t, &user, oauthClientId, oauthClientSecret, oauthTokenRequestURL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, userConfig.Path)
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(userConfig.Profile)) + executeShowSessionParameter(),
			},
		},
	})
}

// TODO: move to common package
func executeShowSessionParameter() string {
	return `
resource snowflake_execute "t" {
    execute = "SELECT 1"
    query = "SHOW PARAMETERS LIKE 'STATEMENT_TIMEOUT_IN_SECONDS' IN SESSION"
    revert        = "SELECT 1"
}`
}
