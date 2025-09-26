package gen

import (
	"text/template"

	_ "embed"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

var (
	//go:embed templates/definition.tmpl
	definitionTemplateContent string
	DefinitionTemplate, _     = template.New("definitionTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
	)).Parse(definitionTemplateContent)

	//go:embed templates/assertions.tmpl
	assertionsTemplateContent string
	AssertionsTemplate, _     = template.New("assertionsTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.TypeWithoutPointerAndBrackets,
		genhelpers.CamelToWords,
		genhelpers.RunMapper,
	)).Parse(assertionsTemplateContent)

	AllTemplates = []*template.Template{genhelpers.PreambleTemplate, DefinitionTemplate, AssertionsTemplate}
)
