package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// WarehouseInteractiveAttrs holds interactive warehouse fields.
// None of WarehouseAttrs match: every description says "interactive warehouse", auto_suspend
// rejects 0, auto_resume is a boolean string with that wording, warehouse_type is computed,
// and tables is interactive-only.
//
// WarehouseSize is the shared enum. sdk.ValidWarehouseSizesString is sdk.AllWarehouseSizesString,
// so enum validation and the listed values match the existing schema.
//
// show_output stays on schemas.ShowWarehouseInteractiveSchema. The resource has no describe_output.
// Parameter output schema stays in the extension because its description and SHOW schema are
// interactive-specific. fallback_warehouse state stays there too: generated set-fields writes
// AccountObjectIdentifier.FullyQualifiedName(), and this resource stores Name().
var WarehouseInteractiveAttrs = struct {
	WarehouseSize      ssot.Attribute
	MaxClusterCount    ssot.Attribute
	MinClusterCount    ssot.Attribute
	AutoSuspend        ssot.Attribute
	AutoResume         ssot.Attribute
	InitiallySuspended ssot.Attribute
	ResourceMonitor    ssot.Attribute
	Comment            ssot.Attribute
	Tables             ssot.Attribute
	WarehouseType      ssot.Attribute
}{
	WarehouseSize: ssot.Attribute{
		Name:                           "warehouse_size",
		Type:                           ssot.AttrEnum,
		Description:                    "Specifies the size of the interactive warehouse.",
		AdditionalDescription:          "Note: changing the size briefly suspends and resumes the warehouse to apply the resize (an interactive warehouse cannot be resized while running); removing the size from config will result in the resource recreation.",
		Enum:                           WarehouseSize,
		IncludeAliasedInTheDescription: true,
		ShowName:                       "size",
		CreateSql:                      ssot.Clause(ssot.SingleQuotes()),
		SetSql:                         ssot.Clause(ssot.SingleQuotes()),
	},
	MaxClusterCount: ssot.Attribute{
		Name:        "max_cluster_count",
		Type:        ssot.AttrInt,
		Description: "Specifies the maximum number of server clusters for the interactive warehouse.",
		Min:         new(1),
	},
	MinClusterCount: ssot.Attribute{
		Name:        "min_cluster_count",
		Type:        ssot.AttrInt,
		Description: "Specifies the minimum number of server clusters for the interactive warehouse (only applies to multi-cluster warehouses).",
		Min:         new(1),
	},
	// Min is 1, so the attribute rule does not emit IntDefault. AttributesModificationExt sets that sentinel.
	AutoSuspend: ssot.Attribute{
		Name:        "auto_suspend",
		Type:        ssot.AttrInt,
		Description: "Specifies the number of seconds of inactivity after which an interactive warehouse is automatically suspended.",
		Min:         new(1),
	},
	AutoResume: ssot.Attribute{
		Name:        "auto_resume",
		Type:        ssot.AttrBool,
		Description: "Specifies whether to automatically resume an interactive warehouse when a SQL statement (e.g. query) is submitted to it.",
	},
	InitiallySuspended: ssot.Attribute{
		Name:        "initially_suspended",
		Type:        ssot.AttrBool,
		CreateOnly:  true,
		NativeBool:  true,
		Description: "Specifies whether the interactive warehouse is created initially in the ‘Suspended’ state.",
	},
	ResourceMonitor: ssot.Attribute{
		Name:            "resource_monitor",
		Type:            ssot.AttrIdentifier,
		IdentifierType:  "AccountObjectIdentifier",
		Description:     "Specifies the name of a resource monitor that is explicitly assigned to the interactive warehouse.",
		RelatedResource: &ResourceMonitor,
		CreateSql:       ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Equals()),
		SetSql:          ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Equals()),
	},
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the interactive warehouse.",
	},
	Tables: ssot.Attribute{
		Name:        "tables",
		Type:        ssot.AttrSet,
		Description: "Specifies the fully qualified names of the tables associated with the interactive warehouse. Changes are applied incrementally (ADD TABLES / DROP TABLES) rather than by full re-association.",
		Elem: &ssot.Attribute{
			Type:           ssot.AttrIdentifier,
			IdentifierType: "SchemaObjectIdentifier",
		},
		CreateSql: ssot.ListParameter("SchemaObjectIdentifier", ssot.Parentheses(), ssot.NoEquals()),
	},
	// Computed in Terraform. The generator emits Optional; AttributesModificationExt clears it.
	// Computed skips output suppression: SHOW already has a column of this name.
	WarehouseType: ssot.Attribute{
		Name:        "warehouse_type",
		Type:        ssot.AttrString,
		Computed:    true,
		Description: "Specifies the type for the interactive warehouse. This field is used for checking external changes and recreating the resource if needed.",
	},
}

// tables is a writable attribute, so the Readonly WarehouseAttrs column cannot be reused.
// TODO [next PRs]: the resource field is a set and this column is an AttrList so SHOW could build a paired struct. A set now maps like a list, so Show.Columns can use WarehouseInteractiveAttrs.Tables and this column can go.
var warehouseInteractiveTablesShow = ssot.Attribute{
	Name: "tables", Type: ssot.AttrList, OptionalInDb: true,
	CustomParser: "ParseCommaSeparatedSchemaObjectIdentifierArray",
	Elem:         &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "SchemaObjectIdentifier"},
}

