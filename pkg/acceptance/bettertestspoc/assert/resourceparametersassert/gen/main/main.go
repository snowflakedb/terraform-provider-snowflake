//go:build exclude

package main

import (
	objectparametersassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/resourceparametersassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/genhelpers"
)

func main() {
	genhelpers.NewGenerator(
		objectparametersassertgen.GetAllSnowflakeObjectParameters,
		gen.ModelFromSnowflakeObjectParameters,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ objectparametersassertgen.SnowflakeObjectParameters, model gen.ResourceParametersAssertionsModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_resource_parameters" + "_gen.go"
}
