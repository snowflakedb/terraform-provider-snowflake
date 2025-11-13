package invokeactionassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func DatabaseRoleDoesNotExist(t *testing.T, id sdk.DatabaseObjectIdentifier) assert.TestCheckFuncProvider {
	t.Helper()
	return newNonExistenceCheck(
		sdk.ObjectTypeDatabaseRole,
		id,
		func(testClient *helpers.TestClient) showByIDFunc[*sdk.DatabaseRole, sdk.DatabaseObjectIdentifier] {
			return testClient.DatabaseRole.Show
		})
}
