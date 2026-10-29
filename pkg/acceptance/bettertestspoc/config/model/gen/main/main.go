//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const (
	name    = "resource model builder"
	version = "0.1.0"
)

func main() {
	genhelpers.NewGenerator(
		genhelpers.NewPreambleModel(name, version).
			WithImport("encoding/json").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk").
			WithNamedImport("tfconfig", "github.com/hashicorp/terraform-plugin-testing/config"),
		gen.ModelBuilderModels,
		gen.WithPreamble,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ gen.ResourceConfigBuilderModel, model gen.ResourceConfigBuilderModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_model" + "_gen.go"
}
