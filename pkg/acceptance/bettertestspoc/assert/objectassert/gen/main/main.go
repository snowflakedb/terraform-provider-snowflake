//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const name = "object assertions"
const version = "0.1.0"

// TODO [this PR]: imports?
//  - "fmt"
//  - "testing"
//  - "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
//  - "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
//  - "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func main() {
	genhelpers.NewGenerator(
		name,
		version,
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
