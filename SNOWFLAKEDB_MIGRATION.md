# Upgrading from Snowflake-Labs to snowflakedb Terraform Registry namespace

As of (03-04-2025) the provider has been transferred from the Snowflake-Labs GitHub organization to snowflakedb. The Terraform Registry namespace has been mirrored to snowflakedb. The new versions will be published only in the snowflakedb namespace. The old namespace will be deprecated and removed in the future - we will share more details soon.
All versions available in Snowflake-Labs (v0.28.6 - v1.0.5) are available in snowflakedb as well.
Before upgrading, please back up your state first. To upgrade the provider, please run the following command:

```shell
terraform state replace-provider Snowflake-Labs/snowflake snowflakedb/snowflake
```

You should also update your lock file / Terraform provider version pinning. From the deprecated configuration:

```hcl
# the old Snowflake-Labs namespace
terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "1.0.5"
    }
  }
}
```

To the new configuration:

```hcl
# the new snowflakedb namespace
terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "1.0.5"
    }
  }
}
```

If you are not pinning your provider versions, you may find it useful to forcefully upgrade providers using the command:

```shell
terraform init -upgrade
```

You can verify the required providers with their namespaces and versions in your deployments with the [terraform providers](https://developer.hashicorp.com/terraform/cli/commands/providers) command.

>**Note**: When the provider was transferred over not all of the older releases were transferred. Only versions 0.28 and newer were transferred (the ones from Snowflake-Labs). If you are using a version older than 0.28, it is highly recommended to upgrade to a newer version and then change to snowflakedb namespace.
