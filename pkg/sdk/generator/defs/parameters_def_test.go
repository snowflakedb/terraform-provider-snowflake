package defs

import (
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/stretchr/testify/require"
)

func TestAllParameters_wellFormed(t *testing.T) {
	for _, p := range AllParameters {
		t.Run(p.SqlName, func(t *testing.T) {
			if p.SqlName == "" {
				t.Error("parameter with empty SqlName")
			}
			if p.Kind == "" {
				t.Error("empty Kind")
			}
			if strings.HasPrefix(p.Kind, "*") {
				t.Errorf("pointer Kind %q; catalog kinds must be base types", p.Kind)
			}
			if len(p.Levels) == 0 {
				t.Error("no levels")
			}
			if slices.Contains(p.Levels, parameterdefs.ParameterLevelAccount) && !slices.Contains(p.Levels, parameterdefs.ParameterLevelAccountExt) {
				t.Error("ParameterLevelAccount without ParameterLevelAccountExt")
			}
		})
	}
}

func TestParameterDefsForLevel(t *testing.T) {
	sqlNames := func(defs []parameterdefs.ParameterDef) []string {
		return collections.Map(defs, func(p parameterdefs.ParameterDef) string { return p.SqlName })
	}

	t.Run("task level", func(t *testing.T) {
		require.Equal(t, []string{
			"LOG_EVENT_LEVEL",
			"LOG_LEVEL",
			"SUSPEND_TASK_AFTER_NUM_FAILURES",
			"TASK_AUTO_RETRY_ATTEMPTS",
			"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE",
			"USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS",
			"USER_TASK_TIMEOUT_MS",
		}, sqlNames(ParameterDefsForLevel(parameterdefs.ParameterLevelTask)))
	})

	t.Run("session level", func(t *testing.T) {
		require.Equal(t, []string{
			"LOG_EVENT_LEVEL",
			"LOG_LEVEL",
			"QUOTED_IDENTIFIERS_IGNORE_CASE",
			"TRACE_LEVEL",
		}, sqlNames(ParameterDefsForLevel(parameterdefs.ParameterLevelSession)))
	})

	t.Run("account level skips parameters that are only settable on the extended set", func(t *testing.T) {
		require.NotContains(t, sqlNames(ParameterDefsForLevel(parameterdefs.ParameterLevelAccount)), "ENABLE_CONSOLE_OUTPUT")
		require.Contains(t, sqlNames(ParameterDefsForLevel(parameterdefs.ParameterLevelAccountExt)), "ENABLE_CONSOLE_OUTPUT")
	})

	t.Run("unknown level", func(t *testing.T) {
		require.Empty(t, ParameterDefsForLevel(parameterdefs.ParameterLevel("NOT_A_LEVEL")))
	})
}
