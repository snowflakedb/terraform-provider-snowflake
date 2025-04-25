package gen

import (
	"log"
	"os"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-1501905]: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
}

type ResourceConfigBuilderModel struct {
	Name       string
	Attributes []ResourceConfigBuilderAttributeModel
	PreambleModel
}

func (m ResourceConfigBuilderModel) SomeFunc() {
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

func ModelFromResourceSchemaDetails(resourceSchemaDetails genhelpers.ResourceSchemaDetails) ResourceConfigBuilderModel {
	additionalImports := []string{"encoding/json"}
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
				attrType, additionalImport := handleAttributeTypeForListsAndSets(attr, resourceSchemaDetails.Name)
				attributeType = attrType
				additionalImports = append(additionalImports, additionalImport)
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

	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return ResourceConfigBuilderModel{
		Name:       resourceSchemaDetails.ObjectName(),
		Attributes: attributes,
		PreambleModel: PreambleModel{
			PackageName:               packageWithGenerateDirective,
			AdditionalStandardImports: additionalImports,
		},
	}
}

// handleAttributeTypeForListsAndSets handles model preparation for list and set attributes.
// For simple types it's handled seamlessly.
// For complex types, we need to define override in complexListAttributesOverrides.
// Also, we need to import package (usually sdk) containing the type representing the given object.
func handleAttributeTypeForListsAndSets(attr genhelpers.SchemaAttribute, resourceName string) (attributeType string, additionalImport string) {
	switch attr.AttributeSubType {
	case schema.TypeBool:
		attributeType = "[]bool"
	case schema.TypeInt:
		attributeType = "[]int"
	case schema.TypeFloat:
		attributeType = "[]float"
	case schema.TypeString:
		attributeType = "[]string"
	case schema.TypeMap:
		if v, ok := complexListAttributesOverrides[resourceName]; ok {
			if t, ok := v[attr.Name]; ok {
				attributeType = "[]" + t
				additionalImport = strings.Split(t, ".")[0]
			} else {
				log.Printf("[WARN] No complex list attribute override found for resource's %s attribute %s", resourceName, attr.Name)
			}
		} else {
			log.Printf("[WARN] No complex list attribute overrides found for resource %s", resourceName)
		}
	default:
		log.Printf("[WARN] Attribute's %s sub type could not be determined", attr.Name)
	}
	return
}
