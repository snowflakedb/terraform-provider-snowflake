package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Users interface {
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateUserOptions) error
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterUserOptions) error
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropUserOptions) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Describe(ctx context.Context, id AccountObjectIdentifier) (*UserDetails, error)
	Show(ctx context.Context, opts *ShowUserOptions) ([]User, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*User, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*User, error)
	ShowParameters(ctx context.Context, id AccountObjectIdentifier) ([]*Parameter, error)

	AddProgrammaticAccessToken(ctx context.Context, request *AddUserProgrammaticAccessTokenRequest) (*AddProgrammaticAccessTokenResult, error)
	ModifyProgrammaticAccessToken(ctx context.Context, request *ModifyUserProgrammaticAccessTokenRequest) error
	RotateProgrammaticAccessToken(ctx context.Context, request *RotateUserProgrammaticAccessTokenRequest) (*RotateProgrammaticAccessTokenResult, error)
	RemoveProgrammaticAccessToken(ctx context.Context, request *RemoveUserProgrammaticAccessTokenRequest) error
	RemoveProgrammaticAccessTokenSafely(ctx context.Context, request *RemoveUserProgrammaticAccessTokenRequest) error
	ShowProgrammaticAccessTokens(ctx context.Context, request *ShowUserProgrammaticAccessTokenRequest) ([]ProgrammaticAccessToken, error)
	ShowProgrammaticAccessTokenByName(ctx context.Context, userId AccountObjectIdentifier, tokenName AccountObjectIdentifier) (*ProgrammaticAccessToken, error)
	ShowProgrammaticAccessTokenByNameSafely(ctx context.Context, userId AccountObjectIdentifier, tokenName AccountObjectIdentifier) (*ProgrammaticAccessToken, error)
	ShowUserWorkloadIdentityAuthenticationMethodOptions(ctx context.Context, userId AccountObjectIdentifier) ([]UserWorkloadIdentityAuthenticationMethod, error)
}

var _ Users = (*users)(nil)

type users struct {
	client *Client
}

// CreateUserOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-user.
type CreateUserOptions struct {
	create            bool                    `ddl:"static" sql:"CREATE"`
	OrReplace         *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	user              bool                    `ddl:"static" sql:"USER"`
	IfNotExists       *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name              AccountObjectIdentifier `ddl:"identifier"`
	ObjectProperties  *UserObjectProperties   `ddl:"keyword"`
	ObjectParameters  *UserObjectParameters   `ddl:"keyword"`
	SessionParameters *SessionParameters      `ddl:"keyword"`
	With              *bool                   `ddl:"keyword" sql:"WITH"`
	Tags              []TagAssociation        `ddl:"keyword,parentheses" sql:"TAG"`
}

type UserTag struct {
	Name  ObjectIdentifier `ddl:"keyword"`
	Value string           `ddl:"parameter,single_quotes"`
}

