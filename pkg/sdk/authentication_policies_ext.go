package sdk

import (
	"errors"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

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

func parseDescribeProperties(properties []AuthenticationPolicyDescription, id AccountObjectIdentifier) (*AuthenticationPolicyDetailsLegacy, error) {
	details := new(AuthenticationPolicyDetailsLegacy)
	var errs []error
	for _, prop := range properties {
		switch strings.ToUpper(prop.Property) {
		case "NAME":
			//details.
		case "OWNER":
		case "COMMENT":
		case "AUTHENTICATION_METHODS":
		case "CLIENT_TYPES":
		case "CLIENT_POLICY":
		case "SECURITY_INTEGRATIONS":
		case "MFA_ENROLLMENT":
		case "MFA_POLICY":
		case "PAT_POLICY":
		case "WORKLOAD_IDENTITY_POLICY":
		}
	}
	return details, errors.Join(errs...)
}
