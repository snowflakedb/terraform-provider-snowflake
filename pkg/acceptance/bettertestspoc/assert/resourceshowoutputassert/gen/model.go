package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type ResourceShowOutputAssertionsModel struct {
	Name       string
	Attributes []ResourceShowOutputAssertionModel

	genhelpers.PreambleModel
}

type ResourceShowOutputAssertionModel struct {
	Name             string
	ConcreteType     string
	AssertionCreator string
	Mapper           genhelpers.Mapper
}

func ModelFromSdkObjectDetails(sdkObject genhelpers.SdkObjectDetails, preamble genhelpers.PreambleModel) ResourceShowOutputAssertionsModel {
	attributes := make([]ResourceShowOutputAssertionModel, len(sdkObject.Fields))
	for idx, field := range sdkObject.Fields {
		attributes[idx] = MapToResourceShowOutputAssertion(field)
	}

	name, _ := strings.CutPrefix(sdkObject.Name, "sdk.")
	return ResourceShowOutputAssertionsModel{
		Name:          name,
		Attributes:    attributes,
		PreambleModel: preamble,
	}
}

func MapToResourceShowOutputAssertion(field genhelpers.Field) ResourceShowOutputAssertionModel {
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")
	// TODO [SNOW-1501905]: get a runtime name for the assertion creator
	var assertionCreator string
	switch {
	case concreteTypeWithoutPtr == "bool":
		assertionCreator = "ResourceShowOutputBoolValue"
	case concreteTypeWithoutPtr == "int":
		assertionCreator = "ResourceShowOutputIntValue"
	case concreteTypeWithoutPtr == "float64":
		assertionCreator = "ResourceShowOutputFloatValue"
	case concreteTypeWithoutPtr == "string":
		assertionCreator = "ResourceShowOutputValue"
	// TODO [SNOW-1501905]: distinguish between different enum types
	case strings.HasPrefix(concreteTypeWithoutPtr, "sdk."):
		assertionCreator = "ResourceShowOutputStringUnderlyingValue"
	default:
		assertionCreator = "ResourceShowOutputValue"
	}

	// TODO [SNOW-1501905]: handle other mappings if needed
	mapper := genhelpers.Identity
	switch concreteTypeWithoutPtr {
	case "sdk.AccountObjectIdentifier":
		mapper = genhelpers.Name
	case "sdk.ObjectIdentifier", "sdk.AccountIdentifier", "sdk.DatabaseObjectIdentifier", "sdk.SchemaObjectIdentifier", "sdk.SchemaObjectIdentifierWithArguments", "sdk.ExternalObjectIdentifier":
		mapper = genhelpers.FullyQualifiedName
	case "time.Time":
		mapper = genhelpers.ToString
	}

	return ResourceShowOutputAssertionModel{
		Name:             field.Name,
		ConcreteType:     concreteTypeWithoutPtr,
		AssertionCreator: assertionCreator,
		Mapper:           mapper,
	}
}
