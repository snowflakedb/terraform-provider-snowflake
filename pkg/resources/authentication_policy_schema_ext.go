package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var authenticationPolicyIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the authentication policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the authentication policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the authentication policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

var authenticationPolicyAttributesSchemaExt = map[string]*schema.Schema{
	"authentication_methods": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToAuthenticationMethodsOption),
		},
		Optional:         true,
		DiffSuppressFunc: NormalizeAndCompareEnumsInSet("authentication_methods", sdk.ToAuthenticationMethodsOption),
		Description:      fmt.Sprintf("A list of authentication methods that are allowed during login. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllAuthenticationMethodsOptions)),
	},
	"mfa_authentication_methods": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToMfaAuthenticationMethodsOption),
		},
		Optional:         true,
		DiffSuppressFunc: ignoreAlwaysSuppressFunc,
		Description:      fmt.Sprintf("A list of authentication methods that enforce multi-factor authentication (MFA) during login. Authentication methods not listed in this parameter do not prompt for multi-factor authentication. Allowed values are %s.", possibleValuesListed(sdk.AllMfaAuthenticationMethodsOptions)),
		Deprecated:       "This field is deprecated and will be removed in the future. Currently, it has no effect. Use the new `enforce_mfa_on_external_authentication` field instead. Read our [BCR Migration Guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/SNOWFLAKE_BCR_MIGRATION_GUIDE.md#changes-in-authentication-policies) for more migration steps and more details.",
	},
	"mfa_enrollment": {
		Type:     schema.TypeString,
		Optional: true,
		Description: fmt.Sprintf("Determines whether a user must enroll in multi-factor authentication. Valid values are (case-insensitive): %s. When REQUIRED is specified, Enforces users to enroll in MFA."+
			" If this value is used, then the `client_types` parameter must include `snowflake_ui`, because Snowsight is the only place users can enroll in multi-factor authentication (MFA)."+
			" Note that when you set this value to OPTIONAL, and your account setup forces users to enroll in MFA, then Snowflake may set quietly this value to `REQUIRED_PASSWORD_ONLY`, which may cause permadiff."+
			" In this case, you may want to adjust this field value.", possibleValuesListed(sdk.AllMfaEnrollmentOptions)),
		ValidateDiagFunc: sdkValidation(sdk.ToMfaEnrollmentOption),
		DiffSuppressFunc: SuppressIfAny(
			NormalizeAndCompare(sdk.ToMfaEnrollmentOption),
			// We need a custom diff because when this field is set to OPTIONAL, the value is returned as REQUIRED_PASSWORD_ONLY.
			func(_, oldRaw, newRaw string, _ *schema.ResourceData) bool {
				old, err := sdk.ToMfaEnrollmentReadOption(oldRaw)
				if err != nil {
					return false
				}
				new, err := sdk.ToMfaEnrollmentOption(newRaw)
				if err != nil {
					return false
				}
				return old == sdk.MfaEnrollmentReadOptionRequiredPasswordOnly && new == sdk.MfaEnrollmentOptionOptional
			},
		),
	},
	"client_types": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToClientTypesOption),
		},
		Optional:         true,
		DiffSuppressFunc: NormalizeAndCompareEnumsInSet("client_types", sdk.ToClientTypesOption),
		Description:      fmt.Sprintf("A list of clients that can authenticate with Snowflake. If a client tries to connect, and the client is not one of the valid `client_types`, then the login attempt fails. Valid values are (case-insensitive): %s. The `client_types` property of an authentication policy is a best effort method to block user logins based on specific clients. It should not be used as the sole control to establish a security boundary.", possibleValuesListed(sdk.AllClientTypesOptions)),
	},
	"client_policy": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Allows to set policies per-client type. Only valid when `client_types` is empty, contains ALL, or contains DRIVERS.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"client_type": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      fmt.Sprintf("The client or driver type. Valid values (case-insensitive): %s.", possibleValuesListed(sdk.AllClientPolicyDriverTypes)),
					ValidateDiagFunc: sdkValidation(sdk.ToClientPolicyDriverType),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToClientPolicyDriverType),
				},
				"minimum_version": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Minimum allowed version for this client/driver type (e.g. '1.14.1').",
				},
			},
		},
	},
	"security_integrations": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("security_integrations"),
		Optional:         true,
		Description:      "A list of security integrations the authentication policy is associated with. This parameter has no effect when `saml` or `oauth` are not in the `authentication_methods` list. All values in the `security_integrations` list must be compatible with the values in the `authentication_methods` list. For example, if `security_integrations` contains a SAML security integration, and `authentication_methods` contains OAUTH, then you cannot create the authentication policy. To allow all security integrations use `ALL` as parameter.",
	},
	"mfa_policy": {
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enforce_mfa_on_external_authentication": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      fmt.Sprintf("Determines whether multi-factor authentication (MFA) is enforced on external authentication. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllEnforceMfaOnExternalAuthenticationOptions)),
					ValidateDiagFunc: sdkValidation(sdk.ToEnforceMfaOnExternalAuthenticationOption),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToEnforceMfaOnExternalAuthenticationOption),
					AtLeastOneOf:     []string{"mfa_policy.0.enforce_mfa_on_external_authentication", "mfa_policy.0.allowed_methods"},
				},
				"allowed_methods": {
					Type:        schema.TypeSet,
					Optional:    true,
					Description: fmt.Sprintf("Specifies the allowed methods for the MFA policy. Valid values are: %s. These values are case-sensitive due to Terraform limitations (it's a nested field). Prefer using uppercased values.", possibleValuesListed(sdk.AllMfaPolicyAllowedMethodsOptions)),
					Elem: &schema.Schema{
						Type:             schema.TypeString,
						ValidateDiagFunc: sdkValidation(sdk.ToMfaPolicyAllowedMethodsOption),
					},
					AtLeastOneOf: []string{"mfa_policy.0.enforce_mfa_on_external_authentication", "mfa_policy.0.allowed_methods"},
				},
			},
		},
		Optional:    true,
		Description: "Specifies the multi-factor authentication (MFA) methods that users can use as a second factor of authentication.",
	},
	"pat_policy": {
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default_expiry_in_days": {
					Type:             schema.TypeInt,
					Optional:         true,
					Description:      "Specifies the default expiration time (in days) for a programmatic access token.",
					AtLeastOneOf:     []string{"pat_policy.0.default_expiry_in_days", "pat_policy.0.max_expiry_in_days", "pat_policy.0.require_role_restriction_for_service_users", "pat_policy.0.network_policy_evaluation"},
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				},
				"max_expiry_in_days": {
					Type:             schema.TypeInt,
					Optional:         true,
					Description:      "Specifies the maximum number of days that can be set for the expiration time for a programmatic access token.",
					AtLeastOneOf:     []string{"pat_policy.0.default_expiry_in_days", "pat_policy.0.max_expiry_in_days", "pat_policy.0.require_role_restriction_for_service_users", "pat_policy.0.network_policy_evaluation"},
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				},
				"require_role_restriction_for_service_users": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      booleanStringFieldDescription("If true, when you generate a programmatic access token for a service user, you must restrict the use of that token to a specific role."),
					AtLeastOneOf:     []string{"pat_policy.0.default_expiry_in_days", "pat_policy.0.max_expiry_in_days", "pat_policy.0.require_role_restriction_for_service_users", "pat_policy.0.network_policy_evaluation"},
					Default:          BooleanDefault,
					ValidateDiagFunc: validateBooleanString,
				},
				"network_policy_evaluation": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: sdkValidation(sdk.ToNetworkPolicyEvaluationOption),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToNetworkPolicyEvaluationOption),
					Description:      fmt.Sprintf("Specifies the network policy evaluation for the PAT. Valid values are: %s.", possibleValuesListed(sdk.AllNetworkPolicyEvaluationOptions)),
					AtLeastOneOf:     []string{"pat_policy.0.default_expiry_in_days", "pat_policy.0.max_expiry_in_days", "pat_policy.0.require_role_restriction_for_service_users", "pat_policy.0.network_policy_evaluation"},
				},
			},
		},
		Optional:    true,
		Description: "Specifies the policy for programmatic access tokens.",
	},
	"workload_identity_policy": {
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allowed_providers": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type:             schema.TypeString,
						ValidateDiagFunc: sdkValidation(sdk.ToAllowedProviderOption),
					},
					Description: fmt.Sprintf("Specifies the allowed providers for the workload identity policy. Valid values are: %s. These values are case-sensitive due to Terraform limitations (it's a nested field). Prefer using uppercased values.", possibleValuesListed(sdk.AllAllowedProviderOptions)),
				},
				"allowed_aws_accounts": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Specifies the list of AWS account IDs allowed by the authentication policy during workload identity authentication of type `AWS`.",
				},
				"allowed_azure_issuers": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Specifies the list of Azure Entra ID issuers allowed by the authentication policy during workload identity authentication of type `AZURE`.",
				},
				"allowed_oidc_issuers": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Specifies the list of OIDC issuers allowed by the authentication policy during workload identity authentication of type `OIDC`.",
				},
			},
		},
		Optional:    true,
		Description: "Specifies the policy for workload identity federation.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the authentication policy.",
	},
}

var authenticationPolicyShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW AUTHENTICATION POLICIES` for the given policy.",
		Elem: &schema.Resource{
			Schema: schemas.ShowAuthenticationPolicySchema,
		},
	},
}

var authenticationPolicyDescribeOutputSchemaExt = map[string]*schema.Schema{
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE AUTHENTICATION POLICY` for the given policy.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeAuthenticationPolicyDescribeDetailsSchema,
		},
	},
}
