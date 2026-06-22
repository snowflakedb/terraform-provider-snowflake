//go:build non_account_level_tests

package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
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
	const mcpOtherPrefix = "https://mcp.example.com/other/"
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

	t.Run("create: aws basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: awsPrefix}}, true).
			WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, apiAwsRoleArn)),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationAwsDetails(t, id).
			HasEnabled(true).
			HasApiKey("").
			HasApiProvider(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway).
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

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: awsPrefix}}, true).
			WithAwsApiProviderParams(
				*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, apiAwsRoleArn).
					WithApiKey("key"),
			).
			WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: awsOtherPrefix}}).
			WithComment("comment"),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment("comment"),
		)
		assertThatObject(t, objectassert.ApiIntegrationAwsDetails(t, id).
			HasEnabled(true).
			HasApiKey("☺☺☺").
			HasApiProvider(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway).
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
			t.Run(string(providerType), func(t *testing.T) {
				if providerType == sdk.ApiIntegrationAwsApiProviderTypeAwsGovApiGateway || providerType == sdk.ApiIntegrationAwsApiProviderTypeAwsGovPrivateApiGateway {
					t.Skip("gov provider types require a GovCloud Snowflake account")
				}
				if providerType == sdk.ApiIntegrationAwsApiProviderTypeAwsPrivateApiGateway && testenvs.GetSnowflakeEnvironmentWithProdDefault() == testenvs.SnowflakeProdEnvironment {
					t.Skip("private api gateway is not supported in prod environment as the feature can be only run on Business Critical Accounts")
				}

				id := testClientHelper().Ids.RandomAccountObjectIdentifier()

				err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: awsPrefix}}, true).
					WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(providerType, apiAwsRoleArn)),
				)
				require.NoError(t, err)
				t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

				assertThatObject(t, objectassert.ApiIntegrationAwsDetails(t, id).
					HasApiProvider(providerType),
				)
			})
		}
	})

	t.Run("create: azure basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: azurePrefix}}, true).
			WithAzureApiProviderParams(*sdk.NewAzureApiParamsRequest(azureTenantId, azureAdApplicationId)),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationAzureDetails(t, id).
			HasEnabled(true).
			HasApiKey("").
			HasApiProvider(sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement).
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

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: azurePrefix}}, true).
			WithAzureApiProviderParams(
				*sdk.NewAzureApiParamsRequest(azureTenantId, azureAdApplicationId).
					WithApiKey("key"),
			).
			WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: azureOtherPrefix}}).
			WithComment("comment"),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment("comment"),
		)
		assertThatObject(t, objectassert.ApiIntegrationAzureDetails(t, id).
			HasEnabled(true).
			HasApiKey("☺☺☺").
			HasApiProvider(sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement).
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

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: googlePrefix}}, true).
			WithGoogleApiProviderParams(*sdk.NewGoogleApiParamsRequest(googleAudience)),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationGoogleDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGoogleApiProviderTypeGoogleApiGateway).
			HasGoogleAudience(googleAudience).
			HasGoogleApiServiceAccountNotEmpty().
			HasAllowedPrefixes(googlePrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("create: google all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: googlePrefix}}, true).
			WithGoogleApiProviderParams(*sdk.NewGoogleApiParamsRequest(googleAudience)).
			WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: googleOtherPrefix}}).
			WithComment("comment"),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment("comment"),
		)
		assertThatObject(t, objectassert.ApiIntegrationGoogleDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGoogleApiProviderTypeGoogleApiGateway).
			HasGoogleAudience(googleAudience).
			HasGoogleApiServiceAccountNotEmpty().
			HasAllowedPrefixes(googlePrefix).
			HasBlockedPrefixes(googleOtherPrefix).
			HasComment("comment"),
		)
	})

	t.Run("create: git https api with token auth (all secrets)", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}, true).
			WithGitHttpsApiTokenBasedProviderParams(
				*sdk.NewGitHttpsApiTokenBasedParamsRequest().
					WithAllowedAuthenticationSecrets(
						*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().
							WithAllSecrets(true),
					),
			),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasUsePrivatelinkEndpoint(false).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitPrefix).
			HasNoBlockedPrefixes().
			HasComment("").
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("create: git https api with token auth (no secrets)", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}, true).
			WithGitHttpsApiTokenBasedProviderParams(
				*sdk.NewGitHttpsApiTokenBasedParamsRequest().
					WithAllowedAuthenticationSecrets(
						*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().
							WithNoSecrets(true),
					),
			),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasUsePrivatelinkEndpoint(false).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitPrefix).
			HasNoBlockedPrefixes().
			HasComment("").
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("create: git https api with token auth (dedicated secrets)", func(t *testing.T) {
		secretId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, secretCleanup := testClientHelper().Secret.CreateWithGenericString(t, secretId, "test_secret_string")
		t.Cleanup(secretCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}, true).
			WithGitHttpsApiTokenBasedProviderParams(
				*sdk.NewGitHttpsApiTokenBasedParamsRequest().
					WithAllowedAuthenticationSecrets(
						*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().
							WithAllowedList([]sdk.SchemaObjectIdentifier{secretId}),
					),
			),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasUsePrivatelinkEndpoint(false).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitPrefix).
			HasNoBlockedPrefixes().
			HasComment("").
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("create: git https api with GitHub App auth", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}, true).
			WithGitHttpsApiGithubAppProviderParams(*sdk.NewGitHttpsApiGithubAppParamsRequest()),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasUsePrivatelinkEndpoint(false).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitPrefix).
			HasNoBlockedPrefixes().
			HasComment("").
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeSnowflakeGithubApp).
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("create: git https api with OAuth2 auth basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		auth := sdk.NewOAuth2GitUserAuthenticationRequest(oauthAuthorizationEndpoint, oauthTokenEndpoint, oauthClientId, oauthClientSecret)

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}, true).
			WithGitHttpsApiOAuth2ProviderParams(
				*sdk.NewGitHttpsApiOAuth2ParamsRequest().
					WithApiUserAuthentication(*auth),
			),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasUsePrivatelinkEndpoint(false).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitPrefix).
			HasNoBlockedPrefixes().
			HasComment("").
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(oauthClientId).
			HasOauthTokenEndpoint(oauthTokenEndpoint).
			HasOauthAuthorizationEndpoint(oauthAuthorizationEndpoint).
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("create: git https api with OAuth2 auth all options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		auth := sdk.NewOAuth2GitUserAuthenticationRequest(oauthAuthorizationEndpoint, oauthTokenEndpoint, oauthClientId, oauthClientSecret).
			WithOauthAccessTokenValidity(3600).
			WithOauthRefreshTokenValidity(86400).
			WithOauthAllowedScopes([]sdk.ApiIntegrationOauthAllowedScopeItem{{Scope: sdk.ApiIntegrationOauthAllowedScopeReadApi}, {Scope: sdk.ApiIntegrationOauthAllowedScopeWriteRepository}}).
			WithOauthUsername("user@example.com")

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}, true).
			WithGitHttpsApiOAuth2ProviderParams(
				*sdk.NewGitHttpsApiOAuth2ParamsRequest().
					WithApiUserAuthentication(*auth),
			).
			WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitOtherPrefix}}).
			WithComment("git oauth2 comment"),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasUsePrivatelinkEndpoint(false).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitPrefix).
			HasBlockedPrefixes(gitOtherPrefix).
			HasComment("git oauth2 comment").
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(oauthClientId).
			HasOauthTokenEndpoint(oauthTokenEndpoint).
			HasOauthAuthorizationEndpoint(oauthAuthorizationEndpoint).
			HasOauthAccessTokenValidity(3600).
			HasOauthRefreshTokenValidity(86400).
			HasOauthAllowedScopes(sdk.ApiIntegrationOauthAllowedScopeReadApi, sdk.ApiIntegrationOauthAllowedScopeWriteRepository).
			HasOauthUsername("user@example.com").
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("create: git https api with private link", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}, true).
			WithGitHttpsApiPrivateLinkProviderParams(
				*sdk.NewGitHttpsApiPrivateLinkParamsRequest(true).
					WithAllowedAuthenticationSecrets(
						*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().
							WithAllSecrets(true),
					),
			),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasUsePrivatelinkEndpoint(true).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitPrefix).
			HasNoBlockedPrefixes().
			HasComment("").
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("create: git https api with private link all options", func(t *testing.T) {
		authSecretId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, authSecretCleanup := testClientHelper().Secret.CreateWithGenericString(t, authSecretId, "test_secret_string")
		t.Cleanup(authSecretCleanup)

		certSecretId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, certSecretCleanup := testClientHelper().Secret.CreateWithGenericString(t, certSecretId, random.GenerateX509(t))
		t.Cleanup(certSecretCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}, true).
			WithGitHttpsApiPrivateLinkProviderParams(
				*sdk.NewGitHttpsApiPrivateLinkParamsRequest(true).
					WithAllowedAuthenticationSecrets(
						*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().
							WithAllowedList([]sdk.SchemaObjectIdentifier{authSecretId}),
					).
					WithTlsTrustedCertificates([]sdk.SchemaObjectIdentifier{certSecretId}),
			).
			WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitOtherPrefix}}).
			WithComment("git private link comment"),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment("git private link comment"),
		)
		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasUsePrivatelinkEndpoint(true).
			HasTlsTrustedCertificates(fmt.Sprintf(`"%s"."%s".%s`, certSecretId.DatabaseName(), certSecretId.SchemaName(), certSecretId.Name())).
			HasAllowedPrefixes(gitPrefix).
			HasBlockedPrefixes(gitOtherPrefix).
			HasComment("git private link comment").
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("create: external mcp with OAuth2 auth basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		auth := sdk.NewOAuth2McpUserAuthenticationRequest(oauthClientId, oauthClientSecret, oauthTokenEndpoint, oauthAuthorizationEndpoint)

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: mcpPrefix}}, true).
			WithExternalMcpOAuth2ProviderParams(
				*sdk.NewExternalMcpOAuth2ParamsRequest().
					WithApiUserAuthentication(*auth),
			),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationExternalMcpDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp).
			HasAllowedPrefixes(mcpPrefix).
			HasNoBlockedPrefixes().
			HasComment("").
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(oauthClientId).
			HasOauthTokenEndpoint(oauthTokenEndpoint).
			HasOauthAuthorizationEndpoint(oauthAuthorizationEndpoint),
		)
	})

	t.Run("create: external mcp with OAuth2 auth all options", func(t *testing.T) {
		t.Skip("TODO(next prs): fix invalid parameter 'API_USER_AUTHENTICATION.OAUTH_DISCOVERY_URL' error when setting OAUTH_DISCOVERY_URL")

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		authMethod := sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost
		auth := sdk.NewOAuth2McpUserAuthenticationRequest(oauthClientId, oauthClientSecret, oauthTokenEndpoint, oauthAuthorizationEndpoint).
			WithOauthClientAuthMethod(authMethod).
			// TODO(next prs): fix invalid parameter 'API_USER_AUTHENTICATION.OAUTH_DISCOVERY_URL' error when setting OAUTH_DISCOVERY_URL
			// WithOauthDiscoveryUrl("https://auth.example.com/.well-known/openid-configuration").
			WithOauthRefreshTokenValidity(86400)

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: mcpPrefix}}, true).
			WithExternalMcpOAuth2ProviderParams(
				*sdk.NewExternalMcpOAuth2ParamsRequest().
					WithApiUserAuthentication(*auth),
			).
			WithComment("mcp oauth2 comment"),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegrationExternalMcpDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp).
			HasAllowedPrefixes(mcpPrefix).
			HasComment("mcp oauth2 comment").
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(oauthClientId).
			HasOauthClientAuthMethod(authMethod).
			HasOauthTokenEndpoint(oauthTokenEndpoint).
			HasOauthAuthorizationEndpoint(oauthAuthorizationEndpoint).
			HasOauthRefreshTokenValidity(86400),
		)
	})

	t.Run("create: external mcp with dynamic client auth", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		auth := sdk.NewDynamicClientMcpUserAuthenticationRequest("https://mcp.atlassian.com/v1/mcp")

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: "https://mcp.atlassian.com/v1/mcp"}}, true).
			WithExternalMcpDynamicClientProviderParams(
				*sdk.NewExternalMcpDynamicClientParamsRequest().
					WithApiUserAuthentication(*auth),
			),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		assertThatObject(t, objectassert.ApiIntegration(t, id).
			HasEnabled(true).
			HasApiTypeExternalApi().
			HasCategoryApi().
			HasComment(""),
		)
		assertThatObject(t, objectassert.ApiIntegrationExternalMcpDetails(t, id).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp).
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauthDynamicClient).
			HasAllowedPrefixes("https://mcp.atlassian.com/v1/mcp"),
		)
	})

	t.Run("alter set: aws", func(t *testing.T) {
		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateAws(t)
		t.Cleanup(integrationCleanup)

		otherRoleArn := "arn:aws:iam::000000000001:/role/other"

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithAwsParams(
					*sdk.NewSetAwsApiParamsRequest().
						WithApiAwsRoleArn(otherRoleArn).
						WithApiKey("key"),
				).
				WithEnabled(true).
				WithApiAllowedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: awsOtherPrefix}}).
				WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: awsPrefix}}).
				WithComment("changed comment")),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAwsDetails(t, integration.ID()).
			HasEnabled(true).
			HasApiKey("☺☺☺").
			HasApiProvider(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway).
			HasApiAwsRoleArn(otherRoleArn).
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty().
			HasAllowedPrefixes(awsOtherPrefix).
			HasBlockedPrefixes(awsPrefix).
			HasComment("changed comment"),
		)
	})

	t.Run("alter set: azure", func(t *testing.T) {
		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateAzure(t)
		t.Cleanup(integrationCleanup)

		otherAdApplicationId := "22222222-2222-2222-2222-222222222222"

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithAzureParams(
					*sdk.NewSetAzureApiParamsRequest().
						WithAzureTenantId(azureOtherTenantId).
						WithAzureAdApplicationId(otherAdApplicationId).
						WithApiKey("key"),
				).
				WithEnabled(true).
				WithApiAllowedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: azureOtherPrefix}}).
				WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: azurePrefix}}).
				WithComment("changed comment")),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAzureDetails(t, integration.ID()).
			HasEnabled(true).
			HasApiKey("☺☺☺").
			HasApiProvider(sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement).
			HasAzureTenantId(azureOtherTenantId).
			HasAzureAdApplicationId(otherAdApplicationId).
			HasAzureMultiTenantAppNameNotEmpty().
			HasAzureConsentUrlNotEmpty().
			HasAllowedPrefixes(azureOtherPrefix).
			HasBlockedPrefixes(azurePrefix).
			HasComment("changed comment"),
		)
	})

	t.Run("alter set: google", func(t *testing.T) {
		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateGoogle(t)
		t.Cleanup(integrationCleanup)

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithGoogleParams(*sdk.NewSetGoogleApiParamsRequest(googleOtherAudience)).
				WithEnabled(true).
				WithApiAllowedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: googleOtherPrefix}}).
				WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: googlePrefix}}).
				WithComment("changed comment")),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationGoogleDetails(t, integration.ID()).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGoogleApiProviderTypeGoogleApiGateway).
			HasGoogleAudience(googleOtherAudience).
			HasGoogleApiServiceAccountNotEmpty().
			HasAllowedPrefixes(googleOtherPrefix).
			HasBlockedPrefixes(googlePrefix).
			HasComment("changed comment"),
		)
	})

	t.Run("alter set: git token", func(t *testing.T) {
		secretId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, secretCleanup := testClientHelper().Secret.CreateWithGenericString(t, secretId, "test_secret_string")
		t.Cleanup(secretCleanup)

		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateGitToken(t)
		t.Cleanup(integrationCleanup)

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithGitHttpsApiTokenBasedParams(
					*sdk.NewSetGitHttpsApiTokenBasedParamsRequest().
						WithAllowedAuthenticationSecrets(
							*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().
								WithAllowedList([]sdk.SchemaObjectIdentifier{secretId}),
						),
				).
				WithEnabled(true).
				WithApiAllowedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitOtherPrefix}}).
				WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}).
				WithComment("changed comment")),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, integration.ID()).
			HasEnabled(true).
			HasUsePrivatelinkEndpoint(false).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitOtherPrefix).
			HasBlockedPrefixes(gitPrefix).
			HasComment("changed comment").
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("alter set: git private link", func(t *testing.T) {
		authSecretId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, authSecretCleanup := testClientHelper().Secret.CreateWithGenericString(t, authSecretId, "test_secret_string")
		t.Cleanup(authSecretCleanup)

		certSecretId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, certSecretCleanup := testClientHelper().Secret.CreateWithGenericString(t, certSecretId, random.GenerateX509(t))
		t.Cleanup(certSecretCleanup)

		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateGitPrivateLink(t)
		t.Cleanup(integrationCleanup)

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithGitHttpsApiPrivateLinkParams(
					*sdk.NewSetGitHttpsApiPrivateLinkParamsRequest().
						WithAllowedAuthenticationSecrets(
							*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().
								WithAllowedList([]sdk.SchemaObjectIdentifier{authSecretId}),
						).
						WithUsePrivatelinkEndpoint(true).
						WithTlsTrustedCertificates([]sdk.SchemaObjectIdentifier{certSecretId}),
				).
				WithEnabled(true).
				WithApiAllowedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitOtherPrefix}}).
				WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}).
				WithComment("changed comment")),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, integration.ID()).
			HasEnabled(true).
			HasUsePrivatelinkEndpoint(true).
			HasTlsTrustedCertificates(fmt.Sprintf(`"%s"."%s".%s`, certSecretId.DatabaseName(), certSecretId.SchemaName(), certSecretId.Name())).
			HasAllowedPrefixes(gitOtherPrefix).
			HasBlockedPrefixes(gitPrefix).
			HasComment("changed comment").
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("alter set: git oauth2", func(t *testing.T) {
		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateGitOAuth2(t)
		t.Cleanup(integrationCleanup)

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithEnabled(true).
				WithApiAllowedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitOtherPrefix}}).
				WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}).
				WithComment("changed comment")),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, integration.ID()).
			HasEnabled(true).
			HasUsePrivatelinkEndpoint(false).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitOtherPrefix).
			HasBlockedPrefixes(gitPrefix).
			HasComment("changed comment").
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(oauthClientId).
			HasOauthTokenEndpoint(oauthTokenEndpoint).
			HasOauthAuthorizationEndpoint(oauthAuthorizationEndpoint).
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("alter set: external mcp oauth2", func(t *testing.T) {
		t.Skip("TODO(next prs): fix invalid parameter 'API_USER_AUTHENTICATION.OAUTH_DISCOVERY_URL' error when setting OAUTH_DISCOVERY_URL")

		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateMcpOAuth2(t)
		t.Cleanup(integrationCleanup)

		authMethod := sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost
		newAuth := sdk.NewOAuth2McpUserAuthenticationRequest("new-id", "new-secret", oauthTokenEndpoint, oauthAuthorizationEndpoint).
			WithOauthClientAuthMethod(authMethod).
			// TODO(next prs): fix invalid parameter 'API_USER_AUTHENTICATION.OAUTH_DISCOVERY_URL' error when setting OAUTH_DISCOVERY_URL
			// WithOauthDiscoveryUrl("https://auth.example.com/.well-known/openid-configuration").
			WithOauthRefreshTokenValidity(3600)

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithExternalMcpOAuth2Params(
					*sdk.NewSetExternalMcpOAuth2ParamsRequest().
						WithApiUserAuthentication(*newAuth),
				).
				WithEnabled(true).
				WithApiAllowedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: mcpOtherPrefix}}).
				WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: mcpPrefix}}).
				WithComment("changed comment")),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationExternalMcpDetails(t, integration.ID()).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp).
			HasAllowedPrefixes(mcpOtherPrefix).
			HasBlockedPrefixes(mcpPrefix).
			HasComment("changed comment").
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId("new-id").
			HasOauthClientAuthMethod(authMethod).
			HasOauthTokenEndpoint(oauthTokenEndpoint).
			HasOauthAuthorizationEndpoint(oauthAuthorizationEndpoint).
			HasOauthRefreshTokenValidity(3600),
		)
	})

	t.Run("alter unset: aws - api_key, enabled, blocked_prefixes, comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: awsPrefix}}, true).
			WithAwsApiProviderParams(
				*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, apiAwsRoleArn).
					WithApiKey("key"),
			).
			WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: awsOtherPrefix}}).
			WithComment("comment"),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		err = client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithUnset(*sdk.NewApiIntegrationUnsetRequest().
				WithAwsParams(*sdk.NewUnsetAwsApiParamsRequest().WithApiKey(true)).
				WithEnabled(true).
				WithApiBlockedPrefixes(true).
				WithComment(true)),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAwsDetails(t, id).
			HasApiKey("").
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("alter unset: azure - api_key, enabled, blocked_prefixes, comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: azurePrefix}}, true).
			WithAzureApiProviderParams(
				*sdk.NewAzureApiParamsRequest(azureTenantId, azureAdApplicationId).
					WithApiKey("key"),
			).
			WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: azureOtherPrefix}}).
			WithComment("comment"),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		err = client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithUnset(*sdk.NewApiIntegrationUnsetRequest().
				WithAzureParams(*sdk.NewUnsetAzureApiParamsRequest().WithApiKey(true)).
				WithEnabled(true).
				WithApiBlockedPrefixes(true).
				WithComment(true)),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationAzureDetails(t, id).
			HasApiKey("").
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("alter unset: google - enabled, blocked_prefixes, comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: googlePrefix}}, true).
			WithGoogleApiProviderParams(*sdk.NewGoogleApiParamsRequest(googleAudience)).
			WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: googleOtherPrefix}}).
			WithComment("comment"),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		err = client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithUnset(*sdk.NewApiIntegrationUnsetRequest().
				WithEnabled(true).
				WithApiBlockedPrefixes(true).
				WithComment(true)),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationGoogleDetails(t, id).
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("alter unset: git token - allowed_authentication_secrets", func(t *testing.T) {
		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateGitToken(t)
		t.Cleanup(integrationCleanup)

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithUnset(*sdk.NewApiIntegrationUnsetRequest().
				WithGitHttpsApiTokenBasedParams(*sdk.NewUnsetGitHttpsApiTokenBasedParamsRequest().WithAllowedAuthenticationSecrets(true))),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, integration.ID()).
			HasEnabled(true).
			HasUsePrivatelinkEndpoint(false).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitPrefix).
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("alter unset: git private link - use_privatelink_endpoint, tls_trusted_certificates", func(t *testing.T) {
		certSecretId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, certSecretCleanup := testClientHelper().Secret.CreateWithGenericString(t, certSecretId, random.GenerateX509(t))
		t.Cleanup(certSecretCleanup)

		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateGitPrivateLink(t)
		t.Cleanup(integrationCleanup)

		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithGitHttpsApiPrivateLinkParams(
					*sdk.NewSetGitHttpsApiPrivateLinkParamsRequest().
						WithTlsTrustedCertificates([]sdk.SchemaObjectIdentifier{certSecretId}),
				)),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, integration.ID()).
			HasUsePrivatelinkEndpoint(true).
			HasTlsTrustedCertificates(fmt.Sprintf(`"%s"."%s".%s`, certSecretId.DatabaseName(), certSecretId.SchemaName(), certSecretId.Name())).
			HasAllowedAuthenticationSecrets(""),
		)

		err = client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithUnset(*sdk.NewApiIntegrationUnsetRequest().
				WithGitHttpsApiPrivateLinkParams(*sdk.NewUnsetGitHttpsApiPrivateLinkParamsRequest().
					WithUsePrivatelinkEndpoint(true).
					WithTlsTrustedCertificates(true))),
		)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, integration.ID()).
			HasEnabled(true).
			HasUsePrivatelinkEndpoint(false).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitPrefix).
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("undocumented: git can use api_blocked_prefixes on create", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: gitPrefix}}, true).
			WithGitHttpsApiTokenBasedProviderParams(
				*sdk.NewGitHttpsApiTokenBasedParamsRequest().
					WithAllowedAuthenticationSecrets(
						*sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest().
							WithAllSecrets(true),
					),
			).
			WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: gitOtherPrefix}}),
		)
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

		req := sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: mcpPrefix}}, true).
			WithExternalMcpDynamicClientProviderParams(
				*sdk.NewExternalMcpDynamicClientParamsRequest().
					WithApiUserAuthentication(*auth),
			).
			WithApiBlockedPrefixes([]sdk.ApiIntegrationEndpointPrefix{{Path: "https://mcp-blocked.example.com/"}})

		err := client.ApiIntegrations.Create(ctx, req)
		if err != nil {
			t.Logf("[UNDOCUMENTED] external_mcp api_blocked_prefixes on create: NOT supported - %v", err)
			return
		}
		t.Logf("[UNDOCUMENTED] external_mcp api_blocked_prefixes on create: SUPPORTED")
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))
	})

	// Prove that OAUTH_CLIENT_AUTH_METHOD and OAUTH_REFRESH_TOKEN_VALIDITY cannot be unset via ALTER once they have been set.
	// There is no UNSET clause for these fields in ApiIntegrationUnset, and omitting them from the auth block
	// in ALTER API INTEGRATION SET does not reset them - Snowflake silently retains the existing values.
	t.Run("external_mcp oauth2 optional auth fields cannot be unset via alter", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		authMethod := sdk.ApiIntegrationOauthClientAuthMethodClientSecretPost
		refreshTokenValidity := 3600

		auth := sdk.NewOAuth2McpUserAuthenticationRequest("oauth-client-id-123", "oauth-client-secret-456", oauthTokenEndpoint, oauthAuthorizationEndpoint).
			WithOauthClientAuthMethod(authMethod).
			WithOauthRefreshTokenValidity(refreshTokenValidity)
		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id,
			[]sdk.ApiIntegrationEndpointPrefix{{Path: mcpPrefix}}, true).
			WithExternalMcpOAuth2ProviderParams(*sdk.NewExternalMcpOAuth2ParamsRequest().WithApiUserAuthentication(*auth)),
		)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ApiIntegration.DropApiIntegrationFunc(t, id))

		// Confirm the initial state: both optional fields are set.
		assertThatObject(t, objectassert.ApiIntegrationExternalMcpDetails(t, id).
			HasOauthClientAuthMethod(authMethod).
			HasOauthRefreshTokenValidity(refreshTokenValidity),
		)

		// Attempt to "unset" by sending an ALTER with required fields only — OAUTH_CLIENT_AUTH_METHOD and
		// OAUTH_REFRESH_TOKEN_VALIDITY are intentionally absent. Snowflake accepts the request without error.
		authRequiredOnly := sdk.NewOAuth2McpUserAuthenticationRequest("oauth-client-id-123", "oauth-client-secret-456", oauthTokenEndpoint, oauthAuthorizationEndpoint)
		err = client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).
			WithSet(*sdk.NewApiIntegrationSetRequest().
				WithExternalMcpOAuth2Params(*sdk.NewSetExternalMcpOAuth2ParamsRequest().WithApiUserAuthentication(*authRequiredOnly))),
		)
		require.NoError(t, err)

		// Despite the ALTER succeeding, the values are unchanged — they cannot be removed this way.
		assertThatObject(t, objectassert.ApiIntegrationExternalMcpDetails(t, id).
			HasOauthClientAuthMethod(authMethod).
			HasOauthRefreshTokenValidity(refreshTokenValidity),
		)
	})

	t.Run("drop: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: awsPrefix}}, true).
			WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, apiAwsRoleArn)),
		)
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

		err := client.ApiIntegrations.Create(ctx, sdk.NewCreateApiIntegrationRequest(id, []sdk.ApiIntegrationEndpointPrefix{{Path: awsPrefix}}, true).
			WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway, apiAwsRoleArn)),
		)
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

	t.Run("show: default", func(t *testing.T) {
		integrationAws, integrationAwsCleanup := testClientHelper().ApiIntegration.CreateAws(t)
		t.Cleanup(integrationAwsCleanup)

		integrationAzure, integrationAzureCleanup := testClientHelper().ApiIntegration.CreateAzure(t)
		t.Cleanup(integrationAzureCleanup)

		integrationGoogle, integrationGoogleCleanup := testClientHelper().ApiIntegration.CreateGoogle(t)
		t.Cleanup(integrationGoogleCleanup)

		returnedIntegrations, err := client.ApiIntegrations.Show(ctx, sdk.NewShowApiIntegrationRequest())
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *integrationAws)
		assert.Contains(t, returnedIntegrations, *integrationAzure)
		assert.Contains(t, returnedIntegrations, *integrationGoogle)
	})

	t.Run("show: with like filter", func(t *testing.T) {
		integrationAws, integrationAwsCleanup := testClientHelper().ApiIntegration.CreateAws(t)
		t.Cleanup(integrationAwsCleanup)

		integrationAzure, integrationAzureCleanup := testClientHelper().ApiIntegration.CreateAzure(t)
		t.Cleanup(integrationAzureCleanup)

		returnedIntegrations, err := client.ApiIntegrations.Show(ctx, sdk.NewShowApiIntegrationRequest().
			WithLike(sdk.Like{Pattern: sdk.String(integrationAws.Name)}))
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *integrationAws)
		assert.NotContains(t, returnedIntegrations, *integrationAzure)
	})

	t.Run("show by id safely: existing", func(t *testing.T) {
		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateAws(t)
		t.Cleanup(integrationCleanup)

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

	t.Run("describe: aws - all fields", func(t *testing.T) {
		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateAws(t)
		t.Cleanup(integrationCleanup)

		assertThatObject(t, objectassert.ApiIntegrationAwsDetails(t, integration.ID()).
			HasEnabled(true).
			HasApiKey("").
			HasApiProvider(sdk.ApiIntegrationAwsApiProviderTypeAwsApiGateway).
			HasApiAwsRoleArn(apiAwsRoleArn).
			HasApiAwsIamUserArnNotEmpty().
			HasApiAwsExternalIdNotEmpty().
			HasAllowedPrefixes(awsPrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("describe: azure - all fields", func(t *testing.T) {
		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateAzure(t)
		t.Cleanup(integrationCleanup)

		assertThatObject(t, objectassert.ApiIntegrationAzureDetails(t, integration.ID()).
			HasEnabled(true).
			HasApiKey("").
			HasApiProvider(sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement).
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
		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateGoogle(t)
		t.Cleanup(integrationCleanup)

		assertThatObject(t, objectassert.ApiIntegrationGoogleDetails(t, integration.ID()).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGoogleApiProviderTypeGoogleApiGateway).
			HasGoogleAudience(googleAudience).
			HasGoogleApiServiceAccountNotEmpty().
			HasAllowedPrefixes(googlePrefix).
			HasNoBlockedPrefixes().
			HasComment(""),
		)
	})

	t.Run("describe: git https api - all fields", func(t *testing.T) {
		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateGitToken(t)
		t.Cleanup(integrationCleanup)

		assertThatObject(t, objectassert.ApiIntegrationGitHttpsApiDetails(t, integration.ID()).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationGitApiProviderTypeGitHttpsApi).
			HasUsePrivatelinkEndpoint(false).
			HasNoTlsTrustedCertificates().
			HasAllowedPrefixes(gitPrefix).
			HasNoBlockedPrefixes().
			HasComment("").
			HasNoUserAuthType().
			HasAllowedAuthenticationSecrets(""),
		)
	})

	t.Run("describe: external mcp - all fields", func(t *testing.T) {
		integration, integrationCleanup := testClientHelper().ApiIntegration.CreateMcpOAuth2(t)
		t.Cleanup(integrationCleanup)

		assertThatObject(t, objectassert.ApiIntegrationExternalMcpDetails(t, integration.ID()).
			HasEnabled(true).
			HasApiProvider(sdk.ApiIntegrationMcpApiProviderTypeExternalMcp).
			HasAllowedPrefixes(mcpPrefix).
			HasComment("").
			HasUserAuthType(sdk.ApiIntegrationUserAuthTypeOauth2).
			HasOauthGrant("AUTHORIZATION_CODE").
			HasOauthClientId(oauthClientId).
			HasOauthTokenEndpoint(oauthTokenEndpoint).
			HasOauthAuthorizationEndpoint(oauthAuthorizationEndpoint),
		)
	})

	t.Run("describe: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier

		_, err := client.ApiIntegrations.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
