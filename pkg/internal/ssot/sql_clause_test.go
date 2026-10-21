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

func Test_SqlClause_optionOverridesDefault(t *testing.T) {
	clause := StringParameter(NoQuotes())

	require.Equal(t, "parameter,no_quotes", clause.Ddl())
	require.Equal(t, SqlQuotesNo, clause.Quotes)
}
