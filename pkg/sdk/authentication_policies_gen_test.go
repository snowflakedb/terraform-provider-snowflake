package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationPolicies_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid CreateAuthenticationPolicyOptions
	defaultOpts := func() *CreateAuthenticationPolicyOptions {
		return &CreateAuthenticationPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateAuthenticationPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateAuthenticationPolicyOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: at least one of the fields [opts.MfaPolicy.EnforceMfaOnExternalAuthentication opts.MfaPolicy.AllowedMethods] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.MfaPolicy = &AuthenticationPolicyMfaPolicy{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("CreateAuthenticationPolicyOptions.MfaPolicy", "EnforceMfaOnExternalAuthentication", "AllowedMethods"))
	})

	t.Run("validation: at least one of the fields [opts.PatPolicy.DefaultExpiryInDays opts.PatPolicy.MaxExpiryInDays opts.PatPolicy.NetworkPolicyEvaluation] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.PatPolicy = &AuthenticationPolicyPatPolicy{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("CreateAuthenticationPolicyOptions.PatPolicy", "DefaultExpiryInDays", "MaxExpiryInDays", "NetworkPolicyEvaluation"))
	})

	t.Run("validation: at least one of the fields [opts.WorkloadIdentityPolicy.AllowedProviders opts.WorkloadIdentityPolicy.AllowedAwsAccounts opts.WorkloadIdentityPolicy.AllowedAzureIssuers opts.WorkloadIdentityPolicy.AllowedOidcIssuers] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.WorkloadIdentityPolicy = &AuthenticationPolicyWorkloadIdentityPolicy{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("CreateAuthenticationPolicyOptions.WorkloadIdentityPolicy", "AllowedProviders", "AllowedAwsAccounts", "AllowedAzureIssuers", "AllowedOidcIssuers"))
	})

	t.Run("validation: exactly one of the fields [opts.SecurityIntegrations.All opts.SecurityIntegrations.SecurityIntegrations] should be set - none set", func(t *testing.T) {
		opts := defaultOpts()
		opts.SecurityIntegrations = &SecurityIntegrationsOption{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateAuthenticationPolicyOptions.SecurityIntegrations", "All", "SecurityIntegrations"))
	})

	t.Run("validation: exactly one of the fields [opts.SecurityIntegrations.All opts.SecurityIntegrations.SecurityIntegrations] should be set - both set", func(t *testing.T) {
		opts := defaultOpts()
		opts.SecurityIntegrations = &SecurityIntegrationsOption{
			All:                  Pointer(true),
			SecurityIntegrations: []AccountObjectIdentifier{NewAccountObjectIdentifier("security_integration")},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateAuthenticationPolicyOptions.SecurityIntegrations", "All", "SecurityIntegrations"))
	})
	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.AuthenticationMethods = []AuthenticationMethods{
			{Method: AuthenticationMethodsAll},
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE AUTHENTICATION POLICY %s AUTHENTICATION_METHODS = ('ALL') COMMENT = 'some comment'", id.FullyQualifiedName())
	})

	t.Run("with security integrations - ALL", func(t *testing.T) {
		opts := defaultOpts()
		opts.SecurityIntegrations = &SecurityIntegrationsOption{
			All: Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE AUTHENTICATION POLICY %s SECURITY_INTEGRATIONS = ('ALL')", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.AuthenticationMethods = []AuthenticationMethods{
			{Method: AuthenticationMethodsSaml},
			{Method: AuthenticationMethodsPassword},
		}
		opts.MfaAuthenticationMethods = []MfaAuthenticationMethods{
			{Method: MfaAuthenticationMethodsPassword},
		}
		opts.MfaEnrollment = Pointer(MfaEnrollmentOptional)
		opts.MfaPolicy = &AuthenticationPolicyMfaPolicy{
			EnforceMfaOnExternalAuthentication: Pointer(EnforceMfaOnExternalAuthenticationAll),
			AllowedMethods: []AuthenticationPolicyMfaPolicyListItem{
				{Method: MfaPolicyPassAllowedMethodPassKey},
			},
		}
		opts.PatPolicy = &AuthenticationPolicyPatPolicy{
			DefaultExpiryInDays:     Int(30),
			MaxExpiryInDays:         Int(90),
			NetworkPolicyEvaluation: Pointer(NetworkPolicyEvaluationEnforcedRequired),
		}
		opts.WorkloadIdentityPolicy = &AuthenticationPolicyWorkloadIdentityPolicy{
			AllowedProviders:    []AuthenticationPolicyAllowedProviderListItem{{Provider: AllowedProviderAll}},
			AllowedAwsAccounts:  []StringListItemWrapper{{Value: "1234567890"}},
			AllowedAzureIssuers: []StringListItemWrapper{{Value: "https://login.microsoftonline.com/1234567890/v2.0"}},
			AllowedOidcIssuers:  []StringListItemWrapper{{Value: "https://oidc.example.com"}},
		}
		opts.ClientTypes = []ClientTypes{
			{ClientType: ClientTypesDrivers},
			{ClientType: ClientTypesSnowSql},
		}
		opts.SecurityIntegrations = &SecurityIntegrationsOption{
			SecurityIntegrations: []AccountObjectIdentifier{
				NewAccountObjectIdentifier("security_integration"),
			},
		}
		opts.Comment = String("some comment")
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE AUTHENTICATION POLICY %s AUTHENTICATION_METHODS = ('SAML', 'PASSWORD')"+
			" MFA_AUTHENTICATION_METHODS = ('PASSWORD') MFA_ENROLLMENT = OPTIONAL MFA_POLICY = (ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION = ALL ALLOWED_METHODS = ('PASSKEY'))"+
			" CLIENT_TYPES = ('DRIVERS', 'SNOWSQL') SECURITY_INTEGRATIONS = (\"security_integration\") PAT_POLICY = (DEFAULT_EXPIRY_IN_DAYS = 30 MAX_EXPIRY_IN_DAYS = 90 NETWORK_POLICY_EVALUATION = ENFORCED_REQUIRED)"+
			" WORKLOAD_IDENTITY_POLICY = (ALLOWED_PROVIDERS = ('ALL') ALLOWED_AWS_ACCOUNTS = ('1234567890') ALLOWED_AZURE_ISSUERS = ('https://login.microsoftonline.com/1234567890/v2.0')"+
			" ALLOWED_OIDC_ISSUERS = ('https://oidc.example.com')) COMMENT = 'some comment'", id.FullyQualifiedName())
	})
}

