package sdk

type NetworkRuleType string

const (
	NetworkRuleTypeIpv4             NetworkRuleType = "IPV4"
	NetworkRuleTypeAwsVpcEndpointId NetworkRuleType = "AWSVPCEID"
	NetworkRuleTypeAzureLinkId      NetworkRuleType = "AZURELINKID"
	NetworkRuleTypeHostPort         NetworkRuleType = "HOST_PORT"
)

type NetworkRuleMode string

const (
	NetworkRuleModeIngress       NetworkRuleMode = "INGRESS"
	NetworkRuleModeInternalStage NetworkRuleMode = "INTERNAL_STAGE"
	NetworkRuleModeEgress        NetworkRuleMode = "EGRESS"
)
