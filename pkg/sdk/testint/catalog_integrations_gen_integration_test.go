//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_CatalogIntegrations(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertCatalogIntegration := func(t *testing.T, s *sdk.CatalogIntegration, name sdk.AccountObjectIdentifier, enabled bool) {
		t.Helper()
		assert.Equal(t, name.Name(), s.Name)
		assert.Equal(t, enabled, s.Enabled)
		assert.Equal(t, "CATALOG", s.Category)
	}

	cleanupCatalogIntegration := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.CatalogIntegrations.Drop(ctx, sdk.NewDropCatalogIntegrationRequest(id).WithIfExists(true))
			require.NoError(t, err)
		}
	}

	createObjectStoreRequest := func(t *testing.T) *sdk.CreateCatalogIntegrationRequest {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		return sdk.NewCreateCatalogIntegrationRequest(id, true).
			WithObjectStoreParams(*sdk.NewObjectStoreParamsRequest(sdk.TableFormatIceberg))
	}

	createWithRequest := func(t *testing.T, request *sdk.CreateCatalogIntegrationRequest) *sdk.CatalogIntegration {
		t.Helper()
		id := request.GetName()

		err := client.CatalogIntegrations.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupCatalogIntegration(id))

		integration, err := client.CatalogIntegrations.ShowByID(ctx, id)
		require.NoError(t, err)

		return integration
	}

	t.Run("create and describe - object store iceberg", func(t *testing.T) {
		request := createObjectStoreRequest(t)

		integration := createWithRequest(t, request)

		assertCatalogIntegration(t, integration, request.GetName(), true)

		details, err := client.CatalogIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.CatalogIntegrationProperty{Name: "ENABLED", Type: "Boolean", Value: "true", Default: "false"})
		assert.Contains(t, details, sdk.CatalogIntegrationProperty{Name: "CATALOG_SOURCE", Type: "String", Value: "OBJECT_STORE", Default: ""})
		assert.Contains(t, details, sdk.CatalogIntegrationProperty{Name: "TABLE_FORMAT", Type: "String", Value: "ICEBERG", Default: ""})
	})

	t.Run("create - object store delta", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		request := sdk.NewCreateCatalogIntegrationRequest(id, true).
			WithObjectStoreParams(*sdk.NewObjectStoreParamsRequest(sdk.TableFormatDelta))

		integration := createWithRequest(t, request)

		assertCatalogIntegration(t, integration, id, true)

		details, err := client.CatalogIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)

		assert.Contains(t, details, sdk.CatalogIntegrationProperty{Name: "TABLE_FORMAT", Type: "String", Value: "DELTA", Default: ""})
	})

	t.Run("create with comment", func(t *testing.T) {
		request := createObjectStoreRequest(t).
			WithComment("test catalog integration")

		integration := createWithRequest(t, request)

		assert.Equal(t, "test catalog integration", integration.Comment)
	})

	// ALTER on catalog integrations currently triggers a Snowflake internal error (incident 000603/XX000).
	// Skipping until the server-side issue is resolved.
	t.Run("alter - set and unset comment", func(t *testing.T) {
		t.Skipf("Skipping due to Snowflake internal error on ALTER CATALOG INTEGRATION")
	})

	t.Run("drop", func(t *testing.T) {
		request := createObjectStoreRequest(t)
		id := request.GetName()

		err := client.CatalogIntegrations.Create(ctx, request)
		require.NoError(t, err)

		err = client.CatalogIntegrations.Drop(ctx, sdk.NewDropCatalogIntegrationRequest(id))
		require.NoError(t, err)

		_, err = client.CatalogIntegrations.ShowByID(ctx, id)
		require.Error(t, err)
	})

	t.Run("drop safely - existing", func(t *testing.T) {
		request := createObjectStoreRequest(t)
		id := request.GetName()

		err := client.CatalogIntegrations.Create(ctx, request)
		require.NoError(t, err)

		err = client.CatalogIntegrations.DropSafely(ctx, id)
		require.NoError(t, err)
	})

	t.Run("drop safely - non existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.CatalogIntegrations.DropSafely(ctx, id)
		require.NoError(t, err)
	})

	t.Run("show", func(t *testing.T) {
		request := createObjectStoreRequest(t)
		integration := createWithRequest(t, request)

		integrations, err := client.CatalogIntegrations.Show(ctx, sdk.NewShowCatalogIntegrationRequest())
		require.NoError(t, err)
		assert.NotEmpty(t, integrations)

		_, err = client.CatalogIntegrations.ShowByID(ctx, integration.ID())
		require.NoError(t, err)
	})

	t.Run("show with like filter", func(t *testing.T) {
		request := createObjectStoreRequest(t)
		integration := createWithRequest(t, request)

		integrations, err := client.CatalogIntegrations.Show(ctx, sdk.NewShowCatalogIntegrationRequest().
			WithLike(sdk.Like{Pattern: sdk.String(integration.Name)}))
		require.NoError(t, err)
		assert.Len(t, integrations, 1)
		assert.Equal(t, integration.Name, integrations[0].Name)
	})

	t.Run("describe", func(t *testing.T) {
		request := createObjectStoreRequest(t)
		integration := createWithRequest(t, request)

		details, err := client.CatalogIntegrations.Describe(ctx, integration.ID())
		require.NoError(t, err)
		assert.NotEmpty(t, details)
	})
}
