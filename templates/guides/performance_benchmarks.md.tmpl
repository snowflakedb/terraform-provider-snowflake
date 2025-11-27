---
page_title: "Performance Analysis"
subcategory: ""
description: |-

---

# Performance Analysis

This document provides a basic performance analysis of the Snowflake Terraform Provider. It is not a complete analysis, but a basic outline, allowing us to give a few recommendations for the current provider versions. The document’s purpose is to set performance expectations for using the provider and give suggestions to users on how to improve its performance. We decided to perform such benchmarks because of concerns reported by our users ([\#3118](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3118), [\#3169](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3169)). They are related to the performance of large workloads (a few thousand resources). These issues have been reported only in recent versions because of the [changes in the reworked objects](https://github.com/snowflakedb/terraform-provider-snowflake/blob/c4b1bebce4bc5a81031248592b34af5e80ca2fc1/v1-preparations/CHANGES_BEFORE_V1.md) (more queries and bigger state sizes in some resources).

## Methodology

We prepared Terraform configurations in [our repository](https://github.com/snowflakedb/terraform-provider-snowflake/tree/main/pkg/manual_tests/benchmarks). We tested different deployment sizes for total execution time and state size. We chose multiple resources, such as tasks, warehouses, and schemas, to test account-level and database-level objects. These resources provide a variety of handling patterns: SHOW, DESCRIBE, and SHOW PARAMETERS outputs, as well as using ALTERS after CREATE because of limitations in Snowflake. All of the tests succeeded.

## Environment

We simply ran `terraform apply` for the given configurations on an empty Terraform state. The tests end when the commands finish successfully. The testing was done on Apple M3 MacBook Pro. We used Terraform CLI v1.10.3 and provider v1.0.1. We used a local (default) backend for storing the state.

## Implementation factors

In this section, we describe additional context and limitations of Terraform and Snowflake. We already addressed some of the changes before v1.

### Terraform workflow

After executing the `terraform apply` command, Terraform reads the module configuration and state file. When the state is kept in a non-local [backend](https://developer.hashicorp.com/terraform/language/backend), the state must be downloaded to be parsed. This alone can cause some delay if the state file is big (see [State size](#state-size)). Then, after computing diffs between the expected and actual state, the provider of a given resource is called to perform planned operations. Next, the provider calls Snowflake to operate (CREATE, ALTER, etc) on the given objects individually.

### External limitations

The logic inside is usually quite simple, but the network latency (we expect \~300-400ms, but it depends on your location and deployment region) can still increase the operation delay. For some resources, like tasks, Snowflake does not provide a way to set a field in CREATE, so ALTER must be called immediately after creation to set given fields. We are tracking such issues internally.

### New output fields

During our road to v1, we decided to include the outputs of SHOW, DESCRIBE, and SHOW PARAMETERS of the reworked objects in the state (read more [here](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md)). This caused the state size to grow (see State size section). These fields are needed to handle changes to fields with default values in Snowflake. We plan to dive deeper into this with [plan modifiers](https://developer.hashicorp.com/terraform/plugin/framework/resources/plan-modification) in Terraform Plugin Framework, but this requires migration to Terraform Plugin Framework, which we are planning after GA.

### Additional requests

Some resources make additional requests to handle all of the object fields in Snowflake, such as getting [policy references](https://docs.snowflake.com/en/sql-reference/account-usage/policy_references) in views or parameters [in users](https://github.com/snowflakedb/terraform-provider-snowflake/blob/c4b1bebce4bc5a81031248592b34af5e80ca2fc1/MIGRATION_GUIDE.md#breaking-change-user-parameters-added-to-snowflake_user-resource). As mentioned above, these requests usually take \~300-400ms, but in large-scale deployments, they can substantially increase the execution time.

## System stability

The tests utilize the network extensively and the results may vary depending on your deployment and the machine it runs on. However, the benchmarks were run consecutively a few times and the results proved to be close enough to call those times stable.

## Execution time

The most requested issue is performance degradation in recent versions ([\#3118](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3118)). We have already prepared basic ideas and suggestions in that thread.

We performed these tests on the local backend (see more in the [Environment](#environment) section).

The results for the `terraform apply` (which results in creating the resources) and the `terraform plan` are presented in the tables below.

Execution time of `terraform apply`:

| Resource count | 1 | 10 | 100 | 1000 | 4000 |
| ----- | ----- | :---- | :---- | :---- | :---- |
| Task | 6s | 9s | 29s | 4m 41s | 38m 9s |
| Schema | 6s | 7s | 28s | 3m 58s | 19m 0s |
| Warehouse | 6s | 8s | 21s | 3m 6s | 33m 41s |

Execution time of `terraform plan`:

| Resource count | 1 | 10 | 100 | 1000 | 4000 |
| ----- | :---- | :---- | :---- | :---- | :---- |
| Task | 3s | 6s | 20s | 2m 24s | 9m 48s |
| Schema | 3s | 7s | 24s | 2m 35s | 9m 57s |
| Warehouse | 3s | 5s | 17s | 2m 5s | 8m 34s |

The execution time does not rise linearly. This may be caused by parallelized operations (10 by default) in Terraform, which are not fully leveraged for lower resource counts and the initial `terraform` binary start time. The execution time of `terraform apply` for 4000 resources is surprisingly high – 5-10x difference for 4x size. The execution time of the `terraform plan` is roughly linear.

The execution time can be greatly affected by the number of resources bigger than \~1000. This has already been discussed in a few Terraform threads:

* [https://discuss.hashicorp.com/t/remote-state-file-size-limit/46324](https://discuss.hashicorp.com/t/remote-state-file-size-limit/46324),
* [https://github.com/hashicorp/terraform/issues/18981](https://github.com/hashicorp/terraform/issues/18981),
* [https://discuss.hashicorp.com/t/seeing-very-bad-performance-when-for-each-3k-resources/52536/7](https://discuss.hashicorp.com/t/seeing-very-bad-performance-when-for-each-3k-resources/52536/7),
* [https://github.com/hashicorp/terraform/issues/26355](https://github.com/hashicorp/terraform/issues/26355),
* [https://www.reddit.com/r/Terraform/comments/st7ohf/comment/hx2o6uw/](https://www.reddit.com/r/Terraform/comments/st7ohf/comment/hx2o6uw/),
* [https://www.reddit.com/r/Terraform/comments/wv0xgd/how\_to\_maximize\_parallelism\_for\_large\_plans/](https://www.reddit.com/r/Terraform/comments/wv0xgd/how_to_maximize_parallelism_for_large_plans/),

A standard solution HashiCorp recommends is splitting the deployments into smaller ones. This causes faster operations on state backends, faster parsing, and fewer resources handled during the `terraform plan`.

**Conclusion**: The execution time of large deployments can be surprisingly high. Consider limiting the number of resources in your deployments.

## State size

We measured the state sizes of the selected resources. We considered allowing a conditional removal of the `parameters` output field from the state. This was one of our ideas to reduce the state size (see [\#3118](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3118#issuecomment-2402618666)). We verified that removing this output is achievable because we do not use this field anymore in handling logic of resource parameters. The results are presented in the table below.

State size with different parallelism values:

| Resource count | 100 | 1000 | 4000 |
| ----- | :---- | :---- | :---- |
| Task  | 2.7MB | 27.1MB | 108MB |
| Task without parameters | 0.5MB | 5.0MB | 19MB |
| Schema | 0.8MB | 8.2MB | 33MB |
| Schema without parameters | 0.2MB | 1.9MB | 7.5MB |
| Warehouse  | 0.4MB | 3.5MB | 14MB |
| Warehouse without parameters | 0.2MB | 2.1MB | 8.3MB |

The resource state grows linearly with resource count, which is expected behavior. For bigger deployments, it may easily exceed tens of megabytes. In these cases, we recommend splitting the deployments to reduce parsing time and network time to your remote backend.

Removing the `parameters` output field greatly reduces the state size, and this idea can be explored further. However, fields `show_output` and `describe_output` can not be removed because they are used to handle default values in Snowflake. This idea can be explored after migrating to the Terraform Plugin Framework because it provides more control over handling default values.

**Conclusion**: State size grows linearly with the number of resources. Consider limiting the number of resources in your deployments. Conditionally clearing values in the `parameters` output field based on a new attribute in the provider configuration could benefit provider performance (while keeping v1 compatibility).

## Skipping refresh before destroying resources

If you are confident your state file reflects the current state of the resources, you can use the `-refresh=false` flag to skip refreshing. We tested destroying 1000 schema resources with and without refresh (see Terraform [docs](https://developer.hashicorp.com/terraform/cloud-docs/run/modes-and-options#skipping-automatic-state-refresh)). The default settings (with refreshing before destroy) resulted in 2m 56s of execution time and without refresh, 1m 46s.

This flag can only be safely used to destroy resources because of the underlying `DROP IF EXISTS`. Even if the given state is outdated, this would not impact the `DROP` operation. Using this flag during `terraform apply` is not recommended because it may result in faulty behavior. Potential external changes on the object (the object can have changed some fields or be present in state and absent from Snowflake and vice versa) will be undetected.

**Conclusion**: Terraform `-refresh=false` flag can be used only to speed up destroying resources. Beware that the state before the operation can be outdated.

## Parallelism

After Terraform computes a resource graph, it traverses the graph concurrently to act on some resources (see more [here](https://developer.hashicorp.com/terraform/internals/graph#walking-the-graph)). The concurrency can be controlled by the  `-parallelism=N` flag for applying and planning. The concurrency is limited by default to 10 to avoid overwhelming the upstream.
Remember that, according to Terraform docs, “Setting \-parallelism is considered an advanced operation and should not be necessary for normal usage of Terraform.”

We tested a few configurations of schema resources with different values of parallelism. The results are presented in the table below.

Execution time with varying values of parallelism:

| Resource count | 100 | 1000 | 4000 |
| ----- | :---- | :---- | :---- |
| \-parallelism=10 (default) | 28s | 3m 58s | 19m 0s |
| \-parallelism=20 | 21s | 2m 19s | 13m 28s |
| \-parallelism=40 | 16s | 1m 24s | 12m 4s |
| \-parallelism=60 | 13s | 1m 8s | 11m 40s |

We can observe that in some scenarios, we have 2x-4x improvement. However, for the biggest resource count, the gain is around 60%. The improvement depends on the resource count itself, but still, it is significant. These numbers may vary depending on your deployment and the machine it runs on. You should perform an investigation on the exact value.

**Conclusion**: Terraform `-parallelism=N` flag can speed up processing the resources. You should check the gain on your concrete deployments. However, providing too much parallelism may cause connection overloading and increased execution time.

## Terraform CLI version

HashiCorp released two performance improvements between Terraform ~1.8 and ~1.13:
* [core: move change encoding outside of core](https://github.com/hashicorp/terraform/pull/35558) - This improvement in Terraform 1.10.0 ([release notes](https://github.com/hashicorp/terraform/releases/tag/v1.10.0)) enhances the performance of state processing, namely internal encoding and decoding of big graphs.
* [reduce instance change lookups in GetResource](https://github.com/hashicorp/terraform/pull/37154) - This improvement in Terraform 1.13.0 ([release notes](https://github.com/hashicorp/terraform/releases/tag/v1.13.0)) improves resource instance lookups during operations.

If the planning stage takes a long time, even before the provider sends a request to Snowflake, consider upgrading to at least version 1.13.0. These improvements may significantly benefit your performance, especially when working with large workloads.

## Suggestions for users

Based on the executed benchmarks, topics on Terraform forums, and GitHub issues, we recommend the following:

* Split deployments into smaller ones (max. 100s of objects).
* Use the latest versions of the terraform binary and the provider.
* Use `-refresh=false` flag for `terraform destroy`.
* Try using `-parallelism=N` flag with different values.

## Summary

Based on our customers’ feedback, we have decided to perform basic benchmarks of Terraform Provider. When we made decisions during rework before v1, we anticipated some performance challenges due to the numerous assumptions and limitations imposed by Snowflake and Terraform.

We have prepared recommendations for the users. We hope that these suggestions will help reduce the provider's execution time. We have also provided an overview of potential future work (conditionally removing `parameters` output values), which is not on our current roadmap.
