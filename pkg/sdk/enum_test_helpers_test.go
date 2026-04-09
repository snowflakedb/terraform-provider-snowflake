package sdk

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type enumTestProvider interface {
	RunTest(t *testing.T)
}

type typedEnumTestProvider[T ~string] struct {
	enumName       string
	allValues      []T
	conversionFunc func(string) (T, error)
}

func (p typedEnumTestProvider[T]) RunTest(t *testing.T) {
	testEnumConversion(t, p.enumName, p.allValues, p.conversionFunc)
}

func Test_AllEnumConversions(t *testing.T) {
	all := []enumTestProvider{
		typedEnumTestProvider[ProgrammaticAccessTokenStatus]{"ProgrammaticAccessTokenStatus", AllProgrammaticAccessTokenStatuses, ToProgrammaticAccessTokenStatus},
	}

	for _, tp := range all {
		tp.RunTest(t)
	}
}

// testEnumConversion is a generic helper that verifies a To<EnumName> conversion function
// against all enum values (currently without aliases).
// TODO [SNOW-2324252]: support aliases
func testEnumConversion[T ~string](t *testing.T, enumName string, allValues []T, convert func(string) (T, error)) {
	t.Helper()

	t.Run(fmt.Sprintf("%s case-insensitive", enumName), func(t *testing.T) {
		require.NotEmpty(t, allValues, "allValidValues must not be empty")

		got, err := convert(strings.ToLower(string(allValues[0])))

		require.NoError(t, err)
		require.Equal(t, allValues[0], got)
	})

	for _, v := range allValues {
		t.Run(fmt.Sprintf("%s: %s", enumName, v), func(t *testing.T) {
			got, err := convert(string(v))

			require.NoError(t, err)
			require.Equal(t, v, got)
		})
	}

	for _, v := range []string{"", "foo"} {
		t.Run(fmt.Sprintf("%s invalid: %s", enumName, v), func(t *testing.T) {
			_, err := convert(v)

			require.Error(t, err)
		})
	}
}
