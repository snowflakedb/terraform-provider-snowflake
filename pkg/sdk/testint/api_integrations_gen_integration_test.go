//go:build non_account_level_tests

package testint

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ApiIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// TODO [SNOW-1017580]: replace with real values when testing with external function invocation.
	const awsPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/dev/"
	const awsOtherPrefix = "https://123456.execute-api.us-west-2.amazonaws.com/prod/"
	const azurePrefix = "https://apim-hello-world.azure-api.net/dev"
	const azureOtherPrefix = "https://apim-hello-world.azure-api.net/prod"
	const googlePrefix = "https://gateway-id-123456.uc.gateway.dev/prod"
	const googleOtherPrefix = "https://gateway-id-123456.uc.gateway.dev/dev"
	const gitPrefix = "https://github.com/my-org/"
	const gitOtherPrefix = "https://github.com/my-org/other/"
	const mcpPrefix = "https://mcp.example.com/api/"
	const apiAwsRoleArn = "arn:aws:iam::000000000001:/role/test"
	const azureTenantId = "00000000-0000-0000-0000-000000000000"
	const azureOtherTenantId = "11111111-1111-1111-1111-111111111111"
	const azureAdApplicationId = "11111111-1111-1111-1111-111111111111"
	const googleAudience = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
	const googleOtherAudience = "api-gateway-id-666777.apigateway.gcp-project.cloud.goog"
	const oauthAuthorizationEndpoint = "https://auth.example.com/authorize"
	const oauthTokenEndpoint = "https://auth.example.com/token"
	const oauthClientId = "oauth-client-id-123"
	const oauthClientSecret = "oauth-client-secret-456"
	const oauthResourceUrl = "https://resource.example.com"

	prefixes := func(prefix string) []sdk.ApiIntegrationEndpointPrefix {
		return []sdk.ApiIntegrationEndpointPrefix{{Path: prefix}}
	}

	// -----------------------------------------------------------------------
	// CREATE
	// -----------------------------------------------------------------------

	t.Run("create: aws basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(awsPrefix), true).
			WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, apiAwsRoleArn)))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiKey("").
			HasApiProvider(strings.ToUpper(string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway))).
			HasApiAwsRoleArn(apiAwsRoleArn).
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty().
			HasAllowedPrefixes(awsPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("create: aws all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(awsPrefix), true).
			WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, apiAwsRoleArn).WithApiKey("key")).
			WithApiBlockedPrefixes(prefixes(awsOtherPrefix)).
			WithComment("comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment("comment"),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiKey("☺☺☺").
			HasApiProvider(strings.ToUpper(string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway))).
			HasApiAwsRoleArn(apiAwsRoleArn).
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty().
			HasAllowedPrefixes(awsPrefix).
			HasBlockedPrefixes(awsOtherPrefix).
			HasComment("comment"),
		)
	})

	t.Run("create: aws all provider type variants", func(t *testing.T) {
		for _, providerType := range sdk.AllApiIntegrationAwsApiProviderTypes {
			providerType := providerType
			t.Run(string(providerType), func(t *testing.T) {
				id := testClientHelper().Ids.RandomAccountObjectIdentifier()
				req := sdk.NewCreateApiIntegrationRequest(id, prefixes(awsPrefix), true).
					WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(providerType, apiAwsRoleArn))
				err := client.ApiIntegrations.Create(ctx, req)
				require.NoError(t, err)
				t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

				assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
					HasApiProvider(strings.ToUpper(string(providerType))),
				)
			})
		}
	})

	t.Run("create: azure basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(azurePrefix), true).
			WithAzureApiProviderParams(*sdk.NewAzureApiParamsRequest(azureTenantId, azureAdApplicationId)))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiKey("").
			HasApiProvider("AZURE_API_MANAGEMENT").
			HasAzureTenantId(azureTenantId).
			HasAzureAdApplicationId(azureAdApplicationId).
			HasAzureMultiTenantAppNameNotEmpty().
			HasAzureConsentUrlNotEmpty().
			HasAllowedPrefixes(azurePrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("create: azure all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(azurePrefix), true).
			WithAzureApiProviderParams(*sdk.NewAzureApiParamsRequest(azureTenantId, azureAdApplicationId).WithApiKey("key")).
			WithApiBlockedPrefixes(prefixes(azureOtherPrefix)).
			WithComment("comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment("comment"),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiKey("☺☺☺").
			HasApiProvider("AZURE_API_MANAGEMENT").
			HasAzureTenantId(azureTenantId).
			HasAzureAdApplicationId(azureAdApplicationId).
			HasAzureMultiTenantAppNameNotEmpty().
			HasAzureConsentUrlNotEmpty().
			HasAllowedPrefixes(azurePrefix).
			HasBlockedPrefixes(azureOtherPrefix).
			HasComment("comment"),
		)
	})

	t.Run("create: google basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(googlePrefix), true).
			WithGoogleApiProviderParams(*sdk.NewGoogleApiParamsRequest(googleAudience)))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProvider("GOOGLE_API_GATEWAY").
			HasGoogleAudience(googleAudience).
			HasAllowedPrefixes(googlePrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("create: google all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(googlePrefix), true).
			WithGoogleApiProviderParams(*sdk.NewGoogleApiParamsRequest(googleAudience)).
			WithApiBlockedPrefixes(prefixes(googleOtherPrefix)).
			WithComment("comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment("comment"),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProvider("GOOGLE_API_GATEWAY").
			HasGoogleAudience(googleAudience).
			HasAllowedPrefixes(googlePrefix).
			HasBlockedPrefixes(googleOtherPrefix).
			HasComment("comment"),
		)
	})

	t.Run("create: git https api with token auth (all secrets)", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(gitPrefix), true).
			WithGitHttpsApiTokenBasedProviderParams(*sdk.NewGitHttpsApiTokenBasedParamsRequest().
				WithAllowedAuthenticationSecrets(*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().WithAllSecrets(true))))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProviderNotEmpty().
			HasAllowedPrefixes(gitPrefix),
		)
	})

	t.Run("create: git https api with token auth (no secrets)", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(gitPrefix), true).
			WithGitHttpsApiTokenBasedProviderParams(*sdk.NewGitHttpsApiTokenBasedParamsRequest().
				WithAllowedAuthenticationSecrets(*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().WithNoSecrets(true))))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasAllowedPrefixes(gitPrefix),
		)
	})

	t.Run("create: git https api with GitHub App auth", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(gitPrefix), true).
			WithGitHttpsApiGithubAppProviderParams(*sdk.NewGitHttpsApiGithubAppParamsRequest()))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProviderNotEmpty().
			HasAllowedPrefixes(gitPrefix),
		)
	})

	t.Run("create: git https api with OAuth2 auth basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		auth := sdk.NewOAuth2GitUserAuthenticationRequest(oauthAuthorizationEndpoint, oauthTokenEndpoint, oauthClientId, oauthClientSecret)
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(gitPrefix), true).
			WithGitHttpsApiOAuth2ProviderParams(*sdk.NewGitHttpsApiOAuth2ParamsRequest().WithApiUserAuthentication(*auth)))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProviderNotEmpty().
			HasAllowedPrefixes(gitPrefix),
		)
	})

	t.Run("create: git https api with OAuth2 auth all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		auth := sdk.NewOAuth2GitUserAuthenticationRequest(oauthAuthorizationEndpoint, oauthTokenEndpoint, oauthClientId, oauthClientSecret).
			WithOauthAccessTokenValidity(3600).
			WithOauthRefreshTokenValidity(86400).
			WithOauthAllowedScopes([]sdk.ApiIntegrationScope{{Scope: "read"}, {Scope: "write"}}).
			WithOauthUsername("user@example.com")
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(gitPrefix), true).
			WithGitHttpsApiOAuth2ProviderParams(*sdk.NewGitHttpsApiOAuth2ParamsRequest().WithApiUserAuthentication(*auth)).
			WithApiBlockedPrefixes(prefixes(gitOtherPrefix)).
			WithComment("git oauth2 comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProviderNotEmpty().
			HasAllowedPrefixes(gitPrefix).
			HasComment("git oauth2 comment"),
		)
	})

	t.Run("create: git https api with private link", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(gitPrefix), true).
			WithGitHttpsApiPrivateLinkProviderParams(*sdk.NewGitHttpsApiPrivateLinkParamsRequest(true).
				WithAllowedAuthenticationSecrets(*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().WithAllSecrets(true))))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProviderNotEmpty().
			HasUsePrivatelinkEndpoint(true).
			HasAllowedPrefixes(gitPrefix),
		)
	})

	t.Run("create: external mcp with OAuth2 auth basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		auth := sdk.NewOAuth2McpUserAuthenticationRequest(oauthClientId, oauthClientSecret, oauthTokenEndpoint, oauthAuthorizationEndpoint)
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(mcpPrefix), true).
			WithExternalMcpOAuth2ProviderParams(*sdk.NewExternalMcpOAuth2ParamsRequest().WithApiUserAuthentication(*auth)))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProviderNotEmpty().
			HasAllowedPrefixes(mcpPrefix),
		)
	})

	t.Run("create: external mcp with OAuth2 auth all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		authMethod := sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost
		auth := sdk.NewOAuth2McpUserAuthenticationRequest(oauthClientId, oauthClientSecret, oauthTokenEndpoint, oauthAuthorizationEndpoint).
			WithOauthClientAuthMethod(authMethod).
			WithOauthDiscoveryUrl("https://auth.example.com/.well-known/openid-configuration").
			WithOauthRefreshTokenValidity(86400)
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(mcpPrefix), true).
			WithExternalMcpOAuth2ProviderParams(*sdk.NewExternalMcpOAuth2ParamsRequest().WithApiUserAuthentication(*auth)).
			WithComment("mcp oauth2 comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasAllowedPrefixes(mcpPrefix).
			HasComment("mcp oauth2 comment"),
		)
	})

	t.Run("create: external mcp with dynamic client auth", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		auth := sdk.NewDynamicClientMcpUserAuthenticationRequest(oauthResourceUrl)
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(mcpPrefix), true).
			WithExternalMcpDynamicClientProviderParams(*sdk.NewExternalMcpDynamicClientParamsRequest().WithApiUserAuthentication(*auth)))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, integration).
			HasEnabled(true).HasApiTypeExternalApi().HasCategoryApi().HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProviderNotEmpty().
			HasAllowedPrefixes(mcpPrefix),
		)
	})

	// -----------------------------------------------------------------------
	// ALTER SET
	// -----------------------------------------------------------------------

	t.Run("alter set: aws", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateAws(t).ID()

		otherRoleArn := "arn:aws:iam::000000000001:/role/other"
		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithAwsParams(*sdk.NewSetAwsApiParamsRequest().WithApiAwsRoleArn(otherRoleArn).WithApiKey("key")).
				WithEnabled(true).
				WithApiAllowedPrefixes(prefixes(awsOtherPrefix)).
				WithApiBlockedPrefixes(prefixes(awsPrefix)).
				WithComment("changed comment")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiKey("☺☺☺").
			HasApiProvider(strings.ToUpper(string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway))).
			HasApiAwsRoleArn(otherRoleArn).
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty().
			HasAllowedPrefixes(awsOtherPrefix).
			HasBlockedPrefixes(awsPrefix).
			HasComment("changed comment"),
		)
	})

	t.Run("alter set: azure", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateAzure(t).ID()

		otherAdApplicationId := "22222222-2222-2222-2222-222222222222"
		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithAzureParams(*sdk.NewSetAzureApiParamsRequest().WithAzureAdApplicationId(otherAdApplicationId).WithApiKey("key")).
				WithEnabled(true).
				WithApiAllowedPrefixes(prefixes(azureOtherPrefix)).
				WithApiBlockedPrefixes(prefixes(azurePrefix)).
				WithComment("changed comment")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiKey("☺☺☺").
			HasApiProvider("AZURE_API_MANAGEMENT").
			HasAzureTenantId(azureTenantId).
			HasAzureAdApplicationId(otherAdApplicationId).
			HasAzureMultiTenantAppNameNotEmpty().
			HasAzureConsentUrlNotEmpty().
			HasAllowedPrefixes(azureOtherPrefix).
			HasBlockedPrefixes(azurePrefix).
			HasComment("changed comment"),
		)
	})

	t.Run("alter set: azure - tenant id only", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateAzure(t).ID()

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithAzureParams(*sdk.NewSetAzureApiParamsRequest().WithAzureTenantId(azureOtherTenantId))))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasAzureTenantId(azureOtherTenantId),
		)
	})

	t.Run("alter set: google", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateGoogle(t).ID()

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithGoogleParams(*sdk.NewSetGoogleApiParamsRequest(googleOtherAudience)).
				WithEnabled(true).
				WithApiAllowedPrefixes(prefixes(googleOtherPrefix)).
				WithApiBlockedPrefixes(prefixes(googlePrefix)).
				WithComment("changed comment")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProvider("GOOGLE_API_GATEWAY").
			HasGoogleAudience(googleOtherAudience).
			HasAllowedPrefixes(googleOtherPrefix).
			HasBlockedPrefixes(googlePrefix).
			HasComment("changed comment"),
		)
	})

	t.Run("alter set: git token - change to no secrets", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateGitToken(t).ID()

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithGitHttpsApiTokenBasedParams(*sdk.NewSetGitHttpsApiTokenBasedParamsRequest().
					WithAllowedAuthenticationSecrets(*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().WithNoSecrets(true))).
				WithComment("updated")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasAllowedPrefixes(gitPrefix).
			HasComment("updated"),
		)
	})

	t.Run("alter set: git private link - disable privatelink", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateGitPrivateLink(t).ID()

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithGitHttpsApiPrivateLinkParams(*sdk.NewSetGitHttpsApiPrivateLinkParamsRequest().
					WithAllowedAuthenticationSecrets(*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().WithNoSecrets(true)).
					WithUsePrivatelinkEndpoint(false)).
				WithComment("updated")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasUsePrivatelinkEndpoint(false).
			HasAllowedPrefixes(gitPrefix).
			HasComment("updated"),
		)
	})

	t.Run("alter set: external mcp oauth2 - update credentials", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateMcpOAuth2(t).ID()

		newAuth := sdk.NewOAuth2McpUserAuthenticationRequest("new-id", "new-secret", oauthTokenEndpoint, oauthAuthorizationEndpoint).
			WithOauthRefreshTokenValidity(3600)
		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithExternalMcpOAuth2Params(*sdk.NewSetExternalMcpOAuth2ParamsRequest().WithApiUserAuthentication(*newAuth)).
				WithComment("updated")))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasAllowedPrefixes(mcpPrefix).
			HasComment("updated"),
		)
	})

	// -----------------------------------------------------------------------
	// ALTER UNSET
	// -----------------------------------------------------------------------

	t.Run("alter unset: aws - api_key, enabled, blocked_prefixes, comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(awsPrefix), true).
			WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, apiAwsRoleArn).WithApiKey("key")).
			WithApiBlockedPrefixes(prefixes(awsOtherPrefix)).
			WithComment("comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		err = client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithUnset(*sdk.NewApiIntegrationUnsetRequest().
				WithApiKey(true).
				WithEnabled(true).
				WithApiBlockedPrefixes(true).
				WithComment(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasApiKey("").
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("alter unset: azure - api_key, enabled, blocked_prefixes, comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(azurePrefix), true).
			WithAzureApiProviderParams(*sdk.NewAzureApiParamsRequest(azureTenantId, azureAdApplicationId).WithApiKey("key")).
			WithApiBlockedPrefixes(prefixes(azureOtherPrefix)).
			WithComment("comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		err = client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithUnset(*sdk.NewApiIntegrationUnsetRequest().
				WithApiKey(true).
				WithEnabled(true).
				WithApiBlockedPrefixes(true).
				WithComment(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasApiKey("").
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("alter unset: google - enabled, blocked_prefixes, comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(googlePrefix), true).
			WithGoogleApiProviderParams(*sdk.NewGoogleApiParamsRequest(googleAudience)).
			WithApiBlockedPrefixes(prefixes(googleOtherPrefix)).
			WithComment("comment"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		err = client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithUnset(*sdk.NewApiIntegrationUnsetRequest().
				WithEnabled(true).
				WithApiBlockedPrefixes(true).
				WithComment(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("alter unset: git token - allowed_authentication_secrets", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateGitToken(t).ID()

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithUnset(*sdk.NewApiIntegrationUnsetRequest().WithAllowedAuthenticationSecrets(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasAllowedPrefixes(gitPrefix),
		)
	})

	t.Run("alter unset: git private link - use_privatelink_endpoint", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateGitPrivateLink(t).ID()

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithUnset(*sdk.NewApiIntegrationUnsetRequest().WithUsePrivatelinkEndpoint(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasAllowedPrefixes(gitPrefix),
		)
	})

	// -----------------------------------------------------------------------
	// UNDOCUMENTED: api_blocked_prefixes support
	// -----------------------------------------------------------------------

	t.Run("undocumented: git can use api_blocked_prefixes on create", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		req := sdk.NewCreateApiIntegrationRequest(id, prefixes(gitPrefix), true).
			WithGitHttpsApiTokenBasedProviderParams(*sdk.NewGitHttpsApiTokenBasedParamsRequest().
				WithAllowedAuthenticationSecrets(*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().WithAllSecrets(true))).
			WithApiBlockedPrefixes(prefixes(gitOtherPrefix))
		err := client.ApiIntegrations.Create(ctx, req)
		if err != nil {
			t.Logf("[UNDOCUMENTED] git api_blocked_prefixes on create: NOT supported - %v", err)
			return
		}
		t.Logf("[UNDOCUMENTED] git api_blocked_prefixes on create: SUPPORTED")
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
		details, descErr := client.ApiIntegrations.DescribeGitHttpsApiDetails(ctx, id)
		require.NoError(t, descErr)
		t.Logf("[UNDOCUMENTED] git blocked_prefixes in describe: %v", details.BlockedPrefixes)
	})

	t.Run("undocumented: external_mcp can use api_blocked_prefixes on create", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		auth := sdk.NewDynamicClientMcpUserAuthenticationRequest(oauthResourceUrl)
		req := sdk.NewCreateApiIntegrationRequest(id, prefixes(mcpPrefix), true).
			WithExternalMcpDynamicClientProviderParams(*sdk.NewExternalMcpDynamicClientParamsRequest().WithApiUserAuthentication(*auth)).
			WithApiBlockedPrefixes(prefixes("https://mcp-blocked.example.com/"))
		err := client.ApiIntegrations.Create(ctx, req)
		if err != nil {
			t.Logf("[UNDOCUMENTED] external_mcp api_blocked_prefixes on create: NOT supported - %v", err)
			return
		}
		t.Logf("[UNDOCUMENTED] external_mcp api_blocked_prefixes on create: SUPPORTED")
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
	})

	// -----------------------------------------------------------------------
	// DROP
	// -----------------------------------------------------------------------

	t.Run("drop: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(awsPrefix), true).
			WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, apiAwsRoleArn)))
		require.NoError(t, err)

		err = client.ApiIntegrations.Drop(ctx, sdk.NewDropApiIntegrationRequest(id))
		require.NoError(t, err)

		_, err = client.ApiIntegrations.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		err := client.ApiIntegrations.Drop(ctx, sdk.NewDropApiIntegrationRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		err := client.ApiIntegrations.Drop(ctx, sdk.NewDropApiIntegrationRequest(id).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop safely: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, prefixes(awsPrefix), true).
			WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, apiAwsRoleArn)))
		require.NoError(t, err)

		err = client.ApiIntegrations.DropSafely(ctx, id)
		require.NoError(t, err)

		_, err = client.ApiIntegrations.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop safely: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		err := client.ApiIntegrations.DropSafely(ctx, id)
		require.NoError(t, err)
	})

	// -----------------------------------------------------------------------
	// SHOW
	// -----------------------------------------------------------------------

	t.Run("show: default", func(t *testing.T) {
		integrationAws := testClientHelper().ApiIntegration.CreateAws(t)
		integrationAzure := testClientHelper().ApiIntegration.CreateAzure(t)
		integrationGoogle := testClientHelper().ApiIntegration.CreateGoogle(t)

		returnedIntegrations, err := client.ApiIntegrations.Show(ctx, sdk.NewShowApiIntegrationRequest())
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *integrationAws)
		assert.Contains(t, returnedIntegrations, *integrationAzure)
		assert.Contains(t, returnedIntegrations, *integrationGoogle)
	})

	t.Run("show: with like filter", func(t *testing.T) {
		integrationAws := testClientHelper().ApiIntegration.CreateAws(t)
		integrationAzure := testClientHelper().ApiIntegration.CreateAzure(t)

		returnedIntegrations, err := client.ApiIntegrations.Show(ctx, sdk.NewShowApiIntegrationRequest().
			WithLike(sdk.Like{Pattern: sdk.String(integrationAws.Name)}))
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *integrationAws)
		assert.NotContains(t, returnedIntegrations, *integrationAzure)
	})

	t.Run("show by id safely: existing", func(t *testing.T) {
		integration := testClientHelper().ApiIntegration.CreateAws(t)

		result, err := client.ApiIntegrations.ShowByIDSafely(ctx, integration.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationFromObject(t, result).
			HasName(integration.Name).
			HasEnabled(true),
		)
	})

	t.Run("show by id safely: non-existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		_, err := client.ApiIntegrations.ShowByIDSafely(ctx, id)
		require.Error(t, err)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	// -----------------------------------------------------------------------
	// DESCRIBE
	// -----------------------------------------------------------------------

	t.Run("describe: aws - all fields", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateAws(t).ID()

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiKey("").
			HasApiProvider(strings.ToUpper(string(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway))).
			HasApiAwsRoleArn(apiAwsRoleArn).
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty().
			HasAllowedPrefixes(awsPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("describe: azure - all fields", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateAzure(t).ID()

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiKey("").
			HasApiProvider("AZURE_API_MANAGEMENT").
			HasAzureTenantId(azureTenantId).
			HasAzureAdApplicationId(azureAdApplicationId).
			HasAzureMultiTenantAppNameNotEmpty().
			HasAzureConsentUrlNotEmpty().
			HasAllowedPrefixes(azurePrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("describe: google - all fields", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateGoogle(t).ID()

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProvider("GOOGLE_API_GATEWAY").
			HasGoogleAudience(googleAudience).
			HasAllowedPrefixes(googlePrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("describe: git https api - all fields", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateGitToken(t).ID()

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProviderNotEmpty().
			HasAllowedPrefixes(gitPrefix).
			HasComment(""),
		)
	})

	t.Run("describe: external mcp - all fields", func(t *testing.T) {
		id := testClientHelper().ApiIntegration.CreateMcpOAuth2(t).ID()

		assertThatObject(t, objectassert.ApiIntegrationAllDetails(t, id).
			HasEnabled(true).
			HasApiProviderNotEmpty().
			HasAllowedPrefixes(mcpPrefix).
			HasComment(""),
		)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		_, err := client.ApiIntegrations.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
