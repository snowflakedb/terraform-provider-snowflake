//go:build exclude

package main

import (
	resourcemodelgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const (
	name    = "provider model builder"
	version = "0.1.0"
)

func main() {
	genhelpers.NewGenerator(
		genhelpers.NewPreambleModel(name, version).
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config").
			WithNamedImport("tfconfig", "github.com/hashicorp/terraform-plugin-testing/config"),
		gen.GetProviderSchemaDetails,
		resourcemodelgen.ModelFromResourceSchemaDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ genhelpers.ResourceSchemaDetails, model resourcemodelgen.ResourceConfigBuilderModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_model" + "_gen.go"
}
