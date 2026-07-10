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
	OptionalText("ExpressionId", g.KeywordOptions()).
	WithValidation(g.ExactlyOneValueSet, "SearchMethodWithTarget", "ColumnName", "ExpressionId")

var tableDropSearchOptimization = g.NewQueryStruct("TableDropSearchOptimization").
	SQL("DROP SEARCH OPTIMIZATION").
	ListQueryStructField("On", tableDropSearchOptimizationOn, g.KeywordOptions().SQL("ON"))

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

var tableSearchOptimizationDetails = g.StructPair("tableSearchOptimizationDetailsRow", "TableSearchOptimizationDetails").
	Number("expression_id").
	Text("method").
	Text("target").
	DataType("target_data_type").
	BoolFromText("active", g.WithBoolTrueValue("true"))

var tableDescribeSearchOptimization = g.NewQueryStruct("DescribeSearchOptimization").
	Describe().
	SQL("SEARCH OPTIMIZATION").
	SQL("ON").
	Name().
	WithValidation(g.ValidIdentifier, "name")

var TableConstraintTypeDef = g.NewEnum(
	"TableConstraintType", "TableConstraintTypes",
	"PRIMARY KEY", "UNIQUE", "FOREIGN KEY",
)

var tableConstraintDetails = g.StructPair("tableConstraintDetailsRow", "TableConstraintDetails").
	Text("CONSTRAINT_CATALOG").
	Text("CONSTRAINT_SCHEMA").
	Text("CONSTRAINT_NAME").
	Text("TABLE_CATALOG").
	Text("TABLE_SCHEMA").
	Text("TABLE_NAME").
	Enum("CONSTRAINT_TYPE", TableConstraintTypeDef).
	BoolFromText("IS_DEFERRABLE", g.WithBoolTrueValue("YES")).
	BoolFromText("INITIALLY_DEFERRED", g.WithBoolTrueValue("YES")).
	OptionalText("COMMENT").
	Time("CREATED").
	Time("LAST_ALTERED").
	BoolFromText("ENFORCED", g.WithBoolTrueValue("YES")).
	BoolFromText("RELY", g.WithBoolTrueValue("YES"))

var tableShowConstraints = g.NewQueryStruct("ShowTableConstraints").
	SQLWithCustomFieldName("selectAll", "SELECT * FROM").
	Identifier("Database", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Required()).
	SQLWithCustomFieldName("dot", ".").
	SQLWithCustomFieldName("informationSchemaTableConstraints", "INFORMATION_SCHEMA.TABLE_CONSTRAINTS").
	SQLWithCustomFieldName("where", "WHERE").
	TextAssignment("TABLE_SCHEMA", g.ParameterOptions().SingleQuotes().Required()).
	SQLWithCustomFieldName("and", "AND").
	TextAssignment("TABLE_NAME", g.ParameterOptions().SingleQuotes().Required()).
	WithValidation(g.ValidIdentifier, "Database")

