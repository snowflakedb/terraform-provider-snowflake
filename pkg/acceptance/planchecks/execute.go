package planchecks

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// Execute is a function that can be used to execute any code at the given phase.
type Execute func()

func (e Execute) CheckPlan(_ context.Context, _ plancheck.CheckPlanRequest, _ *plancheck.CheckPlanResponse) {
	e()
}
