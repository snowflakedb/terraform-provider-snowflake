package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToImageRepositoryEncryptionType(t *testing.T) {
	type test struct {
		input string
		want  ImageRepositoryEncryptionType
	}

	valid := []test{
		// case insensitive
		{input: "snowflake_full", want: ImageRepositoryEncryptionTypeSnowflakeFull},

		// Supported Values
		{input: "SNOWFLAKE_FULL", want: ImageRepositoryEncryptionTypeSnowflakeFull},
		{input: "SNOWFLAKE_SSE", want: ImageRepositoryEncryptionTypeSnowflakeSse},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToImageRepositoryEncryptionType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToImageRepositoryEncryptionType(tc.input)
			require.Error(t, err)
		})
	}
}
