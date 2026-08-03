resource "snowflake_session_policy" "sp" {
  database = "prod"
  schema   = "security"
  name     = "default_session_policy"
}

# Attach the session policy account-wide (default behavior).
resource "snowflake_account_session_policy_attachment" "attachment" {
  session_policy_name = snowflake_session_policy.sp.fully_qualified_name
}

resource "snowflake_session_policy" "service_users" {
  database = "prod"
  schema   = "security"
  name     = "service_users_session_policy"
}

# Attach the session policy to all service users only.
# Use for_all_person_users = true to target all person users instead.
# The two fields are mutually exclusive; when neither is set, the policy is attached account-wide.
resource "snowflake_account_session_policy_attachment" "attachment_service_users" {
  session_policy_name   = snowflake_session_policy.service_users.fully_qualified_name
  for_all_service_users = true
}
