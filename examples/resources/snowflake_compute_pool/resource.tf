# basic resource
resource "snowflake_compute_pool" "basic" {
  name            = "COMPUTE_POOL"
  min_nodes       = 1
  max_nodes       = 2
  instance_family = "CPU_X64_S"
}

# complete resource
resource "snowflake_compute_pool" "complete" {
  name                = "COMPUTE_POOL"
  for_application     = "APPLICATION_NAME"
  min_nodes           = 1
  max_nodes           = 2
  instance_family     = "CPU_X64_S"
  auto_resume         = "true"
  initially_suspended = "true"
  auto_suspend_secs   = 1200
  comment             = "A compute pool."
}
