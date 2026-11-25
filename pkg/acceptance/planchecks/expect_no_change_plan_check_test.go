package planchecks

import (
	"context"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"
)

func Test_ExpectNoChangeOnField(t *testing.T) {
	check := func(attribute string, before, after map[string]any) error {
		req := plancheck.CheckPlanRequest{Plan: &tfjson.Plan{ResourceChanges: []*tfjson.ResourceChange{
			setTestResourceChange(tfjson.Actions{tfjson.ActionUpdate}, before, after, nil),
		}}}
		resp := &plancheck.CheckPlanResponse{}
		ExpectNoChangeOnField(setTestAddress, attribute).CheckPlan(context.Background(), req, resp)
		return resp.Error
	}

	t.Run("scalar: same value", func(t *testing.T) {
		require.NoError(t, check("comment", map[string]any{"comment": "a"}, map[string]any{"comment": "a"}))
	})

	t.Run("scalar: changed value", func(t *testing.T) {
		err := check("comment", map[string]any{"comment": "a"}, map[string]any{"comment": "b"})
		require.ErrorContains(t, err, "attribute comment before=a, after=b")
	})

	t.Run("set: same elements", func(t *testing.T) {
		require.NotPanics(t, func() {
			require.NoError(t, check("allowed_databases", map[string]any{"allowed_databases": []any{"a", "b"}}, map[string]any{"allowed_databases": []any{"a", "b"}}))
		})
	})

	t.Run("set: changed elements", func(t *testing.T) {
		require.NotPanics(t, func() {
			err := check("allowed_databases", map[string]any{"allowed_databases": []any{"a"}}, map[string]any{"allowed_databases": []any{"a", "b"}})
			require.ErrorContains(t, err, "attribute allowed_databases before=[a], after=[a b]")
		})
	})

	t.Run("block: same value", func(t *testing.T) {
		require.NotPanics(t, func() {
			block := func() map[string]any { return map[string]any{"show_output": []any{map[string]any{"name": "a"}}} }
			require.NoError(t, check("show_output", block(), block()))
		})
	})

	t.Run("nested path into a list", func(t *testing.T) {
		before := map[string]any{"show_output": []any{map[string]any{"name": "a"}}}
		after := map[string]any{"show_output": []any{map[string]any{"name": "b"}}}
		require.ErrorContains(t, check("show_output.0.name", before, after), "attribute show_output.0.name before=a, after=b")
	})

	t.Run("resource without before or after", func(t *testing.T) {
		require.NoError(t, check("comment", nil, nil))
	})
}
