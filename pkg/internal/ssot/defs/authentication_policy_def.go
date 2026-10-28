package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var AuthenticationMethodsOption = ssot.NewEnum(
	"AuthenticationMethodsOption", "AuthenticationMethodsOptions",
	"ALL", "SAML", "PASSWORD", "OAUTH", "KEYPAIR", "PROGRAMMATIC_ACCESS_TOKEN", "WORKLOAD_IDENTITY",
)

var MfaAuthenticationMethodsOption = ssot.NewEnum(
	"MfaAuthenticationMethodsOption", "MfaAuthenticationMethodsOptions",
	"ALL", "SAML", "PASSWORD",
)

var MfaAuthenticationMethodsReadOption = ssot.NewEnum(
	"MfaAuthenticationMethodsReadOption", "MfaAuthenticationMethodsReadOptions",
	"ALL", "SAML", "PASSWORD", "OIDC",
)

var MfaEnrollmentOption = ssot.NewEnum(
	"MfaEnrollmentOption", "MfaEnrollmentOptions",
	"REQUIRED", "REQUIRED_PASSWORD_ONLY", "OPTIONAL",
)

var MfaEnrollmentReadOption = ssot.NewEnum(
	"MfaEnrollmentReadOption", "MfaEnrollmentReadOptions",
	"REQUIRED", "REQUIRED_PASSWORD_ONLY", "OPTIONAL", "REQUIRED_SNOWFLAKE_UI_PASSWORD_ONLY",
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

// The nested blocks are written the same way in CREATE and in ALTER SET.
var (
	authenticationPolicyMfaPolicyClause              = ssot.ObjectParameter(ssot.Sql("MFA_POLICY ="), ssot.Parentheses())
	authenticationPolicyPatPolicyClause              = ssot.ObjectParameter(ssot.Sql("PAT_POLICY ="), ssot.Parentheses())
	authenticationPolicyWorkloadIdentityPolicyClause = ssot.ObjectParameter(ssot.Sql("WORKLOAD_IDENTITY_POLICY ="), ssot.Parentheses())
	authenticationPolicySecurityIntegrationsClause   = ssot.ObjectParameter(ssot.Sql("SECURITY_INTEGRATIONS"), ssot.FieldName("SecurityIntegrations"), ssot.AsParameter())
)

var authenticationPolicyMfaPolicyChildren = []ssot.Attribute{
	{
		Name:         "enforce_mfa_on_external_authentication",
		Type:         ssot.AttrEnum,
		Description:  "Determines whether multi-factor authentication (MFA) is enforced on external authentication.",
		Enum:         EnforceMfaOnExternalAuthenticationOption,
		CreateSql:    ssot.PlainParameter(),
		AtLeastOneOf: authenticationPolicyMfaAtLeastOneOf,
	},
	{
		Name:        "allowed_methods",
		Type:        ssot.AttrSet,
		Description: "Specifies the allowed methods for the MFA policy.",
		Elem: &ssot.Attribute{
			Type: ssot.AttrEnum,
			Enum: MfaPolicyAllowedMethodsOption,
		},
		AtLeastOneOf: authenticationPolicyMfaAtLeastOneOf,
	},
}

var authenticationPolicyPatPolicyChildren = []ssot.Attribute{
	{
		Name:         "default_expiry_in_days",
		Type:         ssot.AttrInt,
		CreateSql:    ssot.IntParameter(ssot.NoQuotes()),
		Description:  "Specifies the default expiration time (in days) for a programmatic access token.",
		Min:          new(1),
		AtLeastOneOf: authenticationPolicyPatAtLeastOneOf,
	},
	{
		Name:         "max_expiry_in_days",
		Type:         ssot.AttrInt,
		CreateSql:    ssot.IntParameter(ssot.NoQuotes()),
		Description:  "Specifies the maximum number of days that can be set for the expiration time for a programmatic access token.",
		Min:          new(1),
		AtLeastOneOf: authenticationPolicyPatAtLeastOneOf,
	},
	{
		Name:         "require_role_restriction_for_service_users",
		Type:         ssot.AttrBool,
		Description:  "If true, when you generate a programmatic access token for a service user, you must restrict the use of that token to a specific role.",
		AtLeastOneOf: authenticationPolicyPatAtLeastOneOf,
	},
	{
		Name:         "network_policy_evaluation",
		Type:         ssot.AttrEnum,
		Description:  "Specifies the network policy evaluation for the PAT.",
		Enum:         NetworkPolicyEvaluationOption,
		AtLeastOneOf: authenticationPolicyPatAtLeastOneOf,
	},
}

var authenticationPolicyWorkloadIdentityPolicyChildren = []ssot.Attribute{
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
		CreateSql:   ssot.ListParameter("StringListItemWrapper"),
	},
	{
		Name:        "allowed_azure_issuers",
		Type:        ssot.AttrSet,
		Description: "Specifies the list of Azure Entra ID issuers allowed by the authentication policy during workload identity authentication of type `AZURE`.",
		Elem:        &ssot.Attribute{Type: ssot.AttrString},
		CreateSql:   ssot.ListParameter("StringListItemWrapper"),
	},
	{
		Name:        "allowed_oidc_issuers",
		Type:        ssot.AttrSet,
		Description: "Specifies the list of OIDC issuers allowed by the authentication policy during workload identity authentication of type `OIDC`.",
		Elem:        &ssot.Attribute{Type: ssot.AttrString},
		CreateSql:   ssot.ListParameter("StringListItemWrapper"),
	},
}

