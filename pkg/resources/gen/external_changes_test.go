package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_ShowExternalChanges(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	kind := ssot.Attribute{Name: "warehouse_type", Type: ssot.AttrEnum, ShowName: "type", Enum: ssot.NewEnum("WarehouseType", "WarehouseTypes", "STANDARD")}
	size := ssot.Attribute{Name: "warehouse_size", Type: ssot.AttrEnum, ShowName: "size", OptionalInShow: true, Enum: ssot.NewEnum("WarehouseSize", "WarehouseSizes", "XSMALL")}
	autoResume := ssot.Attribute{Name: "auto_resume", Type: ssot.AttrBool}
	count := ssot.Attribute{Name: "max_cluster_count", Type: ssot.AttrInt, OptionalInShow: true, Min: new(0)}
	monitor := ssot.Attribute{Name: "resource_monitor", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"}
	created := ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true}

	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name:   "Warehouse",
		IdType: "AccountObjectIdentifier",
		Attributes: []ssot.Attribute{
			kind, size, autoResume, count, monitor, comment,
		},
		Show: ssot.OutputSpec{Columns: []*ssot.Attribute{
			&created, &kind, &size, &autoResume, &count, &monitor, &comment,
		}},
		ResourceHooks: []ssot.ResourceHookOption{
			ssot.WithHandleExternalChanges(ssot.Generated()),
		},
	}}, nil, nil)

	require.Equal(t, ExternalChangeBuild{
		Prepare: []string{
			"warehouseSizeVal, warehouseSizeStr := optionalStringOutputMapping(warehouse.Size)",
			"maxClusterCount := optionalIntOutputMappingIntDefault(warehouse.MaxClusterCount)",
		},
		Mappings: []ExternalChangeMapping{
			{OutputKey: "type", ConfigKey: "warehouse_type", Compare: "string(warehouse.Type)", Set: "warehouse.Type"},
			{OutputKey: "size", ConfigKey: "warehouse_size", Compare: "warehouseSizeStr", Set: "warehouseSizeVal"},
			{OutputKey: "auto_resume", ConfigKey: "auto_resume", Compare: "warehouse.AutoResume", Set: "booleanStringFromBool(warehouse.AutoResume)"},
			{OutputKey: "max_cluster_count", ConfigKey: "max_cluster_count", Compare: "maxClusterCount", Set: "maxClusterCount"},
			{OutputKey: "resource_monitor", ConfigKey: "resource_monitor", Compare: "warehouse.ResourceMonitor.Name()", Set: "warehouse.ResourceMonitor.Name()"},
		},
	}, model.ShowExternalChanges())

	got := renderTemplate(t, "generatedHandleExternalChanges", model)
	require.Contains(t, got, "func warehouseHandleExternalChanges(d *tfschema.ResourceData, warehouse *sdk.Warehouse, withExternalChangesMarking bool) error {")
	require.Contains(t, got, "\twarehouseSizeVal, warehouseSizeStr := optionalStringOutputMapping(warehouse.Size)\n")
	require.Contains(t, got, `outputMapping{"type", "warehouse_type", string(warehouse.Type), warehouse.Type, nil},`)
	require.Contains(t, got, "return handleExternalChangesToObjectInShow(d,")
	require.NotContains(t, got, "}warehouse")
	require.NotContains(t, got, "comment")
}

func Test_DescribeExternalChanges_catalogIntegrations(t *testing.T) {
	refresh := ExternalChangeMapping{
		OutputKey: "refresh_interval_seconds",
		ConfigKey: "refresh_interval_seconds",
		Compare:   "details.RefreshIntervalSeconds",
		Set:       "details.RefreshIntervalSeconds",
	}
	glue := ModelFromInputObject(ResourceObject{Object: defs.CatalogIntegrationAwsGlue}, nil, nil)
	require.Equal(t, []ExternalChangeMapping{
		refresh,
		{OutputKey: "glue_region", ConfigKey: "glue_region", Compare: "details.GlueRegion", Set: "details.GlueRegion"},
	}, glue.DescribeExternalChanges().Mappings)

	for _, obj := range []ssot.Object{
		defs.CatalogIntegrationObjectStorage,
		defs.CatalogIntegrationIcebergRest,
		defs.CatalogIntegrationOpenCatalog,
	} {
		model := ModelFromInputObject(ResourceObject{Object: obj}, nil, nil)
		require.Equal(t, []ExternalChangeMapping{refresh}, model.DescribeExternalChanges().Mappings)
	}

	got := renderTemplate(t, "generatedHandleExternalChangesInDescribe", glue)
	require.Contains(t, got, "func catalogIntegrationAwsGlueHandleExternalChangesInDescribe(d *tfschema.ResourceData, details *sdk.CatalogIntegrationAwsGlueDetails, withExternalChangesMarking bool) diag.Diagnostics {")
	require.Contains(t, got, `outputMapping{"glue_region", "glue_region", details.GlueRegion, details.GlueRegion, nil},`)
	require.NotContains(t, got, "comment")
	require.NotContains(t, got, "catalog_namespace")

	read := renderTemplate(t, "lifecycleTemplate", glue)
	require.Contains(t, read, "diags = catalogIntegrationAwsGlueHandleExternalChangesInDescribe(d, catalogIntegrationAwsGlueDetails, withExternalChangesMarking)")
	require.Contains(t, read, "diags = catalogIntegrationAwsGlueHandleExternalChangesInDescribeExt(d, catalogIntegrationAwsGlueDetails, withExternalChangesMarking)")
}

