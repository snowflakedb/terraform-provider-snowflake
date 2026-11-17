package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithDataSource(t *testing.T) {
	t.Run("is off by default", func(t *testing.T) {
		require.False(t, ObjectGenerators().DataSource.Enabled)
	})
	t.Run("with sets Enabled", func(t *testing.T) {
		require.True(t, ObjectGenerators(WithDataSource()).DataSource.Enabled)
	})
	t.Run("without turns it off and a later option wins", func(t *testing.T) {
		require.False(t, ObjectGenerators(WithDataSource(), WithoutDataSource()).DataSource.Enabled)
		require.True(t, ObjectGenerators(WithoutDataSource(), WithDataSource()).DataSource.Enabled)
	})
	t.Run("is not part of WithAllGenerators yet", func(t *testing.T) {
		require.False(t, ObjectGenerators(WithAllGenerators()).DataSource.Enabled)
	})
	t.Run("does not touch the other generators", func(t *testing.T) {
		g := ObjectGenerators(WithAllGenerators(), WithDataSource())
		require.True(t, g.Resource.Enabled)
		require.True(t, g.DataSource.Enabled)
		g = ObjectGenerators(WithAllGenerators(WithoutDataSource()))
		require.True(t, g.Resource.Enabled)
	})
}
