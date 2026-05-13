# Simple usage
data "snowflake_session_policies" "simple" {
}

output "simple_output" {
  value = data.snowflake_session_policies.simple.session_policies
}

# Filtering (like)
data "snowflake_session_policies" "like" {
  like = "session-policy-name"
}

output "like_output" {
  value = data.snowflake_session_policies.like.session_policies
}

# Filtering by prefix (like)
data "snowflake_session_policies" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_session_policies.like_prefix.session_policies
}

# Filtering (starts_with)
data "snowflake_session_policies" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_session_policies.starts_with.session_policies
}

# Filtering (in)
data "snowflake_session_policies" "in_account" {
  in {
    account = true
  }
}

data "snowflake_session_policies" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_session_policies" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "account" : data.snowflake_session_policies.in_account.session_policies,
    "database" : data.snowflake_session_policies.in_database.session_policies,
    "schema" : data.snowflake_session_policies.in_schema.session_policies,
  }
}

# Filtering (on)
data "snowflake_session_policies" "on_account" {
  on {
    account = true
  }
}

data "snowflake_session_policies" "on_user" {
  on {
    user = "<user_name>"
  }
}

output "on_filtered" {
  value = {
    "account" : data.snowflake_session_policies.on_account.session_policies,
    "user" : data.snowflake_session_policies.on_user.session_policies,
  }
}

# Filtering (limit)
data "snowflake_session_policies" "limit" {
  limit {
    rows = 1
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_session_policies.limit.session_policies
}

# Without additional data (to limit the number of calls make for every found session policy)
data "snowflake_session_policies" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE SESSION POLICY for every session policy found and attaches its output to session_policies.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_session_policies.only_show.session_policies
}

# Ensure the number of session policies is equal to at least one element (with the use of postcondition)
data "snowflake_session_policies" "assert_with_postcondition" {
  like = "session-policy-name%"
  lifecycle {
    postcondition {
      condition     = length(self.session_policies) > 0
      error_message = "there should be at least one session policy"
    }
  }
}

# Ensure the number of session policies is equal to exactly one element (with the use of check block)
check "session_policy_check" {
  data "snowflake_session_policies" "assert_with_check_block" {
    like = "session-policy-name"
  }

  assert {
    condition     = length(data.snowflake_session_policies.assert_with_check_block.session_policies) == 1
    error_message = "session policies filtered by '${data.snowflake_session_policies.assert_with_check_block.like}' returned ${length(data.snowflake_session_policies.assert_with_check_block.session_policies)} session policies where one was expected"
  }
}
