package defs

import (
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
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

// TestAllParameters_haveParsers asserts every catalog entry resolves to everything the generated read
// and write accessors are built from. Without this, adding a parameter with a new kind would fail
// during generation rather than here.
func TestAllParameters_haveParsers(t *testing.T) {
	for _, p := range AllParameters {
		t.Run(p.SqlName, func(t *testing.T) {
			info, err := g.InfoForKind(p.Kind)
			require.NoError(t, err)
			require.NotEmpty(t, p.SqlName)
			require.NotEmpty(t, g.ParameterSqlToFieldName(p))
			require.NotEmpty(t, info.GoType)
			require.NotEmpty(t, info.FieldType)
			require.NotEmpty(t, info.ReadParser)
			require.NotEmpty(t, info.WriteParser)
		})
	}
}

func TestParameterDefsForLevel(t *testing.T) {
	sqlNames := func(defs []parameterdefs.ParameterDef) []string {
		return collections.Map(defs, func(p parameterdefs.ParameterDef) string { return p.SqlName })
	}

	t.Run("task level", func(t *testing.T) {
		require.Equal(t, []string{
			"ABORT_DETACHED_QUERY",
			"AUTOCOMMIT",
			"BINARY_INPUT_FORMAT",
			"BINARY_OUTPUT_FORMAT",
			"CLIENT_MEMORY_LIMIT",
			"CLIENT_METADATA_REQUEST_USE_CONNECTION_CTX",
			"CLIENT_PREFETCH_THREADS",
			"CLIENT_RESULT_CHUNK_SIZE",
			"CLIENT_RESULT_COLUMN_CASE_INSENSITIVE",
			"CLIENT_SESSION_KEEP_ALIVE",
			"CLIENT_SESSION_KEEP_ALIVE_HEARTBEAT_FREQUENCY",
			"CLIENT_TIMESTAMP_TYPE_MAPPING",
			"DATE_INPUT_FORMAT",
			"DATE_OUTPUT_FORMAT",
			"ENABLE_UNLOAD_PHYSICAL_TYPE_OPTIMIZATION",
			"ERROR_ON_NONDETERMINISTIC_MERGE",
			"ERROR_ON_NONDETERMINISTIC_UPDATE",
			"GEOGRAPHY_OUTPUT_FORMAT",
			"GEOMETRY_OUTPUT_FORMAT",
			"JDBC_TREAT_TIMESTAMP_NTZ_AS_UTC",
			"JDBC_USE_SESSION_TIMEZONE",
			"JSON_INDENT",
			"LOCK_TIMEOUT",
			"LOG_EVENT_LEVEL",
			"LOG_LEVEL",
			"MULTI_STATEMENT_COUNT",
			"NOORDER_SEQUENCE_AS_DEFAULT",
			"ODBC_TREAT_DECIMAL_AS_INT",
			"QUERY_TAG",
			"QUOTED_IDENTIFIERS_IGNORE_CASE",
			"ROWS_PER_RESULTSET",
			"S3_STAGE_VPCE_DNS_NAME",
			"SERVERLESS_TASK_MAX_STATEMENT_SIZE",
			"SERVERLESS_TASK_MIN_STATEMENT_SIZE",
			"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS",
			"STATEMENT_TIMEOUT_IN_SECONDS",
			"STRICT_JSON_OUTPUT",
			"SUSPEND_TASK_AFTER_NUM_FAILURES",
			"TASK_AUTO_RETRY_ATTEMPTS",
			"TIME_INPUT_FORMAT",
			"TIME_OUTPUT_FORMAT",
			"TIMESTAMP_DAY_IS_ALWAYS_24H",
			"TIMESTAMP_INPUT_FORMAT",
			"TIMESTAMP_LTZ_OUTPUT_FORMAT",
			"TIMESTAMP_NTZ_OUTPUT_FORMAT",
			"TIMESTAMP_OUTPUT_FORMAT",
			"TIMESTAMP_TYPE_MAPPING",
			"TIMESTAMP_TZ_OUTPUT_FORMAT",
			"TIMEZONE",
			"TRACE_LEVEL",
			"TRANSACTION_ABORT_ON_ERROR",
			"TRANSACTION_DEFAULT_ISOLATION_LEVEL",
			"TWO_DIGIT_CENTURY_START",
			"UNSUPPORTED_DDL_ACTION",
			"USE_CACHED_RESULT",
			"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE",
			"USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS",
			"USER_TASK_TIMEOUT_MS",
			"WEEK_OF_YEAR_POLICY",
			"WEEK_START",
		}, sqlNames(ParameterDefsForLevel(parameterdefs.ParameterLevelTask)))
	})

	t.Run("warehouse level", func(t *testing.T) {
		require.Equal(t, []string{
			"MAX_CONCURRENCY_LEVEL",
			"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS",
			"STATEMENT_TIMEOUT_IN_SECONDS",
		}, sqlNames(ParameterDefsForLevel(parameterdefs.ParameterLevelWarehouse)))
	})

	t.Run("warehouse adaptive level", func(t *testing.T) {
		require.Equal(t, []string{
			"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS",
			"STATEMENT_TIMEOUT_IN_SECONDS",
		}, sqlNames(ParameterDefsForLevel(parameterdefs.ParameterLevelWarehouseAdaptive)))
	})

	t.Run("warehouse interactive level", func(t *testing.T) {
		require.Equal(t, []string{
			"FALLBACK_WAREHOUSE",
			"MAX_CONCURRENCY_LEVEL",
			"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS",
			"STATEMENT_TIMEOUT_IN_SECONDS",
		}, sqlNames(ParameterDefsForLevel(parameterdefs.ParameterLevelWarehouseInteractive)))
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

	t.Run("openflow deployment level", func(t *testing.T) {
		require.Equal(t, []string{
			"EVENT_TABLE",
		}, sqlNames(ParameterDefsForLevel(parameterdefs.ParameterLevelOpenflowDeployment)))
	})

	t.Run("unknown level", func(t *testing.T) {
		require.Empty(t, ParameterDefsForLevel(parameterdefs.ParameterLevel("NOT_A_LEVEL")))
	})
}
