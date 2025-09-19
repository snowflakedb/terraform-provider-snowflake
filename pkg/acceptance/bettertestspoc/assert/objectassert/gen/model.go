package gen

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type SnowflakeObjectAssertionsModel struct {
	Name    string
	SdkType string
	IdType  string
	Fields  []SnowflakeObjectFieldAssertion

	*genhelpers.PreambleModel
}

type SnowflakeObjectFieldAssertion struct {
	Name                  string
	ConcreteType          string
	IsOriginalTypePointer bool
	IsOriginalTypeSlice   bool
	Mapper                genhelpers.Mapper
	ExpectedValueMapper   genhelpers.Mapper
}

func ModelFromSdkObjectDetails(sdkObject genhelpers.SdkObjectDetails, preamble *genhelpers.PreambleModel) SnowflakeObjectAssertionsModel {
	name, _ := strings.CutPrefix(sdkObject.Name, "sdk.")
	fields := make([]SnowflakeObjectFieldAssertion, len(sdkObject.Fields))
	for idx, field := range sdkObject.Fields {
		fields[idx] = MapToSnowflakeObjectFieldAssertion(field)
	}

	return SnowflakeObjectAssertionsModel{
		Name:          name,
		SdkType:       sdkObject.Name,
		IdType:        sdkObject.IdType,
		Fields:        fields,
		PreambleModel: preamble,
	}
}

func MapToSnowflakeObjectFieldAssertion(field genhelpers.Field) SnowflakeObjectFieldAssertion {
	concreteTypeWithoutPtrAndBrackets := field.ConcreteTypeNoPointerNoArray()

	mapper := genhelpers.Identity
	if field.IsPointer() {
		mapper = genhelpers.Dereference
	}
	expectedValueMapper := genhelpers.Identity

	// TODO [SNOW-1501905]: handle other mappings if needed
	if concreteTypeWithoutPtrAndBrackets == "sdk.AccountObjectIdentifier" {
		mapper = genhelpers.Name
		if field.IsPointer() {
			mapper = func(s string) string {
				return genhelpers.Name(genhelpers.Parentheses(genhelpers.Dereference(s)))
			}
		}
		expectedValueMapper = genhelpers.Name
	}
	if concreteTypeWithoutPtrAndBrackets == "sdk.SchemaObjectIdentifier" {
		mapper = genhelpers.FullyQualifiedName
		if field.IsPointer() {
			mapper = func(s string) string {
				return genhelpers.FullyQualifiedName(genhelpers.Parentheses(genhelpers.Dereference(s)))
			}
		}
		expectedValueMapper = genhelpers.FullyQualifiedName
	}

	return SnowflakeObjectFieldAssertion{
		Name:                  field.Name,
		ConcreteType:          field.ConcreteType,
		IsOriginalTypePointer: field.IsPointer(),
		IsOriginalTypeSlice:   field.IsSlice(),
		Mapper:                mapper,
		ExpectedValueMapper:   expectedValueMapper,
	}
}
