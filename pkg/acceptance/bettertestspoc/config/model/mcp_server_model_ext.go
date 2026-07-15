package model

import (
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
func DefaultSpecAsYamlencodeHCL() string {
	return `yamlencode({
  tools = [
    {
      title       = "SQL Execution Tool"
      name        = "sql_exec_tool"
      type        = "SYSTEM_EXECUTE_SQL"
      description = "For acceptance tests."
    }
  ]
})`
}

// AltSpecAsYamlencodeHCL returns an alternate HCL yamlencode expression for update tests.
func AltSpecAsYamlencodeHCL() string {
	return `yamlencode({
  tools = [
    {
      title       = "SQL Execution Tool"
      name        = "sql_exec_tool"
      type        = "SYSTEM_EXECUTE_SQL"
      description = "Updated description for acceptance tests."
    }
  ]
})`
}
