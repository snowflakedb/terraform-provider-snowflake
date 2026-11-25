package resources

import (
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_virtualColumnExpressionsEqual(t *testing.T) {
	testCases := []struct {
		a, b     string
		expected bool
	}{
		{a: "ID * 2", b: "id*2", expected: true},
		{a: "ID + 1", b: `( "ID"  +   1 )`, expected: true},
		{a: "ID * 2", b: "((id * 2))", expected: true},
		{a: "(A) + (B)", b: "A + B", expected: false},
		{a: "ID * 2", b: "ID * 3", expected: false},
		{a: "", b: "", expected: true},
		{a: "", b: "ID", expected: false},
	}
	for _, tc := range testCases {
		t.Run(tc.a+" vs "+tc.b, func(t *testing.T) {
			assert.Equal(t, tc.expected, virtualColumnExpressionsEqual(tc.a, tc.b))
		})
	}
}

func Test_withRecreatedVirtualColumns(t *testing.T) {
	old := columns{{name: "ID"}, {name: "V", as: "ID * 2"}, {name: "W", as: "ID * 4"}}
	new := columns{{name: "ID"}, {name: "V", as: "id*3"}, {name: "W", as: "id * 4", comment: "changed"}}
	removed, added, changed := old.diffs(new)

	removed, added, changed = withRecreatedVirtualColumns(old, new, removed, added, changed)

	assert.Equal(t, columns{{name: "V", as: "ID * 2"}}, removed)
	assert.Equal(t, columns{{name: "V", as: "id*3"}}, added)
	for _, c := range changed {
		assert.NotEqual(t, "V", c.newColumn.name)
	}
	assert.True(t, slices.ContainsFunc(changed, func(c changedColumn) bool { return c.newColumn.name == "W" && c.changedComment }))
}

func Test_toColumnConfig_virtualColumn(t *testing.T) {
	got := toColumnConfig([]sdk.TableColumnDetails{
		{Name: "ID", Kind: sdk.TableColumnKindColumn, TypeRaw: "NUMBER(38,0)", Null: true},
		{Name: "V", Kind: sdk.TableColumnKindVirtual, TypeRaw: "NUMBER(38,0)", Null: true, Expression: sdk.String("ID * 2")},
		{Name: "OTHER", Kind: "SOMETHING_ELSE", TypeRaw: "NUMBER(38,0)"},
	})

	require.Len(t, got, 2)
	assert.NotContains(t, got[0].(map[string]any), "as")
	assert.Equal(t, "ID * 2", got[1].(map[string]any)["as"])
	assert.NotContains(t, got[1].(map[string]any), "default")
}
