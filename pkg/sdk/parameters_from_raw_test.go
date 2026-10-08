package sdk

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetParameterFromRaw(t *testing.T) {
	id := randomAccountObjectIdentifier()

	setSql := func(t *testing.T, request *DatabaseSetRequest) string {
		t.Helper()
		sql, err := structToSQL(NewAlterDatabaseRequest(id).WithSet(*request).toOpts())
		require.NoError(t, err)
		return sql
	}

	// Parsing a raw string on the write path must produce the same typed value the read path produces
	// for the same string, otherwise a written value would not match what is later read back.
	t.Run("agrees with the read path", func(t *testing.T) {
		tests := []struct {
			key      string
			raw      string
			fromSet  func(*DatabaseSetRequest) any
			fromRead func(*DatabaseParametersDetails) any
		}{
			{
				key:      "DATA_RETENTION_TIME_IN_DAYS",
				raw:      "10",
				fromSet:  func(r *DatabaseSetRequest) any { return *r.DataRetentionTimeInDays },
				fromRead: func(d *DatabaseParametersDetails) any { return d.DataRetentionTimeInDays.Value },
			},
			{
				key:      "ENABLE_CONSOLE_OUTPUT",
				raw:      "true",
				fromSet:  func(r *DatabaseSetRequest) any { return *r.EnableConsoleOutput },
				fromRead: func(d *DatabaseParametersDetails) any { return d.EnableConsoleOutput.Value },
			},
			{
				key:      "DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU",
				raw:      "MY_POOL",
				fromSet:  func(r *DatabaseSetRequest) any { return *r.DefaultNotebookComputePoolCpu },
				fromRead: func(d *DatabaseParametersDetails) any { return d.DefaultNotebookComputePoolCpu.Value },
			},
			{
				// The one kind whose sides differ: the write field wraps the string, the read side does not.
				key:      "DEFAULT_DDL_COLLATION",
				raw:      "en_US",
				fromSet:  func(r *DatabaseSetRequest) any { return r.DefaultDdlCollation.Value },
				fromRead: func(d *DatabaseParametersDetails) any { return d.DefaultDdlCollation.Value },
			},
			{
				key:      "LOG_LEVEL",
				raw:      "INFO",
				fromSet:  func(r *DatabaseSetRequest) any { return *r.LogLevel },
				fromRead: func(d *DatabaseParametersDetails) any { return d.LogLevel.Value },
			},
			{
				key:      "STORAGE_SERIALIZATION_POLICY",
				raw:      "COMPATIBLE",
				fromSet:  func(r *DatabaseSetRequest) any { return *r.StorageSerializationPolicy },
				fromRead: func(d *DatabaseParametersDetails) any { return d.StorageSerializationPolicy.Value },
			},
			{
				key:      "USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE",
				raw:      "XXXLARGE",
				fromSet:  func(r *DatabaseSetRequest) any { return *r.UserTaskManagedInitialWarehouseSize },
				fromRead: func(d *DatabaseParametersDetails) any { return d.UserTaskManagedInitialWarehouseSize.Value },
			},
			{
				key:      "CATALOG",
				raw:      "MY_CATALOG",
				fromSet:  func(r *DatabaseSetRequest) any { return *r.Catalog },
				fromRead: func(d *DatabaseParametersDetails) any { return d.Catalog.Value },
			},
		}

		for _, tt := range tests {
			t.Run(tt.key, func(t *testing.T) {
				request := NewDatabaseSetRequest()
				require.NoError(t, request.SetParameterFromRaw(tt.key, tt.raw))

				details, err := ToDatabaseParametersDetails([]*Parameter{{Key: tt.key, Value: tt.raw}})
				require.NoError(t, err)

				require.EqualValues(t, tt.fromRead(details), tt.fromSet(request))
			})
		}
	})

	// Going through the writer must render the same SQL as assigning the field, so quoting cannot drift
	// away from the ddl tags.
	t.Run("renders the same SQL as a direct assignment", func(t *testing.T) {
		tests := []struct {
			key      string
			raw      string
			assigned func(*DatabaseSetRequest)
		}{
			{"DATA_RETENTION_TIME_IN_DAYS", "10", func(r *DatabaseSetRequest) { r.DataRetentionTimeInDays = Int(10) }},
			{"DEFAULT_DDL_COLLATION", "en_US", func(r *DatabaseSetRequest) {
				r.DefaultDdlCollation = &StringAllowEmpty{Value: "en_US"}
			}},
			{"LOG_LEVEL", "INFO", func(r *DatabaseSetRequest) { r.LogLevel = Pointer(LogLevelInfo) }},
			// Identifiers render through a different ddl tag than value parameters.
			{"CATALOG", "MY_CATALOG", func(r *DatabaseSetRequest) {
				r.Catalog = Pointer(NewAccountObjectIdentifier("MY_CATALOG"))
			}},
		}

		for _, tt := range tests {
			t.Run(tt.key, func(t *testing.T) {
				fromRaw := NewDatabaseSetRequest()
				require.NoError(t, fromRaw.SetParameterFromRaw(tt.key, tt.raw))

				direct := NewDatabaseSetRequest()
				tt.assigned(direct)

				require.Equal(t, setSql(t, direct), setSql(t, fromRaw))
			})
		}
	})

	// A caller falls back to the untyped setter on the sentinel, so a bad value must never produce it.
	t.Run("distinguishes an unsupported parameter from a bad value", func(t *testing.T) {
		tests := []struct {
			name         string
			key          string
			value        string
			wantSentinel bool
		}{
			{name: "unknown key", key: "NOT_A_PARAMETER", value: "whatever", wantSentinel: true},
			// Comment sits on the same struct but is not a parameter.
			{name: "non-parameter field", key: "COMMENT", value: "x", wantSentinel: true},
			{name: "unparsable int", key: "DATA_RETENTION_TIME_IN_DAYS", value: "abc"},
			{name: "invalid enum", key: "LOG_LEVEL", value: "NOT_A_LEVEL"},
			{name: "invalid identifier", key: "CATALOG", value: `"a"."b"`},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := NewDatabaseSetRequest().SetParameterFromRaw(tt.key, tt.value)
				require.Error(t, err)
				if tt.wantSentinel {
					require.ErrorIs(t, err, ErrParameterNotSupported)
				} else {
					require.NotErrorIs(t, err, ErrParameterNotSupported)
				}
			})
		}
	})

	t.Run("is generated for every object with catalog parameters", func(t *testing.T) {
		require.NoError(t, NewSchemaSetRequest().SetParameterFromRaw("LOG_LEVEL", "INFO"))
		require.ErrorIs(t, NewSchemaSetRequest().SetParameterFromRaw("COMMENT", "x"), ErrParameterNotSupported)
	})

	// INITIAL_REPLICATION_SIZE_LIMIT_IN_TB is a decimal number whose SQL grammar requires an
	// unquoted literal, so it is a KindFloat parameter rather than KindString.
	t.Run("float parameters render without quotes through the raw path", func(t *testing.T) {
		setAccountSql := func(t *testing.T, request *AccountParametersRequest) string {
			t.Helper()
			sql, err := structToSQL(NewAlterAccountRequest().WithSet(*NewAccountSetRequest().WithParameters(*request)).toOpts())
			require.NoError(t, err)
			return sql
		}

		fromRaw := NewAccountParametersRequest()
		require.NoError(t, fromRaw.SetParameterFromRaw("INITIAL_REPLICATION_SIZE_LIMIT_IN_TB", "9.9"))

		direct := NewAccountParametersRequest().WithInitialReplicationSizeLimitInTb(9.9)

		sql := setAccountSql(t, fromRaw)
		require.Equal(t, setAccountSql(t, direct), sql)
		require.Contains(t, sql, "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB = 9.9")
		require.NotContains(t, sql, "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB = '9.9'")
	})

	t.Run("accepts a lowercase parameter name", func(t *testing.T) {
		request := NewDatabaseSetRequest()
		require.NoError(t, request.SetParameterFromRaw("data_retention_time_in_days", "10"))
		require.NotNil(t, request.DataRetentionTimeInDays)
		require.Equal(t, 10, *request.DataRetentionTimeInDays)
	})

	t.Run("accepts an empty string for a StringAllowEmpty parameter", func(t *testing.T) {
		request := NewAccountParametersRequest()
		require.NoError(t, request.SetParameterFromRaw("DEFAULT_DDL_COLLATION", ""))
		require.NotNil(t, request.DefaultDdlCollation)
		require.Equal(t, "", request.DefaultDdlCollation.Value)
	})
}

