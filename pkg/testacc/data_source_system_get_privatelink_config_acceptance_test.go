//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SystemGetPrivateLinkConfig_aws(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: privateLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Common fields — always present
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "account_name"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "account_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "ocsp_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "regionless_account_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "regionless_snowsight_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "snowsight_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "regionless_privatelink_ocsp_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "privatelink_dashed_urls_for_duo"),
					// AWS-specific fields — present on AWS
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "aws_vpce_id"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "privatelink_account_principal"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "app_service_privatelink_url"),
					// Azure-specific fields — absent on AWS
					resource.TestCheckResourceAttr("data.snowflake_system_get_privatelink_config.p", "azure_pls_id", ""),
					resource.TestCheckResourceAttr("data.snowflake_system_get_privatelink_config.p", "internal_stage", ""),
					resource.TestCheckResourceAttr("data.snowflake_system_get_privatelink_config.p", "privatelink_snowflake_managed_storage_volume_nfs", ""),
					resource.TestCheckResourceAttr("data.snowflake_system_get_privatelink_config.p", "privatelink_snowflake_managed_storage_volume_fs", ""),
					// GCP-specific fields — absent on AWS
					resource.TestCheckResourceAttr("data.snowflake_system_get_privatelink_config.p", "privatelink_gcp_service_attachment", ""),
					// Client redirect fields — absent when client redirect is not configured
					resource.TestCheckResourceAttr("data.snowflake_system_get_privatelink_config.p", "privatelink_connection_urls", ""),
					resource.TestCheckResourceAttr("data.snowflake_system_get_privatelink_config.p", "privatelink_connection_ocsp_urls", ""),
				),
			},
		},
	})
}

func privateLinkConfig() string {
	s := `
	data snowflake_system_get_privatelink_config p {}
	`
	return s
}
