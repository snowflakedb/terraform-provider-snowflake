package resources

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_warehouseTypeCreateBuilder(t *testing.T) {
	warehouseTypeSchema := map[string]*schema.Schema{"warehouse_type": {Type: schema.TypeString, Optional: true}}

	tests := []struct {
		name   string
		config map[string]any
		want   sdk.WarehouseType
	}{
		{name: "not set", config: map[string]any{}, want: sdk.WarehouseTypeStandard},
		{name: "standard", config: map[string]any{"warehouse_type": "STANDARD"}, want: sdk.WarehouseTypeStandard},
		{name: "snowpark-optimized lowercase", config: map[string]any{"warehouse_type": "snowpark-optimized"}, want: sdk.WarehouseTypeSnowparkOptimized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, warehouseTypeSchema, tt.config)
			request := sdk.NewCreateWarehouseRequest(sdk.NewAccountObjectIdentifier("wh"))
			require.NoError(t, warehouseTypeCreateBuilder(d, request))
			require.NotNil(t, request.WarehouseType)
			assert.Equal(t, tt.want, *request.WarehouseType)
		})
	}

	t.Run("invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, warehouseTypeSchema, map[string]any{"warehouse_type": "ADAPTIVE"})
		require.Error(t, warehouseTypeCreateBuilder(d, sdk.NewCreateWarehouseRequest(sdk.NewAccountObjectIdentifier("wh"))))
	})
}

func Test_warehouseTypeUpdate(t *testing.T) {
	testResource := &schema.Resource{
		Schema: map[string]*schema.Schema{"warehouse_type": {Type: schema.TypeString, Optional: true}},
	}

	tests := []struct {
		name   string
		state  map[string]string
		config map[string]any
		want   *sdk.WarehouseType
	}{
		{name: "no change", state: map[string]string{"warehouse_type": "STANDARD"}, config: map[string]any{"warehouse_type": "STANDARD"}},
		{name: "not set before and now", state: map[string]string{}, config: map[string]any{}},
		{name: "removed, sets standard instead of unset", state: map[string]string{"warehouse_type": "SNOWPARK-OPTIMIZED"}, config: map[string]any{}, want: sdk.Pointer(sdk.WarehouseTypeStandard)},
		{name: "changed", state: map[string]string{"warehouse_type": "STANDARD"}, config: map[string]any{"warehouse_type": "SNOWPARK-OPTIMIZED"}, want: sdk.Pointer(sdk.WarehouseTypeSnowparkOptimized)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, err := testResource.Diff(context.Background(), &terraform.InstanceState{ID: "wh", Attributes: tt.state}, &terraform.ResourceConfig{Config: tt.config}, nil)
			require.NoError(t, err)
			d, err := schema.InternalMap(testResource.Schema).Data(&terraform.InstanceState{ID: "wh", Attributes: tt.state}, diff)
			require.NoError(t, err)

			set := sdk.NewWarehouseSetRequest()
			require.NoError(t, warehouseTypeUpdate(d, set, sdk.NewWarehouseUnsetRequest()))
			assert.Equal(t, tt.want, set.WarehouseType)
		})
	}
}

func Test_warehouseApplySetUnset_waitsForTheResize(t *testing.T) {
	t.Run("a size change waits for completion", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, warehouseSchema, map[string]any{"warehouse_size": "SMALL"})
		requests, err := warehouseApplySetUnset(d)
		require.NoError(t, err)
		require.NotNil(t, requests.set.WarehouseSize)
		require.NotNil(t, requests.set.WaitForCompletion)
		require.True(t, *requests.set.WaitForCompletion)
	})

	t.Run("without a size change it does not", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, warehouseSchema, map[string]any{"comment": "c"})
		requests, err := warehouseApplySetUnset(d)
		require.NoError(t, err)
		require.Nil(t, requests.set.WarehouseSize)
		require.Nil(t, requests.set.WaitForCompletion)
		require.NotNil(t, requests.set.Comment)
	})
}

func Test_warehouseApplySetUnset_removedTypeSetsStandard(t *testing.T) {
	state := &terraform.InstanceState{ID: "wh", Attributes: map[string]string{"name": "wh", "warehouse_type": "SNOWPARK-OPTIMIZED"}}
	diff, err := schema.InternalMap(warehouseSchema).Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]any{"name": "wh"}), nil, nil, true)
	require.NoError(t, err)
	d, err := schema.InternalMap(warehouseSchema).Data(state, diff)
	require.NoError(t, err)

	requests, err := warehouseApplySetUnset(d)
	require.NoError(t, err)
	// UNSET WAREHOUSE_TYPE would reset the type to the Snowflake default, so the generated update sets STANDARD.
	require.NotNil(t, requests.set.WarehouseType)
	assert.Equal(t, sdk.WarehouseTypeStandard, *requests.set.WarehouseType)
	assert.Nil(t, requests.unset.WarehouseType)
}
