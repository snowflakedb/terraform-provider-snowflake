//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const name = "resource assertions"
const version = "0.1.0"

// TODO [this PR]: imports?
//  - "testing"
//  - "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func main() {
	genhelpers.NewGenerator(
		name,
		version,
		gen.GetResourceSchemaDetails,
		gen.ModelFromResourceSchemaDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ genhelpers.ResourceSchemaDetails, model gen.ResourceAssertionsModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_resource" + "_gen.go"
}
