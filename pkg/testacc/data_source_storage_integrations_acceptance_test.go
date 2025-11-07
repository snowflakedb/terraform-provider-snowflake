//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_StorageIntegrations_basic(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"name": config.StringVariable(id.Name()),
					"allowed_locations": config.SetVariable(
						config.StringVariable("gcs://foo/"),
						config.StringVariable("gcs://bar/"),
					),
					"blocked_locations": config.SetVariable(
						config.StringVariable("gcs://foo/"),
						config.StringVariable("gcs://bar/"),
					),
					"comment": config.StringVariable("some comment"),
				},
				ConfigDirectory: config.TestNameDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_storage_integrations.test", "storage_integrations.#"),
					containsStorageIntegration(id, true, "some comment"),
				),
			},
		},
	})
}

func containsStorageIntegration(id sdk.AccountObjectIdentifier, enabled bool, comment string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "snowflake_storage_integrations" {
				continue
			}
			iter, err := strconv.ParseInt(rs.Primary.Attributes["storage_integrations.#"], 10, 32)
			if err != nil {
				return err
			}

			for i := 0; i < int(iter); i++ {
				if rs.Primary.Attributes[fmt.Sprintf("storage_integrations.%d.name", i)] == id.Name() {
					actualEnabled, err := strconv.ParseBool(rs.Primary.Attributes[fmt.Sprintf("storage_integrations.%d.enabled", i)])
					if err != nil {
						return err
					}

					if actualEnabled != enabled {
						return fmt.Errorf("expected comment: %v, but got: %v", enabled, actualEnabled)
					}

					actualComment := rs.Primary.Attributes[fmt.Sprintf("storage_integrations.%d.comment", i)]
					if actualComment != comment {
						return fmt.Errorf("expected comment: %s, but got: %s", comment, actualComment)
					}

					return nil
				}
			}

			return fmt.Errorf("storage integration (%s) not found", id.FullyQualifiedName())
		}
		return nil
	}
}
