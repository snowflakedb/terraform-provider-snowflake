package testint

import (
	"database/sql"
	assertions "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Secrets(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	integrationId := testClientHelper().Ids.RandomAccountObjectIdentifier()

	refreshTokenExpiryTime := time.Now().Add(24 * time.Hour).Format(time.DateOnly)

	cleanupIntegration := func(t *testing.T, integrationId sdk.AccountObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(integrationId).WithIfExists(true))
			require.NoError(t, err)
		}
	}

	err := client.SecurityIntegrations.CreateApiAuthenticationWithClientCredentialsFlow(
		ctx,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "foo", "foo").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "foo"}, {Scope: "bar"}}),
	)
	require.NoError(t, err)
	t.Cleanup(cleanupIntegration(t, integrationId))

	stringDateToSnowflakeTimeFormat := func(inputLayout, date string) *time.Time {
		parsedTime, err := time.Parse(inputLayout, date)
		require.NoError(t, err)

		loc, err := time.LoadLocation("America/Los_Angeles")
		require.NoError(t, err)

		adjustedTime := parsedTime.In(loc)
		return &adjustedTime
	}

	createSecretWithOAuthClientCredentialsFlow := func(t *testing.T, integrationId sdk.AccountObjectIdentifier, scopes []sdk.ApiIntegrationScope, with func(*sdk.CreateWithOAuthClientCredentialsFlowSecretRequest)) (*sdk.Secret, sdk.SchemaObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, integrationId, scopes)
		if with != nil {
			with(request)
		}
		err := client.Secrets.CreateWithOAuthClientCredentialsFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret, id
	}

	createSecretWithOAuthAuthorizationCodeFlow := func(t *testing.T, integrationId sdk.AccountObjectIdentifier, refreshToken, refreshTokenExpiryTime string, with func(*sdk.CreateWithOAuthAuthorizationCodeFlowSecretRequest)) (*sdk.Secret, sdk.SchemaObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithOAuthAuthorizationCodeFlowSecretRequest(id, refreshToken, refreshTokenExpiryTime, integrationId)
		if with != nil {
			with(request)
		}
		err := client.Secrets.CreateWithOAuthAuthorizationCodeFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret, id
	}

	createSecretWithBasicAuthentication := func(t *testing.T, username, password string, with func(*sdk.CreateWithBasicAuthenticationSecretRequest)) (*sdk.Secret, sdk.SchemaObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithBasicAuthenticationSecretRequest(id, username, password)
		if with != nil {
			with(request)
		}
		err := client.Secrets.CreateWithBasicAuthentication(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret, id
	}

	createSecretWithGenericString := func(t *testing.T, secretString string, with func(options *sdk.CreateWithGenericStringSecretRequest)) (*sdk.Secret, sdk.SchemaObjectIdentifier) {
		t.Helper()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithGenericStringSecretRequest(id, secretString)
		if with != nil {
			with(request)
		}
		err := client.Secrets.CreateWithGenericString(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret, id
	}

	createSecretWithId := func(t *testing.T, id sdk.SchemaObjectIdentifier) *sdk.Secret {
		t.Helper()
		request := sdk.NewCreateWithGenericStringSecretRequest(id, "foo")
		err := client.Secrets.CreateWithGenericString(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		secret, err := client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		return secret
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
		request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}, {Scope: "bar"}}).
			WithComment("a").
			WithIfNotExists(true)

		err := client.Secrets.CreateWithOAuthClientCredentialsFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		assertions.AssertThat(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasComment("a").
				HasSecretType("OAUTH2").
				HasOauthScopes("[foo, bar]").
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:            id.Name(),
			Comment:         sdk.String("a"),
			SecretType:      "OAUTH2",
			OauthScopes:     []string{"foo", "bar"},
			IntegrationName: sdk.String(integrationId.Name()),
		})
	})

	// It is possible to create secret without specifying both refresh token properties and scopes
	// Scopes are not being inherited from the security_integration what is tested further
	t.Run("Create: secretWithOAuth - minimal, without token and scopes", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, integrationId, nil)

		err := client.Secrets.CreateWithOAuthClientCredentialsFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		assertions.AssertThat(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:            id.Name(),
			SecretType:      "OAUTH2",
			IntegrationName: sdk.String(integrationId.Name()),
		})
	})

	// regarding the https://docs.snowflake.com/en/sql-reference/sql/create-secret secret should inherit security_integration scopes, but it does not do so
	t.Run("Create: SecretWithOAuthClientCredentialsFlow - No Scopes Specified", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateWithOAuthClientCredentialsFlowSecretRequest(id, integrationId, []sdk.ApiIntegrationScope{})

		err := client.Secrets.CreateWithOAuthClientCredentialsFlow(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Secret.DropFunc(t, id))

		assertions.AssertThat(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasOauthScopes("[]").
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)

		securityIntegrationProperties, _ := client.SecurityIntegrations.Describe(ctx, integrationId)
		assert.Contains(t, securityIntegrationProperties, sdk.SecurityIntegrationProperty{Name: "OAUTH_ALLOWED_SCOPES", Type: "List", Value: "[foo, bar]", Default: "[]"})

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)
		assert.NotEqual(t, details.OauthScopes, "[foo, bar]")
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

		assertions.AssertThat(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasComment("a").
				HasSecretType("OAUTH2").
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:                        id.Name(),
			SecretType:                  "OAUTH2",
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
			SecretType:                  "OAUTH2",
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

		_, err = client.Secrets.ShowByID(ctx, id)
		require.NoError(t, err)

		assertions.AssertThat(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasComment(comment).
				HasSecretType("PASSWORD").
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			Comment:    sdk.String(comment),
			SecretType: "PASSWORD",
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
			SecretType: "PASSWORD",
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

		assertions.AssertThat(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasComment(comment).
				HasSecretType("GENERIC_STRING").
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

		assertions.AssertThat(t,
			objectassert.Secret(t, id).
				HasName(id.Name()).
				HasSecretType("GENERIC_STRING").
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()),
		)
	})

	t.Run("Alter: SecretWithOAuthClientCredentials", func(t *testing.T) {
		comment := random.Comment()
		_, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}}, nil)
		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForOAuthClientCredentialsFlow(
						*sdk.NewSetForOAuthClientCredentialsFlowRequest(
							[]sdk.ApiIntegrationScope{{Scope: "foo"}, {Scope: "bar"}},
						),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:            id.Name(),
			SecretType:      "OAUTH2",
			Comment:         sdk.String(comment),
			OauthScopes:     sdk.String("[foo, bar]"),
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

		assert.Equal(t, details.Comment, sql.NullString{String: "", Valid: false})
	})

	t.Run("Alter: SecretWithOAuthAuthorizationCode", func(t *testing.T) {
		comment := random.Comment()
		alteredRefreshTokenExpiryTime := time.Now().Add(4 * 24 * time.Hour).Format(time.DateOnly)

		_, id := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo", refreshTokenExpiryTime, nil)
		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForOAuthAuthorizationFlow(
						*sdk.NewSetForOAuthAuthorizationFlowRequest().
							WithOauthRefreshToken("bar").
							WithOauthRefreshTokenExpiryTime(alteredRefreshTokenExpiryTime),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:                        id.Name(),
			SecretType:                  "OAUTH2",
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

		assert.Equal(t, details.Comment, sql.NullString{String: "", Valid: false})
	})

	t.Run("Alter: SecretWithBasicAuthorization", func(t *testing.T) {
		comment := random.Comment()

		_, id := createSecretWithBasicAuthentication(t, "foo", "foo", nil)
		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForBasicAuthentication(
						*sdk.NewSetForBasicAuthenticationRequest().
							WithUsername("bar").
							WithPassword("bar"),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		// Cannot check password property since show and describe on secret do not have access to it
		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			SecretType: "PASSWORD",
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

		assert.Equal(t, details.Comment, sql.NullString{String: "", Valid: false})
	})

	t.Run("Alter: SecretWithGenericString", func(t *testing.T) {
		comment := random.Comment()
		_, id := createSecretWithGenericString(t, "foo", nil)
		setRequest := sdk.NewAlterSecretRequest(id).
			WithSet(
				*sdk.NewSecretSetRequest().
					WithComment(comment).
					WithSetForGenericString(
						*sdk.NewSetForGenericStringRequest().
							WithSecretString("bar"),
					),
			)
		err := client.Secrets.Alter(ctx, setRequest)
		require.NoError(t, err)

		unsetRequest := sdk.NewAlterSecretRequest(id).
			WithUnset(
				*sdk.NewSecretUnsetRequest().
					WithComment(true),
			)
		err = client.Secrets.Alter(ctx, unsetRequest)
		require.NoError(t, err)

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, details.Comment, sql.NullString{String: "", Valid: false})
	})

	t.Run("Drop", func(t *testing.T) {
		_, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}}, nil)

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
		secretOAuthClientCredentials, _ := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}}, nil)
		secretOAuthAuthorizationCode, _ := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo", refreshTokenExpiryTime, nil)
		secretBasicAuthentication, _ := createSecretWithBasicAuthentication(t, "foo", "bar", nil)
		secretGenericString, _ := createSecretWithGenericString(t, "foo", nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest())
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secretOAuthClientCredentials)
		require.Contains(t, returnedSecrets, *secretOAuthAuthorizationCode)
		require.Contains(t, returnedSecrets, *secretBasicAuthentication)
		require.Contains(t, returnedSecrets, *secretGenericString)
	})

	t.Run("Show: SecretWithOAuthClientCredentialsFlow with Like", func(t *testing.T) {
		secret1, id1 := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}}, nil)
		secret2, _ := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.ApiIntegrationScope{{Scope: "bar"}}, nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.String(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithOAuthAuthorization with Like", func(t *testing.T) {
		secret1, id1 := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo_1", refreshTokenExpiryTime, nil)
		secret2, _ := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo_2", refreshTokenExpiryTime, nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.String(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithBasicAuthentication with Like", func(t *testing.T) {
		secret1, id1 := createSecretWithBasicAuthentication(t, "foo_1", "bar_1", nil)
		secret2, _ := createSecretWithBasicAuthentication(t, "foo_2", "bar_2", nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.String(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithGenericString with Like", func(t *testing.T) {
		secret1, id1 := createSecretWithGenericString(t, "foo_1", nil)
		secret2, _ := createSecretWithGenericString(t, "foo_2", nil)

		returnedSecrets, err := client.Secrets.Show(ctx, sdk.NewShowSecretRequest().WithLike(sdk.Like{
			Pattern: sdk.String(id1.Name()),
		}))
		require.NoError(t, err)
		require.Contains(t, returnedSecrets, *secret1)
		require.NotContains(t, returnedSecrets, *secret2)
	})

	t.Run("Show: SecretWithOAuthClientCredentialsFlow with In", func(t *testing.T) {
		secret, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}}, nil)

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
		secret, id := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo", refreshTokenExpiryTime, nil)

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
		secretOAuthClientCredentials, id := createSecretWithOAuthClientCredentialsFlow(t, integrationId, []sdk.ApiIntegrationScope{{Scope: "foo"}}, nil)
		secretOAuthAuthorizationCode, _ := createSecretWithOAuthAuthorizationCodeFlow(t, integrationId, "foo", refreshTokenExpiryTime, nil)
		secretBasicAuthentication, _ := createSecretWithBasicAuthentication(t, "foo", "bar", nil)
		secretGenericString, _ := createSecretWithGenericString(t, "foo", nil)

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
		secret, id := createSecretWithGenericString(t, "foo", nil)

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

		createSecretWithId(t, id1)
		createSecretWithId(t, id2)

		secretShowResult1, err := client.Secrets.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, secretShowResult1.ID())

		secretShowResult2, err := client.Secrets.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, secretShowResult2.ID())
	})

	t.Run("Describe", func(t *testing.T) {
		_, id := createSecretWithGenericString(t, "foo", func(req *sdk.CreateWithGenericStringSecretRequest) {
			req.WithComment("Lorem ipsum")
		})

		details, err := client.Secrets.Describe(ctx, id)
		require.NoError(t, err)

		assertSecretDetails(details, secretDetails{
			Name:       id.Name(),
			Comment:    sdk.String("Lorem ipsum"),
			SecretType: "GENERIC_STRING",
		})
	})
}
