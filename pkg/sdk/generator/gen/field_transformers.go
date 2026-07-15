package gen

import "slices"

type FieldTransformer interface {
	Transform(f *Field) *Field
}

func StaticOptions() *StaticTransformer {
	return new(StaticTransformer)
}

type StaticTransformer struct {
	sqlPrefix string
}

func (v *StaticTransformer) SQL(sqlPrefix string) *StaticTransformer {
	v.sqlPrefix = sqlPrefix
	return v
}

func (v *StaticTransformer) Transform(f *Field) *Field {
	addTagIfMissing(f.Tags, "ddl", "static")
	addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	return f
}

type KeywordTransformer struct {
	required    bool
	sqlPrefix   string
	quotes      string
	parentheses string
	comma       string
}

func KeywordOptions() *KeywordTransformer {
	return new(KeywordTransformer)
}

func (v *KeywordTransformer) Required() *KeywordTransformer {
	v.required = true
	return v
}

func (v *KeywordTransformer) SQL(sqlPrefix string) *KeywordTransformer {
	v.sqlPrefix = sqlPrefix
	return v
}

func (v *KeywordTransformer) NoQuotes() *KeywordTransformer {
	v.quotes = "no_quotes"
	return v
}

func (v *KeywordTransformer) SingleQuotes() *KeywordTransformer {
	v.quotes = "single_quotes"
	return v
}

func (v *KeywordTransformer) DoubleQuotes() *KeywordTransformer {
	v.quotes = "double_quotes"
	return v
}

func (v *KeywordTransformer) DoubleDollarQuotes() *KeywordTransformer {
	v.quotes = "double_dollar_quotes"
	return v
}

func (v *KeywordTransformer) Parentheses() *KeywordTransformer {
	v.parentheses = "parentheses"
	return v
}

func (v *KeywordTransformer) MustParentheses() *KeywordTransformer {
	v.parentheses = "must_parentheses"
	return v
}

func (v *KeywordTransformer) NoComma() *KeywordTransformer {
	v.comma = "no_comma"
	return v
}

func (v *KeywordTransformer) Transform(f *Field) *Field {
	addTagIfMissing(f.Tags, "ddl", "keyword")
	if v.required {
		f.Required = true
	}
	addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	addTagIfMissing(f.Tags, "ddl", v.quotes)
	addTagIfMissing(f.Tags, "ddl", v.parentheses)
	addTagIfMissing(f.Tags, "ddl", v.comma)
	return f
}

type ParameterTransformer struct {
	required    bool
	sqlPrefix   string
	quotes      string
	parentheses string
	equals      string
	reverse     string
}

func ParameterOptions() *ParameterTransformer {
	return new(ParameterTransformer)
}

func (v *ParameterTransformer) Required() *ParameterTransformer {
	v.required = true
	return v
}

func (v *ParameterTransformer) SQL(sqlPrefix string) *ParameterTransformer {
	v.sqlPrefix = sqlPrefix
	return v
}

func (v *ParameterTransformer) NoQuotes() *ParameterTransformer {
	v.quotes = "no_quotes"
	return v
}

func (v *ListTransformer) MustParentheses() *ListTransformer {
	v.parentheses = "must_parentheses"
	return v
}

func (v *ParameterTransformer) NoEquals() *ParameterTransformer {
	v.equals = "no_equals"
	return v
}

func (v *ParameterTransformer) ArrowEquals() *ParameterTransformer {
	v.equals = "arrow_equals"
	return v
}

func (v *ParameterTransformer) SingleQuotes() *ParameterTransformer {
	v.quotes = "single_quotes"
	return v
}

func (v *ParameterTransformer) DoubleQuotes() *ParameterTransformer {
	v.quotes = "double_quotes"
	return v
}

// Warning: Double dollar quotes tag assumes that passed input is sanitated.
// Double dollar quotes cannot be escaped in Snowflake.
// This is why the caller is resposible to allow only sanitized input for fields with this tag.
func (v *ParameterTransformer) DoubleDollarQuotes() *ParameterTransformer {
	v.quotes = "double_dollar_quotes"
	return v
}

func (v *ParameterTransformer) NoParentheses() *ParameterTransformer {
	v.parentheses = "no_parentheses"
	return v
}

func (v *ParameterTransformer) Parentheses() *ParameterTransformer {
	v.parentheses = "parentheses"
	return v
}

func (v *ParameterTransformer) MustParentheses() *ParameterTransformer {
	v.parentheses = "must_parentheses"
	return v
}

func (v *ParameterTransformer) Reverse() *ParameterTransformer {
	v.reverse = "reverse"
	return v
}

