package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestCreateTableColumnMaskingPolicyApplication(t *testing.T) {
	r := require.New(t)

	input := &snowflake.TableColumnMaskingPolicyApplicationCreateInput{
		TableColumnMaskingPolicyApplication: snowflake.TableColumnMaskingPolicyApplication{
			Table: &snowflake.SchemaObjectIdentifier{
				Database:   "db",
				Schema:     "schema",
				ObjectName: "table",
			},
			Column: "column",
			MaskingPolicy: &snowflake.SchemaObjectIdentifier{
				Database:   "db",
				Schema:     "schema",
				ObjectName: "mymaskingpolicy",
			},
		},
	}

	mb := snowflake.NewTableColumnMaskingPolicyApplicationManager()
	createStmt := mb.Create(input)
	r.Equal(`ALTER TABLE IF EXISTS "db"."schema"."table" MODIFY COLUMN "column" SET MASKING POLICY "db"."schema"."mymaskingpolicy";`, createStmt)
}

func TestDeleteTableColumnMaskingPolicyApplication(t *testing.T) {
	r := require.New(t)

	input := &snowflake.TableColumnMaskingPolicyApplicationDeleteInput{
		TableColumn: snowflake.TableColumn{
			Table: &snowflake.SchemaObjectIdentifier{
				Database:   "db",
				Schema:     "schema",
				ObjectName: "table",
			},
			Column: "column",
		},
	}

	mb := snowflake.NewTableColumnMaskingPolicyApplicationManager()
	dropStmt := mb.Delete(input)
	r.Equal(`ALTER TABLE IF EXISTS "db"."schema"."table" MODIFY COLUMN "column" UNSET MASKING POLICY;`, dropStmt)
}

func TestReadTableColumnMaskingPolicyApplication(t *testing.T) {
	r := require.New(t)

	input := &snowflake.TableColumnMaskingPolicyApplicationReadInput{
		Table: &snowflake.SchemaObjectIdentifier{
			Database:   "db",
			Schema:     "schema",
			ObjectName: "table",
		},
	}

	mb := snowflake.NewTableColumnMaskingPolicyApplicationManager()
	describeStmt := mb.Read(input)
	r.Equal(`DESCRIBE TABLE "db"."schema"."table" TYPE = COLUMNS;`, describeStmt)
}
