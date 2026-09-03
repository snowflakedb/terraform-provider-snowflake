//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Definitions are global, immutable and Snowflake-managed, so there is nothing to create or clean up.
// Which ones exist differs between accounts, so these assert shape rather than specific names.
func TestInt_OpenflowConnectorDefinitions(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	client := testClient(t)
	ctx := testContext(t)

	t.Run("show", func(t *testing.T) {
		definitions, err := client.OpenflowConnectorDefinitions.Show(ctx, sdk.NewShowOpenflowConnectorDefinitionRequest())
		require.NoError(t, err)
		require.NotEmpty(t, definitions)

		// NotEmpty rather than HasName(d.Name), which would compare a field against itself and pass even if
		// the column never mapped. Every field below is populated on every definition, max_node_count and
		// categories included, so this runs over all of them rather than an arbitrary one.
		for _, d := range definitions {
			assertThatObject(
				t, objectassert.OpenflowConnectorDefinitionFromObject(t, &d).
					HasNameNotEmpty().
					HasProviderNotEmpty().
					HasVersionNotEmpty().
					HasDescriptionNotEmpty().
					HasDisplayNameNotEmpty().
					HasCategoriesNotEmpty().
					HasMaxNodeCountPositive(),
			)
		}

		// min_runtime_node_type is NULL for some definitions, so this is the one field that cannot be asserted
		// per row. Requiring at least one mapped value still fails on a wrong db tag, which would leave every
		// row nil.
		withNodeType := collections.Filter(definitions, func(d sdk.OpenflowConnectorDefinition) bool {
			return d.MinRuntimeNodeType != nil && *d.MinRuntimeNodeType != ""
		})
		assert.NotEmpty(t, withNodeType, "no definition mapped min_runtime_node_type, which a wrong db tag would also produce")
	})

	t.Run("show: with like", func(t *testing.T) {
		definitions, err := client.OpenflowConnectorDefinitions.Show(ctx, sdk.NewShowOpenflowConnectorDefinitionRequest())
		require.NoError(t, err)
		require.NotEmpty(t, definitions)
		expected := definitions[0]

		filtered, err := client.OpenflowConnectorDefinitions.Show(ctx, sdk.NewShowOpenflowConnectorDefinitionRequest().
			WithLike(sdk.Like{Pattern: sdk.String(expected.Name)}))
		require.NoError(t, err)
		require.NotEmpty(t, filtered)
		assert.Contains(t, collections.Map(filtered, func(d sdk.OpenflowConnectorDefinition) string { return d.Name }), expected.Name)
	})

	t.Run("show: with limit", func(t *testing.T) {
		definitions, err := client.OpenflowConnectorDefinitions.Show(ctx, sdk.NewShowOpenflowConnectorDefinitionRequest())
		require.NoError(t, err)
		require.NotEmpty(t, definitions)

		limited, err := client.OpenflowConnectorDefinitions.Show(ctx, sdk.NewShowOpenflowConnectorDefinitionRequest().
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		assert.Len(t, limited, 1)
		assert.LessOrEqual(t, len(limited), len(definitions), "LIMIT 1 returned more rows than an unlimited SHOW")
	})

	// Snowflake accepts STARTS WITH here but does not filter on it. LIKE and LIMIT filter on the same command
	// and STARTS WITH works on the other three SHOWs, so the clause reaches the server intact. This fails
	// once it is applied, which is the signal to assert real filtering instead.
	t.Run("show: with starts with is accepted but not applied", func(t *testing.T) {
		all, err := client.OpenflowConnectorDefinitions.Show(ctx, sdk.NewShowOpenflowConnectorDefinitionRequest())
		require.NoError(t, err)
		require.NotEmpty(t, all)

		unmatched, err := client.OpenflowConnectorDefinitions.Show(ctx, sdk.NewShowOpenflowConnectorDefinitionRequest().
			WithStartsWith("ZZ_NO_DEFINITION_STARTS_WITH_THIS"))
		require.NoError(t, err)
		assert.Len(t, unmatched, len(all),
			"STARTS WITH is currently ignored by SHOW OPENFLOW CONNECTOR DEFINITIONS; if this now filters, the "+
				"server gap is fixed and this case should assert real filtering instead")

		// The rows still have to convert.
		for _, d := range unmatched {
			assert.NotEmpty(t, d.Name)
		}
	})

	t.Run("show: with like matching nothing", func(t *testing.T) {
		definitions, err := client.OpenflowConnectorDefinitions.Show(ctx, sdk.NewShowOpenflowConnectorDefinitionRequest().
			WithLike(sdk.Like{Pattern: sdk.String("does_not_exist_%")}))
		require.NoError(t, err)
		assert.Empty(t, definitions)
	})
}
