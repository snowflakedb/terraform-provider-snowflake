package gen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_defaultExtHooks_options(t *testing.T) {
	hooks := defaultExtHooks(
		WithBeforeCreate(generated()),
		WithFqnSet(omitted()),
	)

	require.Equal(t, generated(), hooks.BeforeCreate)
	require.Equal(t, omitted(), hooks.FqnSet)
	require.Equal(t, extCall(), hooks.ParseId)
}
