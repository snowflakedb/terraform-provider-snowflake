package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNetworkPolicyDescribeProjection(t *testing.T) {
	allowed := NetworkPolicyProperty{Name: "ALLOWED_IP_LIST", Value: "1.1.1.1"}
	blocked := NetworkPolicyProperty{Name: "BLOCKED_IP_LIST", Value: "2.2.2.2"}
	unknown := NetworkPolicyProperty{Name: "UNEXPECTED", Value: "x"}

	details := AsNetworkPolicyDescribe([]NetworkPolicyProperty{allowed, blocked, unknown})
	require.Equal(t, "1.1.1.1", details.AllowedIpList)
	require.Equal(t, "2.2.2.2", details.BlockedIpList)
	require.Empty(t, details.AllowedNetworkRuleList)
	require.Nil(t, AsNetworkPolicyDescribe(nil))
}