func TestAuthenticationPolicies_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid AlterAuthenticationPolicyOptions
	defaultOpts := func() *AlterAuthenticationPolicyOptions {
		return &AlterAuthenticationPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterAuthenticationPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.RenameTo] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterAuthenticationPolicyOptions", "Set", "Unset", "RenameTo"))
	})

	t.Run("validation: valid identifier for [opts.RenameTo] if set", func(t *testing.T) {
		opts := defaultOpts()
		opts.RenameTo = &emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.Set.AuthenticationMethods opts.Set.MfaAuthenticationMethods opts.Set.MfaEnrollment opts.Set.ClientTypes opts.Set.SecurityIntegrations opts.Set.Comment opts.Set.MfaPolicy opts.Set.PatPolicy opts.Set.WorkloadIdentityPolicy] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AuthenticationPolicySet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Set", "AuthenticationMethods", "MfaAuthenticationMethods", "MfaEnrollment", "ClientTypes", "SecurityIntegrations", "Comment", "MfaPolicy", "PatPolicy", "WorkloadIdentityPolicy"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.ClientTypes opts.Unset.AuthenticationMethods opts.Unset.Comment opts.Unset.SecurityIntegrations opts.Unset.MfaAuthenticationMethods opts.Unset.MfaEnrollment opts.Unset.MfaPolicy opts.Unset.PatPolicy opts.Unset.WorkloadIdentityPolicy] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &AuthenticationPolicyUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Unset", "ClientTypes", "AuthenticationMethods", "Comment", "SecurityIntegrations", "MfaAuthenticationMethods", "MfaEnrollment", "MfaPolicy", "PatPolicy", "WorkloadIdentityPolicy"))
	})

	t.Run("validation: at least one of the fields [opts.Set.MfaPolicy.EnforceMfaOnExternalAuthentication opts.Set.MfaPolicy.AllowedMethods] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AuthenticationPolicySet{
			MfaPolicy: &AuthenticationPolicyMfaPolicy{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Set.MfaPolicy", "EnforceMfaOnExternalAuthentication", "AllowedMethods"))
	})

	t.Run("validation: at least one of the fields [opts.Set.PatPolicy.DefaultExpiryInDays opts.Set.PatPolicy.MaxExpiryInDays opts.Set.PatPolicy.NetworkPolicyEvaluation] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AuthenticationPolicySet{
			PatPolicy: &AuthenticationPolicyPatPolicy{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Set.PatPolicy", "DefaultExpiryInDays", "MaxExpiryInDays", "NetworkPolicyEvaluation"))
	})

	t.Run("validation: at least one of the fields [opts.Set.WorkloadIdentityPolicy.AllowedProviders opts.Set.WorkloadIdentityPolicy.AllowedAwsAccounts opts.Set.WorkloadIdentityPolicy.AllowedAzureIssuers opts.Set.WorkloadIdentityPolicy.AllowedOidcIssuers] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AuthenticationPolicySet{
			WorkloadIdentityPolicy: &AuthenticationPolicyWorkloadIdentityPolicy{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterAuthenticationPolicyOptions.Set.WorkloadIdentityPolicy", "AllowedProviders", "AllowedAwsAccounts", "AllowedAzureIssuers", "AllowedOidcIssuers"))
	})

	t.Run("validation: exactly one of the fields [opts.Set.SecurityIntegrations.All opts.Set.SecurityIntegrations.SecurityIntegrations] should be set - none set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AuthenticationPolicySet{
			SecurityIntegrations: &SecurityIntegrationsOption{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterAuthenticationPolicyOptions.Set.SecurityIntegrations", "All", "SecurityIntegrations"))
	})

	t.Run("validation: exactly one of the fields [opts.Set.SecurityIntegrations.All opts.Set.SecurityIntegrations.SecurityIntegrations] should be set - both set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AuthenticationPolicySet{
			SecurityIntegrations: &SecurityIntegrationsOption{
				All:                  Pointer(true),
				SecurityIntegrations: []AccountObjectIdentifier{NewAccountObjectIdentifier("security_integration")},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterAuthenticationPolicyOptions.Set.SecurityIntegrations", "All", "SecurityIntegrations"))
	})

	t.Run("alter: set basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &AuthenticationPolicySet{
			AuthenticationMethods: []AuthenticationMethods{
				{Method: AuthenticationMethodsSaml},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER AUTHENTICATION POLICY %s SET AUTHENTICATION_METHODS = ('SAML')", id.FullyQualifiedName())
	})

	t.Run("alter: set all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Set = &AuthenticationPolicySet{
			AuthenticationMethods: []AuthenticationMethods{
				{Method: AuthenticationMethodsSaml},
			},
			MfaAuthenticationMethods: []MfaAuthenticationMethods{
				{Method: MfaAuthenticationMethodsPassword},
			},
			MfaEnrollment: Pointer(MfaEnrollmentOptional),
			MfaPolicy: &AuthenticationPolicyMfaPolicy{
				EnforceMfaOnExternalAuthentication: Pointer(EnforceMfaOnExternalAuthenticationAll),
				AllowedMethods: []AuthenticationPolicyMfaPolicyListItem{
					{Method: MfaPolicyPassAllowedMethodPassKey},
				},
			},
			PatPolicy: &AuthenticationPolicyPatPolicy{
				DefaultExpiryInDays:     Int(30),
				MaxExpiryInDays:         Int(90),
				NetworkPolicyEvaluation: Pointer(NetworkPolicyEvaluationEnforcedRequired),
			},
			WorkloadIdentityPolicy: &AuthenticationPolicyWorkloadIdentityPolicy{
				AllowedProviders:    []AuthenticationPolicyAllowedProviderListItem{{Provider: AllowedProviderAll}},
				AllowedAwsAccounts:  []StringListItemWrapper{{Value: "1234567890"}},
				AllowedAzureIssuers: []StringListItemWrapper{{Value: "https://login.microsoftonline.com/1234567890/v2.0"}},
				AllowedOidcIssuers:  []StringListItemWrapper{{Value: "https://oidc.example.com"}},
			},
			ClientTypes: []ClientTypes{
				{ClientType: ClientTypesDrivers},
				{ClientType: ClientTypesSnowSql},
			},
			SecurityIntegrations: &SecurityIntegrationsOption{
				SecurityIntegrations: []AccountObjectIdentifier{
					NewAccountObjectIdentifier("security_integration"),
				},
			},
			Comment: String("some comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER AUTHENTICATION POLICY IF EXISTS %s SET AUTHENTICATION_METHODS = ('SAML') MFA_AUTHENTICATION_METHODS = ('PASSWORD')"+
			" MFA_ENROLLMENT = OPTIONAL MFA_POLICY = (ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION = ALL ALLOWED_METHODS = ('PASSKEY')) CLIENT_TYPES = ('DRIVERS', 'SNOWSQL')"+
			" SECURITY_INTEGRATIONS = (\"security_integration\") PAT_POLICY = (DEFAULT_EXPIRY_IN_DAYS = 30 MAX_EXPIRY_IN_DAYS = 90 NETWORK_POLICY_EVALUATION = ENFORCED_REQUIRED)"+
			" WORKLOAD_IDENTITY_POLICY = (ALLOWED_PROVIDERS = ('ALL') ALLOWED_AWS_ACCOUNTS = ('1234567890') ALLOWED_AZURE_ISSUERS = ('https://login.microsoftonline.com/1234567890/v2.0')"+
			" ALLOWED_OIDC_ISSUERS = ('https://oidc.example.com')) COMMENT = 'some comment'", id.FullyQualifiedName())
	})

	t.Run("alter: unset basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &AuthenticationPolicyUnset{
			ClientTypes: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER AUTHENTICATION POLICY %s UNSET CLIENT_TYPES", id.FullyQualifiedName())
	})

	t.Run("alter: unset all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Unset = &AuthenticationPolicyUnset{
			ClientTypes:              Bool(true),
			AuthenticationMethods:    Bool(true),
			SecurityIntegrations:     Bool(true),
			MfaAuthenticationMethods: Bool(true),
			MfaEnrollment:            Bool(true),
			MfaPolicy:                Bool(true),
			PatPolicy:                Bool(true),
			WorkloadIdentityPolicy:   Bool(true),
			Comment:                  Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER AUTHENTICATION POLICY IF EXISTS %s UNSET CLIENT_TYPES, AUTHENTICATION_METHODS, SECURITY_INTEGRATIONS, MFA_AUTHENTICATION_METHODS, MFA_ENROLLMENT, MFA_POLICY, PAT_POLICY, WORKLOAD_IDENTITY_POLICY, COMMENT", id.FullyQualifiedName())
	})

	t.Run("alter: renameTo", func(t *testing.T) {
		opts := defaultOpts()
		target := randomSchemaObjectIdentifier()
		opts.RenameTo = &target
		assertOptsValidAndSQLEquals(t, opts, "ALTER AUTHENTICATION POLICY %s RENAME TO %s", id.FullyQualifiedName(), opts.RenameTo.FullyQualifiedName())
	})
}

