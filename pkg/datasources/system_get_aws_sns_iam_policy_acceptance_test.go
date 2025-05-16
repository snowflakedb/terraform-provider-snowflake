//go:build !account_level_tests

package datasources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SystemGetAWSSNSIAMPolicy_basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: policyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_system_get_aws_sns_iam_policy.p", "aws_sns_topic_arn", "arn:aws:sns:us-east-1:1234567890123456:mytopic"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_aws_sns_iam_policy.p", "aws_sns_topic_policy_json"),
				),
			},
		},
	})
}

func policyConfig() string {
	s := `
	data snowflake_system_get_aws_sns_iam_policy p {
		aws_sns_topic_arn = "arn:aws:sns:us-east-1:1234567890123456:mytopic"
	}
	`
	return s
}
