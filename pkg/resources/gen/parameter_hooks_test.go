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

	require.Contains(t, got, `"catalog": parameterSchemaFromDefinition(parameterSchemaDefinition{`)
	require.Contains(t, got, `sqlName:     "CATALOG",`)
	require.Contains(t, got, `kind:        "AccountObjectIdentifier",`)
	require.Contains(t, got, `"catalog",`)
	require.Contains(t, got, `handleParameterCreateWithMapping(d, "catalog", &req.Catalog, stringToAccountObjectIdentifier),`)
	require.Contains(t, got, `handleParameterCreate(d, "data_retention_time_in_days", &req.DataRetentionTimeInDays),`)
	require.Contains(t, got, `handleParameterCreateWithMapping(d, "default_ddl_collation", &req.DefaultDdlCollation, stringToStringAllowEmpty),`)
	require.Contains(t, got, "sdk.ToLogLevel),")
	require.Contains(t, got, "sdk.ToWarehouseSize),")
	require.Contains(t, got, `handleParameterUpdateWithMapping(d, "catalog", &set.Catalog, &unset.Catalog, stringToAccountObjectIdentifier),`)
	require.Contains(t, got, `setResourceData(d, "catalog", parameters.Catalog.Value.FullyQualifiedName()),`)
	require.Contains(t, got, `setResourceData(d, "data_retention_time_in_days", parameters.DataRetentionTimeInDays.Value),`)
	require.Contains(t, got, `IdentifierTypedParameterValueComputedIf("catalog", parameters.Catalog, sdk.ParameterTypeSchema),`)
	require.Contains(t, got, `IntTypedParameterValueComputedIf("data_retention_time_in_days", parameters.DataRetentionTimeInDays, sdk.ParameterTypeSchema),`)
	require.Contains(t, got, `BoolTypedParameterValueComputedIf("enable_console_output", parameters.EnableConsoleOutput, sdk.ParameterTypeSchema),`)
	require.Contains(t, got, `StringAllowEmptyTypedParameterValueComputedIf("default_ddl_collation", parameters.DefaultDdlCollation, sdk.ParameterTypeSchema),`)
	require.Contains(t, got, "schemaParseIdExt(d.Id())")
	require.NotContains(t, got, "ssot.")
	require.NotContains(t, got, "ssotdefs.")
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

func TestParameterHooks_variantUsesBaseSdkNames(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: ssotdefs.WarehouseAdaptive}, nil, nil)

	var parametersRendered bytes.Buffer
	require.NoError(t, ParametersTemplate.Execute(&parametersRendered, model))
	parameters := parametersRendered.String()
	require.Contains(t, parameters, "*sdk.WarehouseParametersDetails")
	require.Contains(t, parameters, "sdk.ParameterTypeWarehouse")
	require.Contains(t, parameters, "*sdk.CreateAdaptiveWarehouseRequest")
	require.Contains(t, parameters, "*sdk.WarehouseAdaptiveSetRequest")
	require.Contains(t, parameters, "*sdk.WarehouseAdaptiveUnsetRequest")
	require.NotContains(t, parameters, "*sdk.WarehouseSetRequest")
	require.NotContains(t, parameters, "*sdk.WarehouseUnsetRequest")
	require.NotContains(t, parameters, "sdk.WarehouseAdaptiveParametersDetails")
	require.NotContains(t, parameters, "sdk.CreateWarehouseAdaptiveRequest")
	// Sharing SDK types must not expose parameters from another variant's catalog scope.
	require.NotContains(t, parameters, "fallback_warehouse")
	require.NotContains(t, parameters, "FallbackWarehouse")

	var lifecycleRendered bytes.Buffer
	require.NoError(t, LifecycleTemplate.Execute(&lifecycleRendered, model))
	require.Contains(t, lifecycleRendered.String(), "sdk.ToWarehouseParametersDetails(rawWarehouseAdaptiveParameters)")
}

func TestParameterHooks_nonVariantUsesOwnSdkNames(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: ssotdefs.Warehouse}, nil, nil)
	var rendered bytes.Buffer
	require.NoError(t, ParametersTemplate.Execute(&rendered, model))
	parameters := rendered.String()
	require.Contains(t, parameters, "*sdk.WarehouseParametersDetails")
	require.Contains(t, parameters, "sdk.ParameterTypeWarehouse")
	require.Contains(t, parameters, "*sdk.CreateWarehouseRequest")
	require.Contains(t, parameters, "*sdk.WarehouseSetRequest")
	require.Contains(t, parameters, "*sdk.WarehouseUnsetRequest")
}

func TestParameterHooks_variantRetainsOwnCatalog(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name:       "WarehouseInteractive",
		NamePlural: "Warehouses",
		IdType:     "AccountObjectIdentifier",
		VariantOf:  &ssotdefs.Warehouse,
		ResourceHooks: []ssot.ResourceHookOption{
			ssot.WithParametersAttributesSchema(ssot.Generated()),
			ssot.WithApplyParametersChanges(ssot.Generated()),
			ssot.WithParametersCustomDiff(ssot.Generated()),
		},
	}}, nil, nil)
	var rendered bytes.Buffer
	require.NoError(t, ParametersTemplate.Execute(&rendered, model))
	parameters := rendered.String()
	// FALLBACK_WAREHOUSE belongs to the interactive variant, not the base warehouse.
	require.Contains(t, parameters, `"fallback_warehouse": parameterSchemaFromDefinition(`)
	require.Contains(t, parameters, `handleParameterUpdateWithMapping(d, "fallback_warehouse", &set.FallbackWarehouse, &unset.FallbackWarehouse, stringToAccountObjectIdentifier),`)
	require.Contains(t, parameters, `IdentifierTypedParameterValueComputedIf("fallback_warehouse", parameters.FallbackWarehouse, sdk.ParameterTypeWarehouse),`)
	require.Contains(t, parameters, "*sdk.WarehouseParametersDetails")
	require.Contains(t, parameters, "*sdk.WarehouseInteractiveSetRequest")
	require.Contains(t, parameters, "*sdk.WarehouseInteractiveUnsetRequest")
	require.NotContains(t, parameters, "*sdk.WarehouseSetRequest")
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
