## Local release testing

This document describes the steps to test a new release of the Snowflake Terraform Provider locally.
To run the GoReleaser locally, you have to have it installed locally. You can find the installation instructions [here](https://goreleaser.com/install/).

Once you have the GoReleaser installed, you can run the following command to release locally:

```bash
make release-local
```

The binary should be placed in the `./dist` directory. It should contain another directory with the provider name
and your current distribution (e.g., `terraform-provider-snowflake_darwin_arm64`). In there, you should find a binary
named `terraform-provider-snowflake_vX.Y.Z`. You can use this binary to test the provider locally. It can be only used
in specific cases as this binary should be generally used by Terraform CLI, because of that you'll get the following error:

```text
This binary is a plugin. These are not meant to be executed directly.
Please execute the program that consumes these plugins, which will
load any plugins automatically
```

To use the locally released binary, run the following command.
> **Note:** Replace the placeholders with the actual values:
> - `CURRENT_OS` - **lowercased** output from `uname -s` command (e.g., `darwin`)
> - `CURRENT_ARCH` - output from `arch` command (e.g., `arm64`, `amd64`)
> - `LATEST_GIT_TAG` - output from `git tag --sort=-version:refname | head -n 1` command (e.g., `v2.0.0`)

```bash
cp ./dist/terraform-provider-snowflake_$(CURRENT_OS)_$(CURRENT_ARCH)/terraform-provider-snowflake_$(LATEST_GIT_TAG) $(HOME)/.terraform.d/plugins/terraform-provider-snowflake
```

Next, edit your `~/.terraformrc` file to include the following:

```hcl
provider_installation {

  dev_overrides {
      "registry.terraform.io/snowflakedb/snowflake" = "<path_to_binary>" # $(HOME)/.terraform.d/plugins/
  }

  direct {}
}
```

Now, you are able to use locally released binary. Run the following configuration with and without the overrides to test the release:

```terraform
terraform {
  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = "=2.0.0"
    }
  }
}

provider "snowflake" {
}
```

Look for `[INFO] Setting provider version: X.X.X` in the output logs to verify that both (overridden and regular) versions are different.

Remember to remove the overrides from the `~/.terraformrc` file after you finish testing the release.
