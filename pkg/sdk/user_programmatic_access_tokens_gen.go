package sdk

import (
	"context"
	"database/sql"
	"time"
)

type UserProgrammaticAccessTokens interface {
	Add(ctx context.Context, request *AddUserProgrammaticAccessTokenRequest) ([]AddProgrammaticAccessTokenResult, error)
	Modify(ctx context.Context, request *ModifyUserProgrammaticAccessTokenRequest) error
	Rotate(ctx context.Context, request *RotateUserProgrammaticAccessTokenRequest) ([]RotateProgrammaticAccessTokenResult, error)
	Remove(ctx context.Context, request *RemoveUserProgrammaticAccessTokenRequest) error
	Show(ctx context.Context, request *ShowUserProgrammaticAccessTokenRequest) ([]ProgrammaticAccessToken, error)
}

// AddUserProgrammaticAccessTokenOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-programmatic-access-token.
type AddUserProgrammaticAccessTokenOptions struct {
	alter                                bool                     `ddl:"static" sql:"ALTER"`
	user                                 bool                     `ddl:"static" sql:"USER"`
	IfExists                             *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	User                                 *AccountObjectIdentifier `ddl:"identifier"`
	addProgrammaticAccessToken           bool                     `ddl:"static" sql:"ADD PROGRAMMATIC ACCESS TOKEN"`
	name                                 string                   `ddl:"keyword,double_quotes"`
	RoleRestriction                      *AccountObjectIdentifier `ddl:"identifier,equals" sql:"ROLE_RESTRICTION"`
	DaysToExpiry                         *int                     `ddl:"parameter" sql:"DAYS_TO_EXPIRY"`
	MinsToBypassNetworkPolicyRequirement *int                     `ddl:"parameter" sql:"MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT"`
	Comment                              *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type addProgrammaticAccessTokenResultDBRow struct {
	TokenName   string `db:"token_name"`
	TokenSecret string `db:"token_secret"`
}

type AddProgrammaticAccessTokenResult struct {
	TokenName   string
	TokenSecret string
}

// ModifyUserProgrammaticAccessTokenOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-user-modify-programmatic-access-token.
type ModifyUserProgrammaticAccessTokenOptions struct {
	alter                         bool                                `ddl:"static" sql:"ALTER"`
	user                          bool                                `ddl:"static" sql:"USER"`
	IfExists                      *bool                               `ddl:"keyword" sql:"IF EXISTS"`
	User                          *AccountObjectIdentifier            `ddl:"identifier"`
	modifyProgrammaticAccessToken bool                                `ddl:"static" sql:"MODIFY PROGRAMMATIC ACCESS TOKEN"`
	name                          string                              `ddl:"keyword,double_quotes"`
	Set                           *ModifyProgrammaticAccessTokenSet   `ddl:"keyword" sql:"SET"`
	Unset                         *ModifyProgrammaticAccessTokenUnset `ddl:"list,no_parentheses" sql:"UNSET"`
	RenameTo                      *string                             `ddl:"parameter,double_quotes,no_equals" sql:"RENAME TO"`
}

type ModifyProgrammaticAccessTokenSet struct {
	Disabled                             *bool   `ddl:"parameter" sql:"DISABLED"`
	MinsToBypassNetworkPolicyRequirement *int    `ddl:"parameter" sql:"MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT"`
	Comment                              *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ModifyProgrammaticAccessTokenUnset struct {
	Disabled                             *bool `ddl:"keyword" sql:"DISABLED"`
	MinsToBypassNetworkPolicyRequirement *bool `ddl:"keyword" sql:"MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT"`
	Comment                              *bool `ddl:"keyword" sql:"COMMENT"`
}

// RotateUserProgrammaticAccessTokenOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-user-rotate-programmatic-access-token.
type RotateUserProgrammaticAccessTokenOptions struct {
	alter                         bool                     `ddl:"static" sql:"ALTER"`
	user                          bool                     `ddl:"static" sql:"USER"`
	IfExists                      *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	User                          *AccountObjectIdentifier `ddl:"identifier"`
	rotateProgrammaticAccessToken bool                     `ddl:"static" sql:"ROTATE PROGRAMMATIC ACCESS TOKEN"`
	name                          string                   `ddl:"keyword,double_quotes"`
	ExpireRotatedTokenAfterHours  *int                     `ddl:"parameter" sql:"EXPIRE_ROTATED_TOKEN_AFTER_HOURS"`
}

type rotateProgrammaticAccessTokenResultDBRow struct {
	TokenName        string `db:"token_name"`
	TokenSecret      string `db:"token_secret"`
	RotatedTokenName string `db:"rotated_token_name"`
}

type RotateProgrammaticAccessTokenResult struct {
	TokenName        string
	TokenSecret      string
	RotatedTokenName string
}

// RemoveUserProgrammaticAccessTokenOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-user-remove-programmatic-access-token.
type RemoveUserProgrammaticAccessTokenOptions struct {
	alter                         bool                     `ddl:"static" sql:"ALTER"`
	user                          bool                     `ddl:"static" sql:"USER"`
	IfExists                      *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	User                          *AccountObjectIdentifier `ddl:"identifier"`
	removeProgrammaticAccessToken bool                     `ddl:"static" sql:"REMOVE PROGRAMMATIC ACCESS TOKEN"`
	name                          string                   `ddl:"keyword,double_quotes"`
}

// ShowUserProgrammaticAccessTokenOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-user-programmatic-access-tokens.
type ShowUserProgrammaticAccessTokenOptions struct {
	show                         bool                     `ddl:"static" sql:"SHOW"`
	userProgrammaticAccessTokens bool                     `ddl:"static" sql:"USER PROGRAMMATIC ACCESS TOKENS"`
	User                         *AccountObjectIdentifier `ddl:"identifier" sql:"FOR USER"`
}

type programmaticAccessTokenResultDBRow struct {
	Name                                 string         `db:"name"`
	UserName                             string         `db:"user_name"`
	RoleRestriction                      string         `db:"role_restriction"`
	ExpiresAt                            time.Time      `db:"expires_at"`
	Status                               string         `db:"status"`
	Comment                              sql.NullString `db:"comment"`
	CreatedOn                            time.Time      `db:"created_on"`
	CreatedBy                            string         `db:"created_by"`
	MinsToBypassNetworkPolicyRequirement sql.NullInt64  `db:"mins_to_bypass_network_policy_requirement"`
	RotatedTo                            sql.NullString `db:"rotated_to"`
}

type ProgrammaticAccessToken struct {
	Name                                 string
	UserName                             AccountObjectIdentifier
	RoleRestriction                      *AccountObjectIdentifier
	ExpiresAt                            time.Time
	Status                               ProgrammaticAccessTokenStatus
	Comment                              *string
	CreatedOn                            time.Time
	CreatedBy                            string
	MinsToBypassNetworkPolicyRequirement *int
	RotatedTo                            *string
}
