//go:build non_account_level_tests

package testint

import (
	"strings"
	"testing"

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
	const apiAwsRoleArn = "arn:aws:iam::000000000001:/role/test"
	const azureTenantId = "00000000-0000-0000-0000-000000000000"
	const azureOtherTenantId = "11111111-1111-1111-1111-111111111111"
	const azureAdApplicationId = "11111111-1111-1111-1111-111111111111"
	const googleAudience = "api-gateway-id-123456.apigateway.gcp-project.cloud.goog"
	const googleOtherAudience = "api-gateway-id-666777.apigateway.gcp-project.cloud.goog"

	prefixes := func(prefix string) []sdk.ApiIntegrationEndpointPrefix {
		return []sdk.ApiIntegrationEndpointPrefix{{Path: prefix}}
	}

	assertApiIntegration := func(t *testing.T, s *sdk.ApiIntegration, name sdk.AccountObjectIdentifier, comment string) {
		t.Helper()
		assert.Equal(t, name.Name(), s.Name)
		assert.True(t, s.Enabled)
		assert.Equal(t, "EXTERNAL_API", s.ApiType)
		assert.Equal(t, "API", s.Category)
		assert.Equal(t, comment, s.Comment)
	}

	cleanupApiIntegrationProvider := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.ApiIntegrations.Drop(ctx, sdk.NewDropApiIntegrationRequest(id))
			require.NoError(t, err)
		}
	}

	createApiIntegrationAwsRequest := func(t *testing.T) *sdk.CreateApiIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateApiIntegrationRequest(id, prefixes(awsPrefix), true).
			WithAwsApiProviderParams(*sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiGateway, apiAwsRoleArn))
	}

	createApiIntegrationAzureRequest := func(t *testing.T) *sdk.CreateApiIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateApiIntegrationRequest(id, prefixes(azurePrefix), true).
			WithAzureApiProviderParams(*sdk.NewAzureApiParamsRequest(azureTenantId, azureAdApplicationId))
	}

	createApiIntegrationGoogleRequest := func(t *testing.T) *sdk.CreateApiIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateApiIntegrationRequest(id, prefixes(googlePrefix), true).
			WithGoogleApiProviderParams(*sdk.NewGoogleApiParamsRequest(googleAudience))
	}

	createApiIntegrationWithRequest := func(t *testing.T, request *sdk.CreateApiIntegrationRequest) *sdk.ApiIntegration {
		t.Helper()
		id := request.GetName()

		err := client.ApiIntegrations.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupApiIntegrationProvider(id))

		integration, err := client.ApiIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration
	}

	createAwsApiIntegration := func(t *testing.T) *sdk.ApiIntegration {
		t.Helper()
		return createApiIntegrationWithRequest(t, createApiIntegrationAwsRequest(t))
	}

	createAzureApiIntegration := func(t *testing.T) *sdk.ApiIntegration {
		t.Helper()
		return createApiIntegrationWithRequest(t, createApiIntegrationAzureRequest(t))
	}

	createGoogleApiIntegration := func(t *testing.T) *sdk.ApiIntegration {
		t.Helper()
		return createApiIntegrationWithRequest(t, createApiIntegrationGoogleRequest(t))
	}

	t.Run("create api integration: aws basic", func(t *testing.T) {
		request := createApiIntegrationAwsRequest(t)

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "")
	})

	t.Run("create api integration: azure basic", func(t *testing.T) {
		request := createApiIntegrationAzureRequest(t)

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "")
	})

	t.Run("create api integration: google basic", func(t *testing.T) {
		request := createApiIntegrationGoogleRequest(t)

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "")
	})

	t.Run("create api integration: aws more options", func(t *testing.T) {
		request := createApiIntegrationAwsRequest(t)

		request = request.
			WithAwsApiProviderParams(*request.AwsApiProviderParams.WithApiKey("key")).
			WithApiBlockedPrefixes(prefixes(awsOtherPrefix)).
			WithComment("comment")

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "comment")

		details, err := client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_KEY", Type: "String", Value: "☺☺☺", Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: awsOtherPrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "comment", Default: ""})
	})

	t.Run("create api integration: azure more options", func(t *testing.T) {
		request := createApiIntegrationAzureRequest(t)

		request = request.
			WithAzureApiProviderParams(*request.AzureApiProviderParams.WithApiKey("key")).
			WithApiBlockedPrefixes(prefixes(azureOtherPrefix)).
			WithComment("comment")

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "comment")

		details, err := client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_KEY", Type: "String", Value: "☺☺☺", Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: azureOtherPrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "comment", Default: ""})
	})

	t.Run("create api integration: google more options", func(t *testing.T) {
		request := createApiIntegrationGoogleRequest(t).
			WithApiBlockedPrefixes(prefixes(googleOtherPrefix)).
			WithComment("comment")

		integration := createApiIntegrationWithRequest(t, request)

		assertApiIntegration(t, integration, request.GetName(), "comment")

		details, err := client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: googleOtherPrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "comment", Default: ""})
	})

	t.Run("alter api integration: aws", func(t *testing.T) {
		integration := createAwsApiIntegration(t)

		otherRoleArn := "arn:aws:iam::000000000001:/role/other"
		setRequest := sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(
				*sdk.NewApiIntegrationSetRequest().
					WithAwsParams(*sdk.NewSetAwsApiParamsRequest().WithApiAwsRoleArn(otherRoleArn).WithApiKey("key")).
					WithEnabled(true).
					WithApiAllowedPrefixes(prefixes(awsOtherPrefix)).
					WithApiBlockedPrefixes(prefixes(awsPrefix)).
					WithComment("changed comment"),
			)
		err := client.ApiIntegrations.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_KEY", Type: "String", Value: "☺☺☺", Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_AWS_ROLE_ARN", Type: "String", Value: otherRoleArn, Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_ALLOWED_PREFIXES", Type: "List", Value: awsOtherPrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: awsPrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "changed comment", Default: ""})

		unsetRequest := sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithUnset(
				*sdk.NewApiIntegrationUnsetRequest().
					WithApiKey(true).
					WithEnabled(true).
					WithApiBlockedPrefixes(true).
					WithComment(true),
			)
		err = client.ApiIntegrations.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_KEY", Type: "String", Value: "", Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_AWS_ROLE_ARN", Type: "String", Value: otherRoleArn, Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_ALLOWED_PREFIXES", Type: "List", Value: awsOtherPrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: "", Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("alter api integration: azure", func(t *testing.T) {
		integration := createAzureApiIntegration(t)

		otherAdApplicationId := "22222222-2222-2222-2222-222222222222"
		setRequest := sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(
				*sdk.NewApiIntegrationSetRequest().
					WithAzureParams(*sdk.NewSetAzureApiParamsRequest().WithAzureAdApplicationId(otherAdApplicationId).WithApiKey("key")).
					WithEnabled(true).
					WithApiAllowedPrefixes(prefixes(azureOtherPrefix)).
					WithApiBlockedPrefixes(prefixes(azurePrefix)).
					WithComment("changed comment"),
			)
		err := client.ApiIntegrations.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_KEY", Type: "String", Value: "☺☺☺", Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "AZURE_AD_APPLICATION_ID", Type: "String", Value: otherAdApplicationId, Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_ALLOWED_PREFIXES", Type: "List", Value: azureOtherPrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: azurePrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "changed comment", Default: ""})

		unsetRequest := sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithUnset(
				*sdk.NewApiIntegrationUnsetRequest().
					WithApiKey(true).
					WithEnabled(true).
					WithApiBlockedPrefixes(true).
					WithComment(true),
			)
		err = client.ApiIntegrations.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_KEY", Type: "String", Value: "", Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "AZURE_AD_APPLICATION_ID", Type: "String", Value: otherAdApplicationId, Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_ALLOWED_PREFIXES", Type: "List", Value: azureOtherPrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: "", Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("alter api integration: azure - missing option", func(t *testing.T) {
		integration := createAzureApiIntegration(t)

		setRequest := sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(*sdk.NewApiIntegrationSetRequest().WithAzureParams(*sdk.NewSetAzureApiParamsRequest().WithAzureTenantId(azureOtherTenantId)))
		err := client.ApiIntegrations.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "AZURE_TENANT_ID", Type: "String", Value: azureOtherTenantId, Default: ""})
	})

	t.Run("alter api integration: google", func(t *testing.T) {
		integration := createGoogleApiIntegration(t)

		setRequest := sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(
				*sdk.NewApiIntegrationSetRequest().
					WithEnabled(true).
					WithApiAllowedPrefixes(prefixes(googleOtherPrefix)).
					WithApiBlockedPrefixes(prefixes(googlePrefix)).
					WithComment("changed comment"),
			)
		err := client.ApiIntegrations.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_ALLOWED_PREFIXES", Type: "List", Value: googleOtherPrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: googlePrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "changed comment", Default: ""})

		unsetRequest := sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithUnset(
				*sdk.NewApiIntegrationUnsetRequest().
					WithApiKey(true).
					WithEnabled(true).
					WithApiBlockedPrefixes(true).
					WithComment(true),
			)
		err = client.ApiIntegrations.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "false", Default: "false"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_ALLOWED_PREFIXES", Type: "List", Value: googleOtherPrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: "", Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("alter api integration: google - missing option", func(t *testing.T) {
		integration := createGoogleApiIntegration(t)

		setRequest := sdk.NewAlterApiIntegrationRequest(integration.ID()).
			WithSet(*sdk.NewApiIntegrationSetRequest().WithGoogleParams(*sdk.NewSetGoogleApiParamsRequest(googleOtherAudience)))
		err := client.ApiIntegrations.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "GOOGLE_AUDIENCE", Type: "String", Value: googleOtherAudience, Default: ""})
	})

	t.Run("drop api integration: existing", func(t *testing.T) {
		request := createApiIntegrationAwsRequest(t)
		id := request.GetName()

		err := client.ApiIntegrations.Create(ctx, request)
		require.NoError(t, err)

		err = client.ApiIntegrations.Drop(ctx, sdk.NewDropApiIntegrationRequest(id))
		require.NoError(t, err)

		_, err = client.ApiIntegrations.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop api integration: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier

		err := client.ApiIntegrations.Drop(ctx, sdk.NewDropApiIntegrationRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show api integration: default", func(t *testing.T) {
		integrationAws := createAwsApiIntegration(t)
		integrationAzure := createAzureApiIntegration(t)
		integrationGoogle := createGoogleApiIntegration(t)

		showRequest := sdk.NewShowApiIntegrationRequest()
		returnedIntegrations, err := client.ApiIntegrations.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *integrationAws)
		assert.Contains(t, returnedIntegrations, *integrationAzure)
		assert.Contains(t, returnedIntegrations, *integrationGoogle)
	})

	t.Run("show api integration: with options", func(t *testing.T) {
		integrationAws := createAwsApiIntegration(t)
		integrationAzure := createAzureApiIntegration(t)

		showRequest := sdk.NewShowApiIntegrationRequest().
			WithLike(sdk.Like{Pattern: &integrationAws.Name})
		returnedIntegrations, err := client.ApiIntegrations.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *integrationAws)
		assert.NotContains(t, returnedIntegrations, *integrationAzure)
	})

	t.Run("describe api integration: aws", func(t *testing.T) {
		integration := createAwsApiIntegration(t)

		details, err := client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_KEY", Type: "String", Value: "", Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_PROVIDER", Type: "String", Value: strings.ToUpper(string(sdk.ApiIntegrationAwsApiGateway)), Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_AWS_ROLE_ARN", Type: "String", Value: apiAwsRoleArn, Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_ALLOWED_PREFIXES", Type: "List", Value: awsPrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: "", Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("describe api integration: azure", func(t *testing.T) {
		integration := createAzureApiIntegration(t)

		details, err := client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "AZURE_TENANT_ID", Type: "String", Value: azureTenantId, Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "AZURE_AD_APPLICATION_ID", Type: "String", Value: azureAdApplicationId, Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_KEY", Type: "String", Value: "", Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_PROVIDER", Type: "String", Value: "AZURE_API_MANAGEMENT", Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_ALLOWED_PREFIXES", Type: "List", Value: azurePrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: "", Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("describe api integration: google", func(t *testing.T) {
		integration := createGoogleApiIntegration(t)

		details, err := client.ApiIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_KEY", Type: "String", Value: "", Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_PROVIDER", Type: "String", Value: "GOOGLE_API_GATEWAY", Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "GOOGLE_AUDIENCE", Type: "String", Value: googleAudience, Default: ""})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_ALLOWED_PREFIXES", Type: "List", Value: googlePrefix, Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "API_BLOCKED_PREFIXES", Type: "List", Value: "", Default: "[]"})
		assert.Contains(t, details, sdk.ApiIntegrationProperty{Name: "COMMENT", Type: "String", Value: "", Default: ""})
	})

	t.Run("describe api integration: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier

		_, err := client.ApiIntegrations.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