// authenticationPolicyClientPolicyMinimumVersion is a field of a client_policy entry in Terraform. In SQL it sits inside the
// parameters of the entry, see authenticationPolicyClientPolicyParams.
var authenticationPolicyClientPolicyMinimumVersion = ssot.Attribute{
	Name:        "minimum_version",
	Type:        ssot.AttrString,
	Required:    true,
	Description: "Minimum allowed version for this client/driver type (e.g. '1.14.1').",
}

// authenticationPolicyClientPolicyParams is the `= (MINIMUM_VERSION = '...')` of a client_policy entry. It is not a Terraform
// field: the entry field it holds is written only inside it.
var authenticationPolicyClientPolicyParams = ssot.Attribute{
	Name:      "client_policy_params",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.Sql("="), ssot.Parentheses(), ssot.FieldName("Params"), ssot.Required()),
	Children:  []ssot.Attribute{authenticationPolicyClientPolicyMinimumVersion},
}

var authenticationPolicySecurityIntegrations = ssot.Attribute{
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
	CreateSql:     ssot.IdentifierListKeyword(),
	CreateMapping: ssot.AttributeMappedValueCreate("ToSecurityIntegrationsRequest"),
}

// authenticationPolicySecurityIntegrationsOption is the SQL struct of SECURITY_INTEGRATIONS: either ALL or a list of
// identifiers. It is not a Terraform field. The list is the security_integrations attribute, and ALL is the value of the
// attribute that the resource maps to the keyword.
var authenticationPolicySecurityIntegrationsOption = ssot.Attribute{
	Name:      "security_integrations_option",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: authenticationPolicySecurityIntegrationsClause,
	SetSql:    authenticationPolicySecurityIntegrationsClause,
	Children: []ssot.Attribute{
		{Name: "security_integrations_all", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("('ALL')"), ssot.FieldName("All"))},
		authenticationPolicySecurityIntegrations,
	},
	ChildrenRule: ssot.ChildrenExactlyOne,
}

// AuthenticationPolicyAttrs holds reusable Authentication Policy fields so SHOW and DESCRIBE
// Columns can point at the same definitions instead of restating names.
//
// mfa_authentication_methods stays in the schema ext. It is a deprecated set with no effect,
// not a CREATE/ALTER column. Nested enum sets stay case-sensitive: the generator normalizes
// them, and the ext clears that and keeps the case-sensitive wording.
//
// Sets are created through hand-written per-element mappers named to<Stem><Field>. The
// security integrations option is mapped as a whole, because the SDK takes one option request.
// The nested blocks are omitted from DESCRIBE Columns. A writable column must be in
// Attributes, and the raw DESCRIBE property strings are not a separate attribute.
var authenticationPolicyMfaAtLeastOneOf = []ssot.Attribute{
	{Name: "enforce_mfa_on_external_authentication"},
	{Name: "allowed_methods"},
}

