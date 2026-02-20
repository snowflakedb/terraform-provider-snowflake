package gen

import (
	"text/template"

	_ "embed"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

var (
	//go:embed templates/definition.tmpl
	definitionTemplateContent string
	DefinitionTemplate, _     = template.New("definitionTemplate").Parse(definitionTemplateContent)

	// TODO [SNOW-3113128]: use .IsCollection logic for string checks
	//go:embed templates/assertions.tmpl
	assertionsTemplateContent string
	AssertionsTemplate, _     = template.New("assertionsTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.SnakeCaseToCamel,
	)).Parse(assertionsTemplateContent)

	AllTemplates = []*template.Template{genhelpers.PreambleTemplate, DefinitionTemplate, AssertionsTemplate}
)
