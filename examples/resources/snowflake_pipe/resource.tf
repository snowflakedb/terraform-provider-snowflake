resource "snowflake_pipe" "pipe" {
  database = snowflake_database.database.name
  schema   = snowflake_schema.schema.name
  name     = "PIPE"

  comment = "A pipe."

  copy_statement = "copy into ${snowflake_table.table.fully_qualified_name} from @${snowflake_stage.stage.fully_qualified_name}"
  auto_ingest    = false

  aws_sns_topic_arn    = "..."
  notification_channel = "..."
}


# Recreate the pipe when any of the cloud parameters of the referenced stage change.
# Read more at https://docs.snowflake.com/en/user-guide/data-load-snowpipe-manage#changing-the-cloud-parameters-of-the-referenced-stage.
resource "snowflake_stage" "stage" {
  database = snowflake_database.database.name
  schema   = snowflake_schema.schema.name
  name     = "STAGE"

  url                 = "s3://com.example.bucket/prefix"
  storage_integration = snowflake_storage_integration.storage_integration.name
  encryption          = "TYPE = 'NONE'"
}

resource "snowflake_pipe" "pipe_with_stage_change_trigger" {
  database = snowflake_stage.stage.database
  schema   = snowflake_stage.stage.schema
  name     = "PIPE_WITH_STAGE_CHANGE_TRIGGER"

  copy_statement = "copy into ${snowflake_table.table.fully_qualified_name} from @${snowflake_stage.stage.fully_qualified_name}"

  # Use the replace_triggered_by meta-argument to recreate the pipe when any of the referenced attributes changes.
  # When the referenced attributes are added to or removed from the stage configuration, the pipe will also be recreated.
  # Read more at https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#replace_triggered_by.
  lifecycle {
    replace_triggered_by = [
      snowflake_stage.stage.url,
      snowflake_stage.stage.storage_integration,
      snowflake_stage.stage.encryption,
    ]
  }
}
