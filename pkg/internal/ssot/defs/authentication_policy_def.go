package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var AuthenticationMethodsOption = ssot.NewEnum(
	"AuthenticationMethodsOption", "AuthenticationMethodsOptions",
	"ALL", "SAML", "PASSWORD", "OAUTH", "KEYPAIR", "PROGRAMMATIC_ACCESS_TOKEN", "WORKLOAD_IDENTITY",
)

var MfaEnrollmentOption = ssot.NewEnum(
	"MfaEnrollmentOption", "MfaEnrollmentOptions",
	"REQUIRED", "REQUIRED_PASSWORD_ONLY", "OPTIONAL",
)

var ClientTypesOption = ssot.NewEnum(
	"ClientTypesOption", "ClientTypesOptions",
	"ALL", "SNOWFLAKE_UI", "DRIVERS", "SNOWSQL", "SNOWFLAKE_CLI",
)

var AuthenticationPolicyTargetScope = ssot.NewEnum(
	"AuthenticationPolicyTargetScope", "AuthenticationPolicyTargetScopes",
	"ACCOUNT", "PERSON_USERS", "SERVICE_USERS",
)

var ClientPolicyDriverType = ssot.NewEnum(
	"ClientPolicyDriverType", "ClientPolicyDriverTypes",
	"JDBC_DRIVER", "ODBC_DRIVER", "PYTHON_DRIVER", "JAVASCRIPT_DRIVER",
	"C_DRIVER", "GO_DRIVER", "PHP_DRIVER", "DOTNET_DRIVER",
	"SQL_API", "SNOWPIPE_STREAMING_CLIENT_SDK", "PY_CORE",
	"SPROC_PYTHON", "PYTHON_SNOWPARK", "SQL_ALCHEMY", "SNOWPARK", "SNOWFLAKE_CLIENT",
)

var EnforceMfaOnExternalAuthenticationOption = ssot.NewEnum(
	"EnforceMfaOnExternalAuthenticationOption", "EnforceMfaOnExternalAuthenticationOptions",
	"ALL", "NONE",
)

var MfaPolicyAllowedMethodsOption = ssot.NewEnum(
	"MfaPolicyAllowedMethodsOption", "MfaPolicyAllowedMethodsOptions",
	"ALL", "PASSKEY", "TOTP", "OTP", "DUO",
)

var NetworkPolicyEvaluationOption = ssot.NewEnum(
	"NetworkPolicyEvaluationOption", "NetworkPolicyEvaluationOptions",
	"ENFORCED_REQUIRED", "ENFORCED_NOT_REQUIRED", "NOT_ENFORCED",
)

var AllowedProviderOption = ssot.NewEnum(
	"AllowedProviderOption", "AllowedProviderOptions",
	"ALL", "AWS", "AZURE", "GCP", "OIDC",
)

