package gen

import (
	"os"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
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
	IsRequired    bool
}

func ModelFromResourceSchemaDetails(resourceSchemaDetails genhelpers.ResourceSchemaDetails) ResourceAssertionsModel {
	attributes := make([]ResourceAttributeAssertionModel, 0)
	for _, attr := range resourceSchemaDetails.Attributes {
		if slices.Contains([]string{resources.ShowOutputAttributeName, resources.ParametersAttributeName, resources.DescribeOutputAttributeName}, attr.Name) {
			continue
		}
		attributes = append(attributes, ResourceAttributeAssertionModel{
			Name:          attr.Name,
			AttributeType: attr.AttributeType.String(),
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
