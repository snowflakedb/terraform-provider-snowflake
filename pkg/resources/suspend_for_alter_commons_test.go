package resources

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_suspendForAlterIf(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{"comment": {Type: schema.TypeString, Optional: true}}, map[string]any{"comment": "c"})
	id := sdk.NewAccountObjectIdentifier("widget")
	calls := 0
	suspendIfStarted := func(context.Context, sdk.AccountObjectIdentifier) (func(context.Context) error, error) {
		calls++
		return func(context.Context) error { return nil }, nil
	}

	// The resource says the change keeps Snowflake as it is, so nothing is suspended although comment changed.
	suspension, err := SuspendForAlterIf(context.Background(), d, false, suspendIfStarted, id, "comment")
	require.NoError(t, err)
	require.Nil(t, suspension.resume)
	require.Zero(t, calls)

	suspension, err = SuspendForAlterIf(context.Background(), d, true, suspendIfStarted, id, "comment")
	require.NoError(t, err)
	require.NotNil(t, suspension.resume)
	require.Equal(t, 1, calls)
}
