package invokeactionassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func ApiIntegrationDoesNotExist(t *testing.T, id sdk.AccountObjectIdentifier) assert.TestCheckFuncProvider {
	t.Helper()
	return newNonExistenceCheck(
		sdk.ObjectTypeApiIntegration,
		id,
		func(testClient *helpers.TestClient) showByIDFunc[*sdk.ApiIntegration, sdk.AccountObjectIdentifier] {
			return testClient.ApiIntegration.Show
		})
}
