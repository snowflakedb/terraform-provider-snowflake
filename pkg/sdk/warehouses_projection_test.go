package sdk

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWarehouseShowProjections(t *testing.T) {
	size := WarehouseSizeSmall
	minCluster := 1
	tables := []SchemaObjectIdentifier{NewSchemaObjectIdentifier("DB", "SCH", "T")}
	level := MaxQueryPerformanceLevelMedium
	w := &Warehouse{
		Name:                     "WH",
		State:                    WarehouseStateStarted,
		Type:                     WarehouseTypeStandard,
		Size:                     &size,
		MinClusterCount:          &minCluster,
		AutoResume:               true,
		CreatedOn:                time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
		ResourceMonitor:          NewAccountObjectIdentifier("MON"),
		OwnerRoleType:            "ROLE",
		MaxQueryPerformanceLevel: &level,
		Tables:                   tables,
		Actives:                  "should-not-project",
	}

	regular := w.AsRegular()
	require.Equal(t, w.Name, regular.Name)
	require.Equal(t, w.Size, regular.Size)
	require.Equal(t, w.MinClusterCount, regular.MinClusterCount)
	require.Equal(t, w.ResourceMonitor, regular.ResourceMonitor)
	require.Nil(t, (*Warehouse)(nil).AsRegular())

	adaptive := w.AsAdaptive()
	require.Equal(t, w.Name, adaptive.Name)
	require.Equal(t, w.MaxQueryPerformanceLevel, adaptive.MaxQueryPerformanceLevel)
	require.Nil(t, (*Warehouse)(nil).AsAdaptive())

	interactive := w.AsInteractive()
	require.Equal(t, w.Name, interactive.Name)
	require.Equal(t, w.Tables, interactive.Tables)
	require.Equal(t, w.Size, interactive.Size)
	require.Nil(t, (*Warehouse)(nil).AsInteractive())
}
