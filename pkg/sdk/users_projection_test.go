package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/stretchr/testify/require"
)

func TestUserDetailsFromRows(t *testing.T) {
	details := userDetailsFromRows([]UserProperty{
		{Property: "NAME", Value: "u1"},
		{Property: "DISABLED", Value: "true"},
		{Property: "MINS_TO_UNLOCK", Value: "14"},
		{Property: "DAYS_TO_EXPIRY", Value: "1.5"},
		{Property: "MINS_TO_BYPASS_MFA", Value: "null"},
		{Property: "PASSWORD", Value: "secret"},
		{Property: "UNEXPECTED", Value: "x"},
	})
	require.Equal(t, "u1", details.Name)
	require.True(t, details.Disabled)
	require.Equal(t, 14, *details.MinsToUnlock)
	require.InDelta(t, 1.5, *details.DaysToExpiry, testvars.FloatEpsilon)
	require.Nil(t, details.MinsToBypassMfa)
	require.Equal(t, "secret", details.Password)
	require.Empty(t, details.LoginName)
}
