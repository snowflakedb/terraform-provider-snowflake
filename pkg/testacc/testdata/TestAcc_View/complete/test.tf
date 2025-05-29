resource "snowflake_view" "test" {
  name            = var.name
  comment         = var.comment
  database        = var.database
  schema          = var.schema
  is_secure       = var.is_secure
  copy_grants     = var.copy_grants
  change_tracking = var.change_tracking
  is_temporary    = var.is_temporary
  column {
    column_name = var.column1_name
    comment     = var.column1_comment
  }
  column {
    column_name = var.column2_name
    projection_policy {
      policy_name = var.column2_projection_policy
    }
    masking_policy {
      policy_name = var.column2_masking_policy
      using       = var.column2_masking_policy_using
    }
  }
  data_metric_function {
    function_name   = var.data_metric_function
    on              = var.data_metric_function_on
    schedule_status = "STARTED"
  }
  data_metric_schedule {
    using_cron = var.data_metric_schedule_using_cron
  }
  row_access_policy {
    policy_name = var.row_access_policy
    on          = var.row_access_policy_on
  }
  aggregation_policy {
    policy_name = var.aggregation_policy
    entity_key  = var.aggregation_policy_entity_key
  }
  statement = var.statement
}
