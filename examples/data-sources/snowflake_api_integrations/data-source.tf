# Simple usage
data "snowflake_api_integrations" "simple" {
}

output "simple_output" {
  value = data.snowflake_api_integrations.simple.api_integrations
}

# Filtering (like)
data "snowflake_api_integrations" "like" {
  like = "api-integration-name"
}

output "like_output" {
  value = data.snowflake_api_integrations.like.api_integrations
}

# Filtering by prefix (like)
data "snowflake_api_integrations" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_api_integrations.like_prefix.api_integrations
}

# Without additional data (to limit the number of calls made for every found API integration)
data "snowflake_api_integrations" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE API INTEGRATION for every integration found and attaches its output to api_integrations.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_api_integrations.only_show.api_integrations
}
