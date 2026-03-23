package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type NetworkRuleType string

const (
	NetworkRuleTypeIpv4             NetworkRuleType = "IPV4"
	NetworkRuleTypeAwsVpcEndpointId NetworkRuleType = "AWSVPCEID"
	NetworkRuleTypeAzureLinkId      NetworkRuleType = "AZURELINKID"
	NetworkRuleTypeGcpPscId         NetworkRuleType = "GCPPSCID"
	NetworkRuleTypeHostPort         NetworkRuleType = "HOST_PORT"
	NetworkRuleTypePrivateHostPort  NetworkRuleType = "PRIVATE_HOST_PORT"
)

var AllNetworkRuleTypes = []NetworkRuleType{
	NetworkRuleTypeIpv4,
	NetworkRuleTypeAwsVpcEndpointId,
	NetworkRuleTypeAzureLinkId,
	NetworkRuleTypeGcpPscId,
	NetworkRuleTypeHostPort,
	NetworkRuleTypePrivateHostPort,
}

func ToNetworkRuleType(s string) (NetworkRuleType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllNetworkRuleTypes, NetworkRuleType(s)) {
		return "", fmt.Errorf("invalid network rule type: %s", s)
	}
	return NetworkRuleType(s), nil
}

type NetworkRuleMode string

const (
	NetworkRuleModeIngress         NetworkRuleMode = "INGRESS"
	NetworkRuleModeInternalStage   NetworkRuleMode = "INTERNAL_STAGE"
	NetworkRuleModeEgress          NetworkRuleMode = "EGRESS"
	NetworkRuleModePostgresIngress NetworkRuleMode = "POSTGRES_INGRESS"
	NetworkRuleModePostgresEgress  NetworkRuleMode = "POSTGRES_EGRESS"
)

var AllNetworkRuleModes = []NetworkRuleMode{
	NetworkRuleModeIngress,
	NetworkRuleModeInternalStage,
	NetworkRuleModeEgress,
	NetworkRuleModePostgresIngress,
	NetworkRuleModePostgresEgress,
}

func ToNetworkRuleMode(s string) (NetworkRuleMode, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllNetworkRuleModes, NetworkRuleMode(s)) {
		return "", fmt.Errorf("invalid network rule mode: %s", s)
	}
	return NetworkRuleMode(s), nil
}
