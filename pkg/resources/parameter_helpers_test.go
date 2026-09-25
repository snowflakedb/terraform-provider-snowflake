package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

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
