package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHybridTableDetailsRow_SplitTypeAndCollation(t *testing.T) {
	t.Run("with utf8", func(t *testing.T) {
		row := hybridTableDetailsRow{
			Type: "VARCHAR(10) COLLATE 'utf8'",
		}

		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, "VARCHAR(10)", actualType)
		assert.Equal(t, "utf8", *actualCollation)
	})

	t.Run("with locale", func(t *testing.T) {
		row := hybridTableDetailsRow{
			Type: "VARCHAR(10) COLLATE 'en_US'",
		}

		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, "VARCHAR(10)", actualType)
		assert.Equal(t, "en_US", *actualCollation)
	})

	t.Run("with multiple specifiers", func(t *testing.T) {
		row := hybridTableDetailsRow{
			Type: "VARCHAR(10) COLLATE 'fr_CA-ai-pi-trim'",
		}

		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, "VARCHAR(10)", actualType)
		assert.Equal(t, "fr_CA-ai-pi-trim", *actualCollation)
	})

	t.Run("with empty collation", func(t *testing.T) {
		row := hybridTableDetailsRow{
			Type: "VARCHAR(10) COLLATE ''",
		}

		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, "VARCHAR(10)", actualType)
		assert.Equal(t, "", *actualCollation)
	})

	t.Run("without collation", func(t *testing.T) {
		row := hybridTableDetailsRow{
			Type: "NUMBER(38, 0)",
		}

		actualType, actualCollation := row.splitTypeAndCollation()
		assert.Equal(t, "NUMBER(38, 0)", actualType)
		assert.Nil(t, actualCollation)
	})
}
