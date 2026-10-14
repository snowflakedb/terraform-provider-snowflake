package gen

import (
	"fmt"
	"strings"

	sdkdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ModelFromSsotObject builds the model-builder input for one opted-in object.
// Field types come from the ssot attribute. Parameter fields come from the parameter catalog when the object generates them.
func ModelFromSsotObject(obj ssot.Object) ResourceConfigBuilderModel {
	if !ssot.ObjectGenerators(obj.Generators...).Model.Enabled {
		panic(fmt.Sprintf("object %q is not enabled for model generation", obj.Name))
	}
	attrs := identityModelAttributes(obj)
	attrs = append(attrs, collections.Map(obj.Attributes, func(attr ssot.Attribute) ResourceConfigBuilderAttributeModel {
		return modelAttributeFromField(ssot.ClassifyAttribute(obj, attr))
	})...)
	attrs = append(attrs, collections.Map(parameterDefsForObject(obj), func(parameter parameterdefs.ParameterDef) ResourceConfigBuilderAttributeModel {
		return modelAttributeFromField(ssot.ClassifyParameter(parameter.FieldName(), parameter.Kind))
	})...)
	return ResourceConfigBuilderModel{
		Name: obj.Name,
		Attributes: genhelpers.OrderByNameWithPreferredFront(attrs, func(attr ResourceConfigBuilderAttributeModel) string {
			return attr.Name
		}),
	}
}

func identityModelAttributes(obj ssot.Object) []ResourceConfigBuilderAttributeModel {
	return modelAttributesFromFields(ssot.IdentityFields(obj))
}

func parameterDefsForObject(obj ssot.Object) []parameterdefs.ParameterDef {
	if !ssot.DefaultResourceExtHooks(obj.ResourceHooks...).ParametersAttributesSchema.Generate {
		return nil
	}
	level := parameterdefs.ParameterLevel(strings.ToUpper(genhelpers.ToSnakeCase(obj.Name)))
	// TODO [next PRs]: parameter definitions must be moved to ssot and referenced from there.
	parameters := sdkdefs.ParameterDefsForLevel(level)
	if len(parameters) == 0 {
		panic(fmt.Sprintf("object %q: parameter attributes are generated but level %q has no parameters", obj.Name, level))
	}
	return parameters
}

func modelAttributesFromFields(fields []ssot.ClassifiedField) []ResourceConfigBuilderAttributeModel {
	return collections.Map(fields, modelAttributeFromField)
}

func modelAttributeFromField(field ssot.ClassifiedField) ResourceConfigBuilderAttributeModel {
	switch field.Storage {
	case ssot.AttrBool:
		return typedAttribute(field.Name, "bool", "BoolVariable", schema.TypeBool, field.Required, field.SdkType, field.SdkKind)
	case ssot.AttrInt:
		return typedAttribute(field.Name, "int", "IntegerVariable", schema.TypeInt, field.Required, field.SdkType, field.SdkKind)
	case ssot.AttrString:
		return typedAttribute(field.Name, "string", "StringVariable", schema.TypeString, field.Required, field.SdkType, field.SdkKind)
	case ssot.AttrList:
		return typedAttribute(field.Name, "", "", schema.TypeList, field.Required, field.SdkType, field.SdkKind)
	case ssot.AttrSet:
		return typedAttribute(field.Name, "", "", schema.TypeSet, field.Required, field.SdkType, field.SdkKind)
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
