package planchecks

import (
	"context"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"
)

const setTestAddress = "snowflake_share.t"

func setTestResourceChange(actions tfjson.Actions, before, after, afterUnknown any) *tfjson.ResourceChange {
	return &tfjson.ResourceChange{
		Address: setTestAddress,
		Change:  &tfjson.Change{Actions: actions, Before: before, After: after, AfterUnknown: afterUnknown},
	}
}

func setTestPlan(change *tfjson.ResourceChange, drift *tfjson.ResourceChange) *tfjson.Plan {
	plan := &tfjson.Plan{}
	if change != nil {
		plan.ResourceChanges = []*tfjson.ResourceChange{change}
	}
	if drift != nil {
		plan.ResourceDrift = []*tfjson.ResourceChange{drift}
	}
	return plan
}

func setTestAttributes(accounts any) map[string]any {
	return map[string]any{"name": "SHARE", "accounts": accounts}
}

func runSetPlanCheck(check plancheck.PlanCheck, plan *tfjson.Plan) error {
	resp := &plancheck.CheckPlanResponse{}
	check.CheckPlan(context.Background(), plancheck.CheckPlanRequest{Plan: plan}, resp)
	return resp.Error
}

func TestExpectSetChange(t *testing.T) {
	update := tfjson.Actions{tfjson.ActionUpdate}
	plan := func(before, after any) *tfjson.Plan {
		return setTestPlan(setTestResourceChange(update, setTestAttributes(before), setTestAttributes(after), nil), nil)
	}

	t.Run("order and duplicates do not matter", func(t *testing.T) {
		err := runSetPlanCheck(
			ExpectSetChange(setTestAddress, "accounts", tfjson.ActionUpdate, []string{"A", "B"}, []string{"C"}),
			plan([]any{"B", "A", "A"}, []any{"C"}),
		)
		require.NoError(t, err)
	})
	t.Run("null and empty are the same empty set", func(t *testing.T) {
		require.NoError(t, runSetPlanCheck(ExpectSetChange(setTestAddress, "accounts", tfjson.ActionUpdate, nil, []string{"A"}), plan(nil, []any{"A"})))
		require.NoError(t, runSetPlanCheck(ExpectSetChange(setTestAddress, "accounts", tfjson.ActionUpdate, []string{}, []string{"A"}), plan([]any{}, []any{"A"})))
		require.NoError(t, runSetPlanCheck(ExpectSetChange(setTestAddress, "accounts", tfjson.ActionUpdate, []string{"A"}, nil), plan([]any{"A"}, []any{})))
	})
	t.Run("values are compared exactly, including the case", func(t *testing.T) {
		err := runSetPlanCheck(
			ExpectSetChange(setTestAddress, "accounts", tfjson.ActionUpdate, []string{"a"}, []string{"b"}),
			plan([]any{"A"}, []any{"b"}),
		)
		require.ErrorContains(t, err, "before expected to be [a], got [A]")
	})
	t.Run("different values", func(t *testing.T) {
		err := runSetPlanCheck(
			ExpectSetChange(setTestAddress, "accounts", tfjson.ActionUpdate, []string{"A"}, []string{"B"}),
			plan([]any{"A"}, []any{"B", "C"}),
		)
		require.ErrorContains(t, err, "after expected to be [B], got [B C]")
	})
	t.Run("wrong action", func(t *testing.T) {
		err := runSetPlanCheck(
			ExpectSetChange(setTestAddress, "accounts", tfjson.ActionCreate, nil, []string{"A"}),
			plan(nil, []any{"A"}),
		)
		require.ErrorContains(t, err, "expected action create")
	})
	t.Run("resource not found", func(t *testing.T) {
		err := runSetPlanCheck(ExpectSetChange("snowflake_share.other", "accounts", tfjson.ActionUpdate, nil, nil), plan(nil, nil))
		require.ErrorContains(t, err, "no resource change found for snowflake_share.other")
	})
	t.Run("unknown value", func(t *testing.T) {
		change := setTestResourceChange(update, setTestAttributes([]any{"A"}), setTestAttributes(nil), map[string]any{"accounts": true})
		err := runSetPlanCheck(ExpectSetChange(setTestAddress, "accounts", tfjson.ActionUpdate, []string{"A"}, nil), setTestPlan(change, nil))
		require.ErrorContains(t, err, "unknown at plan time")
	})
	t.Run("not a collection", func(t *testing.T) {
		err := runSetPlanCheck(ExpectSetChange(setTestAddress, "name", tfjson.ActionUpdate, nil, nil), plan([]any{}, []any{}))
		require.ErrorContains(t, err, "expected a collection")
	})
	t.Run("collection of objects", func(t *testing.T) {
		err := runSetPlanCheck(ExpectSetChange(setTestAddress, "accounts", tfjson.ActionUpdate, nil, nil), plan([]any{}, []any{map[string]any{"a": "b"}}))
		require.ErrorContains(t, err, "collection of scalar values")
	})
	t.Run("nested path", func(t *testing.T) {
		nested := func(values ...any) map[string]any {
			return map[string]any{"block": []any{map[string]any{"items": values}}}
		}
		change := setTestResourceChange(update, nested("A"), nested("A", "B"), nil)
		err := runSetPlanCheck(ExpectSetChange(setTestAddress, "block.0.items", tfjson.ActionUpdate, []string{"A"}, []string{"A", "B"}), setTestPlan(change, nil))
		require.NoError(t, err)
	})
}

