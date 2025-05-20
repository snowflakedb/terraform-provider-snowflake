resource "snowflake_external_volume" "complete" {
  name = var.name
  dynamic "storage_location" {
    for_each = var.storage_location
    content {
      storage_location_name = storage_location.value["storage_location_name"]
      storage_provider      = storage_location.value["storage_provider"]
      storage_base_url      = storage_location.value["storage_base_url"]
      storage_aws_role_arn  = try(storage_location.value["storage_aws_role_arn"], null)
      encryption_type       = try(storage_location.value["encryption_type"], null)
      encryption_kms_key_id = try(storage_location.value["encryption_kms_key_id"], null)
      azure_tenant_id       = try(storage_location.value["azure_tenant_id"], null)
    }
  }
}
