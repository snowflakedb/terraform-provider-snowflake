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
