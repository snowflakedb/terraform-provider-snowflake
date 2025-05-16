package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func (s *SecretWithAuthorizationCodeGrantResourceAssert) HasOauthRefreshTokenExpiryTimeNotEmpty() *SecretWithAuthorizationCodeGrantResourceAssert {
	s.AddAssertion(assert.ValuePresent("oauth_refresh_token_expiry_time"))
	return s
}
