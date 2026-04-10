package sdk

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// allEnumConversionTests keeps all enums which conversion needs to be tested
// Currently, it's populated from the list inside Test_AllEnumConversions.
// Ultimately, it will be populated also through init method in the object's test file from each enum existing in the definition.
var allEnumConversionTests []enumTestProvider

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
	additional := []enumTestProvider{
		typedEnumTestProvider[NetworkRuleMode]{"NetworkRuleMode", AllNetworkRuleModes, ToNetworkRuleMode},
	}

	require.Positive(t, len(allEnumConversionTests))

	allEnumConversionTests = append(allEnumConversionTests, additional...)

	for _, tp := range allEnumConversionTests {
		tp.RunTest(t)
	}
}
