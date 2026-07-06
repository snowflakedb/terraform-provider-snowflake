package snowflake_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestCreateTableColumnMaskingPolicyApplication(t *testing.T) {
	tests := []struct {
		name     string
		input    *snowflake.TableColumnMaskingPolicyApplicationCreateInput
		expected string
	}{
		{
			name: "basic",
			input: &snowflake.TableColumnMaskingPolicyApplicationCreateInput{
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
			},
			expected: `ALTER TABLE IF EXISTS "db"."schema"."table" MODIFY COLUMN "column" SET MASKING POLICY "db"."schema"."mymaskingpolicy";`,
		},
		{
			name: "identifiers with double quotes are escaped",
			input: &snowflake.TableColumnMaskingPolicyApplicationCreateInput{
				TableColumnMaskingPolicyApplication: snowflake.TableColumnMaskingPolicyApplication{
					Table: &snowflake.SchemaObjectIdentifier{
						Database:   `d"b`,
						Schema:     `sch"ema`,
						ObjectName: `tab"le`,
					},
					Column: `col"umn`,
					MaskingPolicy: &snowflake.SchemaObjectIdentifier{
						Database:   `d"b`,
						Schema:     `sch"ema`,
						ObjectName: `mymasking"policy`,
					},
				},
			},
			expected: `ALTER TABLE IF EXISTS "d""b"."sch""ema"."tab""le" MODIFY COLUMN "col""umn" SET MASKING POLICY "d""b"."sch""ema"."mymasking""policy";`,
		},
	}

	mb := snowflake.NewTableColumnMaskingPolicyApplicationManager()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, mb.Create(tt.input))
		})
	}
}

func TestDeleteTableColumnMaskingPolicyApplication(t *testing.T) {
	tests := []struct {
		name     string
		input    *snowflake.TableColumnMaskingPolicyApplicationDeleteInput
		expected string
	}{
		{
			name: "basic",
			input: &snowflake.TableColumnMaskingPolicyApplicationDeleteInput{
				TableColumn: snowflake.TableColumn{
					Table: &snowflake.SchemaObjectIdentifier{
						Database:   "db",
						Schema:     "schema",
						ObjectName: "table",
					},
					Column: "column",
				},
			},
			expected: `ALTER TABLE IF EXISTS "db"."schema"."table" MODIFY COLUMN "column" UNSET MASKING POLICY;`,
		},
		{
			name: "identifiers with double quotes are escaped",
			input: &snowflake.TableColumnMaskingPolicyApplicationDeleteInput{
				TableColumn: snowflake.TableColumn{
					Table: &snowflake.SchemaObjectIdentifier{
						Database:   `d"b`,
						Schema:     `sch"ema`,
						ObjectName: `tab"le`,
					},
					Column: `col"umn`,
				},
			},
			expected: `ALTER TABLE IF EXISTS "d""b"."sch""ema"."tab""le" MODIFY COLUMN "col""umn" UNSET MASKING POLICY;`,
		},
	}

	mb := snowflake.NewTableColumnMaskingPolicyApplicationManager()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, mb.Delete(tt.input))
		})
	}
}

func TestReadTableColumnMaskingPolicyApplication(t *testing.T) {
	tests := []struct {
		name     string
		input    *snowflake.TableColumnMaskingPolicyApplicationReadInput
		expected string
	}{
		{
			name: "basic",
			input: &snowflake.TableColumnMaskingPolicyApplicationReadInput{
				Table: &snowflake.SchemaObjectIdentifier{
					Database:   "db",
					Schema:     "schema",
					ObjectName: "table",
				},
			},
			expected: `DESCRIBE TABLE "db"."schema"."table" TYPE = COLUMNS;`,
		},
		{
			name: "identifiers with double quotes are escaped",
			input: &snowflake.TableColumnMaskingPolicyApplicationReadInput{
				Table: &snowflake.SchemaObjectIdentifier{
					Database:   `d"b`,
					Schema:     `sch"ema`,
					ObjectName: `tab"le`,
				},
			},
			expected: `DESCRIBE TABLE "d""b"."sch""ema"."tab""le" TYPE = COLUMNS;`,
		},
	}

	mb := snowflake.NewTableColumnMaskingPolicyApplicationManager()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, mb.Read(tt.input))
		})
	}
}
