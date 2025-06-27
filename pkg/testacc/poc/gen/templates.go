package gen

import (
	"text/template"

	_ "embed"
)

var (
	//go:embed templates/preamble.tmpl
	preambleTemplateContent string
	PreambleTemplate, _     = template.New("preambleTemplate").Parse(preambleTemplateContent)

	//go:embed templates/model.tmpl
	modelTemplateContent string
	ModelTemplate, _     = template.New("modelTemplate").Parse(modelTemplateContent)

	//go:embed templates/schema.tmpl
	schemaTemplateContent string
	SchemaTemplate, _     = template.New("schemaTemplate").Parse(schemaTemplateContent)

	AllTemplates = []*template.Template{PreambleTemplate, ModelTemplate, SchemaTemplate}
)
