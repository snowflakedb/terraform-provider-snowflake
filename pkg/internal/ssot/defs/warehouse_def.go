package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// WarehouseTypeUserSettable excludes ADAPTIVE and INTERACTIVE. Those types are separate resources.
// SdkType is WarehouseType: sdk.ToWarehouseTypeUserSettable returns that type.
var WarehouseTypeUserSettable = ssot.NewEnum(
	"WarehouseTypeUserSettable",
	"WarehouseTypesUserSettable",
	"STANDARD", "SNOWPARK-OPTIMIZED",
).WithSdkType("WarehouseType")

var WarehouseSize = ssot.NewEnum(
	"WarehouseSize", "WarehouseSizes",
	"XSMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE", "XXLARGE", "XXXLARGE", "X4LARGE", "X5LARGE", "X6LARGE",
).WithAliases("XSMALL", "X-SMALL").
	WithAliases("XLARGE", "X-LARGE").
	WithAliases("XXLARGE", "X2LARGE", "2X-LARGE").
	WithAliases("XXXLARGE", "X3LARGE", "3X-LARGE").
	WithAliases("X4LARGE", "4X-LARGE").
	WithAliases("X5LARGE", "5X-LARGE").
	WithAliases("X6LARGE", "6X-LARGE")

var ScalingPolicy = ssot.NewEnum(
	"ScalingPolicy", "ScalingPolicies",
	"STANDARD", "ECONOMY",
)

var WarehouseResourceConstraint = ssot.NewEnum(
	"WarehouseResourceConstraint", "WarehouseResourceConstraints",
	"MEMORY_1X", "MEMORY_1X_x86", "MEMORY_16X", "MEMORY_16X_x86", "MEMORY_64X", "MEMORY_64X_x86",
)

// WithValue keeps WarehouseGenerationStandardGen1 and WarehouseGenerationStandardGen2.
// The SQL values are "1" and "2", which would otherwise generate WarehouseGeneration1 and WarehouseGeneration2.
var WarehouseGeneration = ssot.NewEnum(
	"WarehouseGeneration", "WarehouseGenerations",
).WithValue("StandardGen1", "1").WithValue("StandardGen2", "2")

// WarehouseType is the SHOW type, including ADAPTIVE and INTERACTIVE.
// WarehouseTypeUserSettable is the regular resource's writable enum.
var WarehouseType = ssot.NewEnum(
	"WarehouseType", "WarehouseTypes",
	"STANDARD", "SNOWPARK-OPTIMIZED", "ADAPTIVE", "INTERACTIVE",
)

// ENABLED is reported for adaptive warehouses, which have no suspend/resume lifecycle.
var WarehouseState = ssot.NewEnum(
	"WarehouseState", "WarehouseStates",
	"SUSPENDED", "SUSPENDING", "STARTED", "RESIZING", "RESUMING", "ENABLED",
)

