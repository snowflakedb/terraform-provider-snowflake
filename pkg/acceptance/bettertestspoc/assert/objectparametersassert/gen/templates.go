package gen

import (
	"text/template"

	_ "embed"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/genhelpers"
)

var (
	//go:embed templates/preamble.tmpl
	preambleTemplateContent string
	PreambleTemplate, _     = template.New("preambleTemplate").Parse(preambleTemplateContent)

	//go:embed templates/definition.tmpl
	definitionTemplateContent string
	DefinitionTemplate, _     = template.New("definitionTemplate").Parse(definitionTemplateContent)

	//go:embed templates/generic_checks.tmpl
	genericChecksTemplateContent string
	GenericChecksTemplate, _     = template.New("genericChecksTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
	)).Parse(genericChecksTemplateContent)

	//go:embed templates/aggregated_generic_checks.tmpl
	aggregatedGenericChecksTemplateContent string
	AggregatedGenericChecksTemplate, _     = template.New("aggregatedGenericChecksTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.SnakeCaseToCamel,
		genhelpers.IsLastItem,
	)).Parse(aggregatedGenericChecksTemplateContent)

	//go:embed templates/specific_checks.tmpl
	specificChecksTemplateContent string
	SpecificChecksTemplate, _     = template.New("specificChecksTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.SnakeCaseToCamel,
	)).Parse(specificChecksTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, DefinitionTemplate, GenericChecksTemplate, AggregatedGenericChecksTemplate, SpecificChecksTemplate}
)
