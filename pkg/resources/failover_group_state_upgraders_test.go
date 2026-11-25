package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_v2_24_0_FailoverGroupStateUpgrader(t *testing.T) {
	upgrade := func(state map[string]any) (map[string]any, error) {
		return v2_24_0_FailoverGroupStateUpgrader(context.Background(), state, nil)
	}

	t.Run("nil state", func(t *testing.T) {
		state, err := upgrade(nil)
		require.NoError(t, err)
		assert.Nil(t, state)
	})

	t.Run("no schedule and no from_replica", func(t *testing.T) {
		state, err := upgrade(map[string]any{
			"name":                 "FG",
			"replication_schedule": []any{},
			"from_replica":         []any{},
			"ignore_edition_check": false,
		})
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"name": "FG", "ignore_edition_check": false}, state)
	})

	t.Run("interval", func(t *testing.T) {
		state, err := upgrade(map[string]any{
			"replication_schedule": []any{map[string]any{"interval": float64(10), "cron": []any{}}},
		})
		require.NoError(t, err)
		assert.Equal(t, "10 MINUTE", state["replication_schedule"])
	})

	t.Run("cron", func(t *testing.T) {
		state, err := upgrade(map[string]any{
			"replication_schedule": []any{map[string]any{
				"interval": float64(0),
				"cron":     []any{map[string]any{"expression": "0 0 10-20 * TUE,THU", "time_zone": "UTC"}},
			}},
		})
		require.NoError(t, err)
		assert.Equal(t, "USING CRON 0 0 10-20 * TUE,THU UTC", state["replication_schedule"])
	})

	t.Run("empty schedule block", func(t *testing.T) {
		state, err := upgrade(map[string]any{
			"replication_schedule": []any{map[string]any{"interval": float64(0), "cron": []any{}}},
		})
		require.NoError(t, err)
		assert.NotContains(t, state, "replication_schedule")
	})

	t.Run("from_replica is refused", func(t *testing.T) {
		_, err := upgrade(map[string]any{
			"from_replica": []any{map[string]any{"name": "FG", "organization_name": "ORG", "source_account_name": "ACC"}},
		})
		require.ErrorContains(t, err, "snowflake_failover_group_secondary")
	})
}
