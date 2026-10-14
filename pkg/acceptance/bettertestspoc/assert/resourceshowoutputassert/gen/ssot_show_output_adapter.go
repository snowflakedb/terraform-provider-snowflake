package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// ShowOutputAssertionsFromSsot builds the show or describe assertion model for one opted-in object.
// Columns become genhelpers.Field values, then MapToResourceShowOutputAssertion, the same mapper the reflected path uses.
func ShowOutputAssertionsFromSsot(obj ssot.Object, isDescribe bool) ResourceShowOutputAssertionsModel {
	return ModelFromSdkObjectDetails(showOutputDetailsFromSsot(obj, isDescribe), nil)
}

func showOutputDetailsFromSsot(obj ssot.Object, isDescribe bool) SdkObjectShowOutputDetails {
	g := ssot.ObjectGenerators(obj.Generators...)
	kind := "SHOW"
	enabled := g.ShowOutputAssert.Show
	attrs := obj.ShowOutputAttributes()
	if isDescribe {
		kind = "DESCRIBE"
		enabled = g.ShowOutputAssert.Describe
		attrs = obj.DescribeOutputAttributes()
	}
	if !enabled {
		panic(fmt.Sprintf("object %q is not enabled for show output assertion %s generation", obj.Name, kind))
	}
	if len(attrs) == 0 {
		panic(fmt.Sprintf("object %q is enabled for show output assertion %s but has no %s fields", obj.Name, kind, kind))
	}
	name := "sdk." + obj.Name
	if isDescribe {
		name += "Details"
	}
	fields := make([]genhelpers.Field, 0, len(attrs))
	for _, attr := range attrs {
		if attr == nil {
			panic(fmt.Sprintf("object %q: nil %s column", obj.Name, kind))
		}
		fields = append(fields, showOutputAttributeField(obj, *attr))
	}
	dataSource, _ := dataSourceMappingNormalized[name]
	return SdkObjectShowOutputDetails{
		dataSourceDef: &dataSource,
		SdkObjectDetails: genhelpers.SdkObjectDetails{
			IsDataSourceOutput: isDescribe,
			StructDetails: genhelpers.StructDetails{
				Name:   name,
				Fields: fields,
			},
		},
	}
}

func showOutputAttributeField(obj ssot.Object, attr ssot.Attribute) genhelpers.Field {
	concrete, underlying := showOutputGoTypes(obj, attr)
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

func showOutputGoTypes(obj ssot.Object, attr ssot.Attribute) (string, string) {
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
		elemConcrete, _ := showOutputGoTypes(obj, *attr.Elem)
		return "[]" + elemConcrete, "slice"
	default:
		panic(fmt.Sprintf("object %q: attribute %q has unknown type %q", obj.Name, attr.Name, attr.Type))
	}
}
