package sdk

import (
	"context"
	"encoding/json"
	"fmt"
)

func (opts *ShowUserOptions) additionalValidations() error {
	if valueSet(opts.Limit) && opts.Limit.From != nil && opts.Limit.Rows == nil {
		return errNotSet("ShowUserOptions.Limit", "Rows")
	}
	return nil
}

func (s *CreateUserRequest) ID() AccountObjectIdentifier {
	return s.name
}

// GetSecondaryRolesOptionFrom returns the SecondaryRolesOption for the given string value
// returned by SHOW USERS.
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

func (v *User) GetSecondaryRolesOption() SecondaryRolesOption {
	return GetSecondaryRolesOptionFrom(v.DefaultSecondaryRoles)
}

// ValidSecondaryRolesOptionsString is a slice of all valid secondary roles option strings.
// Kept for backward compatibility with resource descriptions.
var ValidSecondaryRolesOptionsString = []string{
	string(SecondaryRolesOptionDefault),
	string(SecondaryRolesOptionNone),
	string(SecondaryRolesOptionAll),
}

// AcceptableUserTypes maps each UserType to the set of raw strings (from Snowflake) that map to it.
var AcceptableUserTypes = map[UserType][]string{
	UserTypePerson:        {"", string(UserTypePerson)},
	UserTypeService:       {string(UserTypeService)},
	UserTypeLegacyService: {string(UserTypeLegacyService)},
}

// userDetailsFromRows builds a UserDetails from the []UserProperty slice returned by Describe.
// Callers that need *UserDetails should use DescribeDetails instead of Describe.
func userDetailsFromRows(rows []UserProperty) *UserDetails {
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
		case "HAS_WORKLOAD_IDENTITY":
			v.HasWorkloadIdentity = row.toBoolProperty()
		}
	}
	return v
}

// UserWorkloadIdentityAuthenticationMethod AdditionalInfo sub-types — kept in ext because the generator
// cannot produce structs with JSON tags.

type UserWorkloadIdentityAuthenticationMethodsAwsAdditionalInfo struct {
	IamRole      string `json:"iamRole"`
	Type         string `json:"type"`
	AwsAccount   string `json:"awsAccount"`
	AwsPartition string `json:"awsPartition"`
	Issuer       string `json:"issuer"`
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

// additionalConvert handles Type (string→WIFType with error) and AdditionalInfo (JSON→typed sub-structs).
// Called by the generated convert() because type and additional_info fields are WithManualConvert().
func (row *userWorkloadIdentityAuthenticationMethodsDBRow) additionalConvert(result *UserWorkloadIdentityAuthenticationMethod) error {
	wifType, err := ToWIFType(row.Type)
	if err != nil {
		return err
	}
	result.Type = wifType
	switch wifType {
	case WIFTypeAws:
		additionalInfo := &UserWorkloadIdentityAuthenticationMethodsAwsAdditionalInfo{}
		if err := json.Unmarshal([]byte(row.AdditionalInfo.String), additionalInfo); err != nil {
			return err
		}
		result.AwsAdditionalInfo = additionalInfo
	case WIFTypeAzure:
		additionalInfo := &UserWorkloadIdentityAuthenticationMethodsAzureAdditionalInfo{}
		if err := json.Unmarshal([]byte(row.AdditionalInfo.String), additionalInfo); err != nil {
			return err
		}
		result.AzureAdditionalInfo = additionalInfo
	case WIFTypeGcp:
		additionalInfo := &UserWorkloadIdentityAuthenticationMethodsGcpAdditionalInfo{}
		if err := json.Unmarshal([]byte(row.AdditionalInfo.String), additionalInfo); err != nil {
			return err
		}
		result.GcpAdditionalInfo = additionalInfo
	case WIFTypeOidc:
		additionalInfo := &UserWorkloadIdentityAuthenticationMethodsOidcAdditionalInfo{}
		if err := json.Unmarshal([]byte(row.AdditionalInfo.String), additionalInfo); err != nil {
			return err
		}
		result.OidcAdditionalInfo = additionalInfo
	}
	return nil
}

// ID returns the identifier for a UserWorkloadIdentityAuthenticationMethod.
// Kept in ext because the generator does not produce ID() for custom show result types.
func (v *UserWorkloadIdentityAuthenticationMethod) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *users) ShowParameters(ctx context.Context, id AccountObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			User: id,
		},
	})
}

// DescribeDetails returns the aggregated UserDetails for a user by calling Describe and converting
// the []UserProperty result. Callers should migrate from Describe to DescribeDetails.
func (v *users) DescribeDetails(ctx context.Context, id AccountObjectIdentifier) (*UserDetails, error) {
	props, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return userDetailsFromRows(props), nil
}

// PAT delegation methods — delegate to client.UserProgrammaticAccessTokens.

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
	return v.client.UserProgrammaticAccessTokens.RemoveByIDSafely(ctx, request)
}

func (v *users) ShowProgrammaticAccessTokens(ctx context.Context, request *ShowUserProgrammaticAccessTokenRequest) ([]ProgrammaticAccessToken, error) {
	return v.client.UserProgrammaticAccessTokens.Show(ctx, request)
}

func (v *users) ShowProgrammaticAccessTokenByName(ctx context.Context, userId AccountObjectIdentifier, tokenName AccountObjectIdentifier) (*ProgrammaticAccessToken, error) {
	return v.client.UserProgrammaticAccessTokens.ShowByID(ctx, userId, tokenName)
}

func (v *users) ShowProgrammaticAccessTokenByNameSafely(ctx context.Context, userId AccountObjectIdentifier, tokenName AccountObjectIdentifier) (*ProgrammaticAccessToken, error) {
	return v.client.UserProgrammaticAccessTokens.ShowByIDSafely(ctx, userId, tokenName)
}

// additionalConvert handles fields that require custom parsing beyond the standard generator output.
// Called by the generated convert() method when any field has WithManualConvert().
func (r *userDBRow) additionalConvert(result *User) error {
	if err := handleNullableBoolString(r.Disabled, &result.Disabled); err != nil {
		return fmt.Errorf("error parsing disabled: %w", err)
	}
	if err := handleNullableBoolString(r.MustChangePassword, &result.MustChangePassword); err != nil {
		return fmt.Errorf("error parsing must change password: %w", err)
	}
	if err := handleNullableBoolString(r.SnowflakeLock, &result.SnowflakeLock); err != nil {
		return fmt.Errorf("error parsing snowflake lock: %w", err)
	}
	if err := handleNullableBoolString(r.ExtAuthnDuo, &result.ExtAuthnDuo); err != nil {
		return fmt.Errorf("error parsing ext authn duo: %w", err)
	}
	return nil
}

// additionalValidations for UserSet — handles the cross-field policy/property exclusion check
// that cannot be expressed as a single generator validation type.
func (opts *UserSet) additionalValidations() error {
	if anyValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy) && anyValueSet(opts.ObjectProperties, opts.ObjectParameters, opts.SessionParameters) {
		return NewError("policies cannot be set with user properties or parameters at the same time")
	}
	return nil
}

// additionalValidations for UserUnset — handles the cross-field policy/property exclusion check
// that cannot be expressed as a single generator validation type.
// TODO [SNOW-1645875]: change validations with policies
func (opts *UserUnset) additionalValidations() error {
	if anyValueSet(opts.PasswordPolicy, opts.SessionPolicy, opts.AuthenticationPolicy) && anyValueSet(opts.ObjectProperties, opts.ObjectParameters, opts.SessionParameters) {
		return NewError("policies cannot be unset with user properties or parameters at the same time")
	}
	return nil
}
