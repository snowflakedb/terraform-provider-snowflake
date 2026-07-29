package sdk

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// All the contents of this file were added to be able to use them outside the sdk package (i.e. integration tests package).
// It was easier to do it that way, so that we do not include big rename changes in the first moving PR.

// ExecForTests is forwarding function for Client.exec (that is unexported), that some integration tests/helpers were using
// TODO: remove after we have all usages covered by SDK (for now it means implementing stages, tables, and tags)
func (c *Client) ExecForTests(ctx context.Context, sql string) (sql.Result, error) {
	return c.exec(ctx, sql)
}

// QueryOneForTests is forwarding function for Client.queryOne (that is unexported), that some integration tests/helpers were using
// TODO: remove after introducing all resources using this
func (c *Client) QueryOneForTests(ctx context.Context, dest interface{}, sql string) error {
	return c.queryOne(ctx, dest, sql)
}

// QueryForTests is forwarding function for Client.query (that is unexported), that some integration tests/helpers were using
// TODO: remove after introducing all resources using this
func (c *Client) QueryForTests(ctx context.Context, dest interface{}, sql string) error {
	return c.query(ctx, dest, sql)
}

func ErrorsEqual(t *testing.T, expected error, actual error) {
	t.Helper()
	var expectedErr *Error
	var actualErr *Error
	if errors.As(expected, &expectedErr) && errors.As(actual, &actualErr) {
		expectedErrorWithoutFileInfo := errorFileInfoRegexp.ReplaceAllString(expectedErr.Error(), "")
		errorWithoutFileInfo := errorFileInfoRegexp.ReplaceAllString(actualErr.Error(), "")
		assert.Equal(t, expectedErrorWithoutFileInfo, errorWithoutFileInfo)
	} else {
		assert.Equal(t, expected, actual)
	}
}

func ErrExactlyOneOf(structName string, fieldNames ...string) error {
	return errExactlyOneOf(structName, fieldNames...)
}

func ErrOneOf(structName string, fieldNames ...string) error {
	return errOneOf(structName, fieldNames...)
}
