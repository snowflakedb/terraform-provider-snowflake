package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeNetworkPolicySchema represents output of DESCRIBE query for the single NetworkPolicy.
var DescribeNetworkPolicySchema = map[string]*schema.Schema{
	"allowed_ip_list": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"blocked_ip_list": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"allowed_network_rule_list": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"blocked_network_rule_list": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = DescribeNetworkPolicySchema

func NetworkPolicyPropertiesToSchema(networkPolicyProperties []sdk.NetworkPolicyProperty) map[string]any {
	networkPolicySchema := make(map[string]any)
	for _, property := range networkPolicyProperties {
		switch property.Name {
		case "ALLOWED_IP_LIST",
			"BLOCKED_IP_LIST",
			"ALLOWED_NETWORK_RULE_LIST",
			"BLOCKED_NETWORK_RULE_LIST":
			networkPolicySchema[strings.ToLower(property.Name)] = property.Value
		}
	}
	return networkPolicySchema
}

var _ = NetworkPolicyPropertiesToSchema