func TestExpectSetDrift(t *testing.T) {
	drift := setTestResourceChange(tfjson.Actions{tfjson.ActionUpdate}, setTestAttributes([]any{"A", "B"}), setTestAttributes([]any{"A"}), nil)

	t.Run("matches the drift and ignores the planned changes", func(t *testing.T) {
		plan := setTestPlan(setTestResourceChange(tfjson.Actions{tfjson.ActionUpdate}, setTestAttributes([]any{"A"}), setTestAttributes([]any{"A", "B"}), nil), drift)
		require.NoError(t, runSetPlanCheck(ExpectSetDrift(setTestAddress, "accounts", []string{"B", "A"}, []string{"A"}), plan))
	})
	t.Run("different values", func(t *testing.T) {
		err := runSetPlanCheck(ExpectSetDrift(setTestAddress, "accounts", []string{"A"}, []string{"A"}), setTestPlan(nil, drift))
		require.ErrorContains(t, err, "before expected to be [A], got [A B]")
	})
	t.Run("no drift", func(t *testing.T) {
		err := runSetPlanCheck(ExpectSetDrift(setTestAddress, "accounts", nil, nil), setTestPlan(nil, nil))
		require.ErrorContains(t, err, "no resource drift found for "+setTestAddress)
	})
}

func TestExpectSetElementsChange(t *testing.T) {
	update := tfjson.Actions{tfjson.ActionUpdate}
	plan := func(before, after any) *tfjson.Plan {
		return setTestPlan(setTestResourceChange(update, setTestAttributes(before), setTestAttributes(after), nil), nil)
	}

	t.Run("added and removed", func(t *testing.T) {
		err := runSetPlanCheck(
			ExpectSetElementsChange(setTestAddress, "accounts", tfjson.ActionUpdate, []string{"C"}, []string{"A"}),
			plan([]any{"A", "B"}, []any{"C", "B"}),
		)
		require.NoError(t, err)
	})
	t.Run("only added", func(t *testing.T) {
		err := runSetPlanCheck(
			ExpectSetElementsChange(setTestAddress, "accounts", tfjson.ActionUpdate, []string{"B"}, nil),
			plan([]any{"A"}, []any{"A", "B"}),
		)
		require.NoError(t, err)
	})
	t.Run("unchanged set has nothing added or removed", func(t *testing.T) {
		require.NoError(t, runSetPlanCheck(ExpectSetElementsChange(setTestAddress, "accounts", tfjson.ActionUpdate, nil, nil), plan([]any{"A", "B"}, []any{"B", "A"})))
	})
	t.Run("different elements", func(t *testing.T) {
		err := runSetPlanCheck(
			ExpectSetElementsChange(setTestAddress, "accounts", tfjson.ActionUpdate, []string{"C"}, []string{"A"}),
			plan([]any{"A", "B"}, []any{"B", "D"}),
		)
		require.ErrorContains(t, err, "added elements expected to be [C], got [D]")
	})
	t.Run("a new resource adds everything", func(t *testing.T) {
		change := setTestResourceChange(tfjson.Actions{tfjson.ActionCreate}, nil, setTestAttributes([]any{"A"}), nil)
		err := runSetPlanCheck(ExpectSetElementsChange(setTestAddress, "accounts", tfjson.ActionCreate, []string{"A"}, nil), setTestPlan(change, nil))
		require.NoError(t, err)
	})
}

func TestExpectSetElementsDrift(t *testing.T) {
	drift := setTestResourceChange(tfjson.Actions{tfjson.ActionUpdate}, setTestAttributes([]any{"A", "B"}), setTestAttributes([]any{"A", "C"}), nil)

	t.Run("matches", func(t *testing.T) {
		require.NoError(t, runSetPlanCheck(ExpectSetElementsDrift(setTestAddress, "accounts", []string{"C"}, []string{"B"}), setTestPlan(nil, drift)))
	})
	t.Run("different elements", func(t *testing.T) {
		err := runSetPlanCheck(ExpectSetElementsDrift(setTestAddress, "accounts", []string{"C"}, nil), setTestPlan(nil, drift))
		require.ErrorContains(t, err, "removed elements expected to be [], got [B]")
	})
}
