package objectassert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func NetworkPolicyIsMissing(t *testing.T, id sdk.AccountObjectIdentifier) *NetworkPolicyAssert {
	t.Helper()
	return &NetworkPolicyAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeNetworkPolicy, id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.NetworkPolicy, sdk.AccountObjectIdentifier] {
			return func(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.NetworkPolicy, error) {
				t.Helper()
				networkPolicy, err := testClient.NetworkPolicy.Show(t, id)
				if err != nil {
					if errors.Is(err, sdk.ErrObjectNotFound) {
						return networkPolicy, nil
					}
					return networkPolicy, nil
				}
				return networkPolicy, fmt.Errorf("expected network policy %s to be missing, but it exists", id)
			}
		}),
	}
}
