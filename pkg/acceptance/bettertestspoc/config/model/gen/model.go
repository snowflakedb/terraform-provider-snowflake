package gen

import (
	"os"
	"slices"

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
}

func ModelFromResourceSchemaDetails(resourceSchemaDetails genhelpers.ResourceSchemaDetails) ResourceConfigBuilderModel {
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
		}

		attributes = append(attributes, ResourceConfigBuilderAttributeModel{
			Name:           name,
			JsonName:       jsonName,
			AttributeType:  attributeType,
			Required:       attr.Required,
			VariableMethod: variableMethod,
			MethodImport:   "tfconfig",
		})
	}

	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return ResourceConfigBuilderModel{
		Name:       resourceSchemaDetails.ObjectName(),
		Attributes: attributes,
		PreambleModel: PreambleModel{
			PackageName:               packageWithGenerateDirective,
			AdditionalStandardImports: []string{"encoding/json"},
		},
	}
}
