package resourceassert

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *AuthenticationPolicyResourceAssert) HasAuthenticationMethodsEnum(expected ...sdk.AuthenticationMethodsOption) *AuthenticationPolicyResourceAssert {
	return s.HasAuthenticationMethods(collections.Map(expected, func(v sdk.AuthenticationMethodsOption) string {
		return string(v)
	})...)
}

func (s *AuthenticationPolicyResourceAssert) HasClientTypesEnum(expected ...sdk.ClientTypesOption) *AuthenticationPolicyResourceAssert {
	return s.HasClientTypes(collections.Map(expected, func(v sdk.ClientTypesOption) string {
		return string(v)
	})...)
}

func (s *AuthenticationPolicyResourceAssert) HasClientPolicyEntries(entries ...sdk.AuthenticationPolicyClientPolicyRequest) *AuthenticationPolicyResourceAssert {
	s.CollectionLength("client_policy", len(entries))

	for i, e := range entries {
		s.StringValueSet(fmt.Sprintf("client_policy.%d.client_type", i), string(e.ClientType))
		s.StringValueSet(fmt.Sprintf("client_policy.%d.minimum_version", i), e.Params.MinimumVersion)
	}

	return s
}

func (s *AuthenticationPolicyResourceAssert) HasPatPolicy(expected sdk.AuthenticationPolicyPatPolicyRequest) *AuthenticationPolicyResourceAssert {
	defaultExpiryInDays := "0"
	if expected.DefaultExpiryInDays != nil {
		defaultExpiryInDays = strconv.Itoa(*expected.DefaultExpiryInDays)
	}
	maxExpiryInDays := "0"
	if expected.MaxExpiryInDays != nil {
		maxExpiryInDays = strconv.Itoa(*expected.MaxExpiryInDays)
	}
	requireRoleRestrictionForServiceUsers := r.BooleanDefault
	if expected.RequireRoleRestrictionForServiceUsers != nil {
		requireRoleRestrictionForServiceUsers = strconv.FormatBool(*expected.RequireRoleRestrictionForServiceUsers)
	}
	requireRoleRestrictionForPersonUsers := r.BooleanDefault
	if expected.RequireRoleRestrictionForPersonUsers != nil {
		requireRoleRestrictionForPersonUsers = strconv.FormatBool(*expected.RequireRoleRestrictionForPersonUsers)
	}
	var networkPolicyEvaluation string
	if expected.NetworkPolicyEvaluation != nil {
		networkPolicyEvaluation = string(*expected.NetworkPolicyEvaluation)
	}
	s.CollectionLength("pat_policy", 1)
	s.ValueSet("pat_policy.0.default_expiry_in_days", defaultExpiryInDays)
	s.ValueSet("pat_policy.0.max_expiry_in_days", maxExpiryInDays)
	s.ValueSet("pat_policy.0.require_role_restriction_for_service_users", requireRoleRestrictionForServiceUsers)
	s.ValueSet("pat_policy.0.require_role_restriction_for_person_users", requireRoleRestrictionForPersonUsers)
	s.SetContainsExactlyStringValues("pat_policy.0.blocked_roles_list", collections.Map(expected.BlockedRolesList, func(v sdk.StringListItemWrapper) string {
		return sdk.NewAccountObjectIdentifierFromFullyQualifiedName(v.Value).Name()
	})...)
	s.ValueSet("pat_policy.0.network_policy_evaluation", networkPolicyEvaluation)
	return s
}
