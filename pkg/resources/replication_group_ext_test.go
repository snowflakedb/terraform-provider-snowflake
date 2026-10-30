package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func Test_omitCurrentAccount(t *testing.T) {
	t.Parallel()

	current := sdk.NewAccountIdentifier("ORG", "SOURCE")
	target := sdk.NewAccountIdentifier("ORG", "TARGET")

	require.Equal(t, []sdk.AccountIdentifier{target}, omitCurrentAccount([]sdk.AccountIdentifier{current, target}, current))
	require.Empty(t, omitCurrentAccount([]sdk.AccountIdentifier{current}, current))
}
