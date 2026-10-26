package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// attributeSqlStatement selects CreateSql, SetSql, or UnsetSql.
type attributeSqlStatement int

const (
	attributeSqlCreate attributeSqlStatement = iota
	attributeSqlSet
	attributeSqlUnset
)

// appendSqlSequence adds list in slice order.
// Tags emits optional TAG on CREATE. SET and UNSET reject tags.
// Parameters spreads params: WithParameters on CREATE and SET, WithParametersUnset on UNSET.
func (v *QueryStruct) appendSqlSequence(list []*ssot.Attribute, statement attributeSqlStatement) *QueryStruct {
	return v.appendSqlSequenceParams(list, statement, nil)
}

func (v *QueryStruct) appendSqlSequenceParams(list []*ssot.Attribute, statement attributeSqlStatement, params []ssot.ParameterDef) *QueryStruct {
	sawParameters := false
	for _, field := range list {
		if field == nil {
			panic(fmt.Sprintf("%s field is empty", statementName(statement)))
		}
		if field.IsTags() {
			if statement != attributeSqlCreate {
				panic(fmt.Sprintf("%s list contains tags", statementName(statement)))
			}
			v = v.OptionalTags()
			continue
		}
		if field.IsParameters() {
			if sawParameters {
				panic(fmt.Sprintf("%s list contains parameters twice", statementName(statement)))
			}
			sawParameters = true
			switch statement {
			case attributeSqlCreate, attributeSqlSet:
				v = v.WithParameters(params...)
			case attributeSqlUnset:
				v = v.WithParametersUnset(params...)
			default:
				panic(fmt.Sprintf("%s list contains parameters", statementName(statement)))
			}
			continue
		}
		clause := attributeClause(*field, statement)
		if !clause.Present {
			panic(fmt.Sprintf("attribute %q has no %s clause", field.Name, statementName(statement)))
		}
		v = v.appendAttributeClause(*field, clause, statement)
	}
	return v
}

func statementName(statement attributeSqlStatement) string {
	switch statement {
	case attributeSqlCreate:
		return "CREATE"
	case attributeSqlSet:
		return "SET"
	case attributeSqlUnset:
		return "UNSET"
	default:
		return "sql"
	}
}

func attributeClause(attr ssot.Attribute, statement attributeSqlStatement) ssot.SqlClause {
	switch statement {
	case attributeSqlCreate:
		return attr.EffectiveSqlClause(ssot.ClauseCreate)
	case attributeSqlSet:
		return attr.EffectiveSqlClause(ssot.ClauseSet)
	case attributeSqlUnset:
		return attr.EffectiveSqlClause(ssot.ClauseUnset)
	default:
		panic(fmt.Sprintf("attribute %q: unknown sql statement", attr.Name))
	}
}

func (v *QueryStruct) appendAttributeClause(attr ssot.Attribute, clause ssot.SqlClause, statement attributeSqlStatement) *QueryStruct {
	switch clause.Kind {
	case ssot.SqlKindParameter:
		return v.appendParameterClause(attr, clause, clauseRequired(attr, clause, statement))
	case ssot.SqlKindKeyword:
		return v.appendKeywordClause(attr, clause, clauseRequired(attr, clause, statement))
	case ssot.SqlKindIdentifier:
		return v.appendIdentifierClause(attr, clause, clauseRequired(attr, clause, statement))
	case ssot.SqlKindStatic:
		return v.appendStaticClause(attr, clause)
	default:
		panic(fmt.Sprintf("attribute %q: sql kind %q", attr.Name, clause.Kind))
	}
}

func (v *QueryStruct) appendStaticClause(attr ssot.Attribute, clause ssot.SqlClause) *QueryStruct {
	if clause.Field == "" || clause.SQL == "" {
		panic(fmt.Sprintf("attribute %q: static clause needs a field name and SQL", attr.Name))
	}
	return v.SQLWithCustomFieldName(clause.Field, clause.SQL)
}

// clauseRequired is Attribute.Required on create when the clause leaves Required nil.
// Set and unset treat that nil as optional. An explicit Required or Optional wins.
func clauseRequired(attr ssot.Attribute, clause ssot.SqlClause, statement attributeSqlStatement) bool {
	if clause.Required != nil {
		return *clause.Required
	}
	if statement == attributeSqlCreate {
		return attr.Required
	}
	return false
}

