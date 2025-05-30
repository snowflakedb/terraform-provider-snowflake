package gen

import (
	"text/template"

	_ "embed"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

var (
	//go:embed templates/preamble.tmpl
	preambleTemplateContent string
	PreambleTemplate, _     = template.New("preambleTemplate").Parse(preambleTemplateContent)

	//go:embed templates/definition.tmpl
	definitionTemplateContent string
	DefinitionTemplate, _     = template.New("definitionTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
	)).Parse(definitionTemplateContent)

	//go:embed templates/assertions.tmpl
	assertionsTemplateContent string
	AssertionsTemplate, _     = template.New("assertionsTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.IsTypeSlice,
		genhelpers.SnakeCase,
		genhelpers.RunMapper,
	)).Parse(assertionsTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, DefinitionTemplate, AssertionsTemplate}
)
