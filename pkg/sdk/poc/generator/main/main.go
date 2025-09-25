// TODO[this PR]: go:build exclude

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

// TODO [this PR]: conversionErrorWrapped in templates?
func main() {
	genhelpers.NewGenerator(
		genhelpers.NewPreambleModel(name, version),
		poc.GetSdkDefinitions,
		poc.WithPreamble,
		filenameFor(""),
		[]*template.Template{genhelpers.PreambleTemplate, generator.InterfaceTemplate, generator.OperationStructIterateTemplate},
	).
		WithGenerationPart(filenameFor("dto"), []*template.Template{genhelpers.PreambleTemplate, generator.DtoTemplate}).
		RunAndHandleOsReturn()
}

func filenameFor(part string) func(_ *generator.Interface, model *generator.Interface) string {
	return func(_ *generator.Interface, model *generator.Interface) string {
		if part != "" {
			part = "_" + part
		}
		return genhelpers.ToSnakeCase(model.Name) + part + "_gen.go"
	}
}
