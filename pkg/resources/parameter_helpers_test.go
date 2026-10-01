package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestHandleAccountParameterRead(t *testing.T) {
	state := make(map[string]any, len(accountParametersSchema))
	for key, s := range accountParametersSchema {
		switch s.Type {
		case schema.TypeInt:
			state[key] = 0
		case schema.TypeBool:
			state[key] = false
		default:
			state[key] = ""
		}
	}
	d := schema.TestResourceDataRaw(t, accountParametersSchema, state)

	diags := handleAccountParameterRead(d, &sdk.AccountParametersDetails{
		AllowBindValuesAccess:             sdk.TypedParameter[bool]{Value: true},
		EventTable:                        sdk.TypedParameter[sdk.SchemaObjectIdentifier]{Value: sdk.NewSchemaObjectIdentifier("database", "schema", "event_table")},
		IcebergVersionDefault:             sdk.TypedParameter[int]{Value: 2},
		InitialReplicationSizeLimitInTb:   sdk.TypedParameter[float64]{Value: 9.9},
		EnableGetDdlUseDataTypeAlias:      sdk.TypedParameter[bool]{Value: true},
		UseWorkspacesForSql:               sdk.TypedParameter[string]{Value: "unset"},
		DefaultStreamlitNotebookWarehouse: sdk.TypedParameter[sdk.AccountObjectIdentifier]{Value: sdk.NewAccountObjectIdentifier("warehouse")},
	})

	require.Empty(t, diags)
	require.Equal(t, true, d.Get("allow_bind_values_access"))
	require.Equal(t, `"database"."schema"."event_table"`, d.Get("event_table"))
	require.Equal(t, 2, d.Get("iceberg_version_default"))
	require.Equal(t, "9.9", d.Get("initial_replication_size_limit_in_tb"))
	require.Equal(t, true, d.Get("enable_get_ddl_use_data_type_alias"))
	require.Equal(t, "unset", d.Get("use_workspaces_for_sql"))
	require.Equal(t, `"warehouse"`, d.Get("default_streamlit_notebook_warehouse"))
}

func TestHandleDatabaseParameterRead(t *testing.T) {
	state := make(map[string]any, len(databaseParametersSchema))
	for key, s := range databaseParametersSchema {
		switch s.Type {
		case schema.TypeInt:
			state[key] = 0
		case schema.TypeBool:
			state[key] = false
		default:
			state[key] = ""
		}
	}
	d := schema.TestResourceDataRaw(t, databaseParametersSchema, state)

	diags := handleDatabaseParameterRead(d, &sdk.DatabaseParametersDetails{
		ExternalVolume:          sdk.TypedParameter[sdk.AccountObjectIdentifier]{Value: sdk.NewAccountObjectIdentifier("external_volume")},
		Catalog:                 sdk.TypedParameter[sdk.AccountObjectIdentifier]{},
		DataRetentionTimeInDays: sdk.TypedParameter[int]{Value: 7},
		EnableConsoleOutput:     sdk.TypedParameter[bool]{Value: true},
		LogLevel:                sdk.TypedParameter[sdk.LogLevel]{Value: sdk.LogLevelInfo},
	})

	require.Empty(t, diags)
	require.Equal(t, `"external_volume"`, d.Get("external_volume"))
	require.Equal(t, "", d.Get("catalog"))
	require.Equal(t, 7, d.Get("data_retention_time_in_days"))
	require.Equal(t, true, d.Get("enable_console_output"))
	require.Equal(t, string(sdk.LogLevelInfo), d.Get("log_level"))
}

func TestHandleSchemaParameterRead(t *testing.T) {
	state := make(map[string]any, len(schemaParametersSchema))
	for key, s := range schemaParametersSchema {
		switch s.Type {
		case schema.TypeInt:
			state[key] = 0
		case schema.TypeBool:
			state[key] = false
		default:
			state[key] = ""
		}
	}
	d := schema.TestResourceDataRaw(t, schemaParametersSchema, state)

	diags := handleSchemaParameterRead(d, &sdk.SchemaParametersDetails{
		ExternalVolume:          sdk.TypedParameter[sdk.AccountObjectIdentifier]{Value: sdk.NewAccountObjectIdentifier("external_volume")},
		Catalog:                 sdk.TypedParameter[sdk.AccountObjectIdentifier]{},
		DataRetentionTimeInDays: sdk.TypedParameter[int]{Value: 7},
		EnableConsoleOutput:     sdk.TypedParameter[bool]{Value: true},
		LogLevel:                sdk.TypedParameter[sdk.LogLevel]{Value: sdk.LogLevelInfo},
		PipeExecutionPaused:     sdk.TypedParameter[bool]{Value: true},
	})

	require.Empty(t, diags)
	require.Equal(t, `"external_volume"`, d.Get("external_volume"))
	require.Equal(t, "", d.Get("catalog"))
	require.Equal(t, 7, d.Get("data_retention_time_in_days"))
	require.Equal(t, true, d.Get("enable_console_output"))
	require.Equal(t, string(sdk.LogLevelInfo), d.Get("log_level"))
	require.Equal(t, true, d.Get("pipe_execution_paused"))
}

func TestParameterSchema_AllCatalogKindsAreHandled(t *testing.T) {
	for _, parameter := range defs.AllParameters {
		t.Run(parameter.SqlName, func(t *testing.T) {
			require.NotPanics(t, func() {
				parameterSchema(parameter)
			})
		})
	}
}

func TestParameterSchema_PanicsForUnhandledKind(t *testing.T) {
	require.PanicsWithValue(
		t,
		`parameter EXAMPLE has unhandled kind "UnhandledKind"; register it as a primitive, enum, or identifier kind`,
		func() {
			parameterSchema(parameterdefs.ParameterDef{SqlName: "EXAMPLE", Kind: "UnhandledKind"})
		},
	)
}

func TestParameterSchema_AccountSpecificValidation(t *testing.T) {
	testCases := []struct {
		name      string
		parameter parameterdefs.ParameterDef
		valid     any
		invalid   any
	}{
		{name: "integer at least zero", parameter: defs.JsonIndent, valid: 0, invalid: -1},
		{name: "iceberg version default at least zero", parameter: defs.IcebergVersionDefault, valid: 0, invalid: -1},
		{name: "integer at least one", parameter: defs.ClientMemoryLimit, valid: 1, invalid: 0},
		{name: "integer at least sixteen", parameter: defs.ClientResultChunkSize, valid: 16, invalid: 15},
		{name: "integer at least nine hundred", parameter: defs.ClientSessionKeepAliveHeartbeatFrequency, valid: 900, invalid: 899},
		{name: "integer at least nineteen hundred", parameter: defs.TwoDigitCenturyStart, valid: 1900, invalid: 1899},
		{name: "float string", parameter: defs.InitialReplicationSizeLimitInTb, valid: "9.9", invalid: "not-a-float"},
		{name: "account identifier stored as string", parameter: defs.CatalogSync, valid: "catalog_sync"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := parameterSchema(tt.parameter)
			require.NotNil(t, s.ValidateDiagFunc)
			require.Empty(t, s.ValidateDiagFunc(tt.valid, cty.IndexStringPath(tt.parameter.FieldName())))
			if tt.invalid != nil {
				require.NotEmpty(t, s.ValidateDiagFunc(tt.invalid, cty.IndexStringPath(tt.parameter.FieldName())))
			}
		})
	}
}
