//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TestInt_HybridTableShowIndexes proves the existing ShowIndexes method returns the
// shapes the resource Read path relies on (settled in the design's step-1 e2e capture):
//   - user secondary indexes report is_unique = "N" (IsUnique points to false);
//   - the PK- and UNIQUE-backed indexes report "Y" (IsUnique points to true);
//   - columns / included_columns are bracketed, uppercase, comma-space lists.
//
// No SDK production code is exercised beyond the already-generated ShowIndexes.
func TestInt_HybridTableShowIndexes(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	id, cleanup := testClientHelper().HybridTable.CreateWithRequest(t,
		testClientHelper().Ids.RandomSchemaObjectIdentifier(),
		sdk.HybridTableColumnsConstraintsAndIndexesRequest{
			Columns: []sdk.HybridTableColumnRequest{
				{
					Name:     "id",
					DataType: sdk.DataType("NUMBER NOT NULL"),
					InlineConstraint: &sdk.ColumnInlineConstraint{
						Name: sdk.String("pk_fixture"),
						Type: sdk.ColumnConstraintTypePrimaryKey,
					},
				},
				{
					Name:     "email",
					DataType: sdk.DataType("VARCHAR(256)"),
					InlineConstraint: &sdk.ColumnInlineConstraint{
						Name: sdk.String("uq_email"),
						Type: sdk.ColumnConstraintTypeUnique,
					},
				},
				{
					Name:     "status",
					DataType: sdk.DataType("VARCHAR(50)"),
				},
				{
					Name:     "region",
					DataType: sdk.DataType("VARCHAR(50)"),
				},
				{
					Name:     "score",
					DataType: sdk.DataType("NUMBER(38,0)"),
				},
			},
			OutOfLineIndex: []sdk.HybridTableOutOfLineIndexRequest{
				*sdk.NewHybridTableOutOfLineIndexRequest("idx_status", []string{"status"}),
				*sdk.NewHybridTableOutOfLineIndexRequest("idx_region_inc", []string{"region"}).
					WithIncludeColumns([]string{"score"}),
			},
		},
	)
	t.Cleanup(cleanup)

	indexes, err := client.HybridTables.ShowIndexes(ctx,
		sdk.NewShowIndexesHybridTableRequest().WithIn(sdk.TableIn{Table: id}))
	require.NoError(t, err)

	t.Run("user secondary indexes report is_unique false", func(t *testing.T) {
		secondary := collections.Filter(indexes, func(i sdk.HybridTableIndex) bool {
			return i.IsUnique != nil && !*i.IsUnique
		})
		require.Len(t, secondary, 2, "expected exactly two user secondary indexes (is_unique=N)")

		status, err := collections.FindFirst(secondary, func(i sdk.HybridTableIndex) bool { return i.Name == "IDX_STATUS" })
		require.NoError(t, err, "expected secondary index IDX_STATUS")
		require.NotNil(t, status.Columns)
		require.Equal(t, "[STATUS]", *status.Columns)
		// Snowflake returns the literal "[]" (not NULL/empty) when an index has no INCLUDE columns.
		require.Equal(t, "[]", status.IncludedColumns)

		regionInc, err := collections.FindFirst(secondary, func(i sdk.HybridTableIndex) bool { return i.Name == "IDX_REGION_INC" })
		require.NoError(t, err, "expected secondary index IDX_REGION_INC")
		require.NotNil(t, regionInc.Columns)
		require.Equal(t, "[REGION]", *regionInc.Columns)
		require.Equal(t, "[SCORE]", regionInc.IncludedColumns)
	})

	t.Run("PK- and UNIQUE-backed indexes report is_unique true", func(t *testing.T) {
		unique := collections.Filter(indexes, func(i sdk.HybridTableIndex) bool {
			return i.IsUnique != nil && *i.IsUnique
		})
		// At least the PK-backed and the UNIQUE-backed indexes. Assert the discriminator
		// excludes them from the secondary-index set rather than an exact count, because
		// constraint-backed index naming is not part of this feature's contract.
		require.GreaterOrEqual(t, len(unique), 2, "expected PK- and UNIQUE-backed indexes to report is_unique=Y")
		for _, i := range unique {
			require.NotEqual(t, "IDX_STATUS", i.Name)
			require.NotEqual(t, "IDX_REGION_INC", i.Name)
		}
	})
}
