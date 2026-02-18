package gen

import (
	"strings"
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
		strings.TrimSuffix,
	)).Parse(definitionTemplateContent)

	//go:embed templates/assertions.tmpl
	assertionsTemplateContent string
	AssertionsTemplate, _     = template.New("assertionsTemplate").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetter,
		genhelpers.IsTypeSlice,
		genhelpers.SnakeCase,
		genhelpers.RunMapper,
		strings.TrimSuffix,
	)).Parse(assertionsTemplateContent)

	AllTemplates = []*template.Template{genhelpers.PreambleTemplate, DefinitionTemplate, AssertionsTemplate}
)
