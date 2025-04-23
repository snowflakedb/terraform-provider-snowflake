## Support

For official support and urgent, production-impacting issues, please [contact Snowflake Support](https://community.snowflake.com/s/article/How-To-Submit-a-Support-Case-in-Snowflake-Lodge).

Please follow [creating issues guidelines](CREATING_ISSUES.md), [FAQ](FAQ.md), and [known issues](KNOWN_ISSUES.md) before submitting an issue on GitHub or directly to Snowflake Support.

# Snowflake Terraform Provider

> ⚠️ **Please note**: If you believe you have found a security issue, _please responsibly disclose_ by contacting us at [triage-terraformprovider-dl@snowflake.com](mailto:triage-terraformprovider-dl@snowflake.com).

> ⚠️ **Disclaimer**: The project is in GA version, but some features are in preview. Such resources and data sources are considered preview features in the provider, regardless of their state in Snowflake. We do not guarantee their stability. They will be reworked and marked as a stable feature in future releases. Breaking changes in these features are expected, even without bumping the major version. They are disabled by default. To use them, add the relevant feature name to `preview_features_enabled` field in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). The list of preview features is available below. Please always refer to the [Getting Help](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#getting-help) section in our Github repo to best determine how to get help for your questions.

----

![.github/workflows/ci.yml](https://github.com/snowflakedb/terraform-provider-snowflake/workflows/.github/workflows/ci.yml/badge.svg)

This is a terraform provider for managing [Snowflake](https://www.snowflake.com/) resources.

## Table of contents
<!-- TOC -->
* [Snowflake Terraform Provider](#snowflake-terraform-provider)
  * [Table of contents](#table-of-contents)
  * [Getting started](#getting-started)
  * [Migration guide](#migration-guide)
  * [Roadmap](#roadmap)
  * [Getting Help](#getting-help)
  * [Would you like to create an issue?](#would-you-like-to-create-an-issue)
  * [Additional SQL Client configuration](#additional-sql-client-configuration)
  * [Contributing](#contributing)
  * [Releases](#releases)
<!-- TOC -->

## Getting started

> The current namespace of the provider in the Terraform Registry is [snowflakedb](https://registry.terraform.io/namespaces/snowflakedb).
>
> If you're still using the `chanzuckerberg/snowflake` source, see [Upgrading from CZI Provider](./CZI_UPGRADE.md) to upgrade to the Snowflake-Labs namespace.
>
> If you're still using the `Snowflake-Labs/snowflake` source, see [Upgrading from Snowflake-Labs Provider](./SNOWFLAKEDB_MIGRATION.md) to upgrade to the snowflakedb namespace.

Install the Snowflake Terraform provider by adding a requirement block and a provider block to your Terraform codebase:
```hcl
terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = ">= 1.0.0"
    }
  }
}

provider "snowflake" {
  organization_name = "organization_name"
  account_name      = "account_name"
  user              = "johndoe"
  password          = var.password
  role              = "PUBLIC"
}

# Remember to provide the password securely.
variable "password" {
  type      = string
  sensitive = true
}
```

For more information on provider configuration see the [provider docs on the Terraform registry](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs).

Don't forget to run `terraform init` and you're ready to go! 🚀

Start browsing the [registry docs](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs) to find resources and data sources to use.

## Migration guide

Please check the [migration guide](./MIGRATION_GUIDE.md) when changing the version of the provider.

## Roadmap

Check [Roadmap](./ROADMAP.md).

## Getting Help

Some links that might help you:

- The [introductory tutorial](https://guides.snowflake.com/guide/terraforming_snowflake/#0) shows how to set up your Snowflake account for Terraform (service user, role, authentication, etc) and how to create your first resources in Terraform.
- The [docs on the Terraform registry](https://registry.terraform.io/providers/snowflakedb/snowflake/latest) are a complete reference of all resources and data sources supported and contain more advanced examples.
- The [discussions area](https://github.com/snowflakedb/terraform-provider-snowflake/discussions) of this repo, we use this forum to discuss new features and changes to the provider.
- **If you are an enterprise customer**, reach out to your account team. This helps us prioritize issues.
- The [issues section](https://github.com/snowflakedb/terraform-provider-snowflake/issues) might already have an issue addressing your question.

## Would you like to create an issue?
If you would like to create a GitHub issue, please read our [guide](./CREATING_ISSUES.md) first.
It contains useful links, FAQ, and commonly known issues with solutions that may already solve your case.

## Additional SQL Client configuration
The provider uses the underlying [gosnowflake](https://github.com/snowflakedb/gosnowflake) driver to send SQL commands to Snowflake.

By default, the underlying driver is set to error level logging. It can be changed by setting `driver_tracing` field in the configuration to one of (from most to least verbose):
- `trace`
- `debug`
- `info`
- `print`
- `warning`
- `error`
- `fatal`
- `panic`

Read more in [provider configuration docs](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema).

## Contributing

Check [Contributing](./CONTRIBUTING.md).

## Releases

Releases will be performed as needed, typically once every 2 weeks.

Releases are published to [the terraform registry](https://registry.terraform.io/providers/snowflakedb/snowflake/latest). Each change has its own release notes (e.g. https://github.com/snowflakedb/terraform-provider-snowflake/releases/tag/v0.89.0) and migration guide if needed (e.g. https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#v0880--v0890).