var tablesDef = g.NewInterface(
	"Tables",
	"Table",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CustomShowOperationWithPairedStructs(
		"DescribeSearchOptimization",
		g.ShowMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-search-optimization",
		tableSearchOptimizationDetails,
		tableDescribeSearchOptimization,
	).
	CustomShowOperationWithPairedStructs(
		"SelectTableConstraints",
		g.ShowMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/info-schema/table_constraints",
		tableConstraintDetails,
		tableShowConstraints,
	).
	WithEnums(TableConstraintTypeDef)

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

func tableOutOfLineUniquePK() *g.QueryStruct {
	return withOutOfLineConstraintTail(
		g.NewQueryStruct("TableOutOfLineUniquePK").
			OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
			OptionalSQL("UNIQUE").
			OptionalSQL("PRIMARY KEY").
			PredefinedQueryStructField("Columns", "[]Column", g.KeywordOptions().Parentheses()),
		// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
		// WithValidation(g.ExactlyOneValueSet, "Unique", "PrimaryKey")
	)
}

func tableOutOfLineFK() *g.QueryStruct {
	return withOutOfLineConstraintTail(
		g.NewQueryStruct("TableOutOfLineFK").
			OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
			SQL("FOREIGN KEY").
			PredefinedQueryStructField("Columns", "[]Column", g.KeywordOptions().Parentheses()).
			Identifier("References", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("REFERENCES").Required()).
			PredefinedQueryStructField("RefColumns", "[]Column", g.KeywordOptions().Parentheses()).
			PredefinedQueryStructField("Match", g.KindOfTPointer[sdkcommons.MatchType](), g.ParameterOptions().NoEquals().SQL("MATCH")).
			PredefinedQueryStructField("On", g.KindOfTPointer[sdkcommons.ForeignKeyOnAction](), g.KeywordOptions()),
		// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
		// WithValidation(g.ValidIdentifier, "References")
	)
}

func tableOutOfLineCH() *g.QueryStruct {
	return g.NewQueryStruct("TableOutOfLineCH").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
		SQL("CHECK").
		SQLWithCustomFieldName("openParen", "(").
		Text("Expression", g.KeywordOptions().NoQuotes().Required()).
		SQLWithCustomFieldName("closeParen", ")").
		OptionalSQL("ENABLE VALIDATE").
		OptionalSQL("ENABLE NOVALIDATE")
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// WithValidation(g.ConflictingFields, "EnableValidate", "EnableNovalidate")
}

func tableOutOfLineConstraint() *g.QueryStruct {
	return g.NewQueryStruct("TableOutOfLineConstraint").
		OptionalQueryStructField("UniquePK", tableOutOfLineUniquePK(), g.KeywordOptions()).
		OptionalQueryStructField("FK", tableOutOfLineFK(), g.KeywordOptions()).
		OptionalQueryStructField("CH", tableOutOfLineCH(), g.KeywordOptions())
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// WithValidation(g.ExactlyOneValueSet, "UniquePK", "FK", "CH")
}

// withOutOfLineConstraintTail appends the tail clauses shared by out-of-line UNIQUE/PK and FK
// constraints (ENFORCED / DEFERRABLE / INITIALLY / ENABLE / VALIDATE / RELY pairs plus COMMENT)
// and their ConflictingFields validations. Applied as a wrapper so tail fields are emitted
// after the caller's struct-specific fields. Out-of-line CHECK uses its own ENABLE VALIDATE
// pair and is not wrapped.
func withOutOfLineConstraintTail(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
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
		OptionalTextAssignment("COMMENT", g.ParameterOptions().NoEquals().SingleQuotes())
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// WithValidation(g.ConflictingFields, "Enforced", "NotEnforced").
	// WithValidation(g.ConflictingFields, "Deferrable", "NotDeferrable").
	// WithValidation(g.ConflictingFields, "InitiallyDeferred", "InitiallyImmediate").
	// WithValidation(g.ConflictingFields, "Enable", "Disable").
	// WithValidation(g.ConflictingFields, "Validate", "Novalidate").
	// WithValidation(g.ConflictingFields, "Rely", "Norely")
}

func tableColumnInlineConstraint() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnInlineConstraint").
		OptionalQueryStructField("UniquePK", tableColumnInlineUniquePK(), g.KeywordOptions()).
		OptionalQueryStructField("FK", tableColumnInlineFK(), g.KeywordOptions()).
		OptionalQueryStructField("CH", tableColumnInlineCH(), g.KeywordOptions())
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// WithValidation(g.ExactlyOneValueSet, "UniquePK", "FK", "CH")
}

func tableColumnInlineUniquePK() *g.QueryStruct {
	return withInlineConstraintTail(
		g.NewQueryStruct("TableColumnInlineUniquePK").
			OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
			OptionalSQL("UNIQUE").
			OptionalSQL("PRIMARY KEY"),
		// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
		// WithValidation(g.ExactlyOneValueSet, "Unique", "PrimaryKey")
	)
}

func tableColumnInlineFK() *g.QueryStruct {
	return withInlineConstraintTail(
		g.NewQueryStruct("TableColumnInlineFK").
			OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
			OptionalSQL("FOREIGN KEY").
			Identifier("References", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("REFERENCES").Required()).
			PredefinedQueryStructField("RefColumn", "[]Column", g.KeywordOptions().Parentheses()).
			PredefinedQueryStructField("Match", g.KindOfTPointer[sdkcommons.MatchType](), g.ParameterOptions().NoEquals().SQL("MATCH")).
			PredefinedQueryStructField("On", g.KindOfTPointer[sdkcommons.ForeignKeyOnAction](), g.KeywordOptions()),
		// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
		// WithValidation(g.ValidIdentifier, "References")
	)
}

func tableColumnInlineCH() *g.QueryStruct {
	return g.NewQueryStruct("TableColumnInlineCH").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", g.ParameterOptions().NoEquals().DoubleQuotes(), "Name").
		SQL("CHECK").
		SQLWithCustomFieldName("openParen", "(").
		Text("Expression", g.KeywordOptions().NoQuotes().Required()).
		SQLWithCustomFieldName("closeParen", ")").
		OptionalSQL("ENABLE VALIDATE").
		OptionalSQL("ENABLE NOVALIDATE")
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// WithValidation(g.ConflictingFields, "EnableValidate", "EnableNovalidate")
}

// withInlineConstraintTail appends the tail clauses shared by inline UNIQUE/PK and FK
// constraints (ENFORCED / DEFERRABLE / INITIALLY / ENABLE / VALIDATE / RELY pairs) and
// their ConflictingFields validations. Applied as a wrapper so tail fields are emitted
// after the caller's struct-specific fields.
func withInlineConstraintTail(qs *g.QueryStruct) *g.QueryStruct {
	return qs.
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
		OptionalSQL("NORELY")
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// WithValidation(g.ConflictingFields, "Enforced", "NotEnforced").
	// WithValidation(g.ConflictingFields, "Deferrable", "NotDeferrable").
	// WithValidation(g.ConflictingFields, "InitiallyDeferred", "InitiallyImmediate").
	// WithValidation(g.ConflictingFields, "Enable", "Disable").
	// WithValidation(g.ConflictingFields, "Validate", "Novalidate").
	// WithValidation(g.ConflictingFields, "Rely", "Norely")
}
