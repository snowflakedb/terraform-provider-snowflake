package model

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
)

func McpServerWithSpecification(resourceName, database, schema, name, specification string) *McpServerModel {
	model := McpServer(resourceName, database, schema, name, "")
	// This prevents double quotes from being added around yamlencode.
	model.WithSpecificationValue(config.UnquotedWrapperVariable(specification))
	return model
}

// DefaultSpecAsYamlencodeHCL returns a HCL yamlencode expression for a default MCP server spec
// suitable for use with McpServerWithSpecification.
// String values use SnowflakeProviderConfigQuoteMarker so they survive JSON encoding
// (json.Marshal would otherwise escape the double quotes to \").
func DefaultSpecAsYamlencodeHCL() string {
	return fmt.Sprintf(`yamlencode({
  tools = [
    {
      title       = %[1]sSQL Execution Tool%[1]s
      name        = %[1]ssql_exec_tool%[1]s
      type        = %[1]sSYSTEM_EXECUTE_SQL%[1]s
      description = %[1]sFor acceptance tests.%[1]s
    }
  ]
})`, string(config.SnowflakeProviderConfigQuoteMarker))
}

// AltSpecAsYamlencodeHCL returns an alternate HCL yamlencode expression for update tests.
func AltSpecAsYamlencodeHCL() string {
	return fmt.Sprintf(`yamlencode({
  tools = [
    {
      title       = %[1]sSQL Execution Tool%[1]s
      name        = %[1]ssql_exec_tool%[1]s
      type        = %[1]sSYSTEM_EXECUTE_SQL%[1]s
      description = %[1]sUpdated description for acceptance tests.%[1]s
    }
  ]
})`, string(config.SnowflakeProviderConfigQuoteMarker))
}
