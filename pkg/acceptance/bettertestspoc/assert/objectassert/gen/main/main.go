//go:build exclude

package main

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const (
	name    = "object assertions"
	version = "0.1.0"
)

func main() {
	genhelpers.NewGenerator(
		genhelpers.NewPreambleModel(name, version).
			WithImport("fmt").
			WithImport("testing").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"),
		gen.GetSdkObjectDetails,
		gen.ModelFromSdkObjectDetails,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ genhelpers.SdkObjectDetails, model gen.SnowflakeObjectAssertionsModel) string {
	if model.IsDataSourceOutput {
		return strings.TrimSuffix(genhelpers.ToSnakeCase(model.Name), "_details") + "_desc_snowflake" + "_gen.go"
	}
	return genhelpers.ToSnakeCase(model.Name) + "_snowflake" + "_gen.go"
}
