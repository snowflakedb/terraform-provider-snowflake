//go:build !account_level_tests

package testint

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Secrets(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	integrationId := testClientHelper().Ids.RandomAccountObjectIdentifier()

	refreshTokenExpiryTime := time.Now().Add(24 * time.Hour).Format(time.DateOnly)

	_, apiIntegrationCleanup := testClientHelper().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "foo", "foo").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "foo"}, {Scope: "bar"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	stringDateToSnowflakeTimeFormat := func(inputLayout, date string) *time.Time {
		parsedTime, err := time.Parse(inputLayout, date)
		require.NoError(t, err)

		loc, err := time.LoadLocation("America/Los_Angeles")
		require.NoError(t, err)

		adjustedTime := parsedTime.In(loc)
		return &adjustedTime
	}

	type secretDetails struct {
		Name                        string
		Comment                     *string
		SecretType                  string
		Username                    *string
		OauthAccessTokenExpiryTime  *time.Time
		OauthRefreshTokenExpiryTime *time.Time
		OauthScopes                 []string
		IntegrationName             *string
	}

	assertSecretDetails := func(actual *sdk.SecretDetails, expected secretDetails) {
		assert.Equal(t, expected.Name, actual.Name)
		assert.EqualValues(t, expected.Comment, actual.Comment)
		assert.Equal(t, expected.SecretType, actual.SecretType)
		assert.EqualValues(t, expected.Username, actual.Username)
		assert.Equal(t, expected.OauthAccessTokenExpiryTime, actual.OauthAccessTokenExpiryTime)
		assert.Equal(t, expected.OauthRefreshTokenExpiryTime, actual.OauthRefreshTokenExpiryTime)
		assert.EqualValues(t, expected.OauthScopes, actual.OauthScopes)
		assert.EqualValues(t, expected.IntegrationName, actual.IntegrationName)
	}

	t.Run("Create: secretWithOAuthClientCredentialsFlow", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, integrationId).
			WithOauthScopes(sdk.OauthScopesListRequest{OauthScopesList: []sdk.ApiIntegrationScope{{Scope: "foo"}, {Scope: "bar"}}}).
			WithComment("a").
			WithIfNotExists(true)

		err := client.Secrets.CreateWithOAuthClientCredentialsFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		assertThatObject(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasComment("a").
				HasSecretType(string(sdk.SecretTypeOAuth2)).
				HasOauthScopes([]string{"foo", "bar"}).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:            id.Name(),
			Comment:         sdk.String("a"),
			SecretType:      string(sdk.SecretTypeOAuth2),
			OauthScopes:     []string{"foo", "bar"},
			IntegrationName: sdk.String(integrationId.Name()),
		})
	})

	// TODO [SNOW-1678767]: possible to create secret without specifying refresh token, token expiry time and OAuth scopes properties
	t.Run("Create: secretWithOAuth - minimal, without token and scopes", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, integrationId)

		err := client.Secrets.CreateWithOAuthClientCredentialsFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		assertThatObject(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:            id.Name(),
			SecretType:      string(sdk.SecretTypeOAuth2),
			IntegrationName: sdk.String(integrationId.Name()),
		})
	})

	// TODO [SNOW-1678756]: oauth_copes not inherited or not displayed correctly when not provided
	t.Run("Create: SecretWithOAuthClientCredentialsFlow - Empty Scopes List", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, integrationId).WithOauthScopes(sdk.OauthScopesListRequest{})

		err := client.Secrets.CreateWithOAuthClientCredentialsFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		assertThatObject(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasOauthScopes([]string{}).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)

		securityIntegrationProperties, _ := client.SecurityIntegrations.Describe(ctx, integrationId)
		assert.Contains(t, securityIntegrationProperties, sdk.SecurityIntegrationProperty{Name: "OAUTH_ALLOWED_SCOPES", Type: "List", Value: "[foo, bar]", Default: "[]"})

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assert.NotContains(t, details.OauthScopes, "foo")
		assert.NotContains(t, details.OauthScopes, "bar")
		assert.Equal(t, []string{}, details.OauthScopes)
	})

	t.Run("Create: SecretWithOAuthAuthorizationCodeFlow - refreshTokenExpiry date format", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithOAuthAuthorizationCodeFlowSecretRequest(id, "foo", refreshTokenExpiryTime, integrationId).
			WithComment("a").
			WithIfNotExists(true)

		err := client.Secrets.CreateWithOAuthAuthorizationCodeFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		_, err = client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasComment("a").
				HasSecretType(string(sdk.SecretTypeOAuth2)).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:                        id.Name(),
			SecretType:                  string(sdk.SecretTypeOAuth2),
			Comment:                     sdk.String("a"),
			OauthRefreshTokenExpiryTime: stringDateToSnowflakeTimeFormat(time.DateOnly, refreshTokenExpiryTime),
			IntegrationName:             sdk.String(integrationId.Name()),
		})
	})

	t.Run("Create: SecretWithOAuthAuthorizationCodeFlow - refreshTokenExpiry datetime format", func(t *testing.T) {
		refreshTokenWithTime := refreshTokenExpiryTime + " 12:00:00"
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		request := sdk.NewCreateWithOAuthAuthorizationCodeFlowSecretRequest(id, "foo", refreshTokenWithTime, integrationId)

		err := client.Secrets.CreateWithOAuthAuthorizationCodeFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:                        id.Name(),
			SecretType:                  string(sdk.SecretTypeOAuth2),
			OauthRefreshTokenExpiryTime: stringDateToSnowflakeTimeFormat(time.DateTime, refreshTokenWithTime),
			IntegrationName:             sdk.String(integrationId.Name()),
		})
	})

	t.Run("Create: WithBasicAuthentication", func(t *testing.T) {
		comment := random.Comment()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithBasicAuthenticationSecretRequest(id, "foo", "foo").
			WithComment(comment).
			WithIfNotExists(true)

		err := client.Secrets.CreateWithBasicAuthentication(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t,
			objectassert.SecretFromObject(t, secret).
				HasName(id.Name()).
				HasComment(comment).
				HasSecretType(string(sdk.SecretTypePassword)).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			Comment:    sdk.String(comment),
			SecretType: string(sdk.SecretTypePassword),
			Username:   sdk.String("foo"),
		})
	})

	t.Run("Create: WithBasicAuthentication - Empty Username and Password", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithBasicAuthenticationSecretRequest(id, "", "").
			WithIfNotExists(true)

		err := client.Secrets.CreateWithBasicAuthentication(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			SecretType: string(sdk.SecretTypePassword),
			Username:   sdk.String(""),
		})
	})

	t.Run("Create: WithGenericString", func(t *testing.T) {
		comment := random.Comment()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithGenericStringSecretRequest(id, "secret").
			WithComment(comment).
			WithIfNotExists(true)

		err := client.Secrets.CreateWithGenericString(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		_, err = client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasComment(comment).
				HasSecretType(string(sdk.SecretTypeGenericString)).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)
	})

	t.Run("Create: WithGenericString - empty secret_string", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithGenericStringSecretRequest(id, "")

		err := client.Secrets.CreateWithGenericString(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		assertThatObject(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasSecretType(string(sdk.SecretTypeGenericString)).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)
	})

	t.Run("Alter: SecretWithOAuthClientCredentials", func(t *testing.T) {
		comment := random.Comment()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, secretDropFunc := testClientHelper().Secret.CreateWithOAuthClientCredentialsFlow(t, id, integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}})
		t.Cleanup(secretDropFunc)

		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForFlow(*sdk.NewSetForFlowRequest().
						WithSetForOAuthClientCredentials(*sdk.NewSetForOAuthClientCredentialsRequest().
							WithOauthScopes(sdk.OauthScopesListRequest{OauthScopesList: []sdk.ApiIntegrationScope{{Scope: "foo"}, {Scope: "bar"}}}),
						),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:            id.Name(),
			SecretType:      string(sdk.SecretTypeOAuth2),
			Comment:         sdk.String(comment),
			OauthScopes:     []string{"foo", "bar"},
			IntegrationName: sdk.String(integrationId.Name()),
		})

		unsetRequest := sdk.NewAlterSecretRequest(id).
			WithUnset(
				*sdk.NewSecretUnsetRequest().
					WithComment(true),
			)
		err = client.Secrets.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assert.Empty(t, details.Comment)
	})

	t.Run("Alter: SecretWithOAuthAuthorizationCode", func(t *testing.T) {
		comment := random.Comment()
		alteredRefreshTokenExpiryTime := time.Now().Add(4 * 24 * time.Hour).Format(time.DateOnly)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, secretCleanup := testClientHelper().Secret.CreateWithOAuthAuthorizationCodeFlow(t, id, integrationId, "foo", refreshTokenExpiryTime)
		t.Cleanup(secretCleanup)

		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForFlow(*sdk.NewSetForFlowRequest().
						WithSetForOAuthAuthorization(*sdk.NewSetForOAuthAuthorizationRequest().
							WithOauthRefreshToken("bar").
							WithOauthRefreshTokenExpiryTime(alteredRefreshTokenExpiryTime),
						),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:                        id.Name(),
			SecretType:                  string(sdk.SecretTypeOAuth2),
			Comment:                     sdk.String(comment),
			OauthRefreshTokenExpiryTime: stringDateToSnowflakeTimeFormat(time.DateOnly, alteredRefreshTokenExpiryTime),
			IntegrationName:             sdk.String(integrationId.Name()),
		})

		unsetRequest := sdk.NewAlterSecretRequest(id).
			WithUnset(
				*sdk.NewSecretUnsetRequest().
					WithComment(true),
			)
		err = client.Secrets.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assert.Empty(t, details.Comment)
	})

	t.Run("Alter: SecretWithBasicAuthorization", func(t *testing.T) {
		comment := random.Comment()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		_, secretCleanup := testClientHelper().Secret.CreateWithBasicAuthenticationFlow(t, id, "foo", "foo")
		t.Cleanup(secretCleanup)

		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForFlow(*sdk.NewSetForFlowRequest().
						WithSetForBasicAuthentication(*sdk.NewSetForBasicAuthenticationRequest().
							WithUsername("bar").
							WithPassword("bar"),
						),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		// Cannot check password property since show and describe on secret do not have access to it
		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			SecretType: string(sdk.SecretTypePassword),
			Comment:    sdk.String(comment),
			Username:   sdk.String("bar"),
		})

		unsetRequest := sdk.NewAlterSecretRequest(id).
			WithUnset(
				*sdk.NewSecretUnsetRequest().
					WithComment(true),
			)
		err = client.Secrets.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assert.Empty(t, details.Comment)
	})

	t.Run("Alter: SecretWithGenericString", func(t *testing.T) {
		comment := random.Comment()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		_, secretCleanup := testClientHelper().Secret.CreateWithGenericString(t, id, "foo")
		t.Cleanup(secretCleanup)
		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForFlow(*sdk.NewSetForFlowRequest().
						WithSetForGenericString(*sdk.NewSetForGenericStringRequest().
							WithSecretString("bar"),
						),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		// Cannot assert secret string because snowflake does not output its value neither with SHOW nor DESCRIBE
		assert.Equal(t, comment, *details.Comment)

		unsetRequest := sdk.NewAlterSecretRequest(id).
			WithUnset(
				*sdk.NewSecretUnsetRequest().
					WithComment(true),
			)
		err = client.Secrets.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err = client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assert.Empty(t, details.Comment)
	})

	t.Run("Drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, secretCleanup := testClientHelper().Secret.CreateWithOAuthClientCredentialsFlow(t, id, integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}})
		t.Cleanup(secretCleanup)

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NotNil(t, secret)
		require.NoError(t, err)

		err = client.Secrets.Drop(ctx, sdk.NewDropSecretRequest(id))
		require.NoError(t, err)

		secret, err = client.Secrets.ShowByID(ctx, id)
		require.Nil(t, secret)
		require.Error(t, err)
	})

	t.Run("Drop: non-existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Secrets.Drop(ctx, sdk.NewDropSecretRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("Show", func(t *testing.T) {
		secretOAuthClientCredentials, secretCleanupClientCredentialsFlow := testClientHelper().Secret.CreateWithOAuthClientCredentialsFlow(t, testClientHelper().Ids.RandomSchemaObjectIdentifier(), integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}})
		t.Cleanup(secretCleanupClientCredentialsFlow)

		secretOAuthAuthorizationCode, secretCleanupAuthorizationCodeFlow := testClientHelper().Secret.CreateWithOAuthAuthorizationCodeFlow(t, testClientHelper().Ids.RandomSchemaObjectIdentifier(), integrationId, "foo", refreshTokenExpiryTime)
		t.Cleanup(secretCleanupAuthorizationCodeFlow)

		secretBasicAuthentication, secretCleanupBasicAuthentication := testClientHelper().Secret.CreateWithBasicAuthenticationFlow(t, testClientHelper().Ids.RandomSchemaObjectIdentifier(), "foo", "bar")
		t.Cleanup(secretCleanupBasicAuthentication)

		secretGenericString, secretCleanupGenericString := testClientHelper().Secret.CreateWithGenericString(t, testClientHelper().Ids.RandomSchemaObjectIdentifier(), "foo")
		t.Cleanup(secretCleanupGenericString)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest())
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secretOAuthClientCredentials)
		require.Contains(t, returnedSecrets, *secretOAuthAuthorizationCode)
		require.Contains(t, returnedSecrets, *secretBasicAuthentication)
		require.Contains(t, returnedSecrets, *secretGenericString)
	})

	t.Run("Show: SecretWithOAuthClientCredentialsFlow with Like", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		secret1, secretCleanup1 := testClientHelper().Secret.CreateWithOAuthClientCredentialsFlow(t, id1, integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}})
		t.Cleanup(secretCleanup1)

		secret2, secretCleanup2 := testClientHelper().Secret.CreateWithOAuthClientCredentialsFlow(t, id2, integrationId, []sdk.ApiIntegrationScope{{Scope: "bar"}})
		t.Cleanup(secretCleanup2)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.String(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithOAuthAuthorization with Like", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		secret1, secretCleanup1 := testClientHelper().Secret.CreateWithOAuthAuthorizationCodeFlow(t, id1, integrationId, "foo", refreshTokenExpiryTime)
		t.Cleanup(secretCleanup1)

		secret2, secretCleanup2 := testClientHelper().Secret.CreateWithOAuthAuthorizationCodeFlow(t, id2, integrationId, "bar", refreshTokenExpiryTime)
		t.Cleanup(secretCleanup2)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.String(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithBasicAuthentication with Like", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		secret1, secretCleanup1 := testClientHelper().Secret.CreateWithBasicAuthenticationFlow(t, id1, "foo", "foo")
		t.Cleanup(secretCleanup1)

		secret2, secretCleanup2 := testClientHelper().Secret.CreateWithBasicAuthenticationFlow(t, id2, "bar", "bar")
		t.Cleanup(secretCleanup2)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.String(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithGenericString with Like", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		secret1, secretCleanup1 := testClientHelper().Secret.CreateWithGenericString(t, id1, "foo")
		t.Cleanup(secretCleanup1)

		secret2, secretCleanup2 := testClientHelper().Secret.CreateWithGenericString(t, id2, "bar")
		t.Cleanup(secretCleanup2)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.String(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithOAuthClientCredentialsFlow with In", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		secret, secretCleanup := testClientHelper().Secret.CreateWithOAuthClientCredentialsFlow(t, id, integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}})
		t.Cleanup(secretCleanup)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Account: sdk.Pointer(true)}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Database: id.DatabaseId()}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)
	})

	t.Run("Show: SecretWithOAuthAuthorizationCodeFlow with In", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		secret, secretCleanup := testClientHelper().Secret.CreateWithOAuthAuthorizationCodeFlow(t, id, integrationId, "foo", refreshTokenExpiryTime)
		t.Cleanup(secretCleanup)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Account: sdk.Pointer(true)}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Database: id.DatabaseId()}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)
	})

	t.Run("Show: with In", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		secretOAuthClientCredentials, secretCleanupClientCredentialsFlow := testClientHelper().Secret.CreateWithOAuthClientCredentialsFlow(t, testClientHelper().Ids.RandomSchemaObjectIdentifier(), integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}})
		t.Cleanup(secretCleanupClientCredentialsFlow)

		secretOAuthAuthorizationCode, secretCleanupAuthorizationFlow := testClientHelper().Secret.CreateWithOAuthAuthorizationCodeFlow(t, testClientHelper().Ids.RandomSchemaObjectIdentifier(), integrationId, "foo", refreshTokenExpiryTime)
		t.Cleanup(secretCleanupAuthorizationFlow)

		secretBasicAuthentication, secretCleanupBasicAuthentication := testClientHelper().Secret.CreateWithBasicAuthenticationFlow(t, testClientHelper().Ids.RandomSchemaObjectIdentifier(), "foo", "foo")
		t.Cleanup(secretCleanupBasicAuthentication)

		secretGenericString, secretCleanupWithGenericString := testClientHelper().Secret.CreateWithGenericString(t, testClientHelper().Ids.RandomSchemaObjectIdentifier(), "foo")
		t.Cleanup(secretCleanupWithGenericString)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Account: sdk.Pointer(true)}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secretOAuthClientCredentials)
		require.Contains(t, returnedSecrets, *secretOAuthAuthorizationCode)
		require.Contains(t, returnedSecrets, *secretBasicAuthentication)
		require.Contains(t, returnedSecrets, *secretGenericString)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Database: id.DatabaseId()}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secretOAuthClientCredentials)
		require.Contains(t, returnedSecrets, *secretOAuthAuthorizationCode)
		require.Contains(t, returnedSecrets, *secretBasicAuthentication)
		require.Contains(t, returnedSecrets, *secretGenericString)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secretOAuthClientCredentials)
		require.Contains(t, returnedSecrets, *secretOAuthAuthorizationCode)
		require.Contains(t, returnedSecrets, *secretBasicAuthentication)
		require.Contains(t, returnedSecrets, *secretGenericString)
	})

	t.Run("Show: SecretWithGenericString with In", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		secret, secretCleanup := testClientHelper().Secret.CreateWithGenericString(t, id, "foo")
		t.Cleanup(secretCleanup)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Account: sdk.Pointer(true)}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Database: id.DatabaseId()}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)

		returnedSecrets, err = client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id.SchemaId()}}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret)
	})

	t.Run("ShowByID - same name different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		_, cleanup1 := testClientHelper().Secret.CreateWithGenericString(t, id1, "foo")
		t.Cleanup(cleanup1)

		_, cleanup2 := testClientHelper().Secret.CreateWithGenericString(t, id2, "bar")
		t.Cleanup(cleanup2)

		secretShowResult1, err := client.Secrets.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, secretShowResult1.ID())

		secretShowResult2, err := client.Secrets.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, secretShowResult2.ID())
	})

	t.Run("Describe", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		_, secretCleanup := testClientHelper().Secret.CreateWithGenericString(t, id, "foo")
		t.Cleanup(secretCleanup)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			SecretType: string(sdk.SecretTypeGenericString),
		})
	})
}
