//go:build exclude

package main

import (
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

const (
	name    = "SDK builder"
	version = "0.1.0"
)

// TODO [SNOW-2324252]: conversionErrorWrapped in templates?
func main() {
	genhelpers.NewGenerator(
		genhelpers.NewPreambleModel(name, version),
		poc.GetSdkDefinitions,
		poc.ExtendInterface("../../../dto-builder-generator/main.go"),
		filenameForPart(""),
		[]*template.Template{genhelpers.PreambleTemplate, generator.InterfaceTemplate, generator.OperationStructIterateTemplate},
	).
		WithGenerationPart("dto", filenameForPart("dto"), []*template.Template{genhelpers.PreambleTemplate, generator.DtoTemplate}).
		WithGenerationPart("dto_builders", filenameForPart("dto_builders"), []*template.Template{genhelpers.PreambleTemplate, generator.DtoBuildersTemplate}).
		WithGenerationPart("impl", filenameForPart("impl"), []*template.Template{genhelpers.PreambleTemplate, generator.ImplementationTemplate}).
		WithGenerationPart("unit_tests", testFilenameForPart(""), []*template.Template{genhelpers.PreambleTemplate, generator.UnitTestsTemplate}).
		WithGenerationPart("validations", filenameForPart("validations"), []*template.Template{genhelpers.PreambleTemplate, generator.ValidationsTemplate}).
		WithDescription("Generate SDK objects based on the SQL definitions provided.").
		WithMakefileCommandPart("sdk").
		RunAndHandleOsReturn()
}

func filenameForPart(part string) func(_ *generator.Interface, model *generator.Interface) string {
	return func(_ *generator.Interface, model *generator.Interface) string {
		var p string
		if part != "" {
			p = "_" + part
		}
		return genhelpers.ToSnakeCase(model.Name) + p + "_gen.go"
	}
}

func testFilenameForPart(part string) func(_ *generator.Interface, model *generator.Interface) string {
	return func(_ *generator.Interface, model *generator.Interface) string {
		var p string
		if part != "" {
			p = "_" + part
		}
		return genhelpers.ToSnakeCase(model.Name) + p + "_gen_test.go"
	}
}
