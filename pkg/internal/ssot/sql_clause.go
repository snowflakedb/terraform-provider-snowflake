package ssot

import (
	"fmt"
	"strings"
)

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

// ClauseUse selects CreateSql, SetSql, or UnsetSql.
type ClauseUse int

const (
	ClauseCreate ClauseUse = iota
	ClauseSet
	ClauseUnset
)

// SqlClause is one attribute's SQL in CREATE, SET, or UNSET.
// A zero SqlClause in that statement's list uses the type default.
// Create and set default from the attribute type. Unset defaults to a bool keyword.
// An explicit constructor replaces that default. Clause applies functional options on top of it.
// A zero SqlClause outside the list is absent. Leave an attribute out of Unset to omit unset.
// Set UnsetSql to a parameter clause to unset by setting a value.
type SqlClause struct {
	Present  bool
	fromType bool
	Kind     SqlKind
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

// Clause is the type default plus options.
// On unset the default is a bool keyword, not the create and set parameter.
func Clause(opts ...SqlClauseOption) SqlClause {
	c := SqlClause{Present: true, fromType: true}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// EffectiveSqlClause is the clause used for one statement list entry.
// An explicit constructor wins. Clause overlays its options on the type default.
// A zero clause uses that default. Call it for attributes that are in the list.
func (a *Attribute) EffectiveSqlClause(use ClauseUse) SqlClause {
	explicit := a.rawSqlClause(use)
	if explicit.Present && !explicit.fromType {
		return explicit
	}
	base := a.defaultSqlClause(use)
	if !explicit.fromType {
		return base
	}
	return overlaySqlClause(base, explicit)
}

func (a *Attribute) rawSqlClause(use ClauseUse) SqlClause {
	switch use {
	case ClauseCreate:
		return a.CreateSql
	case ClauseSet:
		return a.SetSql
	case ClauseUnset:
		return a.UnsetSql
	default:
		panic(fmt.Sprintf("attribute %q: unknown sql clause", a.Name))
	}
}

func (a *Attribute) defaultSqlClause(use ClauseUse) SqlClause {
	if use == ClauseUnset {
		return BoolKeyword()
	}
	switch a.Type {
	case AttrInt:
		return IntParameter()
	case AttrBool:
		return BoolParameter()
	case AttrString:
		return StringParameter()
	case AttrEnum:
		return EnumParameter()
	case AttrList:
		return ListParameter("")
	case AttrIdentifier:
		return IdentifierKeyword("")
	default:
		panic(fmt.Sprintf("attribute %q: no default SQL clause for type %q", a.Name, a.Type))
	}
}

func overlaySqlClause(base, over SqlClause) SqlClause {
	if over.SQL != "" {
		base.SQL = over.SQL
	}
	if over.Quotes != "" {
		base.Quotes = over.Quotes
	}
	if over.Parentheses != "" {
		base.Parentheses = over.Parentheses
	}
	if over.Equals != "" {
		base.Equals = over.Equals
	}
	if over.Required != nil {
		base.Required = over.Required
	}
	if over.ItemStruct != "" {
		base.ItemStruct = over.ItemStruct
	}
	if over.IdType != "" {
		base.IdType = over.IdType
	}
	return base
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
