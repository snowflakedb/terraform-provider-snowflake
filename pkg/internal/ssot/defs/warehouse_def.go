package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// WarehouseTypeUserSettable excludes ADAPTIVE and INTERACTIVE. Those types are separate resources.
// NamePlural is empty so the schema description is not built from sdk.AllWarehouseTypes.
var WarehouseTypeUserSettable = ssot.NewEnum(
	"WarehouseTypeUserSettable",
	"",
	"STANDARD", "SNOWPARK-OPTIMIZED",
)

var WarehouseSize = ssot.NewEnum(
	"WarehouseSize", "WarehouseSizes",
	"XSMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE", "XXLARGE", "XXXLARGE", "X4LARGE", "X5LARGE", "X6LARGE",
)

var ScalingPolicy = ssot.NewEnum(
	"ScalingPolicy", "ScalingPolicies",
	"STANDARD", "ECONOMY",
)

var WarehouseResourceConstraint = ssot.NewEnum(
	"WarehouseResourceConstraint", "WarehouseResourceConstraints",
	"MEMORY_1X", "MEMORY_1X_x86", "MEMORY_16X", "MEMORY_16X_x86", "MEMORY_64X", "MEMORY_64X_x86",
)

var WarehouseGeneration = ssot.NewEnum(
	"WarehouseGeneration", "WarehouseGenerations",
	"1", "2",
)

// WarehouseAttrs holds reusable warehouse fields.
//
// Omitted from Attributes:
//   - initially_suspended: create-only native bool. Optional bools become boolean strings, so the field is added in warehouseAttributesModificationExt.
//
// show_output stays on schemas.ShowWarehouseRegularSchema. The resource has no describe_output.
var WarehouseAttrs = struct {
	WarehouseType                   ssot.Attribute
	WarehouseSize                   ssot.Attribute
	MaxClusterCount                 ssot.Attribute
	MinClusterCount                 ssot.Attribute
	ScalingPolicy                   ssot.Attribute
	AutoSuspend                     ssot.Attribute
	AutoResume                      ssot.Attribute
	ResourceMonitor                 ssot.Attribute
	Comment                         ssot.Attribute
	EnableQueryAcceleration         ssot.Attribute
	QueryAccelerationMaxScaleFactor ssot.Attribute
	ResourceConstraint              ssot.Attribute
	Generation                      ssot.Attribute
}{
	WarehouseType: ssot.Attribute{
		Name:        "warehouse_type",
		Type:        ssot.AttrEnum,
		Description: "Specifies warehouse type.",
		Enum:        WarehouseTypeUserSettable,
		ShowName:    "type",
	},
	WarehouseSize: ssot.Attribute{
		Name:        "warehouse_size",
		Type:        ssot.AttrEnum,
		Description: "Specifies the size of the virtual warehouse.",
		Enum:        WarehouseSize,
		ShowName:    "size",
	},
	MaxClusterCount: ssot.Attribute{
		Name:        "max_cluster_count",
		Type:        ssot.AttrInt,
		Description: "Specifies the maximum number of server clusters for the warehouse.",
		Min:         new(1),
	},
	MinClusterCount: ssot.Attribute{
		Name:        "min_cluster_count",
		Type:        ssot.AttrInt,
		Description: "Specifies the minimum number of server clusters for the warehouse (only applies to multi-cluster warehouses).",
		Min:         new(1),
	},
	ScalingPolicy: ssot.Attribute{
		Name:        "scaling_policy",
		Type:        ssot.AttrEnum,
		Description: "Specifies the policy for automatically starting and shutting down clusters in a multi-cluster warehouse running in Auto-scale mode.",
		Enum:        ScalingPolicy,
	},
	AutoSuspend: ssot.Attribute{
		Name:        "auto_suspend",
		Type:        ssot.AttrInt,
		Description: "Specifies the number of seconds of inactivity after which a warehouse is automatically suspended.",
		Min:         new(0),
	},
	AutoResume: ssot.Attribute{
		Name:        "auto_resume",
		Type:        ssot.AttrBool,
		Description: "Specifies whether to automatically resume a warehouse when a SQL statement (e.g. query) is submitted to it.",
	},
	ResourceMonitor: ssot.Attribute{
		Name:        "resource_monitor",
		Type:        ssot.AttrIdentifier,
		Description: "Specifies the name of a resource monitor that is explicitly assigned to the warehouse.",
	},
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the warehouse.",
	},
	EnableQueryAcceleration: ssot.Attribute{
		Name:        "enable_query_acceleration",
		Type:        ssot.AttrBool,
		Description: "Specifies whether to enable the query acceleration service for queries that rely on this warehouse for compute resources.",
	},
	QueryAccelerationMaxScaleFactor: ssot.Attribute{
		Name:        "query_acceleration_max_scale_factor",
		Type:        ssot.AttrInt,
		Description: "Specifies the maximum scale factor for leasing compute resources for query acceleration. The scale factor is used as a multiplier based on warehouse size.",
		Min:         new(0),
	},
	ResourceConstraint: ssot.Attribute{
		Name:        "resource_constraint",
		Type:        ssot.AttrEnum,
		Description: "Specifies the resource constraint for the warehouse.",
		Enum:        WarehouseResourceConstraint,
	},
	Generation: ssot.Attribute{
		Name:                  "generation",
		Type:                  ssot.AttrEnum,
		Description:           "Specifies the generation for the warehouse. Only available for standard warehouses.",
		AdditionalDescription: "Gen2 warehouses are not available in all regions. Please consult the [Snowflake Gen2 Region Availability documentation](https://docs.snowflake.com/en/user-guide/warehouses-gen2#region-availability) prior to configuration.",
		Enum:                  WarehouseGeneration,
	},
}

