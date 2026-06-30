---
page_title: "Object Renaming"
subcategory: ""
description: |-

---

# Object Renaming Support

The Terraform Provider team recently started a short research project on object renaming and other similar topics. This document will cover the topics we looked into, explain how we tested them, and discuss their effects on the provider. We'll also list the topics we want to explore more in our next research.

## Topics

### Renaming higher-hierarchy objects

**Description:** This problem relates to renaming objects that are higher in the object hierarchy (e.g. database or schema) and how this affects the lower hierarchy objects created on them (e.g. schema or table) while they are present in the Terraform configuration. We decided to deeply test this problem, as from time to time we got issues related to it. We wanted to get a better understanding of it and how currently our provider is handling such situations to provide appropriate fixes if necessary.

**Tests:** We prepared a [set of test cases](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/pkg/resources/object_renaming_acceptance_test.go) by combining permutations of different aspects like:

- Depth (shallow connections like database and schema or deep connections like database, schema, and table).
- Higher hierarchy object placement (inside or outside of the Terraform configuration)
- Resource dependency (implicit, [depends\_on](https://developer.hashicorp.com/terraform/language/meta-arguments/depends_on), or no dependency)
- Place of rename execution (within the Terraform configuration or manually outside)

**Impact:** We decided that we will provide additional documentation on:

- Best practices
- Guide on dealing with certain errors connected to object renaming
- Guidelines that may be useful in certain scenarios connected to object renaming

In addition to improved documentation, the tests showed us that we need to improve our error handling in the Read and Delete operations. In certain scenarios, the resources failed to remove themselves from the state when they should. This change should decrease the chances of resources trapping themselves in the infinite plan state where the only way out is through [manual state manipulation](https://developer.hashicorp.com/terraform/cli/commands/state).

**Implementation:** Based on the research findings, we implemented the `HIERARCHY_RENAMES` experimental feature that enables in-place handling of hierarchy renames and moves. The feature currently supports both 2-level hierarchies (e.g. `snowflake_schema` within a database) and 3-level hierarchies (e.g. `snowflake_table` within a schema within a database). The provider automatically detects whether a parent object was renamed or whether the child object should be moved, and handles each case without resource recreation. For full details on the supported scenarios and how to enable the feature, see the [object renaming guide](./object_renaming_guide#experimental-in-place-hierarchy-renames).

### Ignoring list order after creation \+ updating list items (mostly related to table columns)

**Description:** The issues with table columns ([\#420](https://github.com/snowflakedb/terraform-provider-snowflake/issues/420), [\#753](https://github.com/snowflakedb/terraform-provider-snowflake/issues/753), [\#2839](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2839)) are something we had in mind for a very long time and finally had a chance to work on a solution that would improve them, and possibly some of the other resources. In short, the use case for table columns consists of two use cases connected to each other:

- Ignoring the order of columns after creation. Users should be able to reorder, add, and remove columns from any place while still having somewhat control over column order on the Snowflake side.
- Updating a given column instead of removing and adding it again. Ignoring the order was an additional challenge because if someone wants to order the columns and change their name in one apply, then we need a way to identify this column to perform the correct action.

**Tests:** The tests were carried out on a resource created only for the purpose of the research ([resource reference](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/pkg/resources/object_renaming_lists_and_sets.go#L125)). Note that it won’t be normally visible and no other changes in other resources were made. We tested a few approaches regarding order ignoring and one on updating items.

**Impact:** The outcomes of the tests showed promising results and potential improvement in how structures like table columns are managed. The adjustment of table columns will be done during the refactoring of the table we do as part of [preparing the GA object for V1](https://github.com/snowflakedb/terraform-provider-snowflake/pull/3147).

Additionally, this gives us more knowledge of how the lists and their items are managed in Terraform SDKv2 and how we can interact with them to achieve certain behaviors. We confirmed their limitations, and found solutions on how to deal with most of them. We also have high hopes that once we can migrate to the newer Terraform Plugin Framework, there will be more tools to support even more demanding use cases.

## Topics for future research

### Diff suppression for the items of lists and sets

We would like to research the usage of DiffSuppressFuncs for items on lists and sets. With some of the resources we found that using them is sometimes tricky and may cause issues with the plan. This could enable us to suppress differences for e.g. quotes in sets of identifiers.

### Computed \+ Optional lists (and sets)

During this research, the concept was only briefly touched on due to its complexity and the need for attention to other topics. However, they may be useful for cases where the list (or set) is optional and is computed on the Snowflake side when not specified. This is an example of recently refactored views where this approach could be used, but it needs further testing before using it in the actual resource.

### Item updates in sets

In Terraform, the indexes of set items are calculated based on the item’s hash. Because of that, it’s hard to handle an item’s update whenever one of the items changes. The topic wasn’t covered in this research, because there’s no real use case for that (yet), but we were already thinking about potentially switching some of the fields from lists to sets where this feature could be useful.

## Summary

This journey has been valuable, enhancing our understanding and guiding future improvements. We also outlined areas for further research, which we believe will bring even more benefits to our users. We are excited to keep improving your experience with our Terraform Provider.
