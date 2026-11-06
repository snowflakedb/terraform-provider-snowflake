package objectassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// PostgresMirrorDetails describes the mirror when the assertion runs. A mirror has no identifier, so the account object
// identifier only carries the mirror name to DESCRIBE_MIRROR.
func PostgresMirrorDetails(t *testing.T, name string) *PostgresMirrorDetailsAssert {
	t.Helper()
	return &PostgresMirrorDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("PostgresMirrorDetails"), sdk.NewAccountObjectIdentifier(name), func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.PostgresMirrorDetails, sdk.AccountObjectIdentifier] {
			return func(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.PostgresMirrorDetails, error) {
				t.Helper()
				return testClient.PostgresMirror.Describe(t, id.Name())
			}
		}),
	}
}
