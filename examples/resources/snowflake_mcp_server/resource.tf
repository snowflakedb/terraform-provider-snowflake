## Minimal using a heredoc string
resource "snowflake_mcp_server" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "mcp_server_name"

  specification = <<-EOT
tools:
  - title: "SQL Execution Tool"
    name: "sql_exec_tool"
    type: "SYSTEM_EXECUTE_SQL"
    description: "Executes SQL queries."
EOT
}

## Complete using yamlencode (with every optional set)
resource "snowflake_mcp_server" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "mcp_server_name"

  specification = yamlencode({
    tools = [
      {
        title       = "SQL Execution Tool"
        name        = "sql_exec_tool"
        type        = "SYSTEM_EXECUTE_SQL"
        description = "Executes SQL queries."
      }
    ]
  })

  comment = "My MCP server"
}
