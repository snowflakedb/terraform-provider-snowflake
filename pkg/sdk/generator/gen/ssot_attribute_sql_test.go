package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_appendAttributeSql(t *testing.T) {
	create, set, unset := sqlClauseSequences()

	t.Run("create", func(t *testing.T) {
		got := NewQueryStruct("clauses").appendSqlSequence(create, attributeSqlCreate)
		want := NewQueryStruct("clauses").
			NumberAssignment("COUNT", ParameterOptions().Required()).
			OptionalNumberAssignment("LIMIT", ParameterOptions().NoEquals()).
			BooleanAssignment("ENABLED", ParameterOptions().Required()).
			OptionalBooleanAssignment("AUTO", ParameterOptions()).
			TextAssignment("LABEL", ParameterOptions().SingleQuotes()).
			OptionalTags().
			OptionalTextAssignment("COMMENT", ParameterOptions().NoQuotes()).
			EnumAssignment("FAMILY", sqlClauseFamily, ParameterOptions().NoQuotes().Required()).
			OptionalEnumAssignment("STATUS", sqlClauseStatus, ParameterOptions().NoQuotes()).
			ListAssignment("ITEMS", "WidgetItem", ParameterOptions().Parentheses()).
			OptionalIdentifier("ForApplication", "AccountObjectIdentifier", IdentifierOptions().SQL("FOR APPLICATION")).
			Identifier("Parent", "DatabaseObjectIdentifier", IdentifierOptions().SingleQuotes().NoEquals().Required().SQL("PARENT")).
			OptionalSQL("TRANSIENT")
		require.Equal(t, want.fields, got.fields)
	})

	t.Run("set", func(t *testing.T) {
		got := NewQueryStruct("clauses").appendSqlSequence(set, attributeSqlSet)
		want := NewQueryStruct("clauses").
			OptionalNumberAssignment("COUNT", ParameterOptions()).
			NumberAssignment("LIMIT", ParameterOptions().Required()).
			OptionalBooleanAssignment("AUTO", ParameterOptions()).
			OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes()).
			OptionalEnumAssignment("STATUS", sqlClauseStatus, ParameterOptions().NoQuotes()).
			ListAssignment("ITEMS", "WidgetItem", ParameterOptions().NoParentheses())
		require.Equal(t, want.fields, got.fields)
	})

	t.Run("unset", func(t *testing.T) {
		got := NewQueryStruct("clauses").appendSqlSequence(unset, attributeSqlUnset)
		want := NewQueryStruct("clauses").
			OptionalSQL("LIMIT").
			OptionalSQL("AUTO RESUME").
			OptionalSQL("COMMENT").
			OptionalSQL("STATUS").
			OptionalSQL("ITEMS")
		require.Equal(t, want.fields, got.fields)
	})
}

func Test_appendAttributeSql_identifierEquals(t *testing.T) {
	monitor := ssot.Attribute{
		Name:      "resource_monitor",
		Type:      ssot.AttrIdentifier,
		CreateSql: ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Equals()),
	}
	got := NewQueryStruct("clauses").appendSqlSequence([]*ssot.Attribute{&monitor}, attributeSqlCreate)
	want := NewQueryStruct("clauses").
		OptionalIdentifier("ResourceMonitor", "AccountObjectIdentifier", IdentifierOptions().SQL("RESOURCE_MONITOR").Equals())
	require.Equal(t, want.fields, got.fields)
}

func Test_appendAttributeSql_showEnum(t *testing.T) {
	resourceEnum := ssot.NewEnum("WidgetTypeUserSettable", "WidgetTypesUserSettable", "STANDARD")
	showEnum := ssot.NewEnum("WidgetType", "WidgetTypes", "STANDARD", "ADAPTIVE")
	kind := ssot.Attribute{
		Name:      "kind",
		Type:      ssot.AttrEnum,
		Enum:      resourceEnum,
		ShowEnum:  showEnum,
		CreateSql: ssot.EnumParameter(ssot.SingleQuotes()),
	}
	got := NewQueryStruct("clauses").appendSqlSequence([]*ssot.Attribute{&kind}, attributeSqlCreate)
	want := NewQueryStruct("clauses").
		OptionalEnumAssignment("KIND", showEnum, ParameterOptions().SingleQuotes())
	require.Equal(t, want.fields, got.fields)
}