var Warehouse = ssot.Object{
	Name:                   "Warehouse",
	NamePlural:             "Warehouses",
	IdType:                 "AccountObjectIdentifier",
	AttributesModification: true,
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithDescribeOutputSet(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
		ssot.WithShowOutputSet(ssot.Omitted()),
		ssot.WithShowOutputSchema(ssot.ExtCall()),
		ssot.WithIdentitySchema(ssot.ExtCall()),
		ssot.WithHierarchyRename(ssot.ExtCall()),
		ssot.WithSchemaVersion(ssot.ExtCall()),
		ssot.WithApplyParametersCreate(ssot.Generated()),
		ssot.WithShowParametersInSdk(ssot.Generated()),
		ssot.WithParametersDetailsFromRaw(ssot.Generated()),
		ssot.WithSetParametersFields(ssot.Generated()),
		ssot.WithParametersOutputSet(ssot.Generated()),
		ssot.WithApplyParametersChanges(ssot.Generated()),
		ssot.WithParametersAttributesSchema(ssot.Generated()),
		ssot.WithParametersOutputSchema(ssot.Generated()),
		ssot.WithParametersFieldNames(ssot.Generated()),
		ssot.WithParametersCustomDiff(ssot.Generated()),
	},
	Attributes: []ssot.Attribute{
		WarehouseAttrs.WarehouseType,
		WarehouseAttrs.WarehouseSize,
		WarehouseAttrs.MaxClusterCount,
		WarehouseAttrs.MinClusterCount,
		WarehouseAttrs.ScalingPolicy,
		WarehouseAttrs.AutoSuspend,
		WarehouseAttrs.AutoResume,
		WarehouseAttrs.ResourceMonitor,
		WarehouseAttrs.Comment,
		WarehouseAttrs.EnableQueryAcceleration,
		WarehouseAttrs.QueryAccelerationMaxScaleFactor,
		WarehouseAttrs.ResourceConstraint,
		WarehouseAttrs.Generation,
	},
}
