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

	AllTemplates = []*template.Template{genhelpers.PreambleTemplate, DefinitionTemplate, MarshalJsonTemplate, BuildersTemplate}
)
