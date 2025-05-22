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
	DefinitionTemplate, _     = template.New("definitionTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.SnakeCaseToCamel,
		genhelpers.RemoveForbiddenAttributeNameSuffix,
	)).Parse(definitionTemplateContent)

	//go:embed templates/marshal_json.tmpl
	marshalJsonTemplateContent string
	MarshalJsonTemplate, _     = template.New("marshalJsonTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.SnakeCaseToCamel,
	)).Parse(marshalJsonTemplateContent)

	//go:embed templates/builders.tmpl
	buildersTemplateContent string
	BuildersTemplate, _     = template.New("buildersTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.SnakeCaseToCamel,
		genhelpers.RemoveForbiddenAttributeNameSuffix,
		genhelpers.ShouldGenerateWithForAttributeType,
	)).Parse(buildersTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, DefinitionTemplate, MarshalJsonTemplate, BuildersTemplate}
)
