package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func init() {
	accountObjectId := randomAccountObjectIdentifier()
	databaseObjectId := randomDatabaseObjectIdentifier()

	parametersTests.SetAccountParameterRaw.
		withDefaultOpts(func() *SetAccountParameterRawParameterOptions {
			return &SetAccountParameterRawParameterOptions{Parameter: "SOME_FUTURE_PARAMETER", Value: "'a b'"}
		}).
		withExpectedSql(case_Parameters_sql_SetAccountParameterRaw_basic, `ALTER ACCOUNT SET SOME_FUTURE_PARAMETER = 'a b'`)

	parametersTests.UnsetAccountParameterRaw.
		withDefaultOpts(func() *UnsetAccountParameterRawParameterOptions {
			return &UnsetAccountParameterRawParameterOptions{Parameter: "SOME_FUTURE_PARAMETER"}
		}).
		withExpectedSql(case_Parameters_sql_UnsetAccountParameterRaw_basic, `ALTER ACCOUNT UNSET SOME_FUTURE_PARAMETER`)

	parametersTests.SetObjectParameterOnObjectRaw.
		withDefaultOpts(func() *SetObjectParameterOnObjectRawParameterOptions {
			return &SetObjectParameterOnObjectRawParameterOptions{
				ObjectType: ObjectTypeUser,
				ObjectName: accountObjectId,
				Parameter:  "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR",
				Value:      "TRUE",
			}
		}).
		withExpectedSqlf(case_Parameters_sql_SetObjectParameterOnObjectRaw_basic,
			`ALTER USER %s SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = TRUE`, accountObjectId.FullyQualifiedName(),
		)

	parametersTests.UnsetObjectParameterOnObjectRaw.
		withDefaultOpts(func() *UnsetObjectParameterOnObjectRawParameterOptions {
			return &UnsetObjectParameterOnObjectRawParameterOptions{
				ObjectType: ObjectTypeUser,
				ObjectName: accountObjectId,
				Parameter:  "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR",
			}
		}).
		withExpectedSqlf(case_Parameters_sql_UnsetObjectParameterOnObjectRaw_basic,
			`ALTER USER %s UNSET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR`, accountObjectId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef("sql_UnsetObjectParameterOnObjectRaw_database",
			func(opts *UnsetObjectParameterOnObjectRawParameterOptions) {
				opts.ObjectType = ObjectTypeDatabase
				opts.Parameter = "USER_TASK_TIMEOUT_MS"
			},
			`ALTER DATABASE %s UNSET USER_TASK_TIMEOUT_MS`, accountObjectId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef("sql_UnsetObjectParameterOnObjectRaw_schema",
			func(opts *UnsetObjectParameterOnObjectRawParameterOptions) {
				opts.ObjectType = ObjectTypeSchema
				opts.ObjectName = databaseObjectId
				opts.Parameter = "DATA_RETENTION_TIME_IN_DAYS"
			},
			`ALTER SCHEMA %s UNSET DATA_RETENTION_TIME_IN_DAYS`, databaseObjectId.FullyQualifiedName(),
		)
}

func TestParameters_ShowParameters(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		assertOptsValidAndSqlEquals(t, &ShowParametersOptions{}, "SHOW PARAMETERS")
	})

	t.Run("like", func(t *testing.T) {
		opts := &ShowParametersOptions{
			Like: &Like{Pattern: String("test")},
		}
		assertOptsValidAndSqlEquals(t, opts, "SHOW PARAMETERS LIKE 'test'")
	})

	t.Run("in account", func(t *testing.T) {
		opts := &ShowParametersOptions{
			In: &ParametersIn{Account: Bool(true)},
		}
		assertOptsValidAndSqlEquals(t, opts, "SHOW PARAMETERS IN ACCOUNT")
	})

	t.Run("like and in account", func(t *testing.T) {
		opts := &ShowParametersOptions{
			Like: &Like{Pattern: String("test")},
			In:   &ParametersIn{Account: Bool(true)},
		}
		assertOptsValidAndSqlEquals(t, opts, "SHOW PARAMETERS LIKE 'test' IN ACCOUNT")
	})

	t.Run("in table", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &ShowParametersOptions{
			In: &ParametersIn{Table: id},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW PARAMETERS IN TABLE %s", id.FullyQualifiedName())
	})

	t.Run("in pipe", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &ShowParametersOptions{
			In: &ParametersIn{Pipe: id},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW PARAMETERS IN PIPE %s", id.FullyQualifiedName())
	})

	t.Run("like log event level in pipe", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &ShowParametersOptions{
			Like: &Like{Pattern: String("LOG_EVENT_LEVEL")},
			In:   &ParametersIn{Pipe: id},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW PARAMETERS LIKE 'LOG_EVENT_LEVEL' IN PIPE %s", id.FullyQualifiedName())
	})

	t.Run("no scope", func(t *testing.T) {
		opts := &ShowParametersOptions{
			In: &ParametersIn{},
		}
		assertOptsInvalid(t, opts, errAtLeastOneOf("Session", "Account", "User", "Warehouse", "Database", "Schema", "Task", "Table", "Function", "Procedure", "OpenflowDeployment", "Service", "Pipe", "Object"))
	})
}

