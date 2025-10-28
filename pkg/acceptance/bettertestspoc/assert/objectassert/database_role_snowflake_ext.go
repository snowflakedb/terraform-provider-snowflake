package objectassert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func DatabaseRoleIsMissing(t *testing.T, id sdk.DatabaseObjectIdentifier) *DatabaseRoleAssert {
	t.Helper()
	return &DatabaseRoleAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeDatabaseRole, id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.DatabaseRole, sdk.DatabaseObjectIdentifier] {
			return func(t *testing.T, id sdk.DatabaseObjectIdentifier) (*sdk.DatabaseRole, error) {
				t.Helper()
				databaseRole, err := testClient.DatabaseRole.Show(t, id)
				if err != nil {
					if errors.Is(err, sdk.ErrObjectNotFound) {
						return databaseRole, nil
					}
					return databaseRole, nil
				}
				return databaseRole, fmt.Errorf("expected database role %s to be missing, but it exists", id)
			}
		}),
	}
}
