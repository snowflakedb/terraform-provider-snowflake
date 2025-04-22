# Known issues

* [General considerations](#general-considerations-)
  * [Debugging provider](#debugging-provider)
  * [Ignore_changes meta-attribute](#ignore_changes-meta-attribute)
  * [Lack of support for the moved block](#lack-of-support-for-the-moved-block)
* [Old Terraform CLI version](#old-terraform-cli-version)
* [Errors with connection to Snowflake](#errors-with-connection-to-snowflake)
* [How to set up the connection with the private key?](#how-to-set-up-the-connection-with-the-private-key)
* [Incorrect identifier (index out of bounds) (even with the old error message)](#incorrect-identifier-index-out-of-bounds-even-with-the-old-error-message)
* [Incorrect account identifier (snowflake_database.from_share)](#incorrect-account-identifier-snowflake_databasefrom_share)
* [Granting on Functions or Procedures](#granting-on-functions-or-procedures)
* [Infinite diffs, empty privileges, errors when revoking on grant resources](#infinite-diffs-empty-privileges-errors-when-revoking-on-grant-resources)
* [Granting PUBLIC role fails](#granting-public-role-fails)
* [Issues with grant_ownership resource](#issues-with-grant_ownership)
* [Using QUOTED_IDENTIFIERS_IGNORE_CASE with the provider](#using-quoted_identifiers_ignore_case-with-the-provider)
* [Experiencing Go related issues (e.g., using Suricata-based firewalls, like AWS Network Firewall, with >=v1.0.4 version of the provider)](#experiencing-go-related-issues-eg-using-suricata-based-firewalls-like-aws-network-firewall-with-v104-version-of-the-provider)

This is a collection of the most common issues (with solutions) that users encounter when using the Snowflake Terraform Provider.

### General considerations

#### Debugging provider
To enable lower levels of logs, follow the official HashiCorp guide on [debugging providers](https://developer.hashicorp.com/terraform/internals/debugging).
The topic of debugging is further described in the [FAQ](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/FAQ.md#how-can-i-debug-the-issue-myself).

#### Ignore_changes meta-attribute
Sometimes if unexpected changes occur, you can use the `ignore_changes` meta-attribute to ignore specific fields.
This is described in the [official Terraform documentation](https://www.terraform.io/docs/language/meta-arguments/resource.html#ignore_changes).

#### Lack of support for the moved block
In the latest Terraform provider framework library, there is a new concept of the moved block.
It can be used to support migrations from deprecated resources to their new counterparts.
We are aware of this feature, but unfortunately, we cannot take advantage of it.
That's because the feature is available in the new [framework](https://developer.hashicorp.com/terraform/plugin/framework) library,
and we are still using the [older one](https://developer.hashicorp.com/terraform/plugin/sdkv2).
HashiCorp still maintains it, but no new features are added.

### Old Terraform CLI version
**Problem:** Sometimes you can get errors similar to:
```text
│ Error: Provider produced invalid plan
│
│ Provider "registry.terraform.io/snowflakedb/snowflake" planned an invalid value for snowflake_schema_grant.schema_grant.on_all: planned value cty.False for a
│ non-computed attribute.
│
│ This is a bug in the provider, which should be reported in the provider's own issue tracker.
│
```
GitHub issue reference: [#2347](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2347)

**Solution:** You have to be using at least 1.1.5 version of the Terraform CLI.

### Errors with connection to Snowflake
**Problem**: If you are getting connection errors with Snowflake error code, similar to this one:
```text
│
│ Error: open snowflake connection: 390144 (08004): JWT token is invalid.
│
```

**Related issues**: [Experiencing Go related issues (e.g., using Suricata-based firewalls, like AWS Network Firewall, with >=v1.0.4 version of the provider)](#experiencing-go-related-issues-eg-using-suricata-based-firewalls-like-aws-network-firewall-with-v104-version-of-the-provider)

**Solution**: Go to the [official Snowflake documentation](https://docs.snowflake.com/en/user-guide/key-pair-auth-troubleshooting#list-of-errors) and search by error code (390144 in this case).

GitHub issue reference: [#2432](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2432#issuecomment-1915074774)

**Problem**: Getting `Error: 260000: account is empty` error with non-empty `account` configuration after upgrading to v1 or v2, with the same provider configuration which worked up to v0.100.0

**Solution**: `account` configuration [has been removed in v1.0.0](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#removed-deprecated-objects). Please specify your organization name and account name separately as mentioned in the [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#removed-deprecated-objects):
* `account_name` (`accountname` if you're sourcing it from `config` TOML)
* `organization_name` (`organizationname` if you're sourcing it from `config` TOML)

Note that in v1.2.0, we have added a new TOML file format which renames the TOML fields to be consistent with the provider schema (see the [migration guide entry](https://github.com/snowflakedb/terraform-provider-snowflake/blob/dev/MIGRATION_GUIDE.md#new-toml-file-schema) for v1.2.0). In v2, this format is used by default and it is recommended (see the [migration guide entry](https://github.com/snowflakedb/terraform-provider-snowflake/blob/dev/MIGRATION_GUIDE.md#breaking-change-changes-in-default-toml-format) for v2.0.0). In this case, the TOML configuration looks like:
```toml
[example]
account_name = 'account_name'
organization_name = 'organization_name'
user = 'user'
...
```

GitHub issue reference: [#3198](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3198), [#3308](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3308)

### How to set up the connection with the private key?
**Problem:** From the version v0.78.0, we introduced a lot of provider configuration changes. One of them was deprecating `private_key_path` in favor of `private_key`.

GitHub issue reference: [#2489](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2489), [Migration Guide reference](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#v0730--v0740)

**Solution:** Use a non-deprecated `private_key` field with the use of the [file](https://developer.hashicorp.com/terraform/language/functions/file) function to pass the private key.

### Incorrect identifier (index out of bounds) (even with the old error message)
**Problem:** When getting stack traces similar to:
```text
│ panic: runtime error: index out of range [2] with length 2
│
│ goroutine 61 [running]:
│ github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake.SchemaObjectIdentifierFromQualifiedName({0x140001c2870?, 0x103adf987?})
│ github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake/identifier.go:58 +0x174
```

GitHub issue reference: [#2224](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2224)

**Solution:** Some fields may expect different types of identifiers, when in doubt check [our documentation](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs) for the field or the [official Snowflake documentation](https://docs.snowflake.com/) what type of identifier is needed.

### Incorrect identifier type (panic: interface conversion)
**Problem** When getting stack traces similar to:
```text
panic: interface conversion: sdk.ObjectIdentifier is sdk.AccountObjectIdentifier, not sdk.DatabaseObjectIdentifier
```

GitHub issue reference: [#2779](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2779)

**Solution:** Some fields may expect different types of identifiers, when in doubt check [our documentation](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs) for the field or the [official Snowflake documentation](https://docs.snowflake.com/) what type of identifier is needed. Quick reference:
- AccountObjectIdentifier - `<name>`
- DatabaseObjectIdentifier - `<database>.<name>`
- SchemaObjectIdentifier - `<database>.<schema>.<name>`
- TableColumnIdentifier - `<database>.<schema>.<table>.<name>`

### Incorrect account identifier (snowflake_database.from_share)
**Problem:** From 0.87.0 version, we are quoting incoming external account identifier correctly, which may break configurations that specified account identifier as `<org_name>.<acc_name>` that worked previously by accident.

GitHub issue reference: [#2590](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2590)

**Solution:** As specified in the [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#behavior-change-external-object-identifier-changes), use account locator instead.

### Granting on Functions or Procedures
**Problem:** Right now, when granting any privilege on Function or Procedure with this or similar configuration:

```terraform
resource "snowflake_grant_privileges_to_account_role" "grant_on_procedure" {
  privileges        = ["USAGE"]
  account_role_name = snowflake_account_role.name
  on_schema_object {
    object_type = "PROCEDURE"
    object_name = "\"${snowflake_database.database.name}\".\"${snowflake_schema.schema.name}\".\"${snowflake_procedure_sql.procedure.name}\""
  }
}
```

You may encounter the following error:
```text
│ Error: 090208 (42601): Argument types of function 'procedure_name' must be
│ specified.
```

**Related issues:** [#2375](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2375), [#2922](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2922)

**Solution:** Specify the arguments in the `object_name`:

```terraform
resource "snowflake_grant_privileges_to_account_role" "grant_on_procedure" {
  privileges        = ["USAGE"]
  account_role_name = snowflake_account_role.name
  on_schema_object {
    object_type = "PROCEDURE"
    object_name = "\"${snowflake_database.database.name}\".\"${snowflake_schema.schema.name}\".\"${snowflake_procedure_sql.procedure.name}\"(NUMBER, VARCHAR)"
  }
}
```

If you manage the procedure in Terraform, you can use `fully_qualified_name` field:

```terraform
resource "snowflake_grant_privileges_to_account_role" "grant_on_procedure" {
  privileges        = ["USAGE"]
  account_role_name = snowflake_account_role.name
  on_schema_object {
    object_type = "PROCEDURE"
    object_name = snowflake_procedure_sql.procedure_name.fully_qualified_name
  }
}
```

### Infinite diffs, empty privileges, errors when revoking on grant resources
**Problem:** If you encountered one of the following issues:
- Issue with revoking: `Error: An error occurred when revoking privileges from an account role.
- Plan in every `terraform plan` run (mostly empty privileges)
It's possible that the `object_type` you specified is "incorrect."
Let's say you would like to grant `SELECT` on event table. In Snowflake, it's possible to specify
`TABLE` object type instead of dedicated `EVENT TABLE` one. As `object_type` is one of the fields
we filter on, it needs to exactly match with the output provided by `SHOW GRANTS` command.

**Related issues:** [#2749](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2749), [#2803](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2803)

**Solution:** Here's a list of things that may help with your issue:
- Firstly, check if the privilege has been granted in Snowflake. If it is, it means the configuration is correct (or at least compliant with Snowflake syntax).
- When granting `IMPORTED PRIVILEGES` on `SNOWFLAKE` database/application, use `object_type = "DATABASE"`.
- Run `SHOW GRANTS` command with the right filters to find the granted privilege and check what is the object type returned of that command. If it doesn't match the one you have in your configuration, then follow those steps:
  - Use state manipulation (no revoking)
    - Remove the resource from your state using `terraform state rm`.
    - Change the `object_type` to correct value.
    - Import the state from Snowflake using `terraform import`.
  - Remove the grant configuration and after `terraform apply` put it back with the correct `object_type` (requires revoking).

### Granting PUBLIC role fails
**Problem:** When you try granting PUBLIC role, like:
```terraform
resource "snowflake_account_role" "any_role" {
  name = "ANY_ROLE"
}

resource "snowflake_grant_account_role" "this_is_a_bug" {
  parent_role_name = snowflake_account_role.any_role.name
  role_name        = "PUBLIC"
}
```
Terraform may fail with:
```
╷
│ Error: Provider produced inconsistent result after apply
│
│ When applying changes to snowflake_grant_account_role.this_is_a_bug, provider "provider["registry.terraform.io/snowflakedb/snowflake"]" produced an
│ unexpected new value: Root object was present, but now absent.
│
│ This is a bug in the provider, which should be reported in the provider's own issue tracker.
╵
```

**Related issues:** [#3001](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3001), [#2848](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2848)

**Solution:** This happens, because the PUBLIC role is a "pseudo-role" (see [docs](https://docs.snowflake.com/en/user-guide/security-access-control-overview#system-defined-roles)) that is already assigned to every role and user, so there is no need to grant it through Terraform. If you have an issue with removing the resources please use `terraform state rm <id>` to remove the resource from the state (and you can safely remove the configuration).

### Issues with grant_ownership

Please read our [guide for grant_ownership](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/grant_ownership_common_use_cases) resource.
It contains common use cases and issues that you may encounter when dealing with ownership transfers.

### Using QUOTED_IDENTIFIERS_IGNORE_CASE with the provider

**Problem:** When `QUOTED_IDENTIFIERS_IGNORE_CASE` parameter is set to true, but resource identifier fields are filled with lowercase letters,
during `terrform apply` they may fail with the `Error: Provider produced inconsistent result after apply` error (removing themselves from the state in the meantime).

**Related issues:** [#2967](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2967)

**Solution:** Either turn off the parameter or adjust your configurations to use only upper-cased names for identifiers and import back the resources.

### Experiencing Go related issues (e.g., using Suricata-based firewalls, like AWS Network Firewall, with >=v1.0.4 version of the provider)

**Problem:** The communication from the provider is dropped, because of the firewalls that are incompatible with the `tlskyber` setting introduced in [Go v1.23](https://go.dev/doc/godebug#go-123).

**Related issues:** [#3421](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3421)

**Solution:** [Solution described in the migration guide for v1.0.3 to v1.0.4 upgrade](./MIGRATION_GUIDE.md#new-go-version-and-conflicts-with-suricata-based-firewalls-like-aws-network-firewall).

### Provider is too slow

**Problem:** The provider is taking too long to perform plan/apply operations.

**Solution:** Refer to our [performance analysis](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/performance_benchmarks) and the optimizations we are proposing.

### Dropping related resources

**Problem:** Sometimes you may seem to have issues with dropping related resources, like in the example below:
```text
│ Error deleting network policy EXAMPLE, err = 001492 (42601): SQL compilation error:
│ Cannot perform Drop operation on network policy EXAMPLE. The policy is attached to INTEGRATION with name EXAMPLE. Unset the network policy from INTEGRATION and try the
│ Drop operation again.
```
That is because some of the Snowflake objects are interdependent, and dropping one requires dropping the other first.

**Solution:** This should be mostly resolved by keeping the dependencies between resources in the configuration code, but also make sure to check our [guide regarding this topic](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/unassigning_policies).