func Test_appendSqlSequence_order(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString, CreateSql: ssot.StringParameter()}
	count := ssot.Attribute{Name: "count", Type: ssot.AttrInt, Required: true, CreateSql: ssot.IntParameter()}
	seq := []*ssot.Attribute{&count, ssot.Tags, &comment}
	got := NewQueryStruct("clauses").appendSqlSequence(seq, attributeSqlCreate)
	want := NewQueryStruct("clauses").
		NumberAssignment("COUNT", ParameterOptions().Required()).
		OptionalTags().
		OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes())
	require.Equal(t, want.fields, got.fields)
}

func Test_appendAttributeSql_panics(t *testing.T) {
	t.Run("parameter type", func(t *testing.T) {
		require.Panics(t, func() {
			NewQueryStruct("clauses").appendSqlSequence(sqlSequence(ssot.Attribute{
				Name:      "created_on",
				Type:      ssot.AttrTime,
				CreateSql: ssot.IntParameter(),
			}), attributeSqlCreate)
		})
	})
	t.Run("enum without enum", func(t *testing.T) {
		require.Panics(t, func() {
			NewQueryStruct("clauses").appendSqlSequence(sqlSequence(ssot.Attribute{
				Name:      "family",
				Type:      ssot.AttrEnum,
				CreateSql: ssot.EnumParameter(),
			}), attributeSqlCreate)
		})
	})
	t.Run("list without item struct", func(t *testing.T) {
		require.Panics(t, func() {
			NewQueryStruct("clauses").appendSqlSequence(sqlSequence(ssot.Attribute{
				Name:      "items",
				Type:      ssot.AttrList,
				CreateSql: ssot.ListParameter(""),
			}), attributeSqlCreate)
		})
	})
	t.Run("required keyword", func(t *testing.T) {
		require.Panics(t, func() {
			NewQueryStruct("clauses").appendSqlSequence(sqlSequence(ssot.Attribute{
				Name:      "flag",
				Type:      ssot.AttrBool,
				CreateSql: ssot.BoolKeyword(ssot.Required()),
			}), attributeSqlCreate)
		})
	})
	t.Run("keyword tags", func(t *testing.T) {
		require.Panics(t, func() {
			NewQueryStruct("clauses").appendSqlSequence(sqlSequence(ssot.Attribute{
				Name:      "flag",
				Type:      ssot.AttrBool,
				CreateSql: ssot.BoolKeyword(ssot.SingleQuotes()),
			}), attributeSqlCreate)
		})
	})
	t.Run("identifier without type", func(t *testing.T) {
		require.Panics(t, func() {
			NewQueryStruct("clauses").appendSqlSequence(sqlSequence(ssot.Attribute{
				Name:      "parent",
				Type:      ssot.AttrIdentifier,
				CreateSql: ssot.IdentifierKeyword(""),
			}), attributeSqlCreate)
		})
	})
	t.Run("identifier parentheses", func(t *testing.T) {
		require.Panics(t, func() {
			NewQueryStruct("clauses").appendSqlSequence(sqlSequence(ssot.Attribute{
				Name:      "parent",
				Type:      ssot.AttrIdentifier,
				CreateSql: ssot.IdentifierKeyword("DatabaseObjectIdentifier", ssot.Parentheses()),
			}), attributeSqlCreate)
		})
	})
	t.Run("unknown kind", func(t *testing.T) {
		require.Panics(t, func() {
			NewQueryStruct("clauses").appendSqlSequence(sqlSequence(ssot.Attribute{
				Name:      "count",
				Type:      ssot.AttrInt,
				CreateSql: ssot.SqlClause{Present: true, Kind: "other"},
			}), attributeSqlCreate)
		})
	})
	t.Run("unknown statement", func(t *testing.T) {
		require.Panics(t, func() {
			NewQueryStruct("clauses").appendSqlSequence(sqlSequence(ssot.Attribute{
				Name:      "count",
				Type:      ssot.AttrInt,
				CreateSql: ssot.IntParameter(),
			}), attributeSqlStatement(9))
		})
	})
}

var (
	sqlClauseFamily = ssot.NewEnum("WidgetFamily", "WidgetFamilies", "SMALL")
	sqlClauseStatus = ssot.NewEnum("WidgetStatus", "WidgetStatuses", "ON")
)

func sqlSequence(attr ssot.Attribute) []*ssot.Attribute {
	return []*ssot.Attribute{&attr}
}

