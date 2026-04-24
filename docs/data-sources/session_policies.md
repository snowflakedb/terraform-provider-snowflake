---
page_title: "snowflake_session_policies Data Source - terraform-provider-snowflake"
subcategory: "Preview"
description: |-
  Data source used to get details of filtered session policies. Filtering is aligned with the current possibilities for SHOW SESSION POLICIES https://docs.snowflake.com/en/sql-reference/sql/show-session-policies query. The results of SHOW and DESCRIBE are encapsulated in one output collection session_policies.
---

!> **Caution: Preview Feature** This feature is considered a preview feature in the provider, regardless of the state of the resource in Snowflake. We do not guarantee its stability. It will be reworked and marked as a stable feature in future releases. Breaking changes are expected, even without bumping the major version. To use this feature, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

# snowflake_session_policies (Data Source)

Data source used to get details of filtered session policies. Filtering is aligned with the current possibilities for [SHOW SESSION POLICIES](https://docs.snowflake.com/en/sql-reference/sql/show-session-policies) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `session_policies`.

## Example Usage

```terraform
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
```

-> **Note** If a field has a default value, it is shown next to the type in the schema.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `in` (Block List, Max: 1) IN clause to filter the list of objects (see [below for nested schema](#nestedblock--in))
- `like` (String) Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).
- `limit` (Block List, Max: 1) Limits the number of rows returned. If the `limit.from` is set, then the limit will start from the first element matched by the expression. The expression is only used to match with the first element, later on the elements are not matched by the prefix, but you can enforce a certain pattern with `starts_with` or `like`. (see [below for nested schema](#nestedblock--limit))
- `on` (Block List, Max: 1) Lists the policies that are effective on the specified object. (see [below for nested schema](#nestedblock--on))
- `starts_with` (String) Filters the output with **case-sensitive** characters indicating the beginning of the object name.
- `with_describe` (Boolean) (Default: `true`) Runs DESC SESSION POLICY for each object returned by SHOW SESSION POLICIES. The output of describe is saved to the describe_output field. By default this value is set to true.

### Read-Only

- `id` (String) The ID of this resource.
- `session_policies` (List of Object) Holds the aggregated output of all session policy details queries. (see [below for nested schema](#nestedatt--session_policies))

<a id="nestedblock--in"></a>
### Nested Schema for `in`

Optional:

- `account` (Boolean) Returns records for the entire account.
- `application` (String) Returns records for the specified application.
- `application_package` (String) Returns records for the specified application package.
- `database` (String) Returns records for the current database in use or for a specified database.
- `schema` (String) Returns records for the current schema in use or a specified schema. Use fully qualified name.


<a id="nestedblock--limit"></a>
### Nested Schema for `limit`

Required:

- `rows` (Number) The maximum number of rows to return.

Optional:

- `from` (String) Specifies a **case-sensitive** pattern that is used to match object name. After the first match, the limit on the number of rows will be applied.


<a id="nestedblock--on"></a>
### Nested Schema for `on`

Optional:

- `account` (Boolean) Returns records for the entire account.
- `user` (String) Returns records for the specified user.


<a id="nestedatt--session_policies"></a>
### Nested Schema for `session_policies`

Read-Only:

- `describe_output` (List of Object) (see [below for nested schema](#nestedobjatt--session_policies--describe_output))
- `show_output` (List of Object) (see [below for nested schema](#nestedobjatt--session_policies--show_output))

<a id="nestedobjatt--session_policies--describe_output"></a>
### Nested Schema for `session_policies.describe_output`

Read-Only:

- `allowed_secondary_roles` (List of String)
- `blocked_secondary_roles` (List of String)
- `comment` (String)
- `id` (String)
- `owner` (String)
- `owner_role_type` (String)
- `session_idle_timeout_mins` (Number)
- `session_ui_idle_timeout_mins` (Number)


<a id="nestedobjatt--session_policies--show_output"></a>
### Nested Schema for `session_policies.show_output`

Read-Only:

- `comment` (String)
- `created_on` (String)
- `database_name` (String)
- `kind` (String)
- `name` (String)
- `options` (String)
- `owner` (String)
- `owner_role_type` (String)
- `schema_name` (String)
