## Minimal using a heredoc string
resource "snowflake_cortex_agent" "basic" {
  database = "database_name"
  schema   = "schema_name"
  name     = "cortex_agent_name"

  specification = <<-EOT
orchestration:
  budget:
    seconds: 30
    tokens: 16000
instructions:
  response: "You are a helpful assistant."
EOT
}

## Complete using yamlencode (with every optional set)
resource "snowflake_cortex_agent" "complete" {
  database = "database_name"
  schema   = "schema_name"
  name     = "cortex_agent_name"

  specification = yamlencode({
    orchestration = {
      budget = {
        seconds = 30
        tokens  = 16000
      }
    }
    instructions = {
      response = "You are a helpful assistant."
    }
  })

  comment = "My Cortex agent"
  profile {
    display_name = "My Helpful Assistant"
    avatar       = "business-icon.png"
    color        = "red"
  }
}
