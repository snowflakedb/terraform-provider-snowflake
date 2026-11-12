package ssot

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/stretchr/testify/require"
)

func Test_SqlClause_zeroIsAbsent(t *testing.T) {
	var clause SqlClause
	require.False(t, clause.Present)
}

func Test_SqlClause_intParameter(t *testing.T) {
	attr := Attribute{Name: "min_nodes", Type: AttrInt, Required: true}
	clause := IntParameter()

	require.True(t, clause.Present)
	require.Equal(t, SqlKindParameter, clause.Kind)
	require.Nil(t, clause.Required)
	require.Equal(t, "parameter", clause.Ddl())
	require.Equal(t, "MIN_NODES", clause.SQLText(attr))
	require.Empty(t, clause.Quotes)
	require.Empty(t, clause.Parentheses)
	require.Empty(t, clause.Equals)
}

func Test_SqlClause_stringParameter(t *testing.T) {
	clause := StringParameter()

	require.Equal(t, "parameter,single_quotes", clause.Ddl())
	require.Equal(t, SqlQuotesSingle, clause.Quotes)
}

func Test_SqlClause_enumParameter(t *testing.T) {
	clause := EnumParameter()

	require.Equal(t, "parameter,no_quotes", clause.Ddl())
	require.Equal(t, SqlQuotesNo, clause.Quotes)
}

func Test_SqlClause_boolKeyword(t *testing.T) {
	attr := Attribute{Name: "auto_resume", Type: AttrBool}
	clause := BoolKeyword()

	require.Equal(t, SqlKindKeyword, clause.Kind)
	require.Equal(t, "keyword", clause.Ddl())
	require.Equal(t, "AUTO_RESUME", clause.SQLText(attr))
}

func Test_SqlClause_listParameter(t *testing.T) {
	clause := ListParameter("PoolBackupItem")

	require.Equal(t, "PoolBackupItem", clause.ItemStruct)
	require.Equal(t, "parameter,parentheses", clause.Ddl())
}

func Test_SqlClause_listParameterMustParentheses(t *testing.T) {
	clause := ListParameter("StringListItemWrapper", MustParentheses())

	require.Equal(t, SqlParensMust, clause.Parentheses)
	require.Equal(t, "parameter,must_parentheses", clause.Ddl())
}

func Test_SqlClause_reverseList(t *testing.T) {
	clause := ListParameter("AccountObjectIdentifier", NoParentheses(), Reverse(), Sql("TO ALLOWED_DATABASES"))

	require.True(t, clause.Reverse)
	require.Equal(t, "parameter,no_parentheses,reverse", clause.Ddl())
	require.Equal(t, "TO ALLOWED_DATABASES", clause.SQLText(Attribute{Name: "allowed_databases"}))
}

func Test_SqlClause_identifierKeyword(t *testing.T) {
	attr := Attribute{Name: "for_application", Type: AttrIdentifier}
	clause := IdentifierKeyword("AccountObjectIdentifier", Sql("FOR APPLICATION"))

	require.Equal(t, "AccountObjectIdentifier", clause.IdType)
	require.Equal(t, "identifier", clause.Ddl())
	require.Equal(t, "FOR APPLICATION", clause.SQLText(attr))
}

func Test_EffectiveSqlClause_defaults(t *testing.T) {
	count := Attribute{Name: "count", Type: AttrInt, Required: true}
	require.Equal(t, IntParameter(), count.EffectiveSqlClause(ClauseCreate))
	require.Equal(t, IntParameter(), count.EffectiveSqlClause(ClauseSet))
	require.Equal(t, BoolKeyword(), count.EffectiveSqlClause(ClauseUnset))

	label := Attribute{Name: "label", Type: AttrString}
	require.Equal(t, StringParameter(), label.EffectiveSqlClause(ClauseCreate))
	label.CreateSql = Clause(NoQuotes())
	require.Equal(t, StringParameter(NoQuotes()), label.EffectiveSqlClause(ClauseCreate))

	family := Attribute{Name: "family", Type: AttrEnum}
	require.Equal(t, EnumParameter(), family.EffectiveSqlClause(ClauseCreate))

	flag := Attribute{Name: "flag", Type: AttrBool, CreateSql: BoolKeyword(Sql("TRANSIENT"))}
	require.Equal(t, BoolKeyword(Sql("TRANSIENT")), flag.EffectiveSqlClause(ClauseCreate))

	limit := Attribute{Name: "limit", Type: AttrInt, UnsetSql: IntParameter()}
	require.Equal(t, IntParameter(), limit.EffectiveSqlClause(ClauseUnset))
	require.False(t, count.UnsetSql.Present)
}

