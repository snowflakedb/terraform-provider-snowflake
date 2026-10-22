package helpers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type CatalogIntegrationClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewCatalogIntegrationClient(context *TestClientContext, idsGenerator *IdsGenerator) *CatalogIntegrationClient {
	return &CatalogIntegrationClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *CatalogIntegrationClient) client() sdk.CatalogIntegrations {
	return c.context.client.CatalogIntegrations
}

func (c *CatalogIntegrationClient) Create(t *testing.T) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()
	id := c.ids.RandomAccountObjectIdentifier()

	err := c.client().Create(ctx, sdk.NewCreateCatalogIntegrationRequest(id, true).
		WithObjectStorageCatalogSourceParams(*sdk.NewObjectStorageParamsRequest(sdk.CatalogIntegrationTableFormatIceberg)))
	require.NoError(t, err)

	return id, c.DropFunc(t, id)
}

func (c *CatalogIntegrationClient) CreateFunc(t *testing.T, request *sdk.CreateCatalogIntegrationRequest) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()
	id := request.GetName()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	return id, c.DropFunc(t, id)
}

func (c *CatalogIntegrationClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropCatalogIntegrationRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *CatalogIntegrationClient) Alter(t *testing.T, request *sdk.AlterCatalogIntegrationRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, request)
	require.NoError(t, err)
}

func (c *CatalogIntegrationClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegration, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *CatalogIntegrationClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) ([]sdk.CatalogIntegrationProperty, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().Describe(ctx, id)
}

func (c *CatalogIntegrationClient) DescribeAwsGlue(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationAwsGlueDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeAwsGlueDetails(ctx, id)
}

func (c *CatalogIntegrationClient) DescribeObjectStorage(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationObjectStorageDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeObjectStorageDetails(ctx, id)
}

func (c *CatalogIntegrationClient) DescribeOpenCatalog(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationOpenCatalogDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeOpenCatalogDetails(ctx, id)
}

func (c *CatalogIntegrationClient) DescribeIcebergRest(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationIcebergRestDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeIcebergRestDetails(ctx, id)
}

func (c *CatalogIntegrationClient) DescribeDetails(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationAllDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeDetails(ctx, id)
}

// RetrieveOpenCatalogBearerToken fetches an OAuth access token from Open Catalog using the
// client-credentials flow, which can then be used as a bearer token for Iceberg REST catalog integrations.
// See: https://docs.snowflake.com/en/user-guide/tables-iceberg-configure-catalog-integration-rest-check-config#step-1-retrieve-an-access-token
func (c *CatalogIntegrationClient) RetrieveOpenCatalogBearerToken(t *testing.T, catalogUri, clientId, clientSecret, scope string) string {
	t.Helper()

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("scope", scope)
	form.Set("client_id", clientId)
	form.Set("client_secret", clientSecret)

	req, err := http.NewRequest(http.MethodPost, catalogUri+"/v1/oauth/tokens", strings.NewReader(form.Encode()))
	require.NoError(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "token request failed: %s", body)

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	require.NoError(t, json.Unmarshal(body, &tokenResponse))
	require.NotEmpty(t, tokenResponse.AccessToken)
	return tokenResponse.AccessToken
}
