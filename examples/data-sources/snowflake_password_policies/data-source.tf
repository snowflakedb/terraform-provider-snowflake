# Simple usage
data "snowflake_password_policies" "simple" {
}

output "simple_output" {
  value = data.snowflake_password_policies.simple.password_policies
}

# Filtering (like)
data "snowflake_password_policies" "like" {
  like = "password-policy-name"
}

output "like_output" {
  value = data.snowflake_password_policies.like.password_policies
}

# Filtering by prefix (like)
data "snowflake_password_policies" "like_prefix" {
  like = "prefix%"
}

output "like_prefix_output" {
  value = data.snowflake_password_policies.like_prefix.password_policies
}

# Filtering (starts_with)
data "snowflake_password_policies" "starts_with" {
  starts_with = "prefix-"
}

output "starts_with_output" {
  value = data.snowflake_password_policies.starts_with.password_policies
}

# Filtering (in)
data "snowflake_password_policies" "in_account" {
  in {
    account = true
  }
}

data "snowflake_password_policies" "in_database" {
  in {
    database = "<database_name>"
  }
}

data "snowflake_password_policies" "in_schema" {
  in {
    schema = "\"<database_name>\".\"<schema_name>\""
  }
}

output "in_filtered" {
  value = {
    "account" : data.snowflake_password_policies.in_account.password_policies,
    "database" : data.snowflake_password_policies.in_database.password_policies,
    "schema" : data.snowflake_password_policies.in_schema.password_policies,
  }
}

# Filtering (on)
data "snowflake_password_policies" "on_account" {
  on {
    account = true
  }
}

data "snowflake_password_policies" "on_user" {
  on {
    user = "<user_name>"
  }
}

output "on_filtered" {
  value = {
    "account" : data.snowflake_password_policies.on_account.password_policies,
    "user" : data.snowflake_password_policies.on_user.password_policies,
  }
}

# Filtering (limit)
data "snowflake_password_policies" "limit" {
  limit {
    rows = 10
    from = "prefix-"
  }
}

output "limit_output" {
  value = data.snowflake_password_policies.limit.password_policies
}

# Without additional data (to limit the number of calls make for every found password policy)
data "snowflake_password_policies" "only_show" {
  # with_describe is turned on by default and it calls DESCRIBE PASSWORD POLICY for every password policy found and attaches its output to password_policies.*.describe_output field
  with_describe = false
}

output "only_show_output" {
  value = data.snowflake_password_policies.only_show.password_policies
}

# Ensure the number of password policies is equal to at least one element (with the use of postcondition)
data "snowflake_password_policies" "assert_with_postcondition" {
  like = "password-policy-name%"
  lifecycle {
    postcondition {
      condition     = length(self.password_policies) > 0
      error_message = "there should be at least one password policy"
    }
  }
}

# Ensure the number of password policies is equal to exactly one element (with the use of check block)
check "password_policy_check" {
  data "snowflake_password_policies" "assert_with_check_block" {
    like = "password-policy-name"
  }

  assert {
    condition     = length(data.snowflake_password_policies.assert_with_check_block.password_policies) == 1
    error_message = "password policies filtered by '${data.snowflake_password_policies.assert_with_check_block.like}' returned ${length(data.snowflake_password_policies.assert_with_check_block.password_policies)} password policies where one was expected"
  }
}