// WarehouseAttrs holds reusable warehouse fields.
//
// initially_suspended is a create-only native bool (CreateOnly and NativeBool). It is not a SHOW column.
// show_output stays on schemas.ShowWarehouseRegularSchema. The resource has no describe_output.
// DESCRIBE columns match sdk.WarehouseDetails.
var WarehouseAttrs = struct {
	WarehouseType                   ssot.Attribute
	WarehouseSize                   ssot.Attribute
	MaxClusterCount                 ssot.Attribute
	MinClusterCount                 ssot.Attribute
	ScalingPolicy                   ssot.Attribute
	AutoSuspend                     ssot.Attribute
	AutoResume                      ssot.Attribute
	InitiallySuspended              ssot.Attribute
	ResourceMonitor                 ssot.Attribute
	Comment                         ssot.Attribute
	EnableQueryAcceleration         ssot.Attribute
	QueryAccelerationMaxScaleFactor ssot.Attribute
	ResourceConstraint              ssot.Attribute
	Generation                      ssot.Attribute
	Name                            ssot.Attribute
	State                           ssot.Attribute
	StartedClusters                 ssot.Attribute
	Running                         ssot.Attribute
	Queued                          ssot.Attribute
	IsDefault                       ssot.Attribute
	IsCurrent                       ssot.Attribute
	Available                       ssot.Attribute
	Provisioning                    ssot.Attribute
	Quiescing                       ssot.Attribute
	Other                           ssot.Attribute
	CreatedOn                       ssot.Attribute
	ResumedOn                       ssot.Attribute
	UpdatedOn                       ssot.Attribute
	Owner                           ssot.Attribute
	Actives                         ssot.Attribute
	Pendings                        ssot.Attribute
	Failed                          ssot.Attribute
	Suspended                       ssot.Attribute
	Uuid                            ssot.Attribute
	OwnerRoleType                   ssot.Attribute
	MaxQueryPerformanceLevel        ssot.Attribute
	QueryThroughputMultiplier       ssot.Attribute
	Tables                          ssot.Attribute
	Kind                            ssot.Attribute
	WaitForCompletion               ssot.Attribute
}{
	WarehouseType: ssot.Attribute{
		Name:                  "warehouse_type",
		Type:                  ssot.AttrEnum,
		Description:           "Specifies warehouse type.",
		AdditionalDescription: "Warehouse needs to be suspended to change its type. Provider will handle automatic suspension and resumption if needed.",
		// Resource validation stays WarehouseTypeUserSettable (STANDARD, SNOWPARK-OPTIMIZED).
		// SHOW type is WarehouseType, which also includes ADAPTIVE and INTERACTIVE.
		Enum:          WarehouseTypeUserSettable,
		ShowName:      "type",
		ShowEnum:      WarehouseType,
		ManualConvert: true,
		DiffSuppress:  ssot.SuppressShow,
		CreateSql:     ssot.Clause(ssot.SingleQuotes()),
		SetSql:        ssot.Clause(ssot.SingleQuotes()),
	},
	WarehouseSize: ssot.Attribute{
		Name:                           "warehouse_size",
		Type:                           ssot.AttrEnum,
		Description:                    "Specifies the size of the virtual warehouse.",
		AdditionalDescription:          "Consult [warehouse documentation](https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties) for the details. Note: removing the size from config will result in the resource recreation.",
		Enum:                           WarehouseSize,
		IncludeAliasedInTheDescription: true,
		ShowName:                       "size",
		DiffSuppress:                   ssot.SuppressShow,
		OptionalInShow:                 true,
		OptionalInDb:                   true,
		CreateSql:                      ssot.Clause(ssot.SingleQuotes()),
		SetSql:                         ssot.Clause(ssot.SingleQuotes()),
	},
	MaxClusterCount: ssot.Attribute{
		Name:           "max_cluster_count",
		Type:           ssot.AttrInt,
		Description:    "Specifies the maximum number of server clusters for the warehouse.",
		Min:            new(1),
		DiffSuppress:   ssot.SuppressShow,
		OptionalInShow: true,
		OptionalInDb:   true,
	},
	MinClusterCount: ssot.Attribute{
		Name:           "min_cluster_count",
		Type:           ssot.AttrInt,
		Description:    "Specifies the minimum number of server clusters for the warehouse (only applies to multi-cluster warehouses).",
		Min:            new(1),
		DiffSuppress:   ssot.SuppressShow,
		OptionalInShow: true,
		OptionalInDb:   true,
	},
	ScalingPolicy: ssot.Attribute{
		Name:           "scaling_policy",
		Type:           ssot.AttrEnum,
		Description:    "Specifies the policy for automatically starting and shutting down clusters in a multi-cluster warehouse running in Auto-scale mode.",
		Enum:           ScalingPolicy,
		DiffSuppress:   ssot.SuppressShow,
		OptionalInShow: true,
		OptionalInDb:   true,
		CreateSql:      ssot.Clause(ssot.SingleQuotes()),
		SetSql:         ssot.Clause(ssot.SingleQuotes()),
	},
	AutoSuspend: ssot.Attribute{
		Name:           "auto_suspend",
		Type:           ssot.AttrInt,
		Description:    "Specifies the number of seconds of inactivity after which a warehouse is automatically suspended.",
		Min:            new(0),
		DiffSuppress:   ssot.SuppressShow,
		OptionalInShow: true,
		OptionalInDb:   true,
	},
	AutoResume: ssot.Attribute{
		Name:         "auto_resume",
		Type:         ssot.AttrBool,
		Description:  "Specifies whether to automatically resume a warehouse when a SQL statement (e.g. query) is submitted to it.",
		DiffSuppress: ssot.SuppressShow,
	},
	InitiallySuspended: ssot.Attribute{
		Name:        "initially_suspended",
		Type:        ssot.AttrBool,
		CreateOnly:  true,
		NativeBool:  true,
		Description: "Specifies whether the warehouse is created initially in the ‘Suspended’ state.",
	},
	ResourceMonitor: ssot.Attribute{
		Name:            "resource_monitor",
		Type:            ssot.AttrIdentifier,
		IdentifierType:  "AccountObjectIdentifier",
		Description:     "Specifies the name of a resource monitor that is explicitly assigned to the warehouse.",
		CustomParser:    "ParseAccountObjectIdentifierExcludingExplicitNullString",
		RelatedResource: &ResourceMonitor,
		DiffSuppress:    ssot.SuppressShow,
		CreateSql:       ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Equals()),
		SetSql:          ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Equals()),
	},
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the warehouse.",
	},
	EnableQueryAcceleration: ssot.Attribute{
		Name:           "enable_query_acceleration",
		Type:           ssot.AttrBool,
		Description:    "Specifies whether to enable the query acceleration service for queries that rely on this warehouse for compute resources.",
		DiffSuppress:   ssot.SuppressShow,
		OptionalInShow: true,
		OptionalInDb:   true,
	},
	QueryAccelerationMaxScaleFactor: ssot.Attribute{
		Name:           "query_acceleration_max_scale_factor",
		Type:           ssot.AttrInt,
		Description:    "Specifies the maximum scale factor for leasing compute resources for query acceleration. The scale factor is used as a multiplier based on warehouse size.",
		Min:            new(0),
		Max:            new(100),
		DiffSuppress:   ssot.SuppressShow,
		OptionalInShow: true,
		OptionalInDb:   true,
	},
	ResourceConstraint: ssot.Attribute{
		Name:                  "resource_constraint",
		Type:                  ssot.AttrEnum,
		Description:           "Specifies the resource constraint for the warehouse. Only available for snowpark-optimized warehouses. For setting generation please use the `generation` field. Please check [Snowflake documentation](https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties) for required warehouse sizes for each resource constraint.",
		AdditionalDescription: "Gen2 warehouses are not available in all regions. Please consult the [Snowflake Gen2 Region Availability documentation](https://docs.snowflake.com/en/user-guide/warehouses-gen2#region-availability) prior to configuration.",
		Enum:                  WarehouseResourceConstraint,
		DiffSuppress:          ssot.SuppressShow,
		OptionalInShow:        true,
		OptionalInDb:          true,
		ManualConvert:         true,
		CreateSql:             ssot.Clause(ssot.SingleQuotes()),
		SetSql:                ssot.Clause(ssot.SingleQuotes()),
	},
	Generation: ssot.Attribute{
		Name:                  "generation",
		Type:                  ssot.AttrEnum,
		Description:           "Specifies the generation for the warehouse. Only available for standard warehouses.",
		AdditionalDescription: "Gen2 warehouses are not available in all regions. Please consult the [Snowflake Gen2 Region Availability documentation](https://docs.snowflake.com/en/user-guide/warehouses-gen2#region-availability) prior to configuration.",
		Enum:                  WarehouseGeneration,
		DiffSuppress:          ssot.SuppressShow,
		OptionalInShow:        true,
		OptionalInDb:          true,
		ManualConvert:         true,
		CreateSql:             ssot.Clause(ssot.SingleQuotes()),
		SetSql:                ssot.Clause(ssot.SingleQuotes()),
	},
	// wait_for_completion is an ALTER SET field. It is not a resource attribute and not a SHOW column.
	WaitForCompletion: ssot.Attribute{
		Name:        "wait_for_completion",
		Type:        ssot.AttrBool,
		Description: "Specifies whether to wait until the warehouse finishes resizing before returning.",
	},
	Name:            ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	State:           ssot.Attribute{Name: "state", Type: ssot.AttrEnum, Enum: WarehouseState, Readonly: true, ManualConvert: true},
	StartedClusters: ssot.Attribute{Name: "started_clusters", Type: ssot.AttrInt, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	Running:         ssot.Attribute{Name: "running", Type: ssot.AttrInt, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	Queued:          ssot.Attribute{Name: "queued", Type: ssot.AttrInt, Readonly: true, OptionalInShow: true, OptionalInDb: true},
	IsDefault:       ssot.Attribute{Name: "is_default", Type: ssot.AttrBool, StoredAsText: true, Readonly: true},
	IsCurrent:       ssot.Attribute{Name: "is_current", Type: ssot.AttrBool, StoredAsText: true, Readonly: true},
	Available:       ssot.Attribute{Name: "available", Type: ssot.AttrFloat, Readonly: true},
	Provisioning:    ssot.Attribute{Name: "provisioning", Type: ssot.AttrFloat, Readonly: true},
	Quiescing:       ssot.Attribute{Name: "quiescing", Type: ssot.AttrFloat, Readonly: true},
	Other:           ssot.Attribute{Name: "other", Type: ssot.AttrFloat, Readonly: true},
	CreatedOn:       ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	ResumedOn:       ssot.Attribute{Name: "resumed_on", Type: ssot.AttrTime, Readonly: true},
	UpdatedOn:       ssot.Attribute{Name: "updated_on", Type: ssot.AttrTime, Readonly: true},
	Owner:           ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	Actives:         ssot.Attribute{Name: "actives", Type: ssot.AttrString, Readonly: true, SkipInSchema: true},
	Pendings:        ssot.Attribute{Name: "pendings", Type: ssot.AttrString, Readonly: true, SkipInSchema: true},
	Failed:          ssot.Attribute{Name: "failed", Type: ssot.AttrString, Readonly: true, SkipInSchema: true},
	Suspended:       ssot.Attribute{Name: "suspended", Type: ssot.AttrString, Readonly: true, SkipInSchema: true},
	Uuid:            ssot.Attribute{Name: "uuid", Type: ssot.AttrString, Readonly: true, SkipInSchema: true},
	OwnerRoleType:   ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, Readonly: true, OptionalInDb: true},
	// SHOW columns for the union warehouse row. Adaptive ALTER writes the adaptive attributes, not these copies.
	MaxQueryPerformanceLevel: ssot.Attribute{
		Name: "max_query_performance_level", Type: ssot.AttrEnum, Enum: MaxQueryPerformanceLevel, Readonly: true, OptionalInShow: true, OptionalInDb: true,
	},
	QueryThroughputMultiplier: ssot.Attribute{
		Name: "query_throughput_multiplier", Type: ssot.AttrInt, Readonly: true, OptionalInShow: true, OptionalInDb: true,
	},
	Tables: ssot.Attribute{
		Name: "tables", Type: ssot.AttrList, Readonly: true, OptionalInDb: true,
		CustomParser: "ParseCommaSeparatedSchemaObjectIdentifierArray",
		Elem:         &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "SchemaObjectIdentifier"},
	},
	Kind: ssot.Attribute{Name: "kind", Type: ssot.AttrString, Readonly: true},
}

