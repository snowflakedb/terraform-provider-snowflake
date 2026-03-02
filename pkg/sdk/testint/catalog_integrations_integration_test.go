//go:build non_account_level_tests

package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_CatalogIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	const (
		glueAwsRoleArn       = "arn:aws:iam::123456789012:role/sqsAccess"
		glueCatalogId        = "123456789012"
		glueRegion           = "us-east-2"
		polarisCatalogUri    = "https://testorg-testacc.snowflakecomputing.com/polaris/api/catalog"
		restCatalogUri       = "https://api.tabular.io/ws"
		sapBdcInvitationLink = "https://example.hanacloud.ondemand.com/?code=123e4567-e89b-12d3-a456-426614174000"
		sapBdcCatalogUri     = "https://example.hanacloud.ondemand.com"
		oAuthClientId        = "my_client_id"
		oAuthClientSecret    = "my_client_secret"
		oAuthAllowedScope    = "PRINCIPAL_ROLE:ALL"
	)

	assertCatalogIntegration := func(t *testing.T, s *sdk.CatalogIntegration, name sdk.AccountObjectIdentifier, comment string) {
		t.Helper()
		assertThatObject(t, objectassert.CatalogIntegrationFromObject(t, s).
			HasName(name.Name()).
			HasEnabled(false).
			HasType("CATALOG").
			HasCategory("CATALOG").
			HasComment(comment))
	}

	assertSharedProperties := func(t *testing.T, details []sdk.CatalogIntegrationProperty, catalogSource sdk.CatalogIntegrationCatalogSourceType, tableFormat sdk.CatalogIntegrationTableFormat, enabled bool, refreshIntervalSeconds int, comment string) {
		t.Helper()
		assert.Contains(t, details, sdk.CatalogIntegrationProperty{Name: "CATALOG_SOURCE", Type: "String", Value: string(catalogSource), Default: ""})
		assert.Contains(t, details, sdk.CatalogIntegrationProperty{Name: "TABLE_FORMAT", Type: "String", Value: string(tableFormat), Default: ""})
		assert.Contains(t, details, sdk.CatalogIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: fmt.Sprintf("%t", enabled), Default: "false"})
		assert.Contains(t, details, sdk.CatalogIntegrationProperty{Name: "REFRESH_INTERVAL_SECONDS", Type: "Integer", Value: fmt.Sprintf("%d", refreshIntervalSeconds), Default: "30"})
		assert.Contains(t, details, sdk.CatalogIntegrationProperty{Name: "COMMENT", Type: "String", Value: comment, Default: ""})
	}

	cleanupCatalogIntegrationProvider := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.CatalogIntegrations.Drop(ctx, sdk.NewDropCatalogIntegrationRequest(id).WithIfExists(true))
			require.NoError(t, err)
		}
	}

	createCatalogIntegrationWithRequest := func(t *testing.T, request *sdk.CreateCatalogIntegrationRequest) *sdk.CatalogIntegration {
		t.Helper()
		id := request.GetName()

		err := client.CatalogIntegrations.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupCatalogIntegrationProvider(id))

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

	createCatalogIntegrationOpenCatalogRequest := func(t *testing.T) *sdk.CreateCatalogIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateCatalogIntegrationRequest(id, false).
			WithOpenCatalogCatalogSourceParams(*sdk.NewOpenCatalogParamsRequest().
				WithRestConfig(sdk.OpenCatalogRestConfigRequest{
					CatalogUri:  polarisCatalogUri,
					CatalogName: "my_catalog_name",
				}).
				WithRestAuthentication(sdk.OAuthRestAuthenticationRequest{
					OauthClientId:      oAuthClientId,
					OauthClientSecret:  oAuthClientSecret,
					OauthAllowedScopes: []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}},
				}))
	}

	createCatalogIntegrationIcebergRestRequest := func(t *testing.T) *sdk.CreateCatalogIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateCatalogIntegrationRequest(id, false).
			WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
				WithRestConfig(sdk.IcebergRestRestConfigRequest{CatalogUri: restCatalogUri}).
				WithBearerRestAuthentication(sdk.BearerRestAuthenticationRequest{BearerToken: "test-token"}))
	}

	createCatalogIntegrationSapBdcRequest := func(t *testing.T) *sdk.CreateCatalogIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		return sdk.NewCreateCatalogIntegrationRequest(id, false).
			WithSapBdcCatalogSourceParams(*sdk.NewSapBdcParamsRequest().
				WithRestConfig(sdk.SapBdcRestConfigRequest{SapBdcInvitationLink: sapBdcInvitationLink}))
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
		return createCatalogIntegrationWithRequest(t, createCatalogIntegrationOpenCatalogRequest(t))
	}

	createIcebergRestCatalogIntegration := func(t *testing.T) *sdk.CatalogIntegration {
		t.Helper()
		return createCatalogIntegrationWithRequest(t, createCatalogIntegrationIcebergRestRequest(t))
	}

	createSapBdcCatalogIntegration := func(t *testing.T) *sdk.CatalogIntegration {
		t.Helper()
		return createCatalogIntegrationWithRequest(t, createCatalogIntegrationSapBdcRequest(t))
	}

	t.Run("create catalog integration: AWS Glue basic", func(t *testing.T) {
		request := createCatalogIntegrationAwsGlueRequest(t)

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "")
		assertThatObject(t, objectassert.AwsGlueParams(t, integration.ID()).
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(glueRegion).
			HasCatalogNamespace(""))
	})

	t.Run("create catalog integration: object storage basic", func(t *testing.T) {
		request := createCatalogIntegrationObjectStorageRequest(t)

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "")
		assertThatObject(t, objectassert.ObjectStorageParams(t, integration.ID()).HasTableFormat(sdk.CatalogIntegrationTableFormatDelta))
	})

	t.Run("create catalog integration: Open Catalog basic", func(t *testing.T) {
		request := createCatalogIntegrationOpenCatalogRequest(t)

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "")

		openCatalogParams, err := client.CatalogIntegrations.DescribeOpenCatalogParams(ctx, integration.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.OpenCatalogParamsFromObject(t, openCatalogParams).
			HasCatalogNamespace(""))
		assertThatObject(t, objectassert.OpenCatalogRestConfigFromObject(t, &openCatalogParams.RestConfig).
			HasCatalogUri(polarisCatalogUri).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasCatalogName("my_catalog_name").
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials))
		assertThatObject(t, objectassert.OAuthRestAuthenticationFromObject(t, &(openCatalogParams.RestAuthentication)).
			HasOauthTokenUri(polarisCatalogUri+"/v1/oauth/tokens").
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(sdk.StringListItemWrapper{oAuthAllowedScope}))
	})

	t.Run("create catalog integration: Iceberg REST basic", func(t *testing.T) {
		request := createCatalogIntegrationIcebergRestRequest(t)

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "")

		icebergRestParams, err := client.CatalogIntegrations.DescribeIcebergRestParams(ctx, integration.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.IcebergRestParamsFromObject(t, icebergRestParams).
			HasCatalogNamespace("").
			HasBearerRestAuthentication())
		assertThatObject(t, objectassert.IcebergRestRestConfigFromObject(t, &icebergRestParams.RestConfig).
			HasCatalogUri(restCatalogUri).
			HasPrefix("").
			HasCatalogName("").
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials))
	})

	t.Run("create catalog integration: SAP Business Data Cloud basic", func(t *testing.T) {
		request := createCatalogIntegrationSapBdcRequest(t)

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "")
	})

	t.Run("create catalog integration: AWS Glue all options", func(t *testing.T) {
		const catalogNamespace = "myNamespace"
		request := createCatalogIntegrationAwsGlueRequest(t).
			WithIfNotExists(true).
			WithAwsGlueCatalogSourceParams(*sdk.NewAwsGlueParamsRequest(glueAwsRoleArn, glueCatalogId).
				WithGlueRegion(glueRegion).
				WithCatalogNamespace(catalogNamespace)).
			WithRefreshIntervalSeconds(120).
			WithComment("test comment")

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "test comment")
		assertThatObject(t, objectassert.AwsGlueParams(t, integration.ID()).
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(glueRegion).
			HasCatalogNamespace(catalogNamespace))
	})

	t.Run("create catalog integration: object storage all options", func(t *testing.T) {
		request := createCatalogIntegrationObjectStorageRequest(t).
			WithIfNotExists(true).
			WithRefreshIntervalSeconds(120).
			WithComment("test comment")

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "test comment")
		assertThatObject(t, objectassert.ObjectStorageParams(t, integration.ID()).HasTableFormat(sdk.CatalogIntegrationTableFormatDelta))
	})

	t.Run("create catalog integration: Open Catalog all options", func(t *testing.T) {
		const catalogNamespace = "myNamespace"
		const catalogName = "my_catalog_name"
		const polarisCatalogUri = "https://testorg-testacc.privatelink.snowflakecomputing.com/polaris/api/catalog"
		const oAuthTokenUri = polarisCatalogUri + "/v2/oauth/tokens"
		request := createCatalogIntegrationOpenCatalogRequest(t).
			WithIfNotExists(true).
			WithOpenCatalogCatalogSourceParams(*sdk.NewOpenCatalogParamsRequest().
				WithCatalogNamespace(catalogNamespace).
				WithRestConfig(sdk.OpenCatalogRestConfigRequest{
					CatalogUri:           polarisCatalogUri,
					CatalogApiType:       sdk.Pointer(sdk.CatalogIntegrationCatalogApiTypePrivate),
					CatalogName:          catalogName,
					AccessDelegationMode: sdk.Pointer(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials),
				}).
				WithRestAuthentication(sdk.OAuthRestAuthenticationRequest{
					OauthTokenUri:      sdk.String(oAuthTokenUri),
					OauthClientId:      oAuthClientId,
					OauthClientSecret:  oAuthClientSecret,
					OauthAllowedScopes: []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}, {Value: "DUMMY"}},
				})).
			WithRefreshIntervalSeconds(120).
			WithComment("test comment")

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "test comment")

		openCatalogParams, err := client.CatalogIntegrations.DescribeOpenCatalogParams(ctx, integration.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.OpenCatalogParamsFromObject(t, openCatalogParams).
			HasCatalogNamespace(catalogNamespace))
		assertThatObject(t, objectassert.OpenCatalogRestConfigFromObject(t, &openCatalogParams.RestConfig).
			HasCatalogUri(polarisCatalogUri).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePrivate).
			HasCatalogName(catalogName).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials))
		assertThatObject(t, objectassert.OAuthRestAuthenticationFromObject(t, &(openCatalogParams.RestAuthentication)).
			HasOauthTokenUri(oAuthTokenUri).
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(sdk.StringListItemWrapper{oAuthAllowedScope}, sdk.StringListItemWrapper{"DUMMY"}))
	})

	t.Run("create catalog integration: Iceberg REST all options", func(t *testing.T) {
		const catalogNamespace = "myNamespace"
		const prefix = "prefix"
		const catalogName = "my_catalog_name"
		const sigV4IamRole = "arn:aws:iam::123456789012:role/my-role"
		const sigV4SigningRole = "us-west-2"
		request := createCatalogIntegrationIcebergRestRequest(t).
			WithIfNotExists(true).
			WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
				WithCatalogNamespace(catalogNamespace).
				WithRestConfig(sdk.IcebergRestRestConfigRequest{
					CatalogUri:           restCatalogUri,
					Prefix:               sdk.String(prefix),
					CatalogName:          sdk.String(catalogName),
					CatalogApiType:       sdk.Pointer(sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway),
					AccessDelegationMode: sdk.Pointer(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials),
				}).
				WithSigV4RestAuthentication(sdk.SigV4RestAuthenticationRequest{
					Sigv4IamRole:       sigV4IamRole,
					Sigv4SigningRegion: sdk.String(sigV4SigningRole),
					Sigv4ExternalId:    sdk.String("external_id"),
				})).
			WithRefreshIntervalSeconds(120).
			WithComment("test comment")

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "test comment")

		icebergRestParams, err := client.CatalogIntegrations.DescribeIcebergRestParams(ctx, integration.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.IcebergRestParamsFromObject(t, icebergRestParams).
			HasCatalogNamespace(catalogNamespace).
			HasSigV4RestAuthentication())
		assertThatObject(t, objectassert.IcebergRestRestConfigFromObject(t, &icebergRestParams.RestConfig).
			HasCatalogUri(restCatalogUri).
			HasPrefix(prefix).
			HasCatalogName(catalogName).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypeAwsApiGateway).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeVendedCredentials))
		assertThatObject(t, objectassert.SigV4RestAuthenticationFromObject(t, icebergRestParams.SigV4RestAuthentication).
			HasSigv4IamRole(sigV4IamRole).
			HasSigv4SigningRegion(sigV4SigningRole))
	})

	t.Run("create catalog integration: SAP Business Data Cloud all options", func(t *testing.T) {
		request := createCatalogIntegrationSapBdcRequest(t).
			WithIfNotExists(true).
			WithRefreshIntervalSeconds(120).
			WithComment("test comment")

		integration := createCatalogIntegrationWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), "test comment")
	})

	t.Run("alter catalog integration: shared options", func(t *testing.T) {
		id := createObjectStorageCatalogIntegration(t).ID()

		err := client.CatalogIntegrations.Alter(ctx, sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*sdk.NewCatalogIntegrationSetRequest().
			WithComment(sdk.StringAllowEmpty{Value: "new comment"}).
			WithEnabled(true).
			WithRefreshIntervalSeconds(120)))
		require.NoError(t, err)

		details, err := client.CatalogIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSharedProperties(t, details, sdk.CatalogIntegrationCatalogSourceTypeObjectStorage, sdk.CatalogIntegrationTableFormatDelta, true, 120, "new comment")
	})

	t.Run("alter catalog integration: bearer token", func(t *testing.T) {
		integrationAwsGlue := createAwsGlueCatalogIntegration(t)
		integrationObjectStorage := createObjectStorageCatalogIntegration(t)
		integrationOpenCatalog := createOpenCatalogCatalogIntegration(t)
		integrationIcebergRest := createIcebergRestCatalogIntegration(t)
		integrationSapBdc := createSapBdcCatalogIntegration(t)

		request := *sdk.NewCatalogIntegrationSetRequest().
			WithSetBearerRestAuthentication(*sdk.NewSetBearerRestAuthenticationRequest("new token"))

		err := client.CatalogIntegrations.Alter(ctx, sdk.NewAlterCatalogIntegrationRequest(integrationIcebergRest.ID()).WithSet(request))
		require.NoError(t, err)

		// Token is not returned by DESCRIBE, nothing to check

		invalid := []*sdk.CatalogIntegration{
			integrationAwsGlue, integrationObjectStorage, integrationOpenCatalog, integrationSapBdc,
		}
		for _, integration := range invalid {
			id := integration.ID()
			err := client.CatalogIntegrations.Alter(ctx, sdk.NewAlterCatalogIntegrationRequest(id).WithSet(request))
			assert.ErrorContains(t, err, "Invalid option")
		}
	})

	t.Run("alter catalog integration: OAuth client secret", func(t *testing.T) {
		integrationAwsGlue := createAwsGlueCatalogIntegration(t)
		integrationObjectStorage := createObjectStorageCatalogIntegration(t)
		integrationOpenCatalog := createOpenCatalogCatalogIntegration(t)

		createRequest := createCatalogIntegrationIcebergRestRequest(t).
			WithIcebergRestCatalogSourceParams(*sdk.NewIcebergRestParamsRequest().
				WithRestConfig(sdk.IcebergRestRestConfigRequest{CatalogUri: restCatalogUri}).
				WithOAuthRestAuthentication(sdk.OAuthRestAuthenticationRequest{
					OauthClientId:      oAuthClientId,
					OauthClientSecret:  oAuthClientSecret,
					OauthAllowedScopes: []sdk.StringListItemWrapper{{Value: oAuthAllowedScope}},
				}))
		integrationIcebergRest := createCatalogIntegrationWithRequest(t, createRequest)

		request := *sdk.NewCatalogIntegrationSetRequest().
			WithSetOAuthRestAuthentication(*sdk.NewSetOAuthRestAuthenticationRequest("new secret"))

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

		details, err := client.CatalogIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSharedProperties(t, details, sdk.CatalogIntegrationCatalogSourceTypeAWSGlue, sdk.CatalogIntegrationTableFormatIceberg, false, 30, "")
		assertThatObject(t, objectassert.AwsGlueParams(t, id).
			HasGlueAwsRoleArn(glueAwsRoleArn).
			HasGlueCatalogId(glueCatalogId).
			HasGlueRegion(glueRegion).
			HasCatalogNamespace(""))
	})

	t.Run("describe catalog integration: object storage", func(t *testing.T) {
		id := createObjectStorageCatalogIntegration(t).ID()

		details, err := client.CatalogIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSharedProperties(t, details, sdk.CatalogIntegrationCatalogSourceTypeObjectStorage, sdk.CatalogIntegrationTableFormatDelta, false, 30, "")
	})

	t.Run("describe catalog integration: Open Catalog", func(t *testing.T) {
		id := createOpenCatalogCatalogIntegration(t).ID()

		details, err := client.CatalogIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSharedProperties(t, details, sdk.CatalogIntegrationCatalogSourceTypePolaris, sdk.CatalogIntegrationTableFormatIceberg, false, 30, "")

		openCatalogParams, err := client.CatalogIntegrations.DescribeOpenCatalogParams(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.OpenCatalogParamsFromObject(t, openCatalogParams).
			HasCatalogNamespace(""))
		assertThatObject(t, objectassert.OpenCatalogRestConfigFromObject(t, &openCatalogParams.RestConfig).
			HasCatalogUri(polarisCatalogUri).
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasCatalogName("my_catalog_name").
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials))
		assertThatObject(t, objectassert.OAuthRestAuthenticationFromObject(t, &(openCatalogParams.RestAuthentication)).
			HasOauthTokenUri(polarisCatalogUri+"/v1/oauth/tokens").
			HasOauthClientId(oAuthClientId).
			HasOauthAllowedScopes(sdk.StringListItemWrapper{oAuthAllowedScope}))
	})

	t.Run("describe catalog integration: Iceberg REST", func(t *testing.T) {
		id := createIcebergRestCatalogIntegration(t).ID()

		details, err := client.CatalogIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSharedProperties(t, details, sdk.CatalogIntegrationCatalogSourceTypeIcebergREST, sdk.CatalogIntegrationTableFormatIceberg, false, 30, "")

		icebergRestParams, err := client.CatalogIntegrations.DescribeIcebergRestParams(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.IcebergRestParamsFromObject(t, icebergRestParams).
			HasCatalogNamespace("").
			HasBearerRestAuthentication())
		assertThatObject(t, objectassert.IcebergRestRestConfigFromObject(t, &icebergRestParams.RestConfig).
			HasCatalogUri(restCatalogUri).
			HasPrefix("").
			HasCatalogName("").
			HasCatalogApiType(sdk.CatalogIntegrationCatalogApiTypePublic).
			HasAccessDelegationMode(sdk.CatalogIntegrationAccessDelegationModeExternalVolumeCredentials))
	})

	t.Run("describe catalog integration: SAP Business Data Cloud", func(t *testing.T) {
		id := createSapBdcCatalogIntegration(t).ID()

		details, err := client.CatalogIntegrations.Describe(ctx, id)
		require.NoError(t, err)

		assertSharedProperties(t, details, sdk.CatalogIntegrationCatalogSourceTypeSAPBusinessDataCloud, sdk.CatalogIntegrationTableFormatDelta, false, 30, "")
	})

	t.Run("describe catalog integration: non-existing", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier

		_, err := client.CatalogIntegrations.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
