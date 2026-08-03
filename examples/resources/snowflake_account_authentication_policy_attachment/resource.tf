resource "snowflake_authentication_policy" "default" {
  database = "prod"
  schema   = "security"
  name     = "default_policy"
}

# Attach the authentication policy account-wide (default behavior).
resource "snowflake_account_authentication_policy_attachment" "attachment" {
  authentication_policy = snowflake_authentication_policy.default.fully_qualified_name
}

resource "snowflake_authentication_policy" "service_users" {
  database = "prod"
  schema   = "security"
  name     = "service_users_policy"
}

# Attach the authentication policy to all service users only.
# Use for_all_person_users = true to target all person users instead.
# The two fields are mutually exclusive; when neither is set, the policy is attached account-wide.
resource "snowflake_account_authentication_policy_attachment" "attachment_service_users" {
  authentication_policy = snowflake_authentication_policy.service_users.fully_qualified_name
  for_all_service_users = true
}
