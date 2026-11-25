package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// appendObjectList adds a list or set of objects. The element is a struct of the Elem Children, with an optional struct
// per variant of the element. See ssot.ListOfObjects and ssot.ObjectParameter.
func (v *QueryStruct) appendObjectList(attr ssot.Attribute, clause ssot.SqlClause, required bool) *QueryStruct {
	name := clause.Field
	prefix := clause.SQL
	if name == "" {
		name = sqlToFieldName(clause.SQLText(attr), true)
		prefix = clause.SQLText(attr)
	}
	var transformer FieldTransformer
	if prefix != "" {
		transformer = parameterOptionsFromClause(clause, required).SQL(prefix)
	} else {
		transformer = listOptionsFromClause(clause, required)
	}
	v = v.ListQueryStructField(name, queryStructFromObjectList(v.owner, attr), transformer)
	if required {
		v = v.WithValidation(ValidateValueSet, name)
	}
	return v
}

// appendWrapStringClause adds a helper struct: static SQL, the string, then static SQL.
func (v *QueryStruct) appendWrapStringClause(attr ssot.Attribute, clause ssot.SqlClause, required bool) *QueryStruct {
	nested := NewQueryStruct(attr.SqlStructName(v.owner))
	if clause.Before.Field != "" {
		nested = nested.SQLWithCustomFieldName(clause.Before.Field, clause.Before.SQL)
	}
	nested = nested.Text(clauseGoFieldName(attr, ssot.SqlClause{}), KeywordOptions().NoQuotes().Required())
	if clause.After.Field != "" {
		nested = nested.SQLWithCustomFieldName(clause.After.Field, clause.After.SQL)
	}
	opts := ListOptions().NoComma()
	if required {
		opts = opts.Required()
	}
	name := clauseGoFieldName(attr, clause)
	if required {
		return v.QueryStructField(name, nested, opts)
	}
	return v.OptionalQueryStructField(name, nested, opts)
}

// appendObjectClause adds a SqlOnly nested struct, such as ENCRYPTION = (...), whose fields are the CREATE Children of attr.
func (v *QueryStruct) appendObjectClause(attr ssot.Attribute, clause ssot.SqlClause, required bool) *QueryStruct {
	nested := queryStructFromObjectSupplied(v.owner, attr, ssot.ClauseCreate, v.suppliedStructs)
	name := clauseGoFieldName(attr, clause)
	if clause.AsParameter {
		opts := ParameterOptions().SQL(clause.SQL)
		if required {
			return v.QueryStructField(name, nested, opts.Required())
		}
		return v.OptionalQueryStructField(name, nested, opts)
	}
	opts := structFieldTransformer(clause, required)
	if required {
		return v.QueryStructField(name, nested, opts)
	}
	return v.OptionalQueryStructField(name, nested, opts)
}

// QueryStructFromAttribute builds the CREATE query struct of a SqlOnly object attribute. owner is the interface name, and the
// struct name is the StructName of its clause or that owner plus the attribute name.
// The SDK definition calls it for one ALTER field and keeps the rest of the operation handwritten. SET and UNSET projections
// use QueryStructFromAttributeUse.
func QueryStructFromAttribute(owner string, attr ssot.Attribute) *QueryStruct {
	return queryStructFromObjectUse(owner, attr, ssot.ClauseCreate)
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
	return queryStructFromObjectSupplied(owner, attr, use, nil)
}

// queryStructFromObjectSupplied is queryStructFromObjectUse with nested structs a child PreexistingStructOf clause can reuse.
func queryStructFromObjectSupplied(owner string, attr ssot.Attribute, use ssot.ClauseUse, supplied map[string]*QueryStruct) *QueryStruct {
	nested := NewQueryStruct(attr.SqlStructNameFor(owner, use))
	nested.owner = owner
	nested.suppliedStructs = supplied
	statement := clauseUseStatement(use)
	fields := attr.NestedSqlFieldsFor(use)
	for i := range fields {
		child := &fields[i]
		if child.IsTags() {
			switch use {
			case ssot.ClauseCreate:
				nested = nested.OptionalTags()
			case ssot.ClauseSet:
				nested = nested.SetTags()
			case ssot.ClauseUnset:
				nested = nested.UnsetTags()
			}
			continue
		}
		nested = nested.appendAttributeClause(*child, child.EffectiveSqlClause(use), statement)
	}
	nested = nested.withChildValidations(attr, use)
	if sharedToOptsClause(attr, use) {
		nested = nested.WithSharedToOpts()
	}
	return nested
}

// sharedToOptsClause reports whether the attribute's clause for this statement asks for a standalone toOpts. A missing
// clause is not shared, so the check reads the clause as written and does not resolve the type default.
func sharedToOptsClause(attr ssot.Attribute, use ssot.ClauseUse) bool {
	switch use {
	case ssot.ClauseCreate:
		return attr.CreateSql.SharedToOpts
	case ssot.ClauseSet:
		return attr.SetSql.SharedToOpts
	case ssot.ClauseUnset:
		return attr.UnsetSql.SharedToOpts
	default:
		panic("unknown sql clause")
	}
}

// withChildValidations adds the Validations of an object attribute over the Children in the struct of one statement, in
// Children order. A validation with fewer than two of its Children in the statement is left out.
func (v *QueryStruct) withChildValidations(attr ssot.Attribute, use ssot.ClauseUse) *QueryStruct {
	statement := clauseUseStatement(use)
	for _, validation := range attr.Validations {
		if validation.NotInSdk {
			continue
		}
		names := collections.Map(attr.ValidationFieldsFor(validation, use), func(child ssot.Attribute) string {
			return clauseFieldName(child, statement)
		})
		v = v.withValidationOf(validation.Kind, names)
	}
	return v
}

// withValidationOf adds the SDK validation of one ssot validation kind over the Go fields. Fewer than two fields add
// nothing. Two conflicting fields are ConflictingFields, and more are MoreThanOneValueSet.
func (v *QueryStruct) withValidationOf(kind ssot.ValidationKind, names []string) *QueryStruct {
	if len(names) < 2 {
		return v
	}
	switch kind {
	case ssot.ValidationAtLeastOne:
		return v.WithValidation(AtLeastOneValueSet, names...)
	case ssot.ValidationExactlyOne:
		return v.WithValidation(ExactlyOneValueSet, names...)
	case ssot.ValidationConflicting:
		if len(names) == 2 {
			return v.WithValidation(ConflictingFields, names...)
		}
		return v.WithValidation(MoreThanOneValueSet, names...)
	default:
		panic(fmt.Sprintf("unknown validation kind %q", kind))
	}
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
	if list.Elem != nil {
		element = element.withChildValidations(*list.Elem, ssot.ClauseCreate)
	}
	if clause.Parentheses != ssot.SqlParensYes {
		return element
	}
	return NewQueryStruct(structName+"Item").
		QueryStructField(structName, element, ListOptions().Parentheses().NoComma().Required())
}
