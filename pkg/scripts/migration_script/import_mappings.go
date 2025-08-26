package main

import (
	"fmt"
	"strings"
	"text/template"

	_ "embed"
)

var (
	//go:embed templates/import_block.tf.tmpl
	ImportBlockTemplateContent string
	ImportBlockTemplate, _     = template.New("import_block_template").Parse(ImportBlockTemplateContent)

	//go:embed templates/import_statement.tf.tmpl
	ImportStatementTemplateContent string
	ImportStatementTemplate, _     = template.New("import_statement_template").Parse(ImportStatementTemplateContent)
)

type ImportModel struct {
	ResourceAddress string
	Id              string
}

// IdEscaped returns the ID with escaped quotes for use in Terraform import blocks.
func (im ImportModel) IdEscaped() string {
	return strings.ReplaceAll(im.Id, "\"", "\\\"")
}

func TransformImportModel(config *Config, m ImportModel) (string, error) {
	stringBuilder := new(strings.Builder)

	switch config.ImportFlag {
	case ImportStatementTypeBlock:
		if err := ImportBlockTemplate.Execute(stringBuilder, m); err != nil {
			return "", err
		}
		stringBuilder.WriteString("\n")
	case ImportStatementTypeStatement:
		if err := ImportStatementTemplate.Execute(stringBuilder, m); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported import statement type: %s", config.ImportFlag)
	}

	return stringBuilder.String(), nil
}
