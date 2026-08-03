# Account-wide attachment:
terraform import snowflake_account_session_policy_attachment.example '"<database_name>"."<schema_name>"."<session_policy_name>"|ACCOUNT'
# For all person users:
terraform import snowflake_account_session_policy_attachment.example '"<database_name>"."<schema_name>"."<session_policy_name>"|PERSON_USERS'
# For all service users:
terraform import snowflake_account_session_policy_attachment.example '"<database_name>"."<schema_name>"."<session_policy_name>"|SERVICE_USERS'
