package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DefaultResourceExtHooks_options(t *testing.T) {
	hooks := DefaultResourceExtHooks(
		WithBeforeCreate(Generated()),
		WithFqnSet(Omitted()),
	)

	require.Equal(t, Generated(), hooks.BeforeCreate)
	require.Equal(t, Omitted(), hooks.FqnSet)
	require.Equal(t, ExtCall(), hooks.ParseId)
}
