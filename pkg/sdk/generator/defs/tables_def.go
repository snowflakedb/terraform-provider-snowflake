package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var tableSetColumnMaskingPolicy = g.NewQueryStruct("TableSetColumnMaskingPolicy").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("SET").
	Identifier("MaskingPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	ListAssignment("USING", "Column", g.ParameterOptions().NoEquals().Parentheses()).
	OptionalSQL("FORCE")

var tableUnsetColumnMaskingPolicy = g.NewQueryStruct("TableUnsetColumnMaskingPolicy").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SQL("UNSET").
	SQL("MASKING POLICY")

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

var tableSetColumnTags = g.NewQueryStruct("TableSetColumnTags").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	SetTags()

var tableUnsetColumnTags = g.NewQueryStruct("TableUnsetColumnTags").
	SQL("ALTER COLUMN").
	Text("Name", g.KeywordOptions().Required().DoubleQuotes()).
	UnsetTags()

var TableSearchMethodEnumDef = g.NewEnum(
	"TableSearchMethod", "TableSearchMethods",
	"SUBSTRING", "EQUALITY", "FULL_TEXT",
)

var tableColumnMaskingPolicy = g.NewQueryStruct("TableColumnMaskingPolicy").
	Identifier("MaskingPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("MASKING POLICY").Required()).
	ListAssignment("USING", "Column", g.ParameterOptions().NoEquals().Parentheses())

var tableColumnProjectionPolicy = g.NewQueryStruct("TableColumnProjectionPolicy").
	Identifier("ProjectionPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PROJECTION POLICY").Required())

var tableDropColumnAction = g.NewQueryStruct("TableDropColumnAction").
	SQL("DROP COLUMN").
	OptionalSQL("IF EXISTS").
	PredefinedQueryStructField("Columns", "[]Column", g.KeywordOptions().Required())

var tableRenameColumnAction = g.NewQueryStruct("TableRenameColumnAction").
	SQL("RENAME COLUMN").
	Text("OldName", g.KeywordOptions().Required().DoubleQuotes()).
	AssignmentWithFieldName("TO", "string", g.ParameterOptions().NoEquals().DoubleQuotes(), "NewName")

func newTableSearchMethodArgs() *g.QueryStruct {
	return g.NewQueryStruct("TableSearchMethodArgs").
		PredefinedQueryStructField("Targets", "[]string", g.KeywordOptions()).
		OptionalTextAssignment("ANALYZER", g.ParameterOptions().ArrowEquals().SingleQuotes())
}

func newTableSearchMethodWithTarget() *g.QueryStruct {
	return g.NewQueryStruct("TableSearchMethodWithTarget").
		PredefinedQueryStructField("Method", TableSearchMethodEnumDef.Kind(), g.KeywordOptions().Required()).
		QueryStructField("Args", newTableSearchMethodArgs(), g.ListOptions().Parentheses())
}

var tableAddSearchOptimization = g.NewQueryStruct("TableAddSearchOptimization").
	SQL("ADD SEARCH OPTIMIZATION").
	ListQueryStructField("On", newTableSearchMethodWithTarget(), g.KeywordOptions().SQL("ON"))

var tableDropSearchOptimizationOn = g.NewQueryStruct("TableDropSearchOptimizationOn").
	OptionalQueryStructField("SearchMethodWithTarget", newTableSearchMethodWithTarget(), g.KeywordOptions()).
	OptionalText("ColumnName", g.KeywordOptions()).
	OptionalText("ExpressionId", g.KeywordOptions())

// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
// .WithValidation(g.ExactlyOneValueSet, "SearchMethodWithTarget", "ColumnName", "ExpressionId")

var tableDropSearchOptimization = g.NewQueryStruct("TableDropSearchOptimization").
	SQL("DROP SEARCH OPTIMIZATION").
	ListQueryStructField("On", tableDropSearchOptimizationOn, g.KeywordOptions().SQL("ON")).
	WithAdditionalValidations()

var tableSearchOptimizationAction = g.NewQueryStruct("TableSearchOptimizationAction").
	OptionalQueryStructField(
		"Add",
		tableAddSearchOptimization,
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Drop",
		tableDropSearchOptimization,
		g.KeywordOptions(),
	).
	WithValidation(g.ExactlyOneValueSet, "Add", "Drop")

var tableSetAggregationPolicy = g.NewQueryStruct("TableSetAggregationPolicy").
	SQL("SET").
	Identifier("AggregationPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("AGGREGATION POLICY").Required()).
	ListAssignment("ENTITY KEY", "Column", g.ParameterOptions().NoEquals().Parentheses()).
	OptionalSQL("FORCE").
	WithValidation(g.ValidIdentifier, "AggregationPolicy")

var tableUnsetAggregationPolicy = g.NewQueryStruct("TableUnsetAggregationPolicy").
	SQL("UNSET AGGREGATION POLICY")

var tableSetJoinPolicy = g.NewQueryStruct("TableSetJoinPolicy").
	SQL("SET").
	Identifier("JoinPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("JOIN POLICY").Required()).
	OptionalSQL("FORCE").
	WithValidation(g.ValidIdentifier, "JoinPolicy")

var tableUnsetJoinPolicy = g.NewQueryStruct("TableUnsetJoinPolicy").
	SQL("UNSET JOIN POLICY")
