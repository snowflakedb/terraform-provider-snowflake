package resourceassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CortexAgentResourceAssert) HasProfile(expected sdk.CortexAgentProfile) *CortexAgentResourceAssert {
	c.CollectionLength("profile", 1)
	if expected.DisplayName != nil {
		c.StringValueSet("profile.0.display_name", *expected.DisplayName)
	} else {
		c.StringValueSet("profile.0.display_name", "")
	}
	if expected.Avatar != nil {
		c.StringValueSet("profile.0.avatar", *expected.Avatar)
	} else {
		c.StringValueSet("profile.0.avatar", "")
	}
	if expected.Color != nil {
		c.StringValueSet("profile.0.color", *expected.Color)
	} else {
		c.StringValueSet("profile.0.color", "")
	}
	return c
}
