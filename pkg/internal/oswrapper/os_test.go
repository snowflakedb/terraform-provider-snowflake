package oswrapper_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testfiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/stretchr/testify/require"
)

func TestReadFileSafeFailsForFileThatIsTooBig(t *testing.T) {
	c := make([]byte, 11*1024*1024)
	configPath := testfiles.TestFile(t, "config", c)

	_, err := oswrapper.ReadFileSafe(configPath, false)
	require.ErrorContains(t, err, fmt.Sprintf("config file %s is too big - maximum allowed size is 10MB", configPath))
}

func TestReadFileSafeCanSkipPermissionVerification(t *testing.T) {
	exp := random.Bytes()
	path := testfiles.TestFileWithCustomPermissions(t, "config", exp, 0o755)
	act, err := oswrapper.ReadFileSafe(path, false)
	require.NoError(t, err)
	require.Equal(t, exp, act)
}

func TestReadFileSafeWithPermissionVerificationFailsForFileThatIsTooBig(t *testing.T) {
	c := make([]byte, 11*1024*1024)
	configPath := testfiles.TestFile(t, "config", c)

	_, err := oswrapper.ReadFileSafe(configPath, true)
	require.ErrorContains(t, err, fmt.Sprintf("config file %s is too big - maximum allowed size is 10MB", configPath))
}

func TestReadFileSafeFailsForFileThatDoesNotExist(t *testing.T) {
	configPath := "non-existing"
	_, err := oswrapper.ReadFileSafe(configPath, true)
	require.ErrorContains(t, err, fmt.Sprintf("reading information about the config file: stat %s: no such file or directory", configPath))
}

func TestReadFileSafeFailsForFileWithTooWidePermissions(t *testing.T) {
	if oswrapper.IsRunningOnWindows() {
		t.Skip("checking file permissions on Windows is currently done in manual tests package")
	}
	tests := []struct {
		permissions fs.FileMode
	}{
		{permissions: 0o707},
		{permissions: 0o706},
		{permissions: 0o705},
		{permissions: 0o704},
		{permissions: 0o703},
		{permissions: 0o702},
		{permissions: 0o701},

		{permissions: 0o770},
		{permissions: 0o760},
		{permissions: 0o750},
		{permissions: 0o740},
		{permissions: 0o730},
		{permissions: 0o720},
		{permissions: 0o710},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("reading file with too wide permissions %#o", tt.permissions), func(t *testing.T) {
			path := testfiles.TestFileWithCustomPermissions(t, "config", random.Bytes(), tt.permissions)
			_, err := oswrapper.ReadFileSafe(path, true)
			require.ErrorContains(t, err, fmt.Sprintf("config file %s has unsafe permissions", path))
		})
	}
}

func TestReadFileSafeFailsForFileWithTooRestrictivePermissions(t *testing.T) {
	if oswrapper.IsRunningOnWindows() {
		t.Skip("checking file permissions on Windows is currently done in manual tests package")
	}
	tests := []struct {
		permissions fs.FileMode
	}{
		{permissions: 0o300},
		{permissions: 0o200},
		{permissions: 0o100},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("reading file with too restrictive permissions %#o", tt.permissions), func(t *testing.T) {
			path := testfiles.TestFileWithCustomPermissions(t, "config", random.Bytes(), tt.permissions)
			_, err := oswrapper.ReadFileSafe(path, true)
			require.ErrorContains(t, err, fmt.Sprintf("open %s: permission denied", path))
		})
	}
}

func TestReadFileSafeReadsFileWithCorrectPermissions(t *testing.T) {
	if oswrapper.IsRunningOnWindows() {
		t.Skip("checking file permissions on Windows is currently done in manual tests package")
	}
	tests := []struct {
		permissions fs.FileMode
	}{
		{permissions: 0o700},
		{permissions: 0o600},
		{permissions: 0o500},
		{permissions: 0o400},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("reading file with correct permissions %#o", tt.permissions), func(t *testing.T) {
			path := testfiles.TestFileWithCustomPermissions(t, "config", random.Bytes(), tt.permissions)
			_, err := oswrapper.ReadFileSafe(path, true)
			require.NoError(t, err)
		})
	}
}

func TestStat(t *testing.T) {
	env := random.AlphaN(10)
	t.Setenv(env, "test")
	require.Equal(t, os.Getenv(env), oswrapper.Getenv(env))
}

func TestStatOnFileThatDoesNotExist(t *testing.T) {
	fileName := random.AlphaN(10)
	expVal, expErr := os.Stat(fileName)
	actVal, actErr := oswrapper.Stat(fileName)
	require.Equal(t, expVal, actVal)
	require.Equal(t, expErr, actErr)
}

func TestGetenv(t *testing.T) {
	env := random.AlphaN(10)
	t.Setenv(env, "test")
	require.Equal(t, os.Getenv(env), oswrapper.Getenv(env))
}

func TestGetenvBool(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{value: "TRUE", expected: true},
		{value: "true", expected: true},
		{value: "1", expected: true},
		{value: "FALSE", expected: false},
		{value: "false", expected: false},
		{value: "0", expected: false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("getting a boolean env with value %v", tt.value), func(t *testing.T) {
			env := random.AlphaN(10)
			t.Setenv(env, tt.value)
			actual, err := oswrapper.GetenvBool(env)
			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestGetenvBoolUnset(t *testing.T) {
	env := random.AlphaN(10)
	value, err := oswrapper.GetenvBool(env)
	require.NoError(t, err)
	require.False(t, value)
}

func TestGetenvBoolEmptyValue(t *testing.T) {
	env := random.AlphaN(10)
	t.Setenv(env, "")
	value, err := oswrapper.GetenvBool(env)
	require.NoError(t, err)
	require.False(t, value)
}

func TestGetenvBoolFailsForInvalidValue(t *testing.T) {
	env := random.AlphaN(10)
	t.Setenv(env, "invalid")
	_, err := oswrapper.GetenvBool(env)
	require.ErrorContains(t, err, "strconv.ParseBool: parsing \"invalid\": invalid syntax")
}

func TestLookupEnvOnSetVariable(t *testing.T) {
	env := random.AlphaN(10)
	t.Setenv(env, "test")
	expVal, expExist := os.LookupEnv(env)
	actVal, actExist := oswrapper.LookupEnv(env)
	require.Equal(t, expVal, actVal)
	require.Equal(t, expExist, actExist)
}

func TestLookupEnvOnUnsetVariable(t *testing.T) {
	env := random.AlphaN(10)
	expVal, expExist := os.LookupEnv(env)
	actVal, actExist := oswrapper.LookupEnv(env)
	require.Equal(t, expVal, actVal)
	require.Equal(t, expExist, actExist)
}

func TestUserHomeDir(t *testing.T) {
	expVal, expExist := os.UserHomeDir()
	actVal, actExist := oswrapper.UserHomeDir()
	require.Equal(t, expVal, actVal)
	require.Equal(t, expExist, actExist)
}
