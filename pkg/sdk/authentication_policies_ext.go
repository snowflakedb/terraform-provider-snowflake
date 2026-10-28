package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (r authenticationPoliciesRow) excludeFromShow() bool {
	return !r.DatabaseName.Valid || !r.SchemaName.Valid
}

func (r authenticationPoliciesRow) additionalConvert(result *AuthenticationPolicy) error {
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

type AuthenticationPolicyDescriptions []AuthenticationPolicyProperty

func (v AuthenticationPolicyDescriptions) Raw(key string) string {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyProperty) bool { return r.Property == key })
	if err != nil {
		return ""
	}
	return raw.Value
}

// AuthenticationPolicyClientPolicyDetails is an element of the client_policy column of the details.
// TODO [next PRs]: generate this element struct from the Elem children of the column, like the object columns. External volume has the same
// hand-written element (ExternalVolumeStorageLocationDetails).
type AuthenticationPolicyClientPolicyDetails struct {
	ClientType     ClientPolicyDriverType
	MinimumVersion string
}

func (d *AuthenticationPolicyDetails) ID() SchemaObjectIdentifier {
	return d.Id
}

// Descriptions returns the DESC rows Snowflake returned. Parsed fields are a projection of these rows.
func (d *AuthenticationPolicyDetails) Descriptions() AuthenticationPolicyDescriptions {
	if d == nil {
		return nil
	}
	return d.descriptions
}

func validateAuthenticationPolicyDescribed(AuthenticationPolicyDetails) error {
	return nil
}

// The column parsers below are called by the generated ParseAuthenticationPolicyDescribed, one for each property.

// TODO [next PRs]: The parsers of authentication_methods, client_types, mfa_enrollment and comment can be generated from the column
// (list of enums, single enum, a string whose null value means empty), see the ssot README Limitations.
func ParseAuthenticationPolicyAuthenticationMethodsDescribed(value string) ([]AuthenticationMethodsOption, error) {
	result, err := parseEnumList(value, ToAuthenticationMethodsOption)
	return result, wrapParsingError("AUTHENTICATION_METHODS", err)
}

func ParseAuthenticationPolicyMfaEnrollmentDescribed(value string) (MfaEnrollmentReadOption, error) {
	result, err := ToMfaEnrollmentReadOption(value)
	return result, wrapParsingError("MFA_ENROLLMENT", err)
}

func ParseAuthenticationPolicyClientTypesDescribed(value string) ([]ClientTypesOption, error) {
	result, err := parseEnumList(value, ToClientTypesOption)
	return result, wrapParsingError("CLIENT_TYPES", err)
}

func ParseAuthenticationPolicySecurityIntegrationsDescribed(value string) (AuthenticationPolicySecurityIntegrationsDetails, error) {
	result, err := parseAuthenticationPolicySecurityIntegrations(value)
	return result, wrapParsingError("SECURITY_INTEGRATIONS", err)
}

func ParseAuthenticationPolicyCommentDescribed(value string) (string, error) {
	return emptyIfNull(value), nil
}

func ParseAuthenticationPolicyClientPolicyDescribed(value string) ([]AuthenticationPolicyClientPolicyDetails, error) {
	result, err := parseAuthenticationPolicyClientPolicy(value)
	return result, wrapParsingError("CLIENT_POLICY", err)
}

func ParseAuthenticationPolicyMfaPolicyDescribed(value string) (AuthenticationPolicyMfaPolicyDetails, error) {
	result, err := parseAuthenticationPolicyMfaPolicy(value)
	return result, wrapParsingError("MFA_POLICY", err)
}

func ParseAuthenticationPolicyPatPolicyDescribed(value string) (AuthenticationPolicyPatPolicyDetails, error) {
	result, err := parseAuthenticationPolicyPatPolicy(value)
	return result, wrapParsingError("PAT_POLICY", err)
}

func ParseAuthenticationPolicyWorkloadIdentityPolicyDescribed(value string) (AuthenticationPolicyWorkloadIdentityPolicyDetails, error) {
	result, err := parseAuthenticationPolicyWorkloadIdentityPolicy(value)
	return result, wrapParsingError("WORKLOAD_IDENTITY_POLICY", err)
}

func wrapParsingError(property string, err error) error {
	if err != nil {
		return fmt.Errorf("parsing %s: %w", property, err)
	}
	return nil
}

func parseEnumList[T ~string](value string, convert func(string) (T, error)) ([]T, error) {
	return collections.MapErr(ParseCommaSeparatedStringArray(value, false), convert)
}

func parseAuthenticationPolicySecurityIntegrations(value string) (AuthenticationPolicySecurityIntegrationsDetails, error) {
	items := ParseCommaSeparatedStringArray(value, false)
	if len(items) == 1 && strings.EqualFold(items[0], "ALL") {
		return AuthenticationPolicySecurityIntegrationsDetails{All: true}, nil
	}
	for _, item := range items {
		if strings.EqualFold(strings.Trim(item, `"`), "ALL") {
			return AuthenticationPolicySecurityIntegrationsDetails{}, fmt.Errorf("cannot combine ALL with other values: %s", value)
		}
	}
	ids, err := collections.MapErr(items, ParseAccountObjectIdentifier)
	if err != nil {
		return AuthenticationPolicySecurityIntegrationsDetails{}, err
	}
	return AuthenticationPolicySecurityIntegrationsDetails{SecurityIntegrations: ids}, nil
}

