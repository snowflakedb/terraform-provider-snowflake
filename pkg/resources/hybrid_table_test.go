package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildHybridAlterColumnTypeAction(t *testing.T) {
	t.Run("valid type produces request with parsed type", func(t *testing.T) {
		col := hybridTableColumn{
			name:     "AGE",
			dataType: "NUMBER(20,5)",
		}

		req, err := buildHybridAlterColumnTypeAction(col)
		require.NoError(t, err)
		require.NotNil(t, req)

		assert.Equal(t, "AGE", req.ColumnName)
		require.NotNil(t, req.DataType)
		assert.Equal(t, sdk.DataType("NUMBER(20, 5)"), *req.DataType)
	})

	t.Run("invalid type returns error mentioning column name", func(t *testing.T) {
		col := hybridTableColumn{
			name:     "WAT",
			dataType: "NOT_A_REAL_TYPE",
		}

		req, err := buildHybridAlterColumnTypeAction(col)
		require.Error(t, err)
		assert.Nil(t, req)
		assert.Contains(t, err.Error(), "WAT")
	})

	t.Run("empty type returns error", func(t *testing.T) {
		col := hybridTableColumn{
			name:     "X",
			dataType: "",
		}

		req, err := buildHybridAlterColumnTypeAction(col)
		require.Error(t, err)
		assert.Nil(t, req)
	})
}

func Test_buildHybridColumnSpec(t *testing.T) {
	t.Run("plain column with no default", func(t *testing.T) {
		col := hybridTableColumn{
			name:     "ID",
			dataType: "NUMBER(38,0)",
			nullable: true,
		}

		spec, err := buildHybridColumnSpec(col)
		require.NoError(t, err)

		require.NotNil(t, spec.dataType)
		assert.Equal(t, "NUMBER(38,0)", spec.dataType.Canonical())
		assert.Nil(t, spec.defaultValue)
		assert.Empty(t, spec.collate)
		assert.Empty(t, spec.comment)
	})

	t.Run("column with constant default and collate", func(t *testing.T) {
		def := "hello"
		col := hybridTableColumn{
			name:     "NAME",
			dataType: "VARCHAR(200)",
			nullable: true,
			_default: &columnDefault{constant: &def},
			collate:  "en-ci",
			comment:  "the name",
		}

		spec, err := buildHybridColumnSpec(col)
		require.NoError(t, err)

		require.NotNil(t, spec.defaultValue)
		require.NotNil(t, spec.defaultValue.Expression)
		assert.Equal(t, "'hello'", *spec.defaultValue.Expression)
		assert.Equal(t, "en-ci", spec.collate)
		assert.Equal(t, "the name", spec.comment)
	})

	t.Run("invalid data type returns error", func(t *testing.T) {
		col := hybridTableColumn{
			name:     "BAD",
			dataType: "NOT_A_REAL_TYPE",
		}

		_, err := buildHybridColumnSpec(col)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "BAD")
	})

	t.Run("default block with multiple fields set returns error", func(t *testing.T) {
		c := "x"
		e := "f()"
		col := hybridTableColumn{
			name:     "X",
			dataType: "VARCHAR(10)",
			nullable: true,
			_default: &columnDefault{constant: &c, expression: &e},
		}
		_, err := buildHybridColumnSpec(col)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exactly one")
	})
}

func Test_hybridTableColumnTypeStateFunc(t *testing.T) {
	cases := []struct {
		input    string
		wantSame bool
	}{
		{"VARCHAR(256)", true},
		{"NUMBER(38,0)", false}, // normalised to "NUMBER(38, 0)" (added space)
		{"TEXT", false},         // expanded to "TEXT(16777216)"
		{"BOOLEAN", true},
		{"VARCHAR", false},
		{"NUMBER", false},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := DataTypeStateFunc(tc.input)
			if tc.wantSame {
				require.Equal(t, tc.input, got)
			} else {
				require.NotEqual(t, tc.input, got)
			}
			gotTwice := DataTypeStateFunc(got)
			require.Equal(t, got, gotTwice)
		})
	}
}
