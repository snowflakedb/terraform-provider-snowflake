package gen

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// ShowPairedStructs builds a SHOW field chain from Columns by calling PairedStructs methods.
// The db row name is the unexported plural plus "Row". The plain name is the object name.
// This function builds the field chain only.
func ShowPairedStructs(obj ssot.Object) *PairedStructs {
	if obj.Name == "" || obj.NamePlural == "" {
		panic(fmt.Sprintf("object %q needs Name and NamePlural to build the SHOW paired struct", obj.Name))
	}
	cfg := ssot.ObjectGenerators(obj.Generators...).Sdk
	pairs := pairedStructsFromColumns(obj, false, unexported(obj.NamePlural)+"Row", obj.Name)
	for _, field := range cfg.ShowPlainOnlyFields {
		pairs = pairs.PlainOnlyField(sqlToFieldName(field.Column.OutputKey(), true), field.Kind)
	}
	if cfg.ShowResultFilterHook {
		pairs = pairs.WithShowResultFilterHook()
	}
	return pairs
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
	if cfg.DescribeDetailsRowsField != nil {
		plain = plain.Field(cfg.DescribeDetailsRowsField.Name, cfg.DescribeDetailsRowsField.Kind)
	}
	return plain
}

func detailsFieldKind(obj ssot.Object, attr ssot.Attribute, nested []*ssot.Attribute) string {
	var kind string
	if slices.ContainsFunc(nested, func(n *ssot.Attribute) bool { return n.Name == attr.Name }) {
		kind = "[]" + ssot.NestedDetailsElem(obj, &attr)
	} else {
		kind = detailsGoType(obj, attr)
	}
	if attr.IsOptionalColumn() {
		kind = "*" + kind
	}
	return kind
}

// detailsGoType is the Go type of a column in a details struct. A list or set of enums is a slice of the enum there,
// unlike in a SHOW row, where it is []string.
func detailsGoType(obj ssot.Object, attr ssot.Attribute) string {
	if (attr.Type == ssot.AttrList || attr.Type == ssot.AttrSet) && attr.Elem != nil && attr.Elem.Type == ssot.AttrEnum && attr.Elem.ColumnEnum() != nil {
		return "[]" + attr.Elem.ColumnEnum().Name
	}
	concrete, _ := ssot.AttributeGoTypes(obj, attr)
	return strings.ReplaceAll(concrete, "sdk.", "")
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
			plain = plain.Field(ssot.ColumnFieldName(child), detailsGoType(obj, child))
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
	plainOnly := collections.Map(generators.Sdk.ShowPlainOnlyFields, func(field ssot.SdkShowPlainOnlyField) *ssot.Attribute { return field.Column })
	for _, column := range columns {
		if column == nil {
			panic(fmt.Sprintf("object %q: nil %s column", obj.Name, kind))
		}
		if !describe && slices.Contains(plainOnly, column) {
			continue
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
			return pairs.Text(key, manualConvertOpt(attr)...)
		case attr.OptionalInShow && attr.OptionalInDb:
			return pairs.OptionalText(key, manualConvertOpt(attr)...)
		case !attr.OptionalInShow && attr.OptionalInDb:
			return pairs.OptionalText(key, append(plainRequiredOpt(attr), manualConvertOpt(attr)...)...)
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
				return pairs.OptionalBoolFromText(key, boolFromTextOpts(attr)...)
			}
			return pairs.OptionalBool(key, plainRequiredOpt(attr)...)
		}
		if attr.StoredAsText {
			return pairs.BoolFromText(key, boolFromTextOpts(attr)...)
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
		if !attr.OptionalInShow && !attr.OptionalInDb {
			return pairs.Field(key, "string", kind, manualConvertOpt(attr)...)
		}
		if !attr.OptionalInShow || !attr.OptionalInDb {
			panicPairedColumn(obj, attr, "an identifier maps only when it is required in SHOW and in the DB row, or optional in both")
		}
		return pairs.Field(key, "sql.NullString", "*"+kind, WithPlainFieldName(sqlToFieldName(key, true)))
	case ssot.AttrList:
		if attr.ManualConvert && attr.Elem != nil && attr.Elem.Type == ssot.AttrIdentifier && !attr.OptionalInShow && !attr.OptionalInDb {
			return pairs.Field(key, "string", "[]"+attr.Elem.IdentifierType, WithManualConvert())
		}
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

func boolFromTextOpts(attr ssot.Attribute) []PairedFieldOption {
	if attr.BoolParsed && !attr.StoredAsText {
		panic(fmt.Sprintf("column %q sets BoolParsed without StoredAsText", attr.OutputKey()))
	}
	opts := plainRequiredOpt(attr)
	if attr.BoolParsed {
		opts = append(opts, WithBoolParsed())
	}
	return opts
}

func manualConvertOpt(attr ssot.Attribute) []PairedFieldOption {
	if attr.ManualConvert {
		return []PairedFieldOption{WithManualConvert()}
	}
	return nil
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
