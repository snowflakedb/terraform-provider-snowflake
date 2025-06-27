//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testacc/poc/gen"
)

func main() {
	genhelpers.NewGenerator(
		getSdkV2ProviderSchemas,
		gen.ModelFromSdkV2Schema,
		getFilename,
		gen.AllTemplates,
	).
		RunAndHandleOsReturn()
}

func getSdkV2ProviderSchemas() []gen.SdkV2ProviderSchema {
	return gen.SdkV2ProviderSchemas
}

// TODO[mux-PR]: Add version?
func getFilename(_ gen.SdkV2ProviderSchema, _ gen.PluginFrameworkProviderModel) string {
	return "13_plugin_framework_model_and_schema_gen.go"
}
