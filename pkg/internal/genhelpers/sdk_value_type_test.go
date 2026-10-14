package genhelpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SdkMethodSuffix(t *testing.T) {
	require.Equal(t, "Bool", SdkMethodSuffix("bool"))
	require.Equal(t, "Sdk", SdkMethodSuffix("sdk.LogLevel"))
	require.Equal(t, "Sdk", SdkMethodSuffix("sdk.AccountObjectIdentifier"))
}

func Test_SdkValueMapper(t *testing.T) {
	require.Equal(t, "strconv.FormatBool(autoResume)", RunMapper(SdkValueMapper("bool"), "autoResume"))
	require.Equal(t, "forApplication.Name()", RunMapper(SdkValueMapper("identifier"), "forApplication"))
	require.Equal(t, "string(instanceFamily)", RunMapper(SdkValueMapper("enum"), "instanceFamily"))
	require.Equal(t, "strconv.FormatBool(expected)", RunMapper(SdkValueMapper("bool"), "expected"))
}
