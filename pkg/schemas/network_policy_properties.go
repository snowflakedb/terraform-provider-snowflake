package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeNetworkPolicySchema is the name-keyed consumer for DESCRIBE NETWORK POLICY.
// The property row is generated as NetworkPolicyPropertySchema / NetworkPolicyPropertyToSchema (UsedAsListEntry).
// Public describe_output stays four TypeString keys (property.Value), not a list of {name, value} rows.
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
