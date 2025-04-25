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
		genhelpers.SnakeCaseToCamel,
	)).Parse(definitionTemplateContent)

	//go:embed templates/builders.tmpl
	buildersTemplateContent string
	BuildersTemplate, _     = template.New("buildersTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.SnakeCaseToCamel,
		genhelpers.RemoveForbiddenAttributeNameSuffix,
	)).Parse(buildersTemplateContent)

	// TODO [SNOW-1501905]: consider duplicating the builders template from resource (currently same template used for datasources and provider which limits the customization possibilities for just one block type)
	AllTemplates = []*template.Template{PreambleTemplate, DefinitionTemplate, BuildersTemplate}
)
