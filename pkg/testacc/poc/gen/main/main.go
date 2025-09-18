//go:build exclude

package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testacc/poc/gen"
)

const name = "PoC plugin framework model and schema"
const version = "0.1.0"

// TODO [this PR]: imports?
//  - "github.com/hashicorp/terraform-plugin-framework/provider/schema"
//  - "github.com/hashicorp/terraform-plugin-framework/types"
//  - "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"

func main() {
	genhelpers.NewGenerator(
		name,
		version,
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

func getFilename(_ gen.SdkV2ProviderSchema, _ gen.PluginFrameworkProviderModel) string {
	return "13_plugin_framework_model_and_schema_gen.go"
}
