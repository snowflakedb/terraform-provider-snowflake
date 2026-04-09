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

// RunTest is a generic helper that verifies a To<EnumName> conversion function
// against all enum values (currently without aliases).
// TODO [SNOW-2324252]: support aliases
func (p typedEnumTestProvider[T]) RunTest(t *testing.T) {
	t.Helper()

	t.Run(fmt.Sprintf("%s case-insensitive", p.enumName), func(t *testing.T) {
		require.NotEmpty(t, p.allValues, "allValidValues must not be empty")

		got, err := p.conversionFunc(strings.ToLower(string(p.allValues[0])))

		require.NoError(t, err)
		require.Equal(t, p.allValues[0], got)
	})

	for _, v := range p.allValues {
		t.Run(fmt.Sprintf("%s: %s", p.enumName, v), func(t *testing.T) {
			got, err := p.conversionFunc(string(v))

			require.NoError(t, err)
			require.Equal(t, v, got)
		})
	}

	for _, v := range []string{"", "foo"} {
		t.Run(fmt.Sprintf("%s invalid: %s", p.enumName, v), func(t *testing.T) {
			_, err := p.conversionFunc(v)

			require.Error(t, err)
		})
	}
}

func Test_AllEnumConversions(t *testing.T) {
	all := []enumTestProvider{
		typedEnumTestProvider[NetworkRuleMode]{"NetworkRuleMode", AllNetworkRuleModes, ToNetworkRuleMode},
		typedEnumTestProvider[ProgrammaticAccessTokenStatus]{"ProgrammaticAccessTokenStatus", AllProgrammaticAccessTokenStatuses, ToProgrammaticAccessTokenStatus},
	}

	for _, tp := range all {
		tp.RunTest(t)
	}
}
