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
		require.NotNil(t, req.Type)
		assert.Equal(t, sdk.DataType("NUMBER(20, 5)"), *req.Type)
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
