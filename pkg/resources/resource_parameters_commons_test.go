package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/stretchr/testify/require"
)

func Test_enrichWithReferenceToParameterDocs(t *testing.T) {
	t.Run("formats the docs reference correctly", func(t *testing.T) {
		description := random.Comment()

		enrichedDescription := enrichWithReferenceToParameterDocs(sdk.UserParameterAbortDetachedQuery, description)

		require.Equal(t, description+" "+"For more information, check [ABORT_DETACHED_QUERY docs](https://docs.snowflake.com/en/sql-reference/parameters#abort-detached-query).", enrichedDescription)
	})
}

func Test_ctyValueToGo(t *testing.T) {
	testCases := []struct {
		Name     string
		Convert  func() (any, error)
		Expected any
	}{
		{Name: "string/non-empty", Convert: func() (any, error) { return ctyValueToGo[string](cty.StringVal("en_nz")) }, Expected: "en_nz"},
		{Name: "string/empty", Convert: func() (any, error) { return ctyValueToGo[string](cty.StringVal("")) }, Expected: ""},
		{Name: "int/positive", Convert: func() (any, error) { return ctyValueToGo[int](cty.NumberIntVal(42)) }, Expected: 42},
		{Name: "int/zero", Convert: func() (any, error) { return ctyValueToGo[int](cty.NumberIntVal(0)) }, Expected: 0},
		{Name: "int/negative", Convert: func() (any, error) { return ctyValueToGo[int](cty.NumberIntVal(-7)) }, Expected: -7},
		{Name: "bool/true", Convert: func() (any, error) { return ctyValueToGo[bool](cty.BoolVal(true)) }, Expected: true},
		{Name: "bool/false", Convert: func() (any, error) { return ctyValueToGo[bool](cty.BoolVal(false)) }, Expected: false},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := tc.Convert()

			require.NoError(t, err)
			require.Equal(t, tc.Expected, got)
		})
	}
}
