package gen

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// ShowPairedStructs builds a SHOW field chain from Columns by calling PairedStructs methods.
// The db row name is the unexported plural plus "Row". The plain name is the object name.
// The query struct, filtering, and description mapping stay in the def.
func ShowPairedStructs(obj ssot.Object) *PairedStructs {
	if obj.Name == "" || obj.NamePlural == "" {
		panic(fmt.Sprintf("object %q needs Name and NamePlural to build the SHOW paired struct", obj.Name))
	}
	return pairedStructsFromColumns(obj, false, unexported(obj.NamePlural)+"Row", obj.Name)
}

// DescribePairedStructs builds a DESCRIBE field chain from Columns by calling PairedStructs methods.
// The db row name is the unexported object name plus "DescRow". The plain name is the object name plus "Details".
// The query struct, filtering, and description mapping stay in the def.
func DescribePairedStructs(obj ssot.Object) *PairedStructs {
	if obj.Name == "" {
		panic(fmt.Sprintf("object %q needs Name to build the DESCRIBE paired struct", obj.Name))
	}
	return pairedStructsFromColumns(obj, true, unexported(obj.Name)+"DescRow", obj.Name+"Details")
}

func unexported(name string) string {
	return strings.ToLower(name[:1]) + name[1:]
}

func pairedStructsFromColumns(obj ssot.Object, describe bool, dbName, plainName string) *PairedStructs {
	generators := ssot.ObjectGenerators(obj.Generators...)
	kind := "SHOW"
	enabled := generators.Sdk.Show
	columns := obj.Show.Columns
	if describe {
		kind = "DESCRIBE"
		enabled = generators.Sdk.Describe
		columns = obj.Describe.Columns
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
			panicPairedColumn(obj, attr, "an optional int has no paired struct mapping")
		}
		return pairs.Number(key)
	case ssot.AttrBool:
		if attr.OptionalInShow || attr.OptionalInDb {
			panicPairedColumn(obj, attr, "an optional bool has no paired struct mapping; BoolFromText is not an ssot type")
		}
		return pairs.Bool(key)
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
		if attr.OptionalInShow || attr.OptionalInDb {
			panicPairedColumn(obj, attr, "an optional enum has no paired struct mapping")
		}
		if attr.Enum == nil {
			panicPairedColumn(obj, attr, "enum has no Enum")
		}
		return pairs.Enum(key, attr.Enum)
	case ssot.AttrIdentifier:
		if !attr.OptionalInShow || !attr.OptionalInDb {
			panicPairedColumn(obj, attr, "an identifier maps only when it is optional in SHOW and in the DB row")
		}
		if obj.IdType == "" {
			panicPairedColumn(obj, attr, "identifier has no IdType")
		}
		return pairs.Field(key, "sql.NullString", "*"+obj.IdType, WithPlainFieldName(sqlToFieldName(key, true)))
	case ssot.AttrList:
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

func panicPairedColumn(obj ssot.Object, attr ssot.Attribute, reason string) {
	panic(fmt.Sprintf("object %q: column %q cannot be a paired struct field: %s", obj.Name, attr.OutputKey(), reason))
}
