package objectassert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func SchemaIsMissing(t *testing.T, id sdk.DatabaseObjectIdentifier) *SchemaAssert {
	t.Helper()
	return &SchemaAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeSchema, id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.Schema, sdk.DatabaseObjectIdentifier] {
			return func(t *testing.T, id sdk.DatabaseObjectIdentifier) (*sdk.Schema, error) {
				t.Helper()
				schema, err := testClient.Schema.Show(t, id)
				if err != nil {
					if errors.Is(err, sdk.ErrObjectNotFound) {
						return schema, nil
					}
					return schema, nil
				}
				return schema, fmt.Errorf("expected schema %s to be missing, but it exists", id)
			}
		}),
	}
}
