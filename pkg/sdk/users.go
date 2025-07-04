package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ validatable = new(CreateUserOptions)
	_ validatable = new(AlterUserOptions)
	_ validatable = new(DropUserOptions)
	_ validatable = new(describeUserOptions)
	_ validatable = new(ShowUserOptions)
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
}

var _ Users = (*users)(nil)

type users struct {
	client *Client
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
}

func (v *User) GetSecondaryRolesOption() SecondaryRolesOption {
	return GetSecondaryRolesOptionFrom(v.DefaultSecondaryRoles)
}

func GetSecondaryRolesOptionFrom(text string) SecondaryRolesOption {
	if text != "" {
		parsedRoles := ParseCommaSeparatedStringArray(text, true)
		if len(parsedRoles) > 0 {
			return SecondaryRolesOptionAll
		} else {
			return SecondaryRolesOptionNone
		}
	}
	return SecondaryRolesOptionDefault
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
}

func (row userDBRow) convert() *User {
	user := &User{
		Name:      row.Name,
		CreatedOn: row.CreatedOn,
		Owner:     row.Owner,
	}
	if row.LoginName.Valid {
		user.LoginName = row.LoginName.String
	}
	if row.DisplayName.Valid {
		user.DisplayName = row.DisplayName.String
	}
	if row.FirstName.Valid {
		user.FirstName = row.FirstName.String
	}
	if row.LastName.Valid {
		user.LastName = row.LastName.String
	}
	if row.Email.Valid {
		user.Email = row.Email.String
	}
	if row.MinsToUnlock.Valid {
		user.MinsToUnlock = row.MinsToUnlock.String
	}
	if row.DaysToExpiry.Valid {
		user.DaysToExpiry = row.DaysToExpiry.String
	}
	if row.Comment.Valid {
		user.Comment = row.Comment.String
	}
	handleNullableBoolString(row.Disabled, &user.Disabled)
	handleNullableBoolString(row.MustChangePassword, &user.MustChangePassword)
	handleNullableBoolString(row.SnowflakeLock, &user.SnowflakeLock)
	handleNullableBoolString(row.ExtAuthnDuo, &user.ExtAuthnDuo)
	if row.ExtAuthnUid.Valid {
		user.ExtAuthnUid = row.ExtAuthnUid.String
	}
	if row.MinsToBypassMfa.Valid {
		user.MinsToBypassMfa = row.MinsToBypassMfa.String
	}
	if row.DefaultWarehouse.Valid {
		user.DefaultWarehouse = row.DefaultWarehouse.String
	}
	if row.DefaultNamespace.Valid {
		user.DefaultNamespace = row.DefaultNamespace.String
	}
	if row.DefaultRole.Valid {
		user.DefaultRole = row.DefaultRole.String
	}
	if row.DefaultSecondaryRoles.Valid {
		user.DefaultSecondaryRoles = row.DefaultSecondaryRoles.String
	}
	if row.LastSuccessLogin.Valid {
		user.LastSuccessLogin = row.LastSuccessLogin.Time
	}
	if row.ExpiresAtTime.Valid {
		user.ExpiresAtTime = row.ExpiresAtTime.Time
	}
	if row.LockedUntilTime.Valid {
		user.LockedUntilTime = row.LockedUntilTime.Time
	}
	if row.HasPassword.Valid {
		user.HasPassword = row.HasPassword.Bool
	}
	if row.HasRsaPublicKey.Valid {
		user.HasRsaPublicKey = row.HasRsaPublicKey.Bool
	}
	if row.Type.Valid {
		user.Type = row.Type.String
	}
	if row.HasMfa.Valid {
		user.HasMfa = row.HasMfa.Bool
	}
	return user
}

func (v *User) ID() AccountObjectIdentifier {
	return AccountObjectIdentifier{v.Name}
}

