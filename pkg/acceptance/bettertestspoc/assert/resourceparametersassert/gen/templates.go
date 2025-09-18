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
		genhelpers.FirstLetter,
	)).Parse(definitionTemplateContent)

	//go:embed templates/specific_checks.tmpl
	specificChecksTemplateContent string
	SpecificChecksTemplate, _     = template.New("specificChecksTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.SnakeCaseToCamel,
	)).Parse(specificChecksTemplateContent)

	AllTemplates = []*template.Template{genhelpers.PreambleTemplate, DefinitionTemplate, SpecificChecksTemplate}
)
