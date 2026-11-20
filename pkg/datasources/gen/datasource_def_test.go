package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_selectsForDataSource(t *testing.T) {
	base := ssot.Object{
		Name:       "Thing",
		NamePlural: "Things",
		Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{{Name: "name", Readonly: true}}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkShow())},
	}

	t.Run("an object with SDK SHOW and show columns", func(t *testing.T) {
		require.True(t, selectsForDataSource(base))
	})
	t.Run("a variant is covered by its base", func(t *testing.T) {
		variant := base
		variant.Name = "ThingSecondary"
		variant.Variant = &ssot.ObjectVariant{Of: &base}
		require.False(t, selectsForDataSource(variant))
	})
	t.Run("a base-only object", func(t *testing.T) {
		baseOnly := base
		baseOnly.BaseVariant = &ssot.BaseVariant{NotResource: true}
		require.True(t, selectsForDataSource(baseOnly))
	})
	t.Run("the paired SHOW union counts as show columns", func(t *testing.T) {
		generic := base
		generic.BaseVariant = &ssot.BaseVariant{ShowGeneric: generic.Show}
		generic.Show = ssot.OutputSpec{}
		require.True(t, selectsForDataSource(generic))
	})
	t.Run("without show columns", func(t *testing.T) {
		noColumns := base
		noColumns.Show = ssot.OutputSpec{}
		require.False(t, selectsForDataSource(noColumns))
	})
	t.Run("without SDK SHOW", func(t *testing.T) {
		noSdk := base
		noSdk.Generators = nil
		require.False(t, selectsForDataSource(noSdk))
	})
}
