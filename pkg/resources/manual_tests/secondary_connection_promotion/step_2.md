To test the secondary connection promotion, run the following SQL in the account where secondary connection is created:
```sql
ALTER CONNECTION <connection_name> PRIMARY;
```

Run `terraform plan` with the same configuration as in the first step.

> Output expectations: You should see both primary and secondary connections plan for re-creation.

After running plan, even if you get back to the previous state where other connection was primary by running in the default account:
```sql
ALTER CONNECTION <connection_name> PRIMARY;
```
You will still see secondary connection re-creation in the plan.

This is because earlier (in 2.5.0 version) we were reading connections incorrectly. 
To fix the issue, you have to bump to at least 2.6.0 version of the provider (or install it locally if not yet released).
Without it, the following steps may not work as expected.

Now, there are two options to test:
1. Fix the state so that we can get back to the original configuration.
2. Proceed by switching primary connection to the secondary one and secondary to primary.

## Option 1: Fix the state

The version bump alone should resolve the issue with the state and terraform plan should show empty plan for both primary and secondary connections.

## Option 2: Switch primary and secondary connections

1. To switch primary and secondary connections, you have to firstly remove them from the state with the following commands:

```bash
terraform state rm snowflake_primary_connection.primary_connection
terraform state rm snowflake_secondary_connection.secondary_connection
```

2. Then, you have to refactor the configuration to switch primary and secondary connections:

```terraform
resource "snowflake_primary_connection" "primary_connection" {
  provider = snowflake.second_account
  name     = "TEST_CONNECTION"
  enable_failover_to_accounts = ["<organization_name>.<seconary_connection_account_name>"]
}

resource "snowflake_secondary_connection" "secondary_connection" {
  name          = snowflake_primary_connection.primary_connection.name
  as_replica_of = "<organization_name>.<primary_connection_account_name>.${snowflake_primary_connection.primary_connection.name}"
}
```

3. After that, you can re-import the resources with the following commands:

```bash
terraform import snowflake_primary_connection.primary_connection '"TEST_CONNECTION"'
terraform import snowflake_secondary_connection.secondary_connection '"TEST_CONNECTION"'
```

> Output expectations: You should see empty plan for both primary and secondary connections.
