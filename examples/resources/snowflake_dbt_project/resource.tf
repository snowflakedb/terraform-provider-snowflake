# basic resource
resource "snowflake_dbt_project" "basic" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "BASIC_DBT_PROJECT"
}

# complete resource with all optional parameters
resource "snowflake_dbt_project" "complete" {
  database        = "DATABASE"
  schema          = "SCHEMA"
  name            = "COMPLETE_DBT_PROJECT"
  from            = "@my_stage/dbt_project"
  default_args    = "--target prod"
  default_version = "LAST"
  comment         = "An example DBT project for data transformations"
}

# resource with git repository source
resource "snowflake_dbt_project" "from_git" {
  database        = "DATABASE"
  schema          = "SCHEMA"
  name            = "GIT_DBT_PROJECT"
  from            = "@git_stage/dbt_project"
  default_args    = "--target production --vars '{\"env\": \"prod\"}'"
  default_version = "VERSION$1"
  comment         = "DBT project sourced from Git repository"
}

# resource with specific version
resource "snowflake_dbt_project" "versioned" {
  database        = "DATABASE"
  schema          = "SCHEMA"
  name            = "VERSIONED_DBT_PROJECT"
  from            = "@internal_stage/dbt_project"
  default_version = "FIRST"
  comment         = "DBT project using first available version"
}
