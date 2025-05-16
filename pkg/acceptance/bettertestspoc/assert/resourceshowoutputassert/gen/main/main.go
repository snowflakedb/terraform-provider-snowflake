//go:build exclude

package main

import (
	objectassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/objectassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/genhelpers"
)

func main() {
	genhelpers.NewGenerator(
		objectassertgen.GetSdkObjectDetails,
		gen.ModelFromSdkObjectDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ genhelpers.SdkObjectDetails, model gen.ResourceShowOutputAssertionsModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_show_output" + "_gen.go"
}
