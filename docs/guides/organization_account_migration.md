---
page_title: "Organization Account Migration"
subcategory: ""
description: |-

---

# Organization Account Migration

## Migrating your configuration for the upgraded organization-level management

Snowflake recently introduced a new [organization management paradigm](https://docs.snowflake.com/en/user-guide/organizations).
The former model is based on the ORGADMIN role in an account.
The new approach transitions org-wide management responsibilities to the organization account and the new GLOBALORGADMIN
role, centralizing control and introducing greater granularity and auditability at the organization level.

Recognizing Terraform's suitability for high-level object management,
we've integrated new functionalities related to the new organization-level management,
by adding the [current_organization_account](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/current_organization_account) resource,
with organization_account to follow. Further enhancements are still planned and will be announced as a
new [roadmap](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md) entry.
We will announce it as a GitHub discussion, as we have done [previously](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/3703).

This guide outlines the migration process for Snowflake organization-level configurations from the ones managed by
ORGADMIN to the organization accounts managed by GLOBALORGADMIN.

For more information on organization management, please visit the [official Snowflake documentation](https://docs.snowflake.com/en/user-guide/organizations).

## Why you should migrate

Beyond the benefits detailed in [Snowflake's documentation](https://docs.snowflake.com/en/user-guide/organizations#benefits),
the migration is important due to future plans to eventually [phase out ORGADMIN support](https://docs.snowflake.com/en/user-guide/organization-administrators#using-the-orgadmin-role)
for multi-account organizations. Organization accounts introduce new features
like [organization users](https://docs.snowflake.com/en/user-guide/organization-users) and [organization user groups](https://docs.snowflake.com/en/user-guide/organization-users#organization-user-groups).
We intend to integrate these capabilities into the provider, allowing you to manage your organizations within the Terraform ecosystem.
However, these features are exclusive to organization accounts, a restriction that also applies to the provider, meaning they will be usable only after migration.
This restriction will also apply to other newly introduced features that further enhance organization-level management in Snowflake.

## The migration process

### Creating an organization account

*Note: If you are already using the organization account in your organization, you can go to the [next section](#migration).*

Use the [CREATE ORGANIZATION ACCOUNT](https://docs.snowflake.com/en/sql-reference/sql/create-organization-account) command from one of the accounts using the ORGADMIN approach either within Snowsight or through the provider’s [execute resource](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/execute).

*Note: We’re currently planning to introduce the `organization_account` resource so that you can create an organization
account with the provider; however, it’s not our highest priority regarding organization-level management features,
as modifications of organization accounts are possible with the [current_organization_account](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/current_organization_account) resource.*

### Migration

#### 1. Identify the provider configurations

*Note: By [provider configuration](https://developer.hashicorp.com/terraform/language/providers/configuration), we mean provider blocks that configure the Snowflake Terraform provider.
Each resource within the Terraform configuration is connected to a given provider configuration.*

Prior to migration, identify provider configurations that utilize either the ORGADMIN role or a role granted with ORGADMIN.

#### 2. Check the objects managed in those configurations

The ORGADMIN's one of the main purposes is to manage accounts (regular ones).
Therefore, provider configurations using ORGADMIN should primarily manage [account resources](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/account).
These are the key objects we will migrate to shift the responsibility of managing these accounts from ORGADMIN to
GLOBALORGADMIN (organization account).

To ensure other resources remain unaffected during the migration, consider separating ORGADMIN provider configurations
that manage objects other than accounts. For further assistance in this regard, consult the [resource migration guide](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/resource_migration).

#### 3. Migrate your provider configurations

As mentioned above, we’ll cover migration for provider configurations containing accounts because that’s the main
purpose of ORGADMIN roles. Take the following configuration as an example configuration you’ll be starting with.
It contains a simple provider configuration using an account with the ORGADMIN role and one account that it manages:

```terraform
provider "snowflake" {
  role    = "ORGADMIN"
  profile = "org_admin_account"
}

resource "snowflake_account" "test" {
  name                 = "account_name"
  admin_name           = "admin_name"
  admin_password = "pass"
  # or`
  admin_rsa_public_key = "rsa_public_key"
  email                = "exampl@email.com"
  edition              = "ENTERPRISE"
  grace_period_in_days = 3
}
```

To simply migrate this configuration, all you need to do is change the provider configuration to match the organization
account credentials and make sure a proper role is used:

```terraform
provider "snowflake" {
  role    = "GLOBALORGADMIN"
  profile = "organization_account"
}
```

Run the terraform plan command to confirm that the migration has been carried out successfully and no changes are planned.
This migration is possible to do this way because the [account](https://docs.snowflake.com/en/sql-reference/sql/create-account) object, unlike other ones in Snowflake,
doesn’t track ownership, so unless you have the right privileges, you can just switch provider configurations as
presented (even if it’s pointing to another account). After adjusting all configurations in a similar fashion,
you should be done with the migration.

#### 4. (Optional) Manage your organization account within the provider

If you haven't already, consider incorporating your organization's account configuration into your Terraform configuration.
To do this, use the newly introduced [current_organization_account](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/current_organization_account) resource.
Remember that you should have only one such resource per organization;
otherwise, multiple instances may compete against each other to set the correct state for the organization account.

## Questions

If you still have questions or suggestions regarding the provided migration guide, please contact us through GitHub by [creating a new issue](https://github.com/snowflakedb/terraform-provider-snowflake/issues/new/choose).
