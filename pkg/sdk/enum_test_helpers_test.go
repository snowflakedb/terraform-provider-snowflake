package sdk

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// testEnumConversion is a generic helper that verifies a To<EnumName> conversion function
// against all enum values (currently without aliases).
// TODO [SNOW-2324252]: support aliases
func testEnumConversion[T ~string](t *testing.T, allValues []T, convert func(string) (T, error)) {
	t.Helper()

	t.Run("case insensitive", func(t *testing.T) {
		require.NotEmpty(t, allValues, "allValidValues must not be empty")

		got, err := convert(strings.ToLower(string(allValues[0])))

		require.NoError(t, err)
		require.Equal(t, allValues[0], got)
	})

	for _, v := range allValues {
		t.Run(fmt.Sprintf("%s", v), func(t *testing.T) {
			got, err := convert(string(v))

			require.NoError(t, err)
			require.Equal(t, v, got)
		})
	}

	for _, v := range []string{"", "foo"} {
		t.Run("invalid: "+v, func(t *testing.T) {
			_, err := convert(v)

			require.Error(t, err)
		})
	}
}
