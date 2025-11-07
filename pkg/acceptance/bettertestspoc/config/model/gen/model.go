package gen

import (
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceConfigBuilderModel struct {
	Name       string
	Attributes []ResourceConfigBuilderAttributeModel

	*genhelpers.PreambleModel
}

type ResourceConfigBuilderAttributeModel struct {
	Name           string
	JsonName       string
	AttributeType  string
	Required       bool
	VariableMethod string
	MethodImport   string
	OriginalType   schema.ValueType
}

func ModelFromResourceSchemaDetails(resourceSchemaDetails genhelpers.ResourceSchemaDetails, preamble *genhelpers.PreambleModel) ResourceConfigBuilderModel {
	attributes := make([]ResourceConfigBuilderAttributeModel, 0)
	for _, attr := range resourceSchemaDetails.Attributes {
		if slices.Contains([]string{resources.ShowOutputAttributeName, resources.ParametersAttributeName, resources.DescribeOutputAttributeName}, attr.Name) {
			continue
		}
		jsonName := attr.Name
		name := genhelpers.SanitizeAttributeName(attr.Name)

		if v, ok := multilineAttributesOverrides[resourceSchemaDetails.Name]; ok && slices.Contains(v, attr.Name) && attr.AttributeType == schema.TypeString {
			attributes = append(attributes, ResourceConfigBuilderAttributeModel{
				Name:           name,
				JsonName:       jsonName,
				AttributeType:  "string",
				Required:       attr.Required,
				VariableMethod: "MultilineWrapperVariable",
				MethodImport:   "config",
				OriginalType:   attr.AttributeType,
			})
			continue
		}

		// TODO [SNOW-1501905]: support the rest of attribute types
		var attributeType string
		var variableMethod string
		switch attr.AttributeType {
		case schema.TypeBool:
			attributeType = "bool"
			variableMethod = "BoolVariable"
		case schema.TypeInt:
			attributeType = "int"
			variableMethod = "IntegerVariable"
		case schema.TypeFloat:
			attributeType = "float"
			variableMethod = "FloatVariable"
		case schema.TypeString:
			attributeType = "string"
			variableMethod = "StringVariable"
		case schema.TypeList, schema.TypeSet:
			// We only run it for the required attributes because the `With` methods are not yet generated; we don't need to set the `variableMethod`.
			// For now, the `With` method for complex object will still need to be added to _ext file.
			if attr.Required {
				attributeType = handleAttributeTypeForListsAndSets(attr, resourceSchemaDetails.Name)
			}
		}

		attributes = append(attributes, ResourceConfigBuilderAttributeModel{
			Name:           name,
			JsonName:       jsonName,
			AttributeType:  attributeType,
			Required:       attr.Required,
			VariableMethod: variableMethod,
			MethodImport:   "tfconfig",
			OriginalType:   attr.AttributeType,
		})
	}

	return ResourceConfigBuilderModel{
		Name:          resourceSchemaDetails.ObjectName(),
		Attributes:    attributes,
		PreambleModel: preamble,
	}
}

// handleAttributeTypeForListsAndSets handles model preparation for list and set attributes.
// For simple types it's handled seamlessly.
// For complex types, we need to define override in complexListAttributesOverrides.
// Also, we need to import package (usually sdk) containing the type representing the given object.
func handleAttributeTypeForListsAndSets(attr genhelpers.SchemaAttribute, resourceName string) string {
	var attributeType string
	switch attr.AttributeSubType {
	case schema.TypeBool:
		attributeType = "[]bool"
	case schema.TypeInt:
		attributeType = "[]int"
	case schema.TypeFloat:
		attributeType = "[]float"
	case schema.TypeString:
		attributeType = handleListTypeOverrides(resourceName, attr.Name)
		if attributeType == "" {
			attributeType = "[]string"
		}
	case schema.TypeMap:
		attributeType = handleListTypeOverrides(resourceName, attr.Name)
	default:
		log.Printf("[WARN] Attribute's %s sub type could not be determined", attr.Name)
	}
	return attributeType
}

// TODO [SNOW-1501905]: handle attribute overriding in one place
func handleListTypeOverrides(resourceName string, attrName string) string {
	var attributeType string
	if v, ok := complexListAttributesOverrides[resourceName]; ok {
		if t, ok := v[attrName]; ok {
			attributeType = "[]" + t
		} else {
			log.Printf("[WARN] No complex list attribute override found for resource's %s attribute %s", resourceName, attrName)
		}
	} else {
		log.Printf("[WARN] No complex list attribute overrides found for resource %s", resourceName)
	}
	return attributeType
}
