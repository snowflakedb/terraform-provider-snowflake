package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var sessionPolicySecondaryRoles = g.NewQueryStruct("SessionPolicySecondaryRoles").
	OptionalSQLWithCustomFieldName("All", "('ALL')").
	OptionalSQLWithCustomFieldName("None", "()").
	List("Roles", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.ListOptions().Parentheses()).
	WithValidation(g.ExactlyOneValueSet, "All", "None", "Roles")

var sessionPoliciesDef = g.NewInterface(
	"SessionPolicies",
	"SessionPolicy",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
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
			OptionalQueryStructField("AllowedSecondaryRoles", sessionPolicySecondaryRoles, g.KeywordOptions().SQL("ALLOWED_SECONDARY_ROLES =")).
			OptionalQueryStructField("BlockedSecondaryRoles", sessionPolicySecondaryRoles, g.KeywordOptions().SQL("BLOCKED_SECONDARY_ROLES =")).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidateValue, "AllowedSecondaryRoles").
			WithValidation(g.ValidateValue, "BlockedSecondaryRoles").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-session-policy",
		g.NewQueryStruct("AlterSessionPolicy").
			Alter().
			SQL("SESSION POLICY").
			IfExists().
			Name().
			OptionalIdentifier("RenameTo", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("SessionPolicySet").
					OptionalNumberAssignment("SESSION_IDLE_TIMEOUT_MINS", g.ParameterOptions().NoQuotes()).
					OptionalNumberAssignment("SESSION_UI_IDLE_TIMEOUT_MINS", g.ParameterOptions().NoQuotes()).
					OptionalQueryStructField("AllowedSecondaryRoles", sessionPolicySecondaryRoles, g.KeywordOptions().SQL("ALLOWED_SECONDARY_ROLES =")).
					OptionalQueryStructField("BlockedSecondaryRoles", sessionPolicySecondaryRoles, g.KeywordOptions().SQL("BLOCKED_SECONDARY_ROLES =")).
					OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
					WithValidation(g.AtLeastOneValueSet, "SessionIdleTimeoutMins", "SessionUiIdleTimeoutMins", "AllowedSecondaryRoles", "BlockedSecondaryRoles", "Comment").
					WithValidation(g.ValidateValue, "AllowedSecondaryRoles").
					WithValidation(g.ValidateValue, "BlockedSecondaryRoles"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("SessionPolicyUnset").
					OptionalSQL("SESSION_IDLE_TIMEOUT_MINS").
					OptionalSQL("SESSION_UI_IDLE_TIMEOUT_MINS").
					OptionalSQL("ALLOWED_SECONDARY_ROLES").
					OptionalSQL("BLOCKED_SECONDARY_ROLES").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "SessionIdleTimeoutMins", "SessionUiIdleTimeoutMins", "AllowedSecondaryRoles", "BlockedSecondaryRoles", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
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
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-session-policies",
		g.StructPair("showSessionPolicyDBRow", "SessionPolicy").
			Text("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("kind").
			Text("owner").
			Text("comment").
			Text("owner_role_type").
			Text("options").
			WithConvertGeneration(),
		g.NewQueryStruct("ShowSessionPolicies").
			Show().
			SQL("SESSION POLICIES").
			OptionalLike().
			OptionalExtendedIn().
			OptionalOn().
			OptionalStartsWith().
			OptionalLimit(),
		g.ShowByIDLikeFiltering,
		g.ShowByIDExtendedInFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-session-policy",
		g.StructPair("describeSessionPolicyDBRow", "SessionPolicyProperty").
			Text("property").
			Text("value").
			Text("default").
			Text("description").
			WithConvertGeneration(),
		g.NewQueryStruct("DescribeSessionPolicy").
			Describe().
			SQL("SESSION POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
		g.PlainStruct("SessionPolicyDetails").
			SchemaObjectIdentifier().
			Text("Owner").
			Text("OwnerRoleType").
			Text("Comment").
			Number("SessionIdleTimeoutMins").
			Number("SessionUiIdleTimeoutMins").
			StringList("AllowedSecondaryRoles").
			StringList("BlockedSecondaryRoles"),
	).
	WithCustomInterfaceMethod(
		"DescribeDetails",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"*SessionPolicyDetails", "error",
	)