// A parameter that the catalog starts to support later must keep sending the same SQL for the same configuration.
func TestParameterValueKind_Render_MatchesCatalogWritePath(t *testing.T) {
	cases := []struct {
		parameter string
		kind      ParameterValueKind
		value     string
	}{
		{"DEFAULT_DDL_COLLATION", ParameterValueKindString, "en_US"},
		{"DEFAULT_DDL_COLLATION", ParameterValueKindString, "it's"},
		{"QUERY_TAG", ParameterValueKindString, "two words"},
		{"LOG_LEVEL", ParameterValueKindString, "INFO"},
		{"TRACE_LEVEL", ParameterValueKindString, "ALWAYS"},
		{"CATALOG", ParameterValueKindIdentifier, "my_catalog"},
		{"NETWORK_POLICY", ParameterValueKindIdentifier, "my_policy"},
		{"EVENT_TABLE", ParameterValueKindIdentifier, "DB.SCHEMA.EVENTS"},
		{"DATA_RETENTION_TIME_IN_DAYS", ParameterValueKindInteger, "042"},
		{"WEEK_START", ParameterValueKindInteger, "1"},
		{"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB", ParameterValueKindFloat, "9.9"},
		{"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB", ParameterValueKindFloat, "3"},
		{"INITIAL_REPLICATION_SIZE_LIMIT_IN_TB", ParameterValueKindFloat, "1e3"},
		{"ENABLE_PERSONAL_DATABASE", ParameterValueKindBoolean, "TRUE"},
		{"ENABLE_PERSONAL_DATABASE", ParameterValueKindBoolean, "false"},
	}
	for _, tc := range cases {
		t.Run(tc.parameter+"="+tc.value, func(t *testing.T) {
			req := NewAccountParametersRequest()
			require.NoError(t, req.SetParameterFromRaw(tc.parameter, tc.value))
			catalogSql, err := structToSQL(NewAlterAccountRequest().WithSet(*NewAccountSetRequest().WithParameters(*req)).toOpts())
			require.NoError(t, err)

			rendered, err := tc.kind.Render(tc.parameter, tc.value)
			require.NoError(t, err)
			rawSql, err := structToSQL(&SetAccountParameterRawParameterOptions{Parameter: tc.parameter, Value: rendered})
			require.NoError(t, err)

			require.Equal(t, catalogSql, rawSql)
		})
	}
}

// Every account parameter of the catalog, with a sample value of its type, is rendered by the matching kind the same way the catalog renders it.
func TestParameterValueKind_Render_MatchesCatalogWritePathForEveryAccountParameter(t *testing.T) {
	samples := map[string]struct {
		kind  ParameterValueKind
		value string
	}{
		"bool":                             {ParameterValueKindBoolean, "TRUE"},
		"int":                              {ParameterValueKindInteger, "042"},
		"float64":                          {ParameterValueKindFloat, "3"},
		"string":                           {ParameterValueKindString, "it's a b"},
		"StringAllowEmpty":                 {ParameterValueKindString, "it's a b"},
		"AccountObjectIdentifier":          {ParameterValueKindIdentifier, "my_object"},
		"SchemaObjectIdentifier":           {ParameterValueKindIdentifier, "db.sc.my_object"},
		"ActivePythonProfiler":             {ParameterValueKindString, "LINE"},
		"BinaryInputFormat":                {ParameterValueKindString, "HEX"},
		"BinaryOutputFormat":               {ParameterValueKindString, "HEX"},
		"ClientTimestampTypeMapping":       {ParameterValueKindString, "TIMESTAMP_LTZ"},
		"DefaultNullOrdering":              {ParameterValueKindString, "FIRST"},
		"GeographyOutputFormat":            {ParameterValueKindString, "WKT"},
		"GeometryOutputFormat":             {ParameterValueKindString, "WKT"},
		"LogLevel":                         {ParameterValueKindString, "INFO"},
		"MetricLevel":                      {ParameterValueKindString, "ALL"},
		"StorageSerializationPolicy":       {ParameterValueKindString, "COMPATIBLE"},
		"TimestampTypeMapping":             {ParameterValueKindString, "TIMESTAMP_LTZ"},
		"TraceLevel":                       {ParameterValueKindString, "ALWAYS"},
		"TransactionDefaultIsolationLevel": {ParameterValueKindString, "READ COMMITTED"},
		"UnsupportedDDLAction":             {ParameterValueKindString, "IGNORE"},
		"WarehouseSize":                    {ParameterValueKindString, "XSMALL"},
	}
	for _, p := range defs.ParameterDefsForLevel(ssot.ParameterLevelAccountExt) {
		t.Run(p.SqlName, func(t *testing.T) {
			sample, ok := samples[p.Kind]
			require.Truef(t, ok, "add a sample value for the catalog kind %s", p.Kind)

			req := NewAccountParametersRequest()
			require.NoError(t, req.SetParameterFromRaw(p.SqlName, sample.value))
			catalogSql, err := structToSQL(NewAlterAccountRequest().WithSet(*NewAccountSetRequest().WithParameters(*req)).toOpts())
			require.NoError(t, err)

			rendered, err := sample.kind.Render(p.SqlName, sample.value)
			require.NoError(t, err)
			rawSql, err := structToSQL(&SetAccountParameterRawParameterOptions{Parameter: p.SqlName, Value: rendered})
			require.NoError(t, err)

			require.Equal(t, catalogSql, rawSql)
		})
	}
}

