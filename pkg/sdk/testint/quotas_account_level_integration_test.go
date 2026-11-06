//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_QuotasAccountLevel(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// an account can have only one Snowflake Intelligence object, so this can't run in parallel with other tests creating one
	t.Run("Snowflake Intelligence as a shared resource", func(t *testing.T) {
		intelligenceId, intelligenceCleanup := testClientHelper().SnowflakeIntelligence.Create(t)
		t.Cleanup(intelligenceCleanup)
		quotaId, quotaCleanup := testClientHelper().Quota.Create(t)
		t.Cleanup(quotaCleanup)

		testClientHelper().Quota.AddSharedResource(t, quotaId, sdk.QuotaSharedResourceDomainSnowflakeIntelligence, intelligenceId.FullyQualifiedName())

		details, err := client.Quotas.DescribeDetails(ctx, quotaId)
		require.NoError(t, err)
		require.Len(t, details.Scope.SharedResources, 1)
		intelligence := details.Scope.SharedResources[0]
		domain, err := intelligence.NormalizedDomain()
		require.NoError(t, err)
		assert.Equal(t, sdk.QuotaSharedResourceDomainSnowflakeIntelligence, domain)
		assert.False(t, intelligence.IsWholeDomain())
		assert.Equal(t, intelligenceId.Name(), intelligence.Name)

		args, err := sdk.NewQuotaRemoveSharedResourceArgsRequestFromInstance(sdk.QuotaSharedResourceDomainSnowflakeIntelligence, intelligenceId.FullyQualifiedName())
		require.NoError(t, err)
		_, err = client.Quotas.RemoveSharedResource(ctx, sdk.NewRemoveSharedResourceQuotaRequest(quotaId, *args))
		require.NoError(t, err)

		details, err = client.Quotas.DescribeDetails(ctx, quotaId)
		require.NoError(t, err)
		assert.Empty(t, details.Scope.SharedResources)
	})
}
