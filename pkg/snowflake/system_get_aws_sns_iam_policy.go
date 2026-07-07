package snowflake

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/jmoiron/sqlx"
)

// SystemGetAWSSNSIAMPolicyBuilder abstracts calling the SYSTEM$GET_AWS_SNS_IAM_POLICY system function.
type SystemGetAWSSNSIAMPolicyBuilder struct {
	awsSnsTopicArn string
}

// SystemGetAWSSNSIAMPolicy returns a pointer to a builder that abstracts calling the the SYSTEM$GET_AWS_SNS_IAM_POLICY system function.
func NewSystemGetAWSSNSIAMPolicyBuilder(awsSnsTopicArn string) *SystemGetAWSSNSIAMPolicyBuilder {
	return &SystemGetAWSSNSIAMPolicyBuilder{
		awsSnsTopicArn: awsSnsTopicArn,
	}
}

// Select generates the select statement for obtaining the aws sns iam policy.
func (pb *SystemGetAWSSNSIAMPolicyBuilder) Select() string {
	value := sdk.SingleQuotes.Modify(pb.awsSnsTopicArn)
	return fmt.Sprintf(`SELECT SYSTEM$GET_AWS_SNS_IAM_POLICY(%s) AS "POLICY"`, value)
}

type AWSSNSIAMPolicy struct {
	Policy string `db:"POLICY"`
}

// ScanAWSSNSIAMPolicy convert a result into a.
func ScanAWSSNSIAMPolicy(row *sqlx.Row) (*AWSSNSIAMPolicy, error) {
	p := &AWSSNSIAMPolicy{}
	e := row.StructScan(p)
	return p, e
}
