package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var SessionPoliciesDef = g.NewInterface(
	"SessionPolicies",
	"SessionPolicy",
	g.KindOfT[SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-session-policy",
		g.NewQueryStruct("CreateSessionPolicy").
			Create().
			OrReplace().
			SQL("SESSION POLICY").
			IfNotExists().
			Name().
			OptionalNumberAssignment("SESSION_IDLE_TIMEOUT_MINS", g.ParameterOptions().NoQuotes()).
			OptionalNumberAssignment("SESSION_UI_IDLE_TIMEOUT_MINS", g.ParameterOptions().NoQuotes()).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-session-policy",
		g.NewQueryStruct("AlterSessionPolicy").
			Alter().
			SQL("SESSION POLICY").
			IfExists().
			Name().
			OptionalIdentifier("RenameTo", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("SessionPolicySet").
					OptionalNumberAssignment("SESSION_IDLE_TIMEOUT_MINS", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("SESSION_UI_IDLE_TIMEOUT_MINS", g.ParameterOptions().NoQuotes()).
					OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
					WithValidation(g.AtLeastOneValueSet, "SessionIdleTimeoutMins", "SessionUiIdleTimeoutMins", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("SessionPolicyUnset").
					OptionalSQL("SESSION_IDLE_TIMEOUT_MINS").
					OptionalSQL("SESSION_UI_IDLE_TIMEOUT_MINS").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "SessionIdleTimeoutMins", "SessionUiIdleTimeoutMins", "Comment"),
				g.KeywordOptions().SQL("UNSET"),
			).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set", "SetTags", "UnsetTags", "Unset"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-session-policy",
		g.NewQueryStruct("DropSessionPolicy").
			Drop().
			SQL("SESSION POLICY").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-session-policies",
		g.DbStruct("showSessionPolicyDBRow").
			Field("created_on", "string").
			Field("name", "string").
			Field("database_name", "string").
			Field("schema_name", "string").
			Field("kind", "string").
			Field("owner", "string").
			Field("comment", "string").
			Field("options", "string").
			Field("owner_role_type", "string"),
		g.PlainStruct("SessionPolicy").
			Field("CreatedOn", "string").
			Field("Name", "string").
			Field("DatabaseName", "string").
			Field("SchemaName", "string").
			Field("Kind", "string").
			Field("Owner", "string").
			Field("Comment", "string").
			Field("Options", "string").
			Field("OwnerRoleType", "string"),
		g.NewQueryStruct("ShowSessionPolicies").
			Show().
			SQL("SESSION POLICIES"),
	).
	ShowByIdOperationWithNoFiltering().
	DescribeOperation(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-session-policy",
		g.DbStruct("describeSessionPolicyDBRow").
			Field("created_on", "string").
			Field("name", "string").
			Field("session_idle_timeout_mins", "int").
			Field("session_ui_idle_timeout_mins", "int").
			Field("comment", "string"),
		g.PlainStruct("SessionPolicyDescription").
			Field("CreatedOn", "string").
			Field("Name", "string").
			Field("SessionIdleTimeoutMins", "int").
			Field("SessionUIIdleTimeoutMins", "int").
			Field("Comment", "string"),
		g.NewQueryStruct("DescribeSessionPolicy").
			Describe().
			SQL("SESSION POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
