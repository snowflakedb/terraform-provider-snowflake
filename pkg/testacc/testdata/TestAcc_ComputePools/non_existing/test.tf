data "snowflake_compute_pools" "test" {
  like = "non-existing-compute-pool"

  lifecycle {
    postcondition {
      condition     = length(self.compute_pools) > 0
      error_message = "there should be at least one compute pool"
    }
  }
}
