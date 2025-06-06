package sdk

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

// assertOptsInvalid could be reused in tests for other interfaces in sdk package.
func assertOptsInvalid(t *testing.T, opts validatable, expectedError error) {
	t.Helper()
	err := opts.validate()
	assert.Error(t, err)
	var sdkErr *Error
	if errors.As(err, &sdkErr) {
		errorWithoutFileInfo := errorFileInfoRegexp.ReplaceAllString(sdkErr.Error(), "")
		expectedErrorWithoutFileInfo := errorFileInfoRegexp.ReplaceAllString(expectedError.Error(), "")
		assert.Equal(t, expectedErrorWithoutFileInfo, errorWithoutFileInfo)
	} else {
		assert.Equal(t, expectedError, err)
	}
}

// assertOptsInvalidJoinedErrors could be reused in tests for other interfaces in sdk package.
func assertOptsInvalidJoinedErrors(t *testing.T, opts validatable, expectedErrors ...error) {
	t.Helper()
	err := opts.validate()
	assert.Error(t, err)
	for _, expectedError := range expectedErrors {
		var sdkErr *Error
		if errors.As(expectedError, &sdkErr) {
			expectedErrorWithoutFileInfo := errorFileInfoRegexp.ReplaceAllString(sdkErr.Error(), "")
			assert.Contains(t, err.Error(), expectedErrorWithoutFileInfo)
		} else {
			assert.Contains(t, err.Error(), expectedError.Error())
		}
	}
}

// assertOptsValid could be reused in tests for other interfaces in sdk package.
func assertOptsValid(t *testing.T, opts validatable) {
	t.Helper()
	err := opts.validate()
	assert.NoError(t, err)
}

// assertSQLEquals could be reused in tests for other interfaces in sdk package.
func assertSQLEquals(t *testing.T, opts any, format string, args ...any) {
	t.Helper()
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(format, args...), actual)
}

// assertOptsValidAndSQLEquals could be reused in tests for other interfaces in sdk package.
// It's a shorthand for assertOptsValid and assertSQLEquals.
func assertOptsValidAndSQLEquals(t *testing.T, opts validatable, format string, args ...any) {
	t.Helper()
	assertOptsValid(t, opts)
	assertSQLEquals(t, opts, format, args...)
}
