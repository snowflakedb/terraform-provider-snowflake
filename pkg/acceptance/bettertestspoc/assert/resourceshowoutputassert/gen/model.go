package gen

import (
	"slices"
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
	IsSdkStruct      bool
	SkipGeneration   bool
}

func ModelFromSdkObjectDetails(sdkObject SdkObjectShowOutputDetails, preamble *genhelpers.PreambleModel) ResourceShowOutputAssertionsModel {
	attributes := make([]ResourceShowOutputAssertionModel, len(sdkObject.Fields))
	for idx, field := range sdkObject.Fields {
		attr := MapToResourceShowOutputAssertion(field)
		if slices.Contains(sdkObject.SkipFields, field.Name) {
			attr.SkipGeneration = true
		}
		attributes[idx] = attr
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
	underlyingTypeWithoutPtr, _ := strings.CutPrefix(field.UnderlyingType, "*")
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
	case strings.HasPrefix(concreteTypeWithoutPtr, "sdk."):
		assertionCreator = "StringValueSet"
		mapper = genhelpers.CastToString
	default:
		assertionCreator = "StringValueSet"
	}

	// isIdentifier tracks whether this sdk type is an identifier type with a dedicated mapper.
	isIdentifier := false
	// TODO [SNOW-1501905]: handle other mappings if needed
	switch concreteTypeWithoutPtr {
	case "sdk.AccountObjectIdentifier":
		mapper = genhelpers.Name
		isIdentifier = true
	case "sdk.ObjectIdentifier", "sdk.AccountIdentifier", "sdk.DatabaseObjectIdentifier", "sdk.SchemaObjectIdentifier", "sdk.SchemaObjectIdentifierWithArguments", "sdk.ExternalObjectIdentifier":
		mapper = genhelpers.FullyQualifiedName
		isIdentifier = true
	case "time.Time":
		mapper = genhelpers.ToString
	}

	// TODO [SNOW-1501905]: currently, assertions for sdk structs and interface types are not properly generated. We mark them to skip them.
	// IsSdkStruct is true for sdk struct types or interface types that cannot be directly cast to string.
	// These are complex types (e.g. sdk.CortexAgentProfile) or interface types (e.g. datatypes.DataType) that require manual handling in _ext.go files.
	// String enums and primitive aliases are kept because their underlying kind is not "struct" or "interface".
	// Identifier types are kept because they have dedicated special mappers (Name/FullyQualifiedName).
	isSdkStruct := (strings.HasPrefix(concreteTypeWithoutPtr, "sdk.") && underlyingTypeWithoutPtr == "struct" && !isIdentifier) || underlyingTypeWithoutPtr == "interface"

	return ResourceShowOutputAssertionModel{
		Name:             field.Name,
		ConcreteType:     concreteTypeWithoutPtr,
		AssertionCreator: assertionCreator,
		Mapper:           mapper,
		IsSdkStruct:      isSdkStruct,
	}
}
