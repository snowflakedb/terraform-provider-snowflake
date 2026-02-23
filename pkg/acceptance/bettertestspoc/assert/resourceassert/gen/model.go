package gen

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceAssertionsModel struct {
	Name       string
	Attributes []ResourceAttributeAssertionModel

	*genhelpers.PreambleModel
}

type ResourceAttributeAssertionModel struct {
	Name         string
	IsCollection bool
	IsRequired   bool

	ExpectedType     string
	AssertionCreator string
}

func ModelFromResourceSchemaDetails(resourceSchemaDetails genhelpers.ResourceSchemaDetails, preamble *genhelpers.PreambleModel) ResourceAssertionsModel {
	attributes := make([]ResourceAttributeAssertionModel, 0)
	for _, attr := range resourceSchemaDetails.Attributes {
		if slices.Contains([]string{resources.ShowOutputAttributeName, resources.ParametersAttributeName, resources.DescribeOutputAttributeName}, attr.Name) {
			continue
		}

		expectedType, assertionCreator := getExpectedTypeAndAssertionCreator(attr)
		attributes = append(attributes, ResourceAttributeAssertionModel{
			Name:         attr.Name,
			IsCollection: attr.AttributeType == schema.TypeList || attr.AttributeType == schema.TypeSet,
			IsRequired:   attr.Required,

			ExpectedType:     expectedType,
			AssertionCreator: assertionCreator,
		})
	}

	return ResourceAssertionsModel{
		Name:          resourceSchemaDetails.ObjectName(),
		Attributes:    attributes,
		PreambleModel: preamble,
	}
}

func getExpectedTypeAndAssertionCreator(attr genhelpers.SchemaAttribute) (expectedType string, assertionCreator string) {
	switch attr.AttributeType {
	case schema.TypeBool:
		expectedType = "bool"
		assertionCreator = "BoolValueSet"
	case schema.TypeInt:
		expectedType = "int"
		assertionCreator = "IntValueSet"
	case schema.TypeFloat:
		expectedType = "float64"
		assertionCreator = "FloatValueSet"
	case schema.TypeString:
		expectedType = "string"
		assertionCreator = "StringValueSet"
	case schema.TypeSet:
		// TODO [SNOW-3113128]: handle/add limitation
	case schema.TypeList:
		// TODO [SNOW-3113128]: handle/add limitation
	case schema.TypeMap:
		// TODO [SNOW-3113128]: handle/add limitation
	case schema.TypeInvalid:
	}
	return
}
