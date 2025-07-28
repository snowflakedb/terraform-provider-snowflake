package sdk

import (
	"context"
	"os"
	"testing"

	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1827331]: Move the rest of the tests to testint package

func TestClient_ping(t *testing.T) {
	client := defaultTestClient(t)
	err := client.Ping()
	require.NoError(t, err)
}

func TestClient_close(t *testing.T) {
	client := defaultTestClient(t)
	err := client.Close()
	require.NoError(t, err)
}

func TestClient_exec(t *testing.T) {
	client := defaultTestClient(t)
	ctx := context.Background()
	_, err := client.exec(ctx, "SELECT 1")
	require.NoError(t, err)
}

func TestClient_query(t *testing.T) {
	client := defaultTestClient(t)
	ctx := context.Background()
	rows := []struct {
		One int `db:"ONE"`
	}{}
	err := client.query(ctx, &rows, "SELECT 1 AS ONE")
	require.NoError(t, err)
	require.NotNil(t, rows)
	require.Len(t, rows, 1)
	require.Equal(t, 1, rows[0].One)
}

func TestClient_queryOne(t *testing.T) {
	client := defaultTestClient(t)
	ctx := context.Background()
	row := struct {
		One int `db:"ONE"`
	}{}
	err := client.queryOne(ctx, &row, "SELECT 1 AS ONE")
	require.NoError(t, err)
	require.Equal(t, 1, row.One)
}

// TODO [SNOW-2054366]: Use dedicated users for these tests.
func TestClient_NewClientDriverLoggingLevel(t *testing.T) {
	t.Run("get default gosnowflake driver logging level", func(t *testing.T) {
		config := DefaultConfig()
		_, err := NewClient(config)
		require.NoError(t, err)

		var expected string
		if os.Getenv("GITHUB_ACTIONS") != "" {
			expected = "fatal"
		} else {
			expected = "error"
		}
		assert.Equal(t, expected, gosnowflake.GetLogger().GetLogLevel())
	})

	t.Run("set gosnowflake driver logging level with config", func(t *testing.T) {
		config := DefaultConfig()
		config.Tracing = "trace"
		_, err := NewClient(config)
		require.NoError(t, err)

		assert.Equal(t, "trace", gosnowflake.GetLogger().GetLogLevel())
	})
}
