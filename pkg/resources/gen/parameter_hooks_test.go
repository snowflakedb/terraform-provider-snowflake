package gen

import (
	"bytes"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func TestParameterHooks_schemaTemplateMatchesHandWrittenMappings(t *testing.T) {
	model := ResourceModel{
		Name:       "Schema",
		NamePlural: "Schemas",
		Hooks:      ssot.DefaultResourceExtHooks(ssotdefs.Schema.ResourceHooks...),
	}
	var rendered bytes.Buffer
	require.NoError(t, ParametersTemplate.Execute(&rendered, model))
	got := rendered.String()

	require.Contains(t, got, "handleParameterCreateWithMapping(d, ssotdefs.Catalog.FieldName(), &req.Catalog, stringToAccountObjectIdentifier),")
	require.Contains(t, got, "handleParameterCreate(d, ssotdefs.DataRetentionTimeInDays.FieldName(), &req.DataRetentionTimeInDays),")
	require.Contains(t, got, "handleParameterCreateWithMapping(d, ssotdefs.DefaultDdlCollation.FieldName(), &req.DefaultDdlCollation, stringToStringAllowEmpty),")
	require.Contains(t, got, "sdk.ToLogLevel),")
	require.Contains(t, got, "sdk.ToWarehouseSize),")
	require.Contains(t, got, "handleParameterUpdateWithMapping(d, ssotdefs.Catalog.FieldName(), &set.Catalog, &unset.Catalog, stringToAccountObjectIdentifier),")
	require.Contains(t, got, "setResourceData(d, ssotdefs.Catalog.FieldName(), parameters.Catalog.Value.FullyQualifiedName()),")
	require.Contains(t, got, "setResourceData(d, ssotdefs.DataRetentionTimeInDays.FieldName(), parameters.DataRetentionTimeInDays.Value),")
	require.Contains(t, got, "IdentifierTypedParameterValueComputedIf(ssotdefs.Catalog.FieldName(), parameters.Catalog, sdk.ParameterTypeSchema),")
	require.Contains(t, got, "IntTypedParameterValueComputedIf(ssotdefs.DataRetentionTimeInDays.FieldName(), parameters.DataRetentionTimeInDays, sdk.ParameterTypeSchema),")
	require.Contains(t, got, "BoolTypedParameterValueComputedIf(ssotdefs.EnableConsoleOutput.FieldName(), parameters.EnableConsoleOutput, sdk.ParameterTypeSchema),")
	require.Contains(t, got, "StringAllowEmptyTypedParameterValueComputedIf(ssotdefs.DefaultDdlCollation.FieldName(), parameters.DefaultDdlCollation, sdk.ParameterTypeSchema),")
	require.Contains(t, got, "schemaParseIdExt(d.Id())")
}

func TestParameterHooks_customDiffUsesConfiguredParseId(t *testing.T) {
	model := ResourceModel{
		Name:       "Schema",
		NamePlural: "Schemas",
		Hooks: ssot.DefaultResourceExtHooks(
			ssot.WithParametersCustomDiff(ssot.Generated()),
			ssot.WithParseId(ssot.Generated()),
		),
	}
	var rendered bytes.Buffer
	require.NoError(t, ParametersTemplate.Execute(&rendered, model))

	require.Contains(t, rendered.String(), "schemaParseId(d.Id())")
	require.NotContains(t, rendered.String(), "schemaParseIdExt(d.Id())")
}

func TestParameterHooks_parametersOutputSetRetainsExtCall(t *testing.T) {
	model := ResourceModel{
		Name:  "Widget",
		Hooks: ssot.DefaultResourceExtHooks(ssot.WithParametersOutputSet(ssot.ExtCall())),
	}
	var rendered bytes.Buffer
	require.NoError(t, LifecycleTemplate.Execute(&rendered, model))

	require.Contains(t, rendered.String(), "widgetParametersOutputSetExt(d, rawWidgetParameters, meta),")
	require.NotContains(t, rendered.String(), "schemas.WidgetParametersToSchema")
}

func TestParameterHookModel_fragments(t *testing.T) {
	tests := []struct {
		sqlName      string
		kind         string
		needsMapping bool
		mapper       string
		read         string
		diffFunc     string
	}{
		{
			sqlName:  "ENABLE_CONSOLE_OUTPUT",
			kind:     string(ssot.AttrBool),
			read:     "parameters.EnableConsoleOutput.Value",
			diffFunc: "BoolTypedParameterValueComputedIf",
		},
		{
			sqlName:  "DATA_RETENTION_TIME_IN_DAYS",
			kind:     string(ssot.AttrInt),
			read:     "parameters.DataRetentionTimeInDays.Value",
			diffFunc: "IntTypedParameterValueComputedIf",
		},
		{
			sqlName:  "DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU",
			kind:     string(ssot.AttrString),
			read:     "parameters.DefaultNotebookComputePoolCpu.Value",
			diffFunc: "StringTypedParameterValueComputedIf",
		},
		{
			sqlName:      "DEFAULT_DDL_COLLATION",
			kind:         "StringAllowEmpty",
			needsMapping: true,
			mapper:       "stringToStringAllowEmpty",
			read:         "parameters.DefaultDdlCollation.Value",
			diffFunc:     "StringAllowEmptyTypedParameterValueComputedIf",
		},
		{
			sqlName:      "INITIAL_REPLICATION_SIZE_LIMIT_IN_TB",
			kind:         string(ssot.AttrFloat),
			needsMapping: true,
			mapper:       "sdk.ToFloat64",
			read:         "formatFloatParameter(parameters.InitialReplicationSizeLimitInTb.Value)",
			diffFunc:     "FloatTypedParameterValueComputedIf",
		},
		{
			sqlName:      "CATALOG",
			kind:         "AccountObjectIdentifier",
			needsMapping: true,
			mapper:       "stringToAccountObjectIdentifier",
			read:         "parameters.Catalog.Value.FullyQualifiedName()",
			diffFunc:     "IdentifierTypedParameterValueComputedIf",
		},
		{
			sqlName:      "EVENT_TABLE",
			kind:         "SchemaObjectIdentifier",
			needsMapping: true,
			mapper:       "sdk.ParseSchemaObjectIdentifier",
			read:         "parameters.EventTable.Value.FullyQualifiedName()",
			diffFunc:     "IdentifierTypedParameterValueComputedIf",
		},
		{
			sqlName:      "LOG_LEVEL",
			kind:         "LogLevel",
			needsMapping: true,
			mapper:       "sdk.ToLogLevel",
			read:         "parameters.LogLevel.Value",
			diffFunc:     "StringTypedParameterValueComputedIf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sqlName, func(t *testing.T) {
			hook := newParameterHookModel(ssot.ParameterDef{SqlName: tt.sqlName, Kind: tt.kind})
			require.Equal(t, tt.needsMapping, hook.NeedsMapping())
			require.Equal(t, tt.mapper, hook.Mapper())
			require.Equal(t, tt.read, hook.Read())
			require.Equal(t, tt.diffFunc, hook.DiffFunc())
		})
	}
}
