package windows_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testfiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/stretchr/testify/require"
)

func TestReadFileSafeWorksOnWindows(t *testing.T) {
	if !oswrapper.IsRunningOnWindows() {
		t.Skip("reading files on other platforms is currently done in the sdk package")
	}
	exp := []byte("content")
	configPath := testfiles.TestFile(t, "config", exp)

	act, err := oswrapper.ReadFileSafe(configPath, true)
	require.NoError(t, err)
	require.Equal(t, exp, act)
}
