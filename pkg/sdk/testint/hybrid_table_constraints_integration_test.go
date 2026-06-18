//go:build non_account_level_tests

package testint

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_HybridTableConstraints(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// Create a parent hybrid table: PK on id.
	parentId, parentCleanup := testClientHelper().HybridTable.CreateWithRequest(t,
		testClientHelper().Ids.RandomSchemaObjectIdentifier(),
		sdk.HybridTableColumnsConstraintsAndIndexesRequest{
			Columns: []sdk.HybridTableColumnRequest{
				{
					Name:     "id",
					DataType: sdk.DataType("NUMBER NOT NULL"),
					InlineConstraint: &sdk.ColumnInlineConstraint{
						Name: sdk.String("pk_parent"),
						Type: sdk.ColumnConstraintTypePrimaryKey,
					},
				},
				{
					Name:     "code",
					DataType: sdk.DataType("VARCHAR(50)"),
				},
			},
		},
	)
	t.Cleanup(parentCleanup)

	// Create the child table via raw SQL so we can control exact constraint names,
	// including an anonymous UNIQUE which Snowflake will assign a SYS_CONSTRAINT_* name.
	childId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
	t.Cleanup(testClientHelper().HybridTable.DropFunc(t, childId))

	createChildSQL := fmt.Sprintf(
		`CREATE HYBRID TABLE %s (
			id NUMBER NOT NULL,
			parent_id NUMBER NOT NULL,
			col_a VARCHAR(50),
			col_b VARCHAR(50),
			CONSTRAINT pk_child PRIMARY KEY (id),
			CONSTRAINT uq_named UNIQUE (col_a),
			UNIQUE (col_b),
			CONSTRAINT fk_named FOREIGN KEY (parent_id) REFERENCES %s (id)
		)`,
		childId.FullyQualifiedName(),
		parentId.FullyQualifiedName(),
	)
	_, err := client.ExecForTests(ctx, createChildSQL)
	require.NoError(t, err)

	constraints, err := client.HybridTables.GetConstraints(ctx, childId)
	require.NoError(t, err)

	t.Run("primary key", func(t *testing.T) {
		pks := collections.Filter(constraints, func(c sdk.HybridTableConstraint) bool { return c.Kind == sdk.ColumnConstraintTypePrimaryKey })
		require.Len(t, pks, 1, "expected exactly one PRIMARY KEY constraint")
		require.Equal(t, "PK_CHILD", pks[0].Name)
		require.Equal(t, []string{"ID"}, pks[0].Columns)
	})

	t.Run("unique keys", func(t *testing.T) {
		uqs := collections.Filter(constraints, func(c sdk.HybridTableConstraint) bool { return c.Kind == sdk.ColumnConstraintTypeUnique })
		require.Len(t, uqs, 2, "expected exactly two UNIQUE constraints")

		// Named UNIQUE
		namedUQ, err := collections.FindFirst(uqs, func(c sdk.HybridTableConstraint) bool { return c.Name == "UQ_NAMED" })
		require.NoError(t, err, "expected UNIQUE constraint UQ_NAMED")
		require.Equal(t, []string{"COL_A"}, namedUQ.Columns)

		// Anonymous UNIQUE — Snowflake assigns a SYS_CONSTRAINT_* name.
		var anonUQ *sdk.HybridTableConstraint
		for i := range uqs {
			if strings.HasPrefix(uqs[i].Name, "SYS_CONSTRAINT_") {
				anonUQ = &uqs[i]
				break
			}
		}
		require.NotNil(t, anonUQ, "expected an anonymous UNIQUE constraint with SYS_CONSTRAINT_* name")
		require.Equal(t, []string{"COL_B"}, anonUQ.Columns)
	})

	t.Run("foreign keys", func(t *testing.T) {
		fks := collections.Filter(constraints, func(c sdk.HybridTableConstraint) bool { return c.Kind == sdk.ColumnConstraintTypeForeignKey })
		require.Len(t, fks, 1, "expected exactly one FOREIGN KEY constraint")
		require.Equal(t, "FK_NAMED", fks[0].Name)
		require.Equal(t, []string{"PARENT_ID"}, fks[0].Columns)
		require.Equal(t, []string{"ID"}, fks[0].ReferencedColumns)
		require.Equal(t, parentId.DatabaseName(), fks[0].ReferencedTable.DatabaseName())
		require.Equal(t, parentId.SchemaName(), fks[0].ReferencedTable.SchemaName())
		require.Equal(t, parentId.Name(), fks[0].ReferencedTable.Name())
	})
}