func TestParameterValueKind_Render(t *testing.T) {
	cases := []struct {
		name    string
		kind    ParameterValueKind
		value   string
		want    string
		wantErr string
	}{
		{name: "string", kind: ParameterValueKindString, value: "a b", want: "'a b'"},
		{name: "string with quote is escaped", kind: ParameterValueKindString, value: "it's", want: `'it\'s'`},
		{name: "identifier", kind: ParameterValueKindIdentifier, value: "my_id", want: `"my_id"`},
		{name: "quoted identifier", kind: ParameterValueKindIdentifier, value: `"My Policy"`, want: `"My Policy"`},
		{name: "fully qualified identifier", kind: ParameterValueKindIdentifier, value: "db.sc.obj", want: `"db"."sc"."obj"`},
		{name: "integer", kind: ParameterValueKindInteger, value: "-5", want: "-5"},
		{name: "float", kind: ParameterValueKindFloat, value: "9.9", want: "9.9"},
		{name: "whole float", kind: ParameterValueKindFloat, value: "3", want: "3.0"},
		{name: "boolean", kind: ParameterValueKindBoolean, value: "TRUE", want: "true"},
		{name: "invalid integer", kind: ParameterValueKindInteger, value: "abc", wantErr: "integer value expected for P parameter, but got abc instead"},
		{name: "invalid float", kind: ParameterValueKindFloat, value: "abc", wantErr: "float value expected for P parameter, but got abc instead"},
		{name: "invalid boolean", kind: ParameterValueKindBoolean, value: "maybe", wantErr: `boolean value ("true"/"false") expected for P parameter, but got maybe instead`},
		{name: "empty identifier", kind: ParameterValueKindIdentifier, value: "", wantErr: "invalid identifier value for P parameter: incompatible identifier: "},
		{name: "unknown kind", kind: ParameterValueKind("OTHER"), value: "x", wantErr: `unknown parameter value kind "OTHER"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.kind.Render("P", tc.value)
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestShowParametersInGenericObject(t *testing.T) {
	opts := &ShowParametersOptions{
		Like: &Like{Pattern: String("SOME_PARAMETER")},
		In: &ParametersIn{
			Object: &ParametersInObject{ObjectType: ObjectType("COMPUTE POOL"), Name: NewAccountObjectIdentifier("my_pool")},
		},
	}
	assertOptsValidAndSqlEquals(t, opts, `SHOW PARAMETERS LIKE 'SOME_PARAMETER' IN COMPUTE POOL "my_pool"`)

	opts = &ShowParametersOptions{In: &ParametersIn{Object: &ParametersInObject{ObjectType: ObjectTypePipe, Name: NewSchemaObjectIdentifier("db", "sc", "p")}}}
	assertOptsValidAndSqlEquals(t, opts, `SHOW PARAMETERS IN PIPE "db"."sc"."p"`)
}

func TestExactParameter(t *testing.T) {
	rows := []*Parameter{{Key: "SOMEXPARAMETER"}, {Key: "SOME_PARAMETER", Value: "1"}}

	p, err := exactParameter(rows, "some_parameter")
	require.NoError(t, err)
	require.Equal(t, "1", p.Value)

	_, err = exactParameter(rows, "OTHER")
	require.EqualError(t, err, "parameter OTHER not found")
}

// No database connection: every case returns before any SQL is sent, and the error shows which write path was taken.
func TestParameters_SetAccountParameterWithGenericFallback(t *testing.T) {
	p := &parameters{client: &Client{}}

	t.Run("supported parameter: the catalog parses the value, and the kind is not used", func(t *testing.T) {
		// STRING would accept "abc", so the integer parse error comes from the catalog.
		err := p.SetAccountParameterWithGenericFallback(t.Context(), "DATA_RETENTION_TIME_IN_DAYS", ParameterValueKindString, "abc")
		require.ErrorContains(t, err, "invalid syntax")
	})

	t.Run("unsupported parameter: the value is rendered by the kind", func(t *testing.T) {
		err := p.SetAccountParameterWithGenericFallback(t.Context(), "SOME_FUTURE_PARAMETER", ParameterValueKindInteger, "abc")
		require.ErrorContains(t, err, "integer value expected for SOME_FUTURE_PARAMETER parameter")
	})

	t.Run("enum value that the catalog does not recognize: the value is rendered by the kind", func(t *testing.T) {
		err := p.SetAccountParameterWithGenericFallback(t.Context(), "LOG_LEVEL", ParameterValueKindIdentifier, "")
		require.ErrorContains(t, err, "invalid identifier value for LOG_LEVEL parameter")
	})
}

// No database connection: every case returns before any SQL is sent, and the error shows which write path was taken.
func TestParameters_SetObjectParameterOnObjectWithGenericFallback(t *testing.T) {
	p := &parameters{client: &Client{}}
	database := Object{ObjectType: ObjectTypeDatabase, Name: randomAccountObjectIdentifier()}
	table := Object{ObjectType: ObjectTypeTable, Name: randomSchemaObjectIdentifier()}
	unknownType := Object{ObjectType: ObjectType("NewObjectType"), Name: randomSchemaObjectIdentifier()}

	t.Run("supported object type and parameter: the catalog parses the value, and the kind is not used", func(t *testing.T) {
		// STRING would accept "abc", so the integer parse error comes from the catalog.
		err := p.SetObjectParameterOnObjectWithGenericFallback(t.Context(), database, "DATA_RETENTION_TIME_IN_DAYS", ParameterValueKindString, "abc")
		require.ErrorContains(t, err, "invalid syntax")
	})

	t.Run("supported object type and unsupported parameter: the value is rendered by the kind", func(t *testing.T) {
		err := p.SetObjectParameterOnObjectWithGenericFallback(t.Context(), database, "SOME_FUTURE_PARAMETER", ParameterValueKindInteger, "abc")
		require.ErrorContains(t, err, "integer value expected for SOME_FUTURE_PARAMETER parameter")
	})

	t.Run("unsupported object type: the value is rendered by the kind", func(t *testing.T) {
		err := p.SetObjectParameterOnObjectWithGenericFallback(t.Context(), table, "DATA_RETENTION_TIME_IN_DAYS", ParameterValueKindBoolean, "maybe")
		require.ErrorContains(t, err, `boolean value ("true"/"false") expected for DATA_RETENTION_TIME_IN_DAYS parameter`)
	})

	t.Run("unknown object type: the value is rendered by the kind", func(t *testing.T) {
		err := p.SetObjectParameterOnObjectWithGenericFallback(t.Context(), unknownType, "SOME_FUTURE_PARAMETER", ParameterValueKindInteger, "abc")
		require.ErrorContains(t, err, "integer value expected for SOME_FUTURE_PARAMETER parameter")
	})
}

func TestSupportedParameters(t *testing.T) {
	t.Run("account", func(t *testing.T) {
		supported := SupportedAccountParameters()
		require.Contains(t, supported, AccountParameterDataRetentionTimeInDays)
		// It is in the generated list, but the catalog write path does not support it.
		require.Contains(t, AllAccountParameters, AccountParameterServiceCallerTokenValiditySecs)
		require.NotContains(t, supported, AccountParameterServiceCallerTokenValiditySecs)
	})

	t.Run("objects", func(t *testing.T) {
		require.Contains(t, SupportedObjectParameters(ObjectTypeDatabase), ObjectParameterDataRetentionTimeInDays)
		require.NotContains(t, SupportedObjectParameters(ObjectTypeDatabase), ObjectParameter("METRIC_LEVEL"))
		require.Empty(t, SupportedObjectParameters(ObjectTypeTable))
		require.Empty(t, SupportedObjectParameters(ObjectType("NewObjectType")))
	})

	t.Run("every object type with a catalog writer is asked", func(t *testing.T) {
		for objectType := range objectParameterCatalogSetters {
			require.Contains(t, objectParameterCatalogSupport, objectType)
		}
		for objectType := range objectParameterCatalogSupport {
			require.Contains(t, objectParameterCatalogSetters, objectType)
		}
	})
}
