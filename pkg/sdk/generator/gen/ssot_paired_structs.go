package gen

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// ShowPairedStructs builds a SHOW field chain from Columns by calling PairedStructs methods.
// The db row name is the unexported plural plus "Row". The plain name is the object name.
// This function builds the field chain only.
func ShowPairedStructs(obj ssot.Object) *PairedStructs {
	if obj.Name == "" || obj.NamePlural == "" {
		panic(fmt.Sprintf("object %q needs Name and NamePlural to build the SHOW paired struct", obj.Name))
	}
	return pairedStructsFromColumns(obj, false, unexported(obj.NamePlural)+"Row", obj.Name)
}

// DescribePairedStructs builds a DESCRIBE field chain from Columns by calling PairedStructs methods.
// The db row name is the unexported object name plus "DescRow". The plain name is the object name plus "Details".
// This function builds the field chain only.
func DescribePairedStructs(obj ssot.Object) *PairedStructs {
	if obj.Name == "" {
		panic(fmt.Sprintf("object %q needs Name to build the DESCRIBE paired struct", obj.Name))
	}
	pairs := pairedStructsFromColumns(obj, true, unexported(obj.Name)+"DescRow", obj.DescribePlainName())
	for _, field := range ssot.ObjectGenerators(obj.Generators...).Sdk.DescribePlainOnlyFields {
		pairs = pairs.PlainOnlyField(field.Name, field.Kind)
	}
	return pairs
}

// DescribeDetailsPlain is the plain details struct built from Describe.Columns. Id comes first unless a column is named id.
// A nested column is a slice of the hand-written element struct named by ssot.NestedDetailsElem.
func DescribeDetailsPlain(obj ssot.Object) *Plain {
	cfg := obj.SdkConfig()
	if !cfg.DescribeDetails {
		panic(fmt.Sprintf("object %q is not enabled for SDK DESCRIBE details generation", obj.Name))
	}
	plain := PlainStruct(obj.Name + "Details")
	hasId := slices.ContainsFunc(obj.Describe.Columns, func(c *ssot.Attribute) bool { return c != nil && c.OutputKey() == "id" })
	if !hasId {
		plain = plain.Field("Id", obj.IdType)
	}
	for _, column := range obj.Describe.Columns {
		if column == nil {
			panic(fmt.Sprintf("object %q: nil DESCRIBE column", obj.Name))
		}
		kind := detailsFieldKind(obj, *column, cfg.DescribeDetailsNestedFields)
		plain = plain.Field(sqlToFieldName(column.OutputKey(), true), kind)
	}
	return plain
}

func detailsFieldKind(obj ssot.Object, attr ssot.Attribute, nested []*ssot.Attribute) string {
	var kind string
	if slices.ContainsFunc(nested, func(n *ssot.Attribute) bool { return n.Name == attr.Name }) {
		kind = "[]" + ssot.NestedDetailsElem(obj, &attr)
	} else {
		concrete, _ := ssot.AttributeGoTypes(obj, attr)
		kind = strings.ReplaceAll(concrete, "sdk.", "")
	}
	if attr.IsOptionalColumn() {
		kind = "*" + kind
	}
	return kind
}

// DescribeNestedPlains are the plain structs of the columns of obj that are objects. The fields are the Children of the
// column, as plain values.
func DescribeNestedPlains(obj ssot.Object) []*Plain {
	var plains []*Plain
	for _, column := range obj.Describe.Columns {
		if column == nil || column.Type != ssot.AttrObject {
			continue
		}
		plain := PlainStruct(ssot.NestedObjectDetails(obj, *column))
		for _, child := range column.NestedFields() {
			if child.Type == ssot.AttrObject {
				panic(fmt.Sprintf("object %q: DESCRIBE column %q has an object child %q", obj.Name, column.Name, child.Name))
			}
			concrete, _ := ssot.AttributeGoTypes(obj, child)
			plain = plain.Field(ssot.ColumnFieldName(child), strings.ReplaceAll(concrete, "sdk.", ""))
		}
		plains = append(plains, plain)
	}
	return plains
}

// plainFromColumns builds the plain side of a SHOW column list. The db row is discarded.
func plainFromColumns(obj ssot.Object, name string, columns []*ssot.Attribute) *Plain {
	if name == "" {
		panic(fmt.Sprintf("object %q: plain SHOW struct name is empty", obj.Name))
	}
	pairs := StructPair(unexported(name)+"Row", name)
	for _, column := range columns {
		if column == nil {
			panic(fmt.Sprintf("object %q: nil plain SHOW column", obj.Name))
		}
		pairs = appendPairedColumn(pairs, obj, *column)
	}
	return pairs.asPlainStruct()
}

func unexported(name string) string {
	return strings.ToLower(name[:1]) + name[1:]
}

