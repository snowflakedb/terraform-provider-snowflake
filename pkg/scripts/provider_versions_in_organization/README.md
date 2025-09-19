## Prerequisites
To use the scripts, generate access token (classic) here: https://github.com/settings/tokens. Select `repo` scope. Authorize it to access the organization.

## Generating the list of Snowflake Terraform Provider versions in repositories in the given GH organization
1. From the main directory of the project run (replacing `<token>` and `<GH organization>` with proper values:
```shell
  SF_TF_SCRIPT_GH_ACCESS_TOKEN=<token> go run ./pkg/scripts/provider_versions_in_organization/main.go <GH organization>
```
2. File `results.csv` should be generated in the main directory.
