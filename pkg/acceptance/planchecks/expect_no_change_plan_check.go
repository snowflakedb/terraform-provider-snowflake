package planchecks

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var _ plancheck.PlanCheck = expectNoChangePlanCheck{}

type expectNoChangePlanCheck struct {
	resourceAddress string
	attribute       string
}

// TODO [SNOW-1473409]: test
func (e expectNoChangePlanCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var result []error

	for _, change := range req.Plan.ResourceChanges {
		if e.resourceAddress != change.Address {
			continue
		}

		var before, after map[string]any
		if change.Change.Before != nil {
			before = change.Change.Before.(map[string]any)
		}
		if change.Change.After != nil {
			after = change.Change.After.(map[string]any)
		}

		attributePathParts := strings.Split(e.attribute, ".")
		attributeRoot := attributePathParts[0]
		valueBefore := before[attributeRoot]
		valueAfter := after[attributeRoot]

		for idx, part := range attributePathParts {
			part := part
			if idx == 0 {
				continue
			}
			partInt, err := strconv.Atoi(part)
			if valueBefore != nil {
				if err != nil {
					valueBefore = valueBefore.(map[string]any)[part]
				} else {
					valueBefore = valueBefore.([]any)[partInt]
				}
			}
			if valueAfter != nil {
				if err != nil {
					valueAfter = valueAfter.(map[string]any)[part]
				} else {
					valueAfter = valueAfter.([]any)[partInt]
				}
			}
		}

		if valueBefore != valueAfter {
			result = append(result, fmt.Errorf("expect no change: attribute %s before=%s, after=%s", e.attribute, valueBefore, valueAfter))
		}
	}

	resp.Error = errors.Join(result...)
}

// TODO [SNOW-1473409]: describe
func ExpectNoChangeOnField(resourceAddress string, attribute string) plancheck.PlanCheck {
	return expectNoChangePlanCheck{
		resourceAddress,
		attribute,
	}
}
