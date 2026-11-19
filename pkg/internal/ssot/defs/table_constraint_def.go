package defs

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// MatchType is the MATCH clause of a foreign key.
var MatchType = ssot.NewEnum("MatchType", "MatchTypes", "FULL", "SIMPLE", "PARTIAL")

// The constraint structs follow https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint.
// Table emits them with SharedToOpts, and iceberg tables reference them. Every struct keeps its name with StructName.
//
// Both wrappers nest UNIQUE or PRIMARY KEY, FOREIGN KEY, and CHECK as children, so one call to
// QueryStructFromAttribute builds all of them.

func constraintName() ssot.Attribute {
	return ssot.Attribute{
		Name: "name", Type: ssot.AttrString, SqlOnly: true,
		CreateSql: ssot.StringParameter(ssot.Sql("CONSTRAINT"), ssot.NoEquals(), ssot.DoubleQuotes(), ssot.FieldName("Name")),
	}
}

func constraintFlag(name, sql string) ssot.Attribute {
	return ssot.Attribute{Name: name, Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql(sql))}
}

// constraintFlagPair is two keywords that cannot be set together, such as ENFORCED and NOT ENFORCED.
func constraintFlagPair(first, firstSql, second, secondSql string) []ssot.Attribute {
	other := constraintFlag(second, secondSql)
	flag := constraintFlag(first, firstSql)
	flag.ConflictsWith = []ssot.Attribute{other}
	return []ssot.Attribute{flag, other}
}

func constraintColumns(name, field string) ssot.Attribute {
	return ssot.Attribute{
		Name: name, Type: ssot.AttrString, SqlOnly: true,
		CreateSql: ssot.PredefinedType("[]Column", ssot.KeywordShape(), ssot.Parentheses(), ssot.FieldName(field)),
	}
}

func constraintReferences() ssot.Attribute {
	return ssot.Attribute{
		Name: "references", Type: ssot.AttrIdentifier, IdentifierType: "SchemaObjectIdentifier", Required: true, SqlOnly: true,
		CreateSql: ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("REFERENCES"), ssot.Required(), ssot.FieldName("References")),
	}
}

func constraintMatch() ssot.Attribute {
	return ssot.Attribute{
		Name: "match", Type: ssot.AttrEnum, Enum: MatchType, SqlOnly: true,
		CreateSql: ssot.PlainParameter(ssot.Sql("MATCH"), ssot.NoEquals(), ssot.FieldName("Match")),
	}
}

func constraintOn() ssot.Attribute {
	return ssot.Attribute{
		Name: "on", Type: ssot.AttrString, SqlOnly: true,
		CreateSql: ssot.PredefinedType("*ForeignKeyOnAction", ssot.KeywordShape(), ssot.FieldName("On")),
	}
}

// constraintProperties are the ENFORCED, DEFERRABLE, INITIALLY, ENABLE, VALIDATE, and RELY pairs that UNIQUE,
// PRIMARY KEY, and FOREIGN KEY constraints share. Each pair conflicts.
func constraintProperties() []ssot.Attribute {
	return slices.Concat(
		constraintFlagPair("enforced", "ENFORCED", "not_enforced", "NOT ENFORCED"),
		constraintFlagPair("deferrable", "DEFERRABLE", "not_deferrable", "NOT DEFERRABLE"),
		constraintFlagPair("initially_deferred", "INITIALLY DEFERRED", "initially_immediate", "INITIALLY IMMEDIATE"),
		constraintFlagPair("enable", "ENABLE", "disable", "DISABLE"),
		constraintFlagPair("validate", "VALIDATE", "novalidate", "NOVALIDATE"),
		constraintFlagPair("rely", "RELY", "norely", "NORELY"),
	)
}

func constraintComment() ssot.Attribute {
	return ssot.Attribute{
		Name: "comment", Type: ssot.AttrString, SqlOnly: true,
		CreateSql: ssot.StringParameter(ssot.Sql("COMMENT"), ssot.NoEquals()),
	}
}

func constraintCheck(structName, field string) ssot.Attribute {
	validate, novalidate := constraintFlag("enable_validate", "ENABLE VALIDATE"), constraintFlag("enable_novalidate", "ENABLE NOVALIDATE")
	validate.ConflictsWith = []ssot.Attribute{novalidate}
	return ssot.Attribute{
		Name:      "check",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: constraintClause(structName, field),
		Children: []ssot.Attribute{
			constraintName(),
			{Name: "check", SqlOnly: true, CreateSql: ssot.StaticSQL("check", "CHECK")},
			{Name: "open_paren", SqlOnly: true, CreateSql: ssot.StaticSQL("openParen", "(")},
			{Name: "expression", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.TextKeyword(ssot.NoQuotes(), ssot.Required(), ssot.FieldName("Expression"))},
			{Name: "close_paren", SqlOnly: true, CreateSql: ssot.StaticSQL("closeParen", ")")},
			validate,
			novalidate,
		},
	}
}

