package defs

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func TestAllValidates(t *testing.T) {
	All()
}

func TestSelectObjects(t *testing.T) {
	enabled := func(g ssot.Generators) bool { return g.ObjectParametersAssert.Enabled }
	t.Run("flag and filter", func(t *testing.T) {
		got := SelectObjects(enabled, func(obj ssot.Object) bool { return obj.Name != "Schema" })
		for _, obj := range got {
			require.NotEqual(t, "Schema", obj.Name)
			require.True(t, ssot.ObjectGenerators(obj.Generators...).ObjectParametersAssert.Enabled)
		}
		require.Contains(t, objectNames(got), "Warehouse")
	})
	t.Run("filter rejects all", func(t *testing.T) {
		require.Empty(t, SelectObjects(enabled, func(ssot.Object) bool { return false }))
	})
}

func objectNames(objs []ssot.Object) []string {
	names := make([]string, 0, len(objs))
	for _, obj := range objs {
		names = append(names, obj.Name)
	}
	return names
}
