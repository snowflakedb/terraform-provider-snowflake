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

// Nested blocks and the deprecated field have no SSOT attribute type.
// Set diffs and the MFA enrollment special case are also applied here.
func authenticationPolicyAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["authentication_methods"].DiffSuppressFunc = NormalizeAndCompareEnumsInSet("authentication_methods", sdk.ToAuthenticationMethodsOption)
	s["client_types"].DiffSuppressFunc = NormalizeAndCompareEnumsInSet("client_types", sdk.ToClientTypesOption)
	s["mfa_enrollment"].DiffSuppressFunc = SuppressIfAny(
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
	)
	s["security_integrations"].Elem.(*schema.Schema).ValidateDiagFunc = IsValidIdentifier[sdk.AccountObjectIdentifier]()
	s["security_integrations"].DiffSuppressFunc = NormalizeAndCompareIdentifiersInSet("security_integrations")

	s["mfa_authentication_methods"] = &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToMfaAuthenticationMethodsOption),
		},
		Optional:         true,
		DiffSuppressFunc: ignoreAlwaysSuppressFunc,
		Description:      fmt.Sprintf("A list of authentication methods that enforce multi-factor authentication (MFA) during login. Authentication methods not listed in this parameter do not prompt for multi-factor authentication. Allowed values are %s.", possibleValuesListed(sdk.AllMfaAuthenticationMethodsOptions)),
		Deprecated:       "This field is deprecated and will be removed in the future. Currently, it has no effect. Use the new `enforce_mfa_on_external_authentication` field instead. Read our [BCR Migration Guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/SNOWFLAKE_BCR_MIGRATION_GUIDE.md#changes-in-authentication-policies) for more migration steps and more details.",
	}
	s["client_policy"] = &schema.Schema{
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
	}
	s["mfa_policy"] = &schema.Schema{
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
	}
	s["pat_policy"] = &schema.Schema{
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
	}
	s["workload_identity_policy"] = &schema.Schema{
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
	}
	return s
}
