//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/objectassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/genhelpers"
)

func main() {
	genhelpers.NewGenerator(
		gen.GetSdkObjectDetails,
		gen.ModelFromSdkObjectDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ genhelpers.SdkObjectDetails, model gen.SnowflakeObjectAssertionsModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_snowflake" + "_gen.go"
}
