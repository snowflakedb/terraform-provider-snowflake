package testfunctional_test

// TestAcc_TerraformPluginFrameworkFunctional_StaleErrorRepro reproduces the
// "Saved plan is stale" error that caused intermittent failures in production
// acceptance tests (e.g. TestAcc_SecretWithBasicAuthentication_BasicUseCase).
//
// Root cause (terraform-plugin-testing v1.14.0):
//   When a destroy step has a non-nil Check function, the framework inserts an
//   explicit Refresh() call between CreatePlan(-destroy) and Apply:
//
//     1. CreatePlan(-destroy) — plan file embeds the current state snapshot (serial N, value V1)
//     2. Refresh()            — re-reads all resources; if any value changed, writes a new
//                               state file (serial N+1, value V2)
//     3. Apply(plan file)     — Terraform checks plan's embedded serial N against the
//                               current serial N+1 → "Saved plan is stale"
//
//   In production, V1 == V2 most of the time (Snowflake returns byte-identical
//   responses), so the failure was intermittent (~50% on slow machines).
//
// Why this test is consistent:
//   The HTTP server handler (randomIntHandler) generates a fresh random int64 on
//   every GET request. Because Read() always returns a different value, step 2
//   always writes a new serial, making the plan stale on every single run.
//
// Fix (terraform-plugin-testing v1.15.0):
//   The framework no longer calls Refresh() in the destroy+Check path.
//   In v1.15.0 the Refresh() is removed → the test passes.

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"testing"

	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// randomIntHandler returns a brand-new random int64 on every GET request.
// This guarantees that the state after Refresh() differs from the state embedded
// in the destroy plan, making "Saved plan is stale" 100% reproducible.
type randomIntHandler struct{}

func (h *randomIntHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(testfunctional.StaleReproRead{RandomInt: rand.Int64()}) // #nosec G404
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
	}
}

func init() {
	allTestHandlers["stale_error_repro"] = &randomIntHandler{}
}

func TestAcc_TerraformPluginFrameworkFunctional_StaleErrorRepro(t *testing.T) {
	id := sdk.NewAccountObjectIdentifier("stale-repro")
	resourceType := fmt.Sprintf("%s_stale_error_repro", PluginFrameworkFunctionalTestsProviderName)
	resourceReference := fmt.Sprintf("%s.test", resourceType)

	tfresource.Test(t, tfresource.TestCase{
		ProtoV6ProviderFactories: providerForPluginFrameworkFunctionalTestsFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []tfresource.TestStep{
			{
				Config: staleErrorReproConfig(id, resourceType),
				Check: tfresource.ComposeTestCheckFunc(
					tfresource.TestCheckResourceAttr(resourceReference, "name", id.Name()),
				),
			},
			// Destroy step with a non-nil Check.
			// In terraform-plugin-testing v1.14.0 this triggers Refresh() between
			// CreatePlan and Apply.  Because random_int changes on every Read, the
			// plan is always stale → test fails with "Saved plan is stale".
			// In v1.15.0 the Refresh() is removed → test passes.
			{
				Destroy: true,
				Config:  staleErrorReproConfig(id, resourceType),
				// The check is intentionally a no-op: its only purpose is to be
				// non-nil, which is the condition terraform-plugin-testing v1.14.0
				// uses to decide whether to insert a Refresh() between plan and apply.
				// A no-op also ensures the test passes on v1.15.0 once upgraded.
				Check: tfresource.ComposeTestCheckFunc(),
			},
		},
	})
}

func staleErrorReproConfig(id sdk.AccountObjectIdentifier, resourceType string) string {
	return fmt.Sprintf(`
resource "%[2]s" "test" {
  provider = "%[3]s"
  name     = "%[1]s"
}
`, id.Name(), resourceType, PluginFrameworkFunctionalTestsProviderName)
}
