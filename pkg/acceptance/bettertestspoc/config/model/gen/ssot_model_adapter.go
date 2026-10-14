package gen

import (
	"fmt"
	"slices"
	"strings"

	sdkdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SdkMethodSuffix names the extra setter that takes the real Go type.
// The existing With method stays and still takes the Terraform storage type, usually string.
// Go cannot overload With, and WithValue is already the arbitrary-value setter, so the real-type method needs its own suffix: Bool for a Go bool, Sdk for an SDK type.
// TODO [next PRs]: migrate ssot-derived model generation to only the new methods.
func SdkMethodSuffix(sdkType string) string {
	if sdkType == "bool" {
		return "Bool"
	}
	return "Sdk"
}

// ModelFromSsotObject builds the model-builder input for one opted-in object.
// Field types come from the ssot attribute. Parameter fields come from the parameter catalog when the object generates them.
func ModelFromSsotObject(obj ssot.Object) ResourceConfigBuilderModel {
	if !ssot.ObjectGenerators(obj.Generators...).Model.Enabled {
		panic(fmt.Sprintf("object %q is not enabled for model generation", obj.Name))
	}
	attrs := identityModelAttributes(obj)
	for _, attr := range obj.Attributes {
		attrs = append(attrs, modelAttributeFromSsot(obj, attr))
	}
	for _, parameter := range parameterDefsForObject(obj) {
		attrs = append(attrs, modelAttributeFromParameter(parameter))
	}
	return ResourceConfigBuilderModel{
		Name: obj.Name,
		Attributes: genhelpers.OrderByNameWithPreferredFront(attrs, func(attr ResourceConfigBuilderAttributeModel) string {
			return attr.Name
		}),
	}
}

func identityModelAttributes(obj ssot.Object) []ResourceConfigBuilderAttributeModel {
	name := plainStringAttribute("name", true)
	database := plainStringAttribute("database", true)
	schemaName := plainStringAttribute("schema", true)
	switch obj.IdType {
	case "AccountObjectIdentifier":
		return []ResourceConfigBuilderAttributeModel{name}
	case "DatabaseObjectIdentifier":
		return []ResourceConfigBuilderAttributeModel{database, name}
	case "SchemaObjectIdentifier":
		return []ResourceConfigBuilderAttributeModel{database, schemaName, name}
	default:
		panic(fmt.Sprintf("object %q: unsupported IdType %q", obj.Name, obj.IdType))
	}
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

func modelAttributeFromSsot(obj ssot.Object, attr ssot.Attribute) ResourceConfigBuilderAttributeModel {
	if slices.Contains(excludedAttributeNames, attr.Name) {
		panic(fmt.Sprintf("object %q: attribute %q is a generated output block", obj.Name, attr.Name))
	}
	switch attr.Type {
	case ssot.AttrInt:
		return typedAttribute(attr.Name, "int", "IntegerVariable", schema.TypeInt, attr.Required, "", "")
	case ssot.AttrBool:
		if attr.Required {
			return typedAttribute(attr.Name, "bool", "BoolVariable", schema.TypeBool, true, "", "")
		}
		return typedAttribute(attr.Name, "string", "StringVariable", schema.TypeString, false, "bool", "bool")
	case ssot.AttrEnum:
		if attr.Enum == nil || attr.Enum.Name == "" {
			panic(fmt.Sprintf("object %q: enum attribute %q has no Enum", obj.Name, attr.Name))
		}
		return typedAttribute(attr.Name, "string", "StringVariable", schema.TypeString, attr.Required, "sdk."+attr.Enum.Name, "enum")
	case ssot.AttrIdentifier:
		if obj.IdType == "" {
			panic(fmt.Sprintf("object %q: identifier attribute %q has no IdType", obj.Name, attr.Name))
		}
		return typedAttribute(attr.Name, "string", "StringVariable", schema.TypeString, attr.Required, "sdk."+obj.IdType, "identifier")
	case ssot.AttrString:
		return plainStringAttribute(attr.Name, attr.Required)
	case ssot.AttrList:
		return typedAttribute(attr.Name, "", "", schema.TypeList, attr.Required, "", "")
	case ssot.AttrSet:
		return typedAttribute(attr.Name, "", "", schema.TypeSet, attr.Required, "", "")
	default:
		panic(fmt.Sprintf("object %q: attribute %q has unknown type %q", obj.Name, attr.Name, attr.Type))
	}
}

func modelAttributeFromParameter(parameter parameterdefs.ParameterDef) ResourceConfigBuilderAttributeModel {
	switch parameter.Kind {
	case "bool":
		return typedAttribute(parameter.FieldName(), "bool", "BoolVariable", schema.TypeBool, false, "", "")
	case "int":
		return typedAttribute(parameter.FieldName(), "int", "IntegerVariable", schema.TypeInt, false, "", "")
	case "string", "StringAllowEmpty", "float64":
		return plainStringAttribute(parameter.FieldName(), false)
	}
	if strings.HasSuffix(parameter.Kind, "Identifier") {
		return typedAttribute(parameter.FieldName(), "string", "StringVariable", schema.TypeString, false, "sdk."+parameter.Kind, "identifier")
	}
	return typedAttribute(parameter.FieldName(), "string", "StringVariable", schema.TypeString, false, "sdk."+parameter.Kind, "enum")
}

func plainStringAttribute(name string, required bool) ResourceConfigBuilderAttributeModel {
	return typedAttribute(name, "string", "StringVariable", schema.TypeString, required, "", "")
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
		SdkType:        sdkType,
		SdkKind:        sdkKind,
	}
}
