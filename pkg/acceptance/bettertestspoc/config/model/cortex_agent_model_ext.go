package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CortexAgentModel) WithProfile(profile sdk.CortexAgentProfile) *CortexAgentModel {
	m := map[string]tfconfig.Variable{}
	if profile.DisplayName != nil {
		m["display_name"] = tfconfig.StringVariable(*profile.DisplayName)
	}
	if profile.Avatar != nil {
		m["avatar"] = tfconfig.StringVariable(*profile.Avatar)
	}
	if profile.Color != nil {
		m["color"] = tfconfig.StringVariable(*profile.Color)
	}
	c.Profile = tfconfig.ListVariable(tfconfig.ObjectVariable(m))
	return c
}
