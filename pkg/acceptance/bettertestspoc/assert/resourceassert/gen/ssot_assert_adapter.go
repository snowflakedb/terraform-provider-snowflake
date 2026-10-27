package gen

import (
	"fmt"
	"strings"

	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AssertionsFromSsotObject builds the assertion model for one opted-in object.
// Field types come from the ssot attribute. Parameter fields come from the parameter catalog when the object generates them.
func AssertionsFromSsotObject(obj ssot.Object) ResourceAssertionsModel {
	return assertionsFromSsotObject(obj, parameterDefsForObject(obj))
}

func assertionsFromSsotObject(obj ssot.Object, parameters []ssot.ParameterDef) ResourceAssertionsModel {
	if !ssot.ObjectGenerators(obj.Generators...).ResourceAssert.Enabled {
		panic(fmt.Sprintf("object %q is not enabled for resource assertion generation", obj.Name))
	}
	identityFields := assertionAttributesFromFields(ssot.IdentityFields(obj))
	attributeFields := collections.Map(obj.ResourceAttributes(), func(attr ssot.Attribute) ResourceAttributeAssertionModel {
		return assertionAttributeFromField(ssot.ClassifyAttribute(obj, attr))
	})
	parameterFields := collections.Map(parameters, func(parameter ssot.ParameterDef) ResourceAttributeAssertionModel {
		return assertionAttributeFromField(ssot.ClassifyParameter(parameter.FieldName(), parameter.Kind))
	})
	fullyQualifiedNameFields := []ResourceAttributeAssertionModel{
		plainStringAssertion(resources.FullyQualifiedNameAttributeName, false),
	}
	attrs := make([]ResourceAttributeAssertionModel, 0, len(identityFields)+len(attributeFields)+len(parameterFields)+len(fullyQualifiedNameFields))
	attrs = append(attrs, identityFields...)
	attrs = append(attrs, attributeFields...)
	for _, attr := range obj.ComputedAttributes {
		attrs = append(attrs, assertionAttributeFromField(ssot.ClassifyAttribute(obj, attr)))
	}
	attrs = append(attrs, parameterFields...)
	attrs = append(attrs, fullyQualifiedNameFields...)
	return ResourceAssertionsModel{
		Name: obj.Name,
		Attributes: genhelpers.OrderByNameWithPreferredFront(attrs, func(attr ResourceAttributeAssertionModel) string {
			return attr.Name
		}),
	}
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

func assertionAttributesFromFields(fields []ssot.ClassifiedField) []ResourceAttributeAssertionModel {
	return collections.Map(fields, assertionAttributeFromField)
}

func assertionAttributeFromField(field ssot.ClassifiedField) ResourceAttributeAssertionModel {
	return assertionAttributeFromSchema(field.Name, genhelpers.SchemaAttribute{
		Name:             field.Name,
		AttributeType:    schemaType(field.Storage),
		AttributeSubType: schemaType(field.Elem),
		Required:         field.Required,
	}, field.SdkType, field.SdkKind)
}

func plainStringAssertion(name string, required bool) ResourceAttributeAssertionModel {
	return assertionAttributeFromField(ssot.ClassifiedField{Name: name, Storage: ssot.AttrString, Required: required})
}

func schemaType(storage ssot.AttributeType) schema.ValueType {
	switch storage {
	case ssot.AttrBool:
		return schema.TypeBool
	case ssot.AttrInt:
		return schema.TypeInt
	case ssot.AttrString:
		return schema.TypeString
	case ssot.AttrList:
		return schema.TypeList
	case ssot.AttrSet:
		return schema.TypeSet
	case "":
		return schema.TypeInvalid
	default:
		panic(fmt.Sprintf("unsupported storage %q", storage))
	}
}

func assertionAttributeFromSchema(name string, attr genhelpers.SchemaAttribute, sdkType, sdkKind string) ResourceAttributeAssertionModel {
	expectedType, assertionCreator := getExpectedTypeAndAssertionCreator(attr)
	var subType string
	if attr.AttributeSubType != schema.TypeInvalid {
		subType = genhelpers.ResourceSchemaTypeToString(attr.AttributeSubType)
	}
	return ResourceAttributeAssertionModel{
		Name:                         name,
		IsCollection:                 attr.AttributeType == schema.TypeList || attr.AttributeType == schema.TypeSet,
		IsRequired:                   attr.Required,
		Type:                         genhelpers.ResourceSchemaTypeToString(attr.AttributeType),
		SubType:                      subType,
		ExpectedType:                 expectedType,
		AssertionCreator:             assertionCreator,
		ShouldGenerateTypedAssertion: assertionCreator != "",
		SdkValueType:                 genhelpers.SdkValueType{SdkType: sdkType, SdkKind: sdkKind},
	}
}
