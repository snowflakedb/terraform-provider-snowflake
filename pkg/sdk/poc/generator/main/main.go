// TODO[this PR]: go:build exclude

package main

import (
	"os"
	"slices"
	"strings"
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
		WithGenerationPart("impl", filenameForPart("impl"), []*template.Template{genhelpers.PreambleTemplate, generator.ImplementationTemplate}).
		WithGenerationPart("unit_tests", testFilenameForPart(""), []*template.Template{genhelpers.PreambleTemplate, generator.UnitTestsTemplate}).
		WithGenerationPart("validations", filenameForPart("validations"), []*template.Template{genhelpers.PreambleTemplate, generator.ValidationsTemplate}).
		WithGenerationPartFilter(filterGenerationPartByNameFromEnv[*generator.Interface, *generator.Interface]).
		RunAndHandleOsReturn()
}

func filenameForPart(part string) func(_ *generator.Interface, model *generator.Interface) string {
	return func(_ *generator.Interface, model *generator.Interface) string {
		if part != "" {
			part = "_" + part
		}
		return genhelpers.ToSnakeCase(model.Name) + part + "_gen.go"
	}
}

func testFilenameForPart(part string) func(_ *generator.Interface, model *generator.Interface) string {
	return func(_ *generator.Interface, model *generator.Interface) string {
		if part != "" {
			part = "_" + part
		}
		return genhelpers.ToSnakeCase(model.Name) + part + "_gen_test.go"
	}
}

// TODO [SNOW-2324252]: move this filter to commons and consider extracting this as a command line param
func filterGenerationPartByNameFromEnv[T genhelpers.ObjectNameProvider, M genhelpers.HasPreambleModel](part genhelpers.GenerationPart[T, M]) bool {
	allowedObjectNamesString := os.Getenv("SF_TF_GENERATOR_EXT_ALLOWED_GENERATION_PARTS_NAMES")
	if allowedObjectNamesString == "" {
		return true
	}
	allowedObjectNames := strings.Split(allowedObjectNamesString, ",")
	return slices.Contains(allowedObjectNames, part.GetName())
}
