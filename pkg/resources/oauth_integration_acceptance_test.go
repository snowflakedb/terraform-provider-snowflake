package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_OAuthIntegration(t *testing.T) {
	oauthIntName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	integrationType := "TABLEAU_DESKTOP"

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: oauthIntegrationConfig(oauthIntName, integrationType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "name", oauthIntName),
					resource.TestCheckResourceAttr("snowflake_oauth_integration.test", "oauth_client", integrationType),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration.test", "oauth_issue_refresh_tokens"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration.test", "oauth_refresh_token_validity"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration.test", "blocked_roles_list"),
					resource.TestCheckResourceAttrSet("snowflake_oauth_integration.test", "created_on"),
				),
			},
			{
				ResourceName:      "snowflake_oauth_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func oauthIntegrationConfig(name string, integrationType string) string {
	return fmt.Sprintf(`
	resource "snowflake_oauth_integration" "test" {
		name = "%s"
		oauth_client = "%s"
		enabled = true
  		oauth_issue_refresh_tokens = true
  		oauth_refresh_token_validity = 3600
  		blocked_roles_list = ["SYSADMIN"]
	}
	`, name, integrationType)
}
