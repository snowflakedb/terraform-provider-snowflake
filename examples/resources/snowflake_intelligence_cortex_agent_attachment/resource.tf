resource "snowflake_cortex_agent" "agent" {
  database = "my_database"
  schema   = "my_schema"
  name     = "my_agent"

  specification = <<-EOT
orchestration:
  budget:
    seconds: 30
    tokens: 16000
instructions:
  response: "You are a helpful assistant."
EOT
}

resource "snowflake_intelligence_cortex_agent_attachment" "agent" {
  snowflake_intelligence_name = "SNOWFLAKE_INTELLIGENCE_OBJECT_DEFAULT"
  cortex_agent_name           = snowflake_cortex_agent.agent.fully_qualified_name
}