// sqlClauseSequences is a fake statement order. Each clause shape the adapter knows is present once.
// note has no clause and is left out of every list. Tags sits between label and comment.
func sqlClauseSequences() (create, set, unset []*ssot.Attribute) {
	count := ssot.Attribute{Name: "count", Type: ssot.AttrInt, Required: true, CreateSql: ssot.IntParameter(), SetSql: ssot.IntParameter()}
	limit := ssot.Attribute{Name: "limit", Type: ssot.AttrInt, CreateSql: ssot.IntParameter(ssot.NoEquals()), SetSql: ssot.IntParameter(ssot.Required()), UnsetSql: ssot.BoolKeyword()}
	enabled := ssot.Attribute{Name: "enabled", Type: ssot.AttrBool, Required: true, CreateSql: ssot.BoolParameter()}
	auto := ssot.Attribute{Name: "auto", Type: ssot.AttrBool, CreateSql: ssot.BoolParameter(), SetSql: ssot.BoolParameter(), UnsetSql: ssot.BoolKeyword(ssot.Sql("AUTO RESUME"))}
	label := ssot.Attribute{Name: "label", Type: ssot.AttrString, Required: true, CreateSql: ssot.StringParameter()}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString, CreateSql: ssot.StringParameter(ssot.NoQuotes()), SetSql: ssot.StringParameter(), UnsetSql: ssot.BoolKeyword()}
	family := ssot.Attribute{Name: "family", Type: ssot.AttrEnum, Required: true, Enum: sqlClauseFamily, CreateSql: ssot.EnumParameter()}
	status := ssot.Attribute{Name: "status", Type: ssot.AttrEnum, Enum: sqlClauseStatus, CreateSql: ssot.EnumParameter(), SetSql: ssot.EnumParameter(), UnsetSql: ssot.BoolKeyword()}
	items := ssot.Attribute{Name: "items", Type: ssot.AttrList, CreateSql: ssot.ListParameter("WidgetItem"), SetSql: ssot.ListParameter("WidgetItem", ssot.NoParentheses()), UnsetSql: ssot.BoolKeyword()}
	target := ssot.Attribute{Name: "target", Type: ssot.AttrIdentifier, CreateSql: ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Sql("FOR APPLICATION"))}
	parent := ssot.Attribute{Name: "parent", Type: ssot.AttrIdentifier, Required: true, CreateSql: ssot.IdentifierKeyword("DatabaseObjectIdentifier", ssot.SingleQuotes(), ssot.NoEquals())}
	flag := ssot.Attribute{Name: "flag", Type: ssot.AttrBool, CreateSql: ssot.BoolKeyword(ssot.Sql("TRANSIENT"))}
	create = []*ssot.Attribute{
		&count, &limit, &enabled, &auto, &label, ssot.Tags, &comment, &family, &status, &items, &target, &parent, &flag,
	}
	set = []*ssot.Attribute{&count, &limit, &auto, &comment, &status, &items}
	unset = []*ssot.Attribute{&limit, &auto, &comment, &status, &items}
	return create, set, unset
}

func Test_appendSqlSequence_typeDefaults(t *testing.T) {
	count := ssot.Attribute{Name: "count", Type: ssot.AttrInt, Required: true}
	label := ssot.Attribute{Name: "label", Type: ssot.AttrString, CreateSql: ssot.Clause(ssot.NoQuotes())}
	flag := ssot.Attribute{Name: "flag", Type: ssot.AttrBool, CreateSql: ssot.BoolKeyword(ssot.Sql("TRANSIENT"))}
	limit := ssot.Attribute{Name: "limit", Type: ssot.AttrInt, UnsetSql: ssot.IntParameter()}

	got := NewQueryStruct("clauses").appendSqlSequence([]*ssot.Attribute{&count, &label, &flag}, attributeSqlCreate)
	want := NewQueryStruct("clauses").
		NumberAssignment("COUNT", ParameterOptions().Required()).
		OptionalTextAssignment("LABEL", ParameterOptions().NoQuotes()).
		OptionalSQL("TRANSIENT")
	require.Equal(t, want.fields, got.fields)

	got = NewQueryStruct("clauses").appendSqlSequence([]*ssot.Attribute{&count, &limit}, attributeSqlUnset)
	want = NewQueryStruct("clauses").
		OptionalSQL("COUNT").
		OptionalNumberAssignment("LIMIT", ParameterOptions())
	require.Equal(t, want.fields, got.fields)
}

func Test_appendSqlSequence_static(t *testing.T) {
	kind := ssot.Attribute{
		Name:      "kind",
		SqlOnly:   true,
		CreateSql: ssot.StaticSQL("kind", "KIND = 'MINI'"),
	}
	got := NewQueryStruct("clauses").appendSqlSequence([]*ssot.Attribute{&kind}, attributeSqlCreate)
	want := NewQueryStruct("clauses").SQLWithCustomFieldName("kind", "KIND = 'MINI'")
	require.Equal(t, want.fields, got.fields)
}
