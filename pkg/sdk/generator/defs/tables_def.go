package defs

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var tableSetColumnMaskingPolicy = g.QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnMaskingPolicy, ssot.ClauseSet).
	WithSharedToOpts()

var tableUnsetColumnMaskingPolicy = g.QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnMaskingPolicy, ssot.ClauseUnset).
	WithSharedToOpts()

var tableSetColumnProjectionPolicy = g.NewQueryStruct("TableSetColumnProjectionPolicy").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("SET").
	Identifier("ProjectionPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PROJECTION POLICY").Required()).
	OptionalSQL("FORCE")

var tableUnsetColumnProjectionPolicy = g.NewQueryStruct("TableUnsetColumnProjectionPolicy").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("UNSET").
	SQL("PROJECTION POLICY")

var tableSetColumnTags = g.QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnTags, ssot.ClauseSet).
	WithSharedToOpts()

var tableUnsetColumnTags = g.QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnTags, ssot.ClauseUnset).
	WithSharedToOpts()

var tableColumnMaskingPolicy = g.NewQueryStruct("TableColumnMaskingPolicy").
	Identifier("MaskingPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	ListAssignment("USING", "Column", g.ParameterOptions().NoEquals().Parentheses()).
	WithValidation(g.ValidIdentifier, "MaskingPolicy").
	WithSharedToOpts()

var tableColumnProjectionPolicy = g.NewQueryStruct("TableColumnProjectionPolicy").
	Identifier("ProjectionPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PROJECTION POLICY").Required())

var tableDropColumnAction = g.QueryStructFromAttribute("Table", ssotdefs.TableDropColumnAction).
	WithSharedToOpts()

var tableRenameColumnAction = g.QueryStructFromAttribute("Table", ssotdefs.TableRenameColumnAction).
	WithSharedToOpts()

func tableSearchOptimizationAction() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableSearchOptimizationAction).WithSharedToOpts()
}

var (
	ReclusterStateEnumDef = ssotdefs.ReclusterState
	MatchTypeEnumDef      = ssotdefs.MatchType
)

func tableColumnsAndConstraints() *g.QueryStruct {
	return g.NewQueryStruct("CreateTableColumnsAndConstraints").
		ListQueryStructField("Columns", tableColumn(), g.KeywordOptions()).
		ListQueryStructField("OutOfLineConstraint", legacyTableOutOfLineConstraint(), g.ListOptions().NoParentheses())
}

func tableColumn() *g.QueryStruct {
	return g.NewQueryStruct("TableColumn").
		Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
		PredefinedQueryStructField("ColumnType", g.KindOfT[sdkcommons.DataType](), g.KeywordOptions().Required()).
		PredefinedQueryStructField("InlineConstraint", g.KindOfTPointer[sdkcommons.ColumnInlineConstraint](), g.KeywordOptions()).
		OptionalSQL("NOT NULL").
		OptionalTextAssignment("COLLATE", g.ParameterOptions().NoEquals().SingleQuotes()).
		OptionalQueryStructField("DefaultValue", tableColumnDefaultValue(), g.KeywordOptions()).
		OptionalSharedQueryStructField("MaskingPolicy", tableColumnMaskingPolicy, g.KeywordOptions()).
		OptionalTags().
		OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes())
}

func tableColumnDefaultValue() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableColumnDefaultValue)
}

