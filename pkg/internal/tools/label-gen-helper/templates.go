package main

import (
	_ "embed"
	"text/template"
)

//go:embed templates/repository_labels.tmpl
var RepositoryLabelsTemplateContent string
var RepositoryLabelsTemplate = template.Must(template.New("repository_labels").Parse(RepositoryLabelsTemplateContent))
