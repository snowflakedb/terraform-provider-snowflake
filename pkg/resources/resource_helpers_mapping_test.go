package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resourceDataWithPlan returns ResourceData for an update plan.
// schema.TestResourceDataRaw does not populate GetChange, and d.Set does not either.
// The state ID is set because an empty ID is a create and marks every attribute changed.
// TypeSet diffs produced here include an extra blank element.
func resourceDataWithPlan(t *testing.T, s map[string]*schema.Schema, stateAttrs map[string]string, config map[string]any) *schema.ResourceData {
	t.Helper()
	res := &schema.Resource{Schema: s}
	state := &terraform.InstanceState{ID: "id", Attributes: stateAttrs}
	diff, err := res.Diff(context.Background(), state, &terraform.ResourceConfig{Config: config}, nil)
	require.NoError(t, err)
	d, err := schema.InternalMap(s).Data(state, diff)
	require.NoError(t, err)
	return d
}

func Test_setMappedValueCreateBuilder(t *testing.T) {
	mapItem := func(v any) (string, error) {
		return v.(string), nil
	}

	d := schema.TestResourceDataRaw(t, stringSetSchema(), map[string]any{
		"items": []any{"b", "a"},
	})

	var got []string
	require.NoError(t, setMappedValueCreateBuilder(d, "items", func(v []string) any {
		got = v
		return nil
	}, mapItem))
	assert.ElementsMatch(t, []string{"a", "b"}, got)

	absent := schema.TestResourceDataRaw(t, stringSetSchema(), map[string]any{})
	var untouched []string
	require.NoError(t, setMappedValueCreateBuilder(absent, "items", func(v []string) any {
		untouched = v
		return nil
	}, mapItem))
	assert.Nil(t, untouched)
}

func Test_listObjects(t *testing.T) {
	got, err := listObjects([]any{
		map[string]any{"name": "a"},
		map[string]any{"name": "b"},
	}, func(block map[string]any) (string, error) {
		return block["name"].(string), nil
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, got)

	_, err = listObjects("nope", func(map[string]any) (string, error) { return "", nil })
	require.Error(t, err)
}

func Test_setValueUpdate(t *testing.T) {
	mapItem := func(v any) (string, error) {
		return v.(string), nil
	}

	t.Run("replaces the set", func(t *testing.T) {
		d := resourceDataWithPlan(t, stringSetSchema(), map[string]string{
			"items.#": "1",
			fmt.Sprintf("items.%d", schema.HashString("old")): "old",
		}, map[string]any{
			"items": []any{"new"},
		})

		var got []string
		var unset *bool
		require.NoError(t, setValueUpdate(d, "items", &got, &unset, mapItem))
		assert.Contains(t, got, "new")
		assert.Nil(t, unset)
	})

	t.Run("unsets when the set is cleared", func(t *testing.T) {
		// resourceDataWithPlan leaves a blank set element, so build a removal diff directly.
		state := &terraform.InstanceState{ID: "id", Attributes: map[string]string{
			"items.#": "1",
			fmt.Sprintf("items.%d", schema.HashString("old")): "old",
		}}
		diff := &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{
			"items.#": {Old: "1", New: "0"},
			fmt.Sprintf("items.%d", schema.HashString("old")): {Old: "old", New: "", NewRemoved: true},
		}}
		d, err := schema.InternalMap(stringSetSchema()).Data(state, diff)
		require.NoError(t, err)

		var got []string
		var unset *bool
		require.NoError(t, setValueUpdate(d, "items", &got, &unset, mapItem))
		assert.Empty(t, got)
		require.NotNil(t, unset)
		assert.True(t, *unset)
	})

	t.Run("leaves the fields untouched when unchanged", func(t *testing.T) {
		d := resourceDataWithPlan(t, stringSetSchema(), map[string]string{
			"items.#": "1",
			fmt.Sprintf("items.%d", schema.HashString("same")): "same",
		}, map[string]any{
			"items": []any{"same"},
		})

		var got []string
		var unset *bool
		require.NoError(t, setValueUpdate(d, "items", &got, &unset, mapItem))
		assert.Nil(t, got)
		assert.Nil(t, unset)
	})
}
