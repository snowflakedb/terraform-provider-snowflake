//go:build exclude

package main

import (
	"text/template"

	_ "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/defs"
	_ "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/internal/example/defs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
)

const (
	name    = "SDK builder"
	version = "0.1.0"
)

func main() {
	genhelpers.NewGenerator(
		genhelpers.NewPreambleModel(name, version),
		gen.GetSdkDefinitions,
		gen.ExtendInterface(),
		filenameForPart(""),
		[]*template.Template{genhelpers.PreambleTemplate, gen.InterfaceTemplate, gen.OperationStructIterateTemplate},
	).
		WithGenerationPart("dto", filenameForPart("dto"), []*template.Template{genhelpers.PreambleTemplate, gen.DtoTemplate}).
		WithGenerationPart("dto_builders", filenameForPart("dto_builders"), []*template.Template{genhelpers.PreambleTemplate, gen.DtoBuildersTemplate}).
		WithGenerationPart("impl", filenameForPart("impl"), []*template.Template{genhelpers.PreambleTemplate, gen.ImplementationTemplate}).
		WithGenerationPart("unit_tests", testFilenameForPart(""), []*template.Template{genhelpers.PreambleTemplate, gen.UnitTestsTemplate}).
		WithGenerationPart("validations", filenameForPart("validations"), []*template.Template{genhelpers.PreambleTemplate, gen.ValidationsTemplate}).
		WithDescription("Generate SDK objects based on the SQL definitions provided.").
		WithMakefileCommandPart("sdk").
		RunAndHandleOsReturn()
}

func filenameForPart(part string) func(_ *gen.Interface, model *gen.Interface) string {
	return func(_ *gen.Interface, model *gen.Interface) string {
		var p string
		if part != "" {
			p = "_" + part
		}
		return genhelpers.ToSnakeCase(model.Name) + p + "_gen.go"
	}
}

func testFilenameForPart(part string) func(_ *gen.Interface, model *gen.Interface) string {
	return func(_ *gen.Interface, model *gen.Interface) string {
		var p string
		if part != "" {
			p = "_" + part
		}
		return genhelpers.ToSnakeCase(model.Name) + p + "_gen_test.go"
	}
}
