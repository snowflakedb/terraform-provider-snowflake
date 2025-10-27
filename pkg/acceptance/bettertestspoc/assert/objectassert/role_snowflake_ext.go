package objectassert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func AccountRoleMissing(t *testing.T, id sdk.AccountObjectIdentifier) *RoleAssert {
	t.Helper()
	return &RoleAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeRole, id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.Role, sdk.AccountObjectIdentifier] {
			return func(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Role, error) {
				accountRole, err := testClient.Role.Show(t, id)
				if err != nil {
					if errors.Is(err, sdk.ErrObjectNotFound) {
						return accountRole, nil
					}
					return accountRole, nil
				}
				return nil, fmt.Errorf()
			}
		}),
	}
}
