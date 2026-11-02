package resources

import (
	"fmt"
	"maps"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_postgresMirrorApplySetUnsetExt(t *testing.T) {
	const mirror = "my_mirror"
	noRequest := func() *sdk.AlterPostgresMirrorArgumentsRequest {
		return sdk.NewAlterPostgresMirrorArgumentsRequest(mirror)
	}
	applySetUnset := func(t *testing.T, state map[string]string, config map[string]any) *postgresMirrorAlterRequestsExt {
		t.Helper()
		d := resourceDataWithPlan(t, postgresMirrorSchema, state, config)
		d.SetId(mirror)
		reqs, err := postgresMirrorApplySetUnsetExt(d)
		require.NoError(t, err)
		return reqs
	}
	interval := func(minutes string) map[string]string {
		return map[string]string{
			"refresh_interval.#":         "1",
			"refresh_interval.0.hours":   "0",
			"refresh_interval.0.minutes": minutes,
			"refresh_interval.0.seconds": "0",
		}
	}
	withState := func(attributes ...map[string]string) map[string]string {
		merged := map[string]string{}
		for _, a := range attributes {
			maps.Copy(merged, a)
		}
		return merged
	}

	t.Run("nothing changed sends nothing", func(t *testing.T) {
		reqs := applySetUnset(t, withState(map[string]string{"warehouse": "WH"}, interval("10")), map[string]any{"warehouse": "WH"})
		assert.Equal(t, noRequest(), reqs.set)
		assert.Equal(t, noRequest(), reqs.unset)
	})

	t.Run("removing the warehouse is an unset only", func(t *testing.T) {
		reqs := applySetUnset(t, map[string]string{"warehouse": "WH"}, map[string]any{})
		assert.Equal(t, noRequest(), reqs.set)
		assert.Equal(t, noRequest().WithUnsetWarehouse(true), reqs.unset)
	})

	t.Run("changing the warehouse is a set only", func(t *testing.T) {
		reqs := applySetUnset(t, map[string]string{"warehouse": "WH"}, map[string]any{"warehouse": "OTHER_WH"})
		assert.Equal(t, noRequest().WithWarehouse("OTHER_WH"), reqs.set)
		assert.Equal(t, noRequest(), reqs.unset)
	})

	t.Run("a set change and an unset are kept apart", func(t *testing.T) {
		reqs := applySetUnset(t, withState(map[string]string{"warehouse": "WH"}, interval("15")), map[string]any{
			"refresh_interval": []any{map[string]any{"minutes": 1}},
		})
		assert.Equal(t, noRequest().WithRefreshInterval("1 minute"), reqs.set)
		assert.Equal(t, noRequest().WithUnsetWarehouse(true), reqs.unset)
	})

	t.Run("removing the refresh interval block keeps the current interval", func(t *testing.T) {
		reqs := applySetUnset(t, interval("15"), map[string]any{})
		assert.Equal(t, noRequest(), reqs.set)
		assert.Equal(t, noRequest(), reqs.unset)
	})

	t.Run("tables and schemas are sent as additions and removals", func(t *testing.T) {
		state := &terraform.InstanceState{ID: mirror, Attributes: map[string]string{
			"postgres_tables.#": "2",
			fmt.Sprintf("postgres_tables.%d", schema.HashString("public.a")): "public.a",
			fmt.Sprintf("postgres_tables.%d", schema.HashString("public.b")): "public.b",
		}}
		diff := &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{
			"postgres_tables.#": {Old: "2", New: "2"},
			fmt.Sprintf("postgres_tables.%d", schema.HashString("public.a")): {Old: "public.a", NewRemoved: true},
			fmt.Sprintf("postgres_tables.%d", schema.HashString("public.b")): {Old: "public.b", New: "public.b"},
			fmt.Sprintf("postgres_tables.%d", schema.HashString("public.c")): {Old: "", New: "public.c"},
		}}
		d, err := schema.InternalMap(postgresMirrorSchema).Data(state, diff)
		require.NoError(t, err)

		reqs, err := postgresMirrorApplySetUnsetExt(d)
		require.NoError(t, err)
		assert.Equal(t, noRequest().
			WithAddTables([]sdk.StringListItemWrapper{{Value: "public.c"}}).
			WithRemoveTables([]sdk.StringListItemWrapper{{Value: "public.a"}}), reqs.set)
		assert.Equal(t, noRequest(), reqs.unset)
	})
}
