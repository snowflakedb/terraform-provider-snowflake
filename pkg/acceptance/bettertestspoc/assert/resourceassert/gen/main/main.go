//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const (
	name    = "resource assertions"
	version = "0.1.0"
)

func main() {
	genhelpers.NewGenerator(
		genhelpers.NewPreambleModel(name, version).
			WithImport("testing").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert").
			WithImport("github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"),
		gen.ResourceAssertionModels,
		gen.WithPreamble,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getFilename(_ gen.ResourceAssertionsModel, model gen.ResourceAssertionsModel) string {
	return genhelpers.ToSnakeCase(model.Name) + "_resource" + "_gen.go"
}
