package sdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (v *authenticationPolicies) DescribeDetails(ctx context.Context, id SchemaObjectIdentifier) (*AuthenticationPolicyDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseAuthenticationPolicyProperties(properties)
}

type AuthenticationPolicyDetailsLegacy []AuthenticationPolicyDescription

func (v AuthenticationPolicyDetailsLegacy) GetAuthenticationMethods() ([]AuthenticationMethodsOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "AUTHENTICATION_METHODS" })
	if err != nil {
		return nil, err
	}
	return collections.MapErr(ParseCommaSeparatedStringArray(raw.Value, false), ToAuthenticationMethodsOption)
}

func (v AuthenticationPolicyDetailsLegacy) Raw(key string) string {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == key })
	if err != nil {
		return ""
	}
	return raw.Value
}

func (v AuthenticationPolicyDetailsLegacy) GetMfaEnrollment() (MfaEnrollmentReadOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "MFA_ENROLLMENT" })
	if err != nil {
		return "", err
	}
	return ToMfaEnrollmentReadOption(raw.Value)
}

func (v AuthenticationPolicyDetailsLegacy) GetClientTypes() ([]ClientTypesOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "CLIENT_TYPES" })
	if err != nil {
		return nil, err
	}
	return collections.MapErr(ParseCommaSeparatedStringArray(raw.Value, false), ToClientTypesOption)
}

func (v AuthenticationPolicyDetailsLegacy) GetClientPolicies() ([]ClientTypesOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "CLIENT_POLICIES" })
	if err != nil {
		return nil, err
	}

	//keyValueUtil := func(s string) map[string]string {
	//	s = strings.TrimPrefix(s, "{")
	//	s = strings.TrimSuffix(s, "}")
	//	result := make(map[string]string)
	//	for _, part := range ParseOuterCommaSeparatedStringArray(fmt.Sprintf("[%s]", s), false) {
	//		key, value, _ := strings.Cut(part, "=")
	//		result[key] = value
	//	}
	//	return result
	//}

	//for driverTypeRaw, driverPoliciesRaw := range keyValueUtil(raw.Value) {
	//	// TODO:
	//}

	return collections.MapErr(ParseCommaSeparatedStringArray(raw.Value, false), ToClientTypesOption)
}

func (v AuthenticationPolicyDetailsLegacy) GetSecurityIntegrations() ([]AccountObjectIdentifier, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "SECURITY_INTEGRATIONS" })
	if err != nil {
		return nil, err
	}
	return ParseCommaSeparatedAccountObjectIdentifierArray(raw.Value)
}

func (v AuthenticationPolicyDetailsLegacy) GetMfaAuthenticationMethods() ([]MfaAuthenticationMethodsReadOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "MFA_AUTHENTICATION_METHODS" })
	if err != nil {
		return nil, err
	}
	return collections.MapErr(ParseCommaSeparatedStringArray(raw.Value, false), ToMfaAuthenticationMethodsReadOption)
}

