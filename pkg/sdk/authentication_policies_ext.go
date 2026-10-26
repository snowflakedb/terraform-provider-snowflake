package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

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

type AuthenticationPolicyDescriptions []AuthenticationPolicyDescription

func (v AuthenticationPolicyDescriptions) Raw(key string) string {
	raw, err := collections.FindFirst(v, func(r AuthenticationPolicyDescription) bool { return r.Property == key })
	if err != nil {
		return ""
	}
	return raw.Value
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

// Describe projects the original DESC rows onto the string columns used by describe_output.
func (d *AuthenticationPolicyDetails) Describe() *AuthenticationPolicyDescribeDetails {
	if d == nil {
		return nil
	}
	return AsAuthenticationPolicyDescribe(d.descriptions)
}

func (v *authenticationPolicies) DescribeDetails(ctx context.Context, id SchemaObjectIdentifier) (*AuthenticationPolicyDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseAuthenticationPolicyDetails(properties, id)
}

func parseAuthenticationPolicyDetails(properties []AuthenticationPolicyDescription, id SchemaObjectIdentifier) (*AuthenticationPolicyDetails, error) {
	details := &AuthenticationPolicyDetails{
		Id:           id,
		descriptions: properties,
	}
	var errs []error
	for _, prop := range properties {
		switch prop.Property {
		case "NAME":
			details.Name = prop.Value
		case "OWNER":
			details.Owner = prop.Value
		case "COMMENT":
			details.Comment = emptyIfNull(prop.Value)
		case "AUTHENTICATION_METHODS":
			if methods, err := parseEnumList(prop.Value, ToAuthenticationMethodsOption); err != nil {
				errs = append(errs, fmt.Errorf("parsing AUTHENTICATION_METHODS: %w", err))
			} else {
				details.AuthenticationMethods = methods
			}
		case "MFA_ENROLLMENT":
			if enrollment, err := ToMfaEnrollmentReadOption(prop.Value); err != nil {
				errs = append(errs, fmt.Errorf("parsing MFA_ENROLLMENT: %w", err))
			} else {
				details.MfaEnrollment = enrollment
			}
		case "CLIENT_TYPES":
			if clientTypes, err := parseEnumList(prop.Value, ToClientTypesOption); err != nil {
				errs = append(errs, fmt.Errorf("parsing CLIENT_TYPES: %w", err))
			} else {
				details.ClientTypes = clientTypes
			}
		case "SECURITY_INTEGRATIONS":
			if securityIntegrations, err := parseAuthenticationPolicySecurityIntegrations(prop.Value); err != nil {
				errs = append(errs, fmt.Errorf("parsing SECURITY_INTEGRATIONS: %w", err))
			} else {
				details.SecurityIntegrations = securityIntegrations
			}
		case "CLIENT_POLICY":
			if clientPolicy, err := parseAuthenticationPolicyClientPolicy(prop.Value); err != nil {
				errs = append(errs, fmt.Errorf("parsing CLIENT_POLICY: %w", err))
			} else {
				details.ClientPolicy = clientPolicy
			}
		case "MFA_POLICY":
			if mfaPolicy, err := parseAuthenticationPolicyMfaPolicy(prop.Value); err != nil {
				errs = append(errs, fmt.Errorf("parsing MFA_POLICY: %w", err))
			} else {
				details.MfaPolicy = mfaPolicy
			}
		case "PAT_POLICY":
			if patPolicy, err := parseAuthenticationPolicyPatPolicy(prop.Value); err != nil {
				errs = append(errs, fmt.Errorf("parsing PAT_POLICY: %w", err))
			} else {
				details.PatPolicy = patPolicy
			}
		case "WORKLOAD_IDENTITY_POLICY":
			if workloadIdentityPolicy, err := parseAuthenticationPolicyWorkloadIdentityPolicy(prop.Value); err != nil {
				errs = append(errs, fmt.Errorf("parsing WORKLOAD_IDENTITY_POLICY: %w", err))
			} else {
				details.WorkloadIdentityPolicy = workloadIdentityPolicy
			}
		}
	}
	return details, errors.Join(errs...)
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

func parseAuthenticationPolicyClientPolicy(value string) ([]AuthenticationPolicyClientPolicyEntryDetails, error) {
	pairs, err := parseSnowflakeMap(value)
	if err != nil {
		return nil, err
	}
	entries := make([]AuthenticationPolicyClientPolicyEntryDetails, 0, len(pairs))
	var errs []error
	for _, pair := range pairs {
		clientType, err := ToClientPolicyDriverType(pair.key)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		entry := AuthenticationPolicyClientPolicyEntryDetails{ClientType: clientType}
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
