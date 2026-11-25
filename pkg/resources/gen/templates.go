package gen

import (
	"embed"
	"strconv"
	"strings"
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

//go:embed templates/sub_templates/*.tmpl
var subTemplatesFS embed.FS

var (
	//go:embed templates/lifecycle.tmpl
	lifecycleTemplateContent string
	LifecycleTemplate        *template.Template

	//go:embed templates/schema.tmpl
	schemaTemplateContent string
	SchemaTemplate        *template.Template

	//go:embed templates/parameters.tmpl
	parametersTemplateContent string
	ParametersTemplate        *template.Template
)

func init() {
	subTemplates := template.New("subTemplates").Funcs(genhelpers.BuildTemplateFuncMap(
		genhelpers.Unimplemented,
		genhelpers.CamelToWords,
		genhelpers.FirstLetterLowercase,
		genhelpers.FirstLetterUppercase,
		genhelpers.SnakeCaseToCamel,
		strconv.Quote,
		strings.ToLower,
		strings.ToUpper,
	))
	genhelpers.ParseSubTemplates(subTemplates, subTemplatesFS, "templates/sub_templates")
	LifecycleTemplate = template.Must(subTemplates.New("lifecycleTemplate").Parse(lifecycleTemplateContent))
	SchemaTemplate = template.Must(subTemplates.New("schemaTemplate").Parse(schemaTemplateContent))
	ParametersTemplate = template.Must(subTemplates.New("parametersTemplate").Parse(parametersTemplateContent))
}
