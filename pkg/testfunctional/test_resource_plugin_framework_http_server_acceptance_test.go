package testfunctional_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var currentResponse *string

func init() {
	allTestHandlers["http_server_example"] = &httpServerExampleHandler{}
}

type httpServerExampleHandler struct{}

func (h *httpServerExampleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	d, err := w.Write([]byte(*currentResponse))
	functionalTestLog.Printf("[DEBUG] Bytes written: %d, err: %v", d, err)
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
					currentResponse = sdk.Pointer("aaa")
				},
				Config: httpServerExampleConfig(id, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "id", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceReference, "response", "aaa"),
				),
			},
			{
				PreConfig: func() {
					currentResponse = sdk.Pointer("bbb")
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
