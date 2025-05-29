resource "snowflake_schema" "test" {
  database = var.database
  name     = var.schema_name
}

resource "snowflake_grant_privileges_to_database_role" "test" {
  depends_on         = [snowflake_schema.test]
  database_role_name = "\"${var.database}\".\"${var.name}\""
  privileges         = var.privileges
  on_schema {
    schema_name = "${var.database}.${var.schema_name}"
  }
}
