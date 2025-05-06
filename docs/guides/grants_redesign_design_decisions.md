---
page_title: "Grant Redesign"
subcategory: ""
description: |-

---

# Grant redesign

## Why did we decide to do a grant redesign?
Multiple factors led us to refactor grant resources, the most notable being:
- Grant problems were the majority in our GitHub issues page. We wanted to resolve all of them by providing a better solution.
- Old grant resources were made by target object and not by grant type. That led to a large number of grants to maintain. In contrast, after the refactor, we ended up with 8 resources which is a significantly lower amount in comparison to around 23 old grants (additionally, they were incomplete, and a lot more should be added to achieve full compatibility with Snowflake capabilities).
- It aligned with our goal of 100% grant feature coverage. When it comes to managing infrastructure in any way, access management is one of the most important things to consider. Thus, the user should be able to perform any granting operation possible to do manually in the worksheet.
- It’s common to use tools that perform data manipulation (like dbt) on infrastructure created by Terraform. Because of that, we should be able to perform granting commands that some of the tools may require to work properly (GRANT OWNERSHIP could be one of them).

## What are the new grant resources?
Here’s a list of resources and data sources we introduced during the grant redesign. Those resources are made to deprecate and eventually fully replace all of the previously existing grant resources.

**Resources**
- [snowflake_grant_privileges_to_database_role](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_privileges_to_database_role)
- [snowflake_grant_privileges_to_account_role](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_privileges_to_account_role)
- [snowflake_grant_account_role](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_account_role)
- [snowflake_grant_database_role](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_database_role)
- [snowflake_grant_application_role](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_application_role)
- [snowflake_grant_privileges_to_share](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_privileges_to_share)
- [snowflake_grant_ownership](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_ownership)

**Data sources**
- [snowflake_grants](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/data-sources/grants)

## Mapping from old grant resources to the new ones

As previous resources had multiple responsibilities within a single entity, we opted to divide them into more specialized resources.
Because of that, they (mostly) cannot be mapped one to one. Here's the mapping of the old resources to the new ones, but for more detailed migration rules head over to [migration notes for grant resource mappings](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#grant-resource-mappings):

