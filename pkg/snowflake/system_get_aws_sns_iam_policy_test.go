package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemGetAWSSNSIAMPolicy(t *testing.T) {
	testCases := []struct {
		name           string
		awsSnsTopicArn string
		expected       string
	}{
		{
			name:           "basic",
			awsSnsTopicArn: "arn:aws:sns:us-east-1:1234567890123456:mytopic",
			expected:       `SELECT SYSTEM$GET_AWS_SNS_IAM_POLICY('arn:aws:sns:us-east-1:1234567890123456:mytopic') AS "POLICY"`,
		},
		{
			name:           "single quote is escaped",
			awsSnsTopicArn: "arn:aws:sns:us-east-1:1234567890123456:my'topic",
			expected:       `SELECT SYSTEM$GET_AWS_SNS_IAM_POLICY('arn:aws:sns:us-east-1:1234567890123456:my\'topic') AS "POLICY"`,
		},
		{
			name:           "backslash is escaped",
			awsSnsTopicArn: `arn:aws:sns\mytopic`,
			expected:       `SELECT SYSTEM$GET_AWS_SNS_IAM_POLICY('arn:aws:sns\\mytopic') AS "POLICY"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			sb := NewSystemGetAWSSNSIAMPolicyBuilder(tc.awsSnsTopicArn)

			r.Equal(tc.expected, sb.Select())
		})
	}
}