// constraintUniquePK is UNIQUE or PRIMARY KEY. Exactly one of the two keywords is set. The name, columns, properties, and
// comment sit beside that choice.
func constraintUniquePK(structName, field string, columns, comment bool) ssot.Attribute {
	children := []ssot.Attribute{
		constraintName(),
		{Name: "unique", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("UNIQUE"))},
		{Name: "primary_key", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("PRIMARY KEY"))},
	}
	if columns {
		children = append(children, constraintColumns("columns", "Columns"))
	}
	children = append(children, constraintProperties()...)
	if comment {
		children = append(children, constraintComment())
	}
	return ssot.Attribute{
		Name:               "unique_pk",
		Type:               ssot.AttrObject,
		SqlOnly:            true,
		ChildrenRule:       ssot.ChildrenExactlyOne,
		ChildrenRuleFields: []string{"unique", "primary_key"},
		CreateSql:          constraintClause(structName, field),
		Children:           children,
	}
}

func constraintClause(structName, field string) ssot.SqlClause {
	return ssot.ObjectParameter(ssot.StructName(structName), ssot.FieldName(field))
}

// TableColumnInlineUniquePK is an inline UNIQUE or PRIMARY KEY constraint.
var TableColumnInlineUniquePK = constraintUniquePK("TableColumnInlineUniquePK", "UniquePK", false, false)

// TableColumnInlineFK is an inline FOREIGN KEY constraint. The FOREIGN KEY keyword is optional inline.
var TableColumnInlineFK = ssot.Attribute{
	Name:      "fk",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: constraintClause("TableColumnInlineFK", "FK"),
	Children: append([]ssot.Attribute{
		constraintName(),
		constraintFlag("foreign_key", "FOREIGN KEY"),
		constraintReferences(),
		constraintColumns("ref_column", "RefColumn"),
		constraintMatch(),
		constraintOn(),
	}, constraintProperties()...),
}

// TableColumnInlineCH is an inline CHECK constraint.
var TableColumnInlineCH = constraintCheck("TableColumnInlineCH", "CH")

// TableColumnInlineConstraint is one inline constraint on a column: UNIQUE or PRIMARY KEY, FOREIGN KEY, or CHECK. The
// three kinds are its children, so the three structs are generated with it.
var TableColumnInlineConstraint = ssot.Attribute{
	Name:         "inline_constraint",
	Type:         ssot.AttrObject,
	SqlOnly:      true,
	ChildrenRule: ssot.ChildrenExactlyOne,
	CreateSql:    ssot.ObjectParameter(ssot.StructName("TableColumnInlineConstraint"), ssot.SharedToOpts()),
	Children: []ssot.Attribute{
		TableColumnInlineUniquePK,
		TableColumnInlineFK,
		TableColumnInlineCH,
	},
}

// TableOutOfLineUniquePK is an out-of-line UNIQUE or PRIMARY KEY constraint.
var TableOutOfLineUniquePK = constraintUniquePK("TableOutOfLineUniquePK", "UniquePK", true, true)

// TableOutOfLineFK is an out-of-line FOREIGN KEY constraint.
var TableOutOfLineFK = ssot.Attribute{
	Name:      "fk",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: constraintClause("TableOutOfLineFK", "FK"),
	Children: append(append([]ssot.Attribute{
		constraintName(),
		{Name: "foreign_key", SqlOnly: true, CreateSql: ssot.StaticSQL("foreignKey", "FOREIGN KEY")},
		constraintColumns("columns", "Columns"),
		constraintReferences(),
		constraintColumns("ref_columns", "RefColumns"),
		constraintMatch(),
		constraintOn(),
	}, constraintProperties()...), constraintComment()),
}

// TableOutOfLineCH is an out-of-line CHECK constraint.
var TableOutOfLineCH = constraintCheck("TableOutOfLineCH", "CH")

// TableOutOfLineConstraint is one out-of-line constraint: UNIQUE or PRIMARY KEY, FOREIGN KEY, or CHECK. The three kinds
// are its children, so the three structs are generated with it.
var TableOutOfLineConstraint = ssot.Attribute{
	Name:         "out_of_line_constraint",
	Type:         ssot.AttrObject,
	SqlOnly:      true,
	ChildrenRule: ssot.ChildrenExactlyOne,
	CreateSql:    ssot.ObjectParameter(ssot.StructName("TableOutOfLineConstraint"), ssot.SharedToOpts()),
	Children: []ssot.Attribute{
		TableOutOfLineUniquePK,
		TableOutOfLineFK,
		TableOutOfLineCH,
	},
}
