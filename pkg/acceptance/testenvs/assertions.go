package testenvs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func AssertEnvNotSet(t *testing.T, envName string) {
	t.Helper()
	require.Emptyf(t, os.Getenv(envName), "environment variable %v should not be set", envName)
}

func AssertEnvSet(t *testing.T, envName string) {
	t.Helper()
	require.NotEmptyf(t, os.Getenv(envName), "environment variable %v should not be empty", envName)
}
