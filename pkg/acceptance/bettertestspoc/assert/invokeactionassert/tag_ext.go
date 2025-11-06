package invokeactionassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func TagDoesNotExist(t *testing.T, id sdk.SchemaObjectIdentifier) assert.TestCheckFuncProvider {
	t.Helper()
	return newNonExistenceCheck(
		sdk.ObjectTypeTag,
		id,
		func(testClient *helpers.TestClient) showByIDFunc[*sdk.Tag, sdk.SchemaObjectIdentifier] {
			return testClient.Tag.Show
		})
}
