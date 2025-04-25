package main

import (
	_ "embed"
	"text/template"
)

//go:embed templates/deprecated_resources.tmpl
var DeprecatedResourcesTemplateContent string
var DeprecatedResourcesTemplate = template.Must(template.New("deprecated_resources").Parse(DeprecatedResourcesTemplateContent))

//go:embed templates/deprecated_datasources.tmpl
var DeprecatedDatasourcesTemplateContent string
var DeprecatedDatasourcesTemplate = template.Must(template.New("deprecated_data_sources").Parse(DeprecatedDatasourcesTemplateContent))

//go:embed templates/feature_stability.tmpl
var FeatureStabilityTemplateContent string
var FeatureStabilityTemplate = template.Must(template.New("stable_resources").Parse(FeatureStabilityTemplateContent))
