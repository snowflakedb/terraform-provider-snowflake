# Snowflake Terraform Provider

[![Unit tests](https://github.com/snowflakedb/terraform-provider-snowflake/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/snowflakedb/terraform-provider-snowflake/actions/workflows/unit-tests.yml)
[![Latest release](https://img.shields.io/github/v/release/snowflakedb/terraform-provider-snowflake)](https://github.com/snowflakedb/terraform-provider-snowflake/releases/latest)
[![Terraform Registry](https://img.shields.io/badge/registry-snowflakedb%2Fsnowflake-623CE4)](https://registry.terraform.io/providers/snowflakedb/snowflake/latest)
[![License](https://img.shields.io/github/license/snowflakedb/terraform-provider-snowflake)](LICENSE)

The Snowflake provider for [Terraform](https://www.terraform.io) manages [Snowflake](https://www.snowflake.com/) objects (databases, warehouses, users, roles, grants, and more) through Snowflake SQL.

Important documents: [Registry docs](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs)
| [Tutorial](https://guides.snowflake.com/guide/terraforming_snowflake/#0)
| [Authentication](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/authentication_methods)
| [FAQ](FAQ.md)
| [Known issues](KNOWN_ISSUES.md)
| [Migration guide](MIGRATION_GUIDE.md)
| [BCR migration guide](SNOWFLAKE_BCR_MIGRATION_GUIDE.md)
| [Roadmap](ROADMAP.md)
| [Changelog](CHANGELOG.md)
| [Contributing](CONTRIBUTING.md)

The provider is generally available. Some resources and data sources are still **preview**: they are disabled by default, may change without a major version bump, and are not officially supported. Enable them with `preview_features_enabled` in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema). See the [preview resources](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#currently-preview-resources) and [preview data sources](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#currently-preview-data-sources) lists. [`experimental_features_enabled`](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#experimental_features_enabled-1) changes behavior of existing resources; it is also a preview feature.

> ⚠️ **Please note**: If you believe you have found a security issue, _please responsibly disclose_ by contacting us at [triage-terraformprovider-dl@snowflake.com](mailto:triage-terraformprovider-dl@snowflake.com). Do not open a public GitHub issue.

## Support

For official support and urgent, production-impacting issues, please [contact Snowflake Support](https://community.snowflake.com/s/article/How-To-Submit-a-Support-Case-in-Snowflake-Lodge).

> ⚠️ **Keep in mind** that the official support starts with the [v2.0.0](https://registry.terraform.io/providers/snowflakedb/snowflake/2.0.0) for **stable resources and data sources only**. All previous versions, preview resources and data sources, and experimental behaviors are not officially supported.

Please follow [creating issues guidelines](CREATING_ISSUES.md), [FAQ](FAQ.md), and [known issues](KNOWN_ISSUES.md) before submitting an issue on GitHub or directly to Snowflake Support.

We welcome contributions. See [CONTRIBUTING.md](CONTRIBUTING.md) before opening a PR.

## Table of contents

* [Requirements](#requirements)
  * [Supported Architectures](#supported-architectures)
* [Getting started](#getting-started)
* [Getting Help](#getting-help)
  * [Which guide do I need?](#which-guide-do-i-need)
* [Additional SQL Client configuration](#additional-sql-client-configuration)
* [Contributing](#contributing)
* [Releases](#releases)
* [License](#license)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install) CLI **1.1.5** or later ([known issues](KNOWN_ISSUES.md#old-terraform-cli-version))
- Provider **v2.0.0** or later for official support of stable resources; prefer the latest `2.x` release
- A Snowflake account and a user/role that can create the objects you want to manage. The [introductory tutorial](https://guides.snowflake.com/guide/terraforming_snowflake/#0) shows how to set up a service user and role.
- [OpenTofu](https://opentofu.org/) is **not supported**. The provider is listed in the OpenTofu Registry, but compatibility is untested. See the [FAQ](FAQ.md#is-this-provider-compatible-with-opentofu).

### Supported Architectures

Official binaries follow what the [gosnowflake driver](https://github.com/snowflakedb/gosnowflake) supports and what [HashiCorp recommends for Terraform providers](https://developer.hashicorp.com/terraform/registry/providers/os-arch):

- Windows: amd64
- Linux: amd64 and arm64
- Darwin: amd64 and arm64

These binaries are also published but are **not** officially supported (fixes are not prioritized):

- Windows: arm64 and 386
- Linux: 386
- Darwin: 386
- Freebsd: any architecture

## Getting started

The Registry source is [`snowflakedb/snowflake`](https://registry.terraform.io/namespaces/snowflakedb).

The example below uses [key-pair (JWT) authentication](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/authentication_methods#jwt-authenticator-flow), which is only one of the authentication methods this provider supports, see the [Authentication Methods](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/authentication_methods) guide. Keep keys and other secrets out of version control (use `file()`, environment variables, or a `~/.snowflake/config` [profile](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#toml-file)).

```hcl
terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "~> 2.20"
    }
  }
}

provider "snowflake" {
  organization_name = "<organization_name>"
  account_name      = "<account_name>"
  user              = "<user_name>"
  authenticator     = "SNOWFLAKE_JWT"
  private_key       = file("~/.ssh/snowflake_key.p8")
  role              = "SYSADMIN"
}

resource "snowflake_database" "example" {
  name = "TF_DEMO"
}
```

```bash
terraform init   # download the provider
terraform plan   # preview changes
terraform apply  # create the resources
```

Next:

- [Introductory tutorial](https://guides.snowflake.com/guide/terraforming_snowflake/#0) — service user, role, and first resources
- [Registry docs](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs) — full resource and data source reference

## Getting Help
<a id="would-you-like-to-create-an-issue"></a>
<a id="reporting-issues"></a>
<a id="migration-guide"></a>
<a id="roadmap"></a>

Start with the matching guide. If you are still blocked, use the support path below.

### Which guide do I need?

| If you need to… | Read |
|---|---|
| Install and configure the provider | [Registry docs](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs) and [authentication methods](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/authentication_methods) |
| Upgrade the provider version (for example 2.19 → 2.20) | [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) |
| Absorb a Snowflake behavior-change (BCR) bundle | [SNOWFLAKE_BCR_MIGRATION_GUIDE.md](SNOWFLAKE_BCR_MIGRATION_GUIDE.md) |
| Change Registry source | [CZI_UPGRADE.md](CZI_UPGRADE.md) then [SNOWFLAKEDB_MIGRATION.md](SNOWFLAKEDB_MIGRATION.md) |
| Import existing Snowflake objects into Terraform | [Resource migration guide](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/resource_migration) |
| See current plans | [ROADMAP.md](ROADMAP.md) |

1. **Production outage / official support** — [Snowflake Support](https://community.snowflake.com/s/article/How-To-Submit-a-Support-Case-in-Snowflake-Lodge) (v2.0.0+ stable resources). Enterprise customers can also contact their account team.
2. **How do I…?** — [FAQ](FAQ.md), [known issues](KNOWN_ISSUES.md), [discussions](https://github.com/snowflakedb/terraform-provider-snowflake/discussions), and the [introductory tutorial](https://guides.snowflake.com/guide/terraforming_snowflake/#0).
3. **Is this already filed?** — Search [open and closed issues](https://github.com/snowflakedb/terraform-provider-snowflake/issues).
4. **Open an issue** — Read [CREATING_ISSUES.md](CREATING_ISSUES.md) first, then use the matching GitHub issue template.

## Additional SQL Client configuration

The provider talks to Snowflake through the [gosnowflake](https://github.com/snowflakedb/gosnowflake) driver. Driver log level is `driver_tracing` in the [provider configuration](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema) (or `SNOWFLAKE_DRIVER_TRACING`). To see SQL that Terraform runs, see [How can I turn on logs?](FAQ.md#how-can-i-turn-on-logs) in the FAQ.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for requirements, the PR workflow and others.

## Releases

Releases land as needed, typically every two weeks, on the [Terraform Registry](https://registry.terraform.io/providers/snowflakedb/snowflake/latest).

- [GitHub Releases](https://github.com/snowflakedb/terraform-provider-snowflake/releases) and [CHANGELOG.md](CHANGELOG.md)
- Breaking changes: [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)
- Snowflake BCR impact: [SNOWFLAKE_BCR_MIGRATION_GUIDE.md](SNOWFLAKE_BCR_MIGRATION_GUIDE.md)

## License

[MIT](LICENSE)
