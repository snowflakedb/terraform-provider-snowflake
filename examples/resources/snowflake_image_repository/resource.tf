# basic resource
resource "snowflake_image_repository" "basic" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "BASIC"
}

# complete resource
resource "snowflake_image_repository" "complete" {
  database = "DATABASE"
  schema   = "SCHEMA"
  name     = "BASIC"
  comment  = "An example image repository"
}
