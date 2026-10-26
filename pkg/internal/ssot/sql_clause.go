package ssot

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// SqlKind is the SQL clause shape: parameter, keyword, identifier, static, or object.
type SqlKind string

const (
	SqlKindParameter  SqlKind = "parameter"
	SqlKindKeyword    SqlKind = "keyword"
	SqlKindIdentifier SqlKind = "identifier"
	SqlKindStatic     SqlKind = "static"
	// SqlKindObject is a nested struct: the fields are the Children of the attribute, or the Children of the Elem object
	// of a list. See ObjectParameter.
	SqlKindObject SqlKind = "object"
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

// SqlEquals is the equals ddl tag. Zero leaves the identifier default, which has no equals tag.
type SqlEquals string

const (
	SqlEqualsYes SqlEquals = "equals"
	SqlEqualsNo  SqlEquals = "no_equals"
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
	// Field is the Go field name. A static clause needs an unexported one (StaticSQL sets it). Any other clause uses it
	// instead of the name derived from the SQL text or the attribute name. FieldName sets it.
	Field string
	// Enum makes the field an enum parameter of this enum, whatever the attribute type is. EnumOf sets it.
	Enum *Enum
	// Variants are the provider-like variants of an object: each is a struct of the attributes it refers to, and the
	// object has one optional field per variant, exactly one of which is set. WithVariants sets it.
	Variants []SqlVariant
}

// SqlVariant is one variant of an object, such as a storage provider. Fields are references to attributes, in SQL order.
// An attribute is a field with the CREATE clause it already has.
type SqlVariant struct {
	Key    string
	Fields []*Attribute
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

// StaticSQL is a private static CREATE field.
// field is the unexported Go name. sql is the ddl text, such as WAREHOUSE_TYPE = 'ADAPTIVE'.
// The attribute is SqlOnly so schema, resource, model, and assertion generation skip it.
func StaticSQL(field, sql string) SqlClause {
	return newClause(SqlKindStatic, func(c *SqlClause) {
		c.Field = field
		c.SQL = sql
	}, nil)
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
		if item := a.enumListItemName(); item != "" {
			return ListParameter(item)
		}
		return ListParameter("")
	case AttrIdentifier:
		return IdentifierKeyword("")
	default:
		panic(fmt.Sprintf("attribute %q: no default SQL clause for type %q", a.Name, a.Type))
	}
}

// EnumListItemName is the query struct for a list whose element is an enum.
// The struct has one field, Value, a required single-quoted keyword of that enum.
func EnumListItemName(enum *Enum) string {
	if enum == nil || enum.Name == "" {
		panic("enum list item needs an enum name")
	}
	return enum.Name + "ListItem"
}

func (a *Attribute) enumListItemName() string {
	if a.Elem == nil || a.Elem.Type != AttrEnum || a.Elem.Enum == nil || a.Elem.Enum.Name == "" {
		return ""
	}
	return EnumListItemName(a.Elem.Enum)
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
	if over.Field != "" {
		base.Field = over.Field
	}
	if over.Enum != nil {
		base.Enum = over.Enum
	}
	if len(over.Variants) > 0 {
		base.Variants = over.Variants
	}
	return base
}

// ObjectParameter is a nested struct clause. On the Elem object of a list it is the shape of each element: Parentheses()
// puts each element in parentheses. On a SqlOnly object attribute it is a nested struct inside a variant: Sql is the text
// before it and Parentheses() wraps its fields. The fields are the Children.
func ObjectParameter(opts ...SqlClauseOption) SqlClause {
	return newClause(SqlKindObject, nil, opts)
}

// ListOfObjects is a parameter clause for a list or set whose Elem is an object. Sql is the text before the list and
// Parentheses() wraps the whole list.
func ListOfObjects(opts ...SqlClauseOption) SqlClause {
	return newClause(SqlKindParameter, nil, opts)
}

// WithVariants declares the variants of an object. The object holds the Children that no variant refers to, and one
// optional struct per variant. Exactly one variant is set.
func WithVariants(variants ...SqlVariant) SqlClauseOption {
	return func(c *SqlClause) { c.Variants = variants }
}

// Variant is one variant of an object: the attributes it consists of, in SQL order. The struct is named after the key, the
// list attribute, and Params (GCS and storage_location give GCSStorageLocationParams).
func Variant(key string, fields ...*Attribute) SqlVariant {
	return SqlVariant{Key: key, Fields: fields}
}

// EnumOf makes the field an enum parameter of enum, whatever the attribute type is.
func EnumOf(enum *Enum) SqlClauseOption {
	return func(c *SqlClause) { c.Enum = enum }
}

// FieldName sets the Go field name when it should not be derived from the SQL text or the attribute name.
func FieldName(name string) SqlClauseOption {
	return func(c *SqlClause) { c.Field = name }
}

// WithCreateSql is a copy of the attribute with another CREATE clause, for the place that needs the same attribute with a
// different shape, such as an enum of another provider. Name, type, and description are not restated.
func WithCreateSql(attr *Attribute, clause SqlClause) Attribute {
	view := *attr
	view.CreateSql = clause
	return view
}

// SqlStructName is the Go struct name of a list of objects or of a nested struct: the object and the attribute,
// ExternalVolume and storage_location give ExternalVolumeStorageLocation.
func (a Attribute) SqlStructName(owner string) string {
	return owner + genhelpers.SnakeCaseToCamel(a.Name)
}

// SqlVariantStructName is the Go struct name of a variant of a list of objects: GCS and storage_location give
// GCSStorageLocationParams.
func (a Attribute) SqlVariantStructName(key string) string {
	return key + genhelpers.SnakeCaseToCamel(a.Name) + "Params"
}

// NestedFields are the fields of a block: the Children of the Elem object of a list or set, or the Children of an object.
func (a Attribute) NestedFields() []Attribute {
	switch {
	case a.IsObjectList():
		return a.Elem.Children
	case a.Type == AttrObject:
		return a.Children
	default:
		return nil
	}
}

// IsObjectList reports whether the attribute is a list or set whose elements are objects.
func (a Attribute) IsObjectList() bool {
	return (a.Type == AttrList || a.Type == AttrSet) && a.Elem != nil && a.Elem.Type == AttrObject
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

func Equals() SqlClauseOption {
	return func(c *SqlClause) { c.Equals = SqlEqualsYes }
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
