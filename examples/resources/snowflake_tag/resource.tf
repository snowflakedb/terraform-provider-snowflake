# basic resource
resource "snowflake_tag" "tag" {
  name     = "tag"
  database = "database"
  schema   = "schema"
}

# complete resource
resource "snowflake_tag" "tag" {
  name                   = "tag"
  database               = "database"
  schema                 = "schema"
  comment                = "comment"
  ordered_allowed_values = ["finance", "engineering", ""]
  masking_policies       = [snowflake_masking_policy.example.fully_qualified_name]
}

# resource with propagation and conflict resolution
resource "snowflake_tag" "tag" {
  name                   = "tag"
  database               = "database"
  schema                 = "schema"
  ordered_allowed_values = ["high", "medium", "low"]
  propagate              = "ON_DEPENDENCY"

  on_conflict {
    allowed_values_sequence = true
  }
}
