//go:build exclude

package main

import (
	objectparametersassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const name = "resource parameter assertions"
const version = "0.1.0"

// TODO [this PR]: imports?
//  - "testing"
//  - "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
//  - "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func main() {
	genhelpers.NewGenerator(
		name,
		version,
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
