//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
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
			WithExternalOauthAllowedRolesList(sdk.AllowedRolesListRequest{AllowedRolesList: []sdk.AccountObjectIdentifier{snowflakeroles.Public}}).
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
				Config: config.FromModels(t, providermodel.SnowflakeProvider().WithProfile(userConfig.Profile), model.ExecuteWithNoOpActions("t")),
			},
		},
	})
}

func TestAcc_Provider_OauthWithAuthorizationCodeSnowflakeIdp(t *testing.T) {
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

func TestAcc_Provider_OauthWithAuthorizationCodeExternalIdp(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	oauthClientId := testenvs.GetOrSkipTest(t, testenvs.OauthWithAuthorizationCodeExternalIdpClientId)
	oauthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OauthWithAuthorizationCodeExternalIdpClientSecret)
	oauthIssuerUrl := testenvs.GetOrSkipTest(t, testenvs.OauthWithAuthorizationCodeExternalIdpIssuer)
	loginName := testenvs.GetOrSkipTest(t, testenvs.OauthWithAuthorizationCodeExternalIdpLoginName)
	oauthAuthorizationUrl := oauthIssuerUrl + "/v1/authorize"
	oauthTokenRequestUrl := oauthIssuerUrl + "/v1/token"
	oauthJwsKeysUrl := oauthIssuerUrl + "/v1/keys"
	password := random.Password()
	id := testClient().Ids.RandomAccountObjectIdentifier()
	url := testClient().Context.AccountURL(t)

	_, userCleanup := testClient().User.CreateUserWithOptions(t, id, &sdk.CreateUserOptions{
		ObjectProperties: &sdk.UserObjectProperties{
			// The service users are prohibited from this type of authentication.
			LoginName:          sdk.String(loginName),
			Password:           sdk.String(password),
			MustChangePassword: sdk.Bool(false),
		},
	})
	t.Cleanup(userCleanup)

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

	userHelper := helpers.TmpServiceUser{
		TmpUser: helpers.TmpUser{
			UserId:    id,
			AccountId: testClient().Context.CurrentAccountId(t),
			RoleId:    snowflakeroles.Public,
		},
		Pass: password,
	}
	userConfig := testClient().TempTomlConfigForServiceUserWithOauthAuthorizationCodeExternalIdp(
		t,
		&userHelper,
		oauthClientId,
		oauthClientSecret,
		oauthTokenRequestUrl,
		oauthAuthorizationUrl,
		"http://localhost:8001",
		"session:role:PUBLIC",
	)

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
