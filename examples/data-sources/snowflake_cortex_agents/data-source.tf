# Simple usage
data "snowflake_cortex_agents" "simple" {
}

output "simple_output" {
  value = data.snowflake_cortex_agents.simple.cortex_agents
}

# Filtering (like)
data "snowflake_cortex_agents" "like" {
  like = "cortex-agent-name"
}

output "like_output" {
  value = data.snowflake_cortex_agents.like.cortex_agents
}

# Filtering by prefix (like)
data "snowflake_cortex_agents" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_cortex_agents.like_prefix.cortex_agents
}

# Filtering (starts_with)
data "snowflake_cortex_agents" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_cortex_agents.starts_with.cortex_agents
}

# Filtering (in)
data "snowflake_cortex_agents" "in_account" {
  in {
    account = true
  }
}

data "snowflake_cortex_agents" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_cortex_agents" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "account" : data.snowflake_cortex_agents.in_account.cortex_agents,
    "database" : data.snowflake_cortex_agents.in_database.cortex_agents,
    "schema" : data.snowflake_cortex_agents.in_schema.cortex_agents,
  }
}

# Filtering (limit)
data "snowflake_cortex_agents" "limit" {
  limit {
    rows = 1
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_cortex_agents.limit.cortex_agents
}

# Without additional data (to limit the number of calls made for every found Cortex agent)
data "snowflake_cortex_agents" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE AGENT for every Cortex agent found and attaches its output to cortex_agents.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_cortex_agents.only_show.cortex_agents
}

# Ensure the number of Cortex agents is equal to at least one element (with the use of postcondition)
data "snowflake_cortex_agents" "assert_with_postcondition" {
  like = "cortex-agent-name%"
  lifecycle {
    postcondition {
      condition     = length(self.cortex_agents) > 0
      error_message = "there should be at least one cortex agent"
    }
  }
}

# Ensure the number of Cortex agents is equal to exactly one element (with the use of check block)
check "cortex_agent_check" {
  data "snowflake_cortex_agents" "assert_with_check_block" {
    like = "cortex-agent-name"
  }

  assert {
    condition     = length(data.snowflake_cortex_agents.assert_with_check_block.cortex_agents) == 1
    error_message = "cortex agents filtered by '${data.snowflake_cortex_agents.assert_with_check_block.like}' returned ${length(data.snowflake_cortex_agents.assert_with_check_block.cortex_agents)} cortex agents where one was expected"
  }
}
