package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

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

	targetScopes, err := parseAuthenticationPolicyTargetScopes(r.Options)
	if err != nil {
		return fmt.Errorf("parsing target scopes for authentication policy with name %s: %w", r.Name, err)
	}
	result.TargetScopes = targetScopes
	return nil
}

// parseAuthenticationPolicyTargetScopes extracts the "target_scopes" JSON array from the raw options column.
func parseAuthenticationPolicyTargetScopes(options string) ([]AuthenticationPolicyTargetScope, error) {
	if options == "" {
		return nil, nil
	}
	var raw struct {
		TargetScopes []AuthenticationPolicyTargetScope `json:"target_scopes"`
	}
	if err := json.Unmarshal([]byte(options), &raw); err != nil {
		return nil, err
	}
	slices.Sort(raw.TargetScopes)
	return raw.TargetScopes, nil
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

// AsAuthenticationPolicyDescribe projects DESCRIBE output onto TypeString describe_output columns.
func AsAuthenticationPolicyDescribe(descriptions []AuthenticationPolicyDescription) *AuthenticationPolicyDescribeDetails {
	if descriptions == nil {
		return nil
	}
	details := &AuthenticationPolicyDescribeDetails{}
	for _, description := range descriptions {
		switch description.Property {
		case "NAME":
			details.Name = description.Value
		case "OWNER":
			details.Owner = description.Value
		case "AUTHENTICATION_METHODS":
			details.AuthenticationMethods = description.Value
		case "MFA_ENROLLMENT":
			details.MfaEnrollment = description.Value
		case "CLIENT_TYPES":
			details.ClientTypes = description.Value
		case "SECURITY_INTEGRATIONS":
			details.SecurityIntegrations = description.Value
		case "COMMENT":
			details.Comment = description.Value
		case "CLIENT_POLICY":
			details.ClientPolicy = description.Value
		case "MFA_POLICY":
			details.MfaPolicy = description.Value
		case "PAT_POLICY":
			details.PatPolicy = description.Value
		case "WORKLOAD_IDENTITY_POLICY":
			details.WorkloadIdentityPolicy = description.Value
		}
	}
	return details
}