func Test_EnumListItemName(t *testing.T) {
	family := NewEnum("WidgetFamily", "WidgetFamilies", "SMALL")
	items := Attribute{
		Name: "items",
		Type: AttrList,
		Elem: &Attribute{Type: AttrEnum, Enum: family},
	}
	require.Equal(t, "WidgetFamilyListItem", EnumListItemName(family))
	require.Equal(t, ListParameter("WidgetFamilyListItem"), items.EffectiveSqlClause(ClauseCreate))
	require.Equal(t, ListParameter("WidgetFamilyListItem"), items.EffectiveSqlClause(ClauseSet))

	set := Attribute{
		Name: "families",
		Type: AttrSet,
		Elem: &Attribute{Type: AttrEnum, Enum: family},
	}
	require.Equal(t, ListParameter("WidgetFamilyListItem"), set.EffectiveSqlClause(ClauseCreate))
}

func Test_SqlClause_arrowEqualsAndValueKeyword(t *testing.T) {
	analyzer := StringParameter(ArrowEquals())
	require.Equal(t, SqlEqualsArrow, analyzer.Equals)
	require.Equal(t, "parameter,single_quotes,arrow_equals", analyzer.Ddl())

	on := ValueKeyword(Sql("ON"), FieldName("On"))
	require.True(t, on.ValueKeyword)
	require.Equal(t, SqlKindKeyword, on.Kind)
	require.Equal(t, "ON", on.SQL)
	require.Equal(t, "On", on.Field)

	args := ObjectParameter(Parentheses(), Required(), Commas())
	require.False(t, args.NoComma)
	require.True(t, ObjectParameter().NoComma)
	require.Equal(t, SqlParensYes, args.Parentheses)
}

func Test_SqlClause_optionOverridesDefault(t *testing.T) {
	clause := StringParameter(NoQuotes())

	require.Equal(t, "parameter,no_quotes", clause.Ddl())
	require.Equal(t, SqlQuotesNo, clause.Quotes)
}

func Test_NestedSqlFieldsFor_createSetUnset(t *testing.T) {
	shared := IdentifierKeyword("SchemaObjectIdentifier", Sql("AGGREGATION POLICY"))
	attr := Attribute{
		Name:      "aggregation_policy",
		Type:      AttrObject,
		CreateSql: ObjectParameter(),
		SetSql:    ObjectParameter(FieldName("SetAggregationPolicy")),
		UnsetSql:  ObjectParameter(FieldName("UnsetAggregationPolicy")),
		Children: []Attribute{
			{Name: "set", SqlOnly: true, SetSql: StaticSQL("set", "SET")},
			{Name: "aggregation_policy", Type: AttrIdentifier, IdentifierType: "SchemaObjectIdentifier", CreateSql: shared, SetSql: shared},
			{Name: "entity_key", Type: AttrList, Elem: &Attribute{Type: AttrString}, CreateSql: ListParameter("Column")},
			{Name: "force", Type: AttrBool, SetSql: BoolKeyword(Sql("FORCE"))},
			{Name: "unset_aggregation_policy", SqlOnly: true, UnsetSql: StaticSQL("unsetAggregationPolicy", "UNSET AGGREGATION POLICY")},
		},
	}
	names := func(use ClauseUse) []string {
		return collections.Map(attr.NestedSqlFieldsFor(use), func(child Attribute) string { return child.Name })
	}
	require.Equal(t, []string{"aggregation_policy", "entity_key"}, names(ClauseCreate))
	require.Equal(t, []string{"set", "aggregation_policy", "force"}, names(ClauseSet))
	require.Equal(t, []string{"unset_aggregation_policy"}, names(ClauseUnset))
	require.Equal(t, "TableAggregationPolicy", attr.SqlStructNameFor("Table", ClauseCreate))
	require.Equal(t, "TableSetAggregationPolicy", attr.SqlStructNameFor("Table", ClauseSet))
	require.Equal(t, "TableUnsetAggregationPolicy", attr.SqlStructNameFor("Table", ClauseUnset))
}
