package resourceassert

import (
	"fmt"

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

	for i, e := range entries {
		s.StringValueSet(fmt.Sprintf("client_policy.%d.client_type", i), string(e.ClientType))
		if e.Params != nil && e.Params.MinimumVersion != nil {
			s.StringValueSet(fmt.Sprintf("client_policy.%d.minimum_version", i), *e.Params.MinimumVersion)
		}
	}

	return s
}
