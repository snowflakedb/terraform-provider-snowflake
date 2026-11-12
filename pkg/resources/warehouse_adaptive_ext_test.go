package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_warehouseAdaptiveApplySetUnset_setsTheTypeWhenItChanges(t *testing.T) {
	t.Run("a changed warehouse type is always ADAPTIVE", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, warehouseAdaptiveSchema, map[string]any{"warehouse_type": "ADAPTIVE", "comment": "c"})
		requests, err := warehouseAdaptiveApplySetUnset(d)
		require.NoError(t, err)
		require.NotNil(t, requests.set.WarehouseType)
		require.Equal(t, sdk.WarehouseTypeAdaptive, *requests.set.WarehouseType)
		require.NotNil(t, requests.set.Comment)
	})

	t.Run("no change leaves the type out", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, warehouseAdaptiveSchema, map[string]any{"comment": "c"})
		requests, err := warehouseAdaptiveApplySetUnset(d)
		require.NoError(t, err)
		require.Nil(t, requests.set.WarehouseType)
	})
}
