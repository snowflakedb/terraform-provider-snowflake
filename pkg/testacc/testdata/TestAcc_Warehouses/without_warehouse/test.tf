data "snowflake_warehouses" "test" {
  like = "non-existing-warehouse"

  lifecycle {
    postcondition {
      condition     = length(self.warehouses) > 0
      error_message = "there should be at least one warehouse"
    }
  }
}
