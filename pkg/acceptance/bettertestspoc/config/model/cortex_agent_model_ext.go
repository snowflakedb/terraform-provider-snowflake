package model

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func CortexAgentWithSpecification(resourceName, database, schema, name, specification string) *CortexAgentModel {
	model := CortexAgent(resourceName, database, schema, name, "")
	// This prevents double quotes from being added around yamlencode.
	model.WithSpecificationValue(config.UnquotedWrapperVariable(specification))
	return model
}

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

// SampleSpecAsYamlencodeHCL returns a multiline Terraform/HCL expression for use with
// model.CortexAgentWithSpecification.
func SampleSpecAsYamlencodeHCL(response string) string {
	return fmt.Sprintf(`yamlencode({
  orchestration = {
    budget = {
      seconds = 30
      tokens  = 16000
    }
  }
  instructions = {
    response = %[1]s%[2]s%[1]s
  }
})`, config.SnowflakeProviderConfigQuoteMarker, response)
}
