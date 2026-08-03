# Account-wide attachment:
terraform import snowflake_account_authentication_policy_attachment.example '"<database_name>"."<schema_name>"."<authentication_policy_name>"|ACCOUNT'
# For all person users:
terraform import snowflake_account_authentication_policy_attachment.example '"<database_name>"."<schema_name>"."<authentication_policy_name>"|PERSON_USERS'
# For all service users:
terraform import snowflake_account_authentication_policy_attachment.example '"<database_name>"."<schema_name>"."<authentication_policy_name>"|SERVICE_USERS'
