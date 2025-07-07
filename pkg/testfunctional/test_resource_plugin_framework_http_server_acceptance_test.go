package testfunctional_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/httpserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var (
	httpServerExampleHandler = common.NewDynamicHandler[httpserver.Read]()
)

func init() {
	allTestHandlers["http_server_example"] = httpServerExampleHandler
}

func TestAcc_TerraformPluginFrameworkFunctional_HttpServer(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("abc")
	resourceType := fmt.Sprintf("%s_http_server", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerForPluginFrameworkFunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					httpServerExampleHandler.SetCurrentValue(httpserver.Read{Msg: "aaa"})
				},
				Config: httpServerExampleConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "response", "aaa"),
				),
			},
			{
				PreConfig: func() {
					httpServerExampleHandler.SetCurrentValue(httpserver.Read{Msg: "bbb"})
				},
				Config: httpServerExampleConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "response", "bbb"),
				),
			},
		},
	})
}

func httpServerExampleConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"

  name = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
