package ssot

import "strings"

// SqlKind is the SQL clause shape: parameter, keyword, or identifier.
type SqlKind string

const (
	SqlKindParameter  SqlKind = "parameter"
	SqlKindKeyword    SqlKind = "keyword"
	SqlKindIdentifier SqlKind = "identifier"
)

// SqlQuotes is the quote ddl tag. Zero means no quote tag.
type SqlQuotes string

const (
	SqlQuotesSingle SqlQuotes = "single_quotes"
	SqlQuotesNo     SqlQuotes = "no_quotes"
)

// SqlParens is the parentheses ddl tag. Zero means no parentheses tag.
type SqlParens string

const (
	SqlParensYes SqlParens = "parentheses"
	SqlParensNo  SqlParens = "no_parentheses"
)

// SqlEquals is the equals ddl tag. Zero means the equals sign stays.
type SqlEquals string

const (
	SqlEqualsNo SqlEquals = "no_equals"
)

// SqlClause is one attribute's SQL in CREATE, SET, or UNSET.
// A zero SqlClause is not present in that statement.
// TODO [next PRs]: Default the create and set clause from the attribute type.
// An enum uses EnumParameter, an int uses IntParameter, and the same for bool, string, list, and identifier.
// Functional options override that default.
type SqlClause struct {
	Present bool
	Kind    SqlKind
	// SQL is the clause text. Empty means the uppercase attribute Name.
	SQL         string
	Quotes      SqlQuotes
	Parentheses SqlParens
	Equals      SqlEquals
	// Required overrides optional-ness. Nil follows Attribute.Required on create and is optional on set.
	Required *bool
	// ItemStruct is the list-item query struct name. ListParameter sets it.
	ItemStruct string
	// IdType is the identifier type, without a pointer. IdentifierKeyword sets it.
	IdType string
}

// SqlClauseOption changes one default on a SqlClause.
// A zero option field keeps the constructor default.
type SqlClauseOption func(*SqlClause)

// SQLText is the clause text. Empty SQL is the uppercase attribute Name.
func (c SqlClause) SQLText(attr Attribute) string {
	if c.SQL != "" {
		return c.SQL
	}
	return strings.ToUpper(attr.Name)
}

// Ddl is the ddl tag body: kind, then quote, parentheses, and equals tags when set.
func (c SqlClause) Ddl() string {
	parts := []string{string(c.Kind)}
	if c.Quotes != "" {
		parts = append(parts, string(c.Quotes))
	}
	if c.Parentheses != "" {
		parts = append(parts, string(c.Parentheses))
	}
	if c.Equals != "" {
		parts = append(parts, string(c.Equals))
	}
	return strings.Join(parts, ",")
}

func IntParameter(opts ...SqlClauseOption) SqlClause {
	return newClause(SqlKindParameter, nil, opts)
}

func BoolParameter(opts ...SqlClauseOption) SqlClause {
	return newClause(SqlKindParameter, nil, opts)
}

func StringParameter(opts ...SqlClauseOption) SqlClause {
	return newClause(SqlKindParameter, func(c *SqlClause) {
		c.Quotes = SqlQuotesSingle
	}, opts)
}

func EnumParameter(opts ...SqlClauseOption) SqlClause {
	return newClause(SqlKindParameter, func(c *SqlClause) {
		c.Quotes = SqlQuotesNo
	}, opts)
}

// ListParameter is a parameter clause for a list.
// itemStruct is the item query struct, which stays outside this clause.
func ListParameter(itemStruct string, opts ...SqlClauseOption) SqlClause {
	return newClause(SqlKindParameter, func(c *SqlClause) {
		c.Parentheses = SqlParensYes
		c.ItemStruct = itemStruct
	}, opts)
}

// BoolKeyword is a bare keyword.
// TODO [next PRs]: Default unset to a bool keyword. Overrides leave the attribute out of unset, or unset it by setting a value.
func BoolKeyword(opts ...SqlClauseOption) SqlClause {
	return newClause(SqlKindKeyword, nil, opts)
}

// IdentifierKeyword is an identifier keyword.
// idType is the identifier type, without a pointer. Pass Sql to set the clause text.
func IdentifierKeyword(idType string, opts ...SqlClauseOption) SqlClause {
	return newClause(SqlKindIdentifier, func(c *SqlClause) {
		c.IdType = idType
	}, opts)
}

func Sql(text string) SqlClauseOption {
	return func(c *SqlClause) { c.SQL = text }
}

func NoQuotes() SqlClauseOption {
	return func(c *SqlClause) { c.Quotes = SqlQuotesNo }
}

func SingleQuotes() SqlClauseOption {
	return func(c *SqlClause) { c.Quotes = SqlQuotesSingle }
}

func Parentheses() SqlClauseOption {
	return func(c *SqlClause) { c.Parentheses = SqlParensYes }
}

func NoParentheses() SqlClauseOption {
	return func(c *SqlClause) { c.Parentheses = SqlParensNo }
}

func NoEquals() SqlClauseOption {
	return func(c *SqlClause) { c.Equals = SqlEqualsNo }
}

func Required() SqlClauseOption {
	return func(c *SqlClause) { c.Required = new(true) }
}

func Optional() SqlClauseOption {
	return func(c *SqlClause) { c.Required = new(false) }
}

func newClause(kind SqlKind, defaults func(*SqlClause), opts []SqlClauseOption) SqlClause {
	c := SqlClause{Present: true, Kind: kind}
	if defaults != nil {
		defaults(&c)
	}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}
