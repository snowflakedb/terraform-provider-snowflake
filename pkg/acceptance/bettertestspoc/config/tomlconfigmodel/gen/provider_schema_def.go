package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type TomlConfigSchemaDef struct {
	name   string
	schema any
}

func GetTomlSchemaDetails() []genhelpers.TomlConfigSchemaDetails {
	allTomlConfigSchemas := allTomlSchemaDefs
	allTomlConfigSchemasDetails := make([]genhelpers.TomlConfigSchemaDetails, len(allTomlConfigSchemas))
	for idx, s := range allTomlConfigSchemas {
		allTomlConfigSchemasDetails[idx] = genhelpers.ExtractTomlConfigSchemaDetails(s.name, s.schema)
	}
	return allTomlConfigSchemasDetails
}

var allTomlSchemaDefs = []TomlConfigSchemaDef{
	{
		name:   "SnowflakeConfig",
		schema: &sdk.ConfigDTO{},
	},
}
