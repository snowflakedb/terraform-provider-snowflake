package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// Aliases are accepted by the SDK converter. The adaptive description lists canonical values
// because IncludeAliasedInTheDescription is unset.
var MaxQueryPerformanceLevel = ssot.NewEnum(
	"MaxQueryPerformanceLevel", "MaxQueryPerformanceLevels",
	"XSMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE", "XXLARGE", "XXXLARGE", "X4LARGE",
).WithAliases("XSMALL", "X-SMALL").
	WithAliases("XLARGE", "X-LARGE").
	WithAliases("XXLARGE", "2X-LARGE").
	WithAliases("XXXLARGE", "3X-LARGE").
	WithAliases("X4LARGE", "4X-LARGE")

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
		CreateSql:   ssot.Clause(ssot.SingleQuotes()),
		SetSql:      ssot.Clause(ssot.SingleQuotes()),
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
		CreateSql:       ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Equals()),
		SetSql:          ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Equals()),
	},
}

// Show columns that cannot reuse WarehouseAttrs.
// auto_resume is writable there and is not an adaptive attribute.
// The performance columns are Readonly there and these names are adaptive attributes.

// warehouseAdaptiveWarehouseType is always emitted on CREATE as a private static field.
// SqlOnly keeps it out of schema, resource, model, and assertions.
// The resource warehouse_type field stays a separate computed attribute.
var warehouseAdaptiveWarehouseType = ssot.Attribute{
	Name:      "warehouse_type_static",
	SqlOnly:   true,
	CreateSql: ssot.StaticSQL("warehouseType", "WAREHOUSE_TYPE = 'ADAPTIVE'"),
}

var (
	warehouseAdaptiveAutoResumeShow = ssot.Attribute{
		Name: "auto_resume", Type: ssot.AttrBool, Readonly: true,
	}
	warehouseAdaptiveMaxQueryPerformanceLevelShow = ssot.Attribute{
		Name: "max_query_performance_level", Type: ssot.AttrEnum, Enum: MaxQueryPerformanceLevel, OptionalInShow: true, OptionalInDb: true,
	}
	warehouseAdaptiveQueryThroughputMultiplierShow = ssot.Attribute{
		Name: "query_throughput_multiplier", Type: ssot.AttrInt, OptionalInShow: true, OptionalInDb: true,
	}
)

var WarehouseAdaptive = ssot.Object{
	Name:                       "WarehouseAdaptive",
	NamePlural:                 "Warehouses",
	IdType:                     "AccountObjectIdentifier",
	SdkCreateRequestObjectName: "AdaptiveWarehouse",
	VariantOf:                  &Warehouse,
	AttributesModification:     true,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		ssot.WithSchema(ssot.SchemaShow()),
		ssot.WithModel(),
		ssot.WithResourceAssert(),
		ssot.WithShowOutputAssert(ssot.ShowOutputShow()),
		ssot.WithObjectAssert(ssot.ObjectAssertShow(
			ssot.ObjectAssertShowFromBase("ShowByIDAdaptive"),
		)),
		ssot.WithObjectParametersAssert(),
		ssot.WithResourceParametersAssert(),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
		ssot.WithShowOutputSet(ssot.Omitted()),
		ssot.WithShowOutputSchema(ssot.ExtCall()),
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
		warehouseAdaptiveWarehouseType,
		WarehouseAdaptiveAttrs.MaxQueryPerformanceLevel,
		WarehouseAdaptiveAttrs.QueryThroughputMultiplier,
		WarehouseAdaptiveAttrs.ResourceMonitor,
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&WarehouseAttrs.Name,
			&WarehouseAttrs.State,
			&WarehouseAttrs.WarehouseType,
			&WarehouseAttrs.Running,
			&WarehouseAttrs.Queued,
			&WarehouseAttrs.IsDefault,
			&WarehouseAttrs.IsCurrent,
			&warehouseAdaptiveAutoResumeShow,
			&WarehouseAttrs.Available,
			&WarehouseAttrs.Provisioning,
			&WarehouseAttrs.Quiescing,
			&WarehouseAttrs.Other,
			&WarehouseAttrs.CreatedOn,
			&WarehouseAttrs.ResumedOn,
			&WarehouseAttrs.UpdatedOn,
			&WarehouseAttrs.Owner,
			&WarehouseAttrs.Comment,
			&WarehouseAttrs.ResourceMonitor,
			&WarehouseAttrs.OwnerRoleType,
			&warehouseAdaptiveMaxQueryPerformanceLevelShow,
			&warehouseAdaptiveQueryThroughputMultiplierShow,
		},
	},
	Create: ssot.CreateSpec{
		Doc:                   "https://docs.snowflake.com/en/user-guide/warehouses-adaptive",
		AdditionalValidations: true,
		Clauses: []*ssot.Attribute{
			&warehouseAdaptiveWarehouseType,
			&WarehouseAdaptiveAttrs.Comment,
			&WarehouseAdaptiveAttrs.MaxQueryPerformanceLevel,
			&WarehouseAdaptiveAttrs.QueryThroughputMultiplier,
			&WarehouseAdaptiveAttrs.ResourceMonitor,
			ssot.Parameters,
			ssot.Tags,
		},
	},
	Set: []*ssot.Attribute{
		&WarehouseAttrs.WarehouseType,
		&WarehouseAdaptiveAttrs.Comment,
		&WarehouseAdaptiveAttrs.ResourceMonitor,
		&WarehouseAdaptiveAttrs.MaxQueryPerformanceLevel,
		&WarehouseAdaptiveAttrs.QueryThroughputMultiplier,
		ssot.Parameters,
	},
	Unset: []*ssot.Attribute{
		&WarehouseAttrs.WarehouseType,
		&WarehouseAdaptiveAttrs.Comment,
		&WarehouseAdaptiveAttrs.ResourceMonitor,
		&WarehouseAdaptiveAttrs.MaxQueryPerformanceLevel,
		&WarehouseAdaptiveAttrs.QueryThroughputMultiplier,
		ssot.Parameters,
	},
	Alter: ssot.AlterSpec{
		Doc:                      "https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse",
		SetAdditionalValidations: true,
		Clauses: []ssot.AlterClause{
			ssot.AlterRenameTo(),
			ssot.AlterSet(),
			ssot.AlterUnset(),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
		},
	},
}
