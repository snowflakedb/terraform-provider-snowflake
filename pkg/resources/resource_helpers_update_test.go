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

func Test_removeThenAdd(t *testing.T) {
	t.Run("removes before adding", func(t *testing.T) {
		var order []string
		err := removeThenAdd([]string{"old"}, []string{"new"},
			func(item string) error {
				order = append(order, "remove:"+item)
				return nil
			},
			func(item string) error {
				order = append(order, "add:"+item)
				return nil
			},
		)
		require.NoError(t, err)
		assert.Equal(t, []string{"remove:old", "add:new"}, order)
	})

	t.Run("stops when a removal fails", func(t *testing.T) {
		var added []string
		err := removeThenAdd([]string{"old"}, []string{"new"},
			func(string) error { return errors.New("remove failed") },
			func(item string) error {
				added = append(added, item)
				return nil
			},
		)
		require.EqualError(t, err, "remove failed")
		assert.Empty(t, added)
	})
}

func Test_removeThenAddWithPlaceholder(t *testing.T) {
	t.Run("removes the placeholder after a failed update", func(t *testing.T) {
		var order []string
		err := removeThenAddWithPlaceholder("temp", []string{"old"}, []string{"new"},
			func(item string) error {
				order = append(order, "remove:"+item)
				if item == "old" {
					return errors.New("remove failed")
				}
				return nil
			},
			func(item string) error {
				order = append(order, "add:"+item)
				return nil
			},
		)
		require.EqualError(t, err, "remove failed")
		assert.Equal(t, []string{"add:temp", "remove:old", "remove:temp"}, order)
	})

	t.Run("does not remove a placeholder that failed to add", func(t *testing.T) {
		var removed []string
		err := removeThenAddWithPlaceholder("temp", []string{"old"}, []string{"new"},
			func(item string) error {
				removed = append(removed, item)
				return nil
			},
			func(item string) error {
				if item == "temp" {
					return errors.New("add failed")
				}
				return nil
			},
		)
		require.EqualError(t, err, "add failed")
		assert.Empty(t, removed)
	})

	t.Run("joins the update error with a placeholder removal error", func(t *testing.T) {
		err := removeThenAddWithPlaceholder("temp", []string{"old"}, nil,
			func(item string) error { return errors.New("remove " + item) },
			func(string) error { return nil },
		)
		require.ErrorContains(t, err, "remove old")
		require.ErrorContains(t, err, "remove temp")
	})
}

func Test_stringAttributeUpdateBuilder(t *testing.T) {
	commentSchema := map[string]*schema.Schema{
		"comment": {Type: schema.TypeString, Optional: true},
	}

	t.Run("sets the new value", func(t *testing.T) {
		d := resourceDataWithPlan(t, commentSchema, map[string]string{"comment": "old"}, map[string]any{"comment": "new"})
		var set *string
		var unset bool
		require.NoError(t, stringAttributeUpdateBuilder(d, "comment", func(v string) any {
			set = &v
			return nil
		}, func(bool) any {
			unset = true
			return nil
		}))
		require.NotNil(t, set)
		assert.Equal(t, "new", *set)
		assert.False(t, unset)
	})

	t.Run("unsets a cleared value", func(t *testing.T) {
		d := resourceDataWithPlan(t, commentSchema, map[string]string{"comment": "old"}, map[string]any{})
		var set *string
		var unset bool
		require.NoError(t, stringAttributeUpdateBuilder(d, "comment", func(v string) any {
			set = &v
			return nil
		}, func(bool) any {
			unset = true
			return nil
		}))
		assert.Nil(t, set)
		assert.True(t, unset)
	})

	t.Run("leaves an unchanged value alone", func(t *testing.T) {
		d := resourceDataWithPlan(t, commentSchema, map[string]string{"comment": "same"}, map[string]any{"comment": "same"})
		called := false
		require.NoError(t, stringAttributeUpdateBuilder(d, "comment", func(string) any {
			called = true
			return nil
		}, func(bool) any {
			called = true
			return nil
		}))
		assert.False(t, called)
	})
}

