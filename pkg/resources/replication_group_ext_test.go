package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_omitCurrentAccount(t *testing.T) {
	t.Parallel()

	current := sdk.NewAccountIdentifier("ORG", "SOURCE")
	target := sdk.NewAccountIdentifier("ORG", "TARGET")

	require.Equal(t, []sdk.AccountIdentifier{target}, omitCurrentAccount([]sdk.AccountIdentifier{current, target}, current))
	require.Empty(t, omitCurrentAccount([]sdk.AccountIdentifier{current}, current))
}

func Test_replicationGroupApplySetUnset_sendsEachStatementInItsOwnRequest(t *testing.T) {
	d := schema.TestResourceDataRaw(t, replicationGroupSchema, map[string]any{
		"object_types":         []any{"DATABASES"},
		"allowed_databases":    []any{"db_one"},
		"comment":              "c",
		"replication_schedule": "10 MINUTE",
		"error_integration":    "notifications",
	})
	requests, err := replicationGroupApplySetUnset(d)
	require.NoError(t, err)

	// OBJECT_TYPES and the allowed lists cannot share a SET with the properties.
	require.Equal(t, sdk.NewReplicationGroupSetRequest().
		WithObjectTypes([]sdk.ReplicationGroupObjectType{sdk.ReplicationGroupObjectTypeDatabases}).
		WithAllowedDatabases([]sdk.AccountObjectIdentifier{sdk.NewAccountObjectIdentifier("db_one")}), requests.setObjectLists)
	require.Equal(t, sdk.NewReplicationGroupSetRequest().
		WithComment("c").
		WithReplicationSchedule("10 MINUTE").
		WithErrorIntegration(sdk.NewAccountObjectIdentifier("notifications")), requests.set)
	require.Equal(t, sdk.NewReplicationGroupUnsetRequest(), requests.unset)
}