var WarehouseInteractive = ssot.Object{
	Name:                       "WarehouseInteractive",
	NamePlural:                 "Warehouses",
	IdType:                     "AccountObjectIdentifier",
	SdkCreateRequestObjectName: "InteractiveWarehouse",
	VariantOf:                  &Warehouse,
	VariantEntityCreateName:    "INTERACTIVE WAREHOUSE",
	AttributesModification:     true,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		ssot.WithSchema(ssot.SchemaShow()),
		ssot.WithModel(),
		ssot.WithResourceAssert(),
		ssot.WithShowOutputAssert(ssot.ShowOutputShow()),
		ssot.WithObjectAssert(ssot.ObjectAssertShow(
			ssot.ObjectAssertShowFromBase("ShowByIDInteractive"),
			ssot.ObjectAssertShowSkipFields(&warehouseInteractiveTablesShow),
		)),
		ssot.WithObjectParametersAssert(),
		ssot.WithResourceParametersAssert(),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseId(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
		ssot.WithShowOutputSet(ssot.Omitted()),
		ssot.WithShowOutputSchema(ssot.ExtCall()),
		ssot.WithHierarchyRename(ssot.ExtCall()),
		ssot.WithBeforeAlter(ssot.ExtCall()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		ssot.WithApplyParametersCreate(ssot.Generated()),
		ssot.WithShowParametersInSdk(ssot.Generated()),
		ssot.WithParametersDetailsFromRaw(ssot.Generated()),
		ssot.WithSetParametersFields(ssot.ExtCall()),
		ssot.WithParametersOutputSet(ssot.Generated()),
		ssot.WithApplyParametersChanges(ssot.Generated()),
		ssot.WithParametersAttributesSchema(ssot.Generated()),
		ssot.WithParametersOutputSchema(ssot.ExtCall()),
		ssot.WithParametersFieldNames(ssot.Generated()),
		ssot.WithParametersCustomDiff(ssot.Generated()),
	},
	Attributes: []ssot.Attribute{
		WarehouseInteractiveAttrs.WarehouseSize,
		WarehouseInteractiveAttrs.MaxClusterCount,
		WarehouseInteractiveAttrs.MinClusterCount,
		WarehouseInteractiveAttrs.AutoSuspend,
		WarehouseInteractiveAttrs.AutoResume,
		WarehouseInteractiveAttrs.InitiallySuspended,
		WarehouseInteractiveAttrs.ResourceMonitor,
		WarehouseInteractiveAttrs.Comment,
		WarehouseInteractiveAttrs.Tables,
		WarehouseInteractiveAttrs.WarehouseType,
	},
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
			&WarehouseAttrs.ResourceMonitor,
			&WarehouseAttrs.OwnerRoleType,
			&warehouseInteractiveTablesShow,
		},
	},
	Create: ssot.CreateSpec{
		Doc:                   "https://docs.snowflake.com/en/user-guide/warehouses-interactive",
		AdditionalValidations: true,
		Clauses: []*ssot.Attribute{
			&WarehouseInteractiveAttrs.Tables,
			&WarehouseInteractiveAttrs.WarehouseSize,
			&WarehouseInteractiveAttrs.MaxClusterCount,
			&WarehouseInteractiveAttrs.MinClusterCount,
			&WarehouseInteractiveAttrs.AutoSuspend,
			&WarehouseInteractiveAttrs.AutoResume,
			&WarehouseInteractiveAttrs.InitiallySuspended,
			&WarehouseInteractiveAttrs.ResourceMonitor,
			&WarehouseInteractiveAttrs.Comment,
			ssot.Parameters,
			ssot.Tags,
		},
	},
	Set: []*ssot.Attribute{
		&WarehouseInteractiveAttrs.WarehouseSize,
		&WarehouseInteractiveAttrs.MaxClusterCount,
		&WarehouseInteractiveAttrs.MinClusterCount,
		&WarehouseInteractiveAttrs.AutoSuspend,
		&WarehouseInteractiveAttrs.AutoResume,
		&WarehouseInteractiveAttrs.ResourceMonitor,
		&WarehouseInteractiveAttrs.Comment,
		ssot.Parameters,
	},
	Unset: []*ssot.Attribute{
		&WarehouseInteractiveAttrs.WarehouseSize,
		&WarehouseInteractiveAttrs.MaxClusterCount,
		&WarehouseInteractiveAttrs.MinClusterCount,
		&WarehouseInteractiveAttrs.AutoSuspend,
		&WarehouseInteractiveAttrs.AutoResume,
		&WarehouseInteractiveAttrs.ResourceMonitor,
		&WarehouseInteractiveAttrs.Comment,
		ssot.Parameters,
	},
	Alter: ssot.AlterSpec{
		Doc:                      "https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse",
		SetAdditionalValidations: true,
		Clauses: []ssot.AlterClause{
			ssot.AlterRenameTo(),
			ssot.AlterSet(),
			ssot.AlterUnset(),
			ssot.AlterList("ADD TABLES", "SchemaObjectIdentifier"),
			ssot.AlterList("DROP TABLES", "SchemaObjectIdentifier"),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
		},
	},
}
