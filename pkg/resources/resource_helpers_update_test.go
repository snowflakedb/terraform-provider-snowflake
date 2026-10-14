package resources

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_attributeMappedValueUpdateSetOnlyAlways(t *testing.T) {
	sizeSchema := map[string]*schema.Schema{
		"size": {Type: schema.TypeString, Optional: true},
	}
	d := resourceDataWithPlan(t, sizeSchema, map[string]string{
		"size": "SMALL",
	}, map[string]any{
		"size": "",
	})

	var got *string
	require.NoError(t, attributeMappedValueUpdateSetOnlyAlways(d, "size", &got, func(v string) (string, error) {
		return v, nil
	}))
	require.NotNil(t, got)
	assert.Equal(t, "", *got)
}

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

func Test_renameAccountLevelResource(t *testing.T) {
	oldId := sdk.NewAccountObjectIdentifier("old")

	t.Run("leaves the identifier unchanged", func(t *testing.T) {
		d := resourceDataWithPlan(t, map[string]*schema.Schema{
			"name": {Type: schema.TypeString, Required: true},
		}, map[string]string{
			"name": "old",
		}, map[string]any{
			"name": "old",
		})

		called := false
		got, err := renameAccountLevelResource(context.Background(), d, oldId, func(context.Context, sdk.AccountObjectIdentifier, sdk.AccountObjectIdentifier) error {
			called = true
			return nil
		})
		require.NoError(t, err)
		assert.False(t, called)
		assert.Equal(t, oldId, got)
		assert.Equal(t, "id", d.Id())
	})

	t.Run("renames an account object", func(t *testing.T) {
		d := resourceDataWithPlan(t, map[string]*schema.Schema{
			"name": {Type: schema.TypeString, Required: true},
		}, map[string]string{
			"name": "old",
		}, map[string]any{
			"name": "new",
		})

		var gotOld, gotNew sdk.AccountObjectIdentifier
		got, err := renameAccountLevelResource(context.Background(), d, oldId, func(_ context.Context, id, newId sdk.AccountObjectIdentifier) error {
			gotOld, gotNew = id, newId
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, oldId, gotOld)
		assert.Equal(t, sdk.NewAccountObjectIdentifier("new"), gotNew)
		assert.Equal(t, gotNew, got)
		assert.Equal(t, helpers.EncodeResourceIdentifier(gotNew), d.Id())
	})

	t.Run("returns the rename error without updating the id", func(t *testing.T) {
		d := resourceDataWithPlan(t, map[string]*schema.Schema{
			"name": {Type: schema.TypeString, Required: true},
		}, map[string]string{
			"name": "old",
		}, map[string]any{
			"name": "new",
		})

		renameErr := errors.New("rename failed")
		got, err := renameAccountLevelResource(context.Background(), d, oldId, func(context.Context, sdk.AccountObjectIdentifier, sdk.AccountObjectIdentifier) error {
			return renameErr
		})
		require.ErrorIs(t, err, renameErr)
		assert.Equal(t, oldId, got)
		assert.Equal(t, "id", d.Id())
	})
}
