//go:build exclude

package main

import (
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const (
	name    = "datasource"
	version = "0.1.0"
)

func main() {
	genhelpers.NewGenerator(
		genhelpers.NewPreambleModel(name, version).
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"),
		gen.GetAllObjects,
		gen.ModelFromInputObject,
		dataSourceFilename,
		[]*template.Template{genhelpers.PreambleTemplate, gen.DataSourceTemplate},
	).
		WithDescription("Generate data sources.").
		WithMakefileCommandPart("datasource").
		RunAndHandleOsReturn()
}

func dataSourceFilename(_ gen.DataSourceObject, model gen.DataSourceModel) string {
	return genhelpers.ToSnakeCase(model.NamePlural) + "_gen.go"
}