// AuthenticationPolicyAttrs holds reusable Authentication Policy fields so SHOW and DESCRIBE
// Columns can point at the same definitions instead of restating names.
//
// mfa_authentication_methods stays in the schema ext. It is a deprecated set with no effect,
// not a CREATE/ALTER column. The ext also sets AtLeastOneOf and keeps the case-sensitive
// wording on the nested enum sets.
//
// Sets are created through hand-written per-element mappers named to<Stem><Field>. The
// security integrations option is mapped as a whole, because the SDK takes one option request.
// The nested blocks are omitted from DESCRIBE Columns. A writable column must be in
// Attributes, and the raw DESCRIBE property strings are not a separate attribute.
var AuthenticationPolicyAttrs = struct {
	AuthenticationMethods  ssot.Attribute
	MfaEnrollment          ssot.Attribute
	ClientTypes            ssot.Attribute
	ClientPolicy           ssot.Attribute
	MfaPolicy              ssot.Attribute
	PatPolicy              ssot.Attribute
	WorkloadIdentityPolicy ssot.Attribute
	SecurityIntegrations   ssot.Attribute
	Comment                ssot.Attribute
	CreatedOn              ssot.Attribute
	Name                   ssot.Attribute
	DatabaseName           ssot.Attribute
	SchemaName             ssot.Attribute
	Kind                   ssot.Attribute
	Owner                  ssot.Attribute
	OwnerRoleType          ssot.Attribute
	Options                ssot.Attribute
	TargetScopes           ssot.Attribute
}{
	AuthenticationMethods: ssot.Attribute{
		Name:        "authentication_methods",
		Type:        ssot.AttrSet,
		Description: "A list of authentication methods that are allowed during login.",
		Elem: &ssot.Attribute{
			Type: ssot.AttrEnum,
			Enum: AuthenticationMethodsOption,
		},
	},
	MfaEnrollment: ssot.Attribute{
		Name:        "mfa_enrollment",
		Type:        ssot.AttrEnum,
		Description: "Determines whether a user must enroll in multi-factor authentication.",
		AdditionalDescription: "When REQUIRED is specified, Enforces users to enroll in MFA. " +
			"If this value is used, then the `client_types` parameter must include `snowflake_ui`, because Snowsight is the only place users can enroll in multi-factor authentication (MFA). " +
			"Note that when you set this value to OPTIONAL, and your account setup forces users to enroll in MFA, then Snowflake may set quietly this value to `REQUIRED_PASSWORD_ONLY`, which may cause permadiff. " +
			"In this case, you may want to adjust this field value.",
		Enum: MfaEnrollmentOption,
	},
	ClientTypes: ssot.Attribute{
		Name:        "client_types",
		Type:        ssot.AttrSet,
		Description: "A list of clients that can authenticate with Snowflake. If a client tries to connect, and the client is not one of the valid `client_types`, then the login attempt fails.",
		AdditionalDescription: "The `client_types` property of an authentication policy is a best effort method to block user logins based on specific clients. " +
			"It should not be used as the sole control to establish a security boundary.",
		Elem: &ssot.Attribute{
			Type: ssot.AttrEnum,
			Enum: ClientTypesOption,
		},
	},
	ClientPolicy: ssot.Attribute{
		Name:        "client_policy",
		Type:        ssot.AttrSet,
		Description: "Allows to set policies per-client type. Only valid when `client_types` is empty, contains ALL, or contains DRIVERS.",
		Elem: ssot.ObjectElem([]ssot.Attribute{
			{
				Name:        "client_type",
				Type:        ssot.AttrEnum,
				Required:    true,
				Description: "The client or driver type.",
				Enum:        ClientPolicyDriverType,
			},
			{
				Name:        "minimum_version",
				Type:        ssot.AttrString,
				Required:    true,
				Description: "Minimum allowed version for this client/driver type (e.g. '1.14.1').",
			},
		}),
	},
	MfaPolicy: ssot.Attribute{
		Name:        "mfa_policy",
		Type:        ssot.AttrList,
		MaxItems:    new(1),
		Description: "Specifies the multi-factor authentication (MFA) methods that users can use as a second factor of authentication.",
		Elem: ssot.ObjectElem([]ssot.Attribute{
			{
				Name:        "enforce_mfa_on_external_authentication",
				Type:        ssot.AttrEnum,
				Description: "Determines whether multi-factor authentication (MFA) is enforced on external authentication.",
				Enum:        EnforceMfaOnExternalAuthenticationOption,
			},
			{
				Name:        "allowed_methods",
				Type:        ssot.AttrSet,
				Description: "Specifies the allowed methods for the MFA policy.",
				Elem: &ssot.Attribute{
					Type: ssot.AttrEnum,
					Enum: MfaPolicyAllowedMethodsOption,
				},
			},
		}),
	},
	PatPolicy: ssot.Attribute{
		Name:        "pat_policy",
		Type:        ssot.AttrList,
		MaxItems:    new(1),
		Description: "Specifies the policy for programmatic access tokens.",
		Elem: ssot.ObjectElem([]ssot.Attribute{
			{
				Name:        "default_expiry_in_days",
				Type:        ssot.AttrInt,
				Description: "Specifies the default expiration time (in days) for a programmatic access token.",
				Min:         new(1),
			},
			{
				Name:        "max_expiry_in_days",
				Type:        ssot.AttrInt,
				Description: "Specifies the maximum number of days that can be set for the expiration time for a programmatic access token.",
				Min:         new(1),
			},
			{
				Name:        "require_role_restriction_for_service_users",
				Type:        ssot.AttrBool,
				Description: "If true, when you generate a programmatic access token for a service user, you must restrict the use of that token to a specific role.",
			},
			{
				Name:        "network_policy_evaluation",
				Type:        ssot.AttrEnum,
				Description: "Specifies the network policy evaluation for the PAT.",
				Enum:        NetworkPolicyEvaluationOption,
			},
		}),
	},
	WorkloadIdentityPolicy: ssot.Attribute{
		Name:        "workload_identity_policy",
		Type:        ssot.AttrList,
		MaxItems:    new(1),
		Description: "Specifies the policy for workload identity federation.",
		Elem: ssot.ObjectElem([]ssot.Attribute{
			{
				Name:        "allowed_providers",
				Type:        ssot.AttrSet,
				Description: "Specifies the allowed providers for the workload identity policy.",
				Elem: &ssot.Attribute{
					Type: ssot.AttrEnum,
					Enum: AllowedProviderOption,
				},
			},
			{
				Name:        "allowed_aws_accounts",
				Type:        ssot.AttrSet,
				Description: "Specifies the list of AWS account IDs allowed by the authentication policy during workload identity authentication of type `AWS`.",
				Elem:        &ssot.Attribute{Type: ssot.AttrString},
			},
			{
				Name:        "allowed_azure_issuers",
				Type:        ssot.AttrSet,
				Description: "Specifies the list of Azure Entra ID issuers allowed by the authentication policy during workload identity authentication of type `AZURE`.",
				Elem:        &ssot.Attribute{Type: ssot.AttrString},
			},
			{
				Name:        "allowed_oidc_issuers",
				Type:        ssot.AttrSet,
				Description: "Specifies the list of OIDC issuers allowed by the authentication policy during workload identity authentication of type `OIDC`.",
				Elem:        &ssot.Attribute{Type: ssot.AttrString},
			},
		}),
	},
	SecurityIntegrations: ssot.Attribute{
		Name: "security_integrations",
		Type: ssot.AttrSet,
		Description: "A list of security integrations the authentication policy is associated with. " +
			"This parameter has no effect when `saml` or `oauth` are not in the `authentication_methods` list. " +
			"All values in the `security_integrations` list must be compatible with the values in the `authentication_methods` list. " +
			"For example, if `security_integrations` contains a SAML security integration, and `authentication_methods` contains OAUTH, then you cannot create the authentication policy. " +
			"To allow all security integrations use `ALL` as parameter.",
		Elem: &ssot.Attribute{
			Type:           ssot.AttrIdentifier,
			IdentifierType: "AccountObjectIdentifier",
		},
		CreateMapping: ssot.AttributeMappedValueCreate("ToSecurityIntegrationsRequest"),
	},
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the authentication policy.",
	},
	CreatedOn:     ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, OptionalInDb: true, Readonly: true},
	Name:          ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	DatabaseName:  ssot.Attribute{Name: "database_name", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
	SchemaName:    ssot.Attribute{Name: "schema_name", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
	Kind:          ssot.Attribute{Name: "kind", Type: ssot.AttrString, Readonly: true},
	Owner:         ssot.Attribute{Name: "owner", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
	OwnerRoleType: ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
	Options:       ssot.Attribute{Name: "options", Type: ssot.AttrString, Readonly: true},
	TargetScopes: ssot.Attribute{
		Name:     "target_scopes",
		Type:     ssot.AttrList,
		Readonly: true,
		Elem: &ssot.Attribute{
			Type: ssot.AttrEnum,
			Enum: AuthenticationPolicyTargetScope,
		},
	},
}

var AuthenticationPolicy = ssot.Object{
	Name:                      "AuthenticationPolicy",
	NamePlural:                "AuthenticationPolicies",
	DescribeFailRead:          true,
	IdType:                    "SchemaObjectIdentifier",
	IdentityHierarchyForceNew: true,
	// Public describe_output uses the projected AuthenticationPolicyDescribeDetails schema.
	DescribeDetailsName:    "AuthenticationPolicyDescribeDetails",
	AttributesModification: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithHandleExternalChangesInDescribe(ssot.ExtCall()),
		ssot.WithHierarchyRename(ssot.ExtCall()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithStateUpgraders(ssot.Generated()),
		ssot.WithCustomizeDiff(ssot.Generated()),
		ssot.WithSetConfigFields(ssot.Generated()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
	},
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"v2_9_0_AuthenticationPolicyStateUpgrader",
	},
	Attributes: []ssot.Attribute{
		AuthenticationPolicyAttrs.AuthenticationMethods,
		AuthenticationPolicyAttrs.MfaEnrollment,
		AuthenticationPolicyAttrs.ClientTypes,
		AuthenticationPolicyAttrs.ClientPolicy,
		AuthenticationPolicyAttrs.MfaPolicy,
		AuthenticationPolicyAttrs.PatPolicy,
		AuthenticationPolicyAttrs.WorkloadIdentityPolicy,
		AuthenticationPolicyAttrs.SecurityIntegrations,
		AuthenticationPolicyAttrs.Comment,
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&AuthenticationPolicyAttrs.CreatedOn,
			&AuthenticationPolicyAttrs.Name,
			&AuthenticationPolicyAttrs.Comment,
			&AuthenticationPolicyAttrs.DatabaseName,
			&AuthenticationPolicyAttrs.SchemaName,
			&AuthenticationPolicyAttrs.Kind,
			&AuthenticationPolicyAttrs.Owner,
			&AuthenticationPolicyAttrs.OwnerRoleType,
			&AuthenticationPolicyAttrs.Options,
			&AuthenticationPolicyAttrs.TargetScopes,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&AuthenticationPolicyAttrs.Name,
			&AuthenticationPolicyAttrs.Owner,
			&AuthenticationPolicyAttrs.AuthenticationMethods,
			&AuthenticationPolicyAttrs.MfaEnrollment,
			&AuthenticationPolicyAttrs.ClientTypes,
			&AuthenticationPolicyAttrs.SecurityIntegrations,
			&AuthenticationPolicyAttrs.Comment,
		},
	},
}
