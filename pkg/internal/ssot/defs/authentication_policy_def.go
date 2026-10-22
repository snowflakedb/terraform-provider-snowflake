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

// AuthenticationPolicyAttrs holds reusable Authentication Policy fields so SHOW and DESCRIBE
// Columns can point at the same definitions instead of restating names.
//
// Omitted from Attributes because SSOT has no nested block type:
//   - mfa_authentication_methods: deprecated schema field; it is not a CREATE/ALTER column.
//   - client_policy: set of {client_type, minimum_version}.
//   - mfa_policy: single block of an enum and a set of enums.
//   - pat_policy: single block of ints, a bool, and an enum.
//   - workload_identity_policy: single block of sets.
//
// Those nested fields are also omitted from DESCRIBE Columns. A writable column must be in
// Attributes, and the raw DESCRIBE property strings are not a separate attribute.
var AuthenticationPolicyAttrs = struct {
	AuthenticationMethods ssot.Attribute
	MfaEnrollment         ssot.Attribute
	ClientTypes           ssot.Attribute
	SecurityIntegrations  ssot.Attribute
	Comment               ssot.Attribute
	CreatedOn             ssot.Attribute
	Name                  ssot.Attribute
	DatabaseName          ssot.Attribute
	SchemaName            ssot.Attribute
	Kind                  ssot.Attribute
	Owner                 ssot.Attribute
	OwnerRoleType         ssot.Attribute
	Options               ssot.Attribute
	TargetScopes          ssot.Attribute
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
	},
	Comment: ssot.Attribute{
		Name:          "comment",
		Type:          ssot.AttrString,
		Description:   "Specifies a comment for the authentication policy.",
		CreateMapping: ssot.StringAttributeCreateBuilder,
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
	Name:                   "AuthenticationPolicy",
	NamePlural:             "AuthenticationPolicies",
	DescribeFailRead:       true,
	IdType:                 "SchemaObjectIdentifier",
	AttributesModification: true,
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithHandleExternalChangesInDescribe(ssot.ExtCall()),
		ssot.WithHierarchyRename(ssot.ExtCall()),
		ssot.WithStateUpgraders(ssot.Generated()),
		// database and schema are ForceNew. The identity template does not emit that yet.
		ssot.WithIdentitySchema(ssot.ExtCall()),
		// The describe elem schema is DescribeAuthenticationPolicyDescribeDetailsSchema.
		ssot.WithDescribeOutputSchema(ssot.ExtCall()),
	},
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"v2_9_0_AuthenticationPolicyStateUpgrader",
	},
	Attributes: []ssot.Attribute{
		AuthenticationPolicyAttrs.AuthenticationMethods,
		AuthenticationPolicyAttrs.MfaEnrollment,
		AuthenticationPolicyAttrs.ClientTypes,
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
