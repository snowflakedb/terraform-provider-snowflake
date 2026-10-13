package resources

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_setValueUpdateFallback(t *testing.T) {
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
		require.NoError(t, setValueUpdateFallback(d, "items", &got, mapItem, []string{"fallback"}))
		assert.Contains(t, got, "new")
		assert.NotContains(t, got, "fallback")
	})

	t.Run("uses the fallback when the set is cleared", func(t *testing.T) {
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
		require.NoError(t, setValueUpdateFallback(d, "items", &got, mapItem, []string{"fallback"}))
		assert.Equal(t, []string{"fallback"}, got)
	})

	t.Run("leaves the field untouched when unchanged", func(t *testing.T) {
		d := resourceDataWithPlan(t, stringSetSchema(), map[string]string{
			"items.#": "1",
			fmt.Sprintf("items.%d", schema.HashString("same")): "same",
		}, map[string]any{
			"items": []any{"same"},
		})

		var got []string
		require.NoError(t, setValueUpdateFallback(d, "items", &got, mapItem, []string{"fallback"}))
		assert.Nil(t, got)
	})
}

func fillExampleBlock(d *schema.ResourceData, prefix string, value *exampleBlock) error {
	var unset *bool
	return stringAttributeUpdate(d, prefix+"name", &value.Name, &unset)
}

func Test_singleBlockAttributeUpdate(t *testing.T) {
	t.Run("changed nested field", func(t *testing.T) {
		d := resourceDataWithPlan(t, exampleBlockSchema(), map[string]string{
			"block.#":      "1",
			"block.0.name": "old",
		}, map[string]any{
			"block": []any{map[string]any{"name": "new"}},
		})

		var got *exampleBlock
		var unset *bool
		require.NoError(t, singleBlockAttributeUpdate(d, "block", &got, &unset, fillExampleBlock))
		require.NotNil(t, got)
		require.NotNil(t, got.Name)
		assert.Equal(t, "new", *got.Name)
		assert.Nil(t, unset)
	})

	t.Run("removed block", func(t *testing.T) {
		d := resourceDataWithPlan(t, exampleBlockSchema(), map[string]string{
			"block.#":      "1",
			"block.0.name": "old",
		}, map[string]any{})

		var got *exampleBlock
		var unset *bool
		require.NoError(t, singleBlockAttributeUpdate(d, "block", &got, &unset, fillExampleBlock))
		assert.Nil(t, got)
		require.NotNil(t, unset)
		assert.True(t, *unset)
	})

	t.Run("unchanged block", func(t *testing.T) {
		d := resourceDataWithPlan(t, exampleBlockSchema(), map[string]string{
			"block.#":      "1",
			"block.0.name": "same",
		}, map[string]any{
			"block": []any{map[string]any{"name": "same"}},
		})

		var got *exampleBlock
		var unset *bool
		require.NoError(t, singleBlockAttributeUpdate(d, "block", &got, &unset, fillExampleBlock))
		assert.Nil(t, got)
		assert.Nil(t, unset)
	})
}
