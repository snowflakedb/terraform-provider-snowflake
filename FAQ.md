# FAQ

* [What are the current/future plans for the provider?](#what-are-the-currentfuture-plans-for-the-provider)
* [When will the Snowflake feature X be available in the provider?](#when-will-the-snowflake-feature-x-be-available-in-the-provider)
* [When will my bug report be fixed/released?](#when-will-my-bug-report-be-fixedreleased)
* [How to migrate from version X to Y?](#how-to-migrate-from-version-x-to-y)
* [How can I contribute?](#how-can-i-contribute)
* [How can I debug the issue myself?](#how-can-i-debug-the-issue-myself)
* [How can I import already existing Snowflake infrastructure into Terraform?](#how-can-i-import-already-existing-snowflake-infrastructure-into-terraform)
* [What identifiers are valid inside the provider and how to reference one resource inside the other one?](#what-identifiers-are-valid-inside-the-provider-and-how-to-reference-one-resource-inside-the-other-one)
* [Is this provider compatible with OpenTofu?](#is-this-provider-compatible-with-opentofu)

### What are the current/future plans for the provider?
Our current plans are documented in the publicly available [roadmap](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md) that you can find in our repository.
We will be updating it to keep you posted on what’s coming for the provider.

### When will the Snowflake feature X be available in the provider?
It depends on the status of the feature. Snowflake marks features as follows:
- Private Preview (PrPr)
- Public Preview (PuPr)
- Generally Available (GA)

Currently, our main focus is on making the provider stable with the most stable GA features,
but please take a closer look at our recently updated [roadmap](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md)
which describes our priorities for the next quarters.

The provider uses SQL under the hood. When requesting a new feature,
make sure all the necessary SQL commands representing CRUD (CREATE/READ/UPDATE/DELETE) operations are available in Snowflake.
If they are not, you can create a feature request (reach out to your account manager) for Snowflake to add the missing functionality.

### When will my bug report be fixed/released?
Our team is checking daily incoming GitHub issues. The resolution depends on the complexity and the topic of a given issue, but the general rules are:
- If the issue is easy enough, we tend to answer it immediately and provide fix depending on the issue and our current workload.
- If the issue needs more insight, we tend to reproduce the issue usually in a matter of days and answer/fix it right away (also very dependent on our current workload).
- If the issue is a part of the incoming topic on the [roadmap](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md), we postpone it to resolve it with the related tasks.

The releases usually happen once every two-three weeks, mostly on Wednesdays or Thursdays.

### How to migrate from version X to Y?
As noted at the top of our [README](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#snowflake-terraform-provider),
the project is still experimental and breaking change may occur. We try to minimize such changes, but with some of the changes required for version 1.0.0, it’s inevitable.
Because of that, whenever we introduce any breaking change, we add it to the [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md).
It’s a document containing every breaking change (starting from around v0.73.0) with additional hints on how to migrate resources between the versions.

### How can I contribute?
If you would like to contribute to the project, please follow our [contribution guidelines](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/CONTRIBUTING.md).

### How can I debug the issue myself?
The provider is simply an abstraction issuing SQL commands through the Go Snowflake driver, so most of the errors will be connected to incorrectly built or executed SQL statements.
To see what SQLs are being run you have to set more verbose logging check the [section below](#how-can-i-turn-on-logs).
To confirm the correctness of the SQLs, refer to the [official Snowflake documentation](https://docs.snowflake.com/).
If the SQLs seem correct, try to run them in the [Snowsight](https://docs.snowflake.com/en/user-guide/ui-snowsight) to confirm it's not a Snowflake issue.

### How can I turn on logs?
The provider offers two main types of logging:
- Terraform execution (check [Terraform Debugging Documentation](https://www.terraform.io/internals/debugging)) - you can set it through the `TF_LOG` environment variable, e.g.: `TF_LOG=DEBUG`; it will make output of the Terraform execution more verbose.
- Snowflake communication (using the logs from the underlying [Go Snowflake driver](https://github.com/snowflakedb/gosnowflake)) - you can set it directly in the provider config ([`driver_tracing`](https://registry.terraform.io/providers/snowflakedb/snowflake/1.0.3/docs#driver_tracing-3) attribute), by `SNOWFLAKE_DRIVER_TRACING` environmental variable (e.g. `SNOWFLAKE_DRIVER_TRACING=info`), or by `drivertracing` field in the TOML file. To see the communication with Snowflake (including the SQL commands run) we recommend setting it to `info`.

As driver logs may seem cluttered, to locate the SQL commands run, search for:
- (preferred) `--terraform_provider_usage_tracking`
- `msg="Query:`
- `msg="Exec:`

### How can I import already existing Snowflake infrastructure into Terraform?
Please refer to [this document](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/docs/guides/resource_migration.md#3-three-options-from-here)
as it describes different approaches of importing the existing Snowflake infrastructure into Terraform as configuration.
One thing worth noting is that some approaches can be automated by scripts interacting with Snowflake and generating needed configuration blocks,
which is highly recommended for large-scale migrations.

### What identifiers are valid inside the provider and how to reference one resource inside the other one?
Please refer to [this document](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/docs/guides/identifiers_rework_design_decisions.md)
- For the recommended identifier format, take a look at the ["Known limitations and identifier recommendations"](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/docs/guides/identifiers_rework_design_decisions.md#known-limitations-and-identifier-recommendations) section.
- For a new way of referencing object identifiers in resources, take a look at the ["New computed fully qualified name field in resources" ](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/docs/guides/identifiers_rework_design_decisions.md#new-computed-fully-qualified-name-field-in-resources) section.

### Is this provider compatible with OpenTofu?
OpenTofu is not currently supported. While it's mostly compatible with Terraform,
OpenTofu's [reactive implementation of new Terraform features based on community demand](https://opentofu.org/faq/#opentofu-compatibility) may decrease compatibility over time. 
We plan to research OpenTofu support in the future, but there's no timeline yet (once planned, it will appear in the [roadmap](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/ROADMAP.md)).
For now, you must research and assess the risk of provider incompatibility.
