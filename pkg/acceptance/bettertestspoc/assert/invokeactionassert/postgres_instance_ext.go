package invokeactionassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func PostgresInstanceDoesNotExist(t *testing.T, id sdk.AccountObjectIdentifier) assert.TestCheckFuncProvider {
	t.Helper()
	return newNonExistenceCheck(
		sdk.ObjectTypePostgresInstance,
		id,
		func(testClient *helpers.TestClient) showByIDFunc[*sdk.PostgresInstance, sdk.AccountObjectIdentifier] {
			return testClient.PostgresInstance.Show
		},
	)
}
