package invokeactionassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func BackupPolicyDoesNotExist(t *testing.T, id sdk.SchemaObjectIdentifier) assert.TestCheckFuncProvider {
	t.Helper()
	return newNonExistenceCheck(
		sdk.ObjectTypeBackupPolicy,
		id,
		func(testClient *helpers.TestClient) showByIDFunc[*sdk.BackupPolicy, sdk.SchemaObjectIdentifier] {
			return testClient.BackupPolicy.Show
		},
	)
}