func (v *ParameterTransformer) Transform(f *Field) *Field {
	addTagIfMissing(f.Tags, "ddl", "parameter")
	if v.required {
		f.Required = true
	}
	addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	addTagIfMissing(f.Tags, "ddl", v.quotes)
	addTagIfMissing(f.Tags, "ddl", v.parentheses)
	addTagIfMissing(f.Tags, "ddl", v.equals)
	addTagIfMissing(f.Tags, "ddl", v.reverse)
	return f
}

type ListTransformer struct {
	required    bool
	sqlPrefix   string
	parentheses string
	equals      string
	comma       string
}

func ListOptions() *ListTransformer {
	return new(ListTransformer)
}

func (v *ListTransformer) Required() *ListTransformer {
	v.required = true
	return v
}

func (v *ListTransformer) Parentheses() *ListTransformer {
	v.parentheses = "parentheses"
	return v
}

func (v *ListTransformer) NoParentheses() *ListTransformer {
	v.parentheses = "no_parentheses"
	return v
}

func (v *ListTransformer) NoEquals() *ListTransformer {
	v.equals = "no_equals"
	return v
}

func (v *ListTransformer) NoComma() *ListTransformer {
	v.comma = "no_comma"
	return v
}

func (v *ListTransformer) SQL(sqlPrefix string) *ListTransformer {
	v.sqlPrefix = sqlPrefix
	return v
}

func (v *ListTransformer) Transform(f *Field) *Field {
	addTagIfMissing(f.Tags, "ddl", "list")
	if v.required {
		f.Required = true
	}
	addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	addTagIfMissing(f.Tags, "ddl", v.parentheses)
	addTagIfMissing(f.Tags, "ddl", v.equals)
	addTagIfMissing(f.Tags, "ddl", v.comma)
	return f
}

type IdentifierTransformer struct {
	required        bool
	sqlPrefix       string
	quotes          string
	equals          string
	systemReference string
}

func IdentifierOptions() *IdentifierTransformer {
	return new(IdentifierTransformer)
}

func (v *IdentifierTransformer) SQL(sqlPrefix string) *IdentifierTransformer {
	v.sqlPrefix = sqlPrefix
	return v
}

func (v *IdentifierTransformer) SingleQuotes() *IdentifierTransformer {
	v.quotes = "single_quotes"
	return v
}

func (v *IdentifierTransformer) DoubleQuotes() *IdentifierTransformer {
	v.quotes = "double_quotes"
	return v
}

func (v *IdentifierTransformer) DoubleDollarQuotes() *IdentifierTransformer {
	v.quotes = "double_dollar_quotes"
	return v
}

func (v *IdentifierTransformer) Required() *IdentifierTransformer {
	v.required = true
	return v
}

func (v *IdentifierTransformer) NoEquals() *IdentifierTransformer {
	v.equals = "no_equals"
	return v
}

func (v *IdentifierTransformer) Equals() *IdentifierTransformer {
	v.equals = "equals"
	return v
}

func (v *IdentifierTransformer) SystemReference(objectType string) *IdentifierTransformer {
	v.systemReference = objectType
	return v
}

func (v *IdentifierTransformer) Transform(f *Field) *Field {
	addTagIfMissing(f.Tags, "ddl", "identifier")
	if v.required {
		f.Required = true
	}
	if v.systemReference != "" {
		addTagIfMissing(f.Tags, "ddl", "system_reference")
		addTagIfMissing(f.Tags, "sql", v.systemReference)
	} else {
		addTagIfMissing(f.Tags, "sql", v.sqlPrefix)
	}
	addTagIfMissing(f.Tags, "ddl", v.quotes)
	addTagIfMissing(f.Tags, "ddl", v.equals)
	return f
}

// InlineTransformer produces ddl:"-" on a struct field, causing the SQL builder
// to render the sub-struct's fields inline without any wrapper keyword.
// This is equivalent to the ddl:"-" tag used directly on Options struct fields.
// Use when a sub-struct provides its own SQL keywords internally (e.g. AccountDrop)
type InlineTransformer struct{}

func InlineOptions() *InlineTransformer {
	return new(InlineTransformer)
}

func (v *InlineTransformer) Transform(f *Field) *Field {
	f.Tags["ddl"] = []string{"-"}
	return f
}

func addTagIfMissing(m map[string][]string, key string, value string) {
	if len(value) > 0 {
		if val, ok := m[key]; ok {
			if !slices.Contains(val, value) {
				m[key] = append(val, value)
			}
		} else {
			m[key] = []string{value}
		}
	}
}
