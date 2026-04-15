//go:build sdk_generation

package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var passwordPoliciesDef = g.NewInterface(
	"PasswordPolicies",
	"PasswordPolicy",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-password-policy",
		g.NewQueryStruct("CreatePasswordPolicy").
			Create().
			OrReplace().
			SQL("PASSWORD POLICY").
			IfNotExists().
			Name().
			OptionalNumberAssignment("PASSWORD_MIN_LENGTH", g.ParameterOptions().NoQuotes()).
			OptionalNumberAssignment("PASSWORD_MAX_LENGTH", g.ParameterOptions().NoQuotes()).
			OptionalNumberAssignment("PASSWORD_MIN_UPPER_CASE_CHARS", g.ParameterOptions().NoQuotes()).
			OptionalNumberAssignment("PASSWORD_MIN_LOWER_CASE_CHARS", g.ParameterOptions().NoQuotes()).
			OptionalNumberAssignment("PASSWORD_MIN_NUMERIC_CHARS", g.ParameterOptions().NoQuotes()).
			OptionalNumberAssignment("PASSWORD_MIN_SPECIAL_CHARS", g.ParameterOptions().NoQuotes()).
			OptionalNumberAssignment("PASSWORD_MIN_AGE_DAYS", g.ParameterOptions().NoQuotes()).
			OptionalNumberAssignment("PASSWORD_MAX_AGE_DAYS", g.ParameterOptions().NoQuotes()).
			OptionalNumberAssignment("PASSWORD_MAX_RETRIES", g.ParameterOptions().NoQuotes()).
			OptionalNumberAssignment("PASSWORD_LOCKOUT_TIME_MINS", g.ParameterOptions().NoQuotes()).
			OptionalNumberAssignment("PASSWORD_HISTORY", g.ParameterOptions().NoQuotes()).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-password-policy",
		g.NewQueryStruct("AlterPasswordPolicy").
			Alter().
			SQL("PASSWORD POLICY").
			IfExists().
			Name().
			OptionalIdentifier("NewName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("PasswordPolicySet").
					OptionalNumberAssignment("PASSWORD_MIN_LENGTH", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("PASSWORD_MAX_LENGTH", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("PASSWORD_MIN_UPPER_CASE_CHARS", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("PASSWORD_MIN_LOWER_CASE_CHARS", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("PASSWORD_MIN_NUMERIC_CHARS", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("PASSWORD_MIN_SPECIAL_CHARS", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("PASSWORD_MIN_AGE_DAYS", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("PASSWORD_MAX_AGE_DAYS", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("PASSWORD_MAX_RETRIES", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("PASSWORD_LOCKOUT_TIME_MINS", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("PASSWORD_HISTORY", g.ParameterOptions().NoQuotes()).
					OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
					WithValidation(g.AtLeastOneValueSet, "PasswordMinLength", "PasswordMaxLength", "PasswordMinUpperCaseChars", "PasswordMinLowerCaseChars", "PasswordMinNumericChars", "PasswordMinSpecialChars", "PasswordMinAgeDays", "PasswordMaxAgeDays", "PasswordMaxRetries", "PasswordLockoutTimeMins", "PasswordHistory", "Comment"),
				g.ListOptions().NoParentheses().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("PasswordPolicyUnset").
					OptionalSQL("PASSWORD_MIN_LENGTH").
					OptionalSQL("PASSWORD_MAX_LENGTH").
					OptionalSQL("PASSWORD_MIN_UPPER_CASE_CHARS").
					OptionalSQL("PASSWORD_MIN_LOWER_CASE_CHARS").
					OptionalSQL("PASSWORD_MIN_NUMERIC_CHARS").
					OptionalSQL("PASSWORD_MIN_SPECIAL_CHARS").
					OptionalSQL("PASSWORD_MIN_AGE_DAYS").
					OptionalSQL("PASSWORD_MAX_AGE_DAYS").
					OptionalSQL("PASSWORD_MAX_RETRIES").
					OptionalSQL("PASSWORD_LOCKOUT_TIME_MINS").
					OptionalSQL("PASSWORD_HISTORY").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "PasswordMinLength", "PasswordMaxLength", "PasswordMinUpperCaseChars", "PasswordMinLowerCaseChars", "PasswordMinNumericChars", "PasswordMinSpecialChars", "PasswordMinAgeDays", "PasswordMaxAgeDays", "PasswordMaxRetries", "PasswordLockoutTimeMins", "PasswordHistory", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "NewName"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-password-policy",
		g.NewQueryStruct("DropPasswordPolicy").
			Drop().
			SQL("PASSWORD POLICY").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-password-policies",
		g.DbStruct("passwordPolicyDBRow").
			Field("created_on", "time.Time").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("kind").
			Text("owner").
			Text("comment").
			Text("owner_role_type").
			Text("options"),
		g.PlainStruct("PasswordPolicy").
			Field("CreatedOn", "time.Time").
			Text("Name").
			Text("DatabaseName").
			Text("SchemaName").
			Text("Kind").
			Text("Owner").
			Text("Comment").
			Text("OwnerRoleType").
			Text("Options"),
		g.NewQueryStruct("ShowPasswordPolicies").
			Show().
			SQL("PASSWORD POLICIES").
			OptionalLike().
			OptionalIn().
			OptionalLimit(),
	).
	ShowByIdOperationWithFiltering(g.ShowByIDLikeFiltering, g.ShowByIDInFiltering).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-password-policy",
		g.DbStruct("describePasswordPolicyDBRow").
			Text("property").
			Text("value").
			Text("default").
			Text("description"),
		g.PlainStruct("PasswordPolicyProperty").
			Text("Property").
			Text("Value").
			Text("Default").
			Text("Description"),
		g.NewQueryStruct("DescribePasswordPolicy").
			Describe().
			SQL("PASSWORD POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
		g.PlainStruct("PasswordPolicyDetails").
			Text("Name").
			Text("Owner").
			Text("Comment").
			Number("PasswordMinLength").
			Number("PasswordMaxLength").
			Number("PasswordMinUpperCaseChars").
			Number("PasswordMinLowerCaseChars").
			Number("PasswordMinNumericChars").
			Number("PasswordMinSpecialChars").
			Number("PasswordMinAgeDays").
			Number("PasswordMaxAgeDays").
			Number("PasswordMaxRetries").
			Number("PasswordLockoutTimeMins").
			Number("PasswordHistory"),
	)