func pairedStructsFromColumns(obj ssot.Object, describe bool, dbName, plainName string) *PairedStructs {
	generators := ssot.ObjectGenerators(obj.Generators...)
	kind := "SHOW"
	enabled := generators.Sdk.Show
	columns := obj.PairedShowColumns()
	if describe {
		kind = "DESCRIBE"
		enabled = generators.Sdk.Describe
		columns = obj.PairedDescribeColumns()
	}
	if !enabled {
		panic(fmt.Sprintf("object %q is not enabled for SDK paired struct %s generation", obj.Name, kind))
	}
	if len(columns) == 0 {
		panic(fmt.Sprintf("object %q is enabled for SDK paired struct %s but has no columns", obj.Name, kind))
	}
	pairs := StructPair(dbName, plainName)
	for _, column := range columns {
		if column == nil {
			panic(fmt.Sprintf("object %q: nil %s column", obj.Name, kind))
		}
		pairs = appendPairedColumn(pairs, obj, *column)
	}
	return pairs
}

func appendPairedColumn(pairs *PairedStructs, obj ssot.Object, attr ssot.Attribute) *PairedStructs {
	key := attr.OutputKey()
	switch attr.Type {
	case ssot.AttrString:
		switch {
		case !attr.OptionalInShow && !attr.OptionalInDb:
			return pairs.Text(key)
		case attr.OptionalInShow && attr.OptionalInDb:
			return pairs.OptionalText(key)
		case !attr.OptionalInShow && attr.OptionalInDb:
			return pairs.OptionalText(key, WithRequiredInPlain())
		default:
			panicPairedColumn(obj, attr, "OptionalInShow without OptionalInDb has no paired struct mapping")
		}
	case ssot.AttrInt:
		if attr.OptionalInShow || attr.OptionalInDb {
			return pairs.OptionalNumber(key, plainRequiredOpt(attr)...)
		}
		return pairs.Number(key)
	case ssot.AttrBool:
		if attr.OptionalInShow || attr.OptionalInDb {
			if attr.StoredAsText {
				return pairs.OptionalBoolFromText(key, plainRequiredOpt(attr)...)
			}
			return pairs.OptionalBool(key, plainRequiredOpt(attr)...)
		}
		if attr.StoredAsText {
			return pairs.BoolFromText(key)
		}
		return pairs.Bool(key)
	case ssot.AttrFloat:
		if attr.OptionalInShow || attr.OptionalInDb {
			panicPairedColumn(obj, attr, "an optional float has no paired struct mapping")
		}
		return pairs.Field(key, "string", "float64", WithManualConvert())
	case ssot.AttrTime:
		switch {
		case !attr.OptionalInShow && !attr.OptionalInDb:
			return pairs.Time(key)
		case !attr.OptionalInShow && attr.OptionalInDb:
			return pairs.OptionalTime(key, WithRequiredInPlain())
		default:
			panicPairedColumn(obj, attr, "this time optionality has no paired struct mapping")
		}
	case ssot.AttrEnum:
		enum := attr.ColumnEnum()
		if attr.OptionalInShow || attr.OptionalInDb {
			opts := plainRequiredOpt(attr)
			if attr.ManualConvert {
				opts = append(opts, WithManualConvert())
			}
			return pairs.OptionalEnum(key, enum, opts...)
		}
		if attr.ManualConvert {
			return pairs.Field(key, "string", enum.Kind(), WithManualConvert())
		}
		return pairs.Enum(key, enum)
	case ssot.AttrIdentifier:
		kind := attr.IdentifierType
		if attr.CustomParser != "" && !attr.OptionalInShow && !attr.OptionalInDb {
			return pairs.Field(key, "string", kind, WithCustomParser(attr.CustomParser))
		}
		if !attr.OptionalInShow || !attr.OptionalInDb {
			panicPairedColumn(obj, attr, "an identifier maps only when it is optional in SHOW and in the DB row")
		}
		return pairs.Field(key, "sql.NullString", "*"+kind, WithPlainFieldName(sqlToFieldName(key, true)))
	case ssot.AttrList:
		if attr.Elem != nil && attr.Elem.Type == ssot.AttrIdentifier && attr.CustomParser != "" && !attr.OptionalInShow && attr.OptionalInDb {
			kind := attr.Elem.IdentifierType
			return pairs.Field(key, "sql.NullString", "[]"+kind, WithCustomParser(attr.CustomParser))
		}
		if attr.OptionalInShow || attr.OptionalInDb || attr.Elem == nil || attr.Elem.Type != ssot.AttrEnum {
			panicPairedColumn(obj, attr, "only a list of enum maps to StringList")
		}
		return pairs.StringList(key)
	case ssot.AttrSet:
		panicPairedColumn(obj, attr, "a set has no paired struct mapping")
	default:
		panicPairedColumn(obj, attr, fmt.Sprintf("type %q has no paired struct mapping", attr.Type))
	}
	return pairs
}

func plainRequiredOpt(attr ssot.Attribute) []PairedFieldOption {
	if !attr.OptionalInShow && attr.OptionalInDb {
		return []PairedFieldOption{WithRequiredInPlain()}
	}
	return nil
}

func panicPairedColumn(obj ssot.Object, attr ssot.Attribute, reason string) {
	panic(fmt.Sprintf("object %q: column %q cannot be a paired struct field: %s", obj.Name, attr.OutputKey(), reason))
}
