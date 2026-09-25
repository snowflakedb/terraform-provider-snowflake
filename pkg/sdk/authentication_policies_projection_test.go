package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthenticationPolicyDescribeProjection(t *testing.T) {
	name := AuthenticationPolicyDescription{Property: "NAME", Value: "ap", Default: "", Description: "name"}
	methods := AuthenticationPolicyDescription{Property: "AUTHENTICATION_METHODS", Value: "[PASSWORD]", Default: "", Description: "methods"}
	comment := AuthenticationPolicyDescription{Property: "COMMENT", Value: "c", Default: "", Description: "comment"}
	unknown := AuthenticationPolicyDescription{Property: "UNEXPECTED", Value: "x", Default: "", Description: ""}

	details := AsAuthenticationPolicyDescribe([]AuthenticationPolicyDescription{name, methods, comment, unknown})
	require.Equal(t, "ap", details.Name)
	require.Equal(t, "[PASSWORD]", details.AuthenticationMethods)
	require.Equal(t, "c", details.Comment)
	require.Empty(t, details.Owner)
	require.Nil(t, AsAuthenticationPolicyDescribe(nil))
}
