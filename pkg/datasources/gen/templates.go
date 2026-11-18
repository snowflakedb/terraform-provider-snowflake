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
	//go:embed templates/datasource.tmpl
	dataSourceTemplateContent string
	DataSourceTemplate        *template.Template
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
	DataSourceTemplate = template.Must(subTemplates.New("dataSourceTemplate").Parse(dataSourceTemplateContent))
}
