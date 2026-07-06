package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToPrivilege(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	valid := []test{
		// Case insensitive (privileges are normalized to upper-case).
		{input: "usage", want: "USAGE"},
		{input: "Create Schema", want: "CREATE SCHEMA"},

		// Supported values.
		{input: "USAGE", want: "USAGE"},
		{input: "SELECT", want: "SELECT"},
		{input: "CREATE SCHEMA", want: "CREATE SCHEMA"},
		{input: "IMPORTED PRIVILEGES", want: "IMPORTED PRIVILEGES"},
		{input: "ALL PRIVILEGES", want: "ALL PRIVILEGES"},

		// Dots and underscores are allowed (e.g. fully qualified privileges).
		{input: "SNOWFLAKE.CORE.SOME_PRIVILEGE", want: "SNOWFLAKE.CORE.SOME_PRIVILEGE"},
		{input: "create_table", want: "CREATE_TABLE"},
	}

	invalid := []test{
		{input: ""},
		{input: "USAGE;"},
		{input: "CREATE SCHEMA--"},
		{input: "SELECT; DROP TABLE users;--"},
		{input: "MONITOR USAGE; SELECT"},
		{input: "privilege,with,comma"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToPrivilege(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToPrivilege(tc.input)
			require.Error(t, err)
		})
	}
}
