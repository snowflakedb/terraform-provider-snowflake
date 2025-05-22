package main

import (
	resourcemodelgen "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/model/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/config/datasourcemodel/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/genhelpers"
)

func main() {
	genhelpers.NewGenerator(
		gen.GetDatasourceSchemaDetails,
		resourcemodelgen.ModelFromResourceSchemaDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ genhelpers.ResourceSchemaDetails, model resourcemodelgen.ResourceConfigBuilderModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_model" + "_gen.go"
}
