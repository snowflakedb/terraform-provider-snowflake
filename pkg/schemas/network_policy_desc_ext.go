package schemas

import (
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var NetworkPolicyPropertiesNames = []string{
	"ALLOWED_IP_LIST",
	"BLOCKED_IP_LIST",
	"ALLOWED_NETWORK_RULE_LIST",
	"BLOCKED_NETWORK_RULE_LIST",
}

func NetworkPolicyPropertiesToSchema(networkPolicyProperties []sdk.NetworkPolicyProperty) map[string]any {
	for _, property := range networkPolicyProperties {
		if !slices.Contains(NetworkPolicyPropertiesNames, property.Name) {
			log.Printf("[WARN] unexpected property %v in network policy returned from Snowflake", property.Name)
		}
	}
	details := sdk.AsNetworkPolicyDescribe(networkPolicyProperties)
	if details == nil {
		return map[string]any{}
	}
	return NetworkPolicyDetailsToSchema(details)
}
