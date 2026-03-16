resource "snowflake_external_volume" "complete" {
  name = var.name
  dynamic "storage_location" {
    for_each = var.storage_location
    content {
      storage_location_name        = storage_location.value["storage_location_name"]
      storage_provider             = storage_location.value["storage_provider"]
      storage_base_url             = storage_location.value["storage_base_url"]
      storage_aws_role_arn         = try(storage_location.value["storage_aws_role_arn"], null)
      storage_aws_external_id      = try(storage_location.value["storage_aws_external_id"], null)
      storage_aws_access_point_arn = try(storage_location.value["storage_aws_access_point_arn"], null)
      use_privatelink_endpoint     = try(storage_location.value["use_privatelink_endpoint"], null)
      encryption_type              = try(storage_location.value["encryption_type"], null)
      encryption_kms_key_id        = try(storage_location.value["encryption_kms_key_id"], null)
      azure_tenant_id              = try(storage_location.value["azure_tenant_id"], null)
      storage_endpoint             = try(storage_location.value["storage_endpoint"], null)
      storage_aws_key_id           = try(storage_location.value["storage_aws_key_id"], null)
      storage_aws_secret_key       = try(storage_location.value["storage_aws_secret_key"], null)
    }
  }
}
