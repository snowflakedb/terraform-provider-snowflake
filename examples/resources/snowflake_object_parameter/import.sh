# The ID is <key>|<object_type>|<object_identifier>, where <object_identifier> is the fully qualified name of the object
# with every part double-quoted: "<name>" for an account object (e.g. a user or a warehouse),
# "<database_name>"."<name>" for a schema, "<database_name>"."<schema_name>"."<name>" for a schema object.
# Quote the whole ID so the shell does not read | as a pipe.
terraform import snowflake_object_parameter.example 'DATA_RETENTION_TIME_IN_DAYS|DATABASE|"<database_name>"'
terraform import snowflake_object_parameter.example 'DATA_RETENTION_TIME_IN_DAYS|TABLE|"<database_name>"."<schema_name>"."<table_name>"'
