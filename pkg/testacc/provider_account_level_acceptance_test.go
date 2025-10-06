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

// This test is marked as account_level_tests because it creates an Oauth security integration with a unique issuer and a user with a unique login name.
func TestAcc_Provider_OauthWithClientCredentials(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	oauthClientId := testenvs.GetOrSkipTest(t, testenvs.OauthWithClientCredentialsClientId)
	oauthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OauthWithClientCredentialsClientSecret)
	oauthIssuerUrl := testenvs.GetOrSkipTest(t, testenvs.OauthWithClientCredentialsIssuer)
	oauthJwsKeysUrl := oauthIssuerUrl + "/v1/keys"
	oauthTokenRequestUrl := oauthIssuerUrl + "/v1/token"

	user, userCleanup := testClient().User.CreateUserWithOptions(t, sdk.NewAccountObjectIdentifier(oauthClientId), &sdk.CreateUserOptions{
		ObjectProperties: &sdk.UserObjectProperties{
			Type:      sdk.Pointer(sdk.UserTypeService),
			LoginName: sdk.String(oauthClientId),
		},
	})
	t.Cleanup(userCleanup)
	url := testClient().Context.AccountURL(t)

	securityIntegrationId := testClient().Ids.RandomAccountObjectIdentifier()
	_, securityIntegrationCleanup := testClient().SecurityIntegration.CreateExternalOauthWithRequest(
		t,
		sdk.NewCreateExternalOauthSecurityIntegrationRequest(
			securityIntegrationId,
			true,
			sdk.ExternalOauthSecurityIntegrationTypeOkta,
			oauthIssuerUrl,
			[]sdk.TokenUserMappingClaim{{Claim: "sub"}},
			sdk.ExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeLoginName,
		).WithExternalOauthJwsKeysUrl([]sdk.JwsKeysUrl{{JwsKeyUrl: oauthJwsKeysUrl}}).
			WithExternalOauthAudienceList(sdk.AudienceListRequest{AudienceList: []sdk.AudienceListItem{{Item: url}}}),
	)
	t.Cleanup(securityIntegrationCleanup)

	userHelper := helpers.TmpUser{
		UserId:    user.ID(),
		AccountId: testClient().Context.CurrentAccountId(t),
		RoleId:    snowflakeroles.Public,
	}
	userConfig := testClient().TempTomlConfigForServiceUserWithOauthClientCredentials(t, &userHelper, oauthClientId, oauthClientSecret, oauthTokenRequestUrl)

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
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(userConfig.Profile)) + helpers.DummyResource(),
			},
		},
	})
}
