---
page_title: "Importing Resources"
subcategory: ""
description: |-

---
# Importing resources

This guide describes how the provider behaves when you bring existing Snowflake objects under Terraform management.

## Import modes

When you import an existing Snowflake object, the provider has to decide what to put in the Terraform state for the attributes that your configuration omits. This decision is the import mode:

- **Enforcing mode** (default): the provider stores the values the object has in Snowflake. The plan after import shows a change for every attribute that your configuration omits but that has a value in Snowflake, whether it's the Snowflake default or not. Applying that plan resets these attributes to their Snowflake defaults.
- **Permissive mode** (opt-in): the provider accepts the values the object already has for the attributes that your configuration omits. The plan after import shows no change for them.

Both modes affect only the import itself. After the import, external changes to the object are detected as usual.

### Example

The examples below import the same compute pool, created outside Terraform with a non-default `AUTO_SUSPEND_SECS`:

```sql
CREATE COMPUTE POOL EXAMPLE_POOL
  MIN_NODES = 1
  MAX_NODES = 1
  INSTANCE_FAMILY = CPU_X64_XS
  AUTO_SUSPEND_SECS = 600;
```

The configuration sets only the required attributes:

```terraform
import {
  to = snowflake_compute_pool.example
  id = "\"EXAMPLE_POOL\""
}

resource "snowflake_compute_pool" "example" {
  name            = "EXAMPLE_POOL"
  min_nodes       = 1
  max_nodes       = 1
  instance_family = "CPU_X64_XS"
}
```

The same applies to the `terraform import` command:

```shell
terraform import snowflake_compute_pool.example '"EXAMPLE_POOL"'
```

#### Enforcing mode

`terraform plan` shows changes for `auto_resume` and `auto_suspend_secs`, which the configuration omits (`show_output` and `describe_output` collapsed):

```text
  # snowflake_compute_pool.example will be updated in-place
  # (imported from ""EXAMPLE_POOL"")
  ~ resource "snowflake_compute_pool" "example" {
      ~ auto_resume              = "true" -> "default"
      ~ auto_suspend_secs        = 600 -> -1
        backup_instance_families = []
        comment                  = null
      ~ describe_output          = [
          - { ... },
        ] -> (known after apply)
        for_application          = null
        fully_qualified_name     = "\"EXAMPLE_POOL\""
        id                       = "\"EXAMPLE_POOL\""
        instance_family          = "CPU_X64_XS"
        max_nodes                = 1
        min_nodes                = 1
        name                     = "EXAMPLE_POOL"
      ~ show_output              = [
          - { ... },
        ] -> (known after apply)
    }

Plan: 1 to import, 0 to add, 1 to change, 0 to destroy.
```

`show_output` and `describe_output` are recomputed on every update, so they're marked as `(known after apply)`.

`"default"` and `-1` are the values the provider uses for an omitted attribute: they mean "use the Snowflake default". Note that `auto_resume` shows a change even though `true` is already the Snowflake default.

Resolve the changes in one of two ways:

1. **Keep the current values:** add them to the configuration. The plan becomes empty and the object isn't changed.

   ```terraform
   resource "snowflake_compute_pool" "example" {
     name              = "EXAMPLE_POOL"
     min_nodes         = 1
     max_nodes         = 1
     instance_family   = "CPU_X64_XS"
     auto_resume       = "true"
     auto_suspend_secs = 600
   }
   ```

2. **Reset to the Snowflake defaults:** apply the plan. The provider unsets both attributes, so `AUTO_SUSPEND_SECS` changes from `600` to the Snowflake default `3600`, and `AUTO_RESUME` stays `true`.

#### Permissive mode

With [permissive mode enabled](#switching-to-permissive-mode), the same import with the same configuration shows no changes (`show_output` and `describe_output` collapsed):

```text
  # snowflake_compute_pool.example will be imported
    resource "snowflake_compute_pool" "example" {
        auto_resume              = "default"
        auto_suspend_secs        = -1
        backup_instance_families = []
        comment                  = null
        describe_output          = [
            { ... },
        ]
        for_application          = null
        fully_qualified_name     = "\"EXAMPLE_POOL\""
        id                       = "\"EXAMPLE_POOL\""
        instance_family          = "CPU_X64_XS"
        max_nodes                = 1
        min_nodes                = 1
        name                     = "EXAMPLE_POOL"
        show_output              = [
            { ... },
        ]
    }

Plan: 1 to import, 0 to add, 0 to change, 0 to destroy.
```

The compute pool keeps `AUTO_SUSPEND_SECS = 600`. The value Snowflake reports is still available in `show_output` and `describe_output` (`auto_suspend_secs = 600`).

Permissive mode only covers the attributes that your configuration omits. If the configuration sets an attribute to a value different from the one in Snowflake (e.g. `auto_suspend_secs = 300`), the plan shows that change in both modes.

### Choosing a mode

Use **enforcing mode** when the configuration should be the full description of the object. The plan after import shows every value that differs from what the configuration implies, so you decide about each one: keep it by adding it to the configuration, or reset it to the Snowflake default by applying the plan. Nothing set outside Terraform is left unnoticed.

Use **permissive mode** when an empty plan after import matters more, for example when you bring many existing objects under Terraform and don't want to review or write down every value set outside Terraform. The objects keep their current values, and your configuration describes only the attributes you set in it.

### Switching to permissive mode

Permissive mode is the [`IMPORT_PERMISSIVE_MODE`](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#import_permissive_mode) opt-in experiment. Add it to `experimental_features_enabled` in the provider configuration:

```terraform
provider "snowflake" {
  experimental_features_enabled = ["IMPORT_PERMISSIVE_MODE"]
}
```

The setting applies only to imports. Resources that are already in the state keep their values, both when you enable the setting and when you remove it later.

Permissive mode is being rolled out resource by resource. The [experiment's description](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#import_permissive_mode) lists the resources that support it.

[`IMPORT_BOOLEAN_DEFAULT`](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#import_boolean_default) is the predecessor of permissive mode, limited to the boolean fields that use the special `"default"` value in the stage, stream, and storage integration resources. These resources honor it until they support `IMPORT_PERMISSIVE_MODE`. See the [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#import_boolean_default-is-deprecated) for details.
