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

type ResourceAssertionsModel struct {
	Name       string
	Attributes []ResourceAttributeAssertionModel
	PreambleModel
}

func (m ResourceAssertionsModel) SomeFunc() {
}

type ResourceAttributeAssertionModel struct {
	Name          string
	AttributeType string
	IsCollection  bool
	IsRequired    bool
}

func ModelFromResourceSchemaDetails(resourceSchemaDetails genhelpers.ResourceSchemaDetails) ResourceAssertionsModel {
	attributes := make([]ResourceAttributeAssertionModel, 0)
	for _, attr := range resourceSchemaDetails.Attributes {
		if slices.Contains([]string{resources.ShowOutputAttributeName, resources.ParametersAttributeName, resources.DescribeOutputAttributeName}, attr.Name) {
			continue
		}
		attributes = append(attributes, ResourceAttributeAssertionModel{
			Name: attr.Name,
			// TODO [SNOW-1501905]: add attribute type logic; allow type safe assertions
			AttributeType: "string",
			IsCollection:  attr.AttributeType == schema.TypeList || attr.AttributeType == schema.TypeSet,
			IsRequired:    attr.Required,
		})
	}

	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return ResourceAssertionsModel{
		Name:       resourceSchemaDetails.ObjectName(),
		Attributes: attributes,
		PreambleModel: PreambleModel{
			PackageName: packageWithGenerateDirective,
		},
	}
}