// ConflictsWith is assigned here. The two attributes reference each other, so setting it inside the
// WarehouseAttrs literal is an initialization cycle.
func init() {
	WarehouseAttrs.ResourceConstraint.ConflictsWith = []ssot.Attribute{WarehouseAttrs.Generation}
	WarehouseAttrs.Generation.ConflictsWith = []ssot.Attribute{WarehouseAttrs.ResourceConstraint}
	for i := range Warehouse.Attributes {
		switch Warehouse.Attributes[i].Name {
		case WarehouseAttrs.ResourceConstraint.Name:
			Warehouse.Attributes[i].ConflictsWith = WarehouseAttrs.ResourceConstraint.ConflictsWith
		case WarehouseAttrs.Generation.Name:
			Warehouse.Attributes[i].ConflictsWith = WarehouseAttrs.Generation.ConflictsWith
		}
	}
	Warehouse.ShowVariants = []*ssot.Object{&WarehouseAdaptive, &WarehouseInteractive}
}

var Warehouse = ssot.Object{
	Name:       "Warehouse",
	NamePlural: "Warehouses",
	IdType:     "AccountObjectIdentifier",
	// ShowParametersDetails lists interactive parameters such as FALLBACK_WAREHOUSE.
	ShowParametersDetailsAdditionalLevels: []ssot.ParameterLevel{ssot.ParameterLevelWarehouseInteractive},
	CustomInterfaceMethods: []ssot.CustomInterfaceMethodSpec{
		{
			Name:    "ShowByIDExperimental",
			Doc:     "ShowByIDExperimental is a show by id function with improved performance (using starts with and limit)",
			Returns: []string{"*Warehouse", "error"},
		},
		{
			Name:    "ShowByIDExperimentalSafely",
			Returns: []string{"*Warehouse", "error"},
		},
		{
			Name: "AlterWithSuspend",
			Doc:  "AlterWithSuspend wraps Alter with automatic suspend/resume when changing warehouse type, or resizing an interactive warehouse",
			Parameters: []ssot.CustomInterfaceMethodParameter{{
				Name: "request",
				Kind: "*AlterWarehouseRequest",
			}},
		},
		{
			Name: "AlterAdaptiveWithSuspend",
			Doc:  "AlterAdaptiveWithSuspend wraps AlterAdaptive and suspends the warehouse before changing its type. Adaptive warehouses are not resumed afterwards",
			Parameters: []ssot.CustomInterfaceMethodParameter{{
				Name: "request",
				Kind: "*AlterAdaptiveWarehouseRequest",
			}},
		},
		{
			Name: "AlterInteractiveWithSuspend",
			Doc:  "AlterInteractiveWithSuspend wraps AlterInteractive and suspends a running interactive warehouse before changing its size",
			Parameters: []ssot.CustomInterfaceMethodParameter{{
				Name: "request",
				Kind: "*AlterInteractiveWarehouseRequest",
			}},
		},
		{
			Name: "CreateInteractivePreservingSession",
			Doc:  "CreateInteractivePreservingSession wraps CreateInteractive to work around a known Snowflake behavior: creating a warehouse implicitly switches the session to it, which for interactive warehouses (5 second statement timeout) can cause subsequent statements in the same session to time out. It restores whichever warehouse (or lack of one) was active in the session before creation",
			Parameters: []ssot.CustomInterfaceMethodParameter{{
				Name: "request",
				Kind: "*CreateInteractiveWarehouseRequest",
			}},
		},
	},
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		ssot.WithSchema(ssot.SchemaShow(), ssot.SchemaDescribe()),
		ssot.WithModel(),
		ssot.WithResourceAssert(),
		ssot.WithShowOutputAssert(
			ssot.ShowOutputShow(ssot.ShowOutputShowSkipFields(&WarehouseAttrs.Tables)),
			ssot.ShowOutputDescribe(),
		),
		ssot.WithObjectAssert(
			ssot.ObjectAssertShow(ssot.ObjectAssertShowSkipFields(&WarehouseAttrs.Tables)),
			ssot.ObjectAssertDescribe(),
		),
		ssot.WithObjectParametersAssert(),
		ssot.WithResourceParametersAssert(),
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkShow(
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalStartsWith(),
				ssot.SdkShowOptionalLimitFrom(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDLike),
			),
			ssot.SdkDescribe(ssot.SdkDescribeQueryStructName("DescribeWarehouse")),
			ssot.SdkCreate(ssot.SdkCreateAdditionalValidations()),
			ssot.SdkCreateForVariants(),
			ssot.SdkAlterForVariants(),
			ssot.SdkAlter(
				ssot.SdkAlterSetAdditionalValidations(),
				ssot.SdkAlterAdditionalValidations(),
			),
			ssot.SdkDrop(),
			ssot.SdkShowParameters(),
			ssot.SdkShowParametersDetails(),
			ssot.SdkCustomInterfaceMethods(),
		),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
		ssot.WithShowOutputSet(ssot.Omitted()),
		ssot.WithShowOutputSchema(ssot.ExtCall()),
		ssot.WithIdentitySchema(ssot.ExtCall()),
		ssot.WithHierarchyRename(ssot.ExtCall()),
		ssot.WithStateUpgraders(ssot.Generated()),
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
		ssot.WithCustomizeDiff(ssot.Generated()),
	},
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"v092WarehouseSizeStateUpgrader",
		"v2_6_0_WarehouseResourceConstraintUpgrader",
	},
	AdditionalCustomizeDiffs: []ssot.CustomizeDiffFunc{
		"warehouseForceNewIfSizeRemoved",
		"warehouseForceNewIfWasInteractive",
	},
	Attributes: []ssot.Attribute{
		WarehouseAttrs.WarehouseType,
		WarehouseAttrs.WarehouseSize,
		WarehouseAttrs.MaxClusterCount,
		WarehouseAttrs.MinClusterCount,
		WarehouseAttrs.ScalingPolicy,
		WarehouseAttrs.AutoSuspend,
		WarehouseAttrs.AutoResume,
		WarehouseAttrs.InitiallySuspended,
		WarehouseAttrs.ResourceMonitor,
		WarehouseAttrs.Comment,
		WarehouseAttrs.EnableQueryAcceleration,
		WarehouseAttrs.QueryAccelerationMaxScaleFactor,
		WarehouseAttrs.ResourceConstraint,
		WarehouseAttrs.Generation,
	},
	Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{
		&WarehouseAttrs.WarehouseType,
		&WarehouseAttrs.WarehouseSize,
		&WarehouseAttrs.MaxClusterCount,
		&WarehouseAttrs.MinClusterCount,
		&WarehouseAttrs.ScalingPolicy,
		&WarehouseAttrs.AutoSuspend,
		&WarehouseAttrs.AutoResume,
		&WarehouseAttrs.InitiallySuspended,
		&WarehouseAttrs.ResourceMonitor,
		&WarehouseAttrs.Comment,
		&WarehouseAttrs.EnableQueryAcceleration,
		&WarehouseAttrs.QueryAccelerationMaxScaleFactor,
		&WarehouseAttrs.ResourceConstraint,
		&WarehouseAttrs.Generation,
		ssot.Parameters,
		ssot.Tags,
	}},
	Set: []*ssot.Attribute{
		&WarehouseAttrs.WarehouseType,
		&WarehouseAttrs.WarehouseSize,
		&WarehouseAttrs.WaitForCompletion,
		&WarehouseAttrs.MaxClusterCount,
		&WarehouseAttrs.MinClusterCount,
		&WarehouseAttrs.ScalingPolicy,
		&WarehouseAttrs.AutoSuspend,
		&WarehouseAttrs.AutoResume,
		&WarehouseAttrs.ResourceMonitor,
		&WarehouseAttrs.Comment,
		&WarehouseAttrs.EnableQueryAcceleration,
		&WarehouseAttrs.QueryAccelerationMaxScaleFactor,
		&WarehouseAttrs.ResourceConstraint,
		&WarehouseAttrs.Generation,
		ssot.Parameters,
	},
	Unset: []*ssot.Attribute{
		&WarehouseAttrs.WarehouseType,
		&WarehouseAttrs.WaitForCompletion,
		&WarehouseAttrs.MaxClusterCount,
		&WarehouseAttrs.MinClusterCount,
		&WarehouseAttrs.ScalingPolicy,
		&WarehouseAttrs.AutoSuspend,
		&WarehouseAttrs.AutoResume,
		&WarehouseAttrs.ResourceMonitor,
		&WarehouseAttrs.Comment,
		&WarehouseAttrs.EnableQueryAcceleration,
		&WarehouseAttrs.QueryAccelerationMaxScaleFactor,
		&WarehouseAttrs.ResourceConstraint,
		&WarehouseAttrs.Generation,
		ssot.Parameters,
	},
	Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
		ssot.AlterKeyword("SUSPEND"),
		ssot.AlterKeyword("RESUME"),
		ssot.AlterKeywordNotExactlyOne("IF SUSPENDED"),
		ssot.AlterKeyword("ABORT ALL QUERIES"),
		ssot.AlterRenameTo(),
		ssot.AlterSet(),
		ssot.AlterUnset(),
		ssot.AlterSetTags(),
		ssot.AlterUnsetTags(),
	}},
	ShowPlainName: "WarehouseRegular",
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&WarehouseAttrs.Name,
			&WarehouseAttrs.State,
			&WarehouseAttrs.WarehouseType,
			&WarehouseAttrs.WarehouseSize,
			&WarehouseAttrs.MinClusterCount,
			&WarehouseAttrs.MaxClusterCount,
			&WarehouseAttrs.StartedClusters,
			&WarehouseAttrs.Running,
			&WarehouseAttrs.Queued,
			&WarehouseAttrs.IsDefault,
			&WarehouseAttrs.IsCurrent,
			&WarehouseAttrs.AutoSuspend,
			&WarehouseAttrs.AutoResume,
			&WarehouseAttrs.Available,
			&WarehouseAttrs.Provisioning,
			&WarehouseAttrs.Quiescing,
			&WarehouseAttrs.Other,
			&WarehouseAttrs.CreatedOn,
			&WarehouseAttrs.ResumedOn,
			&WarehouseAttrs.UpdatedOn,
			&WarehouseAttrs.Owner,
			&WarehouseAttrs.Comment,
			&WarehouseAttrs.EnableQueryAcceleration,
			&WarehouseAttrs.QueryAccelerationMaxScaleFactor,
			&WarehouseAttrs.ResourceMonitor,
			&WarehouseAttrs.ScalingPolicy,
			&WarehouseAttrs.OwnerRoleType,
			&WarehouseAttrs.ResourceConstraint,
			&WarehouseAttrs.Generation,
		},
	},
	// Union of all SHOW WAREHOUSES columns (AllDetails analog).
	// Keep omitted: actives / pendings / failed / suspended / uuid (Snowflake: internal use, will be removed).
	ShowGeneric: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&WarehouseAttrs.Name,
			&WarehouseAttrs.State,
			&WarehouseAttrs.WarehouseType,
			&WarehouseAttrs.WarehouseSize,
			&WarehouseAttrs.MinClusterCount,
			&WarehouseAttrs.MaxClusterCount,
			&WarehouseAttrs.StartedClusters,
			&WarehouseAttrs.Running,
			&WarehouseAttrs.Queued,
			&WarehouseAttrs.IsDefault,
			&WarehouseAttrs.IsCurrent,
			&WarehouseAttrs.AutoSuspend,
			&WarehouseAttrs.AutoResume,
			&WarehouseAttrs.Available,
			&WarehouseAttrs.Provisioning,
			&WarehouseAttrs.Quiescing,
			&WarehouseAttrs.Other,
			&WarehouseAttrs.CreatedOn,
			&WarehouseAttrs.ResumedOn,
			&WarehouseAttrs.UpdatedOn,
			&WarehouseAttrs.Owner,
			&WarehouseAttrs.Comment,
			&WarehouseAttrs.EnableQueryAcceleration,
			&WarehouseAttrs.QueryAccelerationMaxScaleFactor,
			&WarehouseAttrs.ResourceMonitor,
			&WarehouseAttrs.Actives,
			&WarehouseAttrs.Pendings,
			&WarehouseAttrs.Failed,
			&WarehouseAttrs.Suspended,
			&WarehouseAttrs.Uuid,
			&WarehouseAttrs.ScalingPolicy,
			&WarehouseAttrs.OwnerRoleType,
			&WarehouseAttrs.ResourceConstraint,
			&WarehouseAttrs.Generation,
			&WarehouseAttrs.MaxQueryPerformanceLevel,
			&WarehouseAttrs.QueryThroughputMultiplier,
			&WarehouseAttrs.Tables,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&WarehouseAttrs.CreatedOn,
			&WarehouseAttrs.Name,
			&WarehouseAttrs.Kind,
		},
	},
}
