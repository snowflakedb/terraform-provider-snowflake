package schemas

import (
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var NetworkPolicyPropertiesNames = []string{
	"ALLOWED_IP_LIST",
	"BLOCKED_IP_LIST",
	"ALLOWED_NETWORK_RULE_LIST",
	"BLOCKED_NETWORK_RULE_LIST",
}

// NetworkPolicyPropertiesToSchema maps DESCRIBE NETWORK POLICY rows onto describe_output.
// Snowflake omits unset lists from DESCRIBE, so only returned properties are written.
// The generated NetworkPolicyDetailsToSchema always emits all four keys (empty string
// for missing ones) and would change state for empty policies.
func NetworkPolicyPropertiesToSchema(networkPolicyProperties []sdk.NetworkPolicyProperty) map[string]any {
	networkPolicySchema := make(map[string]any)
	for _, property := range networkPolicyProperties {
		if !slices.Contains(NetworkPolicyPropertiesNames, property.Name) {
			log.Printf("[WARN] unexpected property %v in network policy returned from Snowflake", property.Name)
			continue
		}
		networkPolicySchema[strings.ToLower(property.Name)] = property.Value
	}
	return networkPolicySchema
}