- [snowflake_account_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/account_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) or/and [snowflake_grant_privileges_to_share](../docs/resources/grant_privileges_to_share) instead (depending on the configuration)
- [snowflake_database_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/database_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_external_table_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/external_table_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) or/and [snowflake_grant_privileges_to_share](../docs/resources/grant_privileges_to_share) instead (depending on the configuration)
- [snowflake_failover_group_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/failover_group_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_file_format_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/file_format_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_function_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/function_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) or/and [snowflake_grant_privileges_to_share](../docs/resources/grant_privileges_to_share) instead (depending on the configuration)
- [snowflake_grant_privileges_to_role](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/grant_privileges_to_role) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_integration_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/integration_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_masking_policy_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/masking_policy_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_materialized_view_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/materialized_view_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) or/and [snowflake_grant_privileges_to_share](../docs/resources/grant_privileges_to_share) instead (depending on the configuration)
- [snowflake_pipe_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/pipe_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_procedure_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/procedure_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) or/and [snowflake_grant_privileges_to_share](../docs/resources/grant_privileges_to_share) instead (depending on the configuration)
- [snowflake_resource_monitor_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/resource_monitor_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_role_grants](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/role_grants) - use [snowflake_grant_account_role](../docs/resources/grant_account_role) instead
- [snowflake_role_ownership_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/role_ownership_grant) - use [snowflake_grant_ownership](../docs/resources/grant_ownership) instead
- [snowflake_row_access_policy_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/row_access_policy_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_schema_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/schema_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) or/and [snowflake_grant_privileges_to_share](../docs/resources/grant_privileges_to_share) instead (depending on the configuration)
- [snowflake_sequence_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/sequence_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_stage_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/stage_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_stream_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/stream_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_table_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/table_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) or/and [snowflake_grant_privileges_to_share](../docs/resources/grant_privileges_to_share) instead (depending on the configuration)
- [snowflake_tag_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/tag_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_task_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/task_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_user_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/user_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead
- [snowflake_user_ownership_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/user_ownership_grant) - use [snowflake_grant_ownership](../docs/resources/grant_ownership) instead
- [snowflake_view_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/view_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) or/and [snowflake_grant_privileges_to_share](../docs/resources/grant_privileges_to_share) instead (depending on the configuration)
- [snowflake_warehouse_grant](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/0.90.0/docs/resources/warehouse_grant) - use [snowflake_grant_privileges_to_account_role](../docs/resources/grant_privileges_to_account_role) instead

## Design decisions

### Implicit enable_multiple_grants enabled by default
With the grant resources replacement, we wanted to change some of the default behaviors.
One of those defaults would be to make grant resources only care about privileges granted by them, mimicking the old **enable_multiple_grants** field enabled.
The motivation behind this was consistency with other resources. Other resources work in a way that they only care about themselves and manage the objects they are configured for.
Additionally, having such a destructive default could lead to some unexpected grants being revoked if someone forgets to set the flag.
Right now, there’s no alternative to the behavior **enable_multiple_grants** set to false, but we’re considering a flag for that case ([see future topics section](#future-topics)).

### Workaround for on_all and all_privileges (the always_apply parameter)
As with the **on_future** field, granting **all_privileges** or granting **on_all** also raised a few questions about tracking granted objects.
Mostly, it boiled down to running the GRANT statement whenever something changes in the Snowflake infrastructure in the case of **on_all** or whenever a privilege was added to/removed from Snowflake in the case of **all_privileges**.
We still have to discuss how and if we would like to have internal tracking of the objects affected by those commands.

For now, we decided to add the **always_apply** parameter that always produces a Terraform plan which re-grants specified privileges per terraform apply command execution.
It’s worth noting that the workaround doesn’t meet with the Terraform idea of providers having an eventually convergent state (after running the “terraform apply” the provider should eventually produce no plan).
Any user relying on this principle in their CI/CD pipelines should have this in mind when using **always_apply**.

### How should we treat the on_future parameter?
In privilege-granting resources, there’s an option to grant specific privileges on objects created in the future.
This raised a question of what we should do to the granted privileges when the resource with the specified **on_future** field is being removed.
Should we track granted privileges and revoke them or don’t track them at all? We ended up with a decision to treat the **on_future** option as a “[trigger](https://en.wikipedia.org/wiki/Database_trigger)”.
There were a lot of benefits that came with that assumption, and also it was already implemented in such a way, so it wasn’t a big surprise to the users already using that feature.
We are removing the **on_future** “trigger” on the Terraform delete operation, leaving affected grants as they were.

[Documentation Reference](https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#optional-parameters)

### Common identifier misuses
A big portion of the issues regarding grants was identifier-related (missing quotes, identifiers with special characters, etc.).
To partially resolve this issue, we introduced better validation on the schema level for most problematic identifier fields.
Some of the fields couldn’t be validated, because they are multipurpose.
Sometimes they expect account-level objects (e.g. database) and sometimes schema-level objects (e.g. table).
Even though not all of the fields are using this validation method, we don’t see those kinds of issues anymore,
indicating the better validation method and clearer error messages helped you with defining some of the grant resources.

### snowflake_grant_privileges_to_role to snowflake_grant_privileges_to_account_role
There are two main reasons why we wanted to create a successor of **snowflake_grant_privileges_to_role** called **snowflake_grant_privileges_to_account_role**:
- We wanted to make a clear distinction between account and database role resources (in the name of resources and fields).
  We are aware that there may be resources/fields which should be adjusted to follow this convention. We will be gradually changing those notifying you through [MIGRATION_GUIDE.md](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md).
- We wanted to change an internal identifier used in the **snowflake_grant_privileges_to_role** to a more import-friendly alternative.
  To follow some level of deprecation and give users time to migrate, It was easier to create another resource with a different name and internals for the same purpose, than change the existing resource and force everyone to migrate their state.
  We are aware that grants are special resources and in some cases, they may be created in the count of hundreds or thousands.
  To migrate this number of resources is not easy and should be performed gradually, and we didn’t want to block users from the latest features of the provider just because of the grant refactoring we were doing at the time.

To name a few lower-priority reasons, we also had in mind that:
- When adding a new resource we could carry out a large code structure refactor without fear of breaking anything already used by the users.
  This increased the maintainability of the resource and made it easier to grasp, so providing new functionality or fixing a bug can be done faster.
- We wanted to address all known edge cases during the refactor mentioned above, making the resource more complete.

### Why we decide to stick with one data source for grants
Grants have one show page in the documentation and were already represented as one data source. The motivation behind following this path was that:
- It's already there, so users wouldn't have to change much.
- It's easier from the usability point of view because it reflects the Snowflake documentation page, so it’s easier to have it opened on the side and create a configuration based on it.
- Performance shouldn't be a concern, because the bottlenecks that would be visible in the single data source, would also apply to the separate data sources approach.
- Even though the data source seems pretty big, the code behind it is relatively small and easy to maintain. Dividing it into smaller parts could negatively affect the maintainability.

[Documentation Reference](https://docs.snowflake.com/en/sql-reference/sql/show-grants)

### A snowflake_grant_privileges_to_application resource won’t be added
We didn’t implement the snowflake_grant_privileges_to_application resource, because granting/revoking privileges to/from application roles is only possible to perform from within the application context.

[Documentation Reference](https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-application-role#usage-notes)

### An application_role_name parameter won’t be added to the grant_application_role resource
Granting an application role to another application role can only be performed within the context of an installed application, e.g. in the application’s setup script.

[Documentation Reference](https://docs.snowflake.com/en/sql-reference/sql/grant-application-role)

### Instance roles won’t be added to the snowflake_grants data source (for now)
We didn’t add the **instance_role** field to the **snowflake_grants** data source, because they would require implementing a new type of identifier, which wouldn’t be that easy.
We decided to tackle **instance_role** identifier later because the topic of identifiers is something we would like to look into soon [after the grant redesign](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#redesigning-grants).
After the identifiers redesign, we will be in a much better position to add new identifier types or functionalities around them.

[Identifier redesign](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md#identifiers-rework)

[Documentation Reference](https://docs.snowflake.com/en/sql-reference/snowflake-db-classes)

## Grant ownership
Granting ownership was something that we discussed a long time ago.
Initially, we decided not to add it and create a document backed by an analysis that would contain reasons why it wouldn’t be possible to cover some cases.
We wanted to be careful with such decisions, and that’s why we asked for your and internal feedback ([grant ownership discussion](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/2235)).

After receiving feedback and doing a deeper analysis of certain cases, we concluded that we can and should have this resource available.
Soon after, we created a design document representing a proposal of the grant ownership resource schema, behavior, and edge cases to cover.
Currently, we are working on the implementation of the [grant ownership resource](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_ownership). After it is ready, we will announce it in our GitHub repository as well as threads mentioning grant ownership.

## Future topics
Even though the grants redesign initiative is lasting longer than we initially expected, there are still things that need to be discussed or discovered.
There may be some unusual cases where the current version of grant resources may struggle with.
With a more stable project, we will be able to do more work dedicated to discovering such edge cases.
For now, we’re relying on you to report such cases that we later analyze, prioritize, fix, and release as soon as possible.
As for the list of things yet to be discussed, we have:
- Right now, there's no way to "have **enable_multiple_grants** turned off" in the new grant resources, but we are considering adding an **authoritative** flag (the name is not chosen yet) that would work oppositely to the **enable_multiple_grants**.
  By enabling the **authoritative** flag, any other privileges granted to the target object will be revoked, making the granting resource the only source of privileges on this object.
- Discussion on granting **all_privileges** and **on_all** where we’ll decide how and if we would like to track changes of:
  - Added or removed privileges by Snowflake in the case of all_privileges
  - Added or removed objects by the user in Snowflake in the case of on_all
- Think about adding a flag to privilege-granting resources to prevent other sources from granting privileges on the same object. It would behave similarly to the enable_multiple_grants field in the old grant resources.