func parseAuthenticationPolicyProperties(properties []AuthenticationPolicyDescription) (*AuthenticationPolicyDetails, error) {
	details := new(AuthenticationPolicyDetails)

	// TODO: Take out as util (see catalog integration)
	keyValueUtil := func(s string) map[string]string {
		s = strings.TrimPrefix(s, "{")
		s = strings.TrimSuffix(s, "}")
		result := make(map[string]string)
		for _, part := range ParseOuterCommaSeparatedStringArray(fmt.Sprintf("[%s]", s), false) {
			key, value, _ := strings.Cut(part, "=")
			result[key] = value
		}
		return result
	}

	var errs []error
	for _, prop := range properties {
		switch strings.ToUpper(prop.Property) {
		case "NAME":
			if prop.Value != "null" {
				name, err := ParseAccountObjectIdentifier(prop.Value)
				if err != nil {
					errs = append(errs, err)
				} else {
					details.Name = &name
				}
			}
		case "OWNER":
			if prop.Value != "null" {
				details.Owner = String(prop.Value)
			}
		case "COMMENT":
			if prop.Value != "null" {
				details.Comment = String(prop.Value)
			}
		case "AUTHENTICATION_METHODS":
			if authenticationMethods, err := collections.MapErr(ParseCommaSeparatedStringArray(prop.Value, false), ToAuthenticationMethodsOption); err != nil {
				errs = append(errs, err)
			} else {
				details.AuthenticationMethods = authenticationMethods
			}
		case "CLIENT_TYPES":
			if clientTypes, err := collections.MapErr(ParseCommaSeparatedStringArray(prop.Value, false), ToClientTypesOption); err != nil {
				errs = append(errs, err)
			} else {
				details.ClientTypes = clientTypes
			}
		case "CLIENT_POLICY":
			//for key, value := range keyValueUtil(prop.Value) {
			//
			//}
		case "SECURITY_INTEGRATIONS":
			if strings.ToUpper(prop.Value) == "[ALL]" {
				details.SecurityIntegrations.All = true
			} else {
				for _, securityIntegrationId := range ParseCommaSeparatedStringArray(prop.Value, false) {
					id, err := ParseAccountObjectIdentifier(securityIntegrationId)
					if err != nil {
						errs = append(errs, err)
					} else {
						details.SecurityIntegrations.SecurityIntegrations = append(details.SecurityIntegrations.SecurityIntegrations, id)
					}
				}
			}
		case "MFA_ENROLLMENT":
			// TODO: Should it be read option ?
			mfaEnrollment, err := ToMfaEnrollmentOption(prop.Value)
			if err != nil {
				errs = append(errs, err)
			} else {
				details.MfaEnrollment = mfaEnrollment
			}
		case "MFA_POLICY":
			for key, value := range keyValueUtil(prop.Value) {
				driverType, driverTypeErr := ToClientPolicyDriverType(key)
				if driverTypeErr != nil {
					errs = append(errs, driverTypeErr)
				} else {
					details.ClientPolicy[driverType] = new(ClientPolicyDetails)

					s := strings.TrimPrefix(value, "{")
					s = strings.TrimSuffix(s, "}")
					parts := ParseOuterCommaSeparatedStringArray(fmt.Sprintf("[%s]", s), false)
					for _, part := range parts {
						key, value, _ := strings.Cut(part, "=")
						switch key {
						case "MINIMUM_VERSION":
							details.ClientPolicy[driverType].MinimumVersion = value
						}
					}
				}
			}
		case "PAT_POLICY":
			for key, value := range keyValueUtil(prop.Value) {
				switch key {
				case "DEFAULT_EXPIRY_IN_DAYS":
					defaultExpiryInDays, err := strconv.Atoi(value)
					if err != nil {
						errs = append(errs, err)
					} else {
						details.PatPolicy.DefaultExpiryInDays = defaultExpiryInDays
					}
				case "MAX_EXPIRY_IN_DAYS":
					maxExpiryInDays, err := strconv.Atoi(value)
					if err != nil {
						errs = append(errs, err)
					} else {
						details.PatPolicy.MaxExpiryInDays = maxExpiryInDays
					}
				case "NETWORK_POLICY_EVALUATION":
					networkPolicyEvaluation, err := ToNetworkPolicyEvaluationOption(value)
					if err != nil {
						errs = append(errs, err)
					} else {
						details.PatPolicy.NetworkPolicyEvaluation = networkPolicyEvaluation
					}
				}
			}
		case "WORKLOAD_IDENTITY_POLICY":
			for key, value := range keyValueUtil(prop.Value) {
				switch key {
				case "ALLOWED_PROVIDERS":
					for _, allowedProvider := range ParseCommaSeparatedStringArray(value, false) {
						allowedProviderEnum, err := ToAllowedProviderOption(allowedProvider)
						if err != nil {
							errs = append(errs, err)
						} else {
							details.WorkloadIdentityPolicy.AllowedProviders = append(details.WorkloadIdentityPolicy.AllowedProviders, allowedProviderEnum)
						}
					}
				case "ALLOWED_AWS_ACCOUNTS":
					if strings.ToUpper(value) == "[ALL]" {
						details.WorkloadIdentityPolicy.AllowedAwsAccounts.All = true
					} else {
						details.WorkloadIdentityPolicy.AllowedAwsAccounts.AllowedAwsAccounts = ParseCommaSeparatedStringArray(value, false)
					}
				case "ALLOWED_AZURE_ISSUERS":
					if strings.ToUpper(value) == "[ALL]" {
						details.WorkloadIdentityPolicy.AllowedAzureIssuers.All = true
					} else {
						details.WorkloadIdentityPolicy.AllowedAzureIssuers.AllowedAzureIssuers = ParseCommaSeparatedStringArray(value, false)
					}
				case "ALLOWED_OIDC_ISSUERS":
					if strings.ToUpper(value) == "[ALL]" {
						details.WorkloadIdentityPolicy.AllowedOidcIssuers.All = true
					} else {
						details.WorkloadIdentityPolicy.AllowedOidcIssuers.AllowedOidcIssuers = ParseCommaSeparatedStringArray(value, false)
					}
				}
			}
		}
	}

	return details, errors.Join(errs...)
}