type UserObjectProperties struct {
	Password              *string                               `ddl:"parameter,single_quotes" sql:"PASSWORD"`
	LoginName             *string                               `ddl:"parameter,single_quotes" sql:"LOGIN_NAME"`
	DisplayName           *string                               `ddl:"parameter,single_quotes" sql:"DISPLAY_NAME"`
	FirstName             *string                               `ddl:"parameter,single_quotes" sql:"FIRST_NAME"`
	MiddleName            *string                               `ddl:"parameter,single_quotes" sql:"MIDDLE_NAME"`
	LastName              *string                               `ddl:"parameter,single_quotes" sql:"LAST_NAME"`
	Email                 *string                               `ddl:"parameter,single_quotes" sql:"EMAIL"`
	MustChangePassword    *bool                                 `ddl:"parameter,no_quotes" sql:"MUST_CHANGE_PASSWORD"`
	Disable               *bool                                 `ddl:"parameter,no_quotes" sql:"DISABLED"`
	DaysToExpiry          *int                                  `ddl:"parameter,no_quotes" sql:"DAYS_TO_EXPIRY"`
	MinsToUnlock          *int                                  `ddl:"parameter,no_quotes" sql:"MINS_TO_UNLOCK"`
	DefaultWarehouse      *AccountObjectIdentifier              `ddl:"identifier,equals" sql:"DEFAULT_WAREHOUSE"`
	DefaultNamespace      *ObjectIdentifier                     `ddl:"identifier,equals" sql:"DEFAULT_NAMESPACE"`
	DefaultRole           *AccountObjectIdentifier              `ddl:"identifier,equals" sql:"DEFAULT_ROLE"`
	DefaultSecondaryRoles *SecondaryRoles                       `ddl:"parameter,equals" sql:"DEFAULT_SECONDARY_ROLES"`
	MinsToBypassMFA       *int                                  `ddl:"parameter,no_quotes" sql:"MINS_TO_BYPASS_MFA"`
	RSAPublicKey          *string                               `ddl:"parameter,single_quotes" sql:"RSA_PUBLIC_KEY"`
	RSAPublicKeyFp        *string                               `ddl:"parameter,single_quotes" sql:"RSA_PUBLIC_KEY_FP"`
	RSAPublicKey2         *string                               `ddl:"parameter,single_quotes" sql:"RSA_PUBLIC_KEY_2"`
	RSAPublicKey2Fp       *string                               `ddl:"parameter,single_quotes" sql:"RSA_PUBLIC_KEY_2_FP"`
	Type                  *UserType                             `ddl:"parameter,no_quotes" sql:"TYPE"`
	WorkloadIdentity      *UserObjectWorkloadIdentityProperties `ddl:"list,parentheses,no_comma" sql:"WORKLOAD_IDENTITY ="`
	Comment               *string                               `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type UserObjectWorkloadIdentityProperties struct {
	AwsType   *UserObjectWorkloadIdentityAws   `ddl:"keyword"`
	AzureType *UserObjectWorkloadIdentityAzure `ddl:"keyword"`
	GcpType   *UserObjectWorkloadIdentityGcp   `ddl:"keyword"`
	OidcType  *UserObjectWorkloadIdentityOidc  `ddl:"keyword"`
}

type UserObjectWorkloadIdentityAws struct {
	wifType string  `ddl:"static" sql:"TYPE = AWS"`
	Arn     *string `ddl:"parameter,single_quotes" sql:"ARN"`
}

type UserObjectWorkloadIdentityAzure struct {
	wifType string  `ddl:"static" sql:"TYPE = AZURE"`
	Issuer  *string `ddl:"parameter,single_quotes" sql:"ISSUER"`
	Subject *string `ddl:"parameter,single_quotes" sql:"SUBJECT"`
}

type UserObjectWorkloadIdentityGcp struct {
	wifType string  `ddl:"static" sql:"TYPE = GCP"`
	Subject *string `ddl:"parameter,single_quotes" sql:"SUBJECT"`
}

type UserObjectWorkloadIdentityOidc struct {
	wifType          string                  `ddl:"static" sql:"TYPE = OIDC"`
	Issuer           *string                 `ddl:"parameter,single_quotes" sql:"ISSUER"`
	Subject          *string                 `ddl:"parameter,single_quotes" sql:"SUBJECT"`
	OidcAudienceList []StringListItemWrapper `ddl:"parameter,parentheses" sql:"OIDC_AUDIENCE_LIST"`
}

type SecondaryRoles struct {
	None *bool `ddl:"static" sql:"()"`
	All  *bool `ddl:"static" sql:"('ALL')"`
}

type SecondaryRole struct {
	Value string `ddl:"keyword,single_quotes"`
}

type UserObjectParameters struct {
	EnableUnredactedQuerySyntaxError *bool                    `ddl:"parameter,no_quotes" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
	NetworkPolicy                    *AccountObjectIdentifier `ddl:"identifier,equals" sql:"NETWORK_POLICY"`
	PreventUnloadToInternalStages    *bool                    `ddl:"parameter,no_quotes" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
}

// AlterUserOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-user.
type AlterUserOptions struct {
	alter    bool                    `ddl:"static" sql:"ALTER"`
	user     bool                    `ddl:"static" sql:"USER"`
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`

	// one of
	NewName                      AccountObjectIdentifier       `ddl:"identifier" sql:"RENAME TO"`
	ResetPassword                *bool                         `ddl:"keyword" sql:"RESET PASSWORD"`
	AbortAllQueries              *bool                         `ddl:"keyword" sql:"ABORT ALL QUERIES"`
	AddDelegatedAuthorization    *AddDelegatedAuthorization    `ddl:"keyword"`
	RemoveDelegatedAuthorization *RemoveDelegatedAuthorization `ddl:"keyword"`
	Set                          *UserSet                      `ddl:"keyword" sql:"SET"`
	Unset                        *UserUnset                    `ddl:"list" sql:"UNSET"`
	SetTag                       []TagAssociation              `ddl:"keyword" sql:"SET TAG"`
	UnsetTag                     []ObjectIdentifier            `ddl:"keyword" sql:"UNSET TAG"`
}

