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
		ListSharedQueryStructField("OutOfLineConstraint", tableOutOfLineConstraint(), g.KeywordOptions())
}

func tableColumn() *g.QueryStruct {
	return g.NewQueryStruct("TableColumn").
		Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
		PredefinedQueryStructField("ColumnType", g.KindOfT[sdkcommons.DataType](), g.KeywordOptions().Required()).
		OptionalSharedQueryStructField("InlineConstraint", tableColumnInlineConstraint(), g.KeywordOptions()).
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

func tableAsSelectColumn() *g.QueryStruct {
	return g.NewQueryStruct("TableAsSelectColumn").
		Text("Name", g.KeywordOptions().Required()).
		PredefinedQueryStructField("ColumnType", g.KindOfTPointer[sdkcommons.DataType](), g.KeywordOptions()).
		OptionalIdentifier("MaskingPolicy", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY"))
}

func tableClonePoint() *g.QueryStruct {
	return g.NewQueryStruct("ClonePoint").
		Enum("Moment", ssotdefs.CloneMoment, g.ParameterOptions().NoEquals().Required()).
		PredefinedQueryStructField("At", "TimeTravel", g.ListOptions().Parentheses().NoComma().Required())
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
	return g.QueryStructFromAttribute("Table", ssotdefs.TableConstraintAction, tableOutOfLineConstraint())
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
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "SwapWith", "ClusteringAction", "ColumnAction", "ConstraintAction", "ExternalTableAction", "SearchOptimizationAction", "Set", "SetTags", "UnsetTags", "Unset", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "AddStorageLifecyclePolicy", "DropStorageLifecyclePolicy", "SetAggregationPolicy", "UnsetAggregationPolicy", "SetJoinPolicy", "UnsetJoinPolicy"),
	).
	DropOperationFromObject(ssotdefs.Table).
	ShowOperationFromObject(ssotdefs.Table).
	// CREATE AS SELECT, USING TEMPLATE, LIKE, CLONE, describe, and select come from ssotdefs.Table.
	// ALTER helper structs stay handwritten, so SdkDefineFromSsot stays off.
	CustomOperationsFromObject(ssotdefs.Table).
	ShowParametersFromObject(ssotdefs.Table).
	WithEnums(ssotdefs.TableConstraintType, ssotdefs.TableScope, ssotdefs.TableKind, ssotdefs.CloneMoment, ReclusterStateEnumDef, MatchTypeEnumDef, ssotdefs.TableSearchMethod)

// tableOutOfLineConstraint is used by iceberg tables (CREATE) and by CREATE TABLE and ALTER TABLE ADD.
// Iceberg tables generate the struct and its toOpts, and the table operations call that toOpts.
func tableOutOfLineConstraint() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableOutOfLineConstraint).WithSharedToOpts()
}

// tableColumnInlineConstraint is used by iceberg tables (CREATE and ADD COLUMN) and by ALTER TABLE ADD COLUMN.
// Iceberg tables generate the struct and its toOpts, and ALTER TABLE ADD COLUMN calls that toOpts.
func tableColumnInlineConstraint() *g.QueryStruct {
	return g.QueryStructFromAttribute("Table", ssotdefs.TableColumnInlineConstraint).WithSharedToOpts()
}
