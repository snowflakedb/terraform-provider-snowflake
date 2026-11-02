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

// The checks below compare a set (or a list, which is compared as a set) of scalar values in the plan. A set has no
// order in the plan, so the values are sorted and de-duplicated before they are compared. A missing, null, and empty
// collection are all the same empty set. The values are compared exactly, including the case.
//
// ExpectSetChange and ExpectSetDrift compare the whole set before and after.
// ExpectSetElementsChange and ExpectSetElementsDrift compare only the difference between before and after.
//
// TODO [SNOW-1473409]: extract the common traversal with the other plan checks.

var (
	_ plancheck.PlanCheck = expectSetPlanCheck{}
	_ plancheck.PlanCheck = expectSetElementsPlanCheck{}
)

type setPlanCheckSource string

const (
	setPlanCheckChange setPlanCheckSource = "change"
	setPlanCheckDrift  setPlanCheckSource = "drift"
)

// ExpectSetChange checks that the planned change of the resource has the given set of values in the attribute before and after,
// and the given action.
func ExpectSetChange(resourceAddress string, attribute string, action tfjson.Action, oldValue []string, newValue []string) plancheck.PlanCheck {
	return expectSetPlanCheck{setPlanCheckChange, resourceAddress, attribute, action, oldValue, newValue}
}

// ExpectSetDrift checks that the refresh found a drift of the attribute from the first set of values to the second one.
func ExpectSetDrift(resourceAddress string, attribute string, oldValue []string, newValue []string) plancheck.PlanCheck {
	return expectSetPlanCheck{setPlanCheckDrift, resourceAddress, attribute, tfjson.ActionUpdate, oldValue, newValue}
}

// ExpectSetElementsChange checks that the planned change of the resource adds and removes exactly the given values of the attribute.
func ExpectSetElementsChange(resourceAddress string, attribute string, action tfjson.Action, added []string, removed []string) plancheck.PlanCheck {
	return expectSetElementsPlanCheck{setPlanCheckChange, resourceAddress, attribute, action, added, removed}
}

// ExpectSetElementsDrift checks that the refresh found that exactly the given values were added to and removed from the attribute.
func ExpectSetElementsDrift(resourceAddress string, attribute string, added []string, removed []string) plancheck.PlanCheck {
	return expectSetElementsPlanCheck{setPlanCheckDrift, resourceAddress, attribute, tfjson.ActionUpdate, added, removed}
}

type expectSetPlanCheck struct {
	source          setPlanCheckSource
	resourceAddress string
	attribute       string
	action          tfjson.Action
	oldValue        []string
	newValue        []string
}

func (e expectSetPlanCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	name := "expect set " + string(e.source)
	resp.Error = checkSetChanges(name, req.Plan, e.source, e.resourceAddress, e.attribute, e.action, func(before, after []string) []error {
		var errs []error
		if want := normalizeSet(e.oldValue); !slices.Equal(before, want) {
			errs = append(errs, fmt.Errorf("%s: attribute %s before expected to be %v, got %v", name, e.attribute, want, before))
		}
		if want := normalizeSet(e.newValue); !slices.Equal(after, want) {
			errs = append(errs, fmt.Errorf("%s: attribute %s after expected to be %v, got %v", name, e.attribute, want, after))
		}
		return errs
	})
}

type expectSetElementsPlanCheck struct {
	source          setPlanCheckSource
	resourceAddress string
	attribute       string
	action          tfjson.Action
	added           []string
	removed         []string
}

func (e expectSetElementsPlanCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	name := "expect set elements " + string(e.source)
	resp.Error = checkSetChanges(name, req.Plan, e.source, e.resourceAddress, e.attribute, e.action, func(before, after []string) []error {
		var errs []error
		if want, got := normalizeSet(e.added), setDifference(after, before); !slices.Equal(got, want) {
			errs = append(errs, fmt.Errorf("%s: attribute %s added elements expected to be %v, got %v", name, e.attribute, want, got))
		}
		if want, got := normalizeSet(e.removed), setDifference(before, after); !slices.Equal(got, want) {
			errs = append(errs, fmt.Errorf("%s: attribute %s removed elements expected to be %v, got %v", name, e.attribute, want, got))
		}
		return errs
	})
}

// checkSetChanges finds the resource in the planned changes or in the drift, reads the attribute before and after as sets,
// and runs compare on them. It also checks the action.
func checkSetChanges(name string, plan *tfjson.Plan, source setPlanCheckSource, resourceAddress string, attribute string, action tfjson.Action, compare func(before, after []string) []error) error {
	changes := plan.ResourceChanges
	if source == setPlanCheckDrift {
		changes = plan.ResourceDrift
	}

	var errs []error
	var resourceFound bool
	for _, change := range changes {
		if change.Address != resourceAddress {
			continue
		}
		resourceFound = true

		before, err := setAtPath(change.Change.Before, nil, attribute)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: attribute %s before: %w", name, attribute, err))
			continue
		}
		after, err := setAtPath(change.Change.After, change.Change.AfterUnknown, attribute)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: attribute %s after: %w", name, attribute, err))
			continue
		}
		errs = append(errs, compare(before, after)...)

		if !slices.Contains(change.Change.Actions, action) {
			errs = append(errs, fmt.Errorf("%s: expected action %s for %s, got: %v", name, action, resourceAddress, change.Change.Actions))
		}
	}
	if !resourceFound {
		errs = append(errs, fmt.Errorf("%s: no resource %s found for %s", name, source, resourceAddress))
	}
	return errors.Join(errs...)
}

// setAtPath reads the attribute from the object of a resource change (the "before" or "after" value, or nil) and returns
// its values as a sorted set without duplicates. The attribute is a path of object keys and list indexes separated with
// dots, like in the other plan checks. unknown is the matching "after_unknown" object; an unknown attribute is an error.
func setAtPath(object any, unknown any, attribute string) ([]string, error) {
	var value any = object
	for part := range strings.SplitSeq(attribute, ".") {
		if unknown == true {
			return nil, fmt.Errorf("the value is unknown at plan time")
		}
		unknown = childAt(unknown, part)
		value = childAt(value, part)
	}
	if unknown == true {
		return nil, fmt.Errorf("the value is unknown at plan time")
	}

	switch typed := value.(type) {
	case nil:
		return []string{}, nil
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			switch item.(type) {
			case string, bool, float64:
				values = append(values, fmt.Sprintf("%v", item))
			default:
				return nil, fmt.Errorf("expected a collection of scalar values, got an element of type %T", item)
			}
		}
		return normalizeSet(values), nil
	default:
		return nil, fmt.Errorf("expected a collection, got %T", value)
	}
}

func childAt(parent any, part string) any {
	switch typed := parent.(type) {
	case map[string]any:
		return typed[part]
	case []any:
		index, err := strconv.Atoi(part)
		if err != nil || index < 0 || index >= len(typed) {
			return nil
		}
		return typed[index]
	default:
		return nil
	}
}

func normalizeSet(values []string) []string {
	result := slices.Clone(values)
	slices.Sort(result)
	result = slices.Compact(result)
	if result == nil {
		return []string{}
	}
	return result
}

// setDifference returns the sorted values of a that are not in b.
func setDifference(a []string, b []string) []string {
	result := []string{}
	for _, value := range a {
		if !slices.Contains(b, value) {
			result = append(result, value)
		}
	}
	return normalizeSet(result)
}
