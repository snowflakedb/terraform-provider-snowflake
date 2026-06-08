package sdk

import (
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (r showAuthenticationPolicyDBRow) excludeFromShow() bool {
	return !r.DatabaseName.Valid || !r.SchemaName.Valid
}

func (r showAuthenticationPolicyDBRow) additionalConvert(result *AuthenticationPolicy) error {
	var errs []error
	if !r.DatabaseName.Valid {
		errs = append(errs, fmt.Errorf("missing database name for authentication policy with name: %s", r.Name))
	}
	if !r.SchemaName.Valid {
		errs = append(errs, fmt.Errorf("missing schema name for authentication policy with name: %s", r.Name))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	mapNullStringToNonNullableField(&result.DatabaseName, r.DatabaseName)
	mapNullStringToNonNullableField(&result.SchemaName, r.SchemaName)
	return nil
}

type AuthenticationPolicyDetails []AuthenticationPolicyDescription

func (v AuthenticationPolicyDetails) GetAuthenticationMethods() ([]AuthenticationMethodsOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "AUTHENTICATION_METHODS" })
	if err != nil {
		return nil, err
	}
	return collections.MapErr(ParseCommaSeparatedStringArray(raw.Value, false), ToAuthenticationMethodsOption)
}

func (v AuthenticationPolicyDetails) Raw(key string) string {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == key })
	if err != nil {
		return ""
	}
	return raw.Value
}

func (v AuthenticationPolicyDetails) GetMfaEnrollment() (MfaEnrollmentReadOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "MFA_ENROLLMENT" })
	if err != nil {
		return "", err
	}
	return ToMfaEnrollmentReadOption(raw.Value)
}

func (v AuthenticationPolicyDetails) GetClientTypes() ([]ClientTypesOption, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "CLIENT_TYPES" })
	if err != nil {
		return nil, err
	}
	return collections.MapErr(ParseCommaSeparatedStringArray(raw.Value, false), ToClientTypesOption)
}

func (v AuthenticationPolicyDetails) GetSecurityIntegrations() ([]AccountObjectIdentifier, error) {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == "SECURITY_INTEGRATIONS" })
	if err != nil {
		return nil, err
	}
	return ParseCommaSeparatedAccountObjectIdentifierArray(raw.Value)
}
