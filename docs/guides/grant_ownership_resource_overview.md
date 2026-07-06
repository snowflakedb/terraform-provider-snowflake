---
page_title: "Grant Ownership"
subcategory: ""
description: |-

---

# Grant ownership

The [grant\_ownership resource](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/resources/grant_ownership) was introduced in version 0.88.0.
Since its release, feedback indicates that it can be challenging to understand and use effectively in certain scenarios.
We would like to give an overview of the grant\_ownership resource so that its limitations are clearly understood and provide some guidance on how it should be used, and what is planned for the future.

Before we get into the details, we want to explain why we were initially hesitant to add this feature to the provider.

During the redesign of grants, we evaluated the use cases for granting ownership and initially decided not to include it.
Terraform is meant to manage the infrastructure. To manage given objects in Snowflake, it’s often required to be the owner of those objects.
This was our main concern; making it work would require numerous assumptions, rules, and, potentially, manual steps, making the resource difficult to manage and use.

To ensure our decision didn't limit customers, we [asked them to share scenarios](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/2235)
where granting ownership is crucial and a role-based approach is not feasible.
After reviewing these use cases, we decided to offer this resource, but with only essential functionalities to keep it simple while meeting necessary requirements.

Over time, we've found it challenging to use, especially when debugging role-based access errors.
Therefore, in the coming week, we will provide examples for common use cases and error handling to help resolve most frustrations that come up when using the grant\_ownership resource.

## Limitations and workarounds

Most of the limitations and usage errors currently present in the grant\_ownership resource result from its uniqueness in both Snowflake and Terraform.
Here’s the list of things that are the most challenging about the resource and how you can mitigate the common issues that result from them.

### Resource as a connection

The resource doesn’t represent an infrastructure object but rather a connection between objects it doesn’t own.
This means that changes in other objects may affect it. This is also true for some of the other resources,
but the difference with granting ownership is that the number of requirements for this operation is sometimes much bigger than in resources of a similar type.
This leads to failures, especially when the dependencies inside the configuration are set up incorrectly.
Usually by keeping the correct order of configuration execution (by either using an [implicit resource reference](https://developer.hashicorp.com/terraform/tutorials/configuration-language/dependencies#manage-implicit-dependencies)
or [explicit depends\_on meta-argument](https://developer.hashicorp.com/terraform/tutorials/configuration-language/dependencies#manage-explicit-dependencies)),
you can achieve more predictable results (e.g. [\#3253](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3253)).

### Object tracking

Partially to the point above, because this resource represents a connection, we should track the affected objects, so that after the resource is deleted the operation could be reversed.
Similarly to the other grant resources, this is unfortunately not true for the \`on\_all\` option.
This may be dangerous if you, for example, grant ownership on all tables in database X with the resource, then create some tables manually, and at the end delete the resource.
As a result, when the provider grants ownership back to the calling role, the manually created tables will be also affected.
Currently, there’s no easy way to achieve this kind of tracking, and the \`on\_all\` option should be used with caution to prevent any unnecessary ownership transfers.

### Grant ownership caller restrictions/requirements

To grant the ownership, you have to own the object, or the current role has to be granted the MANAGE GRANTS privilege.
If the currently used role is not granted the MANAGE GRANTS privilege and doesn’t have access to the granted object (outside of owning it initially),
it can cause an error when deleting the resource as the user is not privileged to get back the ownership
(e.g. [\#3220](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3220), [\#3317 comment](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3317#issuecomment-2593541448)).
For now, the best practice is to either use a highly privileged role like ACCOUNTADMIN or a role that has MANAGE GRANTS granted.

### Grant ownership object restrictions/requirements

It is not trivial to automatically transfer the ownership of a pipe or task.
Both objects have their own rules and privileges that you have to have to successfully transfer ownership.
Most of the rules are described in the [official documentation for grant ownership,](https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#usage-notes)
which we highly recommend reading. Due to the complexity and risk of partial or full operation failures during creation or deletion,
it is advised to use ownership grants on these objects with caution and report any unexpected errors.

## Future plans

We already have some ideas about how to expand this resource, but we wanted to start with only essential features and extend it according to further customer requirements.
We are considering introducing the following changes (the order doesn’t matter).

### Research the effect and ways to define outbound privileges during create and delete operations

Currently, the outbound privileges are always taken from the configuration and marked as [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew).
We wanted to analyze and see if some of the use cases require an additional parameter that would be responsible for setting different outbound privilege options in the delete operation than in the create operation.
Right now, it’s not possible to call different outbound privilege options during create and delete.

### Granting ownership back to the original role instead of the current one

To support ownership transferring from/to different roles than the current one,
we could add a parameter to the resource that would be responsible for granting ownership to the passed role during the delete operation.

### Granting ownership back to the original role with the use of the granted role

To eliminate the need for highly privileged role usage,we consider providing additional functionality that would enable the provider to use only the role that owns the given object to transfer ownership ([use case](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3317#issuecomment-2593541448)).
This introduces an additional assumption (or internal provider validation) that the user used for running the Terraform configuration is granted the role that initially owns the object and the role to which the ownership is transferred.

## Summary

When we embarked on our journey with the grant ownership resource, we anticipated encountering some challenges due to the numerous assumptions, rules, and limitations imposed by Snowflake and Terraform.
We hope that the additional explanations will aid in its usage and encourage you to provide feedback. This will allow us to enhance the documentation further.
