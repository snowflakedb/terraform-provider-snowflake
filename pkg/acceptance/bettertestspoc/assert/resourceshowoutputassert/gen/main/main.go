//go:build exclude

package main

import (
	objectassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const (
	name    = "resource show output assertions"
	version = "0.1.0"
)

func main() {
	genhelpers.NewGenerator(
		genhelpers.NewPreambleModel(name, version).
			WithImport("testing").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"),
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
