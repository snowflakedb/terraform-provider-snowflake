package example

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/poc/generator"
)

//go:generate go run ../main.go

var (
	dbRoleRename = g.NewQueryStruct("DatabaseRoleRename").
		// Fields
		Identifier("Name", g.KindOfT[DatabaseObjectIdentifier](), g.IdentifierOptions().Required()).
		// Validations
		WithValidation(g.ValidIdentifier, "Name")

	nestedThirdLevel = g.NewQueryStruct("NestedThirdLevel").
		// Fields
		Identifier("Field", g.KindOfT[DatabaseObjectIdentifier](), g.IdentifierOptions().Required()).
		// Validations
		WithValidation(g.AtLeastOneValueSet, "Field")

	dbRoleSet = g.NewQueryStruct("DatabaseRoleSet").
		// Fields
		TextAssignment("COMMENT", g.ParameterOptions().SingleQuotes().Required()).
		OptionalQueryStructField("NestedThirdLevel", nestedThirdLevel, g.ListOptions().NoParentheses().SQL("NESTED"))

	dbRoleUnset = g.NewQueryStruct("DatabaseRoleUnset").
		// Fields
		OptionalSQL("COMMENT").
		// Validations
		WithValidation(g.AtLeastOneValueSet, "Comment")

	DatabaseRole = g.NewInterface(
		"DatabaseRoles",
		"DatabaseRole",
		"DatabaseObjectIdentifier",
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-database-role",
			g.NewQueryStruct("CreateDatabaseRole").
				// Fields
				Create().
				OrReplace().
				SQL("DATABASE ROLE").
				IfNotExists().
				Name().
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				// Validations
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-database-role",
			g.NewQueryStruct("AlterDatabaseRole").
				// Fields
				Alter().
				SQL("DATABASE ROLE").
				IfExists().
				Name().
				OptionalQueryStructField("Rename", dbRoleRename, g.ListOptions().NoParentheses().SQL("RENAME TO")).
				OptionalQueryStructField("Set", dbRoleSet, g.ListOptions().NoParentheses().SQL("SET")).
				OptionalQueryStructField("Unset", dbRoleUnset, g.ListOptions().NoParentheses().SQL("UNSET")).
				// Validations
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ExactlyOneValueSet, "Rename", "Set", "Unset"),
		)
)