func TestAuthenticationPolicies_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid DropAuthenticationPolicyOptions
	defaultOpts := func() *DropAuthenticationPolicyOptions {
		return &DropAuthenticationPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropAuthenticationPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP AUTHENTICATION POLICY IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestAuthenticationPolicies_Show(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid ShowAuthenticationPolicyOptions
	defaultOpts := func() *ShowAuthenticationPolicyOptions {
		return &ShowAuthenticationPolicyOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowAuthenticationPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW AUTHENTICATION POLICIES")
	})

	t.Run("show on account", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = &On{
			Account: Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW AUTHENTICATION POLICIES ON ACCOUNT")
	})

	t.Run("show on user", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = &On{
			User: NewAccountObjectIdentifier("user_name"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW AUTHENTICATION POLICIES ON USER "user_name"`)
	})

	t.Run("complete", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("like-pattern"),
		}
		opts.StartsWith = String("starts-with-pattern")
		opts.In = &ExtendedIn{
			In: In{
				Schema: id.SchemaId(),
			},
		}
		opts.Limit = &LimitFrom{
			Rows: Int(10),
			From: String("limit-from"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW AUTHENTICATION POLICIES LIKE 'like-pattern' IN SCHEMA %s STARTS WITH 'starts-with-pattern' LIMIT 10 FROM 'limit-from'", id.SchemaId().FullyQualifiedName())
	})
}

func TestAuthenticationPolicies_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid DescribeAuthenticationPolicyOptions
	defaultOpts := func() *DescribeAuthenticationPolicyOptions {
		return &DescribeAuthenticationPolicyOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeAuthenticationPolicyOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE AUTHENTICATION POLICY %s", id.FullyQualifiedName())
	})
}

func Test_ToAuthenticationMethodsOption(t *testing.T) {
	type test struct {
		input string
		want  AuthenticationMethodsOption
	}
	valid := []test{
		// case insensitive.
		{input: "all", want: AuthenticationMethodsAll},

		// supported values.
		{input: "ALL", want: AuthenticationMethodsAll},
		{input: "SAML", want: AuthenticationMethodsSaml},
		{input: "PASSWORD", want: AuthenticationMethodsPassword},
		{input: "OAUTH", want: AuthenticationMethodsOauth},
		{input: "KEYPAIR", want: AuthenticationMethodsKeyPair},
		{input: "PROGRAMMATIC_ACCESS_TOKEN", want: AuthenticationMethodsProgrammaticAccessToken},
		{input: "WORKLOAD_IDENTITY", want: AuthenticationMethodsWorkloadIdentity},
	}
	invalid := []test{
		{input: "foo"},
	}
	for _, tt := range valid {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ToAuthenticationMethodsOption(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
	for _, tt := range invalid {
		t.Run(tt.input, func(t *testing.T) {
			_, err := ToAuthenticationMethodsOption(tt.input)
			require.Error(t, err)
		})
	}
}

func Test_ToMfaAuthenticationMethodsOption(t *testing.T) {
	type test struct {
		input string
		want  MfaAuthenticationMethodsOption
	}
	valid := []test{
		// case insensitive.
		{input: "all", want: MfaAuthenticationMethodsAll},

		// supported values.
		{input: "ALL", want: MfaAuthenticationMethodsAll},
		{input: "SAML", want: MfaAuthenticationMethodsSaml},
		{input: "PASSWORD", want: MfaAuthenticationMethodsPassword},
	}
	invalid := []test{
		{input: "foo"},
		{input: "OAUTH"},
	}
	for _, tt := range valid {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ToMfaAuthenticationMethodsOption(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
	for _, tt := range invalid {
		t.Run(tt.input, func(t *testing.T) {
			_, err := ToMfaAuthenticationMethodsOption(tt.input)
			require.Error(t, err)
		})
	}
}

func Test_ToMfaEnrollmentOption(t *testing.T) {
	type test struct {
		input string
		want  MfaEnrollmentOption
	}
	valid := []test{
		// case insensitive.
		{input: "required", want: MfaEnrollmentRequired},
		{input: "required_password_only", want: MfaEnrollmentRequiredPasswordOnly},
		{input: "optional", want: MfaEnrollmentOptional},

		// supported values.
		{input: "REQUIRED", want: MfaEnrollmentRequired},
		{input: "REQUIRED_PASSWORD_ONLY", want: MfaEnrollmentRequiredPasswordOnly},
		{input: "OPTIONAL", want: MfaEnrollmentOptional},
	}
	invalid := []test{
		{input: "foo"},
		{input: "ALL"},
	}
	for _, tt := range valid {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ToMfaEnrollmentOption(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
	for _, tt := range invalid {
		t.Run(tt.input, func(t *testing.T) {
			_, err := ToMfaEnrollmentOption(tt.input)
			require.Error(t, err)
		})
	}
}

func Test_ToClientTypesOption(t *testing.T) {
	type test struct {
		input string
		want  ClientTypesOption
	}
	valid := []test{
		// case insensitive.
		{input: "all", want: ClientTypesAll},
		{input: "snowflake_ui", want: ClientTypesSnowflakeUi},
		{input: "drivers", want: ClientTypesDrivers},
		{input: "snowsql", want: ClientTypesSnowSql},
		{input: "snowflake_cli", want: ClientTypesSnowflakeCli},

		// supported values.
		{input: "ALL", want: ClientTypesAll},
		{input: "SNOWFLAKE_UI", want: ClientTypesSnowflakeUi},
		{input: "DRIVERS", want: ClientTypesDrivers},
		{input: "SNOWSQL", want: ClientTypesSnowSql},
		{input: "SNOWFLAKE_CLI", want: ClientTypesSnowflakeCli},
	}
	invalid := []test{
		{input: "foo"},
		{input: "PASSWORD"},
	}
	for _, tt := range valid {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ToClientTypesOption(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
	for _, tt := range invalid {
		t.Run(tt.input, func(t *testing.T) {
			_, err := ToClientTypesOption(tt.input)
			require.Error(t, err)
		})
	}
}

func Test_ToAllowedProviderOption(t *testing.T) {
	type test struct {
		input string
		want  AllowedProviderOption
	}
	valid := []test{
		// case insensitive.
		{input: "all", want: AllowedProviderAll},

		// supported values.
		{input: "ALL", want: AllowedProviderAll},
		{input: "AWS", want: AllowedProviderAws},
		{input: "AZURE", want: AllowedProviderAzure},
		{input: "GCP", want: AllowedProviderGcp},
		{input: "OIDC", want: AllowedProviderOidc},
	}

	invalid := []test{
		{input: "foo"},
	}
	for _, tt := range valid {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ToAllowedProviderOption(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
	for _, tt := range invalid {
		t.Run(tt.input, func(t *testing.T) {
			_, err := ToAllowedProviderOption(tt.input)
			require.Error(t, err)
		})
	}
}

func Test_ToNetworkPolicyEvaluationOption(t *testing.T) {
	type test struct {
		input string
		want  NetworkPolicyEvaluationOption
	}
	valid := []test{
		// case insensitive.
		{input: "enforced_required", want: NetworkPolicyEvaluationEnforcedRequired},

		// supported values.
		{input: "ENFORCED_REQUIRED", want: NetworkPolicyEvaluationEnforcedRequired},
		{input: "ENFORCED_NOT_REQUIRED", want: NetworkPolicyEvaluationEnforcedNotRequired},
		{input: "NOT_ENFORCED", want: NetworkPolicyEvaluationNotEnforced},
	}
	invalid := []test{
		{input: "foo"},
	}
	for _, tt := range valid {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ToNetworkPolicyEvaluationOption(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
	for _, tt := range invalid {
		t.Run(tt.input, func(t *testing.T) {
			_, err := ToNetworkPolicyEvaluationOption(tt.input)
			require.Error(t, err)
		})
	}
}

func Test_ToEnforceMfaOnExternalAuthenticationOption(t *testing.T) {
	type test struct {
		input string
		want  EnforceMfaOnExternalAuthenticationOption
	}
	valid := []test{
		// case insensitive.
		{input: "all", want: EnforceMfaOnExternalAuthenticationAll},

		// supported values.
		{input: "ALL", want: EnforceMfaOnExternalAuthenticationAll},
		{input: "NONE", want: EnforceMfaOnExternalAuthenticationNone},
	}
	invalid := []test{
		{input: "foo"},
	}
	for _, tt := range valid {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ToEnforceMfaOnExternalAuthenticationOption(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
	for _, tt := range invalid {
		t.Run(tt.input, func(t *testing.T) {
			_, err := ToEnforceMfaOnExternalAuthenticationOption(tt.input)
			require.Error(t, err)
		})
	}
}
