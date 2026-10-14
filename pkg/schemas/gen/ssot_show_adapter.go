package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

func showResultSchemaDetailsFromSsot(obj ssot.Object, isDescribe bool) ShowResultSchemaDetails {
	name := "sdk." + obj.Name
	attrs := obj.ShowOutputAttributes()
	if isDescribe {
		name += "Details"
		attrs = obj.DescribeOutputAttributes()
	}
	fields := make([]genhelpers.Field, 0, len(attrs))
	for _, attr := range attrs {
		if attr == nil {
			kind := "SHOW"
			if isDescribe {
				kind = "DESCRIBE"
			}
			panic(fmt.Sprintf("object %q: nil %s column", obj.Name, kind))
		}
		fields = append(fields, ssotAttributeField(obj, *attr))
	}
	return ShowResultSchemaDetails{
		IsDescribe:    isDescribe,
		StructDetails: genhelpers.StructDetails{Name: name, Fields: fields},
	}
}

func ssotAttributeField(obj ssot.Object, attr ssot.Attribute) genhelpers.Field {
	concrete, underlying := ssotGoTypes(obj, attr)
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

func ssotGoTypes(obj ssot.Object, attr ssot.Attribute) (string, string) {
	switch attr.Type {
	case ssot.AttrString:
		return "string", "string"
	case ssot.AttrInt:
		return "int", "int"
	case ssot.AttrBool:
		return "bool", "bool"
	case ssot.AttrTime:
		return "time.Time", "struct"
	case ssot.AttrIdentifier:
		if obj.IdType == "" {
			panic(fmt.Sprintf("object %q: identifier attribute %q has no IdType", obj.Name, attr.Name))
		}
		return "sdk." + obj.IdType, "struct"
	case ssot.AttrEnum:
		if attr.Enum == nil || attr.Enum.Name == "" {
			panic(fmt.Sprintf("object %q: enum attribute %q has no Enum", obj.Name, attr.Name))
		}
		return "sdk." + attr.Enum.Name, "string"
	case ssot.AttrList, ssot.AttrSet:
		if attr.Elem == nil {
			panic(fmt.Sprintf("object %q: collection attribute %q has no Elem", obj.Name, attr.Name))
		}
		if attr.Elem.Type == ssot.AttrEnum {
			return "[]string", "slice"
		}
		elemConcrete, _ := ssotGoTypes(obj, *attr.Elem)
		return "[]" + elemConcrete, "slice"
	default:
		panic(fmt.Sprintf("object %q: attribute %q has unknown type %q", obj.Name, attr.Name, attr.Type))
	}
}
