package sdk

// SecondaryRoleOption is based on https://docs.snowflake.com/en/sql-reference/sql/use-secondary-roles.
type SecondaryRoleOption string

const (
	SecondaryRolesAll  SecondaryRoleOption = "ALL"
	SecondaryRolesNone SecondaryRoleOption = "NONE"
)
