package gen

import (
	"fmt"
	"strings"

	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ModelFromSsotObject builds the model-builder input for one opted-in object.
// Field types come from the ssot attribute. Parameter fields come from the parameter catalog when the object generates them.
func ModelFromSsotObject(obj ssot.Object) ResourceConfigBuilderModel {
	g := ssot.ObjectGenerators(obj.Generators...)
	if !g.Model.Enabled {
		panic(fmt.Sprintf("object %q is not enabled for model generation", obj.Name))
	}
	attrs := identityModelAttributes(obj)
	attrs = append(attrs, collections.Map(obj.ResourceAttributes(), func(attr ssot.Attribute) ResourceConfigBuilderAttributeModel {
		field := modelAttributeFromField(ssot.ClassifyAttribute(obj, attr), g.Model.RequestTypeOverrides[attr.Name])
		if modelMultiline(g.Model.Multiline, attr.Name) {
			if attr.Type != ssot.AttrString {
				panic("multiline model attribute must be a string")
			}
			field.VariableMethod, field.MethodImport = "MultilineWrapperVariable", "config"
		}
		return field
	})...)
	attrs = append(attrs, collections.Map(parameterDefsForObject(obj), func(parameter ssot.ParameterDef) ResourceConfigBuilderAttributeModel {
		return modelAttributeFromField(ssot.ClassifyParameter(parameter.FieldName(), parameter.Kind), "")
	})...)
	return ResourceConfigBuilderModel{
		Name: obj.Name,
		Attributes: genhelpers.OrderByNameWithPreferredFront(attrs, func(attr ResourceConfigBuilderAttributeModel) string {
			return attr.Name
		}),
	}
}

func modelMultiline(multiline []*ssot.Attribute, name string) bool {
	for _, attr := range multiline {
		if attr != nil && attr.Name == name {
			return true
		}
	}
	return false
}

func identityModelAttributes(obj ssot.Object) []ResourceConfigBuilderAttributeModel {
	return collections.Map(ssot.IdentityFields(obj), func(field ssot.ClassifiedField) ResourceConfigBuilderAttributeModel {
		return modelAttributeFromField(field, "")
	})
}

func parameterDefsForObject(obj ssot.Object) []ssot.ParameterDef {
	if !ssot.DefaultResourceExtHooks(obj.ResourceHooks...).ParametersAttributesSchema.Generate {
		return nil
	}
	level := ssot.ParameterLevel(strings.ToUpper(genhelpers.ToSnakeCase(obj.Name)))
	parameters := ssotdefs.ParameterDefsForLevel(level)
	if len(parameters) == 0 {
		panic(fmt.Sprintf("object %q: parameter attributes are generated but level %q has no parameters", obj.Name, level))
	}
	return parameters
}

// modelAttributeFromField builds one model attribute from a classified field. requestTypeOverride
// is the ModelRequestType value for this field's Name, or "" when none was configured; it only
// applies to a Required AttrList/AttrSet field with no Elem (a nested block with no mechanical type).
func modelAttributeFromField(field ssot.ClassifiedField, requestTypeOverride string) ResourceConfigBuilderAttributeModel {
	switch field.Storage {
	case ssot.AttrBool:
		return typedAttribute(field.Name, "bool", "BoolVariable", schema.TypeBool, field.Required, field.SdkType, field.SdkKind)
	case ssot.AttrInt:
		return typedAttribute(field.Name, "int", "IntegerVariable", schema.TypeInt, field.Required, field.SdkType, field.SdkKind)
	case ssot.AttrString:
		return typedAttribute(field.Name, "string", "StringVariable", schema.TypeString, field.Required, field.SdkType, field.SdkKind)
	case ssot.AttrList:
		return typedAttribute(field.Name, requestTypeOverride, "", schema.TypeList, field.Required, field.SdkType, field.SdkKind)
	case ssot.AttrSet:
		return typedAttribute(field.Name, requestTypeOverride, "", schema.TypeSet, field.Required, field.SdkType, field.SdkKind)
	default:
		panic(fmt.Sprintf("attribute %q has unsupported storage %q", field.Name, field.Storage))
	}
}

func typedAttribute(name, attributeType, variableMethod string, original schema.ValueType, required bool, sdkType, sdkKind string) ResourceConfigBuilderAttributeModel {
	return ResourceConfigBuilderAttributeModel{
		Name:           name,
		JsonName:       name,
		AttributeType:  attributeType,
		Required:       required,
		VariableMethod: variableMethod,
		MethodImport:   "tfconfig",
		OriginalType:   original,
		SdkValueType:   genhelpers.SdkValueType{SdkType: sdkType, SdkKind: sdkKind},
	}
}
