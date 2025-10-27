package objectassert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func DatabaseIsMissing(t *testing.T, id sdk.AccountObjectIdentifier) *DatabaseAssert {
	t.Helper()
	return &DatabaseAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeDatabase, id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.Database, sdk.AccountObjectIdentifier] {
			return func(t *testing.T, identifier sdk.AccountObjectIdentifier) (*sdk.Database, error) {
				db, err := testClient.Database.Show(t, id)
				if err != nil {
					if errors.Is(err, sdk.ErrObjectNotFound) {
						return db, nil
					}
					return db, err
				}
				return db, fmt.Errorf("expected database %s to be missing, but it exists", id)
			}
		}),
	}
}
