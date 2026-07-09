---
page_title: "Object Renaming Guide"
subcategory: ""
description: |-

---

# Object Renaming Guide

Recently, we conducted research on object renaming and published a document summarizing the results.
To leverage the knowledge we gained from this research, we wanted to provide a follow-up document that would help you understand the current best practices for tackling object renaming-related topics.
In this document, we propose recommendations and solutions for the issues identified through our research, as well as those previously reported in our GitHub repository.

## Renaming objects in the hierarchy

This problem relates to renaming objects that are higher in the object hierarchy (e.g. database or schema) and how this affects the lower hierarchy objects created on them (e.g. schema or table) while they are present in the Terraform configuration.
In our research, we tested different sets of configurations described [here](./object_renaming_research_summary#renaming-higher-hierarchy-objects).

### Recommendations

The primary recommendation is to keep your objects in correct relations. Use the following order:
- [Implicit dependency](https://developer.hashicorp.com/terraform/tutorials/configuration-language/dependencies#manage-implicit-dependencies)
- [Explicit dependency (depends_on)](https://developer.hashicorp.com/terraform/tutorials/configuration-language/dependencies#manage-explicit-dependencies)
- No dependency

Maintaining proper resource dependencies is essential for the provider to accurately determine the appropriate actions when a high-level object is renamed.

If you prefer to handle hierarchy renames without resource recreation, consider enabling the `HIERARCHY_RENAMES` experimental feature (see below).
Otherwise, if you need to perform a database rename with other resources referencing its name, you can first remove the dependent objects from the state.
Then, perform the actual rename, and after that, you can import the dependent objects back to the state, but with a different database.
This is very time-consuming, so only consider this when the number of objects dependent on the object you want to rename is low.
To see more or less how this could be implemented, take a look at the [migration guide](./resource_migration) we already described which has similar steps of execution.

### Experimental: In-place hierarchy renames

The `HIERARCHY_RENAMES` experiment enables in-place handling of hierarchy renames and moves for supported resources.
To enable it, add `HIERARCHY_RENAMES` to the `experimental_features_enabled` list in your provider configuration:

```hcl
provider "snowflake" {
  experimental_features_enabled = ["HIERARCHY_RENAMES"]
}
```

Currently supported by: `snowflake_schema`, `snowflake_table`.

When enabled, changing the parent identifier fields (e.g. `database` on `snowflake_schema`, or `database`/`schema` on `snowflake_table`) no longer forces resource recreation. Instead, the provider detects the scenario and handles it accordingly.

#### 2-level hierarchy (database-level objects)

This applies to resources like `snowflake_schema` where the object lives directly inside a database.

##### Database rename (parent renamed)

If the parent database was renamed (e.g. from `A` to `B`), the provider detects that the old database no longer exists while the new database and schema already exist under the new name. In this case, the provider simply updates the resource ID to reflect the new database name — no Snowflake modification is performed.

**Conditions detected:**
- New database `B` exists
- Old database `A` does not exist
- Schema `B.X` exists

##### Schema move (move to a different database)

If both the old and new databases exist, the provider treats this as a request to move the schema from one database to another. It executes `ALTER SCHEMA A.X RENAME TO B.X` in Snowflake, then updates the resource ID.

**Conditions detected:**
- New database `B` exists
- Old database `A` exists
- Schema `A.X` exists

#### 3-level hierarchy (schema-level objects)

This applies to resources like `snowflake_table` where the object lives inside a schema, which itself lives inside a database. When either or both of the `database` and `schema` fields change, the provider determines the correct action.

##### Only database changes

The provider checks the database level (same logic as the 2-level case above) to determine whether the database was renamed or the table should be moved to a different database.

**Database rename:** Old database does not exist, new database exists, and the schema (with original name) exists in the new database. The provider updates the resource ID only — no SQL is executed.

**Table move:** Both databases exist and the schema (with original name) exists in the old database. The provider executes `ALTER TABLE A.S.T RENAME TO B.S.T`.

##### Only schema changes

The provider checks the schema level to determine whether the schema was renamed or the table should be moved to a different schema.

**Schema rename:** Old schema does not exist, new schema exists, and the table already exists at the target location. The provider updates the resource ID only — no SQL is executed.

**Table move:** Both schemas exist and the table still exists at the source location. The provider executes `ALTER TABLE D.A.T RENAME TO D.B.T`.

##### Both database and schema change

When both fields change simultaneously, the provider evaluates all combinations of database and schema existence to determine which of the following four scenarios applies:

| Scenario                  | Database                   | Schema                                       | Action                        |
|---------------------------|----------------------------|----------------------------------------------|-------------------------------|
| DB rename + Schema rename | Old DB gone, new DB exists | Schema with new name exists in new DB        | Update resource ID only       |
| DB rename + Schema move   | Old DB gone, new DB exists | Schema with old name exists in new DB        | Execute `ALTER TABLE` to move |
| DB move + Schema rename   | Both DBs exist             | Old schema gone, new schema exists in old DB | Execute `ALTER TABLE` to move |
| DB move + Schema move     | Both DBs exist             | Old schema exists in old DB                  | Execute `ALTER TABLE` to move |

#### Detection logic

For each scenario, the provider probes the existence of parent objects and the object itself using `SHOW` commands:

1. **Parent probing** — checks whether the old and new parent objects (database, schema) still exist.
2. **Object probing** — checks whether the object exists at the source and/or target location.

Based on the combination of existence checks, the provider determines whether a parent was renamed (requiring only an ID update) or whether the object should be actively moved (requiring an `ALTER ... RENAME TO` statement).

#### Error handling

If the provider cannot determine the rename scenario (e.g. neither the old nor the new parent exists, or the expected object locations do not match any known pattern), it returns an error. The error message includes the existence status of all probed objects and a link to this guide for reference.

#### Requirements

- Proper resource dependencies (implicit or explicit) between parent and child resources are strongly recommended for correct behavior.
- A simultaneous name change of the object itself (e.g. changing `name` together with `database`) is supported — the provider handles the hierarchy change first, then applies the name change.
- If the provider cannot determine the rename scenario, it returns an error with guidance rather than performing an incorrect action.

## Issues with lists and sets

Currently, we have limited capabilities when it comes to certain operations on lists and sets.
An example of such a limitation could be detecting whether a collection item was updated or one item was removed and the new one was put in its place.
This is mainly due to how the Terraform SDKv2 handles changes for collections.
So far, the most challenging case was columns on tables, as Snowflake has its own limitations preventing us from reaching the correct state.
Here are some of the issues pointing to the limitations we are talking about:
- [terraform-plugin-sdk#133](https://github.com/hashicorp/terraform-plugin-sdk/issues/133)
- [terraform-plugin-sdk#196](https://github.com/hashicorp/terraform-plugin-sdk/issues/196) (this is regarding the testing framework, but the issue persists on the provider-level code as well)
- [terraform-plugin-sdk#447](https://github.com/hashicorp/terraform-plugin-sdk/issues/447)
- [terraform-plugin-sdk#1103](https://github.com/hashicorp/terraform-plugin-sdk/issues/1103)

There is more, but the real issue is that those problems overlap, making it really difficult to provide any custom functionality that wasn’t considered when designing the Terraform SDKv2.

### Recommendations

It's important to align your needs with the capabilities of the provider's resources and choose the appropriate tool for the task.
This is particularly crucial for lower-level objects like tables, which are subject to frequent changes and may pose challenges when being provisioned in Terraform.
Tables are unique as they are infrastructure objects that contain data, so modifications need to be considered carefully.
Due to current limitations, it might be impractical to provision tables with the provider, as some table parameter changes require dropping and recreating the table, resulting in data loss.
In Terraform, this approach is common to ultimately achieve the desired infrastructure state with the specified objects.
After the research, we have some upcoming improvements in handling changes in lists and sets, but they won’t resolve all the issues, and the above remains.

### Future plans

As mentioned in the [research summary](./object_renaming_research_summary#ignoring-list-order-after-creation--updating-list-items-mostly-related-to-table-columns), we plan to improve the table resource with all the findings, which will mostly affect the list of columns and how we detect/plan changes for them.
Once implemented, all the details will be available in the documentation for the table resource and in the [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md).

## Summary

We hope that the additional recommendations derived from our research will assist you in making informed decisions regarding the use of our provider.
If you have any questions or need further clarification, we encourage you to create issues in our [GitHub repository](https://github.com/snowflakedb/terraform-provider-snowflake).
Your feedback is invaluable and will contribute to further improving our documentation.