func Test_DerivedOutputDiffSuppress(t *testing.T) {
	byName := func(model ResourceModel) map[string]string {
		got := make(map[string]string, len(model.Attributes))
		for _, attr := range model.Attributes {
			got[attr.Name] = attr.DiffSuppressExpr
		}
		return got
	}

	warehouse := byName(ModelFromInputObject(ResourceObject{Object: defs.Warehouse}, nil, nil))
	require.Contains(t, warehouse["auto_resume"], `IgnoreChangeToCurrentSnowflakeValueInShow("auto_resume")`)
	require.Contains(t, warehouse["resource_monitor"], `IgnoreChangeToCurrentSnowflakeValueInShow("resource_monitor")`)
	require.NotContains(t, warehouse["comment"], "IgnoreChangeToCurrentSnowflakeValueInShow")
	require.NotContains(t, warehouse["initially_suspended"], "IgnoreChangeToCurrentSnowflakeValueInShow")

	adaptive := byName(ModelFromInputObject(ResourceObject{Object: defs.WarehouseAdaptive}, nil, nil))
	require.Contains(t, adaptive["max_query_performance_level"], `IgnoreChangeToCurrentSnowflakeValueInShow("max_query_performance_level")`)
	require.Contains(t, adaptive["query_throughput_multiplier"], `IgnoreChangeToCurrentSnowflakeValueInShow("query_throughput_multiplier")`)
	require.NotContains(t, adaptive["comment"], "IgnoreChangeToCurrentSnowflakeValueInShow")
	require.NotContains(t, adaptive["warehouse_type"], "IgnoreChangeToCurrentSnowflakeValueInShow")

	volume := byName(ModelFromInputObject(ResourceObject{Object: defs.ExternalVolumeObject}, nil, nil))
	require.Contains(t, volume["allow_writes"], `IgnoreChangeToCurrentSnowflakeValueInShow("allow_writes")`)
	require.NotContains(t, volume["comment"], "IgnoreChangeToCurrentSnowflakeValueInShow")

	glue := byName(ModelFromInputObject(ResourceObject{Object: defs.CatalogIntegrationAwsGlue}, nil, nil))
	require.Contains(t, glue["refresh_interval_seconds"], `IgnoreChangeToCurrentSnowflakeValueInDescribe("refresh_interval_seconds")`)
	require.Contains(t, glue["glue_region"], `IgnoreChangeToCurrentSnowflakeValueInDescribe("glue_region")`)
	require.Empty(t, glue["catalog_namespace"])
	require.Empty(t, glue["glue_aws_role_arn"])

	pool := byName(ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil))
	require.Contains(t, pool["for_application"], `IgnoreChangeToCurrentSnowflakeValueInShow("application")`)
	require.NotContains(t, pool["min_nodes"], "IgnoreChangeToCurrentSnowflake")
}

func Test_HandleExternalChanges_extCallUnchanged(t *testing.T) {
	model := lifecycleModel()
	require.Contains(t, renderTemplate(t, "handleExternalChanges", model), "poolHandleExternalChangesExt")
	require.Empty(t, model.ShowExternalChanges().Mappings)
	require.NotContains(t, renderLifecycle(t, model), "func poolHandleExternalChanges(")
}

func Test_ExternalChanges_panicsOnNestedDescribeColumn(t *testing.T) {
	nested := ssot.Attribute{
		Name:          "mfa_policy",
		Type:          ssot.AttrList,
		ServerDefault: true,
		Children: []ssot.Attribute{
			{Name: "allowed_methods", Type: ssot.AttrSet},
		},
	}
	require.Panics(t, func() {
		ModelFromInputObject(ResourceObject{Object: ssot.Object{
			Name:       "Policy",
			IdType:     "SchemaObjectIdentifier",
			Attributes: []ssot.Attribute{nested},
			Describe:   ssot.OutputSpec{Columns: []*ssot.Attribute{&nested}},
			ResourceHooks: []ssot.ResourceHookOption{
				ssot.WithHandleExternalChangesInDescribe(ssot.Generated()),
			},
		}}, nil, nil)
	})
}
