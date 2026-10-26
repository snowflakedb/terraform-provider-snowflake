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
	require.Equal(t, Omitted(), hooks.DropSafelyInSdk)
	require.Equal(t, Generated(), DefaultResourceExtHooks(WithDropSafelyInSdk(Generated())).DropSafelyInSdk)
	require.Equal(t, Omitted(), hooks.StateUpgraders)
	require.Equal(t, Generated(), DefaultResourceExtHooks(WithStateUpgraders(Generated())).StateUpgraders)
	require.Equal(t, Omitted(), hooks.AfterCreate)
	require.Equal(t, Omitted(), hooks.AfterAlter)
	require.Equal(t, ExtCall(), hooks.CustomizeDiff)
	require.Equal(t, Generated(), DefaultResourceExtHooks(WithCustomizeDiff(Generated())).CustomizeDiff)
}

func Test_DefaultResourceExtHooks_afterHooks(t *testing.T) {
	hooks := DefaultResourceExtHooks(
		WithAfterCreate(ExtCall()),
		WithAfterAlter(ExtCall()),
	)

	require.Equal(t, ExtCall(), hooks.AfterCreate)
	require.Equal(t, ExtCall(), hooks.AfterAlter)
}

func Test_DefaultResourceExtHooks_outputSetDefaults(t *testing.T) {
	hooks := DefaultResourceExtHooks()
	require.Equal(t, Generated(), hooks.ShowOutputSet)
	require.Equal(t, Generated(), hooks.DescribeOutputSet)
	require.Equal(t, Generated(), hooks.FqnSet)

	t.Run("show output is omitted without a show output schema", func(t *testing.T) {
		hooks := DefaultResourceExtHooks(WithShowOutputSchema(Omitted()))
		require.Equal(t, Omitted(), hooks.ShowOutputSet)
		require.Equal(t, Generated(), hooks.DescribeOutputSet)
	})

	t.Run("describe output is omitted without the describe step", func(t *testing.T) {
		require.Equal(t, Omitted(), DefaultResourceExtHooks(WithDescribeInSdk(Omitted())).DescribeOutputSet)
	})

	t.Run("describe output is omitted without its schema", func(t *testing.T) {
		require.Equal(t, Omitted(), DefaultResourceExtHooks(WithDescribeOutputSchema(Omitted())).DescribeOutputSet)
	})

	t.Run("a hand-written set describe output keeps the generated set", func(t *testing.T) {
		require.Equal(t, Generated(), DefaultResourceExtHooks(WithSetDescribeOutput(ExtCall())).DescribeOutputSet)
	})

	t.Run("a custom describe schema keeps the generated set", func(t *testing.T) {
		require.Equal(t, Generated(), DefaultResourceExtHooks(WithDescribeOutputSchema(ExtCall())).DescribeOutputSet)
	})

	t.Run("explicit ext call stays", func(t *testing.T) {
		hooks := DefaultResourceExtHooks(
			WithShowOutputSchema(Omitted()),
			WithShowOutputSet(ExtCall()),
			WithDescribeInSdk(Omitted()),
			WithDescribeOutputSet(ExtCall()),
		)
		require.Equal(t, ExtCall(), hooks.ShowOutputSet)
		require.Equal(t, ExtCall(), hooks.DescribeOutputSet)
	})
}
