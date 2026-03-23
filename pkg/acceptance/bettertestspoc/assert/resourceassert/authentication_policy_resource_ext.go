package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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

func (s *AuthenticationPolicyResourceAssert) HasClientPolicyEntries(entries ...sdk.AuthenticationPolicyClientPolicyEntry) *AuthenticationPolicyResourceAssert {
	s.CollectionLength("client_policy", len(entries))

	for _, e := range entries {
		s.StringValueSet("client_policy.*.client_type", string(e.ClientType))
		if e.Params != nil && e.Params.MinimumVersion != nil {
			s.StringValueSet("client_policy.*.minimum_version", *e.Params.MinimumVersion)
		}
	}

	return s
}