// legacyTableOutOfLineConstraint emits `CONSTRAINT my_pk PRIMARY KEY (...)`.
// tableOutOfLineConstraint emits `CONSTRAINT "my_pk" PRIMARY KEY (...)`.
// snowflake_table_constraint and snowflake_table.primary_key would change: Snowflake stores the
// unquoted name folded to uppercase, and a quoted lowercase name in an existing config no longer
// matches that stored name. CREATE and ALTER ADD both keep this builder.
func legacyTableOutOfLineConstraint() *g.QueryStruct {
	return g.NewQueryStruct("OutOfLineConstraint").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals(), "Name").
		WithField(g.EnumLegacy[sdkcommons.ColumnConstraintType]("ConstraintType", g.KeywordOptions().Required())).
		PredefinedQueryStructField("Columns", "[]Column", g.KeywordOptions().Parentheses()).
		OptionalQueryStructField("ForeignKey", tableOutOfLineForeignKey(), g.KeywordOptions()).
		OptionalSQL("ENFORCED").
		OptionalSQL("NOT ENFORCED").
		OptionalSQL("DEFERRABLE").
		OptionalSQL("NOT DEFERRABLE").
		OptionalSQL("INITIALLY DEFERRED").
		OptionalSQL("INITIALLY IMMEDIATE").
		OptionalSQL("ENABLE").
		OptionalSQL("DISABLE").
		OptionalSQL("VALIDATE").
		OptionalSQL("NOVALIDATE").
		OptionalSQL("RELY").
		OptionalSQL("NORELY").
		WithValidation(g.ValidateValueSet, "Columns").
		WithValidation(g.ConflictingFields, "Enforced", "NotEnforced").
		WithValidation(g.ConflictingFields, "Deferrable", "NotDeferrable").
		WithValidation(g.ConflictingFields, "InitiallyDeferred", "InitiallyImmediate").
		WithValidation(g.ConflictingFields, "Enable", "Disable").
		WithValidation(g.ConflictingFields, "Validate", "Novalidate").
		WithValidation(g.ConflictingFields, "Rely", "Norely")
}

func tableOutOfLineForeignKey() *g.QueryStruct {
	return g.NewQueryStruct("OutOfLineForeignKey").
		SQL("REFERENCES").
		Identifier("TableName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("ColumnNames", "[]Column", g.ParameterOptions().NoEquals().Parentheses()).
		OptionalEnum("Match", MatchTypeEnumDef, g.ParameterOptions().NoEquals().SQL("MATCH")).
		PredefinedQueryStructField("On", g.KindOfTPointer[sdkcommons.ForeignKeyOnAction](), g.KeywordOptions())
}

func tableClusteringAction() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableClusteringAction)
}

func tableColumnAction() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableColumn,
		tableColumnInlineConstraint(),
		tableColumnMaskingPolicy,
		tableSetColumnMaskingPolicy,
		tableUnsetColumnMaskingPolicy,
		tableRenameColumnAction,
		tableSetColumnTags,
		tableUnsetColumnTags,
		tableDropColumnAction,
	)
}

func tableConstraintAction() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableConstraintAction, legacyTableOutOfLineConstraint())
}

func tableExternalTableAction() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableExternalTableAction).
		WithValidation(g.ValidateValueSet, "Add")
}

func tableSet() *g.QueryStruct {
	return g.AlterSetQueryStruct(ssotdefs.Table, ssot.SdkGenerator{})
}

func tableUnset() *g.QueryStruct {
	return g.AlterUnsetQueryStruct(ssotdefs.Table)
}

func tableAddRowAccessPolicy() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableAddRowAccessPolicy).WithSharedToOpts()
}

func tableDropRowAccessPolicy() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableDropRowAccessPolicy).WithSharedToOpts()
}

func tableDropAndAddRowAccessPolicy() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableDropAndAddRowAccessPolicy).WithSharedToOpts()
}

func tableAddStorageLifecyclePolicy() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableStorageLifecyclePolicy)
}

func tableSetAggregationPolicy() *g.QueryStruct {
	return g.QueryStructFromAttributeUse("Table", ssotdefs.TableAggregationPolicy, ssot.ClauseSet).WithSharedToOpts()
}

func tableUnsetAggregationPolicy() *g.QueryStruct {
	return g.QueryStructFromAttributeUse("Table", ssotdefs.TableAggregationPolicy, ssot.ClauseUnset).WithSharedToOpts()
}

func tableSetJoinPolicy() *g.QueryStruct {
	return g.QueryStructFromAttributeUse("Table", ssotdefs.TableJoinPolicy, ssot.ClauseSet).WithSharedToOpts()
}

func tableUnsetJoinPolicy() *g.QueryStruct {
	return g.QueryStructFromAttributeUse("Table", ssotdefs.TableJoinPolicy, ssot.ClauseUnset).WithSharedToOpts()
}

// DESCRIBE TABLE returns the collation glued onto the "type" column (e.g. "VARCHAR(200) COLLATE 'en-ci'"),
// so Type is converted manually and Collation is derived from it in additionalConvert().

