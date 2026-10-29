package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	id := authenticationPoliciesTestIdSchemaObjectIdentifier
	renameTarget := randomSchemaObjectIdentifier()
	securityIntegrationId := NewAccountObjectIdentifier("security_integration")
	userName := NewAccountObjectIdentifier("user_name")

	mfaPolicy := &AuthenticationPolicyMfaPolicy{
		EnforceMfaOnExternalAuthentication: new(EnforceMfaOnExternalAuthenticationOptionAll),
		AllowedMethods: []MfaPolicyAllowedMethodsOptionListItem{
			{Value: MfaPolicyAllowedMethodsOptionPasskey},
		},
	}
	patPolicy := &AuthenticationPolicyPatPolicy{
		DefaultExpiryInDays:                   new(30),
		MaxExpiryInDays:                       new(90),
		RequireRoleRestrictionForServiceUsers: new(true),
		NetworkPolicyEvaluation:               new(NetworkPolicyEvaluationOptionEnforcedRequired),
	}
	workloadIdentityPolicy := &AuthenticationPolicyWorkloadIdentityPolicy{
		AllowedProviders:    []AllowedProviderOptionListItem{{Value: AllowedProviderOptionAll}},
		AllowedAwsAccounts:  []StringListItemWrapper{{Value: "1234567890"}},
		AllowedAzureIssuers: []StringListItemWrapper{{Value: "https://login.microsoftonline.com/1234567890/v2.0"}},
		AllowedOidcIssuers:  []StringListItemWrapper{{Value: "https://oidc.example.com"}},
	}
	clientTypes := []ClientTypesOptionListItem{
		{Value: ClientTypesOptionDrivers},
		{Value: ClientTypesOptionSnowsql},
	}
	clientPolicy := []AuthenticationPolicyClientPolicy{
		{ClientType: ClientPolicyDriverTypeGoDriver, Params: AuthenticationPolicyClientPolicyParams{MinimumVersion: "1.14.1"}},
		{ClientType: ClientPolicyDriverTypeJdbcDriver, Params: AuthenticationPolicyClientPolicyParams{MinimumVersion: "3.25.0"}},
	}
	securityIntegrations := &AuthenticationPolicySecurityIntegrationsOption{
		SecurityIntegrations: []AccountObjectIdentifier{securityIntegrationId},
	}
	createAllSql := "CREATE AUTHENTICATION POLICY IF NOT EXISTS %s AUTHENTICATION_METHODS = ('SAML', 'PASSWORD')" +
		" MFA_ENROLLMENT = OPTIONAL MFA_POLICY = (ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION = ALL ALLOWED_METHODS = ('PASSKEY'))" +
		" CLIENT_TYPES = ('DRIVERS', 'SNOWSQL') CLIENT_POLICY = (GO_DRIVER = (MINIMUM_VERSION = '1.14.1'), JDBC_DRIVER = (MINIMUM_VERSION = '3.25.0'))" +
		" SECURITY_INTEGRATIONS = (\"security_integration\") PAT_POLICY = (DEFAULT_EXPIRY_IN_DAYS = 30 MAX_EXPIRY_IN_DAYS = 90 REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS = true NETWORK_POLICY_EVALUATION = ENFORCED_REQUIRED)" +
		" WORKLOAD_IDENTITY_POLICY = (ALLOWED_PROVIDERS = ('ALL') ALLOWED_AWS_ACCOUNTS = ('1234567890') ALLOWED_AZURE_ISSUERS = ('https://login.microsoftonline.com/1234567890/v2.0')" +
		" ALLOWED_OIDC_ISSUERS = ('https://oidc.example.com')) COMMENT = 'some comment'"
	alterSetAllSql := "ALTER AUTHENTICATION POLICY IF EXISTS %s SET AUTHENTICATION_METHODS = ('SAML')" +
		" MFA_ENROLLMENT = OPTIONAL MFA_POLICY = (ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION = ALL ALLOWED_METHODS = ('PASSKEY')) CLIENT_TYPES = ('DRIVERS', 'SNOWSQL')" +
		" CLIENT_POLICY = (GO_DRIVER = (MINIMUM_VERSION = '1.14.1'), JDBC_DRIVER = (MINIMUM_VERSION = '3.25.0'))" +
		" SECURITY_INTEGRATIONS = (\"security_integration\") PAT_POLICY = (DEFAULT_EXPIRY_IN_DAYS = 30 MAX_EXPIRY_IN_DAYS = 90 REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS = true NETWORK_POLICY_EVALUATION = ENFORCED_REQUIRED)" +
		" WORKLOAD_IDENTITY_POLICY = (ALLOWED_PROVIDERS = ('ALL') ALLOWED_AWS_ACCOUNTS = ('1234567890') ALLOWED_AZURE_ISSUERS = ('https://login.microsoftonline.com/1234567890/v2.0')" +
		" ALLOWED_OIDC_ISSUERS = ('https://oidc.example.com')) COMMENT = 'some comment'"

	authenticationPoliciesTests.Create.
		withModify(
			case_AuthenticationPolicies_validation_Create_opts_SecurityIntegrations_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateAuthenticationPolicyOptions) {
				opts.SecurityIntegrations = &AuthenticationPolicySecurityIntegrationsOption{
					All:                  new(true),
					SecurityIntegrations: []AccountObjectIdentifier{securityIntegrationId},
				}
			},
		).
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Create_basic,
			func(opts *CreateAuthenticationPolicyOptions) {
				opts.AuthenticationMethods = []AuthenticationMethodsOptionListItem{{Value: AuthenticationMethodsOptionAll}}
				opts.Comment = new("some comment")
			},
			"CREATE AUTHENTICATION POLICY %s AUTHENTICATION_METHODS = ('ALL') COMMENT = 'some comment'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Create_all,
			func(opts *CreateAuthenticationPolicyOptions) {
				opts.IfNotExists = new(true)
				opts.AuthenticationMethods = []AuthenticationMethodsOptionListItem{
					{Value: AuthenticationMethodsOptionSaml},
					{Value: AuthenticationMethodsOptionPassword},
				}
				opts.MfaEnrollment = new(MfaEnrollmentOptionOptional)
				opts.MfaPolicy = mfaPolicy
				opts.PatPolicy = patPolicy
				opts.WorkloadIdentityPolicy = workloadIdentityPolicy
				opts.ClientTypes = clientTypes
				opts.ClientPolicy = clientPolicy
				opts.SecurityIntegrations = securityIntegrations
				opts.Comment = new("some comment")
			},
			createAllSql, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_securityIntegrations_all",
			func(opts *CreateAuthenticationPolicyOptions) {
				opts.SecurityIntegrations = &AuthenticationPolicySecurityIntegrationsOption{All: new(true)}
			},
			"CREATE AUTHENTICATION POLICY %s SECURITY_INTEGRATIONS = ('ALL')", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateAuthenticationPolicyOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE AUTHENTICATION POLICY %s", id.FullyQualifiedName(),
		)

	authenticationPoliciesTests.Alter.
		withModify(
			case_AuthenticationPolicies_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterAuthenticationPolicyOptions) {
				opts.Set = &AuthenticationPolicySet{Comment: new("some comment")}
				opts.Unset = &AuthenticationPolicyUnset{Comment: new(true)}
			},
		).
		withModify(
			case_AuthenticationPolicies_validation_Alter_opts_Set_SecurityIntegrations_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterAuthenticationPolicyOptions) {
				opts.Set = &AuthenticationPolicySet{
					SecurityIntegrations: &AuthenticationPolicySecurityIntegrationsOption{
						All:                  new(true),
						SecurityIntegrations: []AccountObjectIdentifier{securityIntegrationId},
					},
				}
			},
		).
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Alter_Set,
			func(opts *AlterAuthenticationPolicyOptions) {
				opts.Set = &AuthenticationPolicySet{
					AuthenticationMethods: []AuthenticationMethodsOptionListItem{{Value: AuthenticationMethodsOptionSaml}},
				}
			},
			"ALTER AUTHENTICATION POLICY %s SET AUTHENTICATION_METHODS = ('SAML')", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_all",
			func(opts *AlterAuthenticationPolicyOptions) {
				opts.IfExists = new(true)
				opts.Set = &AuthenticationPolicySet{
					AuthenticationMethods: []AuthenticationMethodsOptionListItem{
						{Value: AuthenticationMethodsOptionSaml},
					},
					MfaEnrollment:          new(MfaEnrollmentOptionOptional),
					MfaPolicy:              mfaPolicy,
					PatPolicy:              patPolicy,
					WorkloadIdentityPolicy: workloadIdentityPolicy,
					ClientTypes:            clientTypes,
					ClientPolicy:           clientPolicy,
					SecurityIntegrations:   securityIntegrations,
					Comment:                new("some comment"),
				}
			},
			alterSetAllSql, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Alter_Unset,
			func(opts *AlterAuthenticationPolicyOptions) {
				opts.Unset = &AuthenticationPolicyUnset{ClientTypes: new(true)}
			},
			"ALTER AUTHENTICATION POLICY %s UNSET CLIENT_TYPES", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_all",
			func(opts *AlterAuthenticationPolicyOptions) {
				opts.IfExists = new(true)
				opts.Unset = &AuthenticationPolicyUnset{
					ClientTypes:            new(true),
					ClientPolicy:           new(true),
					AuthenticationMethods:  new(true),
					SecurityIntegrations:   new(true),
					MfaEnrollment:          new(true),
					MfaPolicy:              new(true),
					PatPolicy:              new(true),
					WorkloadIdentityPolicy: new(true),
					Comment:                new(true),
				}
			},
			"ALTER AUTHENTICATION POLICY IF EXISTS %s UNSET CLIENT_TYPES, CLIENT_POLICY, AUTHENTICATION_METHODS, SECURITY_INTEGRATIONS, MFA_ENROLLMENT, MFA_POLICY, PAT_POLICY, WORKLOAD_IDENTITY_POLICY, COMMENT",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Alter_RenameTo,
			func(opts *AlterAuthenticationPolicyOptions) { opts.RenameTo = &renameTarget },
			"ALTER AUTHENTICATION POLICY %s RENAME TO %s", id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		)

	authenticationPoliciesTests.Drop.
		withExpectedSqlf(
			case_AuthenticationPolicies_sql_Drop_basic,
			"DROP AUTHENTICATION POLICY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Drop_all,
			func(opts *DropAuthenticationPolicyOptions) { opts.IfExists = new(true) },
			"DROP AUTHENTICATION POLICY IF EXISTS %s", id.FullyQualifiedName(),
		)

	authenticationPoliciesTests.Show.
		withExpectedSql(case_AuthenticationPolicies_sql_Show_basic, "SHOW AUTHENTICATION POLICIES").
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Show_all,
			func(opts *ShowAuthenticationPolicyOptions) {
				opts.Like = &Like{Pattern: new("like-pattern")}
				opts.In = &ExtendedIn{In: In{Schema: id.SchemaId()}}
				opts.StartsWith = new("starts-with-pattern")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("limit-from")}
			},
			"SHOW AUTHENTICATION POLICIES LIKE 'like-pattern' IN SCHEMA %s STARTS WITH 'starts-with-pattern' LIMIT 10 FROM 'limit-from'",
			id.SchemaId().FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Show_Like,
			func(opts *ShowAuthenticationPolicyOptions) { opts.Like = &Like{Pattern: new("like-pattern")} },
			"SHOW AUTHENTICATION POLICIES LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Show_In,
			func(opts *ShowAuthenticationPolicyOptions) {
				opts.In = &ExtendedIn{In: In{Schema: id.SchemaId()}}
			},
			"SHOW AUTHENTICATION POLICIES IN SCHEMA %s", id.SchemaId().FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Show_On,
			func(opts *ShowAuthenticationPolicyOptions) { opts.On = &On{Account: new(true)} },
			"SHOW AUTHENTICATION POLICIES ON ACCOUNT",
		).
		withAdditionalSqlCasef(
			"sql_Show_On_user",
			func(opts *ShowAuthenticationPolicyOptions) { opts.On = &On{User: userName} },
			`SHOW AUTHENTICATION POLICIES ON USER "user_name"`,
		).
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Show_StartsWith,
			func(opts *ShowAuthenticationPolicyOptions) { opts.StartsWith = new("starts-with-pattern") },
			"SHOW AUTHENTICATION POLICIES STARTS WITH 'starts-with-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_AuthenticationPolicies_sql_Show_Limit,
			func(opts *ShowAuthenticationPolicyOptions) {
				opts.Limit = &LimitFrom{Rows: new(10), From: new("limit-from")}
			},
			"SHOW AUTHENTICATION POLICIES LIMIT 10 FROM 'limit-from'",
		)

	authenticationPoliciesTests.Describe.
		withExpectedSqlf(
			case_AuthenticationPolicies_sql_Describe_basic,
			"DESCRIBE AUTHENTICATION POLICY %s", id.FullyQualifiedName(),
		)
}

