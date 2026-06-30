package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type ResourceShowOutputAssertionsModel struct {
	Name             string
	DataSourceName   string
	IsDescribeOutput bool
	Attributes       []ResourceShowOutputAssertionModel

	*genhelpers.PreambleModel
}

type ResourceShowOutputAssertionModel struct {
	Name             string
	ConcreteType     string
	AssertionCreator string
	Mapper           genhelpers.Mapper
}

func ModelFromSdkObjectDetails(sdkObject SdkObjectShowOutputDetails, preamble *genhelpers.PreambleModel) ResourceShowOutputAssertionsModel {
	attributes := make([]ResourceShowOutputAssertionModel, len(sdkObject.Fields))
	for idx, field := range sdkObject.Fields {
		attributes[idx] = MapToResourceShowOutputAssertion(field)
	}

	name, _ := strings.CutPrefix(sdkObject.Name, "sdk.")
	var dataSourceName string
	if sdkObject.dataSourceDef != nil {
		dataSourceName = sdkObject.dataSourceDef.pluralName
	}
	return ResourceShowOutputAssertionsModel{
		Name:             name,
		DataSourceName:   dataSourceName,
		IsDescribeOutput: sdkObject.IsDataSourceOutput,
		Attributes:       attributes,
		PreambleModel:    preamble,
	}
}

func MapToResourceShowOutputAssertion(field genhelpers.Field) ResourceShowOutputAssertionModel {
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")
	mapper := genhelpers.Identity
	// TODO [SNOW-1501905]: get a runtime name for the assertion creator
	var assertionCreator string
	switch {
	case concreteTypeWithoutPtr == "bool":
		assertionCreator = "BoolValueSet"
	case concreteTypeWithoutPtr == "int":
		assertionCreator = "IntValueSet"
	case concreteTypeWithoutPtr == "float64":
		assertionCreator = "FloatValueSet"
	case concreteTypeWithoutPtr == "string":
		assertionCreator = "StringValueSet"
	// TODO [SNOW-1501905]: distinguish between different enum types
	// TODO [SNOW-1501905]: currently, it also generates this assertion type for sdk structs
	case strings.HasPrefix(concreteTypeWithoutPtr, "sdk."):
		assertionCreator = "StringValueSet"
		mapper = genhelpers.CastToString
	default:
		assertionCreator = "StringValueSet"
	}

	// TODO [SNOW-1501905]: handle other mappings if needed
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
