package gen

import (
	"text/template"

	_ "embed"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

var (
	//go:embed templates/model.tmpl
	modelTemplateContent string
	ModelTemplate, _     = template.New("modelTemplate").Parse(modelTemplateContent)

	//go:embed templates/schema.tmpl
	schemaTemplateContent string
	SchemaTemplate, _     = template.New("schemaTemplate").Parse(schemaTemplateContent)

	AllTemplates = []*template.Template{genhelpers.PreambleTemplate, ModelTemplate, SchemaTemplate}
)
