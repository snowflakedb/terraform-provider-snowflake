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
	t.Run("converts non-empty string", func(t *testing.T) {
		got, err := ctyValueToGo[string](cty.StringVal("en_nz"))

		require.NoError(t, err)
		require.Equal(t, "en_nz", got)
	})

	t.Run("converts empty string", func(t *testing.T) {
		got, err := ctyValueToGo[string](cty.StringVal(""))

		require.NoError(t, err)
		require.Equal(t, "", got)
	})

	t.Run("converts positive int", func(t *testing.T) {
		got, err := ctyValueToGo[int](cty.NumberIntVal(42))

		require.NoError(t, err)
		require.Equal(t, 42, got)
	})

	t.Run("converts zero int", func(t *testing.T) {
		got, err := ctyValueToGo[int](cty.NumberIntVal(0))

		require.NoError(t, err)
		require.Equal(t, 0, got)
	})

	t.Run("converts negative int", func(t *testing.T) {
		got, err := ctyValueToGo[int](cty.NumberIntVal(-7))

		require.NoError(t, err)
		require.Equal(t, -7, got)
	})

	t.Run("converts true bool", func(t *testing.T) {
		got, err := ctyValueToGo[bool](cty.BoolVal(true))

		require.NoError(t, err)
		require.True(t, got)
	})

	t.Run("converts false bool", func(t *testing.T) {
		got, err := ctyValueToGo[bool](cty.BoolVal(false))

		require.NoError(t, err)
		require.False(t, got)
	})

	t.Run("returns error for unsupported type", func(t *testing.T) {
		_, err := ctyValueToGo[float64](cty.NumberFloatVal(1.5))

		require.ErrorContains(t, err, "unsupported type")
	})
}