type AddDelegatedAuthorization struct {
	Role        string `ddl:"parameter,no_equals" sql:"ADD DELEGATED AUTHORIZATION OF ROLE"`
	Integration string `ddl:"parameter,no_equals" sql:"TO SECURITY INTEGRATION"`
}

type RemoveDelegatedAuthorization struct {
	// one of
	Role           *string `ddl:"parameter,no_equals" sql:"REMOVE DELEGATED AUTHORIZATION OF ROLE"`
	Authorizations *bool   `ddl:"parameter,no_equals" sql:"REMOVE DELEGATED AUTHORIZATIONS"`

	Integration string `ddl:"parameter,no_equals" sql:"FROM SECURITY INTEGRATION"`
}

type UserAlterObjectProperties struct {
	UserObjectProperties
	DisableMfa *bool `ddl:"parameter,no_quotes" sql:"DISABLE_MFA"`
}

type UserObjectPropertiesUnset struct {
	Password              *bool `ddl:"keyword" sql:"PASSWORD"`
	LoginName             *bool `ddl:"keyword" sql:"LOGIN_NAME"`
	DisplayName           *bool `ddl:"keyword" sql:"DISPLAY_NAME"`
	FirstName             *bool `ddl:"keyword" sql:"FIRST_NAME"`
	MiddleName            *bool `ddl:"keyword" sql:"MIDDLE_NAME"`
	LastName              *bool `ddl:"keyword" sql:"LAST_NAME"`
	Email                 *bool `ddl:"keyword" sql:"EMAIL"`
	MustChangePassword    *bool `ddl:"keyword" sql:"MUST_CHANGE_PASSWORD"`
	Disable               *bool `ddl:"keyword" sql:"DISABLED"`
	DaysToExpiry          *bool `ddl:"keyword" sql:"DAYS_TO_EXPIRY"`
	MinsToUnlock          *bool `ddl:"keyword" sql:"MINS_TO_UNLOCK"`
	DefaultWarehouse      *bool `ddl:"keyword" sql:"DEFAULT_WAREHOUSE"`
	DefaultNamespace      *bool `ddl:"keyword" sql:"DEFAULT_NAMESPACE"`
	DefaultRole           *bool `ddl:"keyword" sql:"DEFAULT_ROLE"`
	DefaultSecondaryRoles *bool `ddl:"keyword" sql:"DEFAULT_SECONDARY_ROLES"`
	MinsToBypassMFA       *bool `ddl:"keyword" sql:"MINS_TO_BYPASS_MFA"`
	DisableMfa            *bool `ddl:"keyword" sql:"DISABLE_MFA"`
	RSAPublicKey          *bool `ddl:"keyword" sql:"RSA_PUBLIC_KEY"`
	RSAPublicKey2         *bool `ddl:"keyword" sql:"RSA_PUBLIC_KEY_2"`
	Type                  *bool `ddl:"keyword" sql:"TYPE"`
	WorkloadIdentity      *bool `ddl:"keyword" sql:"WORKLOAD_IDENTITY"`
	Comment               *bool `ddl:"keyword" sql:"COMMENT"`
}

type UserObjectParametersUnset struct {
	EnableUnredactedQuerySyntaxError *bool `ddl:"keyword" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
	NetworkPolicy                    *bool `ddl:"keyword" sql:"NETWORK_POLICY"`
	PreventUnloadToInternalStages    *bool `ddl:"keyword" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
}

type UserSet struct {
	PasswordPolicy       *SchemaObjectIdentifier    `ddl:"identifier" sql:"PASSWORD POLICY"`
	SessionPolicy        *SchemaObjectIdentifier    `ddl:"identifier" sql:"SESSION POLICY"`
	AuthenticationPolicy *SchemaObjectIdentifier    `ddl:"identifier" sql:"AUTHENTICATION POLICY"`
	ObjectProperties     *UserAlterObjectProperties `ddl:"keyword"`
	ObjectParameters     *UserObjectParameters      `ddl:"keyword"`
	SessionParameters    *SessionParameters         `ddl:"keyword"`
}

