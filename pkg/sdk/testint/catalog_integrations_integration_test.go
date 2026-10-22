//go:build non_account_level_tests

package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_CatalogIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	const (
		glueAwsRoleArn              = "arn:aws:iam::123456789012:role/sqsAccess"
		glueCatalogId               = "123456789012"
		glueRegion                  = "us-east-2"
		catalogName                 = "TEST_CATALOG"
		catalogNamespace            = "TEST_NAMESPACE"
		oAuthAllowedScope           = "PRINCIPAL_ROLE:PRINCIPAL"
		additionalOAuthAllowedScope = "PRINCIPAL_ROLE:SECONDARY"
	)

	loadOpenCatalogConfig := func(t *testing.T) (catalogUri, privateCatalogUri, oAuthClientId, oAuthClientSecret string) {
		t.Helper()
		accountLocator := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogAccountLocator)
		oAuthClientId = testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientId)
		oAuthClientSecret = testenvs.GetOrSkipTest(t, testenvs.OpenCatalogPrimaryOAuthClientSecret)
		catalogUri = fmt.Sprintf("https://%s.snowflakecomputing.com/polaris/api/catalog", accountLocator)
		privateCatalogUri = fmt.Sprintf("https://%s.privatelink.snowflakecomputing.com/polaris/api/catalog", accountLocator)
		return
	}

	assertCatalogIntegration := func(t *testing.T, s *sdk.CatalogIntegration, name sdk.AccountObjectIdentifier, comment string) {
		t.Helper()
		assertThatObject(t, objectassert.CatalogIntegrationFromObject(t, s).
			HasName(name.Name()).
			HasEnabled(false).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasComment(comment))
	}

	createCatalogIntegrationWithRequest := func(t *testing.T, request *sdk.CreateCatalogIntegrationRequest) *sdk.CatalogIntegration {
		t.Helper()
		id := request.GetName()

		err := client.CatalogIntegrations.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().CatalogIntegration.DropFunc(t, id))

		integration, err := client.CatalogIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration
	}

	createCatalogIntegrationAwsGlueRequest := func(t *testing.T) *sdk.CreateCatalogIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateCatalogIntegrationRequest(id, false).
			WithAwsGlueCatalogSourceParams(*sdk.NewAwsGlueParamsRequest(glueAwsRoleArn, glueCatalogId).
				WithGlueRegion(glueRegion))
	}

	createCatalogIntegrationObjectStorageRequest := func(t *testing.T) *sdk.CreateCatalogIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateCatalogIntegrationRequest(id, false).
			WithObjectStorageCatalogSourceParams(*sdk.NewObjectStorageParamsRequest(sdk.CatalogIntegrationTableFormatDelta))
	}

	createCatalogIntegrationOpenCatalogRequest := func(t *testing.T, catalogUri, oAuthClientId, oAuthClientSecret string) *sdk.CreateCatalogIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateCatalogIntegrationRequest(id, false).
			WithOpenCatalogCatalogSourceParams(*sdk.NewOpenCatalogParamsRequest().
				WithRestConfig(sdk.OpenCatalogRestConfigRequest{
					CatalogUri:  catalogUri,
					CatalogName: catalogName,
				}).
				WithRestAuthentication(sdk.OAuthRestAuthenticationRequest{
					OauthClientId:      oAuthClientId,
					OauthClientSecret:  oAuthClientSecret,
					OauthAllowedScopes: []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}},
				}))
	}

	createCatalogIntegrationIcebergRestRequest := func(t *testing.T, catalogUri, oAuthClientId, oAuthClientSecret string) *sdk.CreateCatalogIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateCatalogIntegrationRequest(id, false).
			WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
				WithRestConfig(sdk.IcebergRestRestConfigRequest{
					CatalogUri:  catalogUri,
					CatalogName: new(catalogName),
				}).
				WithOAuthRestAuthentication(sdk.OAuthRestAuthenticationRequest{
					OauthClientId:      oAuthClientId,
					OauthClientSecret:  oAuthClientSecret,
					OauthAllowedScopes: []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}},
				}))
	}

	createAwsGlueCatalogIntegration := func(t *testing.T) *sdk.CatalogIntegration {
		t.Helper()
		return createCatalogIntegrationWithRequest(t, createCatalogIntegrationAwsGlueRequest(t))
	}

	createObjectStorageCatalogIntegration := func(t *testing.T) *sdk.CatalogIntegration {
		t.Helper()
		return createCatalogIntegrationWithRequest(t, createCatalogIntegrationObjectStorageRequest(t))
	}

	createOpenCatalogCatalogIntegration := func(t *testing.T) *sdk.CatalogIntegration {
		t.Helper()
		catalogUri, _, oAuthClientId, oAuthClientSecret := loadOpenCatalogConfig(t)
		return createCatalogIntegrationWithRequest(t, createCatalogIntegrationOpenCatalogRequest(t, catalogUri, oAuthClientId, oAuthClientSecret))
	}

	createIcebergRestCatalogIntegration := func(t *testing.T) *sdk.CatalogIntegration {
		t.Helper()
		catalogUri, _, oAuthClientId, oAuthClientSecret := loadOpenCatalogConfig(t)
		return createCatalogIntegrationWithRequest(t, createCatalogIntegrationIcebergRestRequest(t, catalogUri, oAuthClientId, oAuthClientSecret))
	}

	t.Run("create catalog integration: AWS Glue basic", func(t *testing.T) {
		request := createCatalogIntegrationAwsGlueRequest(t)

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "")
		assertThatObject(t, objectassert.CatalogIntegrationAwsGlueDetails(t, integration.ID()).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeGlue).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(glueRegion).
			HasCatalogNamespace("").
			HasGlueAwsIamUserArnNotEmpty().
			HasGlueAwsExternalIdNotEmpty())
	})

	t.Run("create catalog integration: object storage basic", func(t *testing.T) {
		request := createCatalogIntegrationObjectStorageRequest(t)

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "")
		assertThatObject(t, objectassert.CatalogIntegrationObjectStorageDetails(t, integration.ID()).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeObjectStore).
			HasTableFormat(sdk.CatalogIntegrationTableFormatDelta).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment(""))
	})

	t.Run("create catalog integration: Open Catalog basic", func(t *testing.T) {
		catalogUri, _, oAuthClientId, oAuthClientSecret := loadOpenCatalogConfig(t)
		request := createCatalogIntegrationOpenCatalogRequest(t, catalogUri, oAuthClientId, oAuthClientSecret)

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "")

		openCatalogDetails, err := client.CatalogIntegrations.DescribeOpenCatalogDetails(ctx, integration.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.CatalogIntegrationOpenCatalogDetailsFromObject(t, openCatalogDetails).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypePolaris).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasCatalogNamespace("").
			HasRestConfigWith(objectassert.NewOpenCatalogRestConfigDetailsAssert().
				HasCatalogUri(catalogUri).
				HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
				HasCatalogName(catalogName).
				HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)).
			HasRestAuthenticationWith(objectassert.NewOAuthRestAuthenticationDetailsAssert().
				HasOauthTokenUri(catalogUri+"/v1/oauth/tokens").
				HasOauthClientId(oAuthClientId).
				HasOauthAllowedScopes(oAuthAllowedScope)))
	})

	t.Run("create catalog integration: Iceberg REST basic", func(t *testing.T) {
		catalogUri, _, oAuthClientId, oAuthClientSecret := loadOpenCatalogConfig(t)
		request := createCatalogIntegrationIcebergRestRequest(t, catalogUri, oAuthClientId, oAuthClientSecret)

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "")

		icebergRestDetails, err := client.CatalogIntegrations.DescribeIcebergRestDetails(ctx, integration.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.CatalogIntegrationIcebergRestDetailsFromObject(t, icebergRestDetails).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasCatalogNamespace("").
			HasOAuthRestAuthenticationWith(objectassert.NewOAuthRestAuthenticationDetailsAssert().
				HasOauthTokenUri(catalogUri+"/v1/oauth/tokens").
				HasOauthClientId(oAuthClientId).
				HasOauthAllowedScopes(oAuthAllowedScope)).
			HasRestConfigWith(objectassert.NewIcebergRestRestConfigDetailsAssert().
				HasCatalogUri(catalogUri).
				HasPrefix("").
				HasCatalogName(catalogName).
				HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
				HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)))
	})

	t.Run("create catalog integration: AWS Glue all options", func(t *testing.T) {
		request := createCatalogIntegrationAwsGlueRequest(t).
			WithIfNotExists(true).
			WithAwsGlueCatalogSourceParams(*sdk.NewAwsGlueParamsRequest(glueAwsRoleArn, glueCatalogId).
				WithGlueRegion(glueRegion).
				WithCatalogNamespace(catalogNamespace)).
			WithRefreshIntervalSeconds(120).
			WithComment("test comment")

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "test comment")
		assertThatObject(t, objectassert.CatalogIntegrationAwsGlueDetails(t, integration.ID()).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeGlue).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(120).
			HasComment("test comment").
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(glueRegion).
			HasCatalogNamespace(catalogNamespace).
			HasGlueAwsIamUserArnNotEmpty().
			HasGlueAwsExternalIdNotEmpty())
	})

	t.Run("create catalog integration: object storage all options", func(t *testing.T) {
		request := createCatalogIntegrationObjectStorageRequest(t).
			WithIfNotExists(true).
			WithRefreshIntervalSeconds(120).
			WithComment("test comment")

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "test comment")
		assertThatObject(t, objectassert.CatalogIntegrationObjectStorageDetails(t, integration.ID()).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeObjectStore).
			HasTableFormat(sdk.CatalogIntegrationTableFormatDelta).
			HasEnabled(false).
			HasRefreshIntervalSeconds(120).
			HasComment("test comment"))
	})

	t.Run("create catalog integration: Open Catalog all options", func(t *testing.T) {
		_, privateCatalogUri, oAuthClientId, oAuthClientSecret := loadOpenCatalogConfig(t)
		privateOAuthTokenUri := privateCatalogUri + "/v1/oauth/tokens"
		request := createCatalogIntegrationOpenCatalogRequest(t, privateCatalogUri, oAuthClientId, oAuthClientSecret).
			WithIfNotExists(true).
			WithOpenCatalogCatalogSourceParams(*sdk.NewOpenCatalogParamsRequest().
				WithCatalogNamespace(catalogNamespace).
				WithRestConfig(sdk.OpenCatalogRestConfigRequest{
					CatalogUri:           privateCatalogUri,
					CatalogApiType:       sdk.Pointer(sdk.CatalogIntegrationCatalogApiTypePrivate),
					CatalogName:          catalogName,
					AccessDelegationMode: sdk.Pointer(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials),
				}).
				WithRestAuthentication(sdk.OAuthRestAuthenticationRequest{
					OauthTokenUri:      sdk.String(privateOAuthTokenUri),
					OauthClientId:      oAuthClientId,
					OauthClientSecret:  oAuthClientSecret,
					OauthAllowedScopes: []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}, {Value: additionalOAuthAllowedScope}},
				})).
			WithRefreshIntervalSeconds(120).
			WithComment("test comment")

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "test comment")

		openCatalogDetails, err := client.CatalogIntegrations.DescribeOpenCatalogDetails(ctx, integration.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.CatalogIntegrationOpenCatalogDetailsFromObject(t, openCatalogDetails).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypePolaris).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(120).
			HasComment("test comment").
			HasCatalogNamespace(catalogNamespace).
			HasRestConfigWith(objectassert.NewOpenCatalogRestConfigDetailsAssert().
				HasCatalogUri(privateCatalogUri).
				HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePrivate).
				HasCatalogName(catalogName).
				HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials)).
			HasRestAuthenticationWith(objectassert.NewOAuthRestAuthenticationDetailsAssert().
				HasOauthTokenUri(privateOAuthTokenUri).
				HasOauthClientId(oAuthClientId).
				HasOauthAllowedScopes(oAuthAllowedScope, additionalOAuthAllowedScope)))
	})

	t.Run("create catalog integration: Iceberg REST all options", func(t *testing.T) {
		catalogUri, _, oAuthClientId, oAuthClientSecret := loadOpenCatalogConfig(t)
		oAuthTokenUri := catalogUri + "/v1/oauth/tokens"
		const prefix = "prefix"
		request := createCatalogIntegrationIcebergRestRequest(t, catalogUri, oAuthClientId, oAuthClientSecret).
			WithIfNotExists(true).
			WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
				WithCatalogNamespace(catalogNamespace).
				WithRestConfig(sdk.IcebergRestRestConfigRequest{
					CatalogUri:           catalogUri,
					Prefix:               sdk.String(prefix),
					CatalogName:          sdk.String(catalogName),
					CatalogApiType:       sdk.Pointer(sdk.CatalogIntegrationCatalogApiTypePublic),
					AccessDelegationMode: sdk.Pointer(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials),
				}).
				WithOAuthRestAuthentication(sdk.OAuthRestAuthenticationRequest{
					OauthTokenUri:      new(oAuthTokenUri),
					OauthClientId:      oAuthClientId,
					OauthClientSecret:  oAuthClientSecret,
					OauthAllowedScopes: []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}, {Value: additionalOAuthAllowedScope}},
				})).
			WithRefreshIntervalSeconds(120).
			WithComment("test comment")

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "test comment")

		icebergRestDetails, err := client.CatalogIntegrations.DescribeIcebergRestDetails(ctx, integration.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.CatalogIntegrationIcebergRestDetailsFromObject(t, icebergRestDetails).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(120).
			HasComment("test comment").
			HasCatalogNamespace(catalogNamespace).
			HasOAuthRestAuthenticationWith(objectassert.NewOAuthRestAuthenticationDetailsAssert().
				HasOauthTokenUri(oAuthTokenUri).
				HasOauthClientId(oAuthClientId).
				HasOauthAllowedScopes(oAuthAllowedScope, additionalOAuthAllowedScope)).
			HasRestConfigWith(objectassert.NewIcebergRestRestConfigDetailsAssert().
				HasCatalogUri(catalogUri).
				HasPrefix(prefix).
				HasCatalogName(catalogName).
				HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
				HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials)))
	})

	t.Run("alter catalog integration: shared options", func(t *testing.T) {
		id := createObjectStorageCatalogIntegration(t).ID()

		err := client.CatalogIntegrations.Alter(ctx, sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*sdk.NewCatalogIntegrationSetRequest().
			WithComment(sdk.StringAllowEmpty{Value: "new comment"}).
			WithEnabled(false).
			WithRefreshIntervalSeconds(120)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.CatalogIntegrationObjectStorageDetails(t, id).
			HasEnabled(false).
			HasRefreshIntervalSeconds(120).
			HasComment("new comment"))
	})

	t.Run("alter catalog integration: bearer token", func(t *testing.T) {
		catalogUri, _, oAuthClientId, oAuthClientSecret := loadOpenCatalogConfig(t)
		newOAuthClientId := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogSecondaryOAuthClientId)
		newOAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogSecondaryOAuthClientSecret)
		bearerToken := testClientHelper().CatalogIntegration.RetrieveOpenCatalogBearerToken(t, catalogUri, oAuthClientId, oAuthClientSecret, oAuthAllowedScope)
		newBearerToken := testClientHelper().CatalogIntegration.RetrieveOpenCatalogBearerToken(t, catalogUri, newOAuthClientId, newOAuthClientSecret, additionalOAuthAllowedScope)

		integrationAwsGlue := createAwsGlueCatalogIntegration(t)
		integrationObjectStorage := createObjectStorageCatalogIntegration(t)
		integrationOpenCatalog := createOpenCatalogCatalogIntegration(t)
		integrationIcebergRest := createCatalogIntegrationWithRequest(t, sdk.NewCreateCatalogIntegrationRequest(testClientHelper().Ids.RandomAccountObjectIdentifier(), false).
			WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
				WithRestConfig(sdk.IcebergRestRestConfigRequest{
					CatalogUri:  catalogUri,
					CatalogName: new(catalogName),
				}).
				WithBearerRestAuthentication(sdk.BearerRestAuthenticationRequest{BearerToken: bearerToken})))

		request := *sdk.NewCatalogIntegrationSetRequest().
			WithSetBearerRestAuthentication(*sdk.NewSetBearerRestAuthenticationRequest(newBearerToken))

		err := client.CatalogIntegrations.Alter(ctx, sdk.NewAlterCatalogIntegrationRequest(integrationIcebergRest.ID()).WithSet(request))
		require.NoError(t, err)

		// Token is not returned by DESCRIBE, nothing to check

		invalid := []*sdk.CatalogIntegration{
			integrationAwsGlue, integrationObjectStorage, integrationOpenCatalog,
		}
		for _, integration := range invalid {
			id := integration.ID()
			err := client.CatalogIntegrations.Alter(ctx, sdk.NewAlterCatalogIntegrationRequest(id).WithSet(request))
			assert.ErrorContains(t, err, "Invalid option")
		}
	})

	t.Run("alter catalog integration: OAuth client secret", func(t *testing.T) {
		newOAuthClientSecret := testenvs.GetOrSkipTest(t, testenvs.OpenCatalogSecondaryOAuthClientSecret)

		integrationAwsGlue := createAwsGlueCatalogIntegration(t)
		integrationObjectStorage := createObjectStorageCatalogIntegration(t)
		integrationOpenCatalog := createOpenCatalogCatalogIntegration(t)
		integrationIcebergRest := createIcebergRestCatalogIntegration(t)

		request := *sdk.NewCatalogIntegrationSetRequest().
			WithSetOAuthRestAuthentication(*sdk.NewSetOAuthRestAuthenticationRequest(newOAuthClientSecret))

		valid := []*sdk.CatalogIntegration{
			integrationOpenCatalog, integrationIcebergRest,
		}
		for _, integration := range valid {
			id := integration.ID()
			err := client.CatalogIntegrations.Alter(ctx, sdk.NewAlterCatalogIntegrationRequest(id).WithSet(request))
			require.NoError(t, err)

			// Client secret is not returned by DESCRIBE, nothing to check
		}

		invalid := []*sdk.CatalogIntegration{
			integrationAwsGlue, integrationObjectStorage,
		}
		for _, integration := range invalid {
			id := integration.ID()
			err := client.CatalogIntegrations.Alter(ctx, sdk.NewAlterCatalogIntegrationRequest(id).WithSet(request))
			assert.ErrorContains(t, err, "Invalid option")
		}
	})

	t.Run("alter catalog integration: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier

		err := client.CatalogIntegrations.Alter(ctx, sdk.NewAlterCatalogIntegrationRequest(id).
			WithSet(*sdk.NewCatalogIntegrationSetRequest().WithEnabled(true)))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter catalog integration: non-existing with if exists option", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier

		err := client.CatalogIntegrations.Alter(ctx, sdk.NewAlterCatalogIntegrationRequest(id).
			WithSet(*sdk.NewCatalogIntegrationSetRequest().WithEnabled(true)).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("drop catalog integration: existing", func(t *testing.T) {
		id := createAwsGlueCatalogIntegration(t).ID()

		err := client.CatalogIntegrations.Drop(ctx, sdk.NewDropCatalogIntegrationRequest(id))
		require.NoError(t, err)

		_, err = client.CatalogIntegrations.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop catalog integration: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier

		err := client.CatalogIntegrations.Drop(ctx, sdk.NewDropCatalogIntegrationRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop catalog integration: non-existing with if exists option", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier

		err := client.CatalogIntegrations.Drop(ctx, sdk.NewDropCatalogIntegrationRequest(id).WithIfExists(true))
		require.NoError(t, err)
	})

	t.Run("show catalog integrations: default", func(t *testing.T) {
		integrationAwsGlue := createAwsGlueCatalogIntegration(t)
		integrationObjectStorage := createObjectStorageCatalogIntegration(t)
		integrationIcebergRest := createIcebergRestCatalogIntegration(t)

		showRequest := sdk.NewShowCatalogIntegrationRequest()
		returnedIntegrations, err := client.CatalogIntegrations.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *integrationAwsGlue)
		assert.Contains(t, returnedIntegrations, *integrationObjectStorage)
		assert.Contains(t, returnedIntegrations, *integrationIcebergRest)
	})

	t.Run("show catalog integrations: with like option", func(t *testing.T) {
		integrationAwsGlue := createAwsGlueCatalogIntegration(t)
		integrationObjectStorage := createObjectStorageCatalogIntegration(t)

		showRequest := sdk.NewShowCatalogIntegrationRequest().
			WithLike(sdk.Like{Pattern: &integrationAwsGlue.Name})
		returnedIntegrations, err := client.CatalogIntegrations.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Contains(t, returnedIntegrations, *integrationAwsGlue)
		assert.NotContains(t, returnedIntegrations, *integrationObjectStorage)
	})

	t.Run("describe catalog integration: AWS Glue", func(t *testing.T) {
		id := createAwsGlueCatalogIntegration(t).ID()

		assertThatObject(t, objectassert.CatalogIntegrationAwsGlueDetails(t, id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeGlue).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(glueRegion).
			HasCatalogNamespace(""))
	})

	t.Run("describe catalog integration: object storage", func(t *testing.T) {
		id := createObjectStorageCatalogIntegration(t).ID()

		assertThatObject(t, objectassert.CatalogIntegrationObjectStorageDetails(t, id).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeObjectStore).
			HasTableFormat(sdk.CatalogIntegrationTableFormatDelta).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment(""))
	})

	t.Run("describe catalog integration: Open Catalog", func(t *testing.T) {
		catalogUri, _, oAuthClientId, oAuthClientSecret := loadOpenCatalogConfig(t)
		id := createCatalogIntegrationWithRequest(t, createCatalogIntegrationOpenCatalogRequest(t, catalogUri, oAuthClientId, oAuthClientSecret)).ID()

		openCatalogDetails, err := client.CatalogIntegrations.DescribeOpenCatalogDetails(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.CatalogIntegrationOpenCatalogDetailsFromObject(t, openCatalogDetails).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypePolaris).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasCatalogNamespace("").
			HasRestConfigWith(objectassert.NewOpenCatalogRestConfigDetailsAssert().
				HasCatalogUri(catalogUri).
				HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
				HasCatalogName(catalogName).
				HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)).
			HasRestAuthenticationWith(objectassert.NewOAuthRestAuthenticationDetailsAssert().
				HasOauthTokenUri(catalogUri+"/v1/oauth/tokens").
				HasOauthClientId(oAuthClientId).
				HasOauthAllowedScopes(oAuthAllowedScope)))
	})

	t.Run("describe catalog integration: Iceberg REST", func(t *testing.T) {
		catalogUri, _, oAuthClientId, oAuthClientSecret := loadOpenCatalogConfig(t)
		id := createCatalogIntegrationWithRequest(t, createCatalogIntegrationIcebergRestRequest(t, catalogUri, oAuthClientId, oAuthClientSecret)).ID()

		icebergRestDetails, err := client.CatalogIntegrations.DescribeIcebergRestDetails(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.CatalogIntegrationIcebergRestDetailsFromObject(t, icebergRestDetails).
			HasCatalogSource(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest).
			HasTableFormat(sdk.CatalogIntegrationTableFormatIceberg).
			HasEnabled(false).
			HasRefreshIntervalSeconds(30).
			HasComment("").
			HasCatalogNamespace("").
			HasOAuthRestAuthenticationWith(objectassert.NewOAuthRestAuthenticationDetailsAssert().
				HasOauthTokenUri(catalogUri+"/v1/oauth/tokens").
				HasOauthClientId(oAuthClientId).
				HasOauthAllowedScopes(oAuthAllowedScope)).
			HasRestConfigWith(objectassert.NewIcebergRestRestConfigDetailsAssert().
				HasCatalogUri(catalogUri).
				HasPrefix("").
				HasCatalogName(catalogName).
				HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
				HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials)))
	})

	t.Run("describe catalog integration: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier

		_, err := client.CatalogIntegrations.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
