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