// TODO [SNOW-1007542]: add missing features:
// - show columns (https://docs.snowflake.com/en/sql-reference/sql/show-columns)
// - show primary keys (https://docs.snowflake.com/en/sql-reference/sql/show-primary-keys)
// - truncate table (https://docs.snowflake.com/en/sql-reference/sql/truncate-table)
// - undrop table (https://docs.snowflake.com/en/sql-reference/sql/undrop-table)
var tablesDef = g.NewInterface(
	"Tables",
	"Table",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	// CREATE comes from ssotdefs.Table. The structs it names stay here because other operations reuse them.
	CreateOperationFromObject(ssotdefs.Table, tableColumnsAndConstraints()).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-table",
		g.NewQueryStruct("AlterTable").
			Alter().
			SQL("TABLE").
			IfExists().
			Name().
			RenameTo().
			OptionalIdentifier("SwapWith", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("SWAP WITH")).
			OptionalQueryStructField("ClusteringAction", tableClusteringAction(), g.KeywordOptions()).
			OptionalQueryStructField("ColumnAction", tableColumnAction(), g.KeywordOptions()).
			OptionalQueryStructField("ConstraintAction", tableConstraintAction(), g.KeywordOptions()).
			OptionalQueryStructField("ExternalTableAction", tableExternalTableAction(), g.KeywordOptions()).
			OptionalQueryStructField("SearchOptimizationAction", tableSearchOptimizationAction(), g.KeywordOptions()).
			OptionalQueryStructField("Set", tableSet(), g.KeywordOptions().SQL("SET")).
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalQueryStructField("Unset", tableUnset(), g.ListOptions().NoParentheses().SQL("UNSET")).
			OptionalQueryStructField("AddRowAccessPolicy", tableAddRowAccessPolicy(), g.KeywordOptions()).
			OptionalQueryStructField("DropRowAccessPolicy", tableDropRowAccessPolicy(), g.KeywordOptions()).
			OptionalQueryStructField("DropAndAddRowAccessPolicy", tableDropAndAddRowAccessPolicy(), g.ListOptions().NoParentheses()).
			OptionalSQL("DROP ALL ROW ACCESS POLICIES").
			OptionalQueryStructField("AddStorageLifecyclePolicy", tableAddStorageLifecyclePolicy(), g.KeywordOptions()).
			OptionalSQL("DROP STORAGE LIFECYCLE POLICY").
			OptionalQueryStructField("SetAggregationPolicy", tableSetAggregationPolicy(), g.KeywordOptions()).
			OptionalQueryStructField("UnsetAggregationPolicy", tableUnsetAggregationPolicy(), g.KeywordOptions()).
			OptionalQueryStructField("SetJoinPolicy", tableSetJoinPolicy(), g.KeywordOptions()).
			OptionalQueryStructField("UnsetJoinPolicy", tableUnsetJoinPolicy(), g.KeywordOptions()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifierIfSet, "RenameTo").
			WithValidation(g.ValidIdentifierIfSet, "SwapWith").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "SwapWith", "ClusteringAction", "ColumnAction", "ConstraintAction", "ExternalTableAction", "SearchOptimizationAction", "Set", "SetTags", "UnsetTags", "Unset", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "AddStorageLifecyclePolicy", "DropStorageLifecyclePolicy", "SetAggregationPolicy", "UnsetAggregationPolicy", "SetJoinPolicy", "UnsetJoinPolicy").
			WithAdditionalValidations(),
	).
	DropOperationFromObject(ssotdefs.Table).
	ShowOperationFromObject(ssotdefs.Table).
	// CREATE AS SELECT, USING TEMPLATE, LIKE, CLONE, describe, and select come from ssotdefs.Table.
	// ALTER helper structs stay handwritten, so SdkDefineFromSsot stays off.
	CustomOperationsFromObject(ssotdefs.Table).
	WithEnums(ssotdefs.TableConstraintType, ssotdefs.TableScope, ssotdefs.TableKind, ssotdefs.CloneMoment, ReclusterStateEnumDef, MatchTypeEnumDef, ssotdefs.TableSearchMethod)

func tableOutOfLineConstraint() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableOutOfLineConstraint)
}

// tableColumnInlineConstraint is used by iceberg tables (CREATE and ADD COLUMN) and by ALTER TABLE ADD COLUMN.
// Iceberg tables generate the struct and its toOpts, and ALTER TABLE ADD COLUMN calls that toOpts.
func tableColumnInlineConstraint() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableColumnInlineConstraint).WithSharedToOpts()
}