func Test_parseAuthenticationPolicyTargetScopes(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected []AuthenticationPolicyTargetScope
	}{
		{
			Name:     "empty options returns nil",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "missing target_scopes key returns nil",
			Input:    `{"other_field":"value"}`,
			Expected: nil,
		},
		{
			Name:     "empty target_scopes array returns empty slice",
			Input:    `{"target_scopes":[]}`,
			Expected: []AuthenticationPolicyTargetScope{},
		},
		{
			Name:     "single scope",
			Input:    `{"target_scopes":["ACCOUNT"]}`,
			Expected: []AuthenticationPolicyTargetScope{AuthenticationPolicyTargetScopeAccount},
		},
		{
			Name:  "multiple scopes are sorted alphabetically",
			Input: `{"target_scopes":["SERVICE_USERS","ACCOUNT","PERSON_USERS"]}`,
			Expected: []AuthenticationPolicyTargetScope{
				AuthenticationPolicyTargetScopeAccount,
				AuthenticationPolicyTargetScopePersonUsers,
				AuthenticationPolicyTargetScopeServiceUsers,
			},
		},
		{
			Name:  "extra attributes are ignored",
			Input: `{"MFA_ENROLLMENT_REQUIREMENT":"REQUIRED_SNOWFLAKE_UI_PASSWORD_ONLY","SYSTEM_MFA_REQUIRED_CLIENT_TYPES":"ALL","target_scopes":["PERSON_USERS"]}`,
			Expected: []AuthenticationPolicyTargetScope{
				AuthenticationPolicyTargetScopePersonUsers,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := parseAuthenticationPolicyTargetScopes(tc.Input)
			require.NoError(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}

	t.Run("invalid json returns error", func(t *testing.T) {
		result, err := parseAuthenticationPolicyTargetScopes(`{"target_scopes":`)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestParseAuthenticationPolicyDetails(t *testing.T) {
	id := NewSchemaObjectIdentifier("DB", "SCHEMA", "POLICY")

	t.Run("defaults", func(t *testing.T) {
		rows := []AuthenticationPolicyProperty{
			{Property: "NAME", Value: "POLICY"},
			{Property: "OWNER", Value: "ACCOUNTADMIN"},
			{Property: "COMMENT", Value: "null"},
			{Property: "MFA_ENROLLMENT", Value: "OPTIONAL"},
			{Property: "SECURITY_INTEGRATIONS", Value: "[ALL]"},
			{Property: "CLIENT_TYPES", Value: "[ALL]"},
			{Property: "CLIENT_POLICY", Value: "{}"},
			{Property: "AUTHENTICATION_METHODS", Value: "[ALL]"},
			{Property: "MFA_POLICY", Value: "{ALLOWED_METHODS=[ALL], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=NONE}"},
			{Property: "PAT_POLICY", Value: "{DEFAULT_EXPIRY_IN_DAYS=15, MAX_EXPIRY_IN_DAYS=365, NETWORK_POLICY_EVALUATION=ENFORCED_REQUIRED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=true, REQUIRE_ROLE_RESTRICTION_FOR_PERSON_USERS=false, BLOCKED_ROLES_LIST=[]}"},
			{Property: "WORKLOAD_IDENTITY_POLICY", Value: "{ALLOWED_PROVIDERS=[ALL], ALLOWED_AWS_ACCOUNTS=[ALL], ALLOWED_AWS_PARTITIONS=[ALL], ALLOWED_AZURE_ISSUERS=[ALL], ALLOWED_OIDC_ISSUERS=[ALL]}"},
		}
		parsed, err := ParseAuthenticationPolicyDescribed(rows)
		details := &parsed
		details.Id = id
		require.NoError(t, err)
		requireAuthenticationPolicyDetailsEqual(t, rows, &AuthenticationPolicyDetails{
			Id:                    id,
			Name:                  "POLICY",
			Owner:                 "ACCOUNTADMIN",
			AuthenticationMethods: []AuthenticationMethodsOption{AuthenticationMethodsOptionAll},
			MfaEnrollment:         MfaEnrollmentReadOptionOptional,
			ClientTypes:           []ClientTypesOption{ClientTypesOptionAll},
			SecurityIntegrations:  AuthenticationPolicySecurityIntegrationsDetails{All: true},
			ClientPolicy:          []AuthenticationPolicyClientPolicyDetails{},
			MfaPolicy: AuthenticationPolicyMfaPolicyDetails{
				AllowedMethods:                     []MfaPolicyAllowedMethodsOption{MfaPolicyAllowedMethodsOptionAll},
				EnforceMfaOnExternalAuthentication: EnforceMfaOnExternalAuthenticationOptionNone,
			},
			PatPolicy: AuthenticationPolicyPatPolicyDetails{
				DefaultExpiryInDays:                   15,
				MaxExpiryInDays:                       365,
				NetworkPolicyEvaluation:               NetworkPolicyEvaluationOptionEnforcedRequired,
				RequireRoleRestrictionForServiceUsers: true,
			},
			WorkloadIdentityPolicy: AuthenticationPolicyWorkloadIdentityPolicyDetails{
				AllowedProviders:    []AllowedProviderOption{AllowedProviderOptionAll},
				AllowedAwsAccounts:  []string{"ALL"},
				AllowedAzureIssuers: []string{"ALL"},
				AllowedOidcIssuers:  []string{"ALL"},
			},
		}, details)
	})

	t.Run("configured values", func(t *testing.T) {
		rows := []AuthenticationPolicyProperty{
			{Property: "NAME", Value: "POLICY"},
			{Property: "OWNER", Value: "ACCOUNTADMIN"},
			{Property: "COMMENT", Value: "some comment"},
			{Property: "MFA_ENROLLMENT", Value: "REQUIRED_SNOWFLAKE_UI_PASSWORD_ONLY"},
			{Property: "SECURITY_INTEGRATIONS", Value: `["My Integration"]`},
			{Property: "CLIENT_TYPES", Value: "[DRIVERS, SNOWSQL]"},
			{Property: "CLIENT_POLICY", Value: "{JDBC_DRIVER={MINIMUM_VERSION=3.25.0}, GO_DRIVER={MINIMUM_VERSION=1.14.1}}"},
			{Property: "AUTHENTICATION_METHODS", Value: "[PASSWORD, SAML]"},
			{Property: "MFA_POLICY", Value: "{ALLOWED_METHODS=[PASSKEY, DUO], ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=ALL}"},
			{Property: "PAT_POLICY", Value: "{DEFAULT_EXPIRY_IN_DAYS=1, MAX_EXPIRY_IN_DAYS=30, NETWORK_POLICY_EVALUATION=NOT_ENFORCED, REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS=false}"},
			{Property: "WORKLOAD_IDENTITY_POLICY", Value: "{ALLOWED_PROVIDERS=[AWS, AZURE], ALLOWED_AWS_ACCOUNTS=[111122223333, 444455556666], ALLOWED_AZURE_ISSUERS=[https://login.microsoftonline.com/tenantid/v2.0], ALLOWED_OIDC_ISSUERS=[https://example.com, https://other.example]}"},
		}
		parsed, err := ParseAuthenticationPolicyDescribed(rows)
		details := &parsed
		details.Id = id
		require.NoError(t, err)
		requireAuthenticationPolicyDetailsEqual(t, rows, &AuthenticationPolicyDetails{
			Id:                    id,
			Name:                  "POLICY",
			Owner:                 "ACCOUNTADMIN",
			AuthenticationMethods: []AuthenticationMethodsOption{AuthenticationMethodsOptionPassword, AuthenticationMethodsOptionSaml},
			MfaEnrollment:         MfaEnrollmentReadOptionRequiredSnowflakeUiPasswordOnly,
			ClientTypes:           []ClientTypesOption{ClientTypesOptionDrivers, ClientTypesOptionSnowsql},
			SecurityIntegrations: AuthenticationPolicySecurityIntegrationsDetails{
				SecurityIntegrations: []AccountObjectIdentifier{NewAccountObjectIdentifier("My Integration")},
			},
			Comment: "some comment",
			ClientPolicy: []AuthenticationPolicyClientPolicyDetails{
				{ClientType: ClientPolicyDriverTypeJdbcDriver, MinimumVersion: "3.25.0"},
				{ClientType: ClientPolicyDriverTypeGoDriver, MinimumVersion: "1.14.1"},
			},
			MfaPolicy: AuthenticationPolicyMfaPolicyDetails{
				AllowedMethods:                     []MfaPolicyAllowedMethodsOption{MfaPolicyAllowedMethodsOptionPasskey, MfaPolicyAllowedMethodsOptionDuo},
				EnforceMfaOnExternalAuthentication: EnforceMfaOnExternalAuthenticationOptionAll,
			},
			PatPolicy: AuthenticationPolicyPatPolicyDetails{
				DefaultExpiryInDays:                   1,
				MaxExpiryInDays:                       30,
				NetworkPolicyEvaluation:               NetworkPolicyEvaluationOptionNotEnforced,
				RequireRoleRestrictionForServiceUsers: false,
			},
			WorkloadIdentityPolicy: AuthenticationPolicyWorkloadIdentityPolicyDetails{
				AllowedProviders:    []AllowedProviderOption{AllowedProviderOptionAws, AllowedProviderOptionAzure},
				AllowedAwsAccounts:  []string{"111122223333", "444455556666"},
				AllowedAzureIssuers: []string{"https://login.microsoftonline.com/tenantid/v2.0"},
				AllowedOidcIssuers:  []string{"https://example.com", "https://other.example"},
			},
		}, details)
	})

	t.Run("invalid authentication method", func(t *testing.T) {
		_, err := ParseAuthenticationPolicyDescribed([]AuthenticationPolicyProperty{
			{Property: "AUTHENTICATION_METHODS", Value: "[NOT_A_METHOD]"},
		})
		require.ErrorContains(t, err, "parsing AUTHENTICATION_METHODS")
	})

	t.Run("invalid pat expiry", func(t *testing.T) {
		_, err := ParseAuthenticationPolicyDescribed([]AuthenticationPolicyProperty{
			{Property: "PAT_POLICY", Value: "{DEFAULT_EXPIRY_IN_DAYS=nope}"},
		})
		require.ErrorContains(t, err, "parsing PAT_POLICY")
	})

	t.Run("all mixed with security integrations", func(t *testing.T) {
		_, err := ParseAuthenticationPolicyDescribed([]AuthenticationPolicyProperty{
			{Property: "SECURITY_INTEGRATIONS", Value: "[ALL, OTHER]"},
		})
		require.ErrorContains(t, err, "parsing SECURITY_INTEGRATIONS")
	})

	t.Run("malformed client policy", func(t *testing.T) {
		_, err := ParseAuthenticationPolicyDescribed([]AuthenticationPolicyProperty{
			{Property: "CLIENT_POLICY", Value: "{GO_DRIVER={MINIMUM_VERSION=1.0.0}"},
		})
		require.ErrorContains(t, err, "parsing CLIENT_POLICY")
	})
}

func requireAuthenticationPolicyDetailsEqual(t *testing.T, rows []AuthenticationPolicyProperty, expected *AuthenticationPolicyDetails, got *AuthenticationPolicyDetails) {
	t.Helper()
	require.Equal(t, AuthenticationPolicyDescriptions(rows), got.Descriptions())
	parsed := *got
	parsed.descriptions = nil
	require.Equal(t, expected, &parsed)
}
