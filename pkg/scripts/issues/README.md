## Prerequisites
To use the scripts, generate access token here: https://github.com/settings/tokens?type=beta. If you use classic tokens, make sure it's authorized to access the organization.

## Generating the list of open issues
1. To get all open issues invoke the [first script](./gh/main.go) setting `SF_TF_SCRIPT_GH_ACCESS_TOKEN`:
```shell
  cd gh && SF_TF_SCRIPT_GH_ACCESS_TOKEN=<YOUR_PERSONAL_ACCESS_TOKEN> go run .
```
2. File `issues.json` should be generated in the `gh` directory. This is the input file for the second script.
3. To get process the issues invoke the [second script](./file/main.go):
```shell
  cd file && go run .
```
4. File `issues.csv` should be generated in the `file` directory. This is the CSV which summarizes all the issues we have.

## Closing old issues (regarding https://github.com/snowflakedb/terraform-provider-snowflake/discussions/2755)
1. First get all open issues by invoking:
```shell
  cd gh && SF_TF_SCRIPT_GH_ACCESS_TOKEN=<YOUR_PERSONAL_ACCESS_TOKEN> go run .
```
2. File `issues.json` should be generated in the `gh` directory. This is the input file for the second script. The next script is based also on `presnowflake_bucket.csv` that was created based on the GH issues filtering.
3. To filter only closeable issues invoke [this script](./filter-closeable-old-issues/main.go):
```shell
  cd filter-closeable-old-issues && go run .
```
4. Script will output files `issues_to_close.csv` and `issues_edited.csv`. There are two files documenting closing action on 30.04.2024 (`20240430 - issues_edited.csv` and `20240430 - issues_to_close.csv`). In `20240430 - notes.MD` there are notes regarding the questionable issues and the decisions taken.
5. To close the issues with the appropriate comment provide `issues_to_close.csv` in `close-with-comment` dir. Example `20240430 - issues_to_close.csv` is given. The run:
```shell
  cd close-with-comment && SF_TF_SCRIPT_GH_ACCESS_TOKEN=<YOUR_PERSONAL_ACCESS_TOKEN> go run .
```

## Creating new labels and assigning them to issues
1. Firstly, make sure all the needed labels exist in the repository, by running:
```shell
  cd create-labels && SF_TF_SCRIPT_GH_ACCESS_TOKEN=<YOUR_PERSONAL_ACCESS_TOKEN> go run .
```
2. Then, we have to get data about the existing issues with:
```shell
  cd gh && SF_TF_SCRIPT_GH_ACCESS_TOKEN=<YOUR_PERSONAL_ACCESS_TOKEN> go run .
```
3. Afterward, we need to process `issues.json` with:
```shell
  cd file && go run .
```
4. Next you have to analyze generated CSV and assign categories in the `Category` column and resource / data source in the `Object` column (the `GitHub issues buckets` Excel should be used here named as `GitHubIssuesBucket.csv`; Update already existing one). The csv document be of a certain format with the following columns (with headers): "A" column with issue ID (in the format of "#<issue_id>"), "B" column with the category that should be assigned to the issue (should be one of the supported categories: "OTHER", "RESOURCE", "DATA_SOURCE", "IMPORT", "SDK", "IDENTIFIERS", "PROVIDER_CONFIG", "GRANTS", and "DOCUMENTATION"), and the "C" column with the object type (should be in the format of the terraform resource, e.g. "snowflake_database"). Then, you'll be able to use this csv (put it next to the `main.go`) to assign labels to the correct issues.
```shell
  cd assign-labels && SF_TF_SCRIPT_GH_ACCESS_TOKEN=<YOUR_PERSONAL_ACCESS_TOKEN> go run .
```
