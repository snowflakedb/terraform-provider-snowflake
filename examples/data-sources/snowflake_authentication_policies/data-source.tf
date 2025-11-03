# Simple usage
data "snowflake_authentication_policies" "simple" {
}

output "simple_output" {
  value = data.snowflake_authentication_policies.simple.authentication_policies
}

# Filtering (like)
data "snowflake_authentication_policies" "like" {
  like = "authentication-policy-name"
}

output "like_output" {
  value = data.snowflake_authentication_policies.like.authentication_policies
}

# Filtering by prefix (like)
data "snowflake_authentication_policies" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_authentication_policies.like_prefix.authentication_policies
}

# Filtering (starts_with)
data "snowflake_authentication_policies" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_authentication_policies.starts_with.authentication_policies
}

# Filtering (in)
data "snowflake_authentication_policies" "in_account" {
  in {
    account = true
  }
}

data "snowflake_authentication_policies" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_authentication_policies" "in_schema" {
  in {
    schema = "<database_name>.<schema_name>"
  }
}

output "in_output" {
  value = {
    "account" : data.snowflake_authentication_policies.in_account.authentication_policies,
    "database" : data.snowflake_authentication_policies.in_database.authentication_policies,
    "schema" : data.snowflake_authentication_policies.in_schema.authentication_policies,
  }
}

# Filtering (on)
data "snowflake_authentication_policies" "on_account" {
  on {
    account = true
  }
}

data "snowflake_authentication_policies" "on_user" {
  on {
    user = "<user_name>"
  }
}

output "on_output" {
  value = {
    "account" : data.snowflake_authentication_policies.on_account.authentication_policies,
    "user" : data.snowflake_authentication_policies.on_user.authentication_policies,
  }
}


# Filtering (limit)
data "snowflake_authentication_policies" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_authentication_policies.limit.authentication_policies
}

# Without additional data (to limit the number of calls make for every found authentication policy)
data "snowflake_authentication_policies" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE AUTHENTICATION POLICY for every authentication policy found and attaches its output to authentication_policies.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_authentication_policies.only_show.authentication_policies
}

# Ensure the number of authentication policies is equal to at least one element (with the use of postcondition)
data "snowflake_authentication_policies" "assert_with_postcondition" {
  like = "authentication-policy-name%"
  lifecycle {
    postcondition {
      condition     = length(self.authentication_policies) > 0
      error_message = "there should be at least one authentication policy"
    }
  }
}

# Ensure the number of authentication policies is equal to exactly one element (with the use of check block)
check "authentication_policy_check" {
  data "snowflake_authentication_policies" "assert_with_check_block" {
    like = "authentication-policy-name"
  }

  assert {
    condition     = length(data.snowflake_authentication_policies.assert_with_check_block.authentication_policies) == 1
    error_message = "authentication policies filtered by '${data.snowflake_authentication_policies.assert_with_check_block.like}' returned ${length(data.snowflake_authentication_policies.assert_with_check_block.authentication_policies)} authentication policies where one was expected"
  }
}