func (v *QueryStruct) appendParameterClause(attr ssot.Attribute, clause ssot.SqlClause, required bool) *QueryStruct {
	sql := clause.SQLText(attr)
	opts := parameterOptionsFromClause(clause, required)
	switch attr.Type {
	case ssot.AttrInt:
		if required {
			return v.NumberAssignment(sql, opts)
		}
		return v.OptionalNumberAssignment(sql, opts)
	case ssot.AttrBool:
		if required {
			return v.BooleanAssignment(sql, opts)
		}
		return v.OptionalBooleanAssignment(sql, opts)
	case ssot.AttrString:
		if required {
			return v.TextAssignment(sql, opts)
		}
		return v.OptionalTextAssignment(sql, opts)
	case ssot.AttrEnum:
		enum := attr.ColumnEnum()
		if enum == nil || enum.Name == "" {
			panic(fmt.Sprintf("attribute %q: enum parameter has no Enum", attr.Name))
		}
		if required {
			return v.EnumAssignment(sql, enum, opts)
		}
		return v.OptionalEnumAssignment(sql, enum, opts)
	case ssot.AttrList, ssot.AttrSet:
		item := clause.ItemStruct
		if item == "" && attr.Elem != nil && attr.Elem.Type == ssot.AttrIdentifier {
			item = attr.Elem.IdentifierType
		}
		if item == "" {
			panic(fmt.Sprintf("attribute %q: list parameter has no item struct", attr.Name))
		}
		if clause.ItemNested {
			// The item is a nested query struct from the SDK def, so the field keeps its Request type, mapping, and validations.
			itemStruct, ok := v.findItemStruct(item)
			if !ok {
				panic(fmt.Sprintf("attribute %q: nested item struct %q was not passed to the generator", attr.Name, item))
			}
			fieldName := sqlToFieldName(sql, true)
			v = v.ListQueryStructField(fieldName, itemStruct, opts.SQL(sql))
			if required {
				v = v.WithValidation(ValidateValueSet, fieldName)
			}
			return v
		}
		v = v.ListAssignment(sql, item, opts)
		if attr.Elem != nil && attr.Elem.Type == ssot.AttrIdentifier {
			v = v.WithValidation(ValidIdentifier, sqlToFieldName(sql, true))
		}
		return v
	default:
		panic(fmt.Sprintf("attribute %q: parameter clause has type %q", attr.Name, attr.Type))
	}
}

func (v *QueryStruct) appendKeywordClause(attr ssot.Attribute, clause ssot.SqlClause, required bool) *QueryStruct {
	if required {
		panic(fmt.Sprintf("attribute %q: keyword clause is required", attr.Name))
	}
	if clause.Quotes != "" || clause.Parentheses != "" || clause.Equals != "" {
		panic(fmt.Sprintf("attribute %q: keyword clause has parameter tags", attr.Name))
	}
	return v.OptionalSQL(clause.SQLText(attr))
}

func (v *QueryStruct) appendIdentifierClause(attr ssot.Attribute, clause ssot.SqlClause, required bool) *QueryStruct {
	if clause.IdType == "" {
		panic(fmt.Sprintf("attribute %q: identifier clause has no IdType", attr.Name))
	}
	if clause.Parentheses != "" {
		panic(fmt.Sprintf("attribute %q: identifier clause has parentheses", attr.Name))
	}
	sql := clause.SQLText(attr)
	opts := identifierOptionsFromClause(clause, required).SQL(sql)
	name := sqlToFieldName(sql, true)
	if required {
		return v.Identifier(name, clause.IdType, opts)
	}
	return v.OptionalIdentifier(name, clause.IdType, opts)
}

func parameterOptionsFromClause(clause ssot.SqlClause, required bool) *ParameterTransformer {
	opts := ParameterOptions()
	switch clause.Quotes {
	case "":
	case ssot.SqlQuotesSingle:
		opts = opts.SingleQuotes()
	case ssot.SqlQuotesNo:
		opts = opts.NoQuotes()
	default:
		panic(fmt.Sprintf("sql clause quotes %q", clause.Quotes))
	}
	switch clause.Parentheses {
	case "":
	case ssot.SqlParensYes:
		opts = opts.Parentheses()
	case ssot.SqlParensNo:
		opts = opts.NoParentheses()
	default:
		panic(fmt.Sprintf("sql clause parentheses %q", clause.Parentheses))
	}
	switch clause.Equals {
	case "":
	case ssot.SqlEqualsNo:
		opts = opts.NoEquals()
	default:
		panic(fmt.Sprintf("sql clause equals %q", clause.Equals))
	}
	if required {
		opts = opts.Required()
	}
	return opts
}

func identifierOptionsFromClause(clause ssot.SqlClause, required bool) *IdentifierTransformer {
	opts := IdentifierOptions()
	switch clause.Quotes {
	case "":
	case ssot.SqlQuotesSingle:
		opts = opts.SingleQuotes()
	default:
		panic(fmt.Sprintf("identifier clause quotes %q", clause.Quotes))
	}
	switch clause.Equals {
	case "":
	case ssot.SqlEqualsYes:
		opts = opts.Equals()
	case ssot.SqlEqualsNo:
		opts = opts.NoEquals()
	default:
		panic(fmt.Sprintf("identifier clause equals %q", clause.Equals))
	}
	if required {
		opts = opts.Required()
	}
	return opts
}

func (v *QueryStruct) findItemStruct(name string) (*QueryStruct, bool) {
	found, err := collections.FindFirst(v.itemStructs, func(s *QueryStruct) bool { return s.name == name })
	if err != nil {
		return nil, false
	}
	return *found, true
}
