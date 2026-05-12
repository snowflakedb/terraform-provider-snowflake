package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemGetAWSSNSIAMPolicy(t *testing.T) {
	r := require.New(t)
	sb := NewSystemGetAWSSNSIAMPolicyBuilder("arn:aws:sns:us-east-1:1234567890123456:mytopic")

	r.Equal(`SELECT SYSTEM$GET_AWS_SNS_IAM_POLICY('arn:aws:sns:us-east-1:1234567890123456:mytopic') AS "POLICY"`, sb.Select())
}
