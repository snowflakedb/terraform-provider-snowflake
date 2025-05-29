package gen

import (
	"os"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// TODO [SNOW-1501905]: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
}

type SnowflakeObjectAssertionsModel struct {
	Name    string
	SdkType string
	IdType  string
	Fields  []SnowflakeObjectFieldAssertion
	PreambleModel
}

func (m SnowflakeObjectAssertionsModel) SomeFunc() {
}

type SnowflakeObjectFieldAssertion struct {
	Name                  string
	ConcreteType          string
	IsOriginalTypePointer bool
	IsOriginalTypeSlice   bool
	Mapper                genhelpers.Mapper
	ExpectedValueMapper   genhelpers.Mapper
}

func ModelFromSdkObjectDetails(sdkObject genhelpers.SdkObjectDetails) SnowflakeObjectAssertionsModel {
	name, _ := strings.CutPrefix(sdkObject.Name, "sdk.")
	fields := make([]SnowflakeObjectFieldAssertion, len(sdkObject.Fields))
	for idx, field := range sdkObject.Fields {
		fields[idx] = MapToSnowflakeObjectFieldAssertion(field)
	}

	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return SnowflakeObjectAssertionsModel{
		Name:    name,
		SdkType: sdkObject.Name,
		IdType:  sdkObject.IdType,
		Fields:  fields,
		PreambleModel: PreambleModel{
			PackageName:               packageWithGenerateDirective,
			AdditionalStandardImports: genhelpers.AdditionalStandardImports(sdkObject.Fields),
		},
	}
}

func MapToSnowflakeObjectFieldAssertion(field genhelpers.Field) SnowflakeObjectFieldAssertion {
	concreteTypeWithoutPtr, _ := strings.CutPrefix(field.ConcreteType, "*")

	mapper := genhelpers.Identity
	if field.IsPointer() {
		mapper = genhelpers.Dereference
	}
	expectedValueMapper := genhelpers.Identity

	// TODO [SNOW-1501905]: handle other mappings if needed
	if concreteTypeWithoutPtr == "sdk.AccountObjectIdentifier" {
		mapper = genhelpers.Name
		if field.IsPointer() {
			mapper = func(s string) string {
				return genhelpers.Name(genhelpers.Parentheses(genhelpers.Dereference(s)))
			}
		}
		expectedValueMapper = genhelpers.Name
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
