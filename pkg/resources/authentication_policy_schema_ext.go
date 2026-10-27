package resources

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Nested blocks are generated. This keeps the deprecated set, AtLeastOneOf, and the
// case-sensitive wording on nested enum sets. Set diffs and the MFA enrollment special case are also applied here.
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
	mfaAtLeastOneOf := []string{"mfa_policy.0.enforce_mfa_on_external_authentication", "mfa_policy.0.allowed_methods"}
	mfa := s["mfa_policy"].Elem.(*schema.Resource).Schema
	mfa["enforce_mfa_on_external_authentication"].AtLeastOneOf = mfaAtLeastOneOf
	mfa["allowed_methods"].AtLeastOneOf = mfaAtLeastOneOf
	// Nested enum sets stay case-sensitive. The generator would normalize them.
	mfa["allowed_methods"].DiffSuppressFunc = nil
	mfa["allowed_methods"].Description = fmt.Sprintf("Specifies the allowed methods for the MFA policy. Valid values are: %s. These values are case-sensitive due to Terraform limitations (it's a nested field). Prefer using uppercased values.", possibleValuesListed(sdk.AllMfaPolicyAllowedMethodsOptions))

	patAtLeastOneOf := []string{
		"pat_policy.0.default_expiry_in_days",
		"pat_policy.0.max_expiry_in_days",
		"pat_policy.0.require_role_restriction_for_service_users",
		"pat_policy.0.network_policy_evaluation",
	}
	pat := s["pat_policy"].Elem.(*schema.Resource).Schema
	for _, name := range []string{"default_expiry_in_days", "max_expiry_in_days", "require_role_restriction_for_service_users", "network_policy_evaluation"} {
		pat[name].AtLeastOneOf = patAtLeastOneOf
	}
	pat["network_policy_evaluation"].Description = fmt.Sprintf("Specifies the network policy evaluation for the PAT. Valid values are: %s.", possibleValuesListed(sdk.AllNetworkPolicyEvaluationOptions))

	providers := s["workload_identity_policy"].Elem.(*schema.Resource).Schema["allowed_providers"]
	providers.DiffSuppressFunc = nil
	providers.Description = fmt.Sprintf("Specifies the allowed providers for the workload identity policy. Valid values are: %s. These values are case-sensitive due to Terraform limitations (it's a nested field). Prefer using uppercased values.", possibleValuesListed(sdk.AllAllowedProviderOptions))
	return s
}
