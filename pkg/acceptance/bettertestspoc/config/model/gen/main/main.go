package main

import (
	resourceassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/resourceassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/genhelpers"
)

func main() {
	genhelpers.NewGenerator(
		resourceassertgen.GetResourceSchemaDetails,
		gen.ModelFromResourceSchemaDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ genhelpers.ResourceSchemaDetails, model gen.ResourceConfigBuilderModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_model" + "_gen.go"
}
