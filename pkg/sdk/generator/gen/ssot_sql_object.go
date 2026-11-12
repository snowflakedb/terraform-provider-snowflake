package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// appendObjectList adds a list or set of objects. The element is a struct of the Elem Children, with an optional struct
// per variant of the element. See ssot.ListOfObjects and ssot.ObjectParameter.
func (v *QueryStruct) appendObjectList(attr ssot.Attribute, clause ssot.SqlClause, required bool) *QueryStruct {
	name := clause.Field
	if name == "" {
		name = sqlToFieldName(clause.SQLText(attr), true)
	}
	transformer := parameterOptionsFromClause(clause, required).SQL(clause.SQLText(attr))
	v = v.ListQueryStructField(name, queryStructFromObjectList(v.owner, attr), transformer)
	if required {
		v = v.WithValidation(ValidateValueSet, name)
	}
	return v
}

// appendObjectClause adds a SqlOnly nested struct, such as ENCRYPTION = (...), whose fields are the CREATE Children of attr.
func (v *QueryStruct) appendObjectClause(attr ssot.Attribute, clause ssot.SqlClause, required bool) *QueryStruct {
	nested := queryStructFromObjectUse(v.owner, attr, ssot.ClauseCreate)
	name := clauseGoFieldName(attr, clause)
	if clause.AsParameter {
		opts := ParameterOptions().SQL(clause.SQL)
		if required {
			return v.QueryStructField(name, nested, opts.Required())
		}
		return v.OptionalQueryStructField(name, nested, opts)
	}
	opts := ListOptions()
	if clause.NoComma {
		opts = opts.NoComma()
	}
	if clause.Parentheses == ssot.SqlParensYes {
		opts = opts.Parentheses()
	}
	if clause.SQL != "" {
		opts = opts.SQL(clause.SQL)
	}
	if required {
		return v.QueryStructField(name, nested, opts.Required())
	}
	return v.OptionalQueryStructField(name, nested, opts)
}

// QueryStructFromAttribute builds the CREATE query struct of a SqlOnly object attribute. owner is the interface name, and the
// struct name is that owner plus the attribute name. The SDK definition calls it for one ALTER field and keeps the rest
// of the operation handwritten. SET and UNSET projections use QueryStructFromAttributeUse.
func QueryStructFromAttribute(owner string, attr ssot.Attribute) *QueryStruct {
	return QueryStructFromAttributeUse(owner, attr, ssot.ClauseCreate)
}

// QueryStructFromAttributeUse builds the query struct of a SqlOnly object attribute for one statement projection.
// CREATE is QueryStructFromAttribute. SET and UNSET keep FieldName on that clause for the struct name, and only the
// children that have that clause. Handwritten ALTER calls this so tables and iceberg tables share SET/UNSET structs.
func QueryStructFromAttributeUse(owner string, attr ssot.Attribute, use ssot.ClauseUse) *QueryStruct {
	return queryStructFromObjectUse(owner, attr, use)
}

func clauseUseStatement(use ssot.ClauseUse) attributeSqlStatement {
	switch use {
	case ssot.ClauseCreate:
		return attributeSqlCreate
	case ssot.ClauseSet:
		return attributeSqlSet
	case ssot.ClauseUnset:
		return attributeSqlUnset
	default:
		panic("unknown sql clause")
	}
}

// queryStructFromObjectUse builds a fresh query struct for a SqlOnly object attribute: the projection's Children are its fields.
func queryStructFromObjectUse(owner string, attr ssot.Attribute, use ssot.ClauseUse) *QueryStruct {
	nested := NewQueryStruct(attr.SqlStructNameFor(owner, use))
	nested.owner = owner
	statement := clauseUseStatement(use)
	fields := attr.NestedSqlFieldsFor(use)
	for i := range fields {
		child := &fields[i]
		nested = nested.appendAttributeClause(*child, child.EffectiveSqlClause(use), statement)
	}
	if attr.ChildrenRule != "" {
		names := collections.Map(attr.RuleFieldsFor(use), func(child ssot.Attribute) string {
			return clauseFieldName(child, statement)
		})
		switch attr.ChildrenRule {
		case ssot.ChildrenAtLeastOne:
			nested = nested.WithValidation(AtLeastOneValueSet, names...)
		case ssot.ChildrenExactlyOne:
			nested = nested.WithValidation(ExactlyOneValueSet, names...)
		}
	}
	return nested
}