func (v *User) ObjectType() ObjectType {
	return ObjectTypeUser
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

func (opts *CreateUserOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.ObjectProperties) && valueSet(opts.ObjectProperties.DefaultSecondaryRoles) {
		if err := opts.ObjectProperties.DefaultSecondaryRoles.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (v *users) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateUserOptions) error {
	if opts == nil {
		opts = &CreateUserOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

type UserObjectProperties struct {
	Password              *string                  `ddl:"parameter,single_quotes" sql:"PASSWORD"`
	LoginName             *string                  `ddl:"parameter,single_quotes" sql:"LOGIN_NAME"`
	DisplayName           *string                  `ddl:"parameter,single_quotes" sql:"DISPLAY_NAME"`
	FirstName             *string                  `ddl:"parameter,single_quotes" sql:"FIRST_NAME"`
	MiddleName            *string                  `ddl:"parameter,single_quotes" sql:"MIDDLE_NAME"`
	LastName              *string                  `ddl:"parameter,single_quotes" sql:"LAST_NAME"`
	Email                 *string                  `ddl:"parameter,single_quotes" sql:"EMAIL"`
	MustChangePassword    *bool                    `ddl:"parameter,no_quotes" sql:"MUST_CHANGE_PASSWORD"`
	Disable               *bool                    `ddl:"parameter,no_quotes" sql:"DISABLED"`
	DaysToExpiry          *int                     `ddl:"parameter,no_quotes" sql:"DAYS_TO_EXPIRY"`
	MinsToUnlock          *int                     `ddl:"parameter,no_quotes" sql:"MINS_TO_UNLOCK"`
	DefaultWarehouse      *AccountObjectIdentifier `ddl:"identifier,equals" sql:"DEFAULT_WAREHOUSE"`
	DefaultNamespace      *ObjectIdentifier        `ddl:"identifier,equals" sql:"DEFAULT_NAMESPACE"`
	DefaultRole           *AccountObjectIdentifier `ddl:"identifier,equals" sql:"DEFAULT_ROLE"`
	DefaultSecondaryRoles *SecondaryRoles          `ddl:"parameter,equals" sql:"DEFAULT_SECONDARY_ROLES"`
	MinsToBypassMFA       *int                     `ddl:"parameter,no_quotes" sql:"MINS_TO_BYPASS_MFA"`
	RSAPublicKey          *string                  `ddl:"parameter,single_quotes" sql:"RSA_PUBLIC_KEY"`
	RSAPublicKeyFp        *string                  `ddl:"parameter,single_quotes" sql:"RSA_PUBLIC_KEY_FP"`
	RSAPublicKey2         *string                  `ddl:"parameter,single_quotes" sql:"RSA_PUBLIC_KEY_2"`
	RSAPublicKey2Fp       *string                  `ddl:"parameter,single_quotes" sql:"RSA_PUBLIC_KEY_2_FP"`
	Type                  *UserType                `ddl:"parameter,no_quotes" sql:"TYPE"`
	Comment               *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type UserAlterObjectProperties struct {
	UserObjectProperties
	DisableMfa *bool `ddl:"parameter,no_quotes" sql:"DISABLE_MFA"`
}

type SecondaryRoles struct {
	None *bool `ddl:"static" sql:"()"`
	All  *bool `ddl:"static" sql:"('ALL')"`
}

type SecondaryRole struct {
	Value string `ddl:"keyword,single_quotes"`
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
	Comment               *bool `ddl:"keyword" sql:"COMMENT"`
}

type UserObjectParameters struct {
	EnableUnredactedQuerySyntaxError *bool                    `ddl:"parameter,no_quotes" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
	NetworkPolicy                    *AccountObjectIdentifier `ddl:"identifier,equals" sql:"NETWORK_POLICY"`
	PreventUnloadToInternalStages    *bool                    `ddl:"parameter,no_quotes" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
}
type UserObjectParametersUnset struct {
	EnableUnredactedQuerySyntaxError *bool `ddl:"keyword" sql:"ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"`
	NetworkPolicy                    *bool `ddl:"keyword" sql:"NETWORK_POLICY"`
	PreventUnloadToInternalStages    *bool `ddl:"keyword" sql:"PREVENT_UNLOAD_TO_INTERNAL_STAGES"`
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

func (opts *AlterUserOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.NewName, opts.ResetPassword, opts.AbortAllQueries, opts.AddDelegatedAuthorization, opts.RemoveDelegatedAuthorization, opts.Set, opts.Unset, opts.SetTag, opts.UnsetTag) {
		errs = append(errs, errExactlyOneOf("AlterUserOptions", "NewName", "ResetPassword", "AbortAllQueries", "AddDelegatedAuthorization", "RemoveDelegatedAuthorization", "Set", "Unset", "SetTag", "UnsetTag"))
	}
	if valueSet(opts.RemoveDelegatedAuthorization) {
		if err := opts.RemoveDelegatedAuthorization.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Set) {
		if err := opts.Set.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.Unset) {
		if err := opts.Unset.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *users) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterUserOptions) error {
	if opts == nil {
		opts = &AlterUserOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
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

func (opts *RemoveDelegatedAuthorization) validate() error {
	var errs []error
	if !exactlyOneValueSet(opts.Role, opts.Authorizations) {
		errs = append(errs, errExactlyOneOf("RemoveDelegatedAuthorization", "Role", "Authorization"))
	}
	if !valueSet(opts.Integration) {
		errs = append(errs, errNotSet("RemoveDelegatedAuthorization", "Integration"))
	}
	return errors.Join(errs...)
}

type UserSet struct {
	PasswordPolicy       *SchemaObjectIdentifier    `ddl:"identifier" sql:"PASSWORD POLICY"`
	SessionPolicy        *string                    `ddl:"parameter" sql:"SESSION POLICY"`
	AuthenticationPolicy *SchemaObjectIdentifier    `ddl:"identifier" sql:"AUTHENTICATION POLICY"`
	ObjectProperties     *UserAlterObjectProperties `ddl:"keyword"`
	ObjectParameters     *UserObjectParameters      `ddl:"keyword"`
	SessionParameters    *SessionParameters         `ddl:"keyword"`
}

func (opts *UserSet) validate() error {
	if !anyValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy, opts.ObjectProperties, opts.ObjectParameters, opts.SessionParameters) {
		return errAtLeastOneOf("UserSet", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ObjectProperties", "ObjectParameters", "SessionParameters")
	}
	if moreThanOneValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy) {
		return errOneOf("UserSet", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy")
	}
	if anyValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy) && anyValueSet(opts.ObjectProperties, opts.ObjectParameters, opts.SessionParameters) {
		return NewError("policies cannot be set with user properties or parameters at the same time")
	}
	if valueSet(opts.ObjectProperties) && valueSet(opts.ObjectProperties.DefaultSecondaryRoles) {
		if err := opts.ObjectProperties.DefaultSecondaryRoles.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (opts *SecondaryRoles) validate() error {
	if !exactlyOneValueSet(opts.All, opts.None) {
		return errExactlyOneOf("SecondaryRoles", "All", "None")
	}
	return nil
}

type UserUnset struct {
	PasswordPolicy       *bool                      `ddl:"keyword" sql:"PASSWORD POLICY"`
	SessionPolicy        *bool                      `ddl:"keyword" sql:"SESSION POLICY"`
	AuthenticationPolicy *bool                      `ddl:"keyword" sql:"AUTHENTICATION POLICY"`
	ObjectProperties     *UserObjectPropertiesUnset `ddl:"list"`
	ObjectParameters     *UserObjectParametersUnset `ddl:"list"`
	SessionParameters    *SessionParametersUnset    `ddl:"list"`
}

func (opts *UserUnset) validate() error {
	// TODO [SNOW-1645875]: change validations with policies
	if !anyValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.ObjectProperties, opts.ObjectParameters, opts.SessionParameters, opts.AuthenticationPolicy) {
		return errAtLeastOneOf("UserUnset", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ObjectProperties", "ObjectParameters", "SessionParameters")
	}
	if moreThanOneValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy) {
		return errOneOf("UserUnset", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy")
	}
	if anyValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy) && anyValueSet(opts.ObjectProperties, opts.ObjectParameters, opts.SessionParameters) {
		return NewError("policies cannot be unset with user properties or parameters at the same time")
	}
	return nil
}

// DropUserOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-user.
type DropUserOptions struct {
	drop     bool                    `ddl:"static" sql:"DROP"`
	user     bool                    `ddl:"static" sql:"USER"`
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *DropUserOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *users) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropUserOptions) error {
	if opts == nil {
		opts = &DropUserOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate drop options: %w", err)
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	if err != nil {
		return err
	}
	return err
}

func (v *users) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropUserOptions{IfExists: Bool(true)}) }, ctx, id)
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
}

func userDetailsFromRows(rows []propertyRow) *UserDetails {
	v := &UserDetails{}
	for _, row := range rows {
		switch row.Property {
		case "NAME":
			v.Name = row.toStringProperty()
		case "COMMENT":
			v.Comment = row.toStringProperty()
		case "DISPLAY_NAME":
			v.DisplayName = row.toStringProperty()
		case "TYPE":
			v.Type = row.toStringProperty()
		case "LOGIN_NAME":
			v.LoginName = row.toStringProperty()
		case "FIRST_NAME":
			v.FirstName = row.toStringProperty()
		case "MIDDLE_NAME":
			v.MiddleName = row.toStringProperty()
		case "LAST_NAME":
			v.LastName = row.toStringProperty()
		case "EMAIL":
			v.Email = row.toStringProperty()
		case "PASSWORD":
			v.Password = row.toStringProperty()
		case "MUST_CHANGE_PASSWORD":
			v.MustChangePassword = row.toBoolProperty()
		case "DISABLED":
			v.Disabled = row.toBoolProperty()
		case "SNOWFLAKE_LOCK":
			v.SnowflakeLock = row.toBoolProperty()
		case "SNOWFLAKE_SUPPORT":
			v.SnowflakeSupport = row.toBoolProperty()
		case "DAYS_TO_EXPIRY":
			v.DaysToExpiry = row.toFloatProperty()
		case "MINS_TO_UNLOCK":
			v.MinsToUnlock = row.toIntProperty()
		case "DEFAULT_WAREHOUSE":
			v.DefaultWarehouse = row.toStringProperty()
		case "DEFAULT_NAMESPACE":
			v.DefaultNamespace = row.toStringProperty()
		case "DEFAULT_ROLE":
			v.DefaultRole = row.toStringProperty()
		case "DEFAULT_SECONDARY_ROLES":
			v.DefaultSecondaryRoles = row.toStringProperty()
		case "EXT_AUTHN_DUO":
			v.ExtAuthnDuo = row.toBoolProperty()
		case "EXT_AUTHN_UID":
			v.ExtAuthnUid = row.toStringProperty()
		case "HAS_MFA":
			v.HasMfa = row.toBoolProperty()
		case "MINS_TO_BYPASS_MFA":
			v.MinsToBypassMfa = row.toIntProperty()
		case "MINS_TO_BYPASS_NETWORK_POLICY":
			v.MinsToBypassNetworkPolicy = row.toIntProperty()
		case "RSA_PUBLIC_KEY":
			v.RsaPublicKey = row.toStringProperty()
		case "RSA_PUBLIC_KEY_FP":
			v.RsaPublicKeyFp = row.toStringProperty()
		case "RSA_PUBLIC_KEY_LAST_SET_TIME":
			v.RsaPublicKeyLastSetTime = row.toStringProperty()
		case "RSA_PUBLIC_KEY_2":
			v.RsaPublicKey2 = row.toStringProperty()
		case "RSA_PUBLIC_KEY_2_FP":
			v.RsaPublicKey2Fp = row.toStringProperty()
		case "RSA_PUBLIC_KEY_2_LAST_SET_TIME":
			v.RsaPublicKey2LastSetTime = row.toStringProperty()
		case "PASSWORD_LAST_SET_TIME":
			v.PasswordLastSetTime = row.toStringProperty()
		case "CUSTOM_LANDING_PAGE_URL":
			v.CustomLandingPageUrl = row.toStringProperty()
		case "CUSTOM_LANDING_PAGE_URL_FLUSH_NEXT_UI_LOAD":
			v.CustomLandingPageUrlFlushNextUiLoad = row.toBoolProperty()
		}
	}
	return v
}

// describeUserOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-user.
type describeUserOptions struct {
	describe bool                    `ddl:"static" sql:"DESCRIBE"`
	user     bool                    `ddl:"static" sql:"USER"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

func (opts *describeUserOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *users) Describe(ctx context.Context, id AccountObjectIdentifier) (*UserDetails, error) {
	opts := &describeUserOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []propertyRow{}
	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}

	return userDetailsFromRows(dest), nil
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

func (opts *ShowUserOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (v *users) Show(ctx context.Context, opts *ShowUserOptions) ([]User, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[userDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[userDBRow, User](dbRows)
	return resultList, nil
}

func (v *users) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*User, error) {
	users, err := v.Show(ctx, &ShowUserOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}

	return collections.FindFirst(users, func(user User) bool {
		return user.ID().Name() == id.Name()
	})
}

func (v *users) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*User, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *users) ShowParameters(ctx context.Context, id AccountObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			User: id,
		},
	})
}

func (v *users) AddProgrammaticAccessToken(ctx context.Context, request *AddUserProgrammaticAccessTokenRequest) (*AddProgrammaticAccessTokenResult, error) {
	return v.client.UserProgrammaticAccessTokens.Add(ctx, request)
}

func (v *users) ModifyProgrammaticAccessToken(ctx context.Context, request *ModifyUserProgrammaticAccessTokenRequest) error {
	return v.client.UserProgrammaticAccessTokens.Modify(ctx, request)
}

func (v *users) RotateProgrammaticAccessToken(ctx context.Context, request *RotateUserProgrammaticAccessTokenRequest) (*RotateProgrammaticAccessTokenResult, error) {
	return v.client.UserProgrammaticAccessTokens.Rotate(ctx, request)
}

func (v *users) RemoveProgrammaticAccessToken(ctx context.Context, request *RemoveUserProgrammaticAccessTokenRequest) error {
	return v.client.UserProgrammaticAccessTokens.Remove(ctx, request)
}

func (v *users) RemoveProgrammaticAccessTokenSafely(ctx context.Context, request *RemoveUserProgrammaticAccessTokenRequest) error {
	removeProgrammaticAccessToken := func() error {
		return v.RemoveProgrammaticAccessToken(ctx, request)
	}
	return SafeRemoveProgrammaticAccessToken(v.client, removeProgrammaticAccessToken, ctx, request.UserName)
}

func (v *users) ShowProgrammaticAccessTokens(ctx context.Context, request *ShowUserProgrammaticAccessTokenRequest) ([]ProgrammaticAccessToken, error) {
	return v.client.UserProgrammaticAccessTokens.Show(ctx, request)
}

func (v *users) ShowProgrammaticAccessTokenByName(ctx context.Context, userId AccountObjectIdentifier, tokenName AccountObjectIdentifier) (*ProgrammaticAccessToken, error) {
	return v.client.UserProgrammaticAccessTokens.ShowByID(ctx, userId, tokenName)
}

func (v *users) ShowProgrammaticAccessTokenByNameSafely(ctx context.Context, userId AccountObjectIdentifier, tokenName AccountObjectIdentifier) (*ProgrammaticAccessToken, error) {
	return SafeShowProgrammaticAccessTokenByName(v.client, v.ShowProgrammaticAccessTokenByName, ctx, userId, tokenName)
}

type SecondaryRolesOption string

const (
	SecondaryRolesOptionDefault SecondaryRolesOption = "DEFAULT"
	SecondaryRolesOptionNone    SecondaryRolesOption = "NONE"
	SecondaryRolesOptionAll     SecondaryRolesOption = "ALL"
)

func ToSecondaryRolesOption(s string) (SecondaryRolesOption, error) {
	switch strings.ToUpper(s) {
	case string(SecondaryRolesOptionDefault):
		return SecondaryRolesOptionDefault, nil
	case string(SecondaryRolesOptionNone):
		return SecondaryRolesOptionNone, nil
	case string(SecondaryRolesOptionAll):
		return SecondaryRolesOptionAll, nil
	default:
		return "", fmt.Errorf("invalid secondary roles option: %s", s)
	}
}

var ValidSecondaryRolesOptionsString = []string{
	string(SecondaryRolesOptionDefault),
	string(SecondaryRolesOptionNone),
	string(SecondaryRolesOptionAll),
}

type UserType string

const (
	UserTypePerson        UserType = "PERSON"
	UserTypeService       UserType = "SERVICE"
	UserTypeLegacyService UserType = "LEGACY_SERVICE"
)

func ToUserType(s string) (UserType, error) {
	switch strings.ToUpper(s) {
	case string(UserTypePerson):
		return UserTypePerson, nil
	case string(UserTypeService):
		return UserTypeService, nil
	case string(UserTypeLegacyService):
		return UserTypeLegacyService, nil
	default:
		return "", fmt.Errorf("invalid user type: %s", s)
	}
}

var AllUserTypes = []UserType{
	UserTypePerson,
	UserTypeService,
	UserTypeLegacyService,
}

var AcceptableUserTypes = map[UserType][]string{
	UserTypePerson:        {"", string(UserTypePerson)},
	UserTypeService:       {string(UserTypeService)},
	UserTypeLegacyService: {string(UserTypeLegacyService)},
}
