package ssot

import (
	"testing"

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

func Test_SqlClause_optionOverridesDefault(t *testing.T) {
	clause := StringParameter(NoQuotes())

	require.Equal(t, "parameter,no_quotes", clause.Ddl())
	require.Equal(t, SqlQuotesNo, clause.Quotes)
}
