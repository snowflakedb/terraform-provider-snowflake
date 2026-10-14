package ssot

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// AttributeGoTypes returns the concrete Go type and underlying kind for one column,
// without the OptionalInShow pointer. ColumnField applies that pointer.
func AttributeGoTypes(obj Object, attr Attribute) (concrete, underlying string) {
	switch attr.Type {
	case AttrString:
		return "string", "string"
	case AttrInt:
		return "int", "int"
	case AttrBool:
		return "bool", "bool"
	case AttrTime:
		return "time.Time", "struct"
	case AttrIdentifier:
		if obj.IdType == "" {
			panic(fmt.Sprintf("object %q: identifier attribute %q has no IdType", obj.Name, attr.Name))
		}
		return "sdk." + obj.IdType, "struct"
	case AttrEnum:
		if attr.Enum == nil || attr.Enum.Name == "" {
			panic(fmt.Sprintf("object %q: enum attribute %q has no Enum", obj.Name, attr.Name))
		}
		return "sdk." + attr.Enum.Name, "string"
	case AttrList, AttrSet:
		if attr.Elem == nil {
			panic(fmt.Sprintf("object %q: collection attribute %q has no Elem", obj.Name, attr.Name))
		}
		if attr.Elem.Type == AttrEnum {
			return "[]string", "slice"
		}
		elemConcrete, _ := AttributeGoTypes(obj, *attr.Elem)
		return "[]" + elemConcrete, "slice"
	default:
		panic(fmt.Sprintf("object %q: attribute %q has unknown type %q", obj.Name, attr.Name, attr.Type))
	}
}

// ColumnField is the field ExtractStructDetails would report for one column.
// The name comes from OutputKey. OptionalInShow prefixes "*" on both type strings.
func ColumnField(obj Object, attr Attribute) genhelpers.Field {
	concrete, underlying := AttributeGoTypes(obj, attr)
	if attr.OptionalInShow {
		concrete = "*" + concrete
		underlying = "*" + underlying
	}
	return genhelpers.Field{
		Name:           genhelpers.SnakeCaseToCamel(attr.OutputKey()),
		ConcreteType:   concrete,
		UnderlyingType: underlying,
	}
}