type UserUnset struct {
	PasswordPolicy       *bool                      `ddl:"keyword" sql:"PASSWORD POLICY"`
	SessionPolicy        *bool                      `ddl:"keyword" sql:"SESSION POLICY"`
	AuthenticationPolicy *bool                      `ddl:"keyword" sql:"AUTHENTICATION POLICY"`
	ObjectProperties     *UserObjectPropertiesUnset `ddl:"list"`
	ObjectParameters     *UserObjectParametersUnset `ddl:"list"`
	SessionParameters    *SessionParametersUnset    `ddl:"list"`
}

// DropUserOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-user.
type DropUserOptions struct {
	drop     bool                    `ddl:"static" sql:"DROP"`
	user     bool                    `ddl:"static" sql:"USER"`
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

// describeUserOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-user.
type describeUserOptions struct {
	describe bool                    `ddl:"static" sql:"DESCRIBE"`
	user     bool                    `ddl:"static" sql:"USER"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

// UserDetails contains details about a user.
type UserDetails struct {
	Name                                *StringProperty
	Comment                             *StringProperty
	DisplayName                         *StringProperty
	Type                                *StringProperty
	LoginName                           *StringProperty
	FirstName                           *StringProperty
	MiddleName                          *StringProperty
	LastName                            *StringProperty
	Email                               *StringProperty
	Password                            *StringProperty
	MustChangePassword                  *BoolProperty
	Disabled                            *BoolProperty
	SnowflakeLock                       *BoolProperty
	SnowflakeSupport                    *BoolProperty
	DaysToExpiry                        *FloatProperty
	MinsToUnlock                        *IntProperty
	DefaultWarehouse                    *StringProperty
	DefaultNamespace                    *StringProperty
	DefaultRole                         *StringProperty
	DefaultSecondaryRoles               *StringProperty
	ExtAuthnDuo                         *BoolProperty
	ExtAuthnUid                         *StringProperty
	MinsToBypassMfa                     *IntProperty
	MinsToBypassNetworkPolicy           *IntProperty
	RsaPublicKey                        *StringProperty
	RsaPublicKeyFp                      *StringProperty
	RsaPublicKeyLastSetTime             *StringProperty
	RsaPublicKey2                       *StringProperty
	RsaPublicKey2Fp                     *StringProperty
	RsaPublicKey2LastSetTime            *StringProperty
	PasswordLastSetTime                 *StringProperty
	CustomLandingPageUrl                *StringProperty
	CustomLandingPageUrlFlushNextUiLoad *BoolProperty
	HasMfa                              *BoolProperty
	HasWorkloadIdentity                 *BoolProperty
}

// ShowUserOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-users.
type ShowUserOptions struct {
	show       bool    `ddl:"static" sql:"SHOW"`
	Terse      *bool   `ddl:"static" sql:"TERSE"`
	users      bool    `ddl:"static" sql:"USERS"`
	Like       *Like   `ddl:"keyword" sql:"LIKE"`
	StartsWith *string `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit      *int    `ddl:"parameter,no_equals" sql:"LIMIT"`
	From       *string `ddl:"parameter,no_equals,single_quotes" sql:"FROM"`
}

type User struct {
	Name                  string
	CreatedOn             time.Time
	LoginName             string
	DisplayName           string
	FirstName             string
	LastName              string
	Email                 string
	MinsToUnlock          string
	DaysToExpiry          string
	Comment               string
	Disabled              bool
	MustChangePassword    bool
	SnowflakeLock         bool
	DefaultWarehouse      string
	DefaultNamespace      string
	DefaultRole           string
	DefaultSecondaryRoles string
	ExtAuthnDuo           bool
	ExtAuthnUid           string
	MinsToBypassMfa       string
	Owner                 string
	LastSuccessLogin      time.Time
	ExpiresAtTime         time.Time
	LockedUntilTime       time.Time
	HasPassword           bool
	HasRsaPublicKey       bool
	Type                  string
	HasMfa                bool
	HasWorkloadIdentity   bool
}

func (v *User) ID() AccountObjectIdentifier {
	return AccountObjectIdentifier{v.Name}
}

func (v *User) ObjectType() ObjectType {
	return ObjectTypeUser
}

type userDBRow struct {
	Name                  string         `db:"name"`
	CreatedOn             time.Time      `db:"created_on"`
	LoginName             sql.NullString `db:"login_name"`
	DisplayName           sql.NullString `db:"display_name"`
	FirstName             sql.NullString `db:"first_name"`
	LastName              sql.NullString `db:"last_name"`
	Email                 sql.NullString `db:"email"`
	MinsToUnlock          sql.NullString `db:"mins_to_unlock"`
	DaysToExpiry          sql.NullString `db:"days_to_expiry"`
	Comment               sql.NullString `db:"comment"`
	Disabled              sql.NullString `db:"disabled"`
	MustChangePassword    sql.NullString `db:"must_change_password"`
	SnowflakeLock         sql.NullString `db:"snowflake_lock"`
	DefaultWarehouse      sql.NullString `db:"default_warehouse"`
	DefaultNamespace      sql.NullString `db:"default_namespace"`
	DefaultRole           sql.NullString `db:"default_role"`
	DefaultSecondaryRoles sql.NullString `db:"default_secondary_roles"`
	ExtAuthnDuo           sql.NullString `db:"ext_authn_duo"`
	ExtAuthnUid           sql.NullString `db:"ext_authn_uid"`
	MinsToBypassMfa       sql.NullString `db:"mins_to_bypass_mfa"`
	Owner                 string         `db:"owner"`
	LastSuccessLogin      sql.NullTime   `db:"last_success_login"`
	ExpiresAtTime         sql.NullTime   `db:"expires_at_time"`
	LockedUntilTime       sql.NullTime   `db:"locked_until_time"`
	HasPassword           sql.NullBool   `db:"has_password"`
	HasRsaPublicKey       sql.NullBool   `db:"has_rsa_public_key"`
	Type                  sql.NullString `db:"type"`
	HasMfa                sql.NullBool   `db:"has_mfa"`
	HasWorkloadIdentity   sql.NullBool   `db:"has_workload_identity"`
}

// showUserAuthenticationMethodOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-user-workload-identity-authentication-methods
type showUserAuthenticationMethodOptions struct {
	show                            bool                    `ddl:"static" sql:"SHOW"`
	userWorkloadIdentityAuthMethods bool                    `ddl:"static" sql:"USER WORKLOAD IDENTITY AUTHENTICATION METHODS"`
	ForUser                         AccountObjectIdentifier `ddl:"identifier,no_equals,no_quotes" sql:"FOR USER"`
}

type userWorkloadIdentityAuthenticationMethodsDBRow struct {
	Name           string         `db:"name"`
	Type           string         `db:"type"`
	Comment        sql.NullString `db:"comment"`
	LastUsed       sql.NullTime   `db:"last_used"`
	CreatedOn      time.Time      `db:"created_on"`
	AdditionalInfo sql.NullString `db:"additional_info"`
}

type UserWorkloadIdentityAuthenticationMethod struct {
	Name                string
	Type                WIFType
	Comment             string
	LastUsed            time.Time
	CreatedOn           time.Time
	AwsAdditionalInfo   *UserWorkloadIdentityAuthenticationMethodsAwsAdditionalInfo
	AzureAdditionalInfo *UserWorkloadIdentityAuthenticationMethodsAzureAdditionalInfo
	GcpAdditionalInfo   *UserWorkloadIdentityAuthenticationMethodsGcpAdditionalInfo
	OidcAdditionalInfo  *UserWorkloadIdentityAuthenticationMethodsOidcAdditionalInfo
}

func (v *UserWorkloadIdentityAuthenticationMethod) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

type UserWorkloadIdentityAuthenticationMethodsAwsAdditionalInfo struct {
	IamRole      string `json:"iamRole"`
	Type         string `json:"type"`
	AwsAccount   string `json:"awsAccount"`
	AwsPartition string `json:"awsPartition"`
}

type UserWorkloadIdentityAuthenticationMethodsAzureAdditionalInfo struct {
	Issuer  string `json:"issuer"`
	Subject string `json:"subject"`
}

type UserWorkloadIdentityAuthenticationMethodsGcpAdditionalInfo struct {
	Subject string `json:"subject"`
}

type UserWorkloadIdentityAuthenticationMethodsOidcAdditionalInfo struct {
	Issuer       string   `json:"issuer"`
	Subject      string   `json:"subject"`
	AudienceList []string `json:"audienceList"`
}