// queryStructFromBlockAttribute builds the struct of an ALTER block: one element of a list of objects, or a SqlOnly object.
// ALTER blocks keep CREATE children, the same as QueryStructFromAttribute.
func queryStructFromBlockAttribute(owner string, attr ssot.Attribute) *QueryStruct {
	if attr.IsObjectList() {
		return queryStructFromObjectList(owner, attr)
	}
	return queryStructFromObjectUse(owner, attr, ssot.ClauseCreate)
}

// queryStructFromObjectList builds a fresh query struct for the element of the list, because a query struct has one parent.
// The element holds the Children that have a CREATE clause and that no variant refers to, in Children order, then one
// optional struct per variant, of which exactly one is set. With Parentheses on the element clause, each element is held by a
// <struct>Item struct that puts it in parentheses.
func queryStructFromObjectList(owner string, list ssot.Attribute) *QueryStruct {
	clause := list.Elem.CreateSql
	inVariant := make(map[string]struct{})
	for _, variant := range clause.Variants {
		for _, field := range variant.Fields {
			inVariant[field.Name] = struct{}{}
			// The fields of a nested struct belong to the variant too.
			for _, nested := range field.Children {
				inVariant[nested.Name] = struct{}{}
			}
		}
	}

	// A child that a sibling SqlOnly nested struct holds is written only inside that struct.
	for _, sibling := range list.NestedSqlFields() {
		if sibling.SqlOnly && sibling.Type == ssot.AttrObject {
			for _, held := range sibling.Children {
				inVariant[held.Name] = struct{}{}
			}
		}
	}

	structName := list.SqlStructName(owner)
	element := NewQueryStruct(structName)
	element.owner = owner
	fields := list.NestedSqlFields()
	for i := range fields {
		child := &fields[i]
		if _, used := inVariant[child.Name]; used || !child.CreateSql.Present {
			continue
		}
		element = element.appendAttributeClause(*child, child.EffectiveSqlClause(ssot.ClauseCreate), attributeSqlCreate)
	}
	variantFields := make([]string, 0, len(clause.Variants))
	for _, variant := range clause.Variants {
		variantName := list.SqlVariantStructName(variant.Key)
		variantStruct := NewQueryStruct(variantName)
		variantStruct.owner = owner
		for _, field := range variant.Fields {
			variantStruct = variantStruct.appendAttributeClause(*field, field.EffectiveSqlClause(ssot.ClauseCreate), attributeSqlCreate)
		}
		element = element.OptionalQueryStructField(variantName, variantStruct, ListOptions().NoComma())
		variantFields = append(variantFields, variantName)
	}
	if len(variantFields) > 1 {
		element = element.WithValidation(ExactlyOneValueSet, variantFields...)
	}
	if list.Elem != nil && list.Elem.ChildrenRule != "" {
		names := collections.Map(list.Elem.RuleFields(), func(child ssot.Attribute) string {
			return clauseFieldName(child, attributeSqlCreate)
		})
		switch list.Elem.ChildrenRule {
		case ssot.ChildrenAtLeastOne:
			element = element.WithValidation(AtLeastOneValueSet, names...)
		case ssot.ChildrenExactlyOne:
			element = element.WithValidation(ExactlyOneValueSet, names...)
		}
	}
	if clause.Parentheses != ssot.SqlParensYes {
		return element
	}
	return NewQueryStruct(structName+"Item").
		QueryStructField(structName, element, ListOptions().Parentheses().NoComma().Required())
}
