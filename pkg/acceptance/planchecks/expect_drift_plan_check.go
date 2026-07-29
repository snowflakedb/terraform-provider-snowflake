package planchecks

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var _ plancheck.PlanCheck = expectDriftPlanCheck{}

type expectDriftPlanCheck struct {
	resourceAddress string
	attribute       string
	oldValue        *string
	newValue        *string
}

// TODO [SNOW-1473409]: test
// TODO [SNOW-1473409]: extract common logic with expectChangePlanCheck
// TODO [SNOW-1473409]: extract traversal for the attribute path
// TODO [SNOW-1473409]: verify that path to attribute results in nil or primitive
// TODO [SNOW-1473409]: check if the nested attributes also have plan
func (e expectDriftPlanCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var result []error
	var resourceFound bool

	for _, change := range req.Plan.ResourceDrift {
		if e.resourceAddress != change.Address {
			continue
		}
		resourceFound = true

		var before, after map[string]any
		if change.Change.Before != nil {
			before = change.Change.Before.(map[string]any)
		}
		if change.Change.After != nil {
			after = change.Change.After.(map[string]any)
		}

		attributePathParts := strings.Split(e.attribute, ".")
		attributeRoot := attributePathParts[0]
		valueBefore, valueBeforeOk := before[attributeRoot]
		valueAfter, valueAfterOk := after[attributeRoot]

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

		if e.oldValue == nil && !(valueBefore == nil || valueBefore == "") {
			result = append(result, fmt.Errorf("expect drift: attribute %s before expected to be empty, got: %s", e.attribute, valueBefore))
		}
		if e.newValue == nil && !(valueAfter == nil || valueAfter == "") {
			result = append(result, fmt.Errorf("expect drift: attribute %s after expected to be empty, got: %s", e.attribute, valueAfter))
		}

		if e.oldValue != nil {
			if !valueBeforeOk {
				result = append(result, fmt.Errorf("expect drift: attribute %s before expected to be %s, got empty", e.attribute, *e.oldValue))
			} else if *e.oldValue != fmt.Sprintf("%v", valueBefore) {
				result = append(result, fmt.Errorf("expect drift: attribute %s before expected to be %s, got %v", e.attribute, *e.oldValue, valueBefore))
			}
		}
		if e.newValue != nil {
			if !valueAfterOk {
				result = append(result, fmt.Errorf("expect drift: attribute %s after expected to be %s, got empty", e.attribute, *e.newValue))
			} else if *e.newValue != fmt.Sprintf("%v", valueAfter) {
				result = append(result, fmt.Errorf("expect drift: attribute %s after expected to be %s, got %v", e.attribute, *e.newValue, valueAfter))
			}
		}

		if !slices.Contains(change.Change.Actions, tfjson.ActionUpdate) {
			result = append(result, fmt.Errorf("expect drift: expected action %s for %s, got: %v", tfjson.ActionUpdate, e.resourceAddress, change.Change.Actions))
		}
	}

	if !resourceFound {
		result = append(result, fmt.Errorf("expect drift: no resource drift found for %s", e.resourceAddress))
	}

	resp.Error = errors.Join(result...)
}

// TODO [SNOW-1473409]: describe
func ExpectDrift(resourceAddress string, attribute string, oldValue *string, newValue *string) plancheck.PlanCheck {
	return expectDriftPlanCheck{
		resourceAddress,
		attribute,
		oldValue,
		newValue,
	}
}
