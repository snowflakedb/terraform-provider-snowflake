package sdk

import (
	"context"
	"time"
)

type AuthenticationPolicies interface {
	Create(ctx context.Context, request *CreateAuthenticationPolicyRequest) error
	Alter(ctx context.Context, request *AlterAuthenticationPolicyRequest) error
	Drop(ctx context.Context, request *DropAuthenticationPolicyRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, request *ShowAuthenticationPolicyRequest) ([]AuthenticationPolicy, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*AuthenticationPolicy, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*AuthenticationPolicy, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) ([]AuthenticationPolicyDescription, error)
}

// CreateAuthenticationPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-authentication-policy.
type CreateAuthenticationPolicyOptions struct {
	create                   bool                                        `ddl:"static" sql:"CREATE"`
	OrReplace                *bool                                       `ddl:"keyword" sql:"OR REPLACE"`
	authenticationPolicy     bool                                        `ddl:"static" sql:"AUTHENTICATION POLICY"`
	IfNotExists              *bool                                       `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                     SchemaObjectIdentifier                      `ddl:"identifier"`
	AuthenticationMethods    []AuthenticationMethods                     `ddl:"parameter,parentheses" sql:"AUTHENTICATION_METHODS"`
	MfaAuthenticationMethods []MfaAuthenticationMethods                  `ddl:"parameter,parentheses" sql:"MFA_AUTHENTICATION_METHODS"`
	MfaEnrollment            *MfaEnrollmentOption                        `ddl:"parameter" sql:"MFA_ENROLLMENT"`
	MfaPolicy                *AuthenticationPolicyMfaPolicy              `ddl:"list,parentheses,no_comma" sql:"MFA_POLICY ="`
	ClientTypes              []ClientTypes                               `ddl:"parameter,parentheses" sql:"CLIENT_TYPES"`
	SecurityIntegrations     *SecurityIntegrationsOption                 `ddl:"parameter" sql:"SECURITY_INTEGRATIONS"`
	PatPolicy                *AuthenticationPolicyPatPolicy              `ddl:"list,parentheses,no_comma" sql:"PAT_POLICY ="`
	WorkloadIdentityPolicy   *AuthenticationPolicyWorkloadIdentityPolicy `ddl:"list,parentheses,no_comma" sql:"WORKLOAD_IDENTITY_POLICY ="`
	Comment                  *string                                     `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type AuthenticationMethods struct {
	Method AuthenticationMethodsOption `ddl:"keyword,single_quotes"`
}

type MfaAuthenticationMethods struct {
	Method MfaAuthenticationMethodsOption `ddl:"keyword,single_quotes"`
}

type ClientTypes struct {
	ClientType ClientTypesOption `ddl:"keyword,single_quotes"`
}

type SecurityIntegrationsOption struct {
	All                  *bool                     `ddl:"keyword" sql:"('ALL')"`
	SecurityIntegrations []AccountObjectIdentifier `ddl:"keyword,parentheses"`
}

type AuthenticationPolicyMfaPolicyListItem struct {
	Method MfaPolicyAllowedMethodsOption `ddl:"keyword,single_quotes"`
}

type AuthenticationPolicyMfaPolicy struct {
	EnforceMfaOnExternalAuthentication *EnforceMfaOnExternalAuthenticationOption `ddl:"parameter" sql:"ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION"`
	AllowedMethods                     []AuthenticationPolicyMfaPolicyListItem   `ddl:"parameter,parentheses" sql:"ALLOWED_METHODS"`
}

type AuthenticationPolicyAllowedProviderListItem struct {
	Provider AllowedProviderOption `ddl:"keyword,single_quotes"`
}

type AuthenticationPolicyWorkloadIdentityPolicy struct {
	AllowedProviders    []AuthenticationPolicyAllowedProviderListItem `ddl:"parameter,parentheses" sql:"ALLOWED_PROVIDERS"`
	AllowedAwsAccounts  []StringListItemWrapper                       `ddl:"parameter,parentheses" sql:"ALLOWED_AWS_ACCOUNTS"`
	AllowedAzureIssuers []StringListItemWrapper                       `ddl:"parameter,parentheses" sql:"ALLOWED_AZURE_ISSUERS"`
	AllowedOidcIssuers  []StringListItemWrapper                       `ddl:"parameter,parentheses" sql:"ALLOWED_OIDC_ISSUERS"`
}

type AuthenticationPolicyPatPolicy struct {
	DefaultExpiryInDays     *int                           `ddl:"parameter,no_quotes" sql:"DEFAULT_EXPIRY_IN_DAYS"`
	MaxExpiryInDays         *int                           `ddl:"parameter,no_quotes" sql:"MAX_EXPIRY_IN_DAYS"`
	NetworkPolicyEvaluation *NetworkPolicyEvaluationOption `ddl:"parameter,no_quotes" sql:"NETWORK_POLICY_EVALUATION"`
}

// AlterAuthenticationPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-authentication-policy.
type AlterAuthenticationPolicyOptions struct {
	alter                bool                       `ddl:"static" sql:"ALTER"`
	authenticationPolicy bool                       `ddl:"static" sql:"AUTHENTICATION POLICY"`
	IfExists             *bool                      `ddl:"keyword" sql:"IF EXISTS"`
	name                 SchemaObjectIdentifier     `ddl:"identifier"`
	Set                  *AuthenticationPolicySet   `ddl:"keyword" sql:"SET"`
	Unset                *AuthenticationPolicyUnset `ddl:"list,no_parentheses" sql:"UNSET"`
	RenameTo             *SchemaObjectIdentifier    `ddl:"identifier" sql:"RENAME TO"`
}

type AuthenticationPolicySet struct {
	AuthenticationMethods    []AuthenticationMethods                     `ddl:"parameter,parentheses" sql:"AUTHENTICATION_METHODS"`
	MfaAuthenticationMethods []MfaAuthenticationMethods                  `ddl:"parameter,parentheses" sql:"MFA_AUTHENTICATION_METHODS"`
	MfaEnrollment            *MfaEnrollmentOption                        `ddl:"parameter" sql:"MFA_ENROLLMENT"`
	MfaPolicy                *AuthenticationPolicyMfaPolicy              `ddl:"list,parentheses,no_comma" sql:"MFA_POLICY ="`
	ClientTypes              []ClientTypes                               `ddl:"parameter,parentheses" sql:"CLIENT_TYPES"`
	SecurityIntegrations     *SecurityIntegrationsOption                 `ddl:"parameter" sql:"SECURITY_INTEGRATIONS"`
	PatPolicy                *AuthenticationPolicyPatPolicy              `ddl:"list,parentheses,no_comma" sql:"PAT_POLICY ="`
	WorkloadIdentityPolicy   *AuthenticationPolicyWorkloadIdentityPolicy `ddl:"list,parentheses,no_comma" sql:"WORKLOAD_IDENTITY_POLICY ="`
	Comment                  *string                                     `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type AuthenticationPolicyUnset struct {
	ClientTypes              *bool `ddl:"keyword" sql:"CLIENT_TYPES"`
	AuthenticationMethods    *bool `ddl:"keyword" sql:"AUTHENTICATION_METHODS"`
	SecurityIntegrations     *bool `ddl:"keyword" sql:"SECURITY_INTEGRATIONS"`
	MfaAuthenticationMethods *bool `ddl:"keyword" sql:"MFA_AUTHENTICATION_METHODS"`
	MfaEnrollment            *bool `ddl:"keyword" sql:"MFA_ENROLLMENT"`
	MfaPolicy                *bool `ddl:"keyword" sql:"MFA_POLICY"`
	PatPolicy                *bool `ddl:"keyword" sql:"PAT_POLICY"`
	WorkloadIdentityPolicy   *bool `ddl:"keyword" sql:"WORKLOAD_IDENTITY_POLICY"`
	Comment                  *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropAuthenticationPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-authentication-policy.
type DropAuthenticationPolicyOptions struct {
	drop                 bool                   `ddl:"static" sql:"DROP"`
	authenticationPolicy bool                   `ddl:"static" sql:"AUTHENTICATION POLICY"`
	IfExists             *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name                 SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowAuthenticationPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-authentication-policies.
type ShowAuthenticationPolicyOptions struct {
	show                   bool        `ddl:"static" sql:"SHOW"`
	authenticationPolicies bool        `ddl:"static" sql:"AUTHENTICATION POLICIES"`
	Like                   *Like       `ddl:"keyword" sql:"LIKE"`
	In                     *ExtendedIn `ddl:"keyword" sql:"IN"`
	On                     *On         `ddl:"keyword" sql:"ON"`
	StartsWith             *string     `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit                  *LimitFrom  `ddl:"keyword" sql:"LIMIT"`
}

type showAuthenticationPolicyDBRow struct {
	CreatedOn     time.Time `db:"created_on"`
	Name          string    `db:"name"`
	Comment       string    `db:"comment"`
	DatabaseName  string    `db:"database_name"`
	SchemaName    string    `db:"schema_name"`
	Kind          string    `db:"kind"`
	Owner         string    `db:"owner"`
	OwnerRoleType string    `db:"owner_role_type"`
	Options       string    `db:"options"`
}

type AuthenticationPolicy struct {
	CreatedOn     time.Time
	Name          string
	Comment       string
	DatabaseName  string
	SchemaName    string
	Kind          string
	Owner         string
	OwnerRoleType string
	Options       string
}

func (v *AuthenticationPolicy) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *AuthenticationPolicy) ObjectType() ObjectType {
	return ObjectTypeAuthenticationPolicy
}

func (v *AuthenticationPolicy) ObjectType() ObjectType {
	return ObjectTypeAuthenticationPolicy
}

// DescribeAuthenticationPolicyOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-authentication-policy.
type DescribeAuthenticationPolicyOptions struct {
	describe             bool                   `ddl:"static" sql:"DESCRIBE"`
	authenticationPolicy bool                   `ddl:"static" sql:"AUTHENTICATION POLICY"`
	name                 SchemaObjectIdentifier `ddl:"identifier"`
}

type describeAuthenticationPolicyDBRow struct {
	Property    string `db:"property"`
	Value       string `db:"value"`
	Default     string `db:"default"`
	Description string `db:"description"`
}

type AuthenticationPolicyDescription struct {
	Property    string
	Value       string
	Default     string
	Description string
}
