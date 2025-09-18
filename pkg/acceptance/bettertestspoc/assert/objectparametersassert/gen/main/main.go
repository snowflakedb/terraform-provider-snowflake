//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const name = "object parameter assertions"
const version = "0.1.0"

func main() {
	genhelpers.NewGenerator(
		name,
		version,
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
