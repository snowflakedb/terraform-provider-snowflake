package gen

import (
	"log"
	"os"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// TODO [SNOW-1501905]: extract to commons?
type PreambleModel struct {
	PackageName               string
	AdditionalStandardImports []string
	AdditionalImports         []string
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
			AdditionalImports:         getAdditionalImports(sdkObject.Fields),
		},
	}
}

func MapToSnowflakeObjectFieldAssertion(field genhelpers.Field) SnowflakeObjectFieldAssertion {
	TypeWithoutPointerAndBrackets := strings.TrimLeft(field.ConcreteType, "*[]")

	mapper := genhelpers.Identity
	if field.IsPointer() {
		mapper = genhelpers.Dereference
	}
	expectedValueMapper := genhelpers.Identity

	// TODO [SNOW-1501905]: handle other mappings if needed
	if TypeWithoutPointerAndBrackets == "sdk.AccountObjectIdentifier" {
		mapper = genhelpers.Name
		if field.IsPointer() {
			mapper = func(s string) string {
				return genhelpers.Name(genhelpers.Parentheses(genhelpers.Dereference(s)))
			}
		}
		expectedValueMapper = genhelpers.Name
	}
	if TypeWithoutPointerAndBrackets == "sdk.SchemaObjectIdentifier" {
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

func getAdditionalImports(fields []genhelpers.Field) []string {
	imports := make(map[string]struct{})
	for _, field := range fields {
		if field.IsSlice() {
			imports["collections"] = struct{}{}
		}
	}
	additionalImports := make([]string, 0)
	for k := range imports {
		if v, ok := genhelpers.PredefinedImports[k]; ok {
			additionalImports = append(additionalImports, v)
		} else {
			log.Printf("[WARN] No predefined import found for %s", k)
		}
	}
	return additionalImports
}
