package main

import (
	resourcemodelgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/tomlconfigmodel/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

func main() {
	genhelpers.NewGenerator(
		gen.GetTomlSchemaDetails,
		resourcemodelgen.ModelFromTomlConfigSchemaDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ genhelpers.TomlConfigSchemaDetails, model resourcemodelgen.TomlConfigBuilderModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_model" + "_gen.go"
}