func Test_booleanAttributeUpdateSetOnlyBuilder(t *testing.T) {
	flagSchema := map[string]*schema.Schema{
		"change_tracking": {Type: schema.TypeBool, Optional: true},
	}
	d := resourceDataWithPlan(t, flagSchema, map[string]string{"change_tracking": "false"}, map[string]any{"change_tracking": true})
	var set *bool
	require.NoError(t, booleanAttributeUpdateSetOnlyBuilder(d, "change_tracking", func(v bool) any {
		set = &v
		return nil
	}))
	require.NotNil(t, set)
	assert.True(t, *set)
}

func Test_intAttributeUpdateSetOnlyBuilder(t *testing.T) {
	nodesSchema := map[string]*schema.Schema{
		"min_nodes": {Type: schema.TypeInt, Required: true},
	}
	d := resourceDataWithPlan(t, nodesSchema, map[string]string{"min_nodes": "1"}, map[string]any{"min_nodes": 2})
	var set *int
	require.NoError(t, intAttributeUpdateSetOnlyBuilder(d, "min_nodes", func(v int) any {
		set = &v
		return nil
	}))
	require.NotNil(t, set)
	assert.Equal(t, 2, *set)
}

func Test_intAttributeWithSpecialDefaultUpdateBuilder(t *testing.T) {
	secsSchema := map[string]*schema.Schema{
		"auto_suspend_secs": {Type: schema.TypeInt, Optional: true, Default: IntDefault},
	}
	d := resourceDataWithPlan(t, secsSchema, map[string]string{"auto_suspend_secs": "60"}, map[string]any{"auto_suspend_secs": IntDefault})
	var set *int
	var unset bool
	require.NoError(t, intAttributeWithSpecialDefaultUpdateBuilder(d, "auto_suspend_secs", func(v int) any {
		set = &v
		return nil
	}, func(bool) any {
		unset = true
		return nil
	}))
	assert.Nil(t, set)
	assert.True(t, unset)
}

func Test_booleanStringAttributeUpdateBuilder(t *testing.T) {
	flagSchema := map[string]*schema.Schema{
		"auto_resume": {Type: schema.TypeString, Optional: true, Default: BooleanDefault},
	}
	d := resourceDataWithPlan(t, flagSchema, map[string]string{"auto_resume": BooleanDefault}, map[string]any{"auto_resume": BooleanTrue})
	var set *bool
	require.NoError(t, booleanStringAttributeUpdateBuilder(d, "auto_resume", func(v bool) any {
		set = &v
		return nil
	}, func(bool) any {
		return nil
	}))
	require.NotNil(t, set)
	assert.True(t, *set)
}

func Test_listValueUpdateBuilder(t *testing.T) {
	listSchema := map[string]*schema.Schema{
		"items": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
	mapItem := func(v any) (string, error) {
		return v.(string), nil
	}

	t.Run("sets the mapped list", func(t *testing.T) {
		d := resourceDataWithPlan(t, listSchema, map[string]string{
			"items.#": "1",
			"items.0": "old",
		}, map[string]any{
			"items": []any{"new"},
		})
		var set []string
		var unset bool
		require.NoError(t, listValueUpdateBuilder(d, "items", func(v []string) any {
			set = v
			return nil
		}, func(bool) any {
			unset = true
			return nil
		}, mapItem))
		assert.Equal(t, []string{"new"}, set)
		assert.False(t, unset)
	})

	t.Run("unsets an empty list", func(t *testing.T) {
		d := resourceDataWithPlan(t, listSchema, map[string]string{
			"items.#": "1",
			"items.0": "old",
		}, map[string]any{
			"items": []any{},
		})
		var set []string
		var unset bool
		require.NoError(t, listValueUpdateBuilder(d, "items", func(v []string) any {
			set = v
			return nil
		}, func(bool) any {
			unset = true
			return nil
		}, mapItem))
		assert.Nil(t, set)
		assert.True(t, unset)
	})
}
