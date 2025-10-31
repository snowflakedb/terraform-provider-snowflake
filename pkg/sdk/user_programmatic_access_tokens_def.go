package sdk

import (
	"fmt"
	"slices"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

type ProgrammaticAccessTokenStatus string

const (
	ProgrammaticAccessTokenStatusActive   ProgrammaticAccessTokenStatus = "ACTIVE"
	ProgrammaticAccessTokenStatusExpired  ProgrammaticAccessTokenStatus = "EXPIRED"
	ProgrammaticAccessTokenStatusDisabled ProgrammaticAccessTokenStatus = "DISABLED"
)

var allProgrammaticAccessTokenStatuses = []ProgrammaticAccessTokenStatus{
	ProgrammaticAccessTokenStatusActive,
	ProgrammaticAccessTokenStatusExpired,
	ProgrammaticAccessTokenStatusDisabled,
}

func toProgrammaticAccessTokenStatus(s string) (ProgrammaticAccessTokenStatus, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allProgrammaticAccessTokenStatuses, ProgrammaticAccessTokenStatus(s)) {
		return "", fmt.Errorf("invalid programmatic access token status: %s", s)
	}
	return ProgrammaticAccessTokenStatus(s), nil
}

var programmaticAccessTokenResultDBRowDef = g.DbStruct("programmaticAccessTokenResultDBRow").
	Text("name").
	Text("user_name").
	Text("role_restriction").
	Time("expires_at").
	Text("status").
	OptionalText("comment").
	Time("created_on").
	Text("created_by").
	OptionalNumber("mins_to_bypass_network_policy_requirement").
	OptionalText("rotated_to")

var programmaticAccessTokenDef = g.PlainStruct("ProgrammaticAccessToken").
	Text("Name").
	Field("UserName", "AccountObjectIdentifier").
	Field("RoleRestriction", "*AccountObjectIdentifier").
	Time("ExpiresAt").
	Field("Status", "ProgrammaticAccessTokenStatus").
	OptionalText("Comment").
	Time("CreatedOn").
	Text("CreatedBy").
	OptionalNumber("MinsToBypassNetworkPolicyRequirement").
	OptionalText("RotatedTo")

var addProgrammaticAccessTokenResultDBRowDef = g.DbStruct("addProgrammaticAccessTokenResultDBRow").
	Text("token_name").
	Text("token_secret")

var addProgrammaticAccessTokenResultDef = g.PlainStruct("AddProgrammaticAccessTokenResult").
	Text("TokenName").
	Text("TokenSecret")

var rotateProgrammaticAccessTokenResultDBRowDef = g.DbStruct("rotateProgrammaticAccessTokenResultDBRow").
	Text("token_name").
	Text("token_secret").
	Text("rotated_token_name")

var rotateProgrammaticAccessTokenResultDef = g.PlainStruct("RotateProgrammaticAccessTokenResult").
	Text("TokenName").
	Text("TokenSecret").
	Text("RotatedTokenName")

var UserProgrammaticAccessTokensDef = g.NewInterface(
	"UserProgrammaticAccessTokens",
	"UserProgrammaticAccessToken",
	// This works on an assumption that every object has an identifier. PATs do not have identifiers, and they cannot be referenced like "USER"."PAT", but their name part behaves like an identifier.
	// This means that we can use double quotes, the name must be non-empty and no longer than 255 characters.
	// We use AccountObjectIdentifier as a kind of identifier for convenience.
	// TODO(SNOW-2183032) Handle objects that do not have identifiers.
	g.KindOfT[AccountObjectIdentifier](),
).CustomShowOperation(
	"Add",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-programmatic-access-token",
	addProgrammaticAccessTokenResultDBRowDef,
	addProgrammaticAccessTokenResultDef,
	g.NewQueryStruct("AddUserProgrammaticAccessToken").
		Alter().
		SQL("USER").
		IfExists().
		Identifier("UserName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		SQL("ADD PROGRAMMATIC ACCESS TOKEN").
		Name().
		OptionalIdentifier("RoleRestriction", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("ROLE_RESTRICTION").Equals()).
		OptionalNumberAssignment("DAYS_TO_EXPIRY", g.ParameterOptions()).
		OptionalNumberAssignment("MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT", g.ParameterOptions()).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "UserName").
		WithValidation(g.ValidIdentifierIfSet, "RoleRestriction"),
).CustomOperation(
	"Modify",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-user-modify-programmatic-access-token",
	g.NewQueryStruct("ModifyUserProgrammaticAccessToken").
		Alter().
		SQL("USER").
		IfExists().
		Identifier("UserName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		SQL("MODIFY PROGRAMMATIC ACCESS TOKEN").
		Name().
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("ModifyProgrammaticAccessTokenSet").
				OptionalBooleanAssignment("DISABLED", g.ParameterOptions()).
				OptionalNumberAssignment("MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT", g.ParameterOptions()).
				OptionalComment(),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("ModifyProgrammaticAccessTokenUnset").
				OptionalSQL("DISABLED").
				OptionalSQL("MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT").
				OptionalSQL("COMMENT"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		OptionalIdentifier("RenameTo", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO").NoEquals()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "UserName").
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "RenameTo"),
).CustomShowOperation(
	"Rotate",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/alter-user-rotate-programmatic-access-token",
	rotateProgrammaticAccessTokenResultDBRowDef,
	rotateProgrammaticAccessTokenResultDef,
	g.NewQueryStruct("RotateUserProgrammaticAccessToken").
		Alter().
		SQL("USER").
		IfExists().
		Identifier("UserName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		SQL("ROTATE PROGRAMMATIC ACCESS TOKEN").
		Name().
		OptionalNumberAssignment("EXPIRE_ROTATED_TOKEN_AFTER_HOURS", g.ParameterOptions()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "UserName"),
).CustomOperation(
	"Remove",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-user-remove-programmatic-access-token",
	g.NewQueryStruct("RemoveUserProgrammaticAccessToken").
		Alter().
		SQL("USER").
		IfExists().
		Identifier("UserName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		SQL("REMOVE PROGRAMMATIC ACCESS TOKEN").
		Name().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "UserName"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-user-programmatic-access-tokens",
	programmaticAccessTokenResultDBRowDef,
	programmaticAccessTokenDef,
	g.NewQueryStruct("ShowUserProgrammaticAccessTokens").
		Show().
		SQL("USER PROGRAMMATIC ACCESS TOKENS").
		OptionalIdentifier("UserName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("FOR USER")),
).ShowByIdOperationWithNoFiltering()
