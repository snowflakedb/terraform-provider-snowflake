package sdk

import "encoding/json"

// NetworkRulesSnowflakeDto is needed to unpack the applied network rules from the JSON response from Snowflake
type NetworkRulesSnowflakeDto struct {
	FullyQualifiedRuleName string
}

func ParseNetworkRulesSnowflakeDto(networkRulesStringValue string) ([]NetworkRulesSnowflakeDto, error) {
	var networkRules []NetworkRulesSnowflakeDto
	err := json.Unmarshal([]byte(networkRulesStringValue), &networkRules)
	if err != nil {
		return nil, err
	}
	return networkRules, nil
}

func (r *CreateNetworkPolicyRequest) GetName() AccountObjectIdentifier {
	return r.name
}

// AsNetworkPolicyDescribe projects DESCRIBE output onto TypeString describe_output columns.
func AsNetworkPolicyDescribe(properties []NetworkPolicyProperty) *NetworkPolicyDetails {
	if properties == nil {
		return nil
	}
	details := &NetworkPolicyDetails{}
	for _, property := range properties {
		switch property.Name {
		case "ALLOWED_IP_LIST":
			details.AllowedIpList = property.Value
		case "BLOCKED_IP_LIST":
			details.BlockedIpList = property.Value
		case "ALLOWED_NETWORK_RULE_LIST":
			details.AllowedNetworkRuleList = property.Value
		case "BLOCKED_NETWORK_RULE_LIST":
			details.BlockedNetworkRuleList = property.Value
		}
	}
	return details
}
