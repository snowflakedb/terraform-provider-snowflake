//go:build exclude

package main

import (
	resourceassertgen "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const name = "resource model builder"
const version = "0.1.0"

// TODO [this PR]: import "encoding/json"? additional imports?
// TODO [this PR]: imports?
//  - tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
//  - "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
//  - "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

func main() {
	genhelpers.NewGenerator(
		name,
		version,
		resourceassertgen.GetResourceSchemaDetails,
		gen.ModelFromResourceSchemaDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ genhelpers.ResourceSchemaDetails, model gen.ResourceConfigBuilderModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_model" + "_gen.go"
}
