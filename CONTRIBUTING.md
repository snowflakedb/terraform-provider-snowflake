> **Before making a contribution**: make sure to [discuss a change with us](#discuss-a-change-with-us) first!

# Contributing

- [Setting up the development environment](#setting-up-the-development-environment)
- [Repository structure](#repository-structure)
- [Running the tests locally](#running-the-tests-locally)
- [Making a contribution](#making-a-contribution)
  - [Discuss a change with us!](#discuss-a-change-with-us)
  - [Follow the code conventions inside the repository](#follow-the-code-conventions-inside-the-repository)
  - [Test the change](#test-the-change)
  - [Describe the breaking changes](#describe-the-breaking-changes)
  - [Before submitting the PR](#before-submitting-the-pr)
  - [Naming and describing the PR](#naming-and-describing-the-pr)
  - [Requesting the review](#requesting-the-review)
  - [Adding support for a new snowflake object](#adding-support-for-a-new-snowflake-object)
    - [Introducing a new part of the SDK](#add-the-object-to-the-sdk)
- [Advanced Debugging](#advanced-debugging)
- [Extending the migration script](#extending-the-migration-script)

## Setting up the development environment

1. Install Golang environment (check instructions on the official page https://go.dev/doc/install depending on your OS).
2. Fork this repo and clone it. Base your changes on the [dev](https://github.com/snowflakedb/terraform-provider-snowflake/tree/dev) branch as it contains the latest unreleased changes.
3. Run `make dev-setup` in the main directory of the cloned repository.
4. You can clean up the dev setup by running `make dev-cleanup`.

## Repository structure

The notable technical files/directories inside the repository:

- `Makefile` - contains instructions to set up the developer's environment, run tests, etc.
- `pkg/provider` - definition of the provider
- `pkg/resources`, `pkg/datasources` - definitions and tests (consult section [Running the tests locally](#running-the-tests-locally) below) for resources and datasources
- `pkg/acceptance` - helpers for acceptance and integration tests
- `pkg/sdk` - definitions of the SDK objects (SDK is our client to Snowflake, using [gosnowflake driver](https://github.com/snowflakedb/gosnowflake) underneath)
- `pkg/sdk/testint` - integration tests for the SDK (consult section [Running the tests locally](#running-the-tests-locally) below)

**⚠️ Important ⚠️** We are in progress of cleaning up the repository structure, so beware of the changes in the packages/directories.

## Running the tests locally

Currently, we have three main types of tests:
- SDK unit tests (in directory `pkg/sdk`, files ending with `_test.go`)
- SDK integration tests (in directory `pkg/sdk/testint`, files ending with `_integration_test.go`)
- resource/datasource acceptance tests (in directories `pkg/resources` and `pkg/datasources`, files ending with `_acceptance_test.go`)

Both integration and acceptance tests require the connection to Snowflake (some of the tests require multiple accounts).

The preferred way of running particular tests locally is to create a config file `~/.snowflake/config`, with the following content.

```sh
[default]
account_name = "<your account name>"
organization_name = "<organization in which your account is located>"
user = "<your user>"
password = "<your password>"
role = "<your role>"
host="<host of your account, e.g. organisation-account_name.snowflakecomputing.com>"
```

To be able to run all the tests you additionally need the second profile `[secondary_test_account]`:

```sh
[secondary_test_account]
account_name = "<your account name>"
organization_name = "<organization in which your account is located>"
user = "<your user on the secondary account>"
password = "<your password on the secondary account>"
role = "<your role on the secondary account>"
host="<host of your account, e.g. organisation-account_name2.snowflakecomputing.com>"
```

**TIP**: check [how-can-i-get-my-organization-name](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/authentication_methods#how-can-i-get-my-organization-name) and [how-can-i-get-my-account-name](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/authentication_methods#how-can-i-get-my-account-name) sections in our guides if you have troubles setting the proper `organization_name` and `account_name`.

We are aware that not everyone has access to two different accounts, so the majority of tests can be run using just one account. The tests setup however, requires both profiles (`default` and `secondary_test_account`) to be present. You can use the same details for `secondary_test_account` as in the `default` one, if you don't plan to run tests requiring multiple accounts. The warning will be logged when setting up tests with just a single account.

There is also environment flag `TEST_SF_TF_SIMPLIFIED_INTEGRATION_TESTS_SETUP` available to set up only the default account for the integration tests. Careful, as tests requiring multiple accounts will fail when using this flag.

**⚠️ Important ⚠️** Some of the tests require the privileged role (like `ACCOUNTADMIN`). Otherwise, the managed objects may not be created. If you want to use lower role, you have to make sure it has all the necessary privileges added.

To run the tests we have the following commands:
- `make test-unit` run unit tests
- `make test-acceptance` run acceptance tests (without account-level ones)
- `make test-integration` run integration tests (without account-level ones)
- `make test-account-level-features` run both integration and acceptance tests verifying account-level features
- `make test-functional` run functional tests of the underlying terraform libraries (currently SDKv2)

The tests distinction between account-level and non-account-level tests is currently achieved by go build directive:
- `//go:build account_level_tests` for account-level tests;
- `//go:build non_account_level_tests` for non-account-level tests.
Make sure you specify the correct directive when adding new integration or acceptance test file.

You can run the particular tests from inside your chosen IDE but remember that you have to set `TF_ACC=1` environment variable to run any acceptance tests (the above commands set it for you). There are more environment variables set in the above Makefile rules, so familiarize with them before using them. It is also worth setting up more verbose logging (check [this section](FAQ.md#how-can-i-turn-on-logs) for more details).

## Making a contribution

### Discuss a change with us!
It's important to establish the scope of the change before the actual implementation. We want to avoid the situations in which the PR is rejected because it contradicts some other change we are introducing.

Remember to consult [our roadmap](ROADMAP.md), maybe we are already working on the issue!

It's best to approach us through the GitHub issues: either by commenting the already existing one or by creating a new one.

### Follow the code conventions inside the repository
We believe that code following the same conventions is easier to maintain and extend. When working on the given part of the provider try to follow the local solutions and not introduce too many new ideas.

### Test the change
Every introduced change should be tested. Depending on the type of the change it may require (any or mix of):
- adding/modifying existing unit tests (e.g. changing the behavior of validation in the SDK)
- adding/modifying existing integration tests (e.g. adding missing SDK invocations)
- adding/modifying existing acceptance tests (e.g. fixing the parameter on the resource level)

When writing acceptance tests, use the configuration and assertion generators instead of manually writing config strings. See [this guide](./pkg/acceptance/bettertestspoc/README.md) for more details.

It's best to discuss with us what checks we expect prior to making the change.

### Describe the breaking changes

If the change requires manual actions when bumping the provider version, they should be added to the [migration guide](MIGRATION_GUIDE.md).

### Before submitting the PR

The documentation for the provider is generated automatically. We follow the few formatting conventions that are automatically checked with every PR. They can fail and delay the resolution of your PR. To make it much less possible, run `make pre-push` before pushing your changes to GH. It will reformat your code (or suggest reformatting), generate all the missing docs, clean the dependencies, etc.

### Naming and describing the PR

We use [Conventional Commits](https://www.conventionalcommits.org/) for commit message formatting and PR titles. Please try to adhere to the standard.

Refer to the [regular expression](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/.github/workflows/title-lint.yml#L17) for PR title validation.

Implemented changes should be described thoroughly (we will prepare PR template for the known use cases soon):
- reference the issue that is addressed with the given change
- summary of changes
- summary of added tests
- (optional) what parts will be covered in the subsequent PRs

### Requesting the review

We check for the new PRs in our repository every day Monday-Friday. We usually need 1-2 days to leave the review comments. However, there are times when you can expect even more than a week response time. In such cases, please be patient, and ping us after a week if we do not post a reason for the delay ourselves. It's possible that we just missed it.

During our review we try to point out the unhandled special cases, missing tests, and deviations from the established conventions. Remember, review comment is like an invitation to dance: you don't have to agree but please provide the substantive reasons.

Please do not resolve our comments. We prefer to resolve ourselves after the comments are followed up by the contributor.

**⚠️ Important ⚠️** Tests and checks are not run automatically after your PR. We run them manually, when we are happy with the state of the change (even if some corrections are still necessary).

## Adding Support for a new Snowflake Object

This guide describes the end-to-end process to add support for a new Snowflake object in the Terraform provider. Work is typically split into multiple PRs:

| Step | Description | Example PR |
|------|-------------|------------|
| 1. SDK | Add SDK definitions and unit tests | [#4084](https://github.com/snowflakedb/terraform-provider-snowflake/pull/4084) |
| 2. Integration Tests | Add SDK integration tests | [#4123](https://github.com/snowflakedb/terraform-provider-snowflake/pull/4123) |
| 3. Resource | Add resource | [#4195](https://github.com/snowflakedb/terraform-provider-snowflake/pull/4195) |
| 4. Data Source | Add data source | [#4209](https://github.com/snowflakedb/terraform-provider-snowflake/pull/4209), [#4237](https://github.com/snowflakedb/terraform-provider-snowflake/pull/4237) |

### 1. Add the object to the SDK

Take a look at an example [SDK implementation for notebooks](https://github.com/snowflakedb/terraform-provider-snowflake/pull/4084).

- Follow the [SDK Generator guide](pkg/sdk/generator/README.md) to generate the object's SDK.

- Implement unit tests.

### 2. Add integration tests

Take a look at an example [Integration tests implementation for notebooks](https://github.com/snowflakedb/terraform-provider-snowflake/pull/4123).

Add integration tests under the SDK’s testint package to validate the SDK behavior against a live Snowflake connection.

- Follow the [Objects assertions guide](pkg/acceptance/bettertestspoc/README.md#adding-new-snowflake-object-assertions) to generate the necessary assertions.

### 3. Add resource

Take a look at an example [Resource implementation for notebooks](https://github.com/snowflakedb/terraform-provider-snowflake/pull/4195).

Implement the resource schema, read/create/update/delete, acceptance tests, and docs. Use the SDK as the source of truth and mirror its SHOW/DESC coverage and validations.

- Schema design
  - Prefer nested blocks for structured inputs. For example, “create from a stage” is modeled as a `from { stage = "<db>.<schema>.<stage>" path = "path/to/file" }` block rather than a flat string, to align with Snowflake semantics and improve validation.

  - Validate identifiers with the provider’s identifier validators (e.g., `IsValidIdentifier[...]`) and suppress quoting-only diffs for identifier fields (`suppressIdentifierQuoting`).

- Update semantics
  - If it's possible, implement rename in-place (`ALTER … RENAME TO …`) rather than ForceNew. Align with how recently refactored resources handle renames.

  - Detect external changes for derived outputs via SHOW/DESC triggers when possible. If a particular field cannot be detected externally (e.g., notebooks “from” location due to Snowflake limitations), document that limitation explicitly in the resource docs.

- Defaults and constraints surfaced in docs
  - Where Snowflake restricts identifier casing (e.g., only upper-case identifiers are valid for specific warehouse references), document it explicitly and add validators to prevent invalid inputs in plans.

- Documentation and migration guide
  - Add a Migration Guide entry under the correct version, grouping object support under a single H3 “(new feature) snowflake_” heading with H4 subsections for “Added resource” and “Added data source”.

  - When server capabilities are incomplete, document current limitations and ensure Update/Create sequences handle supported paths without requiring double-applies. Remember to use the model builder and assertions that you can automatically generate.

  - Add an example usage of the object. It should be auto-included in the generated documentation via the `.md.tmpl` file.

  - The resource's documentation should also include a `Preview feature` section. It could be a good idea to copy-paste and modify one of the existing files located in `templates/resources/*.md.tmpl`.

  - use `make docs` to generate documentation based on the `.md.tmpl` file (which is the file you should edit instead of `.md` file).

- Implement acceptance tests
  - Provide “basic” and “complete” cases; test rename, validations, and plan drift (ConfigPlanChecks). Avoid relying on “Safe” client wrappers for correctness checks; validate against the same paths real users hit.

- Follow the [Schemas guide](pkg/schemas/gen/README.md) to generate show schemas.

- Follow the [Resource assertions guide](pkg/acceptance/bettertestspoc/README.md#adding-new-resource-assertions) to generate the necessary assertions.

### 4. Add data source

Take a look at an example [Data source implementation for notebooks](https://github.com/snowflakedb/terraform-provider-snowflake/pull/4209) and its follow-up with extra tests [Extended test coverage for notebooks](https://github.com/snowflakedb/terraform-provider-snowflake/pull/4237)

While not strictly required to “support” the object, a data source improves discoverability and enables read-only use cases. For parity with other objects, we recommend adding one.

Example patterns validated by the data source:
- Filtering aligned to SHOW
  - Support `like`, `starts_with`, and `limit { rows, from }` to mirror SHOW filters; include `with_describe` to optionally call DESCRIBE for each item. Keep `with_describe` default-on but allow turning it off to reduce calls in large accounts.

- Output shape
  - Aggregate into a single `<object_name_plural>` collection with nested `show_output` (SHOW) and `describe_output` (DESCRIBE) blocks containing fields as strings/numbers

- Documentation and examples
  - Provide simple, filter, and pagination examples; include a note about default behavior of `with_describe`.

- Provider preview gate and migration guide
  - Add the “Added data source” H4 subsection under the same feature entry in the Migration Guide and link Snowflake’s SHOW docs where appropriate.

- Follow the [Data source config guide](pkg/acceptance/bettertestspoc/README.md#adding-new-datasource-config-model-builders) to generate config model.

## Advanced Debugging

If you want to build and test the provider locally (manually, not through acceptance tests), build the binary first using `make build-local` or install to the proper local directory by invoking `make install-tf` (to uninstall run `make uninstall-tf`).

Next, edit your `~/.terraformrc` file to include the following:

```hcl
provider_installation {

  dev_overrides {
      "registry.terraform.io/snowflakedb/snowflake" = "<path_to_binary>"
  }

  direct {}
}
```

To debug the provider with a debugger:

1. Launch the provider with the `-debug` command line argument in your debugger session. Once the provider starts, it will print instructions on setting the `TF_REATTACH_PROVIDERS` environment variable.

   ```sh
   Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable with the following:

   Command Prompt:	set "TF_REATTACH_PROVIDERS={"registry.terraform.io/snowflakedb/snowflake":{"Protocol":"grpc","ProtocolVersion":5,"Pid":35140,"Test":true,"Addr": {"Network":"tcp","String":"127.0.0.1:54706"}}}"
   PowerShell:	$env:TF_REATTACH_PROVIDERS='{"registry.terraform.io/snowflakedb/snowflake":{"Protocol":"grpc","ProtocolVersion":5,"Pid":35140,"Test":true,"Addr":{"Network":"tcp","String":"127.0.0.1:54706"}}}'
   ```

2. Open a terminal where you will execute Terraform and set the `TF_REATTACH_PROVIDERS` environment variable using the command from the first step.
3. Run Terraform as usual from this terminal. Any breakpoints you set will halt execution, and you can troubleshoot the provider from your debugger.

**Note**: The `TF_REATTACH_PROVIDERS` environment variable needs to be set every time you restart your debugger session as some values like the `Pid` or the TCP port will change with every execution.

For further instructions, please check the official [Terraform Plugin Development guide](https://www.terraform.io/plugin/debugging#starting-a-provider-in-debug-mode).

## Extending the migration script

If you wish to extend the [migration script](./pkg/scripts/migration_script/README.md)
please check the dedicated [contribution guide](./pkg/scripts/migration_script/CONTRIBUTING.md).
