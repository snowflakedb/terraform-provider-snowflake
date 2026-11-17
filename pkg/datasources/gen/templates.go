package gen

import (
	"strconv"
	"strings"
	"text/template"

	_ "embed"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

var (
	//go:embed templates/datasource.tmpl
	dataSourceTemplateContent string
	DataSourceTemplate        *template.Template

	//go:embed templates/sub_templates/switches_schema.tmpl
	switchesSchemaTemplateContent string

	//go:embed templates/sub_templates/filters_schema.tmpl
	filtersSchemaTemplateContent string

	//go:embed templates/sub_templates/item_schema.tmpl
	itemSchemaTemplateContent string

	//go:embed templates/sub_templates/show_output_schema.tmpl
	showOutputSchemaTemplateContent string

	//go:embed templates/sub_templates/describe_output_schema.tmpl
	describeOutputSchemaTemplateContent string

	//go:embed templates/sub_templates/parameters_output_schema.tmpl
	parametersOutputSchemaTemplateContent string

	//go:embed templates/sub_templates/additional_schema.tmpl
	additionalSchemaTemplateContent string

	//go:embed templates/sub_templates/read_wrapper.tmpl
	readWrapperTemplateContent string

	//go:embed templates/sub_templates/description.tmpl
	descriptionTemplateContent string

	//go:embed templates/sub_templates/apply_filters.tmpl
	applyFiltersTemplateContent string

	//go:embed templates/sub_templates/show_in_sdk.tmpl
	showInSdkTemplateContent string

	//go:embed templates/sub_templates/set_id.tmpl
	setIdTemplateContent string

	//go:embed templates/sub_templates/describe_in_sdk.tmpl
	describeInSdkTemplateContent string

	//go:embed templates/sub_templates/describe_output_set.tmpl
	describeOutputSetTemplateContent string

	//go:embed templates/sub_templates/show_parameters_in_sdk.tmpl
	showParametersInSdkTemplateContent string

	//go:embed templates/sub_templates/parameters_output_set.tmpl
	parametersOutputSetTemplateContent string

	//go:embed templates/sub_templates/item_extra.tmpl
	itemExtraTemplateContent string
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
	for name, content := range map[string]string{
		"switchesSchema":         switchesSchemaTemplateContent,
		"filtersSchema":          filtersSchemaTemplateContent,
		"itemSchema":             itemSchemaTemplateContent,
		"showOutputSchema":       showOutputSchemaTemplateContent,
		"describeOutputSchema":   describeOutputSchemaTemplateContent,
		"parametersOutputSchema": parametersOutputSchemaTemplateContent,
		"additionalSchema":       additionalSchemaTemplateContent,
		"readWrapper":            readWrapperTemplateContent,
		"description":            descriptionTemplateContent,
		"applyFilters":           applyFiltersTemplateContent,
		"showInSdk":              showInSdkTemplateContent,
		"setId":                  setIdTemplateContent,
		"describeInSdk":          describeInSdkTemplateContent,
		"describeOutputSet":      describeOutputSetTemplateContent,
		"showParametersInSdk":    showParametersInSdkTemplateContent,
		"parametersOutputSet":    parametersOutputSetTemplateContent,
		"itemExtra":              itemExtraTemplateContent,
	} {
		template.Must(subTemplates.New(name).Parse(content))
	}
	DataSourceTemplate = template.Must(subTemplates.New("dataSourceTemplate").Parse(dataSourceTemplateContent))
}
