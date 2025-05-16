//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/genhelpers"
)

func main() {
	genhelpers.NewGenerator(
		gen.GetAllSnowflakeObjectParameters,
		gen.ModelFromSnowflakeObjectParameters,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ gen.SnowflakeObjectParameters, model gen.SnowflakeObjectParametersAssertionsModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_parameters_snowflake" + "_gen.go"
}