func TestUnsetParameterFromRaw(t *testing.T) {
	id := randomAccountObjectIdentifier()

	unsetSql := func(t *testing.T, request *DatabaseUnsetRequest) string {
		t.Helper()
		sql, err := structToSQL(NewAlterDatabaseRequest(id).WithUnset(*request).toOpts())
		require.NoError(t, err)
		return sql
	}

	t.Run("renders the same SQL as a direct assignment", func(t *testing.T) {
		fromRaw := NewDatabaseUnsetRequest()
		require.NoError(t, fromRaw.UnsetParameterFromRaw("DATA_RETENTION_TIME_IN_DAYS"))

		direct := NewDatabaseUnsetRequest()
		direct.DataRetentionTimeInDays = Bool(true)

		require.Equal(t, unsetSql(t, direct), unsetSql(t, fromRaw))
	})

	t.Run("rejects parameters it cannot unset", func(t *testing.T) {
		for _, key := range []string{"NOT_A_PARAMETER", "COMMENT"} {
			t.Run(key, func(t *testing.T) {
				require.ErrorIs(t, NewDatabaseUnsetRequest().UnsetParameterFromRaw(key), ErrParameterNotSupported)
				require.ErrorIs(t, NewSchemaUnsetRequest().UnsetParameterFromRaw(key), ErrParameterNotSupported)
			})
		}
	})
}

func TestAssignParsedParameter(t *testing.T) {
	t.Run("assigns the parsed value", func(t *testing.T) {
		var target *LogLevel
		require.NoError(t, assignParsedParameter("INFO", ToLogLevel, &target))
		require.NotNil(t, target)
		require.Equal(t, LogLevelInfo, *target)
	})

	t.Run("propagates the parse error and leaves the target unset", func(t *testing.T) {
		parseErr := errors.New("boom")
		var target *int
		err := assignParsedParameter("abc", func(string) (int, error) { return 0, parseErr }, &target)
		require.ErrorIs(t, err, parseErr)
		require.Nil(t, target)
	})
}
