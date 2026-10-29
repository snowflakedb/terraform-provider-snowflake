package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SuppressIfAny(t *testing.T) {
	require.Nil(t, SuppressIfAny())
	require.Nil(t, SuppressIfAny(nil, SuppressIfAny()))
	require.Equal(t, SuppressShow.Suppressors(), SuppressIfAny(SuppressShow).Suppressors())

	combined := SuppressIfAny(
		SuppressStatement,
		SuppressIfAny(SuppressShow, nil, SuppressDescribe.WithKey("definition")),
	)
	require.Equal(t, []DiffSuppressor{
		{Kind: SuppressStatement},
		{Kind: SuppressShow},
		{Kind: SuppressDescribe, Key: "definition"},
	}, combined.Suppressors())

	extended := SuppressIfAny(combined, SuppressShow.WithKey("other"))
	require.Len(t, combined.Suppressors(), 3)
	require.Len(t, extended.Suppressors(), 4)
	returned := combined.Suppressors()
	returned[0].Key = "changed"
	require.Empty(t, combined.Suppressors()[0].Key)
}

func Test_DiffSuppressKind_WithKey(t *testing.T) {
	suppressor := SuppressShow.WithKey("type")
	require.Equal(t, []DiffSuppressor{{Kind: SuppressShow, Key: "type"}}, suppressor.Suppressors())
	require.Equal(t, []DiffSuppressor{{Kind: SuppressShow}}, SuppressShow.Suppressors())
}
