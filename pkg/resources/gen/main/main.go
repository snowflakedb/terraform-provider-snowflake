//go:build exclude

package main

import (
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources/gen"
)

const (
	name    = "resource"
	version = "0.1.0"
)

func main() {
	genhelpers.NewGenerator(
		genhelpers.NewPreambleModel(name, version).
			WithImport("context").
			WithImport("errors").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas").
			WithImport("github.com/hashicorp/terraform-plugin-sdk/v2/diag").
			WithImport("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"),
		gen.GetAllObjects,
		gen.ModelFromInputObject,
		lifecycleFilename,
		[]*template.Template{genhelpers.PreambleTemplate, gen.LifecycleTemplate},
	).
		WithGenerationPart(gen.PartSchema, schemaFilename, []*template.Template{genhelpers.PreambleTemplate, gen.SchemaTemplate}).
		WithDescription("Generate resource lifecycle and schema skeletons.").
		WithMakefileCommandPart("resource").
		RunAndHandleOsReturn()
}

func lifecycleFilename(_ gen.ResourceDef, model gen.ResourceModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_gen.go"
}

func schemaFilename(_ gen.ResourceDef, model gen.ResourceModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_schema_gen.go"
}
