package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// MaxQueryPerformanceLevel lists canonical values only. The adaptive schema description uses
// sdk.AllMaxQueryPerformanceLevels, not the alias string slice.
var MaxQueryPerformanceLevel = ssot.NewEnum(
	"MaxQueryPerformanceLevel", "MaxQueryPerformanceLevels",
	"XSMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE", "XXLARGE", "XXXLARGE", "X4LARGE",
)

// WarehouseAdaptiveAttrs holds adaptive warehouse fields. None of them match WarehouseAttrs:
// comment and resource_monitor use adaptive wording, and warehouse_type is computed.
//
// show_output stays on schemas.ShowWarehouseAdaptiveSchema. The resource has no describe_output.
var WarehouseAdaptiveAttrs = struct {
	Comment                   ssot.Attribute
	WarehouseType             ssot.Attribute
	MaxQueryPerformanceLevel  ssot.Attribute
	QueryThroughputMultiplier ssot.Attribute
	ResourceMonitor           ssot.Attribute
}{
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the adaptive warehouse.",
	},
	// Computed in Terraform. The generator emits Optional; AttributesModificationExt clears it.
	WarehouseType: ssot.Attribute{
		Name:        "warehouse_type",
		Type:        ssot.AttrString,
		Description: "Specifies the type for the adaptive warehouse. This field is used for checking external changes and recreating the resource if needed.",
	},
	MaxQueryPerformanceLevel: ssot.Attribute{
		Name:        "max_query_performance_level",
		Type:        ssot.AttrEnum,
		Description: "Specifies the maximum query performance level for the adaptive warehouse. Determines the initial compute capacity.",
		Enum:        MaxQueryPerformanceLevel,
	},
	QueryThroughputMultiplier: ssot.Attribute{
		Name:        "query_throughput_multiplier",
		Type:        ssot.AttrInt,
		Description: "Specifies the query throughput multiplier for the adaptive warehouse.",
		Min:         new(0),
	},
	ResourceMonitor: ssot.Attribute{
		Name:            "resource_monitor",
		Type:            ssot.AttrIdentifier,
		IdentifierType:  "AccountObjectIdentifier",
		Description:     "Specifies the name of a resource monitor that is explicitly assigned to the adaptive warehouse.",
		RelatedResource: &ResourceMonitor,
		DiffSuppress:    ssot.SuppressShow,
	},
}

var WarehouseAdaptive = ssot.Object{
	Name:                       "WarehouseAdaptive",
	NamePlural:                 "Warehouses",
	IdType:                     "AccountObjectIdentifier",
	SdkCreateRequestObjectName: "AdaptiveWarehouse",
	VariantOf:                  &Warehouse,
	AttributesModification:     true,
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithDescribeOutputSet(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
		ssot.WithShowOutputSet(ssot.Omitted()),
		ssot.WithShowOutputSchema(ssot.ExtCall()),
		ssot.WithIdentitySchema(ssot.ExtCall()),
		ssot.WithSetStateToValuesFromConfig(ssot.Omitted()),
		ssot.WithHierarchyRename(ssot.ExtCall()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		ssot.WithApplyParametersCreate(ssot.Generated()),
		ssot.WithShowParametersInSdk(ssot.Generated()),
		ssot.WithParametersDetailsFromRaw(ssot.Generated()),
		ssot.WithSetParametersFields(ssot.Generated()),
		ssot.WithParametersOutputSet(ssot.Generated()),
		ssot.WithApplyParametersChanges(ssot.Generated()),
		ssot.WithParametersAttributesSchema(ssot.Generated()),
		ssot.WithParametersOutputSchema(ssot.ExtCall()),
		ssot.WithParametersFieldNames(ssot.Generated()),
		ssot.WithParametersCustomDiff(ssot.Generated()),
	},
	Attributes: []ssot.Attribute{
		WarehouseAdaptiveAttrs.Comment,
		WarehouseAdaptiveAttrs.WarehouseType,
		WarehouseAdaptiveAttrs.MaxQueryPerformanceLevel,
		WarehouseAdaptiveAttrs.QueryThroughputMultiplier,
		WarehouseAdaptiveAttrs.ResourceMonitor,
	},
}
