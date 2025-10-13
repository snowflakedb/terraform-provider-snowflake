# Upgrading from CZI Provider

As of (5/25/2022) the provider has been transferred from the Chan Zuckerberg Initiative (CZI) GitHub organization to Snowflake-Labs.
To upgrade from CZI, please run the following command:

```shell
terraform state replace-provider chanzuckerberg/snowflake Snowflake-Labs/snowflake
```

You should also update your lock file / Terraform provider version pinning. From the deprecated source:

```hcl
# deprecated source
terraform {
  required_providers {
    snowflake = {
      source  = "chanzuckerberg/snowflake"
      version = "0.36.0"
    }
  }
}
```

To new source:

```hcl
# new source
terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "0.36.0"
    }
  }
}
```

If you are not pinning your provider versions, you may find it useful to forcefully upgrade providers using the command:

```sh
terraform init -upgrade
```

To verify if all namespaces in the current working directory have been replaced, you can run the [terraform providers](https://developer.hashicorp.com/terraform/cli/commands/providers) command. Then, filter the outputs and ensure there are no `Snowflake-Labs` entries.

>**Note**:  0.34 is the first version published after the transfer. When the provider was transferred over not all of the older releases were transferred. Only versions 0.28 and newer were transferred. If you are using a version older than 0.28, it is highly recommended to upgrade to a newer version.