var authenticationPolicyPatAtLeastOneOf = []ssot.Attribute{
	{Name: "default_expiry_in_days"},
	{Name: "max_expiry_in_days"},
	{Name: "require_role_restriction_for_service_users"},
	{Name: "network_policy_evaluation"},
}

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
	// DESCRIBE-only columns: a read enum, and the values DESCRIBE returns as nested structs. They share the Children of the attributes.
	DescribeMfaEnrollment          ssot.Attribute
	DescribeSecurityIntegrations   ssot.Attribute
	DescribeClientPolicy           ssot.Attribute
	DescribeMfaPolicy              ssot.Attribute
	DescribePatPolicy              ssot.Attribute
	DescribeWorkloadIdentityPolicy ssot.Attribute
	// The DESCRIBE rows are a property list, which the SDK parses into the details that Describe.Columns models.
	Property            ssot.Attribute
	PropertyValue       ssot.Attribute
	PropertyDefault     ssot.Attribute
	PropertyDescription ssot.Attribute
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
		Enum:         MfaEnrollmentOption,
		DiffSuppress: ssot.SuppressDescribe,
		CreateSql:    ssot.PlainParameter(),
		SetSql:       ssot.PlainParameter(),
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
		CreateSql:   ssot.ListOfObjects(ssot.Sql("CLIENT_POLICY"), ssot.Parentheses()),
		Elem: &ssot.Attribute{
			Type:      ssot.AttrObject,
			CreateSql: ssot.ObjectParameter(),
			Children: []ssot.Attribute{
				{
					Name:        "client_type",
					Type:        ssot.AttrEnum,
					Required:    true,
					Description: "The client or driver type.",
					Enum:        ClientPolicyDriverType,
					CreateSql:   ssot.EnumKeyword(ssot.NoQuotes(), ssot.Required()),
				},
				authenticationPolicyClientPolicyMinimumVersion,
				authenticationPolicyClientPolicyParams,
			},
		},
	},
	MfaPolicy: ssot.Attribute{
		Name:         "mfa_policy",
		Type:         ssot.AttrObject,
		CreateSql:    authenticationPolicyMfaPolicyClause,
		SetSql:       authenticationPolicyMfaPolicyClause,
		ChildrenRule: ssot.ChildrenAtLeastOne,
		Description:  "Specifies the multi-factor authentication (MFA) methods that users can use as a second factor of authentication.",
		Children:     authenticationPolicyMfaPolicyChildren,
	},
	PatPolicy: ssot.Attribute{
		Name:         "pat_policy",
		Type:         ssot.AttrObject,
		CreateSql:    authenticationPolicyPatPolicyClause,
		SetSql:       authenticationPolicyPatPolicyClause,
		ChildrenRule: ssot.ChildrenAtLeastOne,
		Description:  "Specifies the policy for programmatic access tokens.",
		Children:     authenticationPolicyPatPolicyChildren,
	},
	WorkloadIdentityPolicy: ssot.Attribute{
		Name:         "workload_identity_policy",
		Type:         ssot.AttrObject,
		CreateSql:    authenticationPolicyWorkloadIdentityPolicyClause,
		SetSql:       authenticationPolicyWorkloadIdentityPolicyClause,
		ChildrenRule: ssot.ChildrenAtLeastOne,
		Description:  "Specifies the policy for workload identity federation.",
		Children:     authenticationPolicyWorkloadIdentityPolicyChildren,
	},
	SecurityIntegrations: authenticationPolicySecurityIntegrations,
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the authentication policy.",
	},
	CreatedOn:             ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, OptionalInDb: true, Readonly: true},
	Name:                  ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	DatabaseName:          ssot.Attribute{Name: "database_name", Type: ssot.AttrString, OptionalInDb: true, ManualConvert: true, Readonly: true},
	SchemaName:            ssot.Attribute{Name: "schema_name", Type: ssot.AttrString, OptionalInDb: true, ManualConvert: true, Readonly: true},
	Kind:                  ssot.Attribute{Name: "kind", Type: ssot.AttrString, Readonly: true},
	Owner:                 ssot.Attribute{Name: "owner", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
	OwnerRoleType:         ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
	Options:               ssot.Attribute{Name: "options", Type: ssot.AttrString, Readonly: true},
	DescribeMfaEnrollment: ssot.Attribute{Name: "describe_mfa_enrollment", ShowName: "mfa_enrollment", Type: ssot.AttrEnum, Enum: MfaEnrollmentReadOption, Readonly: true},
	DescribeSecurityIntegrations: ssot.Attribute{
		Name: "describe_security_integrations", ShowName: "security_integrations", Type: ssot.AttrObject, Required: true, Readonly: true,
		Children: []ssot.Attribute{
			{Name: "all", Type: ssot.AttrBool},
			{Name: "security_integrations", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"}},
		},
	},
	// The elements are the hand-written AuthenticationPolicyClientPolicyDetails, see SdkDescribeDetailsNestedField.
	DescribeClientPolicy:           ssot.Attribute{Name: "describe_client_policy", ShowName: "client_policy", Type: ssot.AttrList, Readonly: true, Elem: &ssot.Attribute{Type: ssot.AttrString}},
	DescribeMfaPolicy:              ssot.Attribute{Name: "describe_mfa_policy", ShowName: "mfa_policy", Type: ssot.AttrObject, Required: true, Readonly: true, Children: authenticationPolicyMfaPolicyChildren},
	DescribePatPolicy:              ssot.Attribute{Name: "describe_pat_policy", ShowName: "pat_policy", Type: ssot.AttrObject, Required: true, Readonly: true, Children: authenticationPolicyPatPolicyChildren},
	DescribeWorkloadIdentityPolicy: ssot.Attribute{Name: "describe_workload_identity_policy", ShowName: "workload_identity_policy", Type: ssot.AttrObject, Required: true, Readonly: true, Children: authenticationPolicyWorkloadIdentityPolicyChildren},
	Property:                       ssot.Attribute{Name: "property", Type: ssot.AttrString, Readonly: true},
	PropertyValue:                  ssot.Attribute{Name: "value", Type: ssot.AttrString, Readonly: true},
	PropertyDefault:                ssot.Attribute{Name: "default", Type: ssot.AttrString, Readonly: true},
	PropertyDescription:            ssot.Attribute{Name: "description", Type: ssot.AttrString, Readonly: true},
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
	AttributesModification: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		// describe_output is a block of the raw DESCRIBE strings. The typed columns and the comment are mapped from the original rows in the ext.
		ssot.WithSchema(ssot.SchemaDescribe(ssot.SchemaDescribeManualFields(
			&AuthenticationPolicyAttrs.AuthenticationMethods,
			&AuthenticationPolicyAttrs.ClientTypes,
			&AuthenticationPolicyAttrs.DescribeSecurityIntegrations,
			&AuthenticationPolicyAttrs.Comment,
			&AuthenticationPolicyAttrs.DescribeClientPolicy,
			&AuthenticationPolicyAttrs.DescribeMfaPolicy,
			&AuthenticationPolicyAttrs.DescribePatPolicy,
			&AuthenticationPolicyAttrs.DescribeWorkloadIdentityPolicy,
		))),
		ssot.WithSdk(
			ssot.SdkShow(
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalExtendedIn(),
				ssot.SdkShowOptionalOn(),
				ssot.SdkShowOptionalStartsWith(),
				ssot.SdkShowOptionalLimit(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDLike, ssot.SdkShowByIDExtendedIn),
				ssot.SdkShowResultFilterHook(),
				ssot.SdkShowPlainOnly(&AuthenticationPolicyAttrs.TargetScopes, "[]AuthenticationPolicyTargetScope"),
			),
			ssot.SdkDescribe(
				ssot.SdkDescribeQueryStructName("DescribeAuthenticationPolicy"),
				ssot.SdkDescribeDetails(
					ssot.SdkDescribeDetailsNestedField(&AuthenticationPolicyAttrs.DescribeClientPolicy),
					ssot.SdkDescribeDetailsMethod(),
					ssot.SdkDescribeDetailsRowsField("descriptions", "AuthenticationPolicyDescriptions"),
					ssot.SdkDescribeDetailsParseIgnoreUnknown(),
					ssot.SdkDescribeDetailsParse(
						ssot.DescribeProperty(&AuthenticationPolicyAttrs.Name, "NAME"),
						ssot.DescribeProperty(&AuthenticationPolicyAttrs.Owner, "OWNER"),
						ssot.DescribePropertyParsed(&AuthenticationPolicyAttrs.AuthenticationMethods, "AUTHENTICATION_METHODS"),
						ssot.DescribePropertyParsed(&AuthenticationPolicyAttrs.DescribeMfaEnrollment, "MFA_ENROLLMENT"),
						ssot.DescribePropertyParsed(&AuthenticationPolicyAttrs.ClientTypes, "CLIENT_TYPES"),
						ssot.DescribePropertyParsed(&AuthenticationPolicyAttrs.DescribeSecurityIntegrations, "SECURITY_INTEGRATIONS"),
						ssot.DescribePropertyParsed(&AuthenticationPolicyAttrs.Comment, "COMMENT"),
						ssot.DescribePropertyParsed(&AuthenticationPolicyAttrs.DescribeClientPolicy, "CLIENT_POLICY"),
						ssot.DescribePropertyParsed(&AuthenticationPolicyAttrs.DescribeMfaPolicy, "MFA_POLICY"),
						ssot.DescribePropertyParsed(&AuthenticationPolicyAttrs.DescribePatPolicy, "PAT_POLICY"),
						ssot.DescribePropertyParsed(&AuthenticationPolicyAttrs.DescribeWorkloadIdentityPolicy, "WORKLOAD_IDENTITY_POLICY"),
					),
				),
			),
			ssot.SdkDrop(),
			ssot.SdkCreate(),
			ssot.SdkAlter(),
		),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithHandleExternalChangesInDescribe(ssot.ExtCall()),
		ssot.WithHierarchyRename(ssot.ExtCall()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithCreateInSdk(ssot.Generated()),
		ssot.WithAlterSetInSdk(ssot.Generated()),
		ssot.WithAlterUnsetInSdk(ssot.Generated()),
		ssot.WithStateUpgraders(ssot.Generated()),
		ssot.WithCustomizeDiff(ssot.Generated()),
		ssot.WithSetConfigFields(ssot.Generated()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		// TODO(SNOW-1818849): unassign policies before dropping. Generated Delete does not do that yet.
		ssot.WithDropSafelyInSdk(ssot.Generated()),
	},
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"v2_9_0_AuthenticationPolicyStateUpgrader",
	},
	// mfa_authentication_methods is a deprecated schema-ext field and DESCRIBE returns the READ variants.
	AdditionalEnums: []*ssot.Enum{MfaAuthenticationMethodsOption, MfaAuthenticationMethodsReadOption, MfaEnrollmentReadOption},
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
	Create: ssot.CreateSpec{
		Clauses: []*ssot.Attribute{
			&AuthenticationPolicyAttrs.AuthenticationMethods,
			&AuthenticationPolicyAttrs.MfaEnrollment,
			&AuthenticationPolicyAttrs.MfaPolicy,
			&AuthenticationPolicyAttrs.ClientTypes,
			&AuthenticationPolicyAttrs.ClientPolicy,
			&authenticationPolicySecurityIntegrationsOption,
			&AuthenticationPolicyAttrs.PatPolicy,
			&AuthenticationPolicyAttrs.WorkloadIdentityPolicy,
			&AuthenticationPolicyAttrs.Comment,
		},
	},
	Set: []*ssot.Attribute{
		&AuthenticationPolicyAttrs.AuthenticationMethods,
		&AuthenticationPolicyAttrs.MfaEnrollment,
		&AuthenticationPolicyAttrs.MfaPolicy,
		&AuthenticationPolicyAttrs.ClientTypes,
		&AuthenticationPolicyAttrs.ClientPolicy,
		&authenticationPolicySecurityIntegrationsOption,
		&AuthenticationPolicyAttrs.PatPolicy,
		&AuthenticationPolicyAttrs.WorkloadIdentityPolicy,
		&AuthenticationPolicyAttrs.Comment,
	},
	Unset: []*ssot.Attribute{
		&AuthenticationPolicyAttrs.ClientTypes,
		&AuthenticationPolicyAttrs.ClientPolicy,
		&AuthenticationPolicyAttrs.AuthenticationMethods,
		&AuthenticationPolicyAttrs.SecurityIntegrations,
		&AuthenticationPolicyAttrs.MfaEnrollment,
		&AuthenticationPolicyAttrs.MfaPolicy,
		&AuthenticationPolicyAttrs.PatPolicy,
		&AuthenticationPolicyAttrs.WorkloadIdentityPolicy,
		&AuthenticationPolicyAttrs.Comment,
	},
	Alter: ssot.AlterSpec{
		Clauses: []ssot.AlterClause{ssot.AlterSet(), ssot.AlterUnset(), ssot.AlterRenameTo()},
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
	DescribeRows: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&AuthenticationPolicyAttrs.Property,
			&AuthenticationPolicyAttrs.PropertyValue,
			&AuthenticationPolicyAttrs.PropertyDefault,
			&AuthenticationPolicyAttrs.PropertyDescription,
		},
	},
	// Order matches sdk.AuthenticationPolicyDetails field order (Id excluded: not a Describe column).
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&AuthenticationPolicyAttrs.Name,
			&AuthenticationPolicyAttrs.Owner,
			&AuthenticationPolicyAttrs.AuthenticationMethods,
			&AuthenticationPolicyAttrs.DescribeMfaEnrollment,
			&AuthenticationPolicyAttrs.ClientTypes,
			&AuthenticationPolicyAttrs.DescribeSecurityIntegrations,
			&AuthenticationPolicyAttrs.Comment,
			&AuthenticationPolicyAttrs.DescribeClientPolicy,
			&AuthenticationPolicyAttrs.DescribeMfaPolicy,
			&AuthenticationPolicyAttrs.DescribePatPolicy,
			&AuthenticationPolicyAttrs.DescribeWorkloadIdentityPolicy,
		},
	},
}