func parseAuthenticationPolicyClientPolicy(value string) ([]AuthenticationPolicyClientPolicyDetails, error) {
	pairs, err := parseSnowflakeMap(value)
	if err != nil {
		return nil, err
	}
	entries := make([]AuthenticationPolicyClientPolicyDetails, 0, len(pairs))
	var errs []error
	for _, pair := range pairs {
		clientType, err := ToClientPolicyDriverType(pair.key)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		entry := AuthenticationPolicyClientPolicyDetails{ClientType: clientType}
		inner, err := parseSnowflakeMap(pair.value)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", pair.key, err))
			continue
		}
		for _, innerPair := range inner {
			if innerPair.key == "MINIMUM_VERSION" {
				entry.MinimumVersion = innerPair.value
			}
		}
		entries = append(entries, entry)
	}
	return entries, errors.Join(errs...)
}

func parseAuthenticationPolicyMfaPolicy(value string) (AuthenticationPolicyMfaPolicyDetails, error) {
	pairs, err := parseSnowflakeMap(value)
	if err != nil {
		return AuthenticationPolicyMfaPolicyDetails{}, err
	}
	details := AuthenticationPolicyMfaPolicyDetails{}
	var errs []error
	for _, pair := range pairs {
		switch pair.key {
		case "ALLOWED_METHODS":
			if methods, err := parseEnumList(pair.value, ToMfaPolicyAllowedMethodsOption); err != nil {
				errs = append(errs, err)
			} else {
				details.AllowedMethods = methods
			}
		case "ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION":
			if enforcement, err := ToEnforceMfaOnExternalAuthenticationOption(pair.value); err != nil {
				errs = append(errs, err)
			} else {
				details.EnforceMfaOnExternalAuthentication = enforcement
			}
		}
	}
	return details, errors.Join(errs...)
}

func parseAuthenticationPolicyPatPolicy(value string) (AuthenticationPolicyPatPolicyDetails, error) {
	pairs, err := parseSnowflakeMap(value)
	if err != nil {
		return AuthenticationPolicyPatPolicyDetails{}, err
	}
	details := AuthenticationPolicyPatPolicyDetails{}
	var errs []error
	for _, pair := range pairs {
		switch pair.key {
		case "DEFAULT_EXPIRY_IN_DAYS":
			if days, err := strconv.Atoi(pair.value); err != nil {
				errs = append(errs, fmt.Errorf("DEFAULT_EXPIRY_IN_DAYS: %w", err))
			} else {
				details.DefaultExpiryInDays = days
			}
		case "MAX_EXPIRY_IN_DAYS":
			if days, err := strconv.Atoi(pair.value); err != nil {
				errs = append(errs, fmt.Errorf("MAX_EXPIRY_IN_DAYS: %w", err))
			} else {
				details.MaxExpiryInDays = days
			}
		case "NETWORK_POLICY_EVALUATION":
			if evaluation, err := ToNetworkPolicyEvaluationOption(pair.value); err != nil {
				errs = append(errs, err)
			} else {
				details.NetworkPolicyEvaluation = evaluation
			}
		case "REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS":
			if required, err := strconv.ParseBool(pair.value); err != nil {
				errs = append(errs, fmt.Errorf("REQUIRE_ROLE_RESTRICTION_FOR_SERVICE_USERS: %w", err))
			} else {
				details.RequireRoleRestrictionForServiceUsers = required
			}
		}
	}
	return details, errors.Join(errs...)
}

func parseAuthenticationPolicyWorkloadIdentityPolicy(value string) (AuthenticationPolicyWorkloadIdentityPolicyDetails, error) {
	pairs, err := parseSnowflakeMap(value)
	if err != nil {
		return AuthenticationPolicyWorkloadIdentityPolicyDetails{}, err
	}
	details := AuthenticationPolicyWorkloadIdentityPolicyDetails{}
	var errs []error
	for _, pair := range pairs {
		switch pair.key {
		case "ALLOWED_PROVIDERS":
			if providers, err := parseEnumList(pair.value, ToAllowedProviderOption); err != nil {
				errs = append(errs, err)
			} else {
				details.AllowedProviders = providers
			}
		case "ALLOWED_AWS_ACCOUNTS":
			details.AllowedAwsAccounts = ParseCommaSeparatedStringArray(pair.value, false)
		case "ALLOWED_AZURE_ISSUERS":
			details.AllowedAzureIssuers = ParseCommaSeparatedStringArray(pair.value, false)
		case "ALLOWED_OIDC_ISSUERS":
			details.AllowedOidcIssuers = ParseCommaSeparatedStringArray(pair.value, false)
		}
	}
	return details, errors.Join(errs...)
}

type snowflakeMapEntry struct {
	key   string
	value string
}

// parseSnowflakeMap splits a Snowflake describe map such as {KEY=value, KEY2=[a, b]}.
// Entry splitting reuses ParseOuterCommaSeparatedStringArray, so commas inside [] and double quotes stay in the value.
// Nested client-policy objects contain a single MINIMUM_VERSION field, so their braces need no extra splitting.
func parseSnowflakeMap(value string) ([]snowflakeMapEntry, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "null" {
		return nil, nil
	}
	if len(value) < 2 || value[0] != '{' || value[len(value)-1] != '}' {
		return nil, fmt.Errorf("expected a map value, got %q", value)
	}
	inner := strings.TrimSpace(value[1 : len(value)-1])
	if inner == "" {
		return nil, nil
	}
	entries := make([]snowflakeMapEntry, 0)
	for _, part := range ParseOuterCommaSeparatedStringArray(inner, false) {
		if part == "" {
			continue
		}
		key, entryValue, ok := strings.Cut(part, "=")
		if !ok {
			return nil, fmt.Errorf("expected key=value entry, got %q", part)
		}
		entries = append(entries, snowflakeMapEntry{
			key:   strings.TrimSpace(key),
			value: strings.TrimSpace(entryValue),
		})
	}
	return entries, nil
}
